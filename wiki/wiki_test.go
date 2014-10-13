package wiki

import (
	"appengine/aetest"
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
		t.Errorf("Returned document doesn't have key set")
	}
}
