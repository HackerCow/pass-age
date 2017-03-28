package main

import (
	"errors"
	"fmt"
	"github.com/hako/durafmt"
	. "github.com/logrusorgru/aurora"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type passwordEntry struct {
	path    string
	pwName  string
	changed time.Time
}

var entries []*passwordEntry

func getFileAbs(commit *object.Commit, base string, path string) (*passwordEntry, error) {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return nil, err
	}

	return getFile(commit, rel)
}

func getFile(commit *object.Commit, path string) (*passwordEntry, error) {
	c, err := git.References(commit, path)
	if err != nil {
		return nil, err
	}

	if c == nil || len(c) == 0 {
		return nil, errors.New("entry not found")
	}

	return &passwordEntry{
		path:    path,
		pwName:  strings.TrimRight(path, ".gpg"),
		changed: c[len(c)-1].Committer.When,
	}, nil
}

func getEntry(commit *object.Commit, pwName string) (*passwordEntry, error) {
	return getFile(commit, pwName+".gpg")
}

func processFile(commit *object.Commit, base string) filepath.WalkFunc {
	return func(path string, _ os.FileInfo, _ error) error {
		if !strings.HasSuffix(path, ".gpg") {
			return nil
		}

		entry, err := getFileAbs(commit, base, path)
		if err != nil {
			return err
		}

		entries = append(entries, entry)
		return nil
	}
}

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	dir := path.Join(usr.HomeDir, ".password-store")

	r, err := git.PlainOpen(dir)
	if err != nil {
		panic(err)
	}

	ref, err := r.Head()
	if err != nil {
		panic(err)
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		filepath.Walk(dir, processFile(commit, dir))
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].changed.Before(entries[j].changed)
		})

		for _, e := range entries {
			durString := durafmt.Parse(time.Now().Sub(e.changed))
			fmt.Printf("%s was last changed %s ago\n", Blue(e.pwName), Red(durString))
		}
	} else {
		entry, err := getEntry(commit, os.Args[1])
		if err != nil {
			fmt.Println(Red("Error:"), err)
			return
		}

		durString := durafmt.Parse(time.Now().Sub(entry.changed))
		fmt.Printf("%s was last changed %s ago\n", Blue(entry.pwName), Red(durString))
	}
}
