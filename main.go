package main

import (
	"context"
	"fmt"

	"os"

	"github.com/google/go-github/v48/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"

	"example/github-caption/routes"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Unable to load env file. Exiting...")
		return
	}

	PAT := os.Getenv("PERSONAL_ACCESS_TOKEN")

	gitContext := context.Background()
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: PAT},
	)
	tokenClient := oauth2.NewClient(gitContext, tokenSource)
	client := github.NewClient(tokenClient)

	routes.Routes(client)

}
