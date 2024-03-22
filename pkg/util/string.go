package util

import (
	"regexp"
	"strings"
)

func TrimDigit(str string) string {
	lastNonDigitIndex := len(str) - 1
	for i := len(str) - 1; i >= 0; i-- {
		if !isDigit(str[i]) {
			lastNonDigitIndex = i
			break
		}
	}

	// 利用字符串切片获取删除结尾数字后的字符串
	result := strings.TrimSpace(str[:lastNonDigitIndex+1])

	return result
}

// 判断字符是否为数字
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func RemoveStrings(slice []string, removeList []string) []string {
	result := []string{}
	for _, element := range slice {
		found := false
		for _, removeElement := range removeList {
			if element == removeElement {
				found = true
				break
			}
		}
		if !found {
			result = append(result, element)
		}
	}
	return result
}

func ValidateDomains(domains []string) bool {
	// 构建一个正则表达式模式，用于匹配域名
	pattern := `^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`

	// 编译正则表达式模式
	regexpPattern, err := regexp.Compile(pattern)
	if err != nil {
		// 如果正则表达式编译失败，则返回 false 或进行错误处理
		return false
	}

	// 遍历传入的域名列表，检查每个域名是否满足正则表达式模式
	for _, domain := range domains {
		if !regexpPattern.MatchString(domain) {
			// 如果有任何一个域名不满足正则表达式模式，则返回 false
			return false
		}
	}

	// 如果所有域名都满足正则表达式模式，则返回 true
	return true
}
