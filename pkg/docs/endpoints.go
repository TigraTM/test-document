package docs

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func NewNewsHandler(subRouter *mux.Router, docsSvc Service) http.Handler {
	subRouter.HandleFunc("/document/set", CreateDocumentHandler(docsSvc)).Methods(http.MethodPost)
	subRouter.HandleFunc("/document/get", GetDocumentsHandler(docsSvc)).Methods(http.MethodGet)

	return subRouter
}

func CreateDocumentHandler(docsSvc Service) func(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Documents []Document `json:"documents"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errorHandler(w, r, http.StatusBadRequest, err)
			return
		}

		if err := docsSvc.CreateDocuments(r.Context(), req.Documents); err != nil {
			errorHandler(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		respond(w, r, http.StatusCreated, nil)
	}
}

func GetDocumentsHandler(docsSvc Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		documents, err := docsSvc.GetDocuments(r.Context())
		if err != nil {
			errorHandler(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		respond(w, r, http.StatusOK, documents)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, code int, err error) {
	respond(w, r, code, map[string]string{"error": err.Error()})
}

func respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			return
		}
	}
}

