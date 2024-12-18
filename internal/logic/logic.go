// Package logic contains the business logic of the application.
package logic

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"

	"github.com/idelchi/wslint/internal/wslint"
)

func usage() {
	log.Println("wslint checks or fixes files with trailing whitespaces and enforces final newlines")
	log.Println("Usage: wslint [flags] [path ...]")
	flag.PrintDefaults()
}

// Run is the main function of the application.
func Run(version string) int {
	// Create the Wslint instance
	wslint := wslint.Wslint{Usage: usage, Version: version}
	wslint.Parse()

	if wslint.Options.Experimental {
		if wslint.Options.Fix {
			log.Println(color.YellowString("Experimental feature may not work as expected"))
			log.Println(color.YellowString("Press [enter] to continue or [ctrl+c] to abort"))

			sigCh := make(chan os.Signal, 1)
			enterCh := make(chan struct{})

			// Register to receive SIGINT (Ctrl+C) signals
			signal.Notify(sigCh, syscall.SIGINT)

			// Goroutine to detect the Enter key
			go func() {
				_, _ = os.Stdin.Read([]byte{0})
				enterCh <- struct{}{}
			}()

			select {
			case <-sigCh:
				log.Println("Ctrl+C was pressed. Aborting...")

				return 0
			case <-enterCh:
				log.Println("Enter was pressed. Continuing...")
			}
		}
	}

	if wslint.Match(); len(wslint.Files) == 0 {
		return 1
	}

	return wslint.Process()
}
