package meta

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zhouyangtingwen/dify-sdk-go"
	"io"
	"net/http"
	"reverse-turing/app/internal/sse"
	"reverse-turing/common/errno"
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

func Publish(w http.ResponseWriter, d *sse.Event) error {
	err := sse.Encode(w, d)
	if err != nil {
		logx.Errorf("encode sse event error: %v", err)
		return err
	}
	w.(http.Flusher).Flush()

	return nil
}

func SendDifyStreamMessage(ctx context.Context, c *dify.Client, req *dify.ChatMessageRequest, w http.ResponseWriter, endMark string, isPublish bool) (recvContent string, err error) {
	ch, err := c.Api().ChatMessagesStream(ctx, req)
	if err != nil {
		return "", err
	}
	w.Header().Set("Connection", "keep-alive")

	var strBuilder strings.Builder
	for {
		select {
		case <-ctx.Done():
			return
		case streamData, isOpen := <-ch:
			streamErr := streamData.Err
			if streamErr != nil {
				if errors.Is(streamErr, io.EOF) {
					break
				} else {
					reason := streamErr.Error()
					logx.Errorf("stream reading failed, error: %s", reason)
					return "", errno.ServerErr
				}
			}

			if !isOpen {
				var err error
				event := &sse.Event{
					Data: []byte(endMark),
				}
				if isPublish {
					err = Publish(w, event)
				}
				return strBuilder.String(), err
			}

			strBuilder.WriteString(streamData.Answer)

			event := &sse.Event{
				Data: []byte(streamData.Answer),
			}

			if isPublish {
				err = Publish(w, event)
			}
		}
	}
}
