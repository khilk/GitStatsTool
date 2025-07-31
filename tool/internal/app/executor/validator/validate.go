package validator

import (
	"fmt"
	"path/filepath"
	"strings"
)

func (v *validator) IsCorrectPath(name string) (bool, error) {
	if len(v.extensions) != 0 {
		tokens := strings.Split(name, ".")
		if len(tokens) < 2 {
			return false, nil
		}
		extension := "." + tokens[len(tokens)-1]
		if _, ok := v.extensions[extension]; !ok {
			return false, nil
		}

	}

	if len(v.restrictTo) != 0 && v.restrictTo[0] != "" {
		matched := false
		for _, p := range v.restrictTo {
			ok, err := filepath.Match(p, name)
			if err != nil {
				return false, fmt.Errorf("failed to match: %s", err)
			}
			if ok {
				matched = true
				break
			}
		}
		if !matched {
			return false, nil
		}
	}

	if len(v.exclude) != 0 {
		for _, p := range v.exclude {
			ok, err := filepath.Match(p, name)
			if err != nil {
				return false, fmt.Errorf("failed to match: %s", err)
			}
			if ok {
				return false, nil
			}
		}
	}

	return true, nil
}
