package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"runtime"
)

type meteredCommand func(cmd *cobra.Command, args []string)

func withMetrics(f meteredCommand) meteredCommand {

	return func(cmd *cobra.Command, args []string) {
		storeMetrics, _ := cmd.Flags().GetBool("debug")
		if storeMetrics {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
			log.Printf("localStore\tStart\tCurrent RAM usage: %d MiB\n\n", bToMb(m.HeapInuse)+bToMb(m.StackInuse))
		}
		f(cmd, args)
		if storeMetrics {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
			log.Printf("localStore\tEnd\tCurrent RAM usage: %d MiB\n", bToMb(m.HeapInuse)+bToMb(m.StackInuse))

			runtime.GC()
			runtime.ReadMemStats(&m)
			log.Printf("localStore\tEnd\tCurrent RAM usage(after garbage collection): %d MiB\n", bToMb(m.HeapInuse)+bToMb(m.StackInuse))
		}
	}
}
