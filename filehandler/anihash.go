package filehandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AniHashResponse is the response from the anihash API.
type AniHashResponse struct {
	File  *AniHashFile `json:"file"`
	State AniHashState `json:"state"`
}

// AniHashFile holds file information from anihash.
type AniHashFile struct {
	AnimeID         int        `json:"AnimeID"`
	AudioBitrate    int        `json:"AudioBitrate"`
	AudioCodec      string     `json:"AudioCodec"`
	CRC             string     `json:"CRC"`
	CreatedAt       time.Time  `json:"CreatedAt"`
	DeletedAt       *time.Time `json:"DeletedAt"`
	Ed2k            string     `json:"Ed2K"`
	EnglishName     string     `json:"EnglishName"`
	EpName          string     `json:"EpName"`
	EpNum           string     `json:"EpNum"`
	EpRomajiName    string     `json:"EpRomajiName"`
	EpisodeID       int        `json:"EpisodeID"`
	Extension       string     `json:"Extension"`
	FileID          int        `json:"FileID"`
	GroupID         int        `json:"GroupID"`
	GroupName       string     `json:"GroupName"`
	ID              int        `json:"ID"`
	MD5             string     `json:"MD5"`
	Quality         string     `json:"Quality"`
	RomajiName      string     `json:"RomajiName"`
	SHA1            string     `json:"SHA1"`
	Size            int64      `json:"Size"`
	Source          string     `json:"Source"`
	State           int        `json:"State"`
	Type            string     `json:"Type"`
	UpdatedAt       time.Time  `json:"UpdatedAt"`
	VideoBitrate    int        `json:"VideoBitrate"`
	VideoCodec      string     `json:"VideoCodec"`
	VideoResolution string     `json:"VideoResolution"`
	Year            string     `json:"Year"`
}

// AniHashState holds the state of the file in anihash.
type AniHashState struct {
	Error  string `json:"error"`
	FileID int    `json:"file_id"`
	State  string `json:"state"`
}

func getAniHashFile(ed2k string, size int64, anihashURL string) (*AniHashFile, error) {
	url := fmt.Sprintf("%s/query/ed2k?ed2k=%s&size=%d", anihashURL, ed2k, size)

	anihashResponse, err := makeRequest(url)
	if err != nil {
		return nil, err
	}

	// If state is FILE_PENDING, do an exponential backoff and retry for 10 minutes
	if anihashResponse.State.State == "FILE_PENDING" {
		backoff := 1 * time.Second
		for backoff < 10*time.Minute {
			time.Sleep(backoff)
			backoff *= 2

			anihashResponse, err = makeRequest(url)
			if err != nil {
				return nil, err
			}

			if anihashResponse.State.State != "FILE_PENDING" {
				break
			}
		}
	}

	if anihashResponse.State.State != "FILE_AVAILABLE" {
		return nil, fmt.Errorf("%s error: %s", anihashResponse.State.State, anihashResponse.State.Error)
	}

	return anihashResponse.File, nil
}

func makeRequest(url string) (*AniHashResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var anihashResponse AniHashResponse
	err = json.Unmarshal(body, &anihashResponse)
	if err != nil {
		return nil, err
	}

	return &anihashResponse, nil
}
