/*
Considered hashed filename schemes:
	Example                         Problem
	main-extra.min.js?v=h45h57r1n6  Query string may be ignored by shitty caches.
	main-extra.min.js/h45h57r1n6    Cannot detect content type by extension and sort by file name.
	main-extra.min.js.h45h57r1n6    Cannot detect content type by extension.
	h45h57r1n6-main-extra.min.js    Cannot sort by file name.
	main-extra-h45h57r1n6.min.js    Cannot sort by complete file name.
	main-extra.min-h45h57r1n6.js    Ugly.
	main-extra.min.h45h57r1n6.js    -
*/

package static

import (
	"os"

	"github.com/gowww/crypto"
)

// isHash tells if s represents a hash.
func isHash(s string) bool {
	if len(s) != 32 {
		return false
	}
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b < '0' || b > '9' && b < 'a' || b > 'z' {
			return false
		}
	}
	return true
}

// hashSplitFilepath returns the part before the hash, the hash, and the filename extension (with dot).
func hashSplitFilepath(path string) (prefix, hash, ext string) {
	prefix = path
	if prefix[len(prefix)-1] == '/' { // Path ends with slash: no hash.
		return
	}
	extSep := extDotIndex(prefix)
	if extSep == -1 { // No dot in last part: no hash or extension.
		return
	}
	ext = prefix[extSep:]    // Put extension (with dot) in its return variable.
	prefix = prefix[:extSep] // Strip extension from prefix.
	hashSep := extDotIndex(prefix)
	if hashSep == -1 { // A single dot in base: see if extension is in fact a hash.
		extWithoutDot := ext[1:]
		if isHash(extWithoutDot) {
			hash, ext = extWithoutDot, ""
		}
		return
	}
	preHash := prefix[hashSep+1:] // Last filename part (without dot) could be a hash.
	if isHash(preHash) {
		hash = preHash            // Put hash in its return variable.
		prefix = prefix[:hashSep] // Strip hash from prefix.
	}
	return
}

func fileHash(path string) (string, error) {
	cfi := hashCacheData.Get(path)
	if cfi != nil { // Check modification time to validate cached hash.
		fi, err := os.Stat(path)
		if err != nil {
			hashCacheData.Del(path)
			return "", err
		}
		if !fi.ModTime().After(cfi.ModTime) {
			return cfi.Hash, nil
		}
		hashCacheData.Del(path)
	}

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	hash, err := crypto.HashMD5(f)
	if err != nil {
		return "", err
	}
	stat, err := f.Stat()
	if err != nil {
		return "", err
	}
	hashCacheData.Set(path, hash, stat.ModTime())
	return hash, nil
}

// extDotIndex returns the last dot index in the last path part.
// It none, -1 is returned.
func extDotIndex(path string) int {
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			return i
		}
	}
	return -1
}
