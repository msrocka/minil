package main

import (
	"path/filepath"
	"strings"
)

func outputFile(origin string, extension string) string {
	outPath := filepath.Base(origin)
	// outPath = outPath + "_" + time.Now().Format(time.RFC822) + extension
	outPath += extension
	return strings.ReplaceAll(strings.ReplaceAll(outPath, ":", "_"), " ", "_")
}
