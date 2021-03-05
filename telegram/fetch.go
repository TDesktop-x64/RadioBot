package telegram

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/c0re100/RadioBot/config"
)

var (
	songList = make(map[int]string)
	mutex    sync.Mutex
)

type playlists struct {
	Playlists []struct {
		ID        string  `json:"id"`
		Index     int     `json:"index"`
		IsCurrent bool    `json:"isCurrent"`
		ItemCount int     `json:"itemCount"`
		Title     string  `json:"title"`
		TotalTime float64 `json:"totalTime"`
	} `json:"playlists"`
}

type playlistColumn struct {
	PlaylistItems struct {
		Items []struct {
			Columns []string `json:"columns"`
		} `json:"items"`
		Offset     int `json:"offset"`
		TotalCount int `json:"totalCount"`
	} `json:"playlistItems"`
}

func getPlaylistItemCount() (string, error) {
	resp, err := http.Get("http://localhost:" + strconv.Itoa(config.GetBeefWebPort()) + "/api/playlists")
	if err != nil {
		return "", errors.New("failed to get api")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("failed to read api response")
	}

	var pl playlists
	if err := json.Unmarshal(body, &pl); err == nil {
		for _, item := range pl.Playlists {
			if item.ID == config.GetPlaylistId() {
				return strconv.Itoa(item.ItemCount), nil
			}
		}
	}

	return "", errors.New("failed to parse api response")
}

func savePlaylistIndexAndName() {
	defer mutex.Unlock()
	mutex.Lock()
	count, err := getPlaylistItemCount()
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := http.Get("http://localhost:" + strconv.Itoa(config.GetBeefWebPort()) + "/api/playlists/" + config.GetPlaylistId() + "/items/0%3A" + count + "?columns=%25artist%25%20-%20%25title%25")
	if err != nil {
		log.Println("playlist not found...")
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("failed to read playlist response...")
		return
	}

	var plc playlistColumn
	if err := json.Unmarshal(body, &plc); err == nil {
		for idx, item := range plc.PlaylistItems.Items {
			if len(item.Columns[0]) > 0 {
				songList[idx] = item.Columns[0]
			}
		}
	}
}
