package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/docker/go-units"
)

const megabyte = 1024 * 1024

func main() {
	var amount string
	flag.StringVar(&amount, "amount", "100M", "Amount of memory to allocate. This is parsed assuming binary prefixes, therefore 1K = 1024 = 2^10. Examples: 1K 1M 1G")

	var help bool
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&help, "h", false, "Show help")

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	targetMemory, err := units.RAMInBytes(amount)
	if err != nil {
		fmt.Printf("Invalid amount specified: %s", err)
		os.Exit(1)
	}

	fmt.Println("Starting simpler memory allocation stress test...")

	// Start goroutine to print memory stats every 5 seconds
	go func() {
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Printf(
				"Stats: HeapSys: %v MB, HeapAlloc: %v MB, HeapIdle: %v MB, HeapReleased: %v MB\n",
				m.HeapSys/megabyte,
				m.HeapAlloc/megabyte,
				m.HeapIdle/megabyte,
				m.HeapReleased/megabyte,
			)
			time.Sleep(5 * time.Second)
		}
	}()

	// Allocate a single large slice of bytes
	largeBuffer := make([]byte, targetMemory)

	// To ensure the memory is actually touched and not optimized away by the compiler
	// (though Go's make typically ensures this for non-zero sized allocations),
	// we can write to each byte in the slice.
	//
	// Setting only largeBuffer[targetAmount-1] = 1 is not sufficient to actually allocate
	// the whole memory.
	fmt.Println("Touching all allocated memory to force physical commit...")
	for i := int64(0); i < targetMemory; i++ {
		largeBuffer[i] = byte(i % 256) // Write a value to each byte
	}

	fmt.Printf("Successfully allocated approximately %d MB of memory.\n", targetMemory/megabyte)

	// Loop endless to keep the program running
	var i int64
	for {
		// Interact with the allocated RAM to avoid any unexpected deallocations
		i += 1
		i %= targetMemory
		largeBuffer[i] = byte(i % 256)

		time.Sleep(50 * time.Millisecond)
	}
}
