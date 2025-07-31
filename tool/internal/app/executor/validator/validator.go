package validator

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

type Validator interface {
	IsCorrectPath(filepath string) (bool, error)
	Init(commandFlags *pflag.FlagSet) error
}

type validator struct {
	extensions map[string]struct{}
	exclude    []string
	restrictTo []string
}

func New() Validator {
	return &validator{
		extensions: make(map[string]struct{}),
	}
}

func (v *validator) Init(commandFlags *pflag.FlagSet) error {
	extensions, err := func(str string, err error) ([]string, error) {
		return strings.Split(str, ","), err
	}(commandFlags.GetString("extensions"))
	if err != nil {
		return fmt.Errorf("can't parse extensions flag: %s", err)
	}

	languages, err := func(str string, err error) ([]string, error) {
		return strings.Split(str, ","), err
	}(commandFlags.GetString("languages"))
	if err != nil {
		return fmt.Errorf("can't parse languages flag: %s", err)
	}

	err = v.setExtensions(languages, extensions)
	if err != nil {
		return fmt.Errorf("can't set extensions: %s", err)
	}

	v.exclude, err = func(str string, err error) ([]string, error) {
		return strings.Split(str, ","), err
	}(commandFlags.GetString("exclude"))
	if err != nil {
		return fmt.Errorf("can't parse exclude flag: %s", err)
	}

	v.restrictTo, err = func(str string, err error) ([]string, error) {
		return strings.Split(str, ","), err
	}(commandFlags.GetString("restrict-to"))
	if err != nil {
		return fmt.Errorf("can't parse restrict-to flag: %s", err)
	}

	return nil
}

//go:embed configs/language_extensions.json
var languageExtensions []byte

func (v *validator) convertLanguagesToExtensions(languages []string) (map[string]struct{}, error) {
	const configLength = 396
	res := make(map[string]struct{}, len(languages))
	languagesToExtensions := make([]struct {
		Name       string   `json:"name"`
		Extensions []string `json:"extensions"`
	}, 0, configLength)

	err := json.Unmarshal(languageExtensions, &languagesToExtensions)
	if err != nil {
		return nil, fmt.Errorf("can't read language-extensions config: %s", err)
	}

	languageToExtensions := make(map[string][]string, len(languagesToExtensions))
	for _, entry := range languagesToExtensions {
		name := strings.ToUpper(entry.Name)
		languageToExtensions[name] = append(languageToExtensions[name], entry.Extensions...)
	}

	for _, l := range languages {
		l = strings.ToUpper(l)
		extensions := languageToExtensions[l]
		for _, e := range extensions {
			res[e] = struct{}{}
		}
	}

	return res, nil
}

func (v *validator) setExtensions(languages []string, extensions []string) error {
	fromLanguages, err := v.convertLanguagesToExtensions(languages)
	if err != nil {
		return fmt.Errorf("can't set extensons: %s", err)
	}

	res := make(map[string]struct{}, len(extensions)+len(languages))
	for _, e := range extensions {
		res[e] = struct{}{}
	}
	for k := range fromLanguages {
		res[k] = struct{}{}
	}

	delete(res, "")

	v.extensions = res

	return nil
}
