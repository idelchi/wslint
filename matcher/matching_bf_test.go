package matcher_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bou.ke/monkey"
)

// Test exists only for the purpose of getting familiar with monkey patching.
func TestExecutableError(t *testing.T) { //nolint:paralleltest // Monkey patching is not thread-safe
	// Make life easier - assume there's an environment variable called 'RUN_MONKEY_PATCH_TESTS'
	if ok, err := strconv.ParseBool(os.Getenv("RUN_MONKEY_PATCH_TESTS")); !ok || err != nil {
		t.Skip("Skipping test because it requires monkey patching (inlining disabled)")
	}

	defer monkey.UnpatchAll()
	monkey.Patch(os.Executable, func() (string, error) {
		return "", assert.AnError
	})

	_, err := os.Executable()

	require.Error(t, err)
}
