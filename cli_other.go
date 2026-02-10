//go:build !windows

package main

func IsDoubleClickRun() bool {
	return false
}

// Fonction stub pour les systèmes non-Windows
func ResizeConsoleWindow() {
	// Ne fait rien sur les systèmes non-Windows
}
