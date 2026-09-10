//go:build !windows

package main

func configureWindowRestore() func() { return nil }
