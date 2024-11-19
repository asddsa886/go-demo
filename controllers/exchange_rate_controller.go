package controllers

import (
	"go-demo/global"
	"go-demo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 添加汇率
func AddExchangeRate(ctx *gin.Context) {
	var exchangeRate models.ExchangeRate
	//绑定请求体
	if err := ctx.ShouldBindJSON(&exchangeRate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//创建表
	if err := global.Db.AutoMigrate(&models.ExchangeRate{}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//创建汇率
	if err := global.Db.Create(&exchangeRate).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "add exchange rate"})
}

// 获取全部汇率
func GetExchangeRate(ctx *gin.Context) {
	var exchangeRate []models.ExchangeRate
	if err := global.Db.Find(&exchangeRate).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": exchangeRate})
}
