package builder

import (
	"context"
	"fmt"
	. "github.com/jiangew/belex"
	"github.com/jiangew/belex/fcoin"
	"github.com/jiangew/belex/huobi"
	"net"
	"net/http"
	"net/url"
	"time"
)

type APIBuilder struct {
	HttpClientConfig *HttpClientConfig
	client           *http.Client
	httpTimeout      time.Duration
	apiKey           string
	secretkey        string
	clientId         string
	apiPassphrase    string
}

type HttpClientConfig struct {
	HttpTimeout  time.Duration
	Proxy        *url.URL
	MaxIdleConns int
}

func (c HttpClientConfig) String() string {
	return fmt.Sprintf("{ProxyUrl:\"%s\",HttpTimeout:%s,MaxIdleConns:%d}", c.Proxy, c.HttpTimeout.String(), c.MaxIdleConns)
}

func (c *HttpClientConfig) SetHttpTimeout(timeout time.Duration) *HttpClientConfig {
	c.HttpTimeout = timeout
	return c
}

func (c *HttpClientConfig) SetProxyUrl(proxyUrl string) *HttpClientConfig {
	if proxyUrl == "" {
		return c
	}
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		return c
	}
	c.Proxy = proxy
	return c
}

func (c *HttpClientConfig) SetMaxIdleConns(max int) *HttpClientConfig {
	c.MaxIdleConns = max
	return c
}

var (
	DefaultHttpClientConfig = &HttpClientConfig{
		Proxy:        nil,
		HttpTimeout:  5 * time.Second,
		MaxIdleConns: 10}
	DefaultAPIBuilder = NewAPIBuilder()
)

func NewAPIBuilder() (builder *APIBuilder) {
	return NewAPIBuilder2(DefaultHttpClientConfig)
}

func NewAPIBuilder2(config *HttpClientConfig) *APIBuilder {
	if config == nil {
		config = DefaultHttpClientConfig
	}
	return &APIBuilder{
		HttpClientConfig: config,
		client: &http.Client{
			Timeout: config.HttpTimeout,
			Transport: &http.Transport{
				Proxy: func(request *http.Request) (*url.URL, error) {
					return config.Proxy, nil
				},
				MaxIdleConns:          config.MaxIdleConns,
				IdleConnTimeout:       5 * config.HttpTimeout,
				TLSHandshakeTimeout:   config.HttpTimeout,
				ResponseHeaderTimeout: config.HttpTimeout,
				ExpectContinueTimeout: config.HttpTimeout,
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.DialTimeout(network, addr, config.HttpTimeout)
				}},
		}}
}

func NewCustomAPIBuilder(client *http.Client) (builder *APIBuilder) {
	return &APIBuilder{client: client}
}

func (builder *APIBuilder) GetHttpClientConfig() *HttpClientConfig {
	return builder.HttpClientConfig
}

func (builder *APIBuilder) GetHttpClient() *http.Client {
	return builder.client
}

func (builder *APIBuilder) HttpProxy(proxyUrl string) (_builder *APIBuilder) {
	if proxyUrl == "" {
		return builder
	}
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		return builder
	}
	builder.HttpClientConfig.Proxy = proxy
	transport := builder.client.Transport.(*http.Transport)
	transport.Proxy = http.ProxyURL(proxy)
	return builder
}

func (builder *APIBuilder) HttpTimeout(timeout time.Duration) (_builder *APIBuilder) {
	builder.HttpClientConfig.HttpTimeout = timeout
	builder.httpTimeout = timeout
	builder.client.Timeout = timeout
	transport := builder.client.Transport.(*http.Transport)
	if transport != nil {
		//transport.ResponseHeaderTimeout = timeout
		//transport.TLSHandshakeTimeout = timeout
		transport.IdleConnTimeout = timeout
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		}
	}
	return builder
}

func (builder *APIBuilder) APIKey(key string) (_builder *APIBuilder) {
	builder.apiKey = key
	return builder
}

func (builder *APIBuilder) APISecretkey(key string) (_builder *APIBuilder) {
	builder.secretkey = key
	return builder
}

func (builder *APIBuilder) ClientID(id string) (_builder *APIBuilder) {
	builder.clientId = id
	return builder
}

func (builder *APIBuilder) ApiPassphrase(apiPassphrase string) (_builder *APIBuilder) {
	builder.apiPassphrase = apiPassphrase
	return builder
}

func (builder *APIBuilder) Build(exName string) (api API) {
	var _api API
	switch exName {
	case HUOBI_PRO:
		_api = huobi.NewHuoBiProSpot(builder.client, builder.apiKey, builder.secretkey)
	case FCOIN:
		_api = fcoin.NewFCoin(builder.client, builder.apiKey, builder.secretkey)
	default:
		println("exchange name error [" + exName + "].")

	}
	return _api
}

func (builder *APIBuilder) BuildFuture(exName string) (api FutureRestAPI) {
	switch exName {
	case HBDM:
		return huobi.NewHbdm(&APIConfig{HttpClient: builder.client, ApiKey: builder.apiKey, ApiSecretKey: builder.secretkey})
	default:
		println(fmt.Sprintf("%s not support future", exName))
		return nil
	}
}
