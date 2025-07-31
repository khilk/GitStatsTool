package printer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

type entry struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

func (p *printer) printResults() error {
	outputStream := os.Stdout

	switch p.format {
	case tabularFormat:
		w := tabwriter.NewWriter(outputStream, 0, 0, 1, ' ', 0)
		fmt.Fprintln(w, "Name\tLines\tCommits\tFiles")
		for _, item := range p.statistics {
			fmt.Fprintln(w,
				strings.Join([]string{item.Name,
					strconv.Itoa(int(item.Statistics.Lines)),
					strconv.Itoa(int(item.Statistics.Commits)),
					strconv.Itoa(int(item.Statistics.Files))},
					"\t"))
		}
		w.Flush()
	case csvFormat:
		w := csv.NewWriter(outputStream)
		err := w.Write([]string{"Name", "Lines", "Commits", "Files"})
		if err != nil {
			return fmt.Errorf("can't print output in csv format: %s", err)
		}
		for _, item := range p.statistics {
			err := w.Write([]string{item.Name,
				strconv.Itoa(int(item.Statistics.Lines)),
				strconv.Itoa(int(item.Statistics.Commits)),
				strconv.Itoa(int(item.Statistics.Files))})
			if err != nil {
				return fmt.Errorf("can't print output in csv format: %s", err)
			}
		}
		w.Flush()
	case jsonFormat:
		entries := make([]string, 0, len(p.statistics))
		for _, item := range p.statistics {
			output := entry{item.Name,
				int(item.Statistics.Lines),
				int(item.Statistics.Commits),
				int(item.Statistics.Files)}
			bytes, err := json.Marshal(output)
			if err != nil {
				return fmt.Errorf("can't print output in json format: %s", err)
			}
			entries = append(entries, string(bytes))
		}
		fmt.Fprintf(outputStream, "[%s]", strings.Join(entries, ","))
	case jsonLinesFormat:
		for _, item := range p.statistics {
			output := entry{item.Name,
				int(item.Statistics.Lines),
				int(item.Statistics.Commits),
				int(item.Statistics.Files)}
			bytes, err := json.Marshal(output)
			if err != nil {
				return fmt.Errorf("can't print output in json-lines format: %s", err)
			}
			bytes = append(bytes, '\n')
			outputStream.Write(bytes)
		}
	default:
	}
	return nil
}
