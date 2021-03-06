package bintest_test

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/buildkite/bintest/v3"
)

func ExampleCompileProxy() {
	// create a proxy for the git command that echos some debug
	p, err := bintest.CompileProxy("git")
	if err != nil {
		log.Fatal(err)
	}

	// call the proxy like a normal binary in the background
	cmd := exec.Command(p.Path, "rev-parse")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// windows needs all the environment variables
	cmd.Env = append(os.Environ(), `MY_MESSAGE=Llama party! 🎉`)

	if err := cmd.Start(); err != nil {
		_ = p.Close()
		log.Fatal(err)
	}

	// handle invocations of the proxy binary
	call := <-p.Ch
	fmt.Fprintln(call.Stdout, call.GetEnv(`MY_MESSAGE`))
	call.Exit(0)

	// wait for the command to finish
	_ = cmd.Wait()

	if err := p.Close(); err != nil {
		log.Fatal(err)
	}

	// Output: Llama party! 🎉
}
