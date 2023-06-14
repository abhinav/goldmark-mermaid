package exectest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandSuccess(t *testing.T) {
	t.Parallel()

	cmd := Act(t, func() {
		fmt.Println("hello world")
	}).Command()

	out, err := cmd.Output()
	require.NoError(t, err)
	assert.Equal(t, "hello world\n", string(out))

	assert.True(t, cmd.ProcessState.Exited(), "must exit")
	assert.Zero(t, cmd.ProcessState.ExitCode(), "exit code")
}

func TestCommandNonZero(t *testing.T) {
	t.Parallel()

	cmd := Act(t, func() {
		fmt.Fprintln(os.Stderr, "great sadness")
		os.Exit(1)
	}).Command()

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	require.Error(t, err, "command must fail")

	assert.Equal(t, "great sadness\n", stderr.String())
	assert.True(t, cmd.ProcessState.Exited(), "must exit")
	assert.Equal(t, 1, cmd.ProcessState.ExitCode(), "exit code")
}

func TestCommandArgs(t *testing.T) {
	t.Parallel()

	actor := Act(t, func() {
		enc := json.NewEncoder(os.Stdout)
		if err := enc.Encode(os.Args[1:]); err != nil {
			log.Fatal(err)
		}
	})

	args := []string{
		"foo", "-bar",
		// Randomly generated argument to ensure that code
		"-random", strconv.Itoa(rand.Int()),
	}
	cmd := actor.Command(args...)

	out, err := cmd.Output()
	require.NoError(t, err)

	var got []string
	require.NoError(t, json.Unmarshal(out, &got))
	assert.Equal(t, args, got)
}
