package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"bitbucket.org/rendon/logger"
	"github.com/revel/revel"
	"identities/app/models"
)

type App struct {
	*revel.Controller
	RequestBody []byte
}

const (
	maxLogLength = 512
)

var (
	log = logger.New("Cache")
)

func (c *App) readRequestBody() revel.Result {
	var err error
	if c.RequestBody, err = ioutil.ReadAll(c.Request.Body); err != nil {
		c.RequestBody = []byte("")
	}

	var method = c.Request.Method
	var path = c.Request.URL.Path
	var addr = c.Request.RemoteAddr
	if len(c.RequestBody) > 0 {
		var data map[string]interface{}
		var printed = false
		err = json.Unmarshal(c.RequestBody, &data)
		if err == nil {
			indented, err := json.MarshalIndent(data, "", "    ")
			if len(indented) > maxLogLength {
				indented = indented[0:maxLogLength]
			}
			if err == nil {
				log.Printf("[%s] %s %s\n%s\n", addr, method, path, indented)
				printed = true
			}
		}
		if !printed {
			log.Printf("[%s] %s %s\n%s\n", addr, method, path, c.RequestBody)
		}
	} else {
		log.Printf("[%s] %s %s\n", addr, method, path)
	}
	return nil
}

func init() {
	revel.InterceptMethod((*App).readRequestBody, revel.BEFORE)
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Data(i interface{}) revel.Result {
	c.Response.Status = http.StatusOK
	c.Response.ContentType = "application/json"

	var res = models.SuccessResponse{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    i,
	}
	return c.RenderJson(res)
}

func (c App) Ok() revel.Result {
	c.Response.Status = http.StatusOK
	c.Response.ContentType = "application/json"

	var res = models.SuccessResponse{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
	}
	log.Printf("OK")
	return c.RenderJson(res)
}

func (c App) Error(e error) revel.Result {
	c.Response.Status = http.StatusBadRequest
	c.Response.ContentType = "application/json"

	log.Errorf("%s", e)
	var res = models.ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: e.Error(),
	}

	return c.RenderJson(res)
}

func (c App) NotFound(e error) revel.Result {
	c.Response.Status = http.StatusNotFound
	c.Response.ContentType = "application/json"

	var res = models.ErrorResponse{
		Status:  http.StatusNotFound,
		Message: e.Error(),
	}

	log.Printf("NOT FOUND")
	return c.RenderJson(res)
}

func (c App) Cors() revel.Result {
	return c.Ok()
}
