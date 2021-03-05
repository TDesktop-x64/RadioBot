package fb2k

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/RadioBot/utils"
	"github.com/c0re100/go-tdlib"
	"github.com/r3labs/sse"
)

var (
	rx         sync.Mutex
	killSwitch = make(chan bool, 1)
)

func isSameAsCurrent(songName string) bool {
	if config.GetStatus().GetCurrent() == songName {
		return true
	}
	return false
}

func sendNewMessage(cId int64, msgText *tdlib.InputMessageText) {
	m, newErr := bot.SendMessage(cId, 0, 0, nil, nil, msgText)
	if newErr != nil {
		log.Println("[Send] Failed to broadcast current song...", newErr)
		return
	}
	bot.PinChatMessage(cId, m.Id, true, false)
	bot.DeleteMessages(cId, []int64{m.Id + 1048576}, true)
	config.SetPinnedMessage(m.Id)
	config.SaveConfig()
}

func GetEvent() {
	fmt.Println("[Player] Update Event Receiver")
	client := sse.NewClient("http://127.0.0.1:" + strconv.Itoa(config.GetBeefWebPort()) + "/api/query/updates?player=true&trcolumns=%25artist%25%20-%20%25title%25,%25artist%25,%25title%25,%25album%25")

	client.Subscribe("messages", func(msg *sse.Event) {
		data := string(msg.Data)
		if data != "{}" {
			var event utils.Event
			if err := json.Unmarshal(msg.Data, &event); err == nil {
				if len(event.Player.ActiveItem.Columns) >= 1 {
					defer func() {
						go func(dur int64) {
							select {
							case <-killSwitch:
								fmt.Printf("Next song monitor: Goroutine #%v killed!\n", getGoId())
								return
							case <-time.After(time.Duration(dur)*time.Second - 500*time.Millisecond):
								checkNextSong()
							}
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
					msgText := tdlib.NewInputMessageText(tdlib.NewFormattedText(text, nil), true, false)
					cId := config.GetChatId()
					mId := config.GetPinnedMessage()

					if mId == 0 {
						sendNewMessage(cId, msgText)
					} else {
						_, getErr := bot.GetMessage(cId, mId)
						if getErr != nil {
							sendNewMessage(cId, msgText)
						} else {
							_, editErr := bot.EditMessageText(cId, mId, nil, msgText)
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
