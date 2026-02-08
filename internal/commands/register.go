package commands

import (
	"errors"
	"redis-go/internal/router"
	"redis-go/internal/store"
	"strconv"
)

const ARGUMENT_ERROR string = "ERR wrong number of arguments for "

func Register(r *router.Router, s *store.Store) {
	r.Register("PING", func(args []string) (string, error) {
		return "PONG", nil
	})

	r.Register("GET", func(args []string) (string, error) {
		if len(args) != 1 {
			return "", errors.New(ARGUMENT_ERROR + "GET")
		}

		val, ok := s.Get(args[0])
		if !ok {
			return "(nil)", nil
		}

		return val, nil
	})

	r.Register("SET", func(args []string) (string, error) {
		if len(args) != 2 {
			return "", errors.New(ARGUMENT_ERROR + "SET")
		}

		s.Set(args[0], args[1])
		return "OK", nil
	})

	r.Register("DELETE", func(args []string) (string, error) {
		if len(args) != 1 {
			return "", errors.New(ARGUMENT_ERROR + "DELETE")
		}
		ok := s.Delete(args[0])
		return strconv.FormatBool(ok), nil
	})
}
