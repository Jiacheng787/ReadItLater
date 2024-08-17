package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/garfield-dev-team/read-it-later/model/dto"
)

var db *gorm.DB
var err error

func initDB() {
	db, err = gorm.Open("mysql", "root:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	db.AutoMigrate(&dto.ArticleDTO{})
}

func createArticle(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL 参数为空"})
		return
	}

	articleTitle, err := crawlArticleTitle(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "爬取文章标题失败"})
		return
	}

	article := dto.ArticleDTO{
		ArticleName: articleTitle,
		URL:         url,
	}

	db.Create(&article)

	c.JSON(http.StatusCreated, article)
}

func crawlArticleTitle(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	title := doc.Find("head > meta[property='twitter:title']").Text()
	return title, nil
}

func readArticle(c *gin.Context) {
	var article dto.ArticleDTO
	id := c.Param("id")

	if err := db.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, article)
}

func main() {
	initDB()

	router := gin.Default()

	// 添加中间件用于验证 API 密钥
	router.Use(authMiddleware)

	router.POST("/articles", createArticle)
	router.GET("/articles/:id", readArticle)

	log.Fatal(router.Run(":8080"))
}

func authMiddleware(c *gin.Context) {
	// 假设我们将 API 密钥放在请求头的 "Authorization" 字段中
	apiKey := c.GetHeader("Authorization")

	// 这里可以根据实际需求检查 API 密钥是否有效
	// 例如，你可以检查它是否与预定义的密钥匹配
	if apiKey != "your_actual_api_key" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// 如果API密钥有效，则继续处理请求
	c.Next()
}
