package utils

import "math"

type Pagination struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int64 `json:"total_pages"`
}

func NewPagination(page, limit int64) *Pagination {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return &Pagination{
		Page:  page,
		Limit: limit,
	}
}

func (p *Pagination) GetOffset() int64 {
	return (p.Page - 1) * p.Limit
}

func (p *Pagination) SetTotal(totalRows int64) {
	p.TotalRows = totalRows
	p.TotalPages = int64(math.Ceil(float64(p.TotalRows) / float64(p.Limit)))
}
