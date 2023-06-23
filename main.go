package main

import (
	"fmt"

	"github.com/eugenshima/myapp/internal/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	fmt.Println("123")
	handlers.HttpHandler(e)

	e.Logger.Fatal(e.Start(":1323"))

}
