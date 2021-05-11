package argoapp

import "testing"

func Test_NewCommand(t *testing.T) {
	// Make sure NewCommand doesn't panic. This may happen when there are
	// typos in the required flags.
	_ = NewCommand()
}
