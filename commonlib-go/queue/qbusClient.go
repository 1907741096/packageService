package queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
)

// 请求路由v/t
type Router struct {
	Method, Action string
}

func getRouter() map[string]Router {
	var router = map[string]Router{
		"ReceiveMessage": {Method: "GET", Action: "/queue/messages/pull"},
		"DeleteMessage":  {Method: "POST", Action: "/queue/messages/ack"},
		"SendMessage":    {Method: "POST", Action: "/queue/messages/send"},
		"PublishMessage": {Method: "POST", Action: "/queue/messages/publish"},
	}
	return router
}

// qbus 配置
type ClientConfig struct {
	Host      string
	UserName  string
	UserToken string
	Version   string
	Timeout   int
}

//Queue Client
type Client struct {
	PathRouter map[string]Router
	config     ClientConfig
}

var client = &Client{}

func (q *Client) Url(uri string) string {
	return q.config.Host + uri
}

func (q *Client) GetHeader() lib.Header {
	return lib.Header{
		"username":   q.config.UserName,
		"token":      q.config.UserToken,
		"sdk":        q.config.Version,
		"Connection": "close",
	}
}

func NewClient(clientConfig ClientConfig) *Client {
	if client != nil {
		return client
	}

	client = &Client{
		PathRouter: getRouter(),
		config:     clientConfig,
	}
	return client
}

// 发送请求
func (q *Client) Request(router string, param interface{}, timeout int) (io.ReadCloser, error) {
	reqRouter := q.PathRouter[router]
	method := reqRouter.Method
	action := reqRouter.Action
	var err error
	var resp io.ReadCloser

	switch method {
	case lib.POST:
		resp, err = lib.PostJson(q, action, param, timeout)
	case lib.GET:
		resp, err = lib.Get(q, action, param, timeout)
	}

	if err != nil {
		return nil, err
	}
	return resp, nil
}

// 发送message
func (q *Client) sendMessage(name string, messages interface{}, delay int) lib.Response {
	params := map[string]interface{}{
		"queueName": name,
		"delayTime": delay,
		"data":      messages,
	}

	r, _ := q.Request("SendMessage", params, 0)
	rr, _ := ioutil.ReadAll(r)
	resp := lib.Response{}
	_ = json.Unmarshal(rr, &resp)

	return resp
}

// 接受message
func (q *Client) receiveMessage(name string, readNum uint8, pollingSecond int) (*QbusResponse, error) {
	params := map[string]string{
		"queueName":          name,
		"numOfMsg":           strconv.FormatUint(uint64(readNum), 10),
		"pollingWaitSeconds": strconv.FormatInt(int64(pollingSecond), 10),
	}
	resp, err := q.Request("ReceiveMessage", params, pollingSecond+10)
	if err != nil {
		return nil, err
	}
	r, _ := ioutil.ReadAll(resp)
	qr := QbusResponse{}
	if err := json.Unmarshal(r, &qr); err != nil {
		return nil, err
	}
	//qr.M2S(r)
	return &qr, nil
}

// delete message
func (q *Client) deleteMessage(name string, data []*QbusData) (*lib.Response, error) {
	if len(data) == 0 {
		return nil, errors.New("success message slice is empty")
	}
	var d []string
	for _, dd := range data {
		d = append(d, dd.ReceiptHandle)
	}

	params := map[string]interface{}{
		"queueName": name,
		"data":      d,
	}
	r, err := q.Request("DeleteMessage", params, 0)
	if err != nil {
		return nil, err
	}

	rr, _ := ioutil.ReadAll(r)
	resp := lib.Response{}
	_ = json.Unmarshal(rr, &resp)
	return &resp, nil
}

func setField(obj interface{}, name string, value interface{}) error {
	defer func() error {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				lib.Error("parse map exception", map[string]interface{}{
					"err":  err,
					"name": name,
				}, "parse")
				return err
			}
		}
		return nil
	}()
	s := reflect.ValueOf(obj)
	sv := reflect.Indirect(s).FieldByName(name)
	if !sv.IsValid() {
		return fmt.Errorf("no such field: %s in obj", name)
	}
	if !sv.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}
	st := sv.Type()
	val := reflect.ValueOf(value)
	if st != val.Type() {
		return errors.New("provided value type didn't match obj field type")
	}
	sv.Set(val)
	return nil
}

func (s *QbusResponse) M2S(params map[string]interface{}) error {
	for k, v := range params {
		err := setField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
