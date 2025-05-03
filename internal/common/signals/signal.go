// Package signals is for stack dump.
package signals

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/mkrtychanr/rag_bot/internal/logger"
)

const bufSize = 1024 * 1024

// SetupStackDump setups a global signal 30 (SIGUSR1) handler
// to print all goroutine stacks to global logger object.
func SetupStackDump() {
	stackDump := make(chan os.Signal, 1)
	signal.Notify(stackDump, syscall.SIGUSR1)

	go func() {
		for {
			<-stackDump

			logger.GetLogger().Error().Msg(PrintStack())
		}
	}()
}

func PrintStack() string {
	buf := make([]byte, bufSize)
	n := runtime.Stack(buf, true)

	return string(buf[:n])
}
