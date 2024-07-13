package util

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type Pagination struct {
	Page   int //页码
	Size   int //每页大小
	Offset int //偏移量
	Total  int //总数
}

// NewPagination 创建分页对象
func NewPagination(c *gin.Context) *Pagination {
	var p = &Pagination{}
	p.Page, p.Size = cast.ToInt(c.Query("page")), cast.ToInt(c.Query("size"))
	//默认分页
	if p.Page == 0 || p.Size == 0 {
		p.Page, p.Size = 1, 10
	}
	//计算偏移量
	p.Offset = (p.Page - 1) * p.Size
	return p
}

// GormPaginate Gorm分页
func (p *Pagination) GormPaginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.Offset).Limit(p.Size)
	}
}
