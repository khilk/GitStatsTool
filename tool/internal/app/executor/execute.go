package executor

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func (g *gitfameExecutor) Execute() error {
	return g.processTree(g.commandFlags.revision, "")
}

func (g *gitfameExecutor) processTree(revision string, path string) error {
	args := []string{"ls-tree", "-l", revision}
	if path != "" {
		args = append(args, path)
	}
	output, err := g.execGitCommand(args...)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()

		name := line[61:]
		entryType := line[7:11]

		if entryType == "tree" {
			if err := g.processTree(revision, name+"/"); err != nil {
				return err
			}
		} else {
			if ok, err := g.validator.IsCorrectPath(name); ok {
				if err != nil {
					return err
				}

				l := 52
				for rn, sz := utf8.DecodeRuneInString(line[l:]); unicode.IsSpace(rn); {
					l += sz
					rn, sz = utf8.DecodeRuneInString(line[l:])
				}
				r := l
				for rn, sz := utf8.DecodeRuneInString(line[r:]); !unicode.IsSpace(rn); {
					r += sz
					rn, sz = utf8.DecodeRuneInString(line[r:])
				}
				size, err := strconv.Atoi(line[l:r])
				if err != nil {
					return err
				}

				if size == 0 {
					err = g.processEmpty(name)
					if err != nil {
						return err
					}
				} else {
					err = g.processFile(name)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (g *gitfameExecutor) execGitCommand(args ...string) ([]byte, error) {
	const name = "git"
	var commonArgs []string
	if g.commandFlags.repository != "." {
		commonArgs = append(commonArgs, "-C", g.commandFlags.repository)
	}
	output, err := exec.Command(name, append(commonArgs, args...)...).Output()
	if err != nil {
		return nil, fmt.Errorf("can't execute git command: %w", err)
	}
	return output, nil
}

func (g *gitfameExecutor) processFile(name string) error {
	output, err := g.execGitCommand("blame", "--porcelain", g.commandFlags.revision, name)
	if err != nil {
		return err
	}
	seen := make(map[string]struct{})
	scanner := bufio.NewScanner(bytes.NewReader(output))
	lines := 0
	commit := ""
	commitsToAuthors := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)
		if len(words) > 1 {
			if len(words) == 4 && len(words[0]) == 40 {
				commit = words[0]
				numberOfLines, err := strconv.Atoi(words[3])
				if err != nil {
					continue
				}
				lines = numberOfLines
				if author, ok := commitsToAuthors[commit]; ok {
					g.results[author].Statistics.Lines += int64(lines)
				}
			}
			if (words[0] == "committer" && g.commandFlags.useCommiter) || (words[0] == "author" && !g.commandFlags.useCommiter) {
				author := line[len(words[0])+1:]
				commitsToAuthors[commit] = author
				if _, ok := g.results[author]; !ok {
					g.results[author] = &Result{Name: author}
				}
				if _, ok := seen[author]; !ok {
					seen[author] = struct{}{}
					g.results[author].Statistics.Files += 1
				}
				g.results[author].Statistics.Lines += int64(lines)
				if _, ok := g.commits[commit]; !ok {
					g.commits[commit] = struct{}{}
					g.results[author].Statistics.Commits += 1
				}
			}
		}
	}
	return nil
}

func (g *gitfameExecutor) processEmpty(name string) error {
	output, err := g.execGitCommand("log", g.commandFlags.revision, "--", name)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewReader(output))
	scanner.Scan()
	line := scanner.Text()
	words := strings.Fields(line)
	currentCommit := words[1]
	scanner.Scan()
	line = scanner.Text()

	var author string
	offset := 8
	for {
		curr, size := utf8.DecodeRuneInString(line[offset:])
		if curr == '<' {
			break
		}
		author += string(curr)
		offset += size
	}
	author = author[:len(author)-1]

	if _, ok := g.results[author]; !ok {
		g.results[author] = &Result{Name: author}
	}
	g.results[author].Statistics.Files += 1
	if _, ok := g.commits[currentCommit]; !ok {
		g.results[author].Statistics.Commits += 1
		g.commits[currentCommit] = struct{}{}
	}
	return nil
}
