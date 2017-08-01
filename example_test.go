package static_test

import (
	"fmt"
	"net/http"

	"github.com/gowww/static"
)

func Example() {
	staticHandler := static.Handle("/static/", "static")

	http.Handle("/static/", staticHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Cacheable asset: %s", staticHandler.Hash("main.js"))
	})

	http.ListenAndServe(":8080", nil)
}

func ExampleHandler_Hash() {
	staticHandler := static.Handle("/static/", ".")

	// File exists: hash will be appended to the file name.
	fmt.Println(staticHandler.Hash("LICENSE"))

	// File doesn't exist: the file name will contain no hash.
	fmt.Println(staticHandler.Hash("unknown"))

	// Output:
	// /static/LICENSE.edbc4f9728a8e311b55e081b27e3caff
	// /static/unknown
}
