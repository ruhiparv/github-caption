package models

type Issue struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Comment struct {
	Body string `json:"body"`
}

type CommentRequest struct {
	Issue IssueRequest `json:"issue_info"`
	Body  string       `json:"body"`
}

type ContainsImageResponse struct {
	ContainsImage bool `json:"containsImage"`
}

type IssueRequest struct {
	Owner    string `json:"owner"`
	Repo     string `json:"repo"`
	IssueNum int    `json:"issueNum"`
}
