package apiserver

import (
	"fmt"
	"kvsql/common"
	"kvsql/postgresql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type KeyObj struct {
	Key string `json:"key" xml:"key" form:"key" query:"key"`
}

type ReqObj struct {
	KeyObj
	Value        string `json:"value" xml:"value" form:"value" query:"value"`
	TtlRemaining int64  `json:"ttl" xml:"ttl" form:"ttl" query:"ttl"`
}

func Test(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func intToTime(t int64) postgresql.TTL {
	retval := postgresql.TTL{}
	if t == 0 {
		return postgresql.TTL{}
	}
	retval.TTL = time.Duration(t) * time.Second
	return retval
}

func PutKey(c echo.Context) error {
	r := new(ReqObj)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	fmt.Println(r)
	ttl := intToTime(r.TtlRemaining)
	err := common.Client.PutKey(r.Key, r.Value, ttl)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "Request completed")
}

func GetKey(c echo.Context) error {
	r := new(KeyObj)
	if err := c.Bind(r); err != nil {
		return err
	}
	val, err := common.Client.GetKey(r.Key)
	if err != nil {
		switch err.Error() {
		case "no rows in result set":
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Key is not present"})
		default:
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	if val.TtlRemaining < 0 {
		err = common.Client.RemoveKey(r.Key)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Key is not present"})
	}
	return c.JSON(http.StatusOK, val)
}

func DeleteKey(c echo.Context) error {
	r := new(KeyObj)
	if err := c.Bind(r); err != nil {
		return err
	}
	err := common.Client.RemoveKey(r.Key)
	if err != nil {
		switch err.Error() {
		case "no rows in result set":
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Key is not present"})
		default:
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "key deleted", "key": r.Key})
}
