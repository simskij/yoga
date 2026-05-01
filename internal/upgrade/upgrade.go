package upgrade

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

type assetPayload struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type releasePayload struct {
	Assets []assetPayload `json:"assets"`
}

// AssetName returns the expected GitHub release asset name for the current platform.
func AssetName() string {
	return fmt.Sprintf("yo_%s_%s", runtime.GOOS, runtime.GOARCH)
}

// Run fetches the latest release from the given GitHub repo (e.g. "simskij/yoga")
// and atomically replaces the running binary.
func Run(repo string) error {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("locate current binary: %w", err)
	}
	return run(apiURL, execPath, http.DefaultClient)
}

func run(apiURL, targetPath string, client *http.Client) error {
	rel, err := fetchLatest(apiURL, client)
	if err != nil {
		return fmt.Errorf("fetch latest release: %w", err)
	}

	name := AssetName()
	var downloadURL string
	for _, a := range rel.Assets {
		if a.Name == name {
			downloadURL = a.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		return fmt.Errorf("no release asset found for %s", name)
	}

	return downloadAndReplace(downloadURL, targetPath, client)
}

func fetchLatest(apiURL string, client *http.Client) (*releasePayload, error) {
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var rel releasePayload
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, err
	}
	return &rel, nil
}

func downloadAndReplace(url, targetPath string, client *http.Client) error {
	dir := filepath.Dir(targetPath)
	tmp, err := os.CreateTemp(dir, "yo-upgrade-*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	defer func() {
		tmp.Close()
		os.Remove(tmpPath)
	}()

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned %d", resp.StatusCode)
	}

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	tmp.Close()

	if err := os.Chmod(tmpPath, 0o755); err != nil {
		return fmt.Errorf("chmod: %w", err)
	}

	return os.Rename(tmpPath, targetPath)
}
