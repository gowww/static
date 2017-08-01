package static

import (
	"os"
	"regexp"
	"testing"
)

var testsFilenameBaseHash = []struct {
	base string
	want string
}{
	{"one.edbc4f9728a8e311b55e081b27e3caff", "edbc4f9728a8e311b55e081b27e3caff"},
	{"one.edbc4f9728a8e311", ""},
	{"one.edbc4f(728a8e&11b55!081b$7e3caff", ""},
}

type argsHashSplitFilepath struct{ prefix, hash, ext string }

var testsHashSplitFilepath = []struct {
	path string
	want argsHashSplitFilepath
}{
	{"main", argsHashSplitFilepath{
		"main",
		"",
		"",
	}}, {"main.js", argsHashSplitFilepath{
		"main",
		"",
		".js",
	}}, {"main.js.", argsHashSplitFilepath{
		"main.js",
		"",
		".",
	}}, {"main.edbc4f9728a8e311b55e081b27e3caff.", argsHashSplitFilepath{
		"main",
		"edbc4f9728a8e311b55e081b27e3caff",
		".",
	}}, {"main.min.js", argsHashSplitFilepath{
		"main.min",
		"",
		".js",
	}}, {"main.min.edbc4f9728a8e311b55e081b27e3caff.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", argsHashSplitFilepath{
		"main.min",
		"edbc4f9728a8e311b55e081b27e3caff",
		".xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
	}}, {"edbc4f9728a8e311b55e081b27e3caff", argsHashSplitFilepath{
		"edbc4f9728a8e311b55e081b27e3caff",
		"",
		"",
	}}, {"edbc4f9728a8e311b55e081b27e3caff.js", argsHashSplitFilepath{
		"edbc4f9728a8e311b55e081b27e3caff",
		"",
		".js",
	}}, {"edbc4f9728a8e311b55e081b27e3caff.min.js", argsHashSplitFilepath{
		"edbc4f9728a8e311b55e081b27e3caff.min",
		"",
		".js",
	}}, {"main-extra.min.edbc4f9728a8e311b55e081b27e3caff", argsHashSplitFilepath{
		"main-extra.min",
		"",
		".edbc4f9728a8e311b55e081b27e3caff",
	}}, {"main-extra.min.edbc4f9728a8e311b55e081b27e3caff.js", argsHashSplitFilepath{
		"main-extra.min",
		"edbc4f9728a8e311b55e081b27e3caff",
		".js",
	}}, {"/styles/main-extra.min.edbc4f9728a8e311b55e081b27e3caff.js", argsHashSplitFilepath{
		"/styles/main-extra.min",
		"edbc4f9728a8e311b55e081b27e3caff",
		".js",
	}}, {"main.edbc4f9728a8e311.js", argsHashSplitFilepath{
		"main.edbc4f9728a8e311",
		"",
		".js",
	}}, {"main.edbc4f(728a8e&11b55!081b$7e3caff.js", argsHashSplitFilepath{
		"main.edbc4f(728a8e&11b55!081b$7e3caff",
		"",
		".js",
	}}, {"main.edbc4f9728a8e311b55e081b27e3caff.min.js", argsHashSplitFilepath{
		"main.edbc4f9728a8e311b55e081b27e3caff.min",
		"",
		".js",
	}},
}

var testsFileHash = []struct {
	path string
	want string
}{
	{"LICENSE", "edbc4f9728a8e311b55e081b27e3caff"},
	{"unknown", ""},
}

func TestHashSplitFilepath(t *testing.T) {
	for _, tt := range testsHashSplitFilepath {
		prefix, hash, ext := hashSplitFilepath(tt.path)
		if prefix != tt.want.prefix || hash != tt.want.hash || ext != tt.want.ext {
			t.Errorf("%q:\nwant: %q %q %q\ngot:  %q %q %q", tt.path, tt.want.prefix, tt.want.hash, tt.want.ext, prefix, hash, ext)
		}
	}
}

func TestFileHash(t *testing.T) {
	for _, tt := range testsFileHash {
		got, err := fileHash(tt.path)
		if err != nil && !os.IsNotExist(err) {
			panic(err)
		}
		if got != tt.want {
			t.Errorf("%q: want %q, got %q", tt.path, tt.want, got)
		}
	}
}

func BenchmarkHashSplitFilepath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tt := range testsHashSplitFilepath {
			hashSplitFilepath(tt.path)
		}
	}
}

func BenchmarkIsHashRegexp(b *testing.B) {
	re := regexp.MustCompile(`^[a-f0-9]{32}$`)
	for i := 0; i < b.N; i++ {
		re.FindStringSubmatch("edbc4f9728a8e311b55e081b27e3caff")
	}
}

func BenchmarkIsHashCustom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isHash("edbc4f9728a8e311b55e081b27e3caff")
	}
}
