package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorcon/rcon"
)

var (
	gConn *rcon.Conn

	ErrNoDefaultConnection      = errors.New("no connection details were specified and no valid default connection exists")
	ErrInvalidConnectionDetails = errors.New("invalid connection details to the rcon server were provided")
	ErrInvalidResponseFromRcon  = errors.New("an invalid response was received from the rcon server")
)

type openConnInfo struct {
	Server   string `form:"server" json:"server" xml:"server"`
	Password string `form:"password" json:"password" xml:"password"`
}
type commandReq struct {
	openConnInfo
	Command string `form:"command" json:"command" xml:"command"  binding:"required"`
}
type commandRes struct {
	Message string `json:"message"`
}

func main() {
	initDefaultRcon()
	initWeb()
}

func initDefaultRcon() {
	var err error
	rconHostPort := net.JoinHostPort(os.Getenv("RCON_SERVER"), os.Getenv("RCON_PORT"))
	rconAdminPass := os.Getenv("RCON_ADMIN_PASSWORD")
	info := openConnInfo{Server: rconHostPort, Password: rconAdminPass}

	//check env vars to construct
	gConn, err = openRcon(info)
	//log error as a warning, if we have bad default info just throw it away
	if err == ErrInvalidConnectionDetails {
		fmt.Println("default RCON server connection details were not provided or invalid")
		return
	} else if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("successfully opened default RCON connection to", rconHostPort)
	}
}

func initWeb() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.ForwardedByClientIP = true
	proxies := strings.Split(os.Getenv("TRUSTED_PROXIES"), ",")
	r.SetTrustedProxies(proxies)

	r.GET("/status", healthCheck)
	r.POST("/command", processCommand)

	bindPort, customBind := os.LookupEnv("PORT")
	if customBind {
		r.Run(":" + bindPort)
	} else {
		r.Run()
	}
}

func processCommand(c *gin.Context) {
	var info commandReq
	err := c.Bind(&info)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// if they passed connection info, we should try to create a new connection
	var conn *rcon.Conn
	if info.openConnInfo.Server != "" {
		fmt.Println("using provided connection info for incoming request")
		conn, err = openRcon(info.openConnInfo)
		if err != nil {
			switch err {
			case ErrInvalidResponseFromRcon:
				c.AbortWithError(http.StatusInternalServerError, err)
			case ErrInvalidConnectionDetails:
				c.AbortWithError(http.StatusUnauthorized, err)
			default: //assume we screwed up
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}
		defer conn.Close()
	} else if gConn != nil {
		fmt.Println("using default server connection for incoming request")
		conn = gConn
	} else {
		// no valid default connection and no credentials provided
		c.AbortWithError(http.StatusBadGateway, ErrNoDefaultConnection)
	}
	msg, err := conn.Execute(info.Command)
	if err != nil {
		switch err {
		case rcon.ErrCommandTooLong:
			fallthrough
		case rcon.ErrCommandEmpty:
			c.AbortWithStatus(http.StatusBadRequest)
		default:
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	res := commandRes{msg}

	c.JSON(http.StatusOK, res)
}

func healthCheck(c *gin.Context) {
	c.Status(http.StatusOK)
}
