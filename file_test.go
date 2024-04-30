package zweb

import "testing"

func TestDownload(t *testing.T) {
	s := NewHTTPServer()
	fd := &FileDownload{Dir: "./testdata/download"}
	s.Get("/download", fd.Handle())
	s.Start(":8081")
}