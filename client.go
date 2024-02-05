package archeryapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http/cookiejar"
	"net/url"

	"github.com/go-resty/resty/v2"
)

const (
	apiURL = "http://sql.archery.com"

	maxRetries = 3
)

func NewClient(username, password string, opts ...ClientOption) *Client {
	u, err := url.Parse(apiURL)

	if err != nil {
		panic(err)
	}

	c := &Client{
		baseUrl:    u,
		maxRetries: maxRetries,
		debugger:   false,
		username:   username,
		password:   password,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.initHttpClient()

	c.Instance = &InstanceClient{apiClient: c}
	c.Database = &DatabaseClient{apiClient: c}
	c.Workflow = &WorkflowClient{apiClient: c}

	return c
}

type Client struct {
	httpClient *resty.Client

	baseUrl    *url.URL
	maxRetries int
	debugger   bool

	username string
	password string

	Database DatabaseService
	Instance InstanceService
	Workflow WorkflowService
}

type ClientOption func(*Client)

func WithMaxRetries(maxRetries int) ClientOption {
	return func(c *Client) {
		c.maxRetries = maxRetries
	}
}

func WithBaseUrl(baseUrl string) ClientOption {
	return func(c *Client) {
		u, err := url.Parse(baseUrl)
		if err != nil {
			panic(err)
		}
		c.baseUrl = u
	}
}

func WithDebugger(debugger bool) ClientOption {
	return func(c *Client) {
		c.debugger = debugger
	}
}

func (c *Client) initHttpClient() {
	jar, _ := cookiejar.New(nil)
	c.httpClient = resty.New()
	c.httpClient.SetCookieJar(jar)
	c.httpClient.SetBaseURL(c.baseUrl.String())
	c.httpClient.SetDebug(c.debugger)
	c.httpClient.SetRetryCount(c.maxRetries)
	err := c.authenticate()
	if err != nil {
		panic(err)
	}
}

func (c *Client) authenticate() error {
	resp, err := c.httpClient.R().
		EnableTrace().
		SetFormData(map[string]string{
			"username": c.username,
			"password": c.password,
		}).
		SetHeader("Cookie", "csrftoken=OTGjsZYttS1f59SmhA7RNgwNa0g0I8mbwEl0qycZQPUTKhLpil7kEEJfSzL2J1Z3").
		SetHeader("X-Csrftoken", "OTGjsZYttS1f59SmhA7RNgwNa0g0I8mbwEl0qycZQPUTKhLpil7kEEJfSzL2J1Z3").
		Post("/authenticate/")
	if err != nil {
		log.Fatal(err)
	}
	// 设置响应cookie
	c.httpClient.SetCookies(resp.Cookies())
	//从cookie 中获取 csrftoken
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "csrftoken" {
			c.httpClient.SetHeader("X-Csrftoken", cookie.Value)
		}
	}

	var result Result
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		log.Fatal(err)
	}
	if result.Status != 0 {
		return errors.New(result.Msg)
	}
	return nil
}

type Result struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
}
