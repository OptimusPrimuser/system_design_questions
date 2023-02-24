package main

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

func userOnline(c echo.Context) error {
	usrName := c.Param("usrName")
	err := redisClient.OnlineUser(usrName)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, usrName+" is online now")
}

func getUsers(c echo.Context) error {
	users, err := redisClient.GetOnlineUsers()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	usrListByte, err := json.Marshal(users)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, string(usrListByte))
}
