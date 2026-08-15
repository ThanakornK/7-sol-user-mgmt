package utils

import "strconv"

type Pagination struct {
	Page     int64
	PageSize int64
	Total    int64
}

func ParsePagination(pageStr, pageSizeStr string) (Pagination, error) {
	if pageStr == "" {
		pageStr = "1"
	}
	if pageSizeStr == "" {
		pageSizeStr = "10"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return Pagination{}, err
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return Pagination{}, err
	}
	return Pagination{
		Page:     int64(page),
		PageSize: int64(pageSize),
	}, nil
}
