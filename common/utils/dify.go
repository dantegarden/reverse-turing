package utils

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zhouyangtingwen/dify-sdk-go"
	"strings"
)

func GetDifyConfig(ctx context.Context, apiHost string, apiKey string) *dify.Client {
	config := &dify.ClientConfig{
		Host:         apiHost,
		ApiSecretKey: apiKey,
	}
	c := dify.NewClientWithConfig(config)
	return c
}

func SendDifyMessage(ctx context.Context, c *dify.Client, req *dify.ChatMessageRequest, handelEveryStreamData func(d string)) (err error) {
	ch, err := c.Api().ChatMessagesStream(ctx, req)
	if err != nil {
		return err
	}

	var strBuilder strings.Builder

	for {
		select {
		case <-ctx.Done():
			return
		case streamData, isOpen := <-ch:
			if err = streamData.Err; err != nil {
				logx.Error(err.Error())
			}
			if !isOpen {
				logx.Info(strBuilder.String())
				return
			}

			strBuilder.WriteString(streamData.Answer)
			handelEveryStreamData(streamData.Answer)
		}
	}
}
