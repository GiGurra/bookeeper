package github

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
)

type Release struct {
	URL             string    `json:"url"`
	AssetsURL       string    `json:"assets_url"`
	UploadURL       string    `json:"upload_url"`
	HTMLURL         string    `json:"html_url"`
	ID              int64     `json:"id"`
	Author          User      `json:"author"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       string    `json:"created_at"`
	PublishedAt     string    `json:"published_at"`
	Assets          []Asset   `json:"assets"`
	TarballURL      string    `json:"tarball_url"`
	ZipballURL      string    `json:"zipball_url"`
	Body            string    `json:"body"`
	Reactions       Reactions `json:"reactions"`
}

type User struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	UserViewType      string `json:"user_view_type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Asset struct {
	URL                string `json:"url"`
	ID                 int64  `json:"id"`
	NodeID             string `json:"node_id"`
	Name               string `json:"name"`
	Label              any    `json:"label"`
	Uploader           User   `json:"uploader"`
	ContentType        string `json:"content_type"`
	State              string `json:"state"`
	Size               int    `json:"size"`
	DownloadCount      int    `json:"download_count"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func (a Asset) DownloadToDir(dir string) string {
	// download the asset to the specified directory
	resp, err := http.Get(a.BrowserDownloadURL)
	if err != nil {
		panic(fmt.Errorf("failed to download asset: %w", err))
	}
	defer func() { _ = resp.Body.Close() }()
	defer func() { _, _ = io.Copy(io.Discard, resp.Body) }()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("failed to download asset: %s", resp.Status))
	}

	contDisp := resp.Header.Get("Content-Disposition")
	fileName := func() string {
		if contDisp != "" {
			_, params, err := mime.ParseMediaType(contDisp)
			if err != nil {
				panic(fmt.Errorf("failed to parse content disposition: %w", err))
			}
			return params["filename"]
		}
		panic("no content disposition header found")
	}()

	filePath := dir + "/" + fileName

	// write the asset to the specified directory
	outFile, err := os.Create(filePath)
	if err != nil {
		panic(fmt.Errorf("failed to create file: %w", err))
	}
	defer func() { _ = outFile.Close() }()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		panic(fmt.Errorf("failed to write asset to file: %w", err))
	}

	return filePath
}

type Reactions struct {
	URL        string `json:"url"`
	TotalCount int    `json:"total_count"`
	PlusOne    int    `json:"+1"`
	MinusOne   int    `json:"-1"`
	Laugh      int    `json:"laugh"`
	Hooray     int    `json:"hooray"`
	Confused   int    `json:"confused"`
	Heart      int    `json:"heart"`
	Rocket     int    `json:"rocket"`
	Eyes       int    `json:"eyes"`
}

func GetLatestRelease(owner string, repo string) Release {
	resp, err := http.Get("https://api.github.com/repos/" + owner + "/" + repo + "/releases/latest")
	if err != nil {
		panic(fmt.Errorf("failed to get latest release: %w", err))
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
	}()

	// decode the json blob into a Release struct
	var release Release
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		panic(fmt.Errorf("failed to decode latest release: %w", err))
	}

	return release
}
