package main

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var commands = []string{
	"/plan",
	"/events",
	"/countdown",
	"/commands",
	"/signup",
}

//Events :: Declare array for Events
var Events = []Event{}

//Event :: Holds information about a planned event
type Event struct {
	Description  string
	Date         time.Time
	Participants []string
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

		if MessageType == "chat" {
		}
		if MessageType == "command" {
			ProcessCommand(update.Message.Text, ChatID, bot, update.Message.From)
		}

	}
}

//ProcessCommand :: Checks which command has been given and calls for the function related to that command
func ProcessCommand(command string, ChatID int64, bot *tgbotapi.BotAPI, user *tgbotapi.User) {

	command = strings.ToLower(command)

	if strings.Contains(command, "/plan") {
		PlanEvent(command, ChatID, bot)
	}

	if strings.Contains(command, "/events") {
		msg := tgbotapi.NewMessage(ChatID, ShowEvents())
		bot.Send(msg)
	}

	if strings.Contains(command, "/countdown") {
		CountdownEvent(command, ChatID, bot)
	}

	if strings.Contains(command, "/commands") {
		msg := tgbotapi.NewMessage(ChatID, CommandList())
		bot.Send(msg)
	}

	if strings.Contains(command, "/signup") {
		SignUpForEvent(user, command, bot, ChatID)
	}
}

//SignUpForEvent :: Adds the user to a specific event
func SignUpForEvent(user *tgbotapi.User, message string, bot *tgbotapi.BotAPI, ChatID int64) {

	text := strings.Split(message, " ")

	for n, x := range Events {
		if x.Description == text[1] {
			Events[n].Participants = append(Events[n].Participants, user.UserName)
			msg := tgbotapi.NewMessage(ChatID, "Je hebt je ingeschreven voor het evenement!")
			bot.Send(msg)
		}
	}

}

//CountdownEvent :: Calculates the amount of hours remaining for a requested event.
func CountdownEvent(command string, ChatID int64, bot *tgbotapi.BotAPI) {

	text := strings.Split(command, " ")

	if len(text) > 1 {
		for _, x := range Events {
			if x.Description == text[1] {
				t := time.Until(x.Date)

				msg := tgbotapi.NewMessage(ChatID, "Het is nog "+strconv.FormatFloat(t.Hours(), 'f', 0, 64)+" uur tot "+x.Description)
				bot.Send(msg)
			}
		}
	}
}

//ShowEvents :: Gets all the currently planned events and returns them in a string to display
func ShowEvents() string {

	var buffer bytes.Buffer

	buffer.WriteString("De volgende items staan er op de planning: \n")

	for _, x := range Events {
		buffer.WriteString("Naam: " + x.Description + " datum: " + x.Date.Format("2006-01-02") + " inschrijvingen: " + strconv.Itoa(len(x.Participants)) + " \n")
	}

	return buffer.String()
}

//CommandList :: Returns a message with the list of currently available commands
func CommandList() string {

	var buffer bytes.Buffer

	buffer.WriteString("De volgende commands zijn beschikbaar: \n")

	for _, x := range commands {
		buffer.WriteString(x + " \n")
	}

	return buffer.String()

}

//PlanEvent :: Creates the new event and adds it to the list of events
func PlanEvent(command string, ChatID int64, bot *tgbotapi.BotAPI) {

	var newEvent Event

	textInfo := strings.Split(command, " ")

	if len(textInfo) > 2 {

		newEvent.Description = textInfo[1]
		newEvent.Date, _ = time.Parse("2006-01-02", textInfo[2])
		newEvent.Participants = []string{}

		Events = append(Events, newEvent)

		msg := tgbotapi.NewMessage(ChatID, newEvent.Description+" succesvol aangemaakt!")
		bot.Send(msg)
	}
}

//messageType :: Checks if the sent message is a regular chat message or a command
func messageType(msg string) string {

	command := getCommand(msg)

	for _, x := range commands {
		if x == command[0] {
			return "command"
		}
	}
	return "chat"
}

//getCommand :: splits the text to return the command
func getCommand(msg string) []string {

	command := strings.Split(msg, " ")[:1]

	return command
}
