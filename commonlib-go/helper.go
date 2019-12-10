package lib

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/rs/xid"
	"hash/crc32"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// 获取执行的绝对路径
func GetAbsolutePath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index]
	return ret
}

/**
struct -> json
*/
func JsonMarshal(params interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		jsonError := map[string]interface{}{
			"params": params,
			"error":  err,
		}
		Error("json解析异常", jsonError, "json")
		return nil, ErrSystemInter
	}

	return jsonData, nil
}

/**
json 转 map
*/
func JsonUnMarshal(params []byte, obj interface{}) error {
	err := json.Unmarshal(params, &obj)
	if err != nil {
		jsonError := map[string]interface{}{
			"params": params,
			"error":  err,
		}
		Error("json解析异常", jsonError, "json")
		return ErrSystemInter
	}

	return nil
}

func IsEmpty(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func SetMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

/**
	如果time1在time2之前，返回true
*/
func ComparedTime(time1 string, time2 string) bool {
	t1, err := time.Parse(TimeFormat, time1)
	t2, err := time.Parse(TimeFormat, time2)
	if err == nil && t1.Before(t2) {
		return true
	}
	return false
}

/**
	带日期的全局唯一标识
*/
func GenerateUniqueCode() string {
	return time.Now().Format("20060102") + GenerateUUID()
}

func GenerateUUID() string {
	guid := xid.New()
	return guid.String()
}

func Crc32IEEE(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

/**
获取整形数组中的最大值
 */
func MaxIntArray(data []int)  int{
	maxVal := data[0]
	for i := 0; i < len(data); i++ {
		//从第二个元素开始循环比较，如果发现有更大的数，则交换
		if maxVal < data[i] {
			maxVal = data[i]
		}
	}
	return maxVal
}

