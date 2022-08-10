package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/tristanisham/bottle/utils"
)

type GithubReleases []struct {
	Name       string `json:"name"`
	ZipballURL string `json:"zipball_url"`
	TarballURL string `json:"tarball_url"`
	Commit     struct {
		Sha string `json:"sha"`
		URL string `json:"url"`
	} `json:"commit"`
	NodeID string `json:"node_id"`
}

func Upgrade() {
	tags := getRepoTags()
	latest := (*tags)[0]
	color.HiBlue("%s -> %s", utils.VERSION, latest)

	tarball, err := os.CreateTemp("", fmt.Sprintf("bottle-%s.tar.gz", utils.VERSION))
	if err != nil {
		log.Fatal("tmp file: unable to create new temporary file", err)
	}
	defer tarball.Close()

	resp, err := http.Get(latest.TarballURL)
	if err != nil {
		log.Fatalln("downloading new version: unable to download bottle", utils.VERSION, err)
	}
	defer resp.Body.Close()

	if n, err := io.Copy(tarball, resp.Body); err != nil {
		log.Fatal("writing bottle: unable to write new version", err)
	} else {
		color.HiGreen("Wrote %d bytes", n)
	}
	tmpDir, err := os.MkdirTemp("", "bottle-*")
	if err != nil {
		log.Fatal("tmp dir: unable to create tmp dir", err)
	}
	
	if err = utils.Untar(tarball, tmpDir); err != nil {
		log.Fatalln("untar:", err)
	}
}

func getRepoTags() *GithubReleases {
	res, err := http.Get("https://api.github.com/repos/tristanisham/bottle/tags")
	if err != nil {
		log.Fatal("client: could not connect to CDN", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("http response: could not read response body", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatal(res.Status, resBody)
	}

	var data GithubReleases
	if err := json.Unmarshal(resBody, &data); err != nil {
		log.Fatal("unmarshal", err)
	}

	return &data
}

