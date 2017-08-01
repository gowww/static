// Package static provides a handler for static file serving with cache control and automatic fingerprinting.
package static

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// A Handler serves files and provides helpers.
type Handler interface {
	http.Handler
	Hash(string) string
}

type handler struct {
	next   http.Handler
	prefix string
	dir    string
}

// Handle returns a handler for static file serving.
func Handle(prefix, dir string) Handler {
	if !strings.HasPrefix(prefix, "/") || !strings.HasSuffix(prefix, "/") {
		panic(fmt.Errorf("static: prefix %q must begin and end with %q", prefix, "/"))
	}
	// TODO: Cache all file hashes from dir recursively in a goroutine.
	return &handler{
		next:   http.FileServer(http.Dir(dir)),
		prefix: prefix,
		dir:    dir,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, h.prefix)
	prefix, reqHash, ext := hashSplitFilepath(r.URL.Path)
	if reqHash == "" { // No hash: serve file as wanted.
		h.next.ServeHTTP(w, r)
		return
	}
	realHash, err := fileHash(filepath.Join(h.dir, prefix+ext))
	if err != nil { // Cannot open file to get real hash.
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}
	if reqHash != realHash { // Hash has changed: redirect to new one.
		http.Redirect(w, r, h.prefix+prefix+"."+realHash+ext, http.StatusMovedPermanently)
		return
	}
	r.URL.Path = prefix + ext
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Header().Set("ETag", `"`+realHash+`"`)
	h.next.ServeHTTP(w, r)
}

// Hash returns the URL path from a prefix and a file path.
// If the file was successfully opened, the file hash is appended to the file name.
// dir is the directory where the file is to be found.
func (h *handler) Hash(path string) string {
	path = strings.TrimPrefix(path, "/") // Avoid double "/" as h.prefix already ends with one.
	hash, err := fileHash(filepath.Join(h.dir, path))
	if err != nil {
		return h.prefix + path
	}
	extDotIdx := extDotIndex(path)
	if extDotIdx == -1 {
		return h.prefix + path + "." + hash
	}
	return h.prefix + path[:extDotIdx] + "." + hash + path[extDotIdx:]
}

// toHTTPError is copied from net/http/fs.go.
func toHTTPError(err error) (msg string, httpStatus int) {
	if os.IsNotExist(err) {
		return "404 page not found", http.StatusNotFound
	}
	if os.IsPermission(err) {
		return "403 Forbidden", http.StatusForbidden
	}
	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}
