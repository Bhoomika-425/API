package handler

import (
	"log"
	"project/internal/auth"
	"project/internal/middleware"
	service "project/internal/service"

	"github.com/gin-gonic/gin"
)

func API(a auth.UserAuth, svc service.UserService) *gin.Engine {
	r := gin.New()

	m, err := middleware.NewMiddleware(a)
	if err != nil {
		log.Panic("middlewares not setup")
		return nil
	}
	h, err := Newhandler(svc)
	if err != nil {
		log.Panic("middlewares not setup")
		return nil
	}

	r.Use(m.Log(), gin.Recovery())

	r.GET("/check", Check)

	r.POST("/signup", h.SignUp)
	r.POST("/signin", h.Login)

	r.POST("/addCompany", m.Authenticate(h.AddCompany))
	r.GET("/view/allcompany", m.Authenticate(h.ViewAllCompanies))
	r.GET("/viewcompany/:id", m.Authenticate(h.ViewCompany))

	r.POST("/addJobs/:cid", m.Authenticate(h.CreateJobs))
	r.GET("/view/alljobs", m.Authenticate(h.AllJobs))

	r.GET("/viewjob/:cid", m.Authenticate(h.JobsByCid)) //

	r.GET("/jobs/:id", m.Authenticate(h.JobByJID)) //

	r.POST("/jobApp", m.Authenticate(h.JobApplicationById)) //null
	r.POST("/forgotpassword", (h.ForgotPassword))
	r.POST("/setpassword", (h.AddingNewPassword))
	return r

}

func Check(c *gin.Context) {
	c.JSON(200, gin.H{
		"Message": "ok",
	})
}
