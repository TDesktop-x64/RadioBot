package fb2k

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/go-tdlib"
)

var (
	bot       *tdlib.Client
	SongQueue = make(chan int, 100)
)

func New(client *tdlib.Client) {
	bot = client
	restoreQueue()
	GetEvent()
}

func restoreQueue() {
	for _, q := range config.GetStatus().GetQueue() {
		SongQueue <- q
	}
}

func Play() {
	http.Post("http://127.0.0.1:8880/api/player/play", "", nil)
}

func PlayNext() {
	http.Post("http://127.0.0.1:8880/api/player/next", "", nil)
}

func PlayRandom() {
	http.Post("http://127.0.0.1:8880/api/player/play/random", "", nil)
}

func Stop() {
	http.Post("http://127.0.0.1:8880/api/player/stop", "", nil)
}

func Pause() {
	http.Post("http://127.0.0.1:8880/api/player/pause", "", nil)
}

func PlaySelected(selectedIdx int) {
	rx.Lock()
	defer rx.Unlock()
	http.Post("http://127.0.0.1:8880/api/player/play/"+config.GetPlaylistId()+"/"+strconv.Itoa(selectedIdx), "", nil)

	rList := config.GetStatus().GetRecent()
	qList := config.GetStatus().GetQueue()

	recent := append(rList, selectedIdx)
	config.SetRecentSong(recent)
	queue := append(qList[:0], qList[1:]...)
	config.SetQueueSong(queue)

	if len(recent) >= config.GetRecentLimit() {
		recent = append(recent[:0], recent[1:]...)
		config.SetRecentSong(recent)
	}

	config.SaveStatus()
}

func PushQueue(selectedIdx int) {
	rx.Lock()
	defer rx.Unlock()
	SongQueue <- selectedIdx

	queue := config.GetStatus().GetQueue()
	queue = append(queue, selectedIdx)
	config.SetQueueSong(queue)

	if len(queue) >= config.GetQueueLimit() {
		queue = append(queue[:0], queue[1:]...)
		config.SetQueueSong(queue)
	}

	config.SaveStatus()
}

func checkNextSong() {
	if len(SongQueue) == 0 {
		return
	}
	next := <-SongQueue
	PlaySelected(next)
}

func SetKillSwitch() {
	killSwitch <- true
}

func getGoId() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
