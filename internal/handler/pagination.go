package handler

import (
	"strconv"
	"strings"
	"time"
)

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

func parseCursor(raw string) (time.Time, uint, bool) {
	if raw == "" {
		return time.Time{}, 0, false
	}
	parts := strings.Split(raw, "_")
	if len(parts) != 2 {
		return time.Time{}, 0, false
	}

	ts, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return time.Time{}, 0, false
	}

	id, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return time.Time{}, 0, false
	}

	return time.Unix(0, ts), uint(id), true
}
