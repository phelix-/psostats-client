// Entrypoint for PSOStats Client
package main

import (
	"github.com/phelix-/psostats/v2/pkg/model"
	"log"
	"os"
	"syscall"

	"github.com/phelix-/psostats/v2/client/internal/client"
)

var (
	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	procSetStdHandle = kernel32.MustFindProc("SetStdHandle")
)

func main() {
	file, err := os.OpenFile("psostats.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	version := model.ClientInfo{
		VersionMajor: 1,
		VersionMinor: 4,
		VersionPatch: 0,
	}

	log.SetOutput(file)
	// Redirecting this app's stderr to psostats.log so we can get debug info from panics
	redirectStderr(file)
	log.Printf("Starting Up version %v", version)

	c, err := client.New(version)
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}
	c.Run()
}

func setStdHandle(stdhandle int32, handle syscall.Handle) error {
	r0, _, e1 := syscall.Syscall(procSetStdHandle.Addr(), 2, uintptr(stdhandle), uintptr(handle), 0)
	if r0 == 0 {
		if e1 != 0 {
			return error(e1)
		}
		return syscall.EINVAL
	}
	return nil
}

// redirectStderr to the file passed in
func redirectStderr(f *os.File) {
	err := setStdHandle(syscall.STD_ERROR_HANDLE, syscall.Handle(f.Fd()))
	if err != nil {
		log.Fatalf("Failed to redirect stderr to file: %v", err)
	}
	// SetStdHandle does not affect prior references to stderr
	os.Stderr = f
}