//go:build windows

package main

import "syscall"

func getWindowSize() (width, height int) {

	const (
		SM_CXSCREEN = uintptr(0) // X Size of screen
		SM_CYSCREEN = uintptr(1) // Y Size of screen
	)
	//w, _, _ := syscall.NewLazyDLL(`User32.dll`).NewProc(`GetSystemMetrics`).Call(SM_CXSCREEN)
	h, _, _ := syscall.NewLazyDLL(`User32.dll`).NewProc(`GetSystemMetrics`).Call(SM_CYSCREEN)

	return int(h * 3 / 2), int(h * 3 / 2)
}
