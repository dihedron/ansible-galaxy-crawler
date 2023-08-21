package main

import (
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Collections *Collections `long:"collections" short:"c" description:"The list of collections to sync." required:"yes" json:"collections" yaml:"collections"`
	Directory   string       `long:"directory" short:"d" description:"The directory into which collections are copies." optional:"yes" default:"_dowloads" json:"directory" yaml:"directory"`
}

type CollectionMetadata struct {
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

func main() {

	options := Options{}
	if _, err := flags.Parse(&options); err != nil {
		os.Exit(1)
	}

	//fmt.Println(helpers.ToPrettyJSON(options))

	client := resty.New()
	for _, collection := range options.Collections.Entries {
		collection.Download(client, options.Directory)
	}
}
