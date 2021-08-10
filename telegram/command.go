package telegram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func playerControl(chatID int64, userID int32, cs int) {
	if chatID == config.GetChatID() || chatID > 0 {
		if isAdmin(config.GetChatID(), userID) {
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
}

func checkQueueSong(chatID, msgID int64) {
	if len(GetQueue()) > 0 {
		list := "Current queue:\n"
		for i, idx := range GetQueue() {
			list += fmt.Sprintf("<b>%v</b>. <code>%v</code>\n", i+1, songList[idx])
		}
		format, err := bot.ParseTextEntities(list, tdlib.NewTextParseModeHTML())
		if err != nil {
			log.Println(err)
			return
		}
		text := tdlib.NewInputMessageText(format, false, false)
		bot.SendMessage(chatID, 0, msgID, tdlib.NewMessageSendOptions(false, true, nil), nil, text)
	} else {
		msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText("No queue song.", nil), true, false)
		bot.SendMessage(chatID, 0, msgID, tdlib.NewMessageSendOptions(false, true, nil), nil, msgText)
	}
}

func checkLatestSong(chatID, msgID int64, offset int) {
	list := "Recently added:\n"
	for i := offset; i < offset+30; i++ {
		if songList[i] == nil {
			continue
		}
		list += fmt.Sprintf("<b>%v</b>. <code>%v</code>\n", i+1, songList[i])
	}
	format, err := bot.ParseTextEntities(list, tdlib.NewTextParseModeHTML())
	if err != nil {
		log.Println(err)
		return
	}
	text := tdlib.NewInputMessageText(format, false, false)
	bot.SendMessage(chatID, 0, msgID, tdlib.NewMessageSendOptions(false, true, nil), nil, text)
}
