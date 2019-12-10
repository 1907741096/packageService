package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"runtime/debug"
)

func Recovery() gin.HandlerFunc {
	raven.SetDSN(SentryConfig.URL)
	return RecoveryWithWriter(raven.DefaultClient, true)
}

func RecoveryWithWriter(client *raven.Client, onlyCrashes bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			flags := map[string]string{
				"environment": gin.Mode(),
				"source":      SystemSource,
				"endpoint":    c.Request.RequestURI,
			}

			if rvr := recover(); rvr != nil {
				if gin.Mode() == gin.DebugMode {
					debug.PrintStack()
					log.Printf("\n\n\x1b[34m %s \n\n", rvr)
				}
				rvrStr := fmt.Sprintf("\n\n\x1b[31m %s", rvr)
				packet := raven.NewPacketWithExtra(
					rvrStr,
					raven.Extra{"params": getRequestParams(c)},
					raven.NewException(errors.New(rvrStr), raven.NewStacktrace(3, 3, nil)),
					raven.NewHttp(c.Request))
				client.Capture(packet, flags)
				FormatFail(c, ErrSystemInter)
				c.Abort()
			}
			if !onlyCrashes {
				for _, item := range c.Errors {
					packet := raven.NewPacket(
						item.Error(),
						&raven.Message{
							Message: item.Error(),
							Params:  []interface{}{item.Meta},
						},
						raven.NewHttp(c.Request))
					client.Capture(packet, flags)
				}
			}
		}()

		c.Next()
	}
}

func getRequestParams(c *gin.Context) (params map[string]interface{}) {
	if requestBody, ok := c.Get(gin.BodyBytesKey); ok {
		if err := json.Unmarshal(requestBody.([]byte), &params); err != nil {
			log.Println("request body params: ", err.Error())
		}
	}
	return
}
