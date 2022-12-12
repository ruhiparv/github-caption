package routes

import (
	"example/github-caption/controllers"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v48/github"
)

func Routes(client *github.Client) {
	route := gin.Default()

	route.GET("/health", controllers.HealthHandler(client))
	route.GET("/api/v1/github/:owner/:repo/issue/:issue_num/", controllers.GetIssueHandler(client))
	route.GET("/api/v1/github/:owner/:repo/issue/:issue_num/image", controllers.ContainsImageHandler(client))
	route.POST("/api/v1/github/:owner/:repo/issues/:issue_num/comment", controllers.CreateCommentHandler(client))
	route.POST("/api/v1/github/:owner/:repo/issues/:issue_num/identify", controllers.IdentifyAndCommentHandler(client))

	route.Run()
}
