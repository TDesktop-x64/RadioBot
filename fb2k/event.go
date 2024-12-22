package fb2k

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/RadioBot/helper"
	"github.com/c0re100/RadioBot/utils"
	tdlib "github.com/c0re100/gotdlib/client"
	"github.com/r3labs/sse/v2"
)

var (
	rx         sync.Mutex
	skipSwitch = make(chan bool, 1)
)

func isSameAsCurrent(songName string) bool {
	if config.GetStatus().GetCurrent() == songName {
		return true
	}
	return false
}

func sendNewMessage(msg *tdlib.SendMessageRequest) {
	m, newErr := bot.SendMessage(msg)
	if newErr != nil {
		log.Println("[Send] Failed to broadcast current song...", newErr)
		return
	}
	_, _ = bot.PinChatMessage(&tdlib.PinChatMessageRequest{ChatId: msg.ChatId, MessageId: m.Id, DisableNotification: true})
	_, _ = bot.DeleteMessages(helper.NewDeleteMessages(msg.ChatId, []int64{m.Id + 1048576}))
	config.SetPinnedMessage(m.Id)
	config.SaveConfig()
}

func getEvent() {
	fmt.Println("[Player] Update Event Receiver")
	client := sse.NewClient("http://127.0.0.1:" + strconv.Itoa(config.GetBeefWebPort()) + "/api/query/updates?player=true&trcolumns=%25artist%25%20-%20%25title%25,%25artist%25,%25title%25,%25album%25")

	routine := []int{}
	manualKill := make(chan bool, 1)

	client.Subscribe("messages", func(msg *sse.Event) {
		data := string(msg.Data)
		if data != "{}" {
			var event utils.Event
			if err := json.Unmarshal(msg.Data, &event); err == nil {
				if len(event.Player.ActiveItem.Columns) >= 1 {
					defer func() {
						go func(dur int64) {
							rID := getGoID()
							routine = append(routine, rID)
							if len(routine) > 1 {
								manualKill <- true
							}

							select {
							case <-manualKill:
								if routine[0] == rID && len(routine) > 1 {
									routine = append(routine[:0], routine[1:]...)
									fmt.Printf("Song change by manually: Goroutine #%v killed!\n", rID)
									return
								}
							case <-skipSwitch:
								routine = append(routine[:0], routine[1:]...)
								fmt.Printf("Song change by vote: Goroutine #%v killed!\n", rID)
								return
							case <-time.After(time.Duration(dur)*time.Second - 500*time.Millisecond):
								checkNextSong()
							}
							routine = append(routine[:0], routine[1:]...)
						}(int64(event.Player.ActiveItem.Duration - event.Player.ActiveItem.Position))
					}()

					songName := event.Player.ActiveItem.Columns[0]
					if isSameAsCurrent(songName) {
						return
					}

					artist := event.Player.ActiveItem.Columns[1]
					track := event.Player.ActiveItem.Columns[2]
					album := event.Player.ActiveItem.Columns[3]
					idx := event.Player.ActiveItem.Index
					text := fmt.Sprintf("Now playing: \n"+
						"Artist: %v\n"+
						"Track: %v\n"+
						"Album: %v\n"+
						"Duration: %v", utils.IsEmpty(artist), utils.IsEmpty(track), utils.IsEmpty(album), utils.SecondsToMinutes(int64(event.Player.ActiveItem.Duration)))
					cID := config.GetChatID()
					mID := config.GetPinnedMessage()

					if mID == 0 {
						msgText := helper.NewSimpleMessage(cID, 0, html.EscapeString(text))
						sendNewMessage(msgText)
					} else {
						_, getErr := bot.GetMessage(helper.NewGetMessage(cID, mID))
						if getErr != nil {
							msgText := helper.NewSimpleMessage(cID, 0, html.EscapeString(text))
							sendNewMessage(msgText)
						} else {
							_, editErr := bot.EditMessageText(helper.NewEditMessageText(cID, mID, nil, helper.NewFormattedText(html.EscapeString(text), nil)))
							if editErr != nil {
								log.Println("[Edit] Failed to broadcast current song...", editErr)
								return
							}
						}
					}

					rx.Lock()
					recent := config.GetStatus().GetRecent()
					recent = append(recent, idx)
					config.SetRecentSong(recent)
					if len(recent) >= config.GetRecentLimit() {
						recent = append(recent[:0], recent[1:]...)
						config.SetRecentSong(recent)
					}
					config.SetCurrentSong(songName)
					config.SaveStatus()
					rx.Unlock()
				}
			}
		}
	})
}
