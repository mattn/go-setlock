package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	flagndelay, flagx := parseOpt()
	argv := flag.Args()

	if len(argv) < 2 {
		// show usage
		fmt.Fprintf(os.Stderr, "setlock: usage: setlock [ -nNxX ] file program [ arg ... ]\n")
		os.Exit(100)
	}

	filePath := argv[0]
	file, err := os.OpenFile(filePath, syscall.O_WRONLY|syscall.O_NONBLOCK|syscall.O_APPEND|syscall.O_CREAT, 0600) // open append
	if err != nil {
		if flagx {
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "setlock: fatal: unable to open %s: temporary failure\n", filePath)
		os.Exit(111)
	}

	var locker locker
	if flagndelay {
		locker = newLockerEXNB()
	} else {
		locker = newLockerEX()
	}
	err = locker.lock(file)
	if err != nil {
		if flagx {
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "setlock: fatal: unable to lock %s: temporary failure\n", filePath)
		os.Exit(111)
	}

	cmd := exec.Command(argv[1])
	for _, arg := range argv[2:] {
		cmd.Args = append(cmd.Args, arg)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "setlock: fatal: unable to run %s: file does not exist\n", argv[1])
		os.Exit(111)
	}
}

func parseOpt() (bool, bool) {
	var n, N, x, X bool
	flag.BoolVar(&n, "n", false, "No delay. If fn is locked by another process, setlock gives up (windows doesn't support this).")
	flag.BoolVar(&N, "N", false, "(Default.) Delay. If fn is locked by another process, setlock waits until it can obtain a new lock.")
	flag.BoolVar(&x, "x", false, "If fn cannot be opened (or created) or locked, setlock exits zero.")
	flag.BoolVar(&X, "X", false, "(Default.) If fn cannot be opened (or created) or locked, setlock prints an error message and exits nonzero.")
	flag.Parse()

	flagndelay := n && !N
	flagx := x && !X

	return flagndelay, flagx
}