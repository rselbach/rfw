package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

type cmd struct {
	name string
	args []string
	l    *logger

	sync.Mutex
	process *exec.Cmd
}

func newCmd(l *logger, args ...string) (*cmd, error) {
	if len(args) == 0 {
		return nil, errors.New("command is required")
	}
	c := &cmd{name: args[0], l: l}
	if len(args) > 1 {
		c.args = args[1:]
	}

	return c, nil
}

func (c *cmd) restart() {
	c.Lock()
	defer c.Unlock()
	c.l.log("restarting")

	if c.process != nil {
		c.terminate()
	}

	cmd := c.makeCmdLine()
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "could not start command: %s\n", err)
		return
	}
	c.l.logf("why....? %d", cmd.Process.Pid)
	c.process = cmd
	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		c.l.log("task finished")
		c.Lock()
		c.process = nil
		c.Unlock()
	}()
}

func (c *cmd) terminate() {
	if c.process == nil {
		return
	}
	c.l.log("terminating")
	if err := syscall.Kill(-c.process.Process.Pid, syscall.SIGTERM); err != nil {
		fmt.Fprintf(os.Stderr, "could not terminate process: %s", err)
		return
	}
	c.process = nil
}

func (c *cmd) makeCmdLine() *exec.Cmd {
	cmd := exec.Command(c.name, c.args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
