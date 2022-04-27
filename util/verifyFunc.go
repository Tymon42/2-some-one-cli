package util

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func AssertErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func AssertConv(ok bool) {
	if !ok {
		log.Panic("invalid widget conversion")
	}
}

func IsUrl(urls string) bool {

	re := regexp.MustCompile("(http|https):\\/\\/[\\w\\-_]+(\\.[\\w\\-_]+)+([\\w\\-\\.,@?^=%&:/~\\+#]*[\\w\\-\\@?^=%&/~\\+#])?")
	result := re.FindAllStringSubmatch(urls, -1)
	if result == nil {
		log.Println("URL不合法")
		return false
	}
	return true

}

func ConvPath(data string) (pathc string) {
	re3, _ := regexp.Compile("/")
	rep := re3.ReplaceAllStringFunc(data, strings.ToUpper)
	fmt.Println(rep)
	rep2 := re3.ReplaceAllString(data, "\\\\")
	fmt.Println(rep2)

	return rep2
}

func ResolveTime(seconds int) (minute int, sec int) {
	const (
		// 定义每分钟的秒数
		SecondsPerMinute = 60
	)
	minute = seconds / SecondsPerMinute
	sec = seconds % SecondsPerMinute

	return
}
