package util

import (
	"sync"
	"time"
	"fmt"
)

//return length of sync map
func len(m sync.Map) int {
	len := 0
	m.Range(func(_, _ interface{}) bool {
		len++
		return true
	})
	return len
}

//return list yyyy_MM from current to prvious
//@param n number of month in list
func GetListYearMonthFromCurrent(n int) ([]string) {
	arr := make([]string, n)
	now := time.Now()
	for i := n; i > 0; i-- {
		year := int(now.Year())
		month := int(now.Month())
		arr = append(arr, fmt.Sprintf("%d_%d", year, month))
		now = now.AddDate(0, -1, 0)
	}
	return arr
}
