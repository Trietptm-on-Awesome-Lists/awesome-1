package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/google/go-github/v24/github"
	"golang.org/x/oauth2"

	"github.com/junary/awesome"
)

const (
	exitCodeSuccess = iota
	exitCodeError
	exitCodeHelp
)

const repository = "sindresorhus/awesome"

type cli struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

func main() {
	cli := &cli{os.Stdin, os.Stdout, os.Stderr}
	os.Exit(cli.run(os.Args))
}

func (c *cli) run(args []string) int {
	f := flag.NewFlagSet(args[0], flag.ContinueOnError)
	f.SetOutput(c.err)

	err := f.Parse(args[1:])
	if err != nil {
		if err == flag.ErrHelp {
			return exitCodeHelp
		}
		fmt.Fprintln(c.err, "cases:", err)
		return exitCodeError
	}

	ctx := context.Background()
	client := newClient(ctx, os.Getenv("AWESOME_GITHUB_TOKEN"))
	content, err := getContent(client, ctx, repository)
	if err != nil {
		fmt.Fprintln(c.err, "cases:", err)
	}

	query := f.Arg(0)
	repos, err := search(content, query)
	if err != nil {
		fmt.Fprintln(c.err, "cases:", err)
	}
	for _, repo := range repos {
		fmt.Fprintln(c.out, repo.Url)
	}

	return exitCodeSuccess
}

func newClient(ctx context.Context, token string) *github.Client {
	if token == "" {
		return github.NewClient(nil)
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func getContent(client *github.Client, ctx context.Context, repository string) (string, error) {
	var owner, repo string
	if strings.Contains(repository, "/") {
		split := strings.SplitN(repository, "/", 2)
		owner = split[0]
		repo = split[1]
	}
	readme, _, err := client.Repositories.GetReadme(ctx, owner, repo, &github.RepositoryContentGetOptions{})
	if err != nil {
		return "", err
	}
	return readme.GetContent()
}

func search(content, query string) ([]*awesome.Repository, error) {
	req := awesome.SearchRequest{Query: query}
	b := []byte(content)
	res, err := awesome.Search(b, req)
	if err != nil {
		return nil, err
	}
	return res.Repositories, nil
}
