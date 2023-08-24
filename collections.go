package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/Masterminds/semver"
	"github.com/cavaliergopher/grab/v3"
	"github.com/dihedron/ansible-galaxy-grabber/helpers"
	"github.com/dihedron/rawdata"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/mattn/go-isatty"
	"github.com/pterm/pterm"
)

type Collections struct {
	Entries []Collection `json:"collections" yaml:"collections"`
}

func (c *Collections) UnmarshalFlag(value string) error {
	tmp := Collections{}
	*c = tmp
	return rawdata.UnmarshalInto(value, &c.Entries)
}

type Collection struct {
	Namespace   string  `json:"namespace" yaml:"namespace"`
	Collection  string  `json:"collection" yaml:"collection"`
	Constraints *string `json:"constraint,omitempty" yaml:"constraint,omitempty"`
}

type Metadata struct {
	Type string `json:"type"`
	Data struct {
		Collection struct {
			ID        int    `json:"id"`
			Created   string `json:"created"`
			Modified  string `json:"modified"`
			Namespace struct {
				ID          int         `json:"id"`
				Created     string      `json:"created"`
				Modified    string      `json:"modified"`
				Description string      `json:"description"`
				Active      bool        `json:"active"`
				Name        string      `json:"name"`
				AvatarURL   string      `json:"avatar_url"`
				Location    interface{} `json:"location"`
				Company     interface{} `json:"company"`
				Email       interface{} `json:"email"`
				HTMLURL     string      `json:"html_url"`
				IsVendor    bool        `json:"is_vendor"`
				Owners      []int       `json:"owners"`
			} `json:"namespace"`
			Name           string  `json:"name"`
			Deprecated     bool    `json:"deprecated"`
			DownloadCount  int     `json:"download_count"`
			CommunityScore float64 `json:"community_score"`
			LatestVersion  struct {
				Pk           int     `json:"pk"`
				Version      string  `json:"version"`
				QualityScore float64 `json:"quality_score"`
				Created      string  `json:"created"`
				Modified     string  `json:"modified"`
				Metadata     struct {
					Name         string   `json:"name"`
					Tags         []string `json:"tags"`
					Issues       string   `json:"issues"`
					Readme       string   `json:"readme"`
					Authors      []string `json:"authors"`
					License      []string `json:"license"`
					Version      string   `json:"version"`
					Homepage     string   `json:"homepage"`
					Namespace    string   `json:"namespace"`
					Repository   string   `json:"repository"`
					Description  string   `json:"description"`
					Dependencies struct {
					} `json:"dependencies"`
					LicenseFile   interface{} `json:"license_file"`
					Documentation string      `json:"documentation"`
				} `json:"metadata"`
				Contents []struct {
					Name   string `json:"name"`
					Scores struct {
						Content       float64     `json:"content"`
						Quality       float64     `json:"quality"`
						Metadata      float64     `json:"metadata"`
						Compatibility interface{} `json:"compatibility"`
					} `json:"scores"`
					Metadata struct {
						ContainerMeta interface{} `json:"container_meta"`
					} `json:"metadata"`
					RoleMeta struct {
						Tags      []string    `json:"tags"`
						Author    string      `json:"author"`
						Company   string      `json:"company"`
						Licenese  interface{} `json:"licenese"`
						Platforms []struct {
							Name    string `json:"name"`
							Release string `json:"release"`
						} `json:"platforms"`
						Dependencies      []interface{} `json:"dependencies"`
						CloudPlatforms    []interface{} `json:"cloud_platforms"`
						MinAnsibleVersion float64       `json:"min_ansible_version"`
					} `json:"role_meta"`
					Description string `json:"description"`
					ContentType string `json:"content_type"`
				} `json:"contents"`
				ReadmeHTML string `json:"readme_html"`
			} `json:"latest_version"`
			CommunitySurveyCount int `json:"community_survey_count"`
			AllVersions          []struct {
				Pk           int     `json:"pk"`
				Version      string  `json:"version"`
				QualityScore float64 `json:"quality_score"`
				Created      string  `json:"created"`
				Modified     string  `json:"modified"`
				DownloadURL  string  `json:"download_url"`
			} `json:"all_versions"`
		} `json:"collection"`
	} `json:"data"`
}

