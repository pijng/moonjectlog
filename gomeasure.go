package gomeasure

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Span struct {
	Name     string
	Start    time.Time
	Children []*Span
	Parent   *Span // Add Parent field
	Ctx      context.Context
	mu       sync.Mutex
}

type traceKeyType string

const traceKey traceKeyType = "trace"

// NewSpan creates a new span
func newSpan(ctx context.Context, name string, parent *Span) *Span {
	return &Span{
		Name:     name,
		Start:    time.Now(),
		Children: make([]*Span, 0),
		Parent:   parent,
		Ctx:      ctx,
	}
}

// AddChild adds a child span to the current span
func (s *Span) addChild(child *Span) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Children = append(s.Children, child)
}

// Print outputs the span and its children with indentation
func (s *Span) print(indent int) {
	fmt.Printf("%s%s - %v\n", strings.Repeat("  ", indent), s.Name, time.Since(s.Start))
	for _, child := range s.Children {
		child.print(indent + 1)
	}
	if indent == 0 {
		fmt.Println("------")
	}
}

// StartSpan creates a new span and attaches it to the parent span in the context
func StartSpan(name string) *Span {
	ctx, err := getContext()
	if err != nil {
		panic(err)
	}

	parentSpan, _ := ctx.Value(traceKey).(*Span)
	newSpan := newSpan(ctx, name, parentSpan)
	if parentSpan != nil {
		newSpan.Ctx = parentSpan.Ctx
		parentSpan.addChild(newSpan)
	} else {
		newSpan.Ctx = ctx
	}

	ctx = context.WithValue(ctx, traceKey, newSpan)
	err = setContext(ctx)
	if err != nil {
		panic(err)
	}

	return newSpan
}

func (s *Span) End() {
	if s.Parent != nil {
		return
	}

	s.print(0)
}
