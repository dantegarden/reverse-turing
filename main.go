package main

import (
	"context"
	dify "github.com/zhouyangtingwen/dify-sdk-go"
	"log"
	"strings"
)

const (
	APIKey  = "app-vYLMNL4icPjR8WqHMGft3g6B"
	APIHost = "http://10.0.1.12:19110"
)

var (
	Characters = []string{"刘备", "曹操", "孙权"}
	AINum      = "2"
)

func a() {
	//var tpl = strings.Replace(PromptTpl, "{{character}}", "曹操", -1)
	//tpl = strings.Replace(PromptTpl, "{{ai_num}}", "2", -1)
	//tpl = strings.Replace(PromptTpl, "{{other_characters}}", "刘备、孙权", -1)

	// 创建请求体
	var (
		ctx    = context.Background()
		config = &dify.ClientConfig{
			Host:         APIHost,
			ApiSecretKey: APIKey,
		}
		c = dify.NewClientWithConfig(config)

		req = &dify.ChatMessageRequest{
			Inputs: map[string]interface{}{
				"character":        "曹操",
				"ai_num":           "2",
				"other_characters": "刘备、孙权",
			},
			ResponseMode: "streaming",
			Query:        "请你说一段开场白",
			User:         "my-user",
		}

		ch  chan dify.ChatMessageStreamChannelResponse
		err error
	)

	if ch, err = c.Api().ChatMessagesStream(ctx, req); err != nil {
		return
	}

	var strBuilder strings.Builder

	for {
		select {
		case <-ctx.Done():
			return
		case streamData, isOpen := <-ch:
			if err = streamData.Err; err != nil {
				log.Println(err.Error())
				return
			}
			if !isOpen {
				log.Println(strBuilder.String())
				return
			}

			strBuilder.WriteString(streamData.Answer)
		}
	}
}

func main() {

}
