package commands

import (
	"errors"
	"redis-go/internal/router"
	"redis-go/internal/store"
)

func Register(r *router.Router, s *store.Store) {
	r.Register("PING", func(args []string) (string, error) {
		return "PONG", nil
	})

	r.Register("GET", func(args []string) (string, error) {
		if len(args) != 1 {
			return "", errors.New("ERR wrong number of arguments for GET")
		}

		val, ok := s.Get(args[0])
		if !ok {
			return "(nil)", nil
		}

		return val, nil
	})

	r.Register("SET", func(args []string) (string, error) {
		if len(args) != 2 {
			return "", errors.New("ERR wrong number of arguments for SET")
		}

		s.Set(args[0], args[1])
		return "OK", nil
	})
}
