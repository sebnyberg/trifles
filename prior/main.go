package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)

	var all []*github.Repository

	for page := 1; ; page++ {
		repos, _, err := client.Repositories.List(ctx, "", &github.RepositoryListOptions{ListOptions: github.ListOptions{Page: page}})
		if err != nil {
			log.Printf("error fetching repositories page %d = %+v\n", page, err)
			break
		}

		log.Printf("fetched %d results", len(repos))

		if len(repos) == 0 {
			log.Printf("done")
			break
		}

		all = append(all, repos...)
	}

	log.Printf("total %d repos", len(all))

	sort.Sort(ByCreatedAt(all))

	for _, r := range all {
		// ignore forks
		if r.GetFork() {
			log.Printf("ignoring fork %q", r.GetFullName())
			continue
		}

		if r.GetDescription() == "" {
			log.Printf("missing description for %q", r.GetFullName())
		}

		fmt.Fprintln(w, r.GetFullName(), "\t", r.GetDescription(), "\t", r.GetCreatedAt().Format("2006-01-02"))
	}

	w.Flush()
}

type ByCreatedAt []*github.Repository

func (b ByCreatedAt) Len() int           { return len(b) }
func (b ByCreatedAt) Less(i, j int) bool { return b[i].GetCreatedAt().Before(b[j].GetCreatedAt().Time) }
func (b ByCreatedAt) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
