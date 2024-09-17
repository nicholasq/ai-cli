package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"nicholasq.xyz/ai/internal/ai"
	"nicholasq.xyz/ai/internal/config"
)

var runCmd = &cobra.Command{
	Use:   "run <query>",
	Short: "Run a query on the AI model",
	Long: `Run a query on the AI model with optional flags:
    --cot: Enable chain of thought reasoning
    --no-web: Disable web search capability
    --no-func: Disable function calling capability
    --model: LLM to use, Default=llama3.1
    --debug: Print debugging logs`,
	Example: `  ai run "What is the capital of France?"
  ai run --cot "Explain the process of photosynthesis"
  ai run --no-web --no-func "Calculate 15% of 85"`,
	Run: runQuery,
}

var cfg config.Config

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVar(&cfg.Model, "model", "llama3.1", "LLM to use. Default=llama3.1")
	runCmd.Flags().BoolVar(&cfg.ChainOfThought, "cot", false, "Enable chain of thought")
	runCmd.Flags().BoolVar(&cfg.WebSearch, "no-web", false, "Disable web search")
	runCmd.Flags().BoolVar(&cfg.FunctionCall, "no-func", false, "Disable function calling")
	runCmd.Flags().BoolVar(&cfg.Debug, "debug", false, "Enable debug logging")
}

func runQuery(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		color.Red("Error: Please provide a query")
		cmd.Usage()
		return
	}

	query := args[0]
	client := ai.NewOllamaClient()

	if cfg.Debug {
		color.Blue("Query:\n%s\n---\n", query)
	}

	response, err := client.Query(cmd.Context(), query, cfg)
	if err != nil {
		color.Red("Error executing query: %v\n", err)
		return
	}

	color.Green("AI Response:\n%s\n", response.Text)
}
