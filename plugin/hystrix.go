package plugin

import (
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/v3/client"
)

type clientWrapper struct {
	client.Client
}

func NewHystrixClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &clientWrapper{c}
	}
}

// Call 熔断逻辑
func (c *clientWrapper) Call(ctx context.Context, req client.Request, resp interface{}, opts ...client.CallOption) error {
	return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
		// 正常逻辑
		fmt.Println("正常调用逻辑")
		return c.Client.Call(ctx, req, resp, opts...)
	}, func(err error) error {
		// 熔断逻辑， 每个服务都不一样
		fmt.Println("熔断处理")
		return err
	})
}
