package config

import "fmt"

type Config struct {
	Model          string
	WebSearch      bool `default:"true"`
	FunctionCall   bool `default:"true"`
	ChainOfThought bool
	Verbose        bool
	Debug          bool
	OutputFormat   string
}

func (c Config) String() string {
	return fmt.Sprintf("Config{Model: %s, WebSearch: %t, FunctionCall: %t, ChainOfThought: %t, Verbose: %t, OutputFormat: %s}",
		c.Model, c.WebSearch, c.FunctionCall, c.ChainOfThought, c.Verbose, c.OutputFormat)
}
