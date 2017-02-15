package main

import (
	"github.com/labstack/echo"
)

func AuthMiddleware(username, password string, c echo.Context) bool {
	return apiKey == username
}
