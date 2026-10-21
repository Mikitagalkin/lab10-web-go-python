package handlers

import (
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
)

// Ping — простой health-check.
func Ping(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message":   "pong from Go",
        "timestamp": time.Now().Unix(),
    })
}

// Hello — приветствие с валидацией параметра :name.
func Hello(c *gin.Context) {
    name := strings.TrimSpace(c.Param("name"))
    if name == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "name is required",
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "message": "Hello, " + name + "!",
    })
}

// EchoRequest — структура запроса.
type EchoRequest struct {
    Text string `json:"text" binding:"required,min=1,max=200"`
}

// Echo — возвращает переданный текст с валидацией.
func Echo(c *gin.Context) {
    var req EchoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "validation failed",
            "details": err.Error(),
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "echo":      req.Text,
        "text_len":  len(req.Text),
        "timestamp": time.Now().Unix(),
    })
}
