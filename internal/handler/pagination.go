package handler

import "strconv"

const (
	defaultPage = 1
	defaultSize = 20
	maxSize     = 50
)

func parsePageSize(pageStr, sizeStr string) (int, int) {
	page := defaultPage
	size := defaultSize

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
		if s > maxSize {
			s = maxSize
		}
		size = s
	}
	return page, size
}
