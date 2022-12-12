package controllers

import (
	"example/github-caption/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v48/github"
)

// Checks to see if authentication works for user by listing current repositories
func HealthHandler(client *github.Client) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		_, resp, err := client.Repositories.List(ctx, "", nil)
		if err != nil {
			ctx.IndentedJSON(resp.StatusCode, gin.H{"status": GetStatusMessage(resp)})
		} else {
			ctx.IndentedJSON(resp.StatusCode, gin.H{"status": GetStatusMessage(resp)})
		}
		ctx.Next()
	}

	return gin.HandlerFunc(fn)
}

// Gets Github Issue by querying from owner, repository, and issue number
func GetIssueHandler(client *github.Client) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {

		issueReq, err := CreateIssueRequest(ctx)
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Next()
			return
		}

		issue, err := GetIssue(client, ctx, *issueReq)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			ctx.Next()
			return
		}

		ctx.IndentedJSON(http.StatusOK, issue)
		ctx.Next()
	}

	return gin.HandlerFunc(fn)
}

// Checks to see if a certain issue contains an image in its body
func ContainsImageHandler(client *github.Client) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {

		issueReq, err := CreateIssueRequest(ctx)
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Next()
			return
		}

		issue, err := GetIssue(client, ctx, *issueReq)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			ctx.Next()
			return
		}

		containsImage := models.ContainsImageResponse{ContainsImage: false}

		if strings.Contains(issue.Body, "![image]") {
			containsImage.ContainsImage = true
		}

		ctx.IndentedJSON(http.StatusOK, containsImage)
		ctx.Next()
	}

	return gin.HandlerFunc(fn)
}

// Posts a comment on a given issue
func CreateCommentHandler(client *github.Client) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {

		commentReq, err := CreateCommentRequest(ctx)
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Next()
			return
		}

		// Comment is not used as API is specified to respond only with success or error string
		_, err = CreateComment(client, ctx, *commentReq)
		if err != nil {
			print(err.Error())
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			ctx.Next()
			return
		}

		ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Success"})
		ctx.Next()
	}

	return gin.HandlerFunc(fn)
}

// Posts a comment on an issue if the issue contains an image
func IdentifyAndCommentHandler(client *github.Client) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {

		issueReq, err := CreateIssueRequest(ctx)
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Next()
			return
		}

		issue, err := GetIssue(client, ctx, *issueReq)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			ctx.Next()
			return
		}

		if strings.Contains(issue.Body, "![image]") {
			time, err := currentTime()
			if err != nil {
				ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				ctx.Next()
				return
			}
			body := fmt.Sprintf("I found an image in your issue at %s", time)
			commentReq := models.CommentRequest{Issue: *issueReq, Body: body}
			_, err = CreateComment(client, ctx, commentReq)
			if err != nil {
				ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				ctx.Next()
				return
			}
			ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully commented on issue"})
			ctx.Next()
			return
		}

		ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Did not find image in issue"})
		ctx.Next()
	}

	return gin.HandlerFunc(fn)
}
