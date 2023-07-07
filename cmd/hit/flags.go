package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const usageText = `
Usage:
	hit [options] url
Options:`

type flags struct {
	url  string
	n, c int
}

// number is a natural number.
type number int

// toNumber is a convenience function for converting p to *number.
func toNumber(p *int) *number {
	return (*number)(p)
}
func (n *number) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	switch {
	case err != nil:
		err = errors.New("parse error")
	case v <= 0:
		err = errors.New("should be positive")
	}
	*n = number(v)
	return err
}

func (n *number) String() string {
	return strconv.Itoa(int(*n))
}

type parseFunc func(string) error

func (f *flags) parse(s *flag.FlagSet, args []string) (err error) {
	s.Usage = func() {
		fmt.Fprintln(s.Output(), usageText[1:])
		s.PrintDefaults()
	}
	s.Var(toNumber(&f.n), "n", "Number of requests to make")
	s.Var(toNumber(&f.c), "c", "Concurrency level")
	if err := s.Parse(args); err != nil {
		return err
	}
	f.url = s.Arg(0)
	if err := f.validate(s); err != nil {
		fmt.Fprintln(s.Output(), err)
		s.Usage()
		return err
	}
	return nil
}

func (f *flags) urlVar(p *string) parseFunc {
	return func(s string) error {
		_, err := url.Parse(s)
		*p = s
		return err
	}
}

func (f *flags) intVar(p *int) parseFunc {
	return func(s string) (err error) {
		*p, err = strconv.Atoi(s)
		return err
	}
}

func (f *flags) validate(s *flag.FlagSet) error {
	if err := validateURL(s.Arg(0)); err != nil {
		return fmt.Errorf("url: %w", err)
	}
	if f.c > f.n {
		return fmt.Errorf("-c=%d: should be less than or equal to -n=%d", f.c, f.n)
	}
	return nil
}

func validateURL(s string) error {
	u, err := url.Parse(s)
	switch {
	case strings.TrimSpace(s) == "":
		err = errors.New("required")
	case err != nil:
		err = errors.New("parse error")
	case u.Scheme != "http":
		err = errors.New("only supported scheme is http")
	case u.Host == "":
		err = errors.New("missing host")
	}
	return err
}
