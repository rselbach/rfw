package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

func main() {
	var (
		watchPaths sliceFlag
		verbose    bool
	)
	flag.Var(&watchPaths, "path", "path(s) to watch; can be specified multiple times")
	flag.BoolVar(&verbose, "v", false, "print details about what is going on")

	flag.Parse()
	logger := logger(verbose)

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "missing command")
		os.Exit(-2)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	cmd, err := newCmd(&logger, flag.Args()...)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {

			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op == fsnotify.Chmod {
					continue
				}

				logger.logf("change detected: %s", event.Name)

				cmd.restart()
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.log("error:", err)
			}
		}
	}()

	if len(watchPaths) == 0 {
		watchPaths = []string{"."}
	}
	for _, p := range watchPaths {
		logger.logf("Watching for changes in %s", p)
		err := watcher.Add(p)
		if err != nil {
			panic(err)
		}
	}

	// run the command for the first time
	cmd.restart()

	<-c
	cmd.terminate()
	logger.log("Bye") // dsfsdd
}
