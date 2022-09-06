//go:build darwin

package main

func getWindowSize() (width, height int) {
	return 1920, 1340
}

func openDirCmd(dir string) (string, []string) {
	return "open", []string{dir}
}

func openTextFileCmd(text string) (string, []string) {
	return "open", []string{"-e", text}
}
