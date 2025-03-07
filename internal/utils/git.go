package utils

import (
	"bytes"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"io"
	"strings"
)

var maxSize = 255 * 1024 // 255 KB
var includeDiffFileLimit = 50
var errStopIteration = errors.New("stop iteration")

var ErrNoCommit = errors.New("no commit")

func OpenGitRepo(dir string) (*GitRepo, error) {
	repo, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, err
	}
	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	rootDir := wt.Filesystem.Root()

	return &GitRepo{Dir: dir, repo: repo, RootDir: rootDir}, nil
}

type GitRepo struct {
	Dir     string
	RootDir string

	repo *git.Repository
}

func (r *GitRepo) Commit(msg string, des string) (string, error) {
	w, _ := r.repo.Worktree()
	if des != "" {
		msg = fmt.Sprintf("%s\n\n%s", msg, des)
	}

	hash, err := w.Commit(msg, &git.CommitOptions{})
	if err != nil {
		return "", err
	}

	return hash.String(), err
}

func (r *GitRepo) AddAll() error {
	w, err := r.repo.Worktree()
	if err != nil {
		return err
	}

	return w.AddWithOptions(&git.AddOptions{All: true})
}

func (r *GitRepo) HasStagedChanges() (bool, error) {
	w, err := r.repo.Worktree()
	if err != nil {
		return false, err
	}

	status, err := w.Status()
	if err != nil {
		return false, err
	}

	return !status.IsClean(), nil
}

func (r *GitRepo) GitDiff() (string, error) {
	w, err := r.repo.Worktree()
	if err != nil {
		return "", err
	}

	status, err := w.Status()
	if err != nil {
		return "", err
	}

	var headTree *object.Tree
	head, err := r.repo.Head()
	if err != nil {
		if !errors.Is(err, plumbing.ErrReferenceNotFound) {
			return "", err
		}
		// This is a new repository with no commits
	} else {
		commit, err := r.repo.CommitObject(head.Hash())
		if err != nil {
			return "", err
		}
		headTree, err = commit.Tree()
		if err != nil {
			return "", err
		}
	}

	var sb strings.Builder

	counter := 0

	for path, fileStatus := range status {
		if counter > includeDiffFileLimit {
			break
		}
		// Get old content (from HEAD)
		var oldContent string
		if headTree != nil && fileStatus.Staging != git.Added {
			file, err := headTree.File(path)
			if err == nil {
				oldContent, err = file.Contents()
				if err != nil {
					fmt.Printf("Error reading file from HEAD: %s\n", err)
				}
			}
		}

		// Get new content (from stage/index)
		var newContent string
		if fileStatus.Staging != git.Deleted {
			// For staged changes, we need to get the blob from the git object database
			// First get the file from the worktree
			f, err := w.Filesystem.Open(path)
			if err != nil {
				fmt.Printf("Error opening file: %s\n", err)
				continue
			}

			content, err := io.ReadAll(f)
			_ = f.Close()
			if err != nil {
				fmt.Printf("Error reading file: %s\n", err)
				continue
			}
			newContent = string(content)
		}

		if isBinary([]byte(oldContent)) || isBinary([]byte(newContent)) {
			sb.WriteString("File: " + path + "\n")
			sb.WriteString("Binary file <changes>\n")
			pterm.Warning.Printf("Binary file: %s\n", path)
			continue
		}
		if len(oldContent) > maxSize || len(newContent) > maxSize {
			sb.WriteString("File: " + path + "\n")
			sb.WriteString("File too large <changes>\n")
			pterm.Info.Printf("File too large: %s\n", path)
			continue
		}

		//d := diff.Diff(oldContent, newContent)
		//sb.WriteString("File: " + path + "\n")
		//sb.WriteString(d)
		//sb.WriteString("\n")
		//sb.WriteString("-----------------------------------------\n")
		edits := myers.ComputeEdits(span.URIFromPath(path), oldContent, newContent)
		d := fmt.Sprint(gotextdiff.ToUnified(path, path, oldContent, edits))
		sb.WriteString("File: " + path + "\n")
		// to large diff no point of showing or sending i guess
		if countLlines(d) > 150 {
			sb.WriteString("large <changes>\n")
		} else {
			sb.WriteString(d)
		}
		sb.WriteString("\n")
		sb.WriteString("-----------------------------------------\n")
	}

	return sb.String(), nil
}

// GetLastNCommitMsg get last n commit message
func (r *GitRepo) GetLastNCommitMsg(n int) ([]string, error) {
	// Get HEAD reference
	ref, err := r.repo.Head()
	if err != nil {
		return nil, ErrNoCommit
	}

	// Get commit iterator
	commitIter, err := r.repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, fmt.Errorf("error getting commit iterator: %w", err)
	}
	defer commitIter.Close()

	// Collect commits
	commits := make([]string, 0, n)
	err = commitIter.ForEach(func(c *object.Commit) error {
		commit := strings.Split(c.Message, "\n")[0]
		commits = append(commits, commit)
		if len(commits) >= n {
			return errStopIteration
		}
		return nil
	})
	if err != nil && !errors.Is(err, errStopIteration) {
		return nil, fmt.Errorf("error iterating commits: %w", err)
	}

	return commits, nil
}

// isBinary checks if content is likely binary by looking for null bytes
// or by checking if it contains a high ratio of non-printable characters
func isBinary(content []byte) bool {
	if len(content) == 0 {
		return false
	}

	// Quick check: if it contains NULL bytes, it's almost certainly binary
	if bytes.Contains(content, []byte{0}) {
		return true
	}

	// Only check a reasonable sample size to improve performance
	sampleSize := 8000
	if len(content) < sampleSize {
		sampleSize = len(content)
	}

	// Count non-printable, non-whitespace characters
	nonPrintable := 0
	for i := 0; i < sampleSize; i++ {
		c := content[i]
		if (c < 32 && c != 9 && c != 10 && c != 13) || c > 126 {
			nonPrintable++
		}
	}

	// If more than 10% is non-printable, consider it binary
	return nonPrintable > sampleSize/10
}

func countLlines(s string) int {
	count := 0
	for _, c := range s {
		if c == '\n' {
			count++
		}
	}
	return count
}
