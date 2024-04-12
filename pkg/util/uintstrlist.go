package util

import (
	"strconv"
	"strings"
)

// UintList2Str .
func UintList2Str(ls []uint) string {

	var strBuffer strings.Builder
	for _, x := range ls {
		strBuffer.WriteString(strconv.Itoa(int(x)))
		strBuffer.WriteByte(',')
	}

	str := strBuffer.String()
	if len(str) > 1 {
		str = str[:len(str)-1]
	}
	return str
}

// StrList2Str .
func StrList2Str(ls []string) string {

	var strBuffer strings.Builder
	if len(ls) == 0 {
		return ""
	}

	for _, x := range ls {
		strBuffer.WriteString(x)
		strBuffer.WriteByte(',')
	}

	str := strBuffer.String()
	if len(str) > 1 {
		str = str[:len(str)-1]
	}
	return str
}

// Str2UintList .
func Str2UintList(str string) []uint {

	strList := strings.Split(str, ",")
	intList := make([]uint, 0)
	for _, s := range strList {
		intVal, err := strconv.Atoi(s)
		if err != nil {
			return []uint{}
		}
		intList = append(intList, uint(intVal))
	}

	return intList
}

// DeleteSubIntList .
func DeleteSubIntList(ls, subls []int) []int {
	for _, sub := range subls {

		j := 0
		for i, s := range ls {
			if s == sub {
				ls[i] = ls[j]
				ls[j] = sub
				j++
			}
		}

		ls = ls[j:]
	}

	return ls
}

// AddSubIntList .
func AddSubIntList(ls, subls []int) []int {
	lsSet := make(map[int]struct{}, 0)
	var null struct{}

	for _, s := range ls {
		lsSet[s] = null
	}

	for _, sub := range subls {
		lsSet[sub] = null
	}

	rtnLs := make([]int, 0)
	for s := range lsSet {
		rtnLs = append(rtnLs, s)
	}

	return rtnLs
}

// DomainStr .
func DomainStr(domainStr string) (string, int, error) {
	if strings.Contains(domainStr, ":") {
		domainString := strings.Split(domainStr, ":")
		domain := domainString[0]
		port, err := strconv.Atoi(domainString[1])
		if err != nil {
			return "", 0, nil
		}
		return domain, port, nil
	} else {
		return domainStr, 80, nil
	}
}

// SliceItemStrExist .
func SliceItemStrExist(target []string, src string) bool {
	for _, t := range target {
		if t == src {
			return true
		}
	}
	return false
}

func AllElementsExist(src, target []string) bool {
	for _, s := range src {
		found := false
		for _, t := range target {
			if s == t {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// 两个字符串列表中是否有至少两个元素相同
func CheckElements(src []string, target []string) bool {
	count := 0

	for _, s := range src {
		for _, t := range target {
			if s == t {
				count++
				if count >= 2 {
					return true
				}
			}
		}
	}
	return false
}

// 两个字符串列表中是否有至少1个元素相同
func HasCommonElement(list1 []string, list2 []string) bool {
	for _, str1 := range list1 {
		for _, str2 := range list2 {
			if str1 == str2 {
				return true
			}
		}
	}
	return false
}

// 判断list2是否是list1的子切片（subslice）
func IsSubSlice(list1, list2 []string) bool {
	if len(list1) < len(list2) {
		return false
	}

	for i := 0; i <= len(list1)-len(list2); i++ {
		match := true
		for j := 0; j < len(list2); j++ {
			if list1[i+j] != list2[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}

	return false
}

// SliceItemIntExist .
func SliceItemIntExist(target []int, src int) bool {
	for _, t := range target {
		if t == src {
			return true
		}
	}
	return false
}

// SliceItemUintExist .
func SliceItemUintExist(target []uint, src uint) bool {
	for _, t := range target {
		if t == src {
			return true
		}
	}
	return false
}

// DelUintSliceItemDuplicate .
func DelUintSliceItemDuplicate(arr []uint) []uint {
	result := make([]uint, 0, len(arr))
	temp := map[uint]struct{}{}
	for _, item := range arr {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// DelStrSliceItemDuplicate .
func DelStrSliceItemDuplicate(arr []string) []string {
	result := make([]string, 0, len(arr))
	temp := map[string]struct{}{}
	for _, item := range arr {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// CompareStrList returns the differences between listA and listB and both-have list.
func CompareStrList(lsA, lsB []string) ([]string, []string, []string) {

	bothSet := make(map[string]struct{})
	var null struct{}

	lsAOnly := make([]string, 0)
	for _, a := range lsA {
		if !InStrList(a, lsB) {
			lsAOnly = append(lsAOnly, a)
		} else {
			bothSet[a] = null
		}
	}

	lsBOnly := make([]string, 0)
	for _, b := range lsB {
		if !InStrList(b, lsA) {
			lsBOnly = append(lsBOnly, b)
		} else {
			bothSet[b] = null
		}
	}

	bothLs := make([]string, 0)
	for ele := range bothSet {
		bothLs = append(bothLs, ele)
	}

	return lsAOnly, lsBOnly, bothLs
}

// CompareUintList returns the differences between listA and listB and both-have list.
func CompareUintList(lsA, lsB []uint) ([]uint, []uint, []uint) {

	bothSet := make(map[uint]struct{})
	var null struct{}

	lsAOnly := make([]uint, 0)
	for _, a := range lsA {
		if !InUintList(a, lsB) {
			lsAOnly = append(lsAOnly, a)
		} else {
			bothSet[a] = null
		}
	}

	lsBOnly := make([]uint, 0)
	for _, b := range lsB {
		if !InUintList(b, lsA) {
			lsBOnly = append(lsBOnly, b)
		} else {
			bothSet[b] = null
		}
	}

	bothLs := make([]uint, 0)
	for ele := range bothSet {
		bothLs = append(bothLs, ele)
	}

	return lsAOnly, lsBOnly, bothLs
}

// InStrList .
func InStrList(ele string, ls []string) bool {
	for _, data := range ls {
		if data == ele {
			return true
		}
	}
	return false
}

// InUintList .
func InUintList(ele uint, ls []uint) bool {
	for _, data := range ls {
		if data == ele {
			return true
		}
	}
	return false
}

// Str2Uint .
func Str2Uint(str string) uint {
	num, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0
	}
	return uint(num)
}

// UintList2StrList .
func UintList2StrList(ids []uint, prefix string) []any {
	var result []any
	for _, v := range ids {
		result = append(result, prefix+strconv.Itoa(int(v)))
	}
	return result
}

func CompareEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// UniqueString 去重
func UniqueString(m []string) []string {
	d := make([]string, 0)
	tempMap := make(map[string]bool, len(m))
	for _, v := range m { // 以值作为键名
		if tempMap[v] == false {
			tempMap[v] = true
			d = append(d, v)
		}
	}
	return d
}

func GetUniqueDomains(hosts []string) []string {
	result := make([]string, 0)

	for _, host := range hosts {
		parts := strings.Split(host, ".")
		if len(parts) == 2 {
			result = append(result, host, "*."+host)
		} else if len(parts) > 2 {
			result = append(result, "*."+strings.Join(parts[1:], "."))
		}
	}
	return UniqueString(result)
}
