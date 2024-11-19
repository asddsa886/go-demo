package controllers

import (
	"encoding/json"
	"go-demo/global"
	"go-demo/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 根据id获取文章
func GetArticleById(ctx *gin.Context) {
	id := ctx.Param("id")

	// 先从redis中查询
	key := "article:" + id
	val, err := global.RedisDb.Get(key).Result()
	if err == nil {
		// redis中存在,直接返回
		log.Println("redis中存在,直接返回")
		var article models.Article
		if err := json.Unmarshal([]byte(val), &article); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": article})
		return
	}

	// redis中不存在,查询数据库
	var article models.Article
	if err := global.Db.Where("id = ?", id).First(&article).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 写入redis缓存
	jsonData, err := json.Marshal(article)
	log.Println("写入redis缓存")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := global.RedisDb.Set(key, jsonData, time.Hour*24).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": article})
}

// 创建文章
func CreateArticle(ctx *gin.Context) {
	var article models.Article
	// 绑定请求体
	if err := ctx.ShouldBindJSON(&article); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	//创建表
	if err := global.Db.AutoMigrate(&models.Article{}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建文章
	if err := global.Db.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回文章
	ctx.JSON(http.StatusOK, gin.H{"message": "create article"})
}

// 获取全部文章
func GetAllArticle(ctx *gin.Context) {
	var articles []models.Article
	// 查询数据库
	if err := global.Db.Find(&articles).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 返回文章
	ctx.JSON(http.StatusOK, gin.H{"data": articles})
}
