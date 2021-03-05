package telegram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/RadioBot/fb2k"
	"github.com/c0re100/RadioBot/utils"
	"github.com/c0re100/go-tdlib"
)

func getCurrentPlaying(chatID, msgID int64) {
	resp, err := http.Get("http://127.0.0.1:" + strconv.Itoa(config.GetBeefWebPort()) + "/api/query?player=true&trcolumns=%25artist%25%20-%20%25title%25")
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var event utils.Event
	if err := json.Unmarshal(body, &event); err == nil {
		if len(event.Player.ActiveItem.Columns) >= 1 {
			songName := event.Player.ActiveItem.Columns[0]
			msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("Now playing: \n"+songName, nil), true, false)
			bot.SendMessage(chatID, 0, msgID, nil, nil, msgText)
		}
	}
}

func nominate(chatID, msgID int64, userID int32, arg string) {
	if arg == "" {
		msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("Track name or Artist name is empty.", nil), true, false)
		bot.SendMessage(chatID, 0, msgID, nil, nil, msgText)
		return
	}

	if ok, sec := canReqSong(userID); !ok {
		msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText(fmt.Sprintf("Please try again in %v seconds.", sec), nil), true, false)
		bot.SendMessage(chatID, 0, msgID, nil, nil, msgText)
	} else {
		list := searchSong(arg)
		if len(list) > 0 {
			sendCustomButtonMessage(chatID, msgID, list)
		} else {
			msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("No result.", nil), true, false)
			bot.SendMessage(chatID, 0, msgID, nil, nil, msgText)
		}
	}
}

func isAdmin(chatID int64, userID int32) bool {
	u, err := bot.GetChatMember(chatID, userID)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if u.Status.GetChatMemberStatusEnum() == "chatMemberStatusAdministrator" || u.Status.GetChatMemberStatusEnum() == "chatMemberStatusCreator" {
		return true
	}

	return false
}

func reload(chatID, msgID int64, userID int32) {
	if isAdmin(chatID, userID) {
		config.LoadConfig()
		savePlaylistIndexAndName()
		text := tdlib.NewInputMessageText(tdlib.NewFormattedText("Config&Playlist reloaded!", nil), false, false)
		bot.SendMessage(chatID, 0, msgID, tdlib.NewMessageSendOptions(false, true, nil), nil, text)
	}
}

func playerControl(chatID int64, userID int32, cs int) {
	if isAdmin(chatID, userID) {
		switch cs {
		case 0:
			fb2k.Play()
		case 1:
			fb2k.Stop()
		case 2:
			fb2k.Pause()
		case 3:
			fb2k.PlayRandom()
		default:

		}
	}
}
