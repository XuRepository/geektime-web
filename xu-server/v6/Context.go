package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{W: w, R: r}
}

func (c *Context) ReadJson(data interface{}) error {
	body, err := io.ReadAll(c.R.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, data)
}

func (c *Context) WriteJson(status int, data interface{}) error {
	c.W.WriteHeader(status)

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = c.W.Write(bytes)
	if err != nil {
		return err
	}

	return err

}

func (c *Context) OKJson(data interface{}) error {
	return c.WriteJson(http.StatusOK, data)
}

func (c *Context) SystemErrJson(data interface{}) error {
	// http 库里面提前定义好了各种响应码
	return c.WriteJson(http.StatusInternalServerError, data)
}

func (c *Context) BadRequestJson(data interface{}) error {
	// http 库里面提前定义好了各种响应码
	return c.WriteJson(http.StatusBadRequest, data)
}
