package wiki

import (
	"appengine"
	"appengine/aetest"
	"appengine/datastore"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestDocumentCreate(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	defer inst.Close()
	req, err := inst.NewRequest("POST", "/documents", nil)
	if err != nil {
		t.Fatalf("Failed to create req: %v", err)
	}
	w := httptest.NewRecorder()
	DocumentCreate(w, req)
	var document *Document
	err = json.Unmarshal(w.Body.Bytes(), &document)
	if err != nil {
		t.Fatalf("Failed to unmarshal document: %v", err)
	}
	if document.Key == 0 {
		t.Error("Returned document doesn't have key set")
	}
}

func TestDocumentsIndex(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	defer inst.Close()
	req, err := inst.NewRequest("GET", "/documents", nil)
	if err != nil {
		t.Fatalf("Failed to create req: %v", err)
	}
	w := httptest.NewRecorder()
	DocumentsIndex(w, req)

	c := appengine.NewContext(req)
	var document Document
	document.Name = "New document"
	ds := datastore.NewIncompleteKey(c, "Document", nil)
	_, err = datastore.Put(c, ds, &document)
	if err != nil {
		t.Fatal(err)
	}

	var documents []*Document
	err = json.Unmarshal(w.Body.Bytes(), &documents)
	if err != nil {
		t.Fatalf("Failed to unmarshal documents: %v", err)
	}
	if len(documents) != 1 {
		t.Errorf("Expected to get 1 document, but got %v", len(documents))
	}
}
