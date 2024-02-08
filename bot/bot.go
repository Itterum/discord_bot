package bot

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Token       string
	ServerID    string
	TargetUsers []string
	Nicknames   []string
)

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}

func Run() {
	// create a session
	discord, err := discordgo.New("Bot " + Token)
	checkNilErr(err)

	loadNicknamesFromFile("nickname_list.txt")

	if len(Nicknames) == 0 {
		fmt.Println("Не удалось загрузить никнеймы из файла.")
		return
	}

	// add event handler
	discord.AddHandler(newMessage)

	// open session
	err = discord.Open()

	if err != nil {
		fmt.Println("Error open session", err)
		return
	}

	go changeNicknames(discord)

	defer func(discord *discordgo.Session) {
		err := discord.Close()
		if err != nil {

		}
	}(discord) // close session, after function termination

	// keep bot running until there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	/* prevent bot responding to its own message
	this is archived by looking into the message author id
	if message.author.id is same as bot.author.id then just return
	*/
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// respond to user message if it contains `!help` or `!bye`
	switch {
	case strings.Contains(message.Content, "!help"):
		_, err := discord.ChannelMessageSend(message.ChannelID, "Hello World😃")
		if err != nil {
			return
		}
	case strings.Contains(message.Content, "!bye"):
		_, err := discord.ChannelMessageSend(message.ChannelID, "Good Bye👋")
		if err != nil {
			return
		}
	}
}

func loadNicknamesFromFile(filename string) {
	file, err := os.Open(filename)

	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Ошибка при закрытии файла:", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		nickname := strings.TrimSpace(scanner.Text())
		if nickname != "" {
			Nicknames = append(Nicknames, nickname)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
	}
}

func changeNicknames(s *discordgo.Session) {
	for {
		for _, userID := range TargetUsers {
			newNickname := getRandomNickname()
			err := s.GuildMemberNickname(ServerID, userID, newNickname)
			if err != nil {
				fmt.Println("Ошибка при изменении никнейма пользователя", userID, ":", err)
			} else {
				fmt.Println("Никнейм пользователя", userID, "изменен на", newNickname)
			}
		}

		time.Sleep(24 * time.Hour)
	}
}

func getRandomNickname() string {
	index := rand.Intn(len(Nicknames))
	return Nicknames[index]
}
