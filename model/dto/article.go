package dto

import "github.com/jinzhu/gorm"

type ArticleDTO struct {
	gorm.Model
	ArticleName string `json:"article_name" gorm:"type:varchar(255);unique_index"`
	URL         string `json:"url" gorm:"type:varchar(255);unique_index"`
	Tag         string `json:"tag" gorm:"type:varchar(255)"`
}

func (ArticleDTO) TableName() string {
	return "rd_article"
}
