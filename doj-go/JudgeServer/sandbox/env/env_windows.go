package env

import (
	"doj-go/JudgeServer/sandbox/env/pool"
	"doj-go/JudgeServer/sandbox/env/winc"
)

// NewBuilder build a environment builder
func NewBuilder(c Config) (pool.EnvBuilder, map[string]any, error) {
	b, err := winc.NewBuilder("")
	if err != nil {
		return nil, nil, err
	}
	c.Info("created winc builder")
	return b, map[string]any{}, nil
}
