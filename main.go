package main

import (
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorcon/rcon"
)

type rconInfo struct {
	Server   string `form:"server" json:"server" xml:"server"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password"  binding:"required"`
	Command  string `form:"command" json:"command" xml:"command"  binding:"required"`
}
type rconReply struct {
	Message string `json:"message"`
}
type errorResponse struct {
	Error string
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.ForwardedByClientIP = true
	proxies := strings.Split(os.Getenv("TRUSTED_PROXIES"), ",")
	r.SetTrustedProxies(proxies)

	r.GET("/status", healthCheck)
	r.POST("/command", processCommand)

	r.Run(":" + os.Getenv("PORT"))
}

func processCommand(c *gin.Context) {
	var info rconInfo
	err := c.Bind(&info)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = info.validateConnectionInfo()
	if err != nil {
		c.JSON(http.StatusBadGateway, errorResponse{"Failed to test TCP connection to provided server"})
		return
	}

	conn, err := rcon.Dial(info.Server, info.Password)
	if err != nil {
		switch err {
		case rcon.ErrAuthNotRCON:
		case rcon.ErrInvalidAuthResponse:
			c.AbortWithStatus(http.StatusInternalServerError)
		case rcon.ErrAuthFailed:
			fallthrough
		default:
			c.AbortWithStatus(http.StatusUnauthorized)
		}
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
	}
	res := rconReply{msg}

	c.JSON(http.StatusOK, res)
}

func (i *rconInfo) validateConnectionInfo() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Var(i.Server, "required,hostname_port")
	if err != nil {
		return err
	}
	dialer := &net.Dialer{Timeout: time.Second * 1}
	_, err = dialer.Dial("tcp", i.Server)
	if err != nil {
		return err
	}
	return nil
}

func healthCheck(c *gin.Context) {
	c.Status(http.StatusOK)
}
