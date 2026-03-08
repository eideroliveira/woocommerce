package woocommerce

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	filesBasePath = "download"
)

// FileService is an interface for interfacing with the file endpoints of
// the WooCommerce files restful API
// https://woocommerce.github.io/woocommerce-rest-api-docs/#files
type FileService interface {
	Get(file string) (*File, error)
	GetStream(file string) (*FileDownload, error)
	GetMeta(file string) (*FileMeta, error)
}

// FileMeta contains file metadata returned by the download-meta endpoint.
type FileMeta struct {
	Filename     string `json:"filename"`
	Size         int64  `json:"size"`
	LastModified int64  `json:"last_modified"`
}

// FileServiceOp handles communication with the files related methods of WooCommerce restful api
type FileServiceOp struct {
	Client *Client
}

// File represent a  wooCommerce file's All  properties columns
type File struct {
	Name    string `json:"filename,omitempty"`
	Content []byte `json:"content,omitempty"`
}

// FileDownload wraps a temporary file containing downloaded content.
// The caller must call Close() to remove the temp file when done.
type FileDownload struct {
	Name    string
	tmpFile *os.File
}

func (fd *FileDownload) Read(p []byte) (int, error) {
	return fd.tmpFile.Read(p)
}

func (fd *FileDownload) Seek(offset int64, whence int) (int64, error) {
	return fd.tmpFile.Seek(offset, whence)
}

func (fd *FileDownload) ReadAt(p []byte, off int64) (int, error) {
	return fd.tmpFile.ReadAt(p, off)
}

// File returns the underlying *os.File for use with APIs that accept it
// (e.g., base.Base.Scan). The caller should NOT close this file directly;
// use FileDownload.Close() instead.
func (fd *FileDownload) File() *os.File {
	return fd.tmpFile
}

// Close closes and removes the temporary file.
func (fd *FileDownload) Close() error {
	name := fd.tmpFile.Name()
	fd.tmpFile.Close()
	return os.Remove(name)
}

// Get retrieves a file using the JSON API. The entire response body is buffered
// in memory. For large files, use GetStream instead.
func (w *FileServiceOp) Get(file string) (*File, error) {
	path := fmt.Sprintf("%s/%s", filesBasePath, file)
	resource := new(File)
	// Use createAndDoGetHeaders to access response headers
	headers, err := w.Client.createAndDoGetHeaders("GET", path, nil, nil, &resource)

	if err == nil {
		w.Client.log.Infof("FileServiceOp.Get success: file=%s, size=%d, headers=%v", file, len(resource.Content), headers)
	}

	return resource, err
}

// GetStream downloads a file by streaming the HTTP response body directly to a
// temporary file on disk. This avoids buffering the entire file in memory.
// Caller must call Close() on the returned FileDownload to clean up.
func (w *FileServiceOp) GetStream(file string) (*FileDownload, error) {
	relPath := fmt.Sprintf("%s/%s", filesBasePath, file)

	req, err := w.Client.NewAPIRequest("GET", relPath, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	w.Client.logRequest(req)

	resp, err := w.Client.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		errBody := make([]byte, 1024)
		n, _ := io.ReadFull(resp.Body, errBody)
		return nil, fmt.Errorf("http %d downloading %s: %s", resp.StatusCode, file, string(errBody[:n]))
	}

	// The download endpoint returns JSON: {"filename":"...", "content":"<base64>"}
	// We stream-decode to avoid buffering the entire file in memory.
	dec := json.NewDecoder(resp.Body)

	// Expect opening {
	if _, err := dec.Token(); err != nil {
		return nil, fmt.Errorf("expected JSON object start: %w", err)
	}

	var filename string
	var tmp *os.File
	var written int64

	for dec.More() {
		// Read field name
		tok, err := dec.Token()
		if err != nil {
			return nil, fmt.Errorf("reading JSON key: %w", err)
		}
		key, ok := tok.(string)
		if !ok {
			continue
		}

		switch key {
		case "filename":
			var name string
			if err := dec.Decode(&name); err != nil {
				return nil, fmt.Errorf("decoding filename: %w", err)
			}
			filename = name

		case "content":
			// Stream-decode the base64 value directly to a temp file.
			// Read the raw base64 string token, decode in chunks.
			tok, err := dec.Token()
			if err != nil {
				return nil, fmt.Errorf("reading content token: %w", err)
			}
			b64str, ok := tok.(string)
			if !ok {
				return nil, fmt.Errorf("expected string for content, got %T", tok)
			}

			tmp, err = os.CreateTemp("", "wc-dl-*")
			if err != nil {
				return nil, fmt.Errorf("creating temp file: %w", err)
			}

			// Decode the base64 string in chunks to the temp file
			b64Reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64str))
			written, err = io.Copy(tmp, b64Reader)
			if err != nil {
				tmp.Close()
				os.Remove(tmp.Name())
				return nil, fmt.Errorf("decoding base64 content: %w", err)
			}

		default:
			// Skip unknown fields
			var discard json.RawMessage
			if err := dec.Decode(&discard); err != nil {
				return nil, fmt.Errorf("skipping field %s: %w", key, err)
			}
		}
	}

	if tmp == nil {
		return nil, fmt.Errorf("no content field found in download response for %s", file)
	}

	if filename == "" {
		filename = filepath.Base(file)
	}

	if _, err := tmp.Seek(0, io.SeekStart); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return nil, fmt.Errorf("seeking temp file: %w", err)
	}

	w.Client.log.Infof("FileServiceOp.GetStream: file=%s, size=%d, tmp=%s", filename, written, tmp.Name())

	return &FileDownload{Name: filename, tmpFile: tmp}, nil
}

// GetMeta retrieves file metadata (size, filename, last_modified) without
// downloading the file content. Uses the download-meta endpoint.
func (w *FileServiceOp) GetMeta(file string) (*FileMeta, error) {
	path := fmt.Sprintf("download-meta/%s", file)
	resource := new(FileMeta)
	_, err := w.Client.createAndDoGetHeaders("GET", path, nil, nil, resource)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// splitDisposition splits a Content-Disposition header value by semicolons and trims whitespace.
func splitDisposition(s string) []string {
	var parts []string
	for _, p := range splitSemicolon(s) {
		p = trimSpace(p)
		if p != "" {
			parts = append(parts, p)
		}
	}
	return parts
}

func splitSemicolon(s string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ';' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimSpace(s string) string {
	for len(s) > 0 && s[0] == ' ' {
		s = s[1:]
	}
	for len(s) > 0 && s[len(s)-1] == ' ' {
		s = s[:len(s)-1]
	}
	return s
}

func trimQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
