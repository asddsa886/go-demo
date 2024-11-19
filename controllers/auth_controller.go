package controllers

import (
	"fmt"
	"go-demo/global"
	"go-demo/models"
	"go-demo/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 注册
func Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//创建表
	if err := global.Db.AutoMigrate(&models.User{}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 判断用户是否存在
	result := global.Db.Where("username = ?", user.Username).First(&models.User{})
	if result.RowsAffected > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户已存在"})
		return
	}

	// 创建用户
	// 加密
	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Password = hashedPwd

	if err := global.Db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "用户创建成功"})
}

// 登录
func Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 验证用户
	var user models.User
	result := global.Db.Where("username = ?", input.Username).First(&user)
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
		return
	}

	fmt.Println("user.Password", user.Password)
	fmt.Println("input.Password", input.Password)
	// 验证密码
	if err := utils.CheckPassword(input.Password, user.Password); !err {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "密码错误"})
		return
	}

	// 生成token
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("generate token", token)
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// 测试
func Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}
