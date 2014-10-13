package wiki

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"github.com/bmizerany/pat"
	"net/http"
	"time"
)

func init() {
	mux := pat.New()
	mux.Get("/documents", http.HandlerFunc(DocumentsIndex))
	mux.Post("/documents", http.HandlerFunc(DocumentCreate))
	http.Handle("/", mux)
}

type Document struct {
	Key       int64     `datastore:"-" json:"id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DocumentsIndex(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Document")
	var documents []*Document
	key, err := q.GetAll(c, &documents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for i := range documents {
		documents[i].Key = key[i].IntID()
	}
	b, err := json.Marshal(documents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(documents) == 0 {
		b = []byte("[]")
	}
	fmt.Fprintf(w, "%s", b)
}

func DocumentCreate(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var document Document
	document.Name = r.FormValue("name")
	document.CreatedAt = time.Now()
	document.UpdatedAt = time.Now()
	ds := datastore.NewIncompleteKey(c, "Document", nil)
	key, err := datastore.Put(c, ds, &document)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	document.Key = key.IntID()
	b, err := json.Marshal(document)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", b)
}
