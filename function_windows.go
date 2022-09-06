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

	w := h * 4 / 3
	return int(w * 93 / 100), int(h * 93 / 100)
}

func openDirCmd(dir string) (string, []string) {
	return "start", []string{dir}
}

func openTextFileCmd(text string) (string, []string) {
	return "notepad", []string{text}
}
