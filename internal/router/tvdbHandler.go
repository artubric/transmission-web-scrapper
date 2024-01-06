package router

import (
	"net/http"
	"strconv"
)

func (rt Router) tmdbSearchHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		searchQuery := r.URL.Query().Get("searchQuery")
		result, err := rt.tmdbAPIService.SearchTVShow(searchQuery)
		rt.handleResult(result, err, w)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (rt Router) tmdbGetTVShow(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		showIdString := rt.getPathParamString(r.URL.Path)
		
		showIdInt, err := strconv.Atoi(showIdString)
		if err != nil {
			rt.handleResult(nil, err, w)
		}

		result, err := rt.tmdbAPIService.GetTVShowDetails(showIdInt)
		rt.handleResult(result, err, w)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
