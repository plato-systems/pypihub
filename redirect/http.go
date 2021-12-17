package redirect

import "net/http"

func HandleHTTP(w http.ResponseWriter, r *http.Request) {
	dest := table[r.URL.Path]
	if dest == "" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, dest, http.StatusTemporaryRedirect)
}
