package ghost

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var BadArgsErr = errors.New("insufficient arguments")

func DefaultCli(args []string, fs *flag.FlagSet) (chan *Event, error) {
	if fs == nil {
		fs = flag.NewFlagSet("ghostrace", flag.ExitOnError)
		fs.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage: %s [options] -p <pid> | <exe> [args...]\n", args[0])
			fs.PrintDefaults()
		}
	}
	follow := fs.Bool("f", false, "follow subprocesses")
	pid := fs.Int("p", -1, "attach to pid")
	fs.Parse(args[1:])
	args = fs.Args()

	var trace chan *Event
	var err error
	tracer := NewTracer()
	if pid != nil && *pid >= 0 {
		trace, err = tracer.Trace(*pid)
	} else {
		if len(args) > 0 {
			trace, err = tracer.Spawn(args[0], args...)
		} else {
			fs.Usage()
			return nil, BadArgsErr
		}
	}
	if err != nil {
		return nil, fmt.Errorf("Error starting trace: %s", err)
	}
	tracer.ExecFilter(func(c *Event) (bool, bool) {
		// fmt.Println("exec filter", c)
		// keepParent, followChild
		return true, *follow
	})
	return trace, nil
}
