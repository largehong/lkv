package command

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

type Command struct {
	name    string
	args    []string
	envs    []string
	workdir string
	timeout int64
}

func New(name string, args ...string) *Command {
	return &Command{
		name: name,
		args: args,
	}
}

func (c *Command) AddArgs(args ...string) *Command {
	c.args = append(c.args, args...)
	return c
}

func (c *Command) AddEnvs(envs ...string) *Command {
	c.envs = append(c.envs, envs...)
	return c
}

func (c *Command) SetWorkDirecotry(dir string) *Command {
	c.workdir = dir
	return c
}

func (c *Command) SetTimeout(timeout int64) *Command {
	c.timeout = timeout
	return c
}

func (c *Command) String() string {
	if len(c.args) == 0 {
		return c.name
	}
	return fmt.Sprintf("%s %s", c.name, strings.Join(c.args, " "))
}

func (c *Command) RunWithPipe(stdout, stderr io.Writer) (err error) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if c.timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(c.timeout)*time.Second)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	cmd := exec.CommandContext(ctx, c.name, c.args...)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	cmd.Env = c.envs
	cmd.Dir = c.workdir
	return cmd.Run()
}

func (c *Command) Run() (output string, err error) {
	buf := &bytes.Buffer{}
	err = c.RunWithPipe(buf, buf)
	return buf.String(), err
}
