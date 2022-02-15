package zeuz_node

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
)

// extractTo extracts the file specified by dirEntry into destDir
func extractTo(embeddedFiles embed.FS, destDir string, dirEntry fs.DirEntry) {
	log.Printf("extracting `%s` into `%s` directory", dirEntry.Name(), destDir)
	embeddedFile, err := embeddedFiles.Open(path.Join("embed", dirEntry.Name()))
	if err != nil {
		log.Fatalf("failed to open embedded file: %s", dirEntry.Name())
	}
	defer embeddedFile.Close()

	osFile, err := os.Create(filepath.Join(destDir, dirEntry.Name()))
	if err != nil {
		log.Fatalf("failed to open file: %s", dirEntry.Name())
	}
	defer osFile.Close()

	_, err = io.Copy(osFile, embeddedFile)
	if err != nil {
		log.Fatalf("failed to extract file `%s` to `%s`", dirEntry.Name(), osFile.Name())
	}
}

// ExtractFiles extracts all the embedded files into the "payload" directory.
func ExtractFiles(embeddedFiles embed.FS, payloadDir string) {
	dir := payloadDir

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create `%s` directory to extract required files", dir)
	}

	// extract all the embedded files
	efDirEntries, err := embeddedFiles.ReadDir("embed")
	if err != nil {
		log.Fatal("failed to read embedded files")
	}
	for _, efDirEntry := range efDirEntries {
		extractTo(embeddedFiles, dir, efDirEntry)
	}
}
