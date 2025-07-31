package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/khilk/GitStatsTool/tool/internal/app"
)

var cmd = &cobra.Command{
	Short: "",
	Run:   app.ProcessRepository,
}

func Execute() {
	cmd.Flags().String("repository", ".", "путь до Git репозитория; по умолчанию текущая директория")
	cmd.Flags().String("revision", "HEAD", "указатель на коммит; HEAD по умолчанию")
	cmd.Flags().String("order-by", "lines", "ключ сортировки результатов; один из lines (дефолт), commits, files.")
	cmd.Flags().Bool("use-committer", false, " булев флаг, заменяющий в расчётах автора (дефолт) на коммиттера")
	cmd.Flags().String("format", "tabular", "формат вывода; один из tabular (дефолт), csv, json, json-lines")
	cmd.Flags().String("extensions", "", "список расширений, сужающий список файлов в расчёте; множество ограничений разделяется запятыми, например, '.go,.md'")
	cmd.Flags().String("languages", "", "список языков (программирования, разметки и др.), сужающий список файлов в расчёте; множество ограничений разделяется запятыми, например 'go,markdown'")
	cmd.Flags().String("exclude", "", "набор Glob паттернов, исключающих файлы из расчёта, например 'foo/*,bar/*'")
	cmd.Flags().String("restrict-to", "", "набор Glob паттернов, исключающий все файлы, не удовлетворяющие ни одному из паттернов набора")

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
