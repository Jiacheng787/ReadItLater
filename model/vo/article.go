package vo

type ArticleVo struct {
	ID          int    `json:"id"`
	ArticleName string `json:"article_name" gorm:"type:varchar(255);unique_index"`
	URL         string `json:"url" gorm:"type:varchar(255);unique_index"`
	Tag         string `json:"tag" gorm:"type:varchar(255)"`
}
