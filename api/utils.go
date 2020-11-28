package api

import (
	"context"
	"io"
	"os"
)

const (
	PluginsDir    = "./plugins"
	CmdSymbolName = "Commands"
	DefaultPrompt = "xsh>"
)

func GetStdout(ctx context.Context) io.Writer {
	var out io.Writer = os.Stdout
	if ctx == nil {
		return out
	}
	if outVal := ctx.Value("xsh.stdout"); outVal != nil {
		if stdout, ok := outVal.(io.Writer); ok {
			out = stdout
		}
	}
	return out
}

func GetPrompt(ctx context.Context) string {
	prompt := DefaultPrompt
	if ctx == nil {
		return prompt
	}
	if promptVal := ctx.Value("xsh.prompt"); promptVal != nil {
		if p, ok := promptVal.(string); ok {
			prompt = p
		}
	}
	return prompt
}
