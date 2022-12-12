package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v48/github"

	"example/github-caption/models"
)

// Constructs issue request struct from request parameters
func CreateIssueRequest(ctx *gin.Context) (*models.IssueRequest, error) {
	owner := ctx.Param("owner")
	repo := ctx.Param("repo")
	issueNum, err := strconv.Atoi(ctx.Param("issue_num"))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to parse issue number %s", err.Error()))
	}

	issueReq := models.IssueRequest{Owner: owner, Repo: repo, IssueNum: issueNum}

	return &issueReq, nil
}

// Calls github api to retrieve issue
func GetIssue(client *github.Client, ctx context.Context, ir models.IssueRequest) (*models.Issue, error) {

	githubIssue, resp, err := client.Issues.Get(ctx, ir.Owner, ir.Repo, ir.IssueNum)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error retrieving issue, status %s: %s", resp.Status, err))
	}

	githubIssueJson, err := json.Marshal(githubIssue)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to create github issue JSON: %s", err))
	}

	var subsetIssue models.Issue
	err = json.Unmarshal([]byte(githubIssueJson), &subsetIssue)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to convert issue JSON to model: %s", err))
	}

	return &subsetIssue, nil

}

// Constructs comment request struct from request parameters
func CreateCommentRequest(ctx *gin.Context) (*models.CommentRequest, error) {
	var newComment models.Comment

	if err := ctx.BindJSON(&newComment); err != nil {
		return nil, errors.New(fmt.Sprintf("Could not retrieve comment body in request: %s", err.Error()))
	}

	issueRequest, err := CreateIssueRequest(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to create issue request: %s", err.Error()))
	}

	commentReq := models.CommentRequest{Issue: *issueRequest, Body: newComment.Body}

	return &commentReq, nil

}

// Calls github api to post a comment on an issue
func CreateComment(client *github.Client, ctx context.Context, cr models.CommentRequest) (*github.IssueComment, error) {

	commentBody := &github.IssueComment{Body: &cr.Body}

	issueComment, response, err := client.Issues.CreateComment(ctx, cr.Issue.Owner, cr.Issue.Repo, cr.Issue.IssueNum, commentBody)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error commenting on issue, status %s: %s", response.Status, err.Error()))
	}

	return issueComment, nil
}

// Provides current time in ET with RFC1123 format
func currentTime() (string, error) {
	newYork, err := time.LoadLocation("America/New_York")
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error retrieving America/New York time zone: %s", err.Error()))
	}
	time := time.Now().In(newYork).Format(time.RFC1123)

	return time, nil
}

// Returns status message from a go-github response
func GetStatusMessage(response *github.Response) string {
	return strings.Split(response.Status, " ")[1]
}
