package base

import (
	"errors"
	"fmt"
)

/*
В этом файле представлены кастомные ошибки для базы движка
*/

var (
	ErrNotComparableStruct = errors.New("left argument is a struct, but it doesn't implement Comapre() method")
)

type ErrInvalidType struct {
	expected any
	got      any
}

func NewErrInvalidType(expected, got any) error {
	return ErrInvalidType{expected: expected, got: got}
}

func (e ErrInvalidType) Error() string {
	return fmt.Sprintf("expected %v, but got %v", e.expected, e.got)
}

type ErrUncomparable struct {
	left  any
	right any
}

func NewErrUncomparable(left, right any) error {
	return ErrUncomparable{left: left, right: right}
}

func (e ErrUncomparable) Error() string {
	return fmt.Sprintf("can't compare %v with %v", e.left, e.right)
}

type ErrInvalidPath struct {
	path    string
	details string
}

func NewErrInvalidPath(path, details string) error {
	return ErrInvalidPath{path: path, details: details}
}

func (e ErrInvalidPath) Error() string {
	return fmt.Sprintf("invalid path %v: %v", e.path, e.details)
}

type ErrInexpectedBehavior struct {
	source  string
	details string
}

func NewErrInexpectedBehavior(source, details string) error {
	return ErrInexpectedBehavior{source: source, details: details}
}

func (e ErrInexpectedBehavior) Error() string {
	return fmt.Sprintf("unexpected behavior in %v : %v", e.source, e.details)
}
