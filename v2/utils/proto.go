package utils

import (
	"io"

	"google.golang.org/protobuf/proto"

	"github.com/gin-gonic/gin"
)

func ParseProto(c *gin.Context, varPtr proto.Message) error {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	if err := proto.Unmarshal(data, varPtr); err != nil {
		return err
	}
	return nil
}

func SendProto(c *gin.Context, statusCode int, varPtr proto.Message) error {
	data, err := proto.Marshal(varPtr)
	if err != nil {
		return err
	}
	c.Data(statusCode, "application/x-protobuf", data)
	return nil
}
