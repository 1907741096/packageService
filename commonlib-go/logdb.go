package lib

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const LOG_DB = "LOG_DB"

func LogDB() gin.HandlerFunc {
	return func(context *gin.Context) {
		start := time.Now()
		context.Next()
		end := time.Now()

		buf := new(bytes.Buffer)
		binaryBuffer(buf, NewLogHead())
		binaryBuffer(buf, NewLogBody(context, start, end))

		go func() {
			conn, err := net.Dial("udp", LogDBConfig.URL)
			if err != nil {
				Warning("LogDb connect error", err.Error(), LOG_DB)
				return
			}
			defer conn.Close()
			conn.Write(buf.Bytes())
		}()
	}
}

type logHead struct {
	Result   uint8  /* 结果，只在返回包中有效，非0 代表失败 */
	LogType  uint8  /* 0: 表示包体中为数据，LogDB 解析数据执行入库操作 1: 表示包体为 SQL 语句，LogDB 执行该SQL 语句 */
	Sequence uint32 /* 序列号，返回包会回带此字段 */
	EchoLen  uint16 /* 回带字段长度 */
}

func NewLogHead() *logHead {
	return &logHead{
		Result:   0,
		LogType:  0,
		Sequence: uint32(rand.Intn(65535) + 1),
		EchoLen:  0,
	}
}

type logBody struct {
	Source     string
	UserId     int32
	ApiName    string
	Method     string
	Param      string
	Status     int32
	ErrCode    string
	ErrMessage string
	Response   string
	UserIp     string
	ServerIp   string
	EndTime    string
	SpendTime  string
}

func NewLogBody(c *gin.Context, start, end time.Time) *logBody {
	requestBytes, ok := c.Get(gin.BodyBytesKey)
	requestBytesSlice, typeOk := requestBytes.([]byte)
	if !ok || !typeOk {
		requestBytesSlice = make([]byte, 0)
	}

	responseData, ok := c.Get(ResponseDataKey)
	responseDataMap, typeOk := responseData.(gin.H)
	if !ok || !typeOk {
		responseDataMap = gin.H{
			"code":    SystemInterError,
			"message": ErrSystemInter.Error(),
		}
	}
	responseBytesSlice, _ := json.Marshal(responseDataMap)
	if len(responseBytesSlice) > 4096 {
		responseBytesSlice = responseBytesSlice[:4096]
	}

	var userId int
	//bindForm, ok := c.Get(utils.BindFormKey)
	//former, typeOk := bindForm.(models.Former)
	//if !ok || !typeOk {
	//	userId = former.GetUserId()
	//}
	return &logBody{
		Source:     SystemSource,
		UserId:     int32(userId),
		ApiName:    c.Request.RequestURI,
		Method:     c.Request.Method,
		Param:      strings.Replace(string(requestBytesSlice), "\n", "", -1),
		Status:     http.StatusOK,
		ErrCode:    strconv.Itoa(responseDataMap["code"].(int)),
		ErrMessage: responseDataMap["message"].(string),
		Response:   string(responseBytesSlice),
		UserIp:     c.ClientIP(),
		ServerIp:   getLocalIp(),
		EndTime:    fmt.Sprintf("%d", end.Unix()),
		SpendTime:  floatToString(float32(end.UnixNano()-start.UnixNano()) / 1e9),
	}
}

func binaryBuffer(buffer *bytes.Buffer, v interface{}) {
	ele := reflect.ValueOf(v).Elem()
	for i := 0; i < ele.NumField(); i++ {
		if !ele.Field(i).CanSet() {
			continue
		}

		switch ele.Field(i).Kind() {
		case reflect.String:
			field := ele.Field(i).Interface().(string)
			byteField := []byte(field)
			binary.Write(buffer, binary.BigEndian, uint16(len(field)))
			binary.Write(buffer, binary.BigEndian, byteField)
		default:
			binary.Write(buffer, binary.BigEndian, ele.Field(i).Interface())
		}
	}
}

func getLocalIp() string {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		Warning("LogDB combine data warning", "not found local IP", LOG_DB)
		return ""
	}

	for _, address := range addresses {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	return ""
}

func floatToString(arg float32) string {
	return strconv.FormatFloat(float64(arg), 'f', 6, 64)
}
