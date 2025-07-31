package printer

import (
	"fmt"
	"sort"

	"github.com/spf13/pflag"
	"github.com/khilk/GitStatsTool/tool/internal/app/executor"
)

const lines = "lines"
const commits = "commits"

const tabularFormat = "tabular"
const csvFormat = "csv"
const jsonFormat = "json"
const jsonLinesFormat = "json-lines"

type Printer interface {
	PrintStatistics(statistics []*executor.Result) error
	ValidateConfig() error
}

type printer struct {
	format     string
	orderBy    string
	statistics []*executor.Result
}

func New(commandFlags *pflag.FlagSet) (Printer, error) {
	p := &printer{}
	var err error
	p.orderBy, err = commandFlags.GetString("order-by")
	if err != nil {
		return nil, fmt.Errorf("can't parse order-by flag: %s", err)
	}

	p.format, err = commandFlags.GetString("format")
	if err != nil {
		return nil, fmt.Errorf("can't parse format flag: %s", err)
	}

	return p, nil
}

func (p *printer) PrintStatistics(statistics []*executor.Result) error {
	p.statistics = statistics
	p.sortStats()
	err := p.printResults()
	if err != nil {
		return fmt.Errorf("can't print statistics: %s", err)
	}
	return nil
}

func (p *printer) sortStats() {
	if p.orderBy == lines {
		sort.Slice(p.statistics, func(i int, j int) bool {
			if p.statistics[i].Statistics.Lines != p.statistics[j].Statistics.Lines {
				return p.statistics[i].Statistics.Lines > p.statistics[j].Statistics.Lines
			}
			if p.statistics[i].Statistics.Commits != p.statistics[j].Statistics.Commits {
				return p.statistics[i].Statistics.Commits > p.statistics[j].Statistics.Commits
			}
			if p.statistics[i].Statistics.Files != p.statistics[j].Statistics.Files {
				return p.statistics[i].Statistics.Files > p.statistics[j].Statistics.Files
			}
			return p.statistics[i].Name < p.statistics[j].Name
		})
		return
	}

	if p.orderBy == commits {
		sort.Slice(p.statistics, func(i int, j int) bool {
			if p.statistics[i].Statistics.Commits != p.statistics[j].Statistics.Commits {
				return p.statistics[i].Statistics.Commits > p.statistics[j].Statistics.Commits
			}
			if p.statistics[i].Statistics.Lines != p.statistics[j].Statistics.Lines {
				return p.statistics[i].Statistics.Lines > p.statistics[j].Statistics.Lines
			}
			if p.statistics[i].Statistics.Files != p.statistics[j].Statistics.Files {
				return p.statistics[i].Statistics.Files > p.statistics[j].Statistics.Files
			}
			return p.statistics[i].Name < p.statistics[j].Name
		})
		return
	}

	sort.Slice(p.statistics, func(i int, j int) bool {
		if p.statistics[i].Statistics.Files != p.statistics[j].Statistics.Files {
			return p.statistics[i].Statistics.Files > p.statistics[j].Statistics.Files
		}
		if p.statistics[i].Statistics.Lines != p.statistics[j].Statistics.Lines {
			return p.statistics[i].Statistics.Lines > p.statistics[j].Statistics.Lines
		}
		if p.statistics[i].Statistics.Commits != p.statistics[j].Statistics.Commits {
			return p.statistics[i].Statistics.Commits > p.statistics[j].Statistics.Commits
		}
		return p.statistics[i].Name < p.statistics[j].Name
	})
}
