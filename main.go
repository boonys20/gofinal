package main

import (
	c "github.com/boonys20/gofinal/customer"
	m "github.com/boonys20/gofinal/middleware"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(m.AuthMiddleware)
	r.GET("/customers", c.GetCustomersHandler)
	r.GET("/customers/:id", c.GetCustomerByIdHandler)
	r.POST("/customers", c.CreateCustomerHandler)
	r.PUT("/customers/:id", c.UpdateCustomersHandler)
	r.DELETE("/customers/:id", c.DeleteCustomersHandler)
	return r
}

func main() {
	r := setupRouter()
	r.Run(":2009")
}
