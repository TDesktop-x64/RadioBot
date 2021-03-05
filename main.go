package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/c0re100/RadioBot/config"
	"github.com/c0re100/RadioBot/fb2k"
	"github.com/c0re100/RadioBot/telegram"
	"github.com/c0re100/RadioBot/web"
	"github.com/c0re100/RadioBot/wrtc"
)

var (
	version = "1.0.1"
)

func main() {
	fmt.Printf("RadioBot v%v\n", version)

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT)
	signal.Notify(ch, os.Interrupt, syscall.SIGKILL)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	signal.Notify(ch, os.Interrupt, syscall.SIGQUIT)
	signal.Notify(ch, os.Interrupt, syscall.SIGSEGV)
	go func() {
		<-ch
		config.Save()
		wrtc.Disconnect()
		fmt.Println("Shutdown...")
		os.Exit(0)
	}()

	config.Read()
	go web.StartServer()
	bot, _ := telegram.New()
	fb2k.New(bot)
}
