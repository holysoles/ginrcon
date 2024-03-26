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
	defConn = false

	ErrMissingDefaultConnectionInfo = errors.New("at least one default connection parameter was specified, but others were missing")
	ErrInvalidDefaultConnection     = errors.New("error opening a connection to the specified default rcon server")
	ErrNoDefaultConnection          = errors.New("no connection details were specified and no valid default connection exists")
	ErrInvalidConnectionDetails     = errors.New("invalid connection details to the rcon server were provided")
	ErrInvalidResponseFromRcon      = errors.New("an invalid response was received from the rcon server")
	ErrUnableToConnectTcpRcon       = errors.New("unable to establish tcp connection to specified rcon server")
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
	go testDefault()
	initWeb()
}

func openDefault() (*rcon.Conn, error) {
	s, defS := os.LookupEnv("RCON_SERVER")
	s = "docker3.donut.lan"
	defS = true //TODO remove
	p, defP := os.LookupEnv("RCON_PORT")
	p = "25575"
	defP = true
	pwd, defPwd := os.LookupEnv("RCON_ADMIN_PASSWORD")
	pwd = "flyhigh"
	defPwd = true

	if !defS && !defP && !defPwd {
		//nothing specified
		fmt.Println("no default connection information specified")
		return nil, nil
	}
	defConn = true // even if we arent able to setup the connection, there was intent to
	if !defS || !defP || !defPwd {
		return nil, ErrMissingDefaultConnectionInfo
	}
	rconHostPort := net.JoinHostPort(s, p)
	info := openConnInfo{Server: rconHostPort, Password: pwd}
	con, err := openRcon(info)
	if err != nil {
		return nil, err
	}
	defConn = true
	return con, nil
}
func testDefault() {
	// errors here are non-fatal, users can still provide specific connection server info later on
	con, err := openDefault()
	if err != nil {
		fmt.Println(err)
		return
	}
	if con == nil {
		return
	}
	err = con.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("successfully tested default RCON connection")
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
			case ErrUnableToConnectTcpRcon:
				c.AbortWithError(http.StatusBadGateway, err)
			default: //assume we screwed up
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}
	} else if defConn {
		fmt.Println("using default server connection for incoming request")
		conn, err = openDefault()
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(http.StatusBadGateway, ErrInvalidDefaultConnection)
			return
		}
	} else {
		// no default connection and no connection info provided
		c.AbortWithError(http.StatusBadRequest, ErrNoDefaultConnection)
		return
	}
	defer conn.Close()

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
