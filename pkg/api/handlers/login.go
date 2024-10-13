package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginPageHandler(c *gin.Context) {

	c.HTML(http.StatusOK, "login.html", nil)

}
