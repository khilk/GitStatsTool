package executor

import (
	_ "embed"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/khilk/GitStatsTool/tool/internal/app/executor/validator"
)

type Result struct {
	Name       string
	Statistics struct {
		Lines   int64
		Commits int64
		Files   int64
	}
}

type GitfameExecutor interface {
	Init(commandFlags *pflag.FlagSet) error
	GetResult() []*Result
	Execute() error
}

type gitfameExecutor struct {
	commits      map[string]struct{}
	results      map[string]*Result
	commandFlags struct {
		repository  string
		revision    string
		useCommiter bool
	}
	validator validator.Validator
}

func New() GitfameExecutor {
	executor := &gitfameExecutor{
		results:   make(map[string]*Result),
		commits:   make(map[string]struct{}),
		validator: validator.New(),
	}
	return executor
}

func (g *gitfameExecutor) Init(commandFlags *pflag.FlagSet) error {
	err := g.setFlags(commandFlags)
	if err != nil {
		return fmt.Errorf("can't set flags: %s", err)
	}

	err = g.validator.Init(commandFlags)
	if err != nil {
		return fmt.Errorf("can't initialize validator: %s", err)
	}

	return nil
}

func (g *gitfameExecutor) setFlags(commandFlags *pflag.FlagSet) error {
	var err error

	g.commandFlags.repository, err = commandFlags.GetString("repository")
	if err != nil {
		return fmt.Errorf("can't parse repository flag: %s", err)
	}

	g.commandFlags.revision, err = commandFlags.GetString("revision")
	if err != nil {
		return fmt.Errorf("can't parse  revision flag: %s", err)
	}

	g.commandFlags.useCommiter, err = commandFlags.GetBool("use-committer")
	if err != nil {
		return fmt.Errorf("can't parse use-committer flag: %s", err)
	}

	return nil
}

func (g *gitfameExecutor) GetResult() []*Result {
	res := make([]*Result, 0, len(g.results))
	for _, result := range g.results {
		res = append(res, result)
	}

	return res
}
