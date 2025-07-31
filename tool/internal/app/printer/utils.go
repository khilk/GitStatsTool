package printer

import "fmt"

func (p *printer) ValidateConfig() error {
	if p.orderBy != "lines" &&
		p.orderBy != "commits" &&
		p.orderBy != "files" {
		return fmt.Errorf("invalid value for orderBy flag: %s", p.orderBy)
	}

	if p.format != "tabular" &&
		p.format != "csv" &&
		p.format != "json" &&
		p.format != "json-lines" {
		return fmt.Errorf("invalid value for format flag: %s", p.format)
	}

	return nil
}
