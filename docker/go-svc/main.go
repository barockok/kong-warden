package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")
		resultData := findRecordById(id)
		if checkNilInterface(resultData) {
			return c.JSON(404, map[string]string{"message": "notfound"})
		}
		return c.JSON(200, resultData)
	})

	e.PUT("/:id", func(c echo.Context) error {
		id := c.Param("id")
		resultData := findRecordById(id)

		if checkNilInterface(resultData) {
			return c.JSON(404, map[string]string{"message": "notfound"})
		}

		p := make(map[string]interface{})
		err := c.Bind(&p)
		if err != nil {
			return err
		}
		return c.JSON(200, p)
	})

	e.POST("/", func(c echo.Context) error {
		p := make(map[string]interface{})
		err := c.Bind(&p)
		if err != nil {
			return err
		}
		return c.JSON(200, p)
	})

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{"data": readDb()})
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}

var db []interface{}

func findRecordById(id string) map[string]interface{} {
	var resultData map[string]interface{}
	for _, d := range readDb() {
		data := d.(map[string]interface{})
		if data["id"] == id {
			resultData = data
		}
	}
	return resultData
}
func readDb() []interface{} {
	if checkNilInterface(db) {
		fmt.Println("log")
		dat, err := os.ReadFile(os.Getenv("DB_FILE"))
		if err != nil {
			panic(err)
		}
		json.Unmarshal(dat, &db)
	}
	return db
}

func checkNilInterface(i interface{}) bool {
	iv := reflect.ValueOf(i)
	if !iv.IsValid() {
		return true
	}
	switch iv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		return iv.IsNil()
	default:
		return false
	}
}
