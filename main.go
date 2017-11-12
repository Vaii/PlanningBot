package main

import (
	"bytes"
	"log"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var commands = []string{
	"/plan",
	"/events",
}

//Events :: Declare array for Events
var Events = []Event{}

//Event :: Holds information about a planned event
type Event struct {
	Description string
	Date        time.Time
}

func main() {
	bot, err := tgbotapi.NewBotAPI("480240226:AAHgzoSDim7aKSJbThkiUUw4wJK-h1D_Yc0")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		ChatID := update.Message.Chat.ID
		MessageType := messageType(update.Message.Text)

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		//t3 := time.Until(t)

		//t2 := t3.Hours()
		//t4 := strconv.FormatFloat(t2, 'f', 0, 64)

		if MessageType == "chat" {
		}
		if MessageType == "command" {
			ProcessCommand(update.Message.Text, ChatID, bot)
		}

	}
}

//ProcessCommand :: Checks which command has been given and calls for the function related to that command
func ProcessCommand(command string, ChatID int64, bot *tgbotapi.BotAPI) {

	command = strings.ToLower(command)

	if strings.Contains(command, "/plan") {
		PlanEvent(command, ChatID, bot)
	}

	if strings.Contains(command, "/events") {
		msg := tgbotapi.NewMessage(ChatID, ShowEvents())
		bot.Send(msg)
	}
}

//ShowEvents :: Gets all the currently planned events and returns them in a string to display
func ShowEvents() string {

	var buffer bytes.Buffer

	buffer.WriteString("De volgende items staan er op de planning: \n")

	for _, x := range Events {
		buffer.WriteString("Naam: " + x.Description + " Datum: " + x.Date.Format("2006-01-02") + " \n")
	}

	return buffer.String()
}

//PlanEvent :: Creates the new event and adds it to the list of events
func PlanEvent(command string, ChatID int64, bot *tgbotapi.BotAPI) {

	var newEvent Event

	textInfo := strings.Split(command, " ")

	if len(textInfo) > 2 {

		newEvent.Description = textInfo[1]
		log.Printf(textInfo[2])
		newEvent.Date, _ = time.Parse("2006-01-02", textInfo[2])

		Events = append(Events, newEvent)

		msg := tgbotapi.NewMessage(ChatID, newEvent.Description+" succesvol aangemaakt!")
		bot.Send(msg)
	}
}

func messageType(msg string) string {

	command := getCommand(msg)

	for _, x := range commands {
		if x == command[0] {
			return "command"
		}
	}
	return "chat"
}

func getCommand(msg string) []string {

	command := strings.Split(msg, " ")[:1]

	return command
}
