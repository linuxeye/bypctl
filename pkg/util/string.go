package util

import (
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
