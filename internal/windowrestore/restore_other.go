//go:build !windows

package windowrestore

// Configure is a no-op on platforms that do not need the Windows workaround.
func Configure() func() { return nil }
