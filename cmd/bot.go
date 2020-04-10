package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bsati/twitchbot/pkg/commands"
	"github.com/bsati/twitchbot/pkg/core"
	"github.com/bsati/twitchbot/pkg/handlers"
	"github.com/bsati/twitchbot/pkg/twitchirc"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Must specify nickname and OAuth-Token")
		os.Exit(1)
	}
	client := twitchirc.NewClient(true)
	env := &core.Env{
		CommandMatcher: commands.NewCommandMatcher("!", " "),
	}
	cmdlist := commands.BuildCommandList()
	for _, v := range cmdlist {
		env.CommandMatcher.Register(v.Aliases, v.Command)
	}
	mh := &handlers.MessageHandler{Env: env}
	client.AddHandler(mh.GetCommandHandler())
	client.AddHandler(func(irc *twitchirc.IRCClient, je *twitchirc.JoinEvent) {
		fmt.Printf("User: %v joined Channel %v\n", je.User, je.Channel)
	})
	err := client.Connect(os.Args[1], os.Args[2])
	client.JoinChannel("bsati")
	if err != nil {
		panic(err)
	}
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
