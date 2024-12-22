package telegram

import (
	"fmt"
	"log"
	"os"

	"github.com/c0re100/RadioBot/config"
	tdlib "github.com/c0re100/gotdlib/client"
)

var (
	bot       *tdlib.Client
	botID     int64
	userBot   *tdlib.Client
	userBotID int64
)

// New create telegram session
func New() (*tdlib.Client, *tdlib.Client) {
	checkPlayerIsActive() // Check music player is running

	tdlib.SetLogLevel(0)
	_, _ = tdlib.SetLogStream(&tdlib.SetLogStreamRequest{
		LogStream: &tdlib.LogStreamFile{
			Path:           "./errors.txt",
			MaxFileSize:    1048576000,
			RedirectStderr: true,
		},
	})

	if _, err := os.Stat("instance"); os.IsNotExist(err) {
		if err := os.Mkdir("instance", 0755); err != nil {
			log.Fatal("Failed to create instance dir...")
		}
	}

	botLogin()
	checkGroupIsExist(bot)

	if !config.IsWebEnabled() {
		userLogin()
		checkGroupIsExist(userBot)
	}

	if listErr := savePlaylistIndexAndName(); listErr != nil {
		log.Println(listErr)
	}
	createReceiver()

	return bot, userBot
}

func GetTdParameters(name string) *tdlib.SetTdlibParametersRequest {
	return &tdlib.SetTdlibParametersRequest{
		UseTestDc:              false,
		DatabaseDirectory:      "./instance/" + name + "-db",
		FilesDirectory:         "./instance/" + name + "-files",
		UseFileDatabase:        true,
		UseChatInfoDatabase:    true,
		UseMessageDatabase:     true,
		UseSecretChats:         false,
		ApiId:                  config.GetAPIID(),
		ApiHash:                config.GetAPIHash(),
		SystemLanguageCode:     "en",
		DeviceModel:            "Radio Controller",
		SystemVersion:          "1.0",
		ApplicationVersion:     "1.0",
	}
}

func botLogin() {
	authorizer := tdlib.BotAuthorizer(config.GetBotToken())

	authorizer.TdlibParameters <- GetTdParameters("bot")

	var err error
	bot, err = tdlib.NewClient(authorizer)
	if err != nil {
		log.Fatal(err)
	}

	me, err := bot.GetMe()
	if err != nil {
		log.Fatal(err)
	}
	botID = me.Id
	fmt.Println(me.Usernames.ActiveUsernames[0] + " connected.")
}

func userLogin() {
	authorizer := tdlib.ClientAuthorizer()
	go tdlib.CliInteractor(authorizer)

	authorizer.TdlibParameters <- GetTdParameters("user")

	var err error
	userBot, err = tdlib.NewClient(authorizer)
	if err != nil {
		log.Fatal(err)
	}

	me, err := userBot.GetMe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nHello!", me.FirstName, me.LastName, "("+me.Usernames.ActiveUsernames[0]+")")
}

func createReceiver() {
	go newMessages()
	go callbackQuery()
	if !config.IsWebEnabled() {
		go newGroupCallUpdate()
		go newGroupCallPtcpUpdate()
		joinGroupCall()
	}
}

func checkGroupIsExist(cl *tdlib.Client) {
	chatID := config.GetChatID()
	if chatID == 0 {
		uName := config.GetChatUsername()
		if uName == "" {
			log.Fatal("Username should not empty.")
		}
		s, err := cl.SearchPublicChat(&tdlib.SearchPublicChatRequest{Username: uName})
		if err != nil {
			log.Fatal("SearchPublicChat error:", err)
		}
		_, err = cl.GetChat(&tdlib.GetChatRequest{ChatId: s.Id})
		if err != nil {
			log.Fatal("GetChat error:", err)
		}
		config.SetChatID(s.Id)
		config.SaveConfig()
	} else {
		_, err := cl.GetChat(&tdlib.GetChatRequest{ChatId: config.GetChatID()})
		if err != nil {
			log.Fatal("GetChat error:", err)
		}
	}
}
