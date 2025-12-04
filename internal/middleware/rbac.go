package middleware

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func check(e *casbin.Enforcer, sub, obj, act string) {
	ok, _ := e.Enforce(sub, obj, act)
	if ok {
		fmt.Printf("%s CAN %s %s\n", sub, act, obj)
	} else {
		fmt.Printf("%s CANNOT %s %s\n", sub, act, obj)
	}
}

func AuthMiddleware(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// method := c.Request.Method
		check(e, "zhangsan", "/index", "GET")
		check(e, "zhangsan1", "/index", "GET")
		// e.AddPolicy("zhangsan1", "/index", "GET")
		// e.SavePolicy()
		e.RemovePolicy("zhangsan1", "/index", "GET")
		e.SavePolicy()
		check(e, "zhangsan1", "/index", "GET")
		c.Next()
	}
}
