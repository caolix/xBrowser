//go:build linux

package main

import (
	"os/exec"
	"strconv"
	"strings"
)

func getWindowSize() (width, height int) {
	// Try to get screen size using xdpyinfo (common on X11)
	cmd := exec.Command("xdpyinfo")
	out, err := cmd.Output()
	if err == nil {
		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(line, "dimensions:") {
				parts := strings.Fields(line)
				for i, part := range parts {
					if part == "dimensions:" && i+1 < len(parts) {
						res := strings.Split(parts[i+1], "x")
						if len(res) == 2 {
							w, _ := strconv.Atoi(res[0])
							h, _ := strconv.Atoi(res[1])
							if w > 0 && h > 0 {
								return w * 93 / 100, h * 93 / 100
							}
						}
					}
				}
			}
		}
	}

	// Fallback to default size
	return 1920, 1340
}

func openDirCmd(dir string) (string, []string) {
	return "xdg-open", []string{dir}
}

func openTextFileCmd(text string) (string, []string) {
	return "xdg-open", []string{text}
}
