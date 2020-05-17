package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	gistID      = flag.String("gist", "", "Gist ID")
	maxComments = flag.Int("max-comments", 10, "Max comments allowed")
)

var client *github.Client
var ctx = context.Background()

func listCommentsID(gistID string) ([]*int64, error) {
	var ids []*int64
	page := 1

	for {
		opt := &github.ListOptions{Page: page, PerPage: 30}
		comments, _, err := client.Gists.ListComments(ctx, gistID, opt)

		if err != nil {
			return nil, err
		}

		if len(comments) == 0 {
			break
		}

		for _, comment := range comments {
			ids = append(ids, comment.ID)
		}

		page++
	}

	return ids, nil
}

func deleteComment(gistID string, comments []*int64) (err error) {

	for _, comment := range comments {
		_, err = client.Gists.DeleteComment(ctx, gistID, *comment)

		if err != nil {
			return
		}
	}

	return
}

func genRandom(total int, count int) []int {
	var list []int

	for {
		idx := rand.Intn(total - 1)
		list = append(list, idx)

		if len(list) == count {
			break
		}
	}

	return list
}

func containsArray(items []int, value int) bool {
	for _, i := range items {
		if i == value {
			return true
		}
	}

	return false
}

func getCommentsToDelete(comments []*int64, skip []int) []*int64 {
	var toDelete []*int64
	count := len(comments)

	for i := 0; i < count; i++ {
		if containsArray(skip, i) {
			continue
		}

		toDelete = append(toDelete, comments[i])
	}

	return toDelete
}

func main() {
	flag.Parse()
	token := os.Getenv("GH_TOKEN")

	if token == "" {
		log.Fatal("Unauthorized: No token present")
		return
	}

	if *gistID == "" {
		log.Fatal("You need to specify a non-empty value for the flags `-gist`")
		return
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	comments, err := listCommentsID(*gistID)

	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	total := len(comments)
	if total <= *maxComments {
		return
	}

	fmt.Println("Total Comments: ", total)
	skipIndex := genRandom(total, *maxComments)

	commentsDelete := getCommentsToDelete(comments, skipIndex)
	fmt.Println("Comments delete: ", len(commentsDelete))
	err = deleteComment(*gistID, commentsDelete)

	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
}
