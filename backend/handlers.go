package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getTasks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, gin.H{"data": tasks})
}
