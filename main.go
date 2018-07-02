package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"os"
	"strings"
	"time"
)

const (
	APP     = "tagzytout v%s\n"
	VERSION = "1.0.0"

	prefix = "tagzytout:"
)

var (
	debug bool
	path  string
)

func logError(msg string, err error) {
	log.Info(msg)
	log.WithError(err).Debug(msg)
}

func createTag(r *git.Repository, hash plumbing.Hash, tagName string) error {
	tag := object.Tag{
		Name:    tagName,
		Message: tagName,
		Tagger: object.Signature{
			Name:  "Siegfried Ehret",
			Email: "siegfried@ehret.me",
			When:  time.Now(),
		},
		PGPSignature: "",
		Target:       hash,
		TargetType:   plumbing.CommitObject,
	}

	e := r.Storer.NewEncodedObject()
	tag.Encode(e)
	hash, err := r.Storer.SetEncodedObject(e)
	if err != nil {
		logError("Failed to set tag object", err)
		return err
	}

	err = r.Storer.SetReference(plumbing.NewReferenceFromStrings("refs/tags/"+tagName, hash.String()))
	if err != nil {
		logError("Failed to set reference to tag object", err)
		return err
	}

	return nil
}

func init() {
	flag.BoolVar(&debug, "debug", false, "Enable debug logs")
	flag.StringVar(&path, "path", "", "A path to a local repository")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(APP, VERSION))
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, fmt.Sprintln("Example: tagzytout -path=/home/user/gitrepo/"))

		return
	}
}

func main() {
	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	tagPrefixLength := len(prefix)

	if path == "" {
		log.Fatal("Please use a valid path. See `tagzytout -help` for details.")
	}

	r, err := git.PlainOpen(path)
	if err != nil {
		logError("Failed to open repository", err)
	}

	ref, err := r.Head()
	if err != nil {
		logError("Failed to get HEAD", err)
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		logError("Failed to get commit history", err)
	}

	var commits []*object.Commit

	cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})

	count := 0
	for i := len(commits) - 1; i >= 0; i-- {
		c := commits[i]

		for _, line := range strings.Split(c.Message, "\n") {

			tagzytoutIndex := strings.Index(line, prefix)

			if tagzytoutIndex == 0 {
				tagzytoutContent := strings.Trim(line[tagPrefixLength:], " ")

				log.WithFields(log.Fields{
					"tagzytout": tagzytoutContent,
					"commit":    c,
				}).Debug("Found a tag")

				err := createTag(r, c.Hash, tagzytoutContent)

				if err == nil {
					log.Infof("Created tag: %s", tagzytoutContent)
				}

				count++
			}
		}
	}

	log.Infof("Found %d tags", count)
}
