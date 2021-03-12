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

func checkPlayerIsActive() {
	_, err := http.Get("http://localhost:" + strconv.Itoa(config.GetBeefWebPort()) + "/api/player")
	if err != nil {
		log.Fatal("BeefWeb is not running?\n",
			"If you're first time to use RadioBot, please read the documentation from this page.\n" +
			"https://github.com/c0re100/RadioBot#quick-start")
	}
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
			if item.ID == config.GetPlaylistID() {
				return strconv.Itoa(item.ItemCount), nil
			}
		}
	}

	return "", errors.New("failed to parse api response")
}

func savePlaylistIndexAndName() error {
	if songList != nil {
		songList = make(map[int]string)
	}
	defer mutex.Unlock()
	mutex.Lock()
	count, err := getPlaylistItemCount()
	if err != nil {
		return err
	}
	resp, err := http.Get("http://localhost:" + strconv.Itoa(config.GetBeefWebPort()) + "/api/playlists/" + config.GetPlaylistID() + "/items/0%3A" + count + "?columns=%25artist%25%20-%20%25title%25")
	if err != nil {
		return errors.New("playlist not found")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("failed to read playlist response")
	}

	var plc playlistColumn
	if err := json.Unmarshal(body, &plc); err == nil {
		for idx, item := range plc.PlaylistItems.Items {
			if len(item.Columns[0]) > 0 {
				songList[idx] = item.Columns[0]
			}
		}
	}
	return nil
}
