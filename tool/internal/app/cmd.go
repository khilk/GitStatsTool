package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/khilk/GitStatsTool/tool/internal/app/executor"
	"github.com/khilk/GitStatsTool/tool/internal/app/printer"
)

func ProcessRepository(cmd *cobra.Command, args []string) {
	errorStream := os.Stderr

	executor, err := inintializeExecutor(cmd)
	if err != nil {
		fmt.Fprintln(errorStream, fmt.Errorf("can't inintialize gitfame executor: %s", err))
		os.Exit(1)
	}

	printer, err := inintializePrinter(cmd)
	if err != nil {
		fmt.Fprintln(errorStream, fmt.Errorf("can't inintialize printer: %s", err))
		os.Exit(1)
	}

	err = executor.Execute()
	if err != nil {
		fmt.Fprintln(errorStream, fmt.Errorf("failed to execute gitfame: %s", err))
		os.Exit(1)
	}
	results := executor.GetResult()

	err = printer.PrintStatistics(results)
	if err != nil {
		fmt.Fprintln(errorStream, fmt.Errorf("can't print statistics: %s", err))
		os.Exit(1)
	}
}

func inintializeExecutor(cmd *cobra.Command) (executor.GitfameExecutor, error) {
	executor := executor.New()
	err := executor.Init(cmd.Flags())
	if err != nil {
		return nil, fmt.Errorf("can't init executor: %s", err)
	}

	return executor, nil
}

func inintializePrinter(cmd *cobra.Command) (printer.Printer, error) {
	printer, err := printer.New(cmd.Flags())
	if err != nil {
		return nil, fmt.Errorf("can't create printer: %s", err)
	}

	err = printer.ValidateConfig()
	if err != nil {
		return nil, fmt.Errorf("printer has invalid config: %s", err)
	}

	return printer, nil
}
