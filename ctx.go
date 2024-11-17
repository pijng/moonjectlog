package gomeasure

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
)

var (
	goroutinePrefix = []byte("goroutine ")
	contexts        map[int]context.Context
)

func setContext(ctx context.Context) error {
	if contexts == nil {
		contexts = make(map[int]context.Context)
	}

	id, err := getID()
	if err != nil {
		return fmt.Errorf("failed setting context: %w", err)
	}

	contexts[id] = ctx

	return nil
}

func getContext() (context.Context, error) {
	if contexts == nil {
		contexts = make(map[int]context.Context)
	}

	id, err := getID()
	if err != nil {
		return nil, fmt.Errorf("failed getting context: %w", err)
	}

	ctx, ok := contexts[id]
	if !ok {
		ctx = context.Background()
		err = setContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed creating context: %w", err)
		}
	}

	return ctx, nil
}

func getID() (int, error) {
	buf := make([]byte, 32)
	n := runtime.Stack(buf, false)
	buf = buf[:n]

	buf, ok := bytes.CutPrefix(buf, goroutinePrefix)
	if !ok {
		return 0, errors.New("failed cutting prefix from Stack")
	}

	i := bytes.IndexByte(buf, ' ')
	if i < 0 {
		return 0, errors.New("failed indexing Stack")
	}

	return strconv.Atoi(string(buf[:i]))
}
