package telegram

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/RadioBot/fb2k"
	"github.com/c0re100/RadioBot/helper"
	"github.com/c0re100/RadioBot/utils"
	tdlib "github.com/c0re100/gotdlib/client"
)

func getCurrentPlaying(chatID, msgID int64) {
	resp, err := http.Get("http://127.0.0.1:" + strconv.Itoa(config.GetBeefWebPort()) + "/api/query?player=true&trcolumns=%25artist%25%20-%20%25title%25")
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var event utils.Event
	if err := json.Unmarshal(body, &event); err == nil {
		if len(event.Player.ActiveItem.Columns) >= 1 {
			songName := html.EscapeString(event.Player.ActiveItem.Columns[0])
			msg := helper.NewSimpleMessage(chatID, 0, msgID, "Now playing: \n"+songName)
			_, _ = bot.SendMessage(msg)
		}
	}
}

func isAdmin(chatID int64, userID int64) bool {
	u, err := bot.GetChatMember(&tdlib.GetChatMemberRequest{ChatId: chatID, MemberId: helper.NewMessageSenderUser(userID)})
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if u.Status.ChatMemberStatusType() == "chatMemberStatusAdministrator" || u.Status.ChatMemberStatusType() == "chatMemberStatusCreator" {
		return true
	}

	return false
}

func playerControl(chatID int64, userID int64, cs int) {
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
		format, err := tdlib.ParseTextEntities(helper.NewHTMLText(list))
		if err != nil {
			log.Println(err)
			return
		}
		msg := helper.NewEnititesMessage(chatID, 0, msgID, format)
		_, _ = bot.SendMessage(msg)
	} else {
		msg := helper.NewSimpleMessage(chatID, 0, msgID, "No queue song.")
		_, _ = bot.SendMessage(msg)
	}
}

func checkLatestSong(chatID, msgID int64, offset int) {
	list := "Recently added:\n"
	for i := offset; i < offset+30; i++ {
		if songList[i] == nil {
			continue
		}
		list += fmt.Sprintf("<b>%v</b>. <code>%v - %v</code>\n", i+1, html.EscapeString(songList[i].Artist), html.EscapeString(songList[i].Track))
	}
	format, err := tdlib.ParseTextEntities(helper.NewHTMLText(list))
	if err != nil {
		log.Println(err)
		return
	}
	msg := helper.NewEnititesMessage(chatID, 0, msgID, format)
	_, _ = bot.SendMessage(msg)
}
