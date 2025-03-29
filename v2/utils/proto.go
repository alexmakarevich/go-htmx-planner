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
