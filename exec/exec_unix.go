// +build !windows

package exec

// No initialization required.
func platformInit() error {
	return nil
}

// Run the application until terminated by a signal.
func Exec() error {
	execSignal()
	return nil
}

// No cleanup necessary.
func Cleanup() {
}
