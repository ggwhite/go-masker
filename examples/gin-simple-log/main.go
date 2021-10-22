package main

import (
	"log"
	"net/http"

	masker "github.com/ggwhite/go-masker"
	"github.com/gin-gonic/gin"
)

func hello() func(c *gin.Context) {

	type Request struct {
		Name  string `json:"name" mask:"name"`
		Email string `json:"email" mask:"email"`
	}

	type Response struct {
		Success bool
		Message string
	}

	return func(c *gin.Context) {

		req := new(Request)
		res := new(Response)

		if err := c.ShouldBindJSON(req); err != nil {
			res.Success = false
			res.Message = err.Error()
			c.IndentedJSON(http.StatusBadRequest, res)
			return
		}

		_req, _ := masker.Struct(req)

		// use the origin struct what you want
		log.Println(req)

		// log the masked struct
		log.Println(_req)

		res.Success = true
		res.Message = "OK"

		c.IndentedJSON(http.StatusOK, res)
	}
}

func main() {
	engine := gin.Default()
	engine.POST("", hello())
	engine.Run(":3000")
}
