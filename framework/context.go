package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context
	handler        ControllerHandler
	// 是否超时标记位
	hasTimeout bool
	// 写保护机制
	writerMux *sync.Mutex // 对responseWriter的写入时保护，防止多个goroutine同时写入乱序
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		writerMux:      &sync.Mutex{},
	}
}

func (c *Context) WriterMux() *sync.Mutex {
	return c.writerMux
}

func (c *Context) GetRequest() *http.Request {
	return c.request
}

func (c *Context) GetResponse() http.ResponseWriter {
	return c.responseWriter
}

func (c *Context) SetHasTimeout() {
	c.hasTimeout = true
}

func (c *Context) HasTimeout() bool {
	return c.hasTimeout
}

func (c *Context) BaseContext() context.Context {
	return c.request.Context()
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.BaseContext().Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.BaseContext().Done()
}

func (c *Context) Err() error {
	return c.BaseContext().Err()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.BaseContext().Value(key)
}

func (c *Context) QueryArray(key string, def []string) []string {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

func (c *Context) QueryAll() map[string][]string {
	if c.request != nil {
		return c.request.URL.Query()
	}
	return map[string][]string{}
}

func (c *Context) QueryInt(key string, def int) int {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		// 返回数组中的最后一个
		if len > 0 {
			intval, err := strconv.Atoi(vals[len-1])
			if err != nil {
				return def
			}
			return intval
		}
	}
	return def
}

func (c *Context) QueryString(key string, def string) string {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			return vals[len-1]
		}
	}
	return def
}

func (c *Context) FormArray(key string, def []string) []string {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

func (c *Context) FormAll() map[string][]string {
	if c.request != nil {
		return c.request.PostForm
	}
	return map[string][]string{}
}

func (c *Context) FormInt(key string, def int) int {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			intval, err := strconv.Atoi(vals[len-1])
			if err != nil {
				return def
			}
			return intval
		}
	}
	return def
}

func (c *Context) FormString(key string, def string) string {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			return vals[len-1]
		}
	}
	return def
}

// BindJson 读出body数据
func (c *Context) BindJson(obj interface{}) error {
	if c.request != nil {
		body, err := ioutil.ReadAll(c.request.Body)
		if err != nil {
			return err
		}
		c.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ctx.request empty")
	}
	return nil
}

func (c *Context) Json(status int, obj interface{}) error {
	if c.HasTimeout() {
		return nil
	}
	c.responseWriter.Header().Set("Content-Type", "application/json")
	c.responseWriter.WriteHeader(status)
	byt, err := json.Marshal(obj)
	if err != nil {
		c.responseWriter.WriteHeader(500)
		return err
	}
	c.responseWriter.Write(byt)
	return nil
}

func (c *Context) HTML(status int, obj interface{}, template string) error {
	return nil
}

func (c *Context) Text(status int, obj string) error {
	return nil
}
