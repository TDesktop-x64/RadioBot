package telegram

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/go-tdlib"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	bot       *tdlib.Client
	botId     int32
	userBot   *tdlib.Client
	userBotId int32
)

func New() (*tdlib.Client, *tdlib.Client) {
	tdlib.SetLogVerbosityLevel(0)
	tdlib.SetFilePath("./errors.txt")

	if _, err := os.Stat("instance"); os.IsNotExist(err) {
		if err := os.Mkdir("instance", 0755); err != nil {
			log.Fatal("Failed to create instance dir...")
		}
	}

	err := botLogin()
	if err != nil {
		log.Fatal("bot login failed:", err)
	}
	checkGroupIsExist(bot)

	if !config.IsWebEnabled() {
		err = userLogin()
		if err != nil {
			log.Fatal("userbot login failed:", err)
		}
		checkGroupIsExist(userBot)
	}

	savePlaylistIndexAndName()
	Receiver()

	return bot, userBot
}

func newClient(name string) *tdlib.Client {
	return tdlib.NewClient(tdlib.Config{
		APIID:               config.GetApiId(),
		APIHash:             config.GetApiHash(),
		SystemLanguageCode:  "en",
		DeviceModel:         "Radio Controller",
		SystemVersion:       "1.0",
		ApplicationVersion:  "1.0",
		UseMessageDatabase:  true,
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseTestDataCenter:   false,
		DatabaseDirectory:   "./instance/" + name + "-db",
		FileDirectory:       "./instance/" + name + "-files",
		IgnoreFileNames:     false,
	})
}

func botLogin() error {
	bot = newClient("bot")

	for {
		currentState, _ := bot.Authorize()
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			_, err := bot.CheckAuthenticationBotToken(config.GetBotToken())
			if err != nil {
				log.Fatal(err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
			me, err := bot.GetMe()
			if err != nil {
				return err
			}
			botId = me.Id
			fmt.Println(me.Username + " connected.")
			break
		}
	}
	return nil
}

func userLogin() error {
	userBot = newClient("user")

	for {
		currentState, _ := userBot.Authorize()
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			fmt.Print("Enter phone: ")
			var number string
			fmt.Scanln(&number)
			_, err := userBot.SendPhoneNumber(number)
			if err != nil {
				fmt.Printf("Error sending phone number: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitCodeType {
			fmt.Print("Enter code: ")
			var code string
			fmt.Scanln(&code)
			_, err := userBot.SendAuthCode(code)
			if err != nil {
				fmt.Printf("Error sending auth code : %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPasswordType {
			fmt.Print("Enter Password: ")
			bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				fmt.Println(err)
			}
			_, err = userBot.SendAuthPassword(string(bytePassword))
			if err != nil {
				fmt.Printf("Error sending auth password: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
			me, err := userBot.GetMe()
			if err != nil {
				return err
			}
			userBotId = me.Id
			fmt.Println("\nHello!", me.FirstName, me.LastName, "("+me.Username+")")
			break
		}
	}

	return nil
}

func Receiver() {
	go newMessages()
	go callbackQuery()
	if !config.IsWebEnabled() {
		go newGroupCallUpdate()
		go newGroupCallPtcpUpdate()
		JoinGroupCall()
	}
}

func checkGroupIsExist(cl *tdlib.Client) {
	chatId := config.GetChatId()
	if chatId == 0 {
		uName := config.GetChatUsername()
		if uName == "" {
			log.Fatal("Username should not empty.")
		}
		s, err := cl.SearchPublicChat(uName)
		if err != nil {
			log.Fatal("SearchPublicChat error:", err)
		}
		_, err = cl.GetChat(s.Id)
		if err != nil {
			log.Fatal("GetChat error:", err)
		}
		config.SetChatId(s.Id)
		config.SaveConfig()
	} else {
		_, err := cl.GetChat(config.GetChatId())
		if err != nil {
			log.Fatal("GetChat error:", err)
		}
	}
}
