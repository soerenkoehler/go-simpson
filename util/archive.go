package util

import (
	"archive/tar"
	"archive/zip"
	"compress/flate"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type archiveBase interface {
	CreateEntry(string, os.FileInfo) (io.Writer, error)
	Close() error
}

type archiveZip struct {
	*zip.Writer
}

func (archive archiveZip) CreateEntry(file string, info os.FileInfo) (io.Writer, error) {
	return archive.Create(file)
}

type archiveTgz struct {
	*tar.Writer
	gzWriter *gzip.Writer
}

func (archive archiveTgz) CreateEntry(file string, info os.FileInfo) (io.Writer, error) {
	header, err := tar.FileInfoHeader(info, info.Name())
	if err == nil {
		return archive.Writer, archive.WriteHeader(header)
	}
	return nil, err
}

func (archive archiveTgz) Close() error {
	if err := archive.Writer.Close(); err != nil {
		return err
	}
	return archive.gzWriter.Close()
}

type ArchiveType struct {
	Creator   func(io.Writer) archiveBase
	Extension string
}

// supported archive types
var (
	ZIP = ArchiveType{
		func(w io.Writer) archiveBase {
			zipWriter := zip.NewWriter(w)
			zipWriter.RegisterCompressor(
				zip.Deflate,
				func(w io.Writer) (io.WriteCloser, error) {
					return flate.NewWriter(w, flate.BestCompression)
				})
			return &archiveZip{zipWriter}
		},
		"zip"}
	TGZ = ArchiveType{
		func(w io.Writer) archiveBase {
			if gzWriter, err := gzip.NewWriterLevel(w, gzip.BestCompression); err == nil {
				return &archiveTgz{
					tar.NewWriter(gzWriter),
					gzWriter}
			}
			return nil
		},
		"tgz"}
)

// Creates an archive of the given type with the content from contentPath.
//
// Return values: archivePath, archiveHash, error
//
func CreateArchive(
	archiveType ArchiveType,
	contentPath string) (string, string, error) {

	archiveFilePath := fmt.Sprintf(
		"%v.%v",
		contentPath,
		archiveType.Extension)
	archiveFile, err := os.Create(archiveFilePath)
	checksum := sha256.New()
	output := io.MultiWriter(archiveFile, checksum)

	if err == nil {

		defer func() {
			archiveFile.Close()
		}()

		archive := archiveType.Creator(output)

		addFiles(archive, contentPath)

		archive.Close()
	}

	return archiveFilePath, hex.EncodeToString(checksum.Sum(nil)), err
}

func addFiles(archive archiveBase, contentPath string) error {
	for file, info := range walkTree(contentPath) {
		if entry, err := archive.CreateEntry(file, info); err == nil {
			data, err := os.Open(filepath.Join(contentPath, file))
			if err == nil {
				io.Copy(entry, data)
				data.Close()
			}
		}
	}

	return nil
}

func walkTree(root string) map[string]os.FileInfo {
	result := make(map[string]os.FileInfo, 0x100)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if !info.IsDir() {
				relPath, _ := filepath.Rel(root, path)
				result[relPath] = info
			}
		}
		return nil
	})
	return result
}
