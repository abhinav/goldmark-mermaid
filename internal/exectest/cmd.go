// Package exectest provides a means of mocking [os/exec.Cmd]s
// allowing injection of arbitrary behavior into an external executable
// from a test.
package exectest

import (
	"context"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// Adapted from
// https://github.com/uber-go/fx/blob/fd5fd36f83c2ac06a1bf2e366766afb0c9fdb2b9/docs/internal/exectest/cmd.go.

// Actor defines a fake behavior (or act) for an external executable.
// Build these with [Act].
type Actor struct {
	testExe  string
	testName string
}

// Act builds an actor that acts as specified in the given function.
//
// Because of how this operates:
//
//   - call this near the top of a test function
//   - avoid non-determistic code before this function is called
//   - do not call this inside subtests
func Act(t testing.TB, main func()) *Actor {
	t.Helper()

	// This messes up the hijacking sometimes.
	// Keep it simple -- only top level tests can do this.
	require.NotContains(t, t.Name(), "/",
		"exectest.Act cannot be used with subtests")

	// We can't get coverage for this block
	// because if the condition is true,
	// we're inside the subprocess.
	if filepath.Base(os.Args[0]) == t.Name() {
		// After test argument parsing,
		// flag.Args holds the arguments
		// that were originally passed to Actor.Command.
		os.Args = append(os.Args[:1], flag.Args()...)
		main()
		os.Exit(0)
	}

	exe, err := os.Executable()
	require.NoError(t, err, "determine executable")

	return &Actor{
		testExe:  exe,
		testName: t.Name(),
	}
}

// CommandContext builds an exec.Cmd that will run this Actor as an external
// executable with the provided arguments.
//
// This operates by re-running the test executable to run only the current
// test, and hijacking that test execution to run the main function.
//
//	actor := exectest.Act(t, func() { fmt.Println("hello") })
//	cmd := actor.CommandContext(ctx, args)
//	got, err := cmd.Output()
//	...
//	fmt.Println(string(got) == "hello\n") // true
func (c *Actor) CommandContext(ctx context.Context, args ...string) *exec.Cmd {
	testArgs := []string{"-test.run", "^" + c.testName + "$"}
	if len(args) > 0 {
		testArgs = append(testArgs, "--")
		testArgs = append(testArgs, args...)
	}
	cmd := exec.CommandContext(ctx, c.testExe, testArgs...)

	// Args[0] is the value of os.Args[0] for the new executable.
	// os.Args[0] is allowed to be different from the command.
	cmd.Args[0] = c.testName
	return cmd
}
