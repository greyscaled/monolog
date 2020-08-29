// Package xgit wraps go-git with exposed services
package xgit

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	xconfig "vapurrmaid/monolog/config"

	"github.com/go-git/go-git/v5"
)

// LogLatest prints the latest commit to stdout
func LogLatest(path string) {
	cf, err := xconfig.New()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal(err)
	}

	opts := git.LogOptions{
		Since: &cf.Latest,
	}

	cItr, err := repo.Log(&opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cItr.Close()

	buf := bufio.NewReader(os.Stdin)

	for {
		cmt, err := cItr.Next()
		if err != nil {
			break
		}

		if cmt == nil {
			break
		}

		printCommit(cmt, repo)

		_, err = buf.ReadString('\n')
		if err != nil {
			break
		}
	}

	cf.Latest = time.Now()
	cf.Save()
}

func printCommit(cmt *object.Commit, repo *git.Repository) {
	id := cmt.ID()
	fmt.Println(id)
	fmt.Println(cmt.Message)
	printRepoURL(&id, repo)
	fmt.Println()
}

func printRepoURL(hash *plumbing.Hash, repo *git.Repository) {
	rmts, err := repo.Remotes()
	if err != nil {
		return
	}

	if len(rmts) <= 0 {
		return
	}

	split := strings.Split(rmts[0].String(), "  ")
	if len(split) < 2 {
		return
	}

	splitURL := strings.Split(split[1], " ")
	if len(splitURL) < 2 {
		return
	}

	url := strings.TrimSpace(splitURL[0])

	fmt.Println(path.Join(url, "commits", hash.String()))
}
