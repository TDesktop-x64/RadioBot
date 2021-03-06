package telegram

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

var (
	sch = gocron.NewScheduler(time.UTC)
)

func startScheduler() {
	sch.StartAsync()
}

func stopScheduler() {
	sch.Stop()
}

func addVoteJob(chatID, msgID int64, updateTime int32) {
	timeLeftJob, err := sch.Every(int(updateTime)).Second().Do(updateVote, chatID, msgID, true)
	if err != nil {
		log.Println("error creating job:", err)
		return
	}
	timeLeftJob.Tag("timeleft")
	timeLeftJob.RemoveAfterLastRun()
}

func addLoadGroupCallPtpcsJob() {
	_, err := sch.Every(1).Minute().Do(loadParticipants)
	if err != nil {
		log.Println("error creating job:", err)
		return
	}
}
