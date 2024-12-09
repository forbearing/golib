package filetype

import (
	"testing"
)

var (
	documentFiles = []string{
		"../../testdata/filetype/sample-documents/sample.doc",
		"../../testdata/filetype/sample-documents/sample.xls",
		"../../testdata/filetype/sample-documents/sample.ppt",
		"../../testdata/filetype/sample-documents/sample.docx",
		"../../testdata/filetype/sample-documents/sample.xlsx",
		"../../testdata/filetype/sample-documents/sample.pptx",
		"../../testdata/filetype/sample-documents/sample.odt",
		"../../testdata/filetype/sample-documents/sample.ods",
		"../../testdata/filetype/sample-documents/sample.odp",
		"../../testdata/filetype/sample-documents/sample.pdf",
		"../../testdata/filetype/sample-documents/sample.rtf",
	}
	textFiles = []string{
		"../../testdata/filetype/sample-text/sample.css",
		"../../testdata/filetype/sample-text/sample.csv",
		"../../testdata/filetype/sample-text/sample.html",
		"../../testdata/filetype/sample-text/sample.js",
		"../../testdata/filetype/sample-text/sample.json",
		"../../testdata/filetype/sample-text/sample.php",
		"../../testdata/filetype/sample-text/sample.sh",
		"../../testdata/filetype/sample-text/sample.txt",
		"../../testdata/filetype/sample-text/sample.xml",
		"../../testdata/filetype/sample-text/sample.yml",
	}
	compressFiles = []string{
		"../../testdata/filetype/sample-compress/sample.zip",
		"../../testdata/filetype/sample-compress/sample.tar",
		"../../testdata/filetype/sample-compress/sample.z",
		"../../testdata/filetype/sample-compress/sample.7z",
		"../../testdata/filetype/sample-compress/sample.gz",
		"../../testdata/filetype/sample-compress/sample.lz",
		"../../testdata/filetype/sample-compress/sample.xz",
		"../../testdata/filetype/sample-compress/sample.bz2",
		"../../testdata/filetype/sample-compress/sample.rar",
		"../../testdata/filetype/sample-compress/sample.zst",
		"../../testdata/filetype/sample-compress/sample.lzma", // unknow
		"../../testdata/filetype/sample-compress/sample.lzop", // unknow
	}
	imageFiles = []string{
		"../../testdata/filetype/sample-images/sample.gif",
		"../../testdata/filetype/sample-images/sample.ico",
		"../../testdata/filetype/sample-images/sample.jpg",
		"../../testdata/filetype/sample-images/sample.png",
		"../../testdata/filetype/sample-images/sample.svg",
		"../../testdata/filetype/sample-images/sample.tiff",
		"../../testdata/filetype/sample-images/sample.webp",
	}
	videoFiles = []string{
		"../../testdata/filetype/sample-videos/sample.avi",
		"../../testdata/filetype/sample-videos/sample.mov",
		"../../testdata/filetype/sample-videos/sample.mp4",
		"../../testdata/filetype/sample-videos/sample.ogg",
		"../../testdata/filetype/sample-videos/sample.webm",
		"../../testdata/filetype/sample-videos/sample.wmv",
	}
	audoFiles = []string{
		"../../testdata/filetype/sample-audio/sample.mp3",
		"../../testdata/filetype/sample-audio/sample.ogg",
		"../../testdata/filetype/sample-audio/sample.wav",
	}
	otherFiles = []string{
		"../../testdata/filetype/sample-others/sample.elf",
		"../../testdata/filetype/sample-others/sample.exe",
		"../../testdata/filetype/sample-others/sample.macho",
		"../../testdata/filetype/sample-others/sample.iso",
		"../../testdata/filetype/sample-others/sample.jar",
	}
)

func TestDetectFiletype(t *testing.T) {
	// Documents
	for _, file := range documentFiles {
		filetype, mime := Detect(file)
		t.Logf("%5s  %s\n", filetype, mime)
	}
	t.Log()

	// Text/Plains
	for _, file := range textFiles {
		filetype, mime := Detect(file)
		t.Logf("%5s  %s\n", filetype, mime)
	}
	t.Log()

	// compress
	for _, file := range compressFiles {
		filetype, mime := Detect(file)
		t.Logf("%5s  %s\n", filetype, mime)
	}
	t.Log()

	// Images
	for _, file := range imageFiles {
		filetype, mime := Detect(file)
		t.Logf("%5s  %s\n", filetype, mime)
	}
	t.Log()

	// Videos
	for _, file := range videoFiles {
		filetype, mime := Detect(file)
		t.Logf("%5s  %s\n", filetype, mime)
	}
	t.Log()

	// Audio
	for _, file := range audoFiles {
		filetype, mime := Detect(file)
		t.Logf("%5s  %s\n", filetype, mime)
	}
	t.Log()

	// others
	for _, file := range otherFiles {
		filetype, mime := Detect(file)
		t.Logf("%5s  %s\n", filetype, mime)
	}
}
