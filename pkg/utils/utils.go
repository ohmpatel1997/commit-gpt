package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Option struct {
	Cmd     *exec.Cmd
	Verbose bool
}

type OptionFns func(*Option)

func WithVerbose() OptionFns {
	return func(c *Option) {
		c.Verbose = true
	}
}

func WithPath(path string) OptionFns {
	return func(c *Option) {
		c.Cmd.Path = path
	}
}

func WithDir(dir string) OptionFns {
	return func(c *Option) {
		filePath, _ := filepath.Abs(dir)
		c.Cmd.Dir = filePath
	}
}

func WithStdOut(stdOut io.Writer) OptionFns {
	return func(c *Option) {
		c.Cmd.Stdout = stdOut
	}
}

func WithStdErr(stdErr io.Writer) OptionFns {
	return func(c *Option) {
		c.Cmd.Stderr = stdErr
	}
}

func WithStdIn(stdIn io.Reader) OptionFns {
	return func(c *Option) {
		c.Cmd.Stdin = stdIn
	}
}

func WithStdOutOrErr(stdOutOrErr io.Writer) OptionFns {
	return func(c *Option) {
		c.Cmd.Stderr = stdOutOrErr
		c.Cmd.Stdout = stdOutOrErr
	}
}

func WithEnv(lines ...string) OptionFns {
	return func(c *Option) {
		for _, env := range lines {
			c.Cmd.Env = append(c.Cmd.Env, env)
		}
	}
}

func WithArgs(args ...string) OptionFns {
	return func(c *Option) {
		c.Cmd.Args = append([]string{c.Cmd.String()}, args...)
	}
}

func Run(cmd string, options ...OptionFns) error {
	c := exec.Command(cmd)
	c.Env = os.Environ()

	cmdOptions := &Option{
		Cmd: c,
	}

	for _, opt := range options {
		opt(cmdOptions)
	}

	cmd = os.ExpandEnv(cmd)

	for i := range c.Args {
		c.Args[i] = os.ExpandEnv(c.Args[i])
	}

	if cmdOptions.Verbose {
		fmt.Fprintf(c.Stdout, "Exec: %s\n", strings.Join(c.Args, " "))
	}

	if err := c.Run(); err != nil {
		return err
	}

	return nil
}
