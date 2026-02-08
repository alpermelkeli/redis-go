package commands

import (
	"redis-go/internal/router"
	"strings"
)

func Register(r *router.Router) {
	r.Register("PING", func(args []string) (string, error) {
		return "PONG", nil
	})

	r.Register("ECHO", func(args []string) (string, error) {
		return strings.Join(args, " "), nil
	})

	r.Register("UPPER", func(args []string) (string, error) {
		return strings.ToUpper(strings.Join(args, " ")), nil
	})
}