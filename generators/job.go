package generators

import (
	"fmt"
	"net/http"
	"time"
)

// PrettyJob wraps a plaintext job with a pretty HTML shell.
func PrettyJob(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		RenderAdminHeader(w)

		title := fmt.Sprintf("Query: %s", r.URL.Path)
		fmt.Fprintf(w, "<h1>%s</h1>", title)
		fmt.Fprintf(w, "<title>%s</title>", title)

		start := time.Now()
		fmt.Fprintf(w, "<p>Started at: %s</p>", start)

		w.Write([]byte(`<script>
		window.scrollerInterval = setInterval(()=>window.scrollTo(0,document.body.offsetHeight),100);
		</script>`))

		w.Write([]byte("<pre>"))

		// Periodically flush the ResponseWriter so it renders on the client in
		// chunks.
		if f, ok := w.(http.Flusher); ok {
			f.Flush()

			t := time.NewTicker(300 * time.Millisecond)
			defer t.Stop()

			go func() {
				for range t.C {
					f.Flush()
				}
			}()
		}

		f(w, r)
		w.Write([]byte("</pre>"))
		fmt.Fprintf(w, "<p>Time taken: %s, %s</p>", time.Since(start), time.Now())
		fmt.Fprintf(w, `<p><a href="%s">Run Again</a></p>`, r.URL.Path)
		w.Write([]byte(`<script>
		window.scrollTo(0,document.body.offsetHeight);
		clearInterval(window.scrollerInterval);
		</script>`))
	}
}

// RenderAdminHeader adds the admin header to the stream.
func RenderAdminHeader(w http.ResponseWriter) {
	if err := Templates.ExecuteTemplate(w, "adminhead.html", nil); err != nil {
		HandleErr(w, err)
		return
	}
}

// HandleErr handles an error and renders it to the writer.
func HandleErr(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Sprintf("%+v", err), 500)
}
