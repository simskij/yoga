package upgrade

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAssetName(t *testing.T) {
	got := AssetName()
	want := fmt.Sprintf("yo_%s_%s", runtime.GOOS, runtime.GOARCH)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRun_DownloadsAndReplacesTarget(t *testing.T) {
	newContent := []byte("#!/bin/sh\necho upgraded\n")
	assetName := AssetName()

	mux := http.NewServeMux()
	mux.HandleFunc("/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		rel := releasePayload{Assets: []assetPayload{{
			Name:               assetName,
			BrowserDownloadURL: "http://" + r.Host + "/download/" + assetName,
		}}}
		json.NewEncoder(w).Encode(rel)
	})
	mux.HandleFunc("/download/"+assetName, func(w http.ResponseWriter, r *http.Request) {
		w.Write(newContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	targetPath := filepath.Join(t.TempDir(), "yo")
	if err := os.WriteFile(targetPath, []byte("old"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := run(srv.URL+"/releases/latest", targetPath, srv.Client()); err != nil {
		t.Fatalf("run: %v", err)
	}

	got, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(newContent) {
		t.Errorf("binary content not replaced: got %q", got)
	}
	info, err := os.Stat(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&0o111 == 0 {
		t.Error("replaced binary is not executable")
	}
}

func TestRun_NoMatchingAsset(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		rel := releasePayload{Assets: []assetPayload{{
			Name:               "yo_plan9_mips",
			BrowserDownloadURL: "http://" + r.Host + "/download/yo_plan9_mips",
		}}}
		json.NewEncoder(w).Encode(rel)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	targetPath := filepath.Join(t.TempDir(), "yo")
	err := run(srv.URL+"/releases/latest", targetPath, srv.Client())
	if err == nil {
		t.Fatal("expected error when no matching asset, got nil")
	}
}

func TestRun_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer srv.Close()

	targetPath := filepath.Join(t.TempDir(), "yo")
	err := run(srv.URL+"/releases/latest", targetPath, srv.Client())
	if err == nil {
		t.Fatal("expected error on API failure, got nil")
	}
}
