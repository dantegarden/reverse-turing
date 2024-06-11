package utils

import (
	"context"
	"google.golang.org/grpc/metadata"
	"reverse-turing/common/errno"
)

func GetDeviceInfo(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		deviceId := md.Get("device_id")[0]
		return deviceId, nil
	}
	return "", errno.DeviceIdNotFoundErr
}