func (c *Collection) Download(client *resty.Client, destination string) error {

	directory := path.Join(destination, c.Namespace, c.Collection)

	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		slog.Error("error creating directory", "namespace", c.Namespace, "collection", c.Collection, "error", err)
		return err
	}

	result := &Metadata{}
	_, err := client.
		R().
		SetQueryParam("namespace", c.Namespace).
		SetQueryParam("name", c.Collection).
		EnableTrace().
		SetResult(result).
		Get("https://galaxy.ansible.com/api/internal/ui/repo-or-collection-detail/")

	if err != nil {
		slog.Error("error downloading collection index", "namespace", c.Namespace, "collection", c.Collection, "error", err)
		return err
	}

	f, err := os.Create(path.Join(directory, "index.json"))
	if err != nil {
		slog.Error("error opening collection index.json for output", "namespace", c.Namespace, "collection", c.Collection, "error", err)
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(helpers.ToPrettyJSON(result)); err != nil {
		slog.Error("error persisting index.json", "namespace", c.Namespace, "collection", c.Collection, "error", err)
		return err
	}

	var filter *semver.Constraints
	if c.Constraints != nil {
		filter, err = semver.NewConstraint(*c.Constraints)
		if err != nil {
			slog.Error("error parsing constraints", "namespace", c.Namespace, "collection", c.Collection, "constraints", *c.Constraints, "error", err)
			return err
		}
	}

	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Printf("collection %s - %s (output: %s):\n", color.HiMagentaString(c.Namespace), color.HiMagentaString(c.Collection), directory)
	} else {
		fmt.Printf("collection %s - %s (output: %s):\n", c.Namespace, c.Collection, directory)
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	for _, version := range result.Data.Collection.AllVersions {
		link := fmt.Sprintf("https://galaxy.ansible.com%s", version.DownloadURL)
		if filter != nil {
			v, err := semver.NewVersion(version.Version)
			if err != nil {
				slog.Error("error parsing version", "namespace", c.Namespace, "collection", c.Collection, "version", version.Version, "error", err)
				return err
			}
			if !filter.Check(v) {
				if isatty.IsTerminal(os.Stdout.Fd()) {
					pterm.DefaultSpinner.InfoPrinter = &pterm.PrefixPrinter{
						MessageStyle: &pterm.Style{pterm.FgLightBlue},
						Prefix: pterm.Prefix{
							Style: &pterm.Style{pterm.FgBlack, pterm.BgLightBlue},
							Text:  "SKIPPED",
						},
					}
					pterm.DefaultSpinner.Info(fmt.Sprintf("v%s: download skipped [URL: %s]", version.Version, link))
					pterm.DefaultSpinner.Stop()
					fmt.Println()
				} else {
					fmt.Printf(" - v%s: %s\n", version.Version, "skipped")
				}
				continue
			}
		}

		client := grab.NewClient()
		request, _ := grab.NewRequest(directory, link)
		resp := client.Do(request)

		if isatty.IsTerminal(os.Stdout.Fd()) {
			err = func() error {
				spinner, _ := pterm.DefaultSpinner.WithShowTimer(true).Start(fmt.Sprintf("v%s: downloading %s...", version.Version, link))
				defer spinner.Stop()
				t := time.NewTicker(500 * time.Millisecond)
				defer t.Stop()

			loop:
				for {
					select {
					case <-signals:
						fmt.Printf("aborting...\n")
						os.Exit(1)
					case <-t.C:
						if resp.BytesComplete() == resp.Size() {
							spinner.Success(fmt.Sprintf("v%s: download succeeded in %s [URL: %s]", version.Version, resp.Duration(), link))
							break loop
						}
					case <-resp.Done:
						if err := resp.Err(); err != nil {
							spinner.Fail(fmt.Sprintf("v%s: download failed [URL: %s]", version.Version, link))
						} else {
							spinner.Success(fmt.Sprintf("v%s: download succeeded in %s [URL: %s]", version.Version, resp.Duration(), link))
							break loop
						}
						break loop
					}
				}
				return nil
			}()

		} else {
			fmt.Printf(" - v%s: downloaded %d bytes from %s (duration: %s)\n", version.Version, resp.Size(), link, resp.Duration())
		}

		//_, err := grab.Get(directory, link)
		// if err:= resp.Err(); err != nil {
		// 	slog.Error("error downloading collection", "namespace", c.Namespace, "collection", c.Collection, "error", err)
		// }
		// if isatty.IsTerminal(os.Stdout.Fd()) {
		// 	fmt.Printf(" - v%s: %s from %s\n", version.Version, color.GreenString("downloaded"), link)
		// } else {
		// 	fmt.Printf(" - v%s: %s from %s\n", version.Version, "downloaded", link)
		// }
	}

	return nil
}
