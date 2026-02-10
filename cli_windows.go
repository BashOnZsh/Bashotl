package main

import (
	"syscall"
	"unsafe"
)

// Structure pour définir les coordonnées de la console
type coord struct {
	X, Y int16
}

// Structure pour définir un rectangle de console
type smallRect struct {
	Left, Top, Right, Bottom int16
}

// Agrandir la fenêtre de console pour mieux afficher l'ASCII art
func ResizeConsoleWindow() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	
	// Obtenir le handle de la console
	getStdHandle := kernel32.NewProc("GetStdHandle")
	setConsoleScreenBufferSize := kernel32.NewProc("SetConsoleScreenBufferSize")
	setConsoleWindowInfo := kernel32.NewProc("SetConsoleWindowInfo")
	
	// Handle de la console de sortie standard
	stdOutputHandle := uintptr(^uint32(10) + 1) // STD_OUTPUT_HANDLE = -11
	handle, _, _ := getStdHandle.Call(stdOutputHandle)
	
	if handle != 0 {
		// Définir une taille de buffer plus grande (120 colonnes x 40 lignes)
		bufferSize := coord{X: 120, Y: 40}
		setConsoleScreenBufferSize.Call(handle, uintptr(*(*int32)(unsafe.Pointer(&bufferSize))))
		
		// Définir la taille de la fenêtre (120 colonnes x 35 lignes visibles)
		windowSize := smallRect{Left: 0, Top: 0, Right: 119, Bottom: 34}
		setConsoleWindowInfo.Call(handle, 1, uintptr(unsafe.Pointer(&windowSize)))
	}
}
	
func IsDoubleClickRun() bool {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	lp := kernel32.NewProc("GetConsoleProcessList")
	if lp != nil {
		var pids [2]uint32
		var maxCount uint32 = 2
		ret, _, _ := lp.Call(uintptr(unsafe.Pointer(&pids)), uintptr(maxCount))
		if ret > 1 {
			return false
		}
	}
	return true
}
