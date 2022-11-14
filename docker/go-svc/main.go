package main

import (
	"encoding/json"
	"fmt"
	"go-geneic-svc/warden"
	"os"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const HeaderWardenPermissionsForward = "X-Warden-Permissions-Forward"
const WardenPermissionEchoKey = "warden.forwardpermissions"

func wardenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rawPermissions := c.Request().Header.Get(HeaderWardenPermissionsForward)
		fmt.Printf("Get Warden Forwarded Permission %v\n", rawPermissions)
		var permissions []string
		json.Unmarshal([]byte(rawPermissions), &permissions)
		c.Set(WardenPermissionEchoKey, warden.NewPrimitiveQuery(permissions))
		return next(c)
	}
}
func AuthorizeResource(c echo.Context, i interface{}) error {
	permissions := c.Get(WardenPermissionEchoKey).(*warden.PrimitiveQuery)
	if !permissions.Match(i.(map[string]interface{})) {
		return echo.ErrForbidden
	}
	return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(wardenMiddleware)
	e.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")
		resultData := findRecordById(id)
		if authzE := AuthorizeResource(c, resultData); authzE != nil {
			return authzE
		}

		if checkNilInterface(resultData) {
			return c.JSON(404, map[string]string{"message": "notfound"})
		}
		return c.JSON(200, resultData)
	})

	e.PUT("/:id", func(c echo.Context) error {
		id := c.Param("id")
		resultData := findRecordById(id)
		if authzE := AuthorizeResource(c, resultData); authzE != nil {
			return authzE
		}
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
		if authzE := AuthorizeResource(c, p); authzE != nil {
			return authzE
		}
		if err != nil {
			return err
		}
		return c.JSON(200, p)
	})

	e.GET("/", func(c echo.Context) error {
		filteredData := []map[string]interface{}{}
		permission := c.Get(WardenPermissionEchoKey).(*warden.PrimitiveQuery)
		for _, d := range readDb() {
			dd := d.(map[string]interface{})
			if permission.Match(dd) {
				filteredData = append(filteredData, dd)
			}
		}
		return c.JSON(200, map[string]interface{}{"data": filteredData})
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
