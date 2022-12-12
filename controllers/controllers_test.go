package controllers

import (
	"context"
	"example/github-caption/models"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"

	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
)

func TestCreateIssueRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	cases := []struct {
		testName    string
		input       []gin.Param
		expected    models.IssueRequest
		shouldError bool
	}{
		{
			testName: "Successful call",
			input: []gin.Param{
				{Key: "owner", Value: "user1"},
				{Key: "repo", Value: "repo1"},
				{Key: "issue_num", Value: "1"},
			},
			expected:    models.IssueRequest{Owner: "user1", Repo: "repo1", IssueNum: 1},
			shouldError: false,
		},
		{
			testName: "Issue number input error",
			input: []gin.Param{
				{Key: "owner", Value: "user1"},
				{Key: "repo", Value: "repo1"},
				{Key: "issue_num", Value: "abcde"},
			},
			expected:    models.IssueRequest{},
			shouldError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.testName, func(t *testing.T) {
			c.Params = tc.input
			actual, err := CreateIssueRequest(c)
			if tc.shouldError {
				if err == nil {
					t.Errorf("Did not return error when it should have")
				}
			} else if !cmp.Equal(actual, &tc.expected) {
				t.Errorf("Recieved %v, Expected %v", actual, tc.expected)
			}
		})
	}

}

func TestGetIssue(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		testName         string
		mockedHTTPClient *http.Client
		expectedID       int64
		shouldError      bool
	}{
		{
			testName: "Successful call",
			mockedHTTPClient: mock.NewMockedHTTPClient(
				mock.WithRequestMatch(
					mock.GetReposIssuesByOwnerByRepoByIssueNumber,
					github.Issue{
						ID: github.Int64(200),
					},
				),
			),
			expectedID:  int64(200),
			shouldError: false,
		},
		{
			testName: "Github API returns error",
			mockedHTTPClient: mock.NewMockedHTTPClient(
				mock.WithRequestMatchHandler(
					mock.GetReposIssuesByOwnerByRepoByIssueNumber,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						mock.WriteError(
							w,
							http.StatusInternalServerError,
							"something went wrong",
						)
					}),
				),
			),
			expectedID:  int64(-1),
			shouldError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.testName, func(t *testing.T) {
			c := github.NewClient(tc.mockedHTTPClient)
			issueReq := models.IssueRequest{Owner: "owner", Repo: "repo", IssueNum: 1}
			actualIssue, err := GetIssue(c, ctx, issueReq)
			if tc.shouldError {
				if err == nil {
					t.Errorf("Did not return error when it should have")
				}
			} else if !cmp.Equal(int64(actualIssue.ID), int64(tc.expectedID)) {
				t.Errorf("Recieved %v, Expected %v", actualIssue.ID, tc.expectedID)
			}
		})
	}
}
