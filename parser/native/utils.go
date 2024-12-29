package native

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func Base64Safe(content string) string {
	urlSafe := strings.NewReplacer("+", "-", "/", "_", "=", "").Replace(content)
	if decode, err := base64.RawURLEncoding.DecodeString(urlSafe); err == nil {
		return string(decode)
	}
	return content
}

func ConvertToStrings(data map[string]any) map[string]string {
	stringMap := make(map[string]string)
	for key, value := range data {
		switch v := value.(type) {
		case string:
			stringMap[key] = v
		case float64:
			stringMap[key] = strconv.Itoa(int(v))
		default:
			stringMap[key] = fmt.Sprintf("%v", v)
		}
	}
	return stringMap
}

func TrimBlank(str string) string {
	return strings.TrimSpace(str)
}

func StringToUint16(content string) uint16 {
	intNum, _ := strconv.Atoi(content)
	return uint16(intNum)
}

func StringToInt64(content string) int64 {
	intNum, _ := strconv.Atoi(content)
	return int64(intNum)
}

func StringToUint32(content string) uint32 {
	intNum, _ := strconv.Atoi(content)
	return uint32(intNum)
}

func DecodeURIComponent(content string) string {
	result, _ := url.QueryUnescape(content)
	return result
}

func SplitKeyValueWithColon(content string) (string, string) {
	if !strings.Contains(content, ":") {
		return TrimBlank(content), "1"
	}
	arr := strings.Split(content, ":")
	return TrimBlank(arr[0]), TrimBlank(arr[1])
}
