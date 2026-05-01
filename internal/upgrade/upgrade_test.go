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

func makeServer(t *testing.T, tag, assetName, assetContent string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		rel := releasePayload{
			TagName: tag,
			Assets: []assetPayload{{
				Name:               assetName,
				BrowserDownloadURL: "http://" + r.Host + "/download/" + assetName,
			}},
		}
		json.NewEncoder(w).Encode(rel)
	})
	mux.HandleFunc("/download/"+assetName, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(assetContent))
	})
	return httptest.NewServer(mux)
}

func TestRun_DownloadsAndReplacesTarget(t *testing.T) {
	newContent := "#!/bin/sh\necho upgraded\n"
	assetName := AssetName()
	srv := makeServer(t, "v2.0.0", assetName, newContent)
	defer srv.Close()

	targetPath := filepath.Join(t.TempDir(), "yo")
	if err := os.WriteFile(targetPath, []byte("old"), 0o755); err != nil {
		t.Fatal(err)
	}

	upgraded, err := run(srv.URL+"/releases/latest", targetPath, "v1.0.0", srv.Client())
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !upgraded {
		t.Error("expected upgraded=true")
	}

	got, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != newContent {
		t.Errorf("binary content not replaced: got %q", got)
	}
	if info, _ := os.Stat(targetPath); info.Mode()&0o111 == 0 {
		t.Error("replaced binary is not executable")
	}
}

func TestRun_AlreadyLatest_SkipsDownload(t *testing.T) {
	assetName := AssetName()
	srv := makeServer(t, "v1.0.0", assetName, "new content")
	defer srv.Close()

	targetPath := filepath.Join(t.TempDir(), "yo")
	original := []byte("original")
	if err := os.WriteFile(targetPath, original, 0o755); err != nil {
		t.Fatal(err)
	}

	upgraded, err := run(srv.URL+"/releases/latest", targetPath, "v1.0.0", srv.Client())
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if upgraded {
		t.Error("expected upgraded=false when already on latest")
	}

	got, _ := os.ReadFile(targetPath)
	if string(got) != string(original) {
		t.Error("binary was replaced when it should have been skipped")
	}
}

func TestRun_DevVersion_AlwaysUpgrades(t *testing.T) {
	assetName := AssetName()
	srv := makeServer(t, "v1.0.0", assetName, "new content")
	defer srv.Close()

	targetPath := filepath.Join(t.TempDir(), "yo")
	if err := os.WriteFile(targetPath, []byte("old"), 0o755); err != nil {
		t.Fatal(err)
	}

	upgraded, err := run(srv.URL+"/releases/latest", targetPath, "dev", srv.Client())
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !upgraded {
		t.Error("expected upgraded=true for dev builds")
	}
}

func TestRun_NoMatchingAsset(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		rel := releasePayload{TagName: "v1.0.0", Assets: []assetPayload{{
			Name:               "yo_plan9_mips",
			BrowserDownloadURL: "http://" + r.Host + "/download/yo_plan9_mips",
		}}}
		json.NewEncoder(w).Encode(rel)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	targetPath := filepath.Join(t.TempDir(), "yo")
	_, err := run(srv.URL+"/releases/latest", targetPath, "dev", srv.Client())
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
	_, err := run(srv.URL+"/releases/latest", targetPath, "dev", srv.Client())
	if err == nil {
		t.Fatal("expected error on API failure, got nil")
	}
}
