package router

import (
	"fmt"
	"strings"
)

type CommandHandler func(args []string) (string, error)

type Router struct {
	handlers map[string]CommandHandler
}

func New() *Router {
	return &Router{
		handlers: make(map[string]CommandHandler),
	}
}

func (r *Router) Register(cmd string, handler CommandHandler) {
	r.handlers[strings.ToUpper(cmd)] = handler
}

func (r *Router) Handle(msg string) (string, error) {
	parts := strings.Fields(msg)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	cmd := strings.ToUpper(parts[0])
	args := parts[1:]

	handler, ok := r.handlers[cmd]
	if !ok {
		return "", fmt.Errorf("unknown command: %s", cmd)
	}

	return handler(args)
}