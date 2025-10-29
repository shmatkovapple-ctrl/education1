package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	tele "gopkg.in/telebot.v4"
)

type User struct {
	Login    string
	Password string
}

var (
	bot        *tele.Bot
	userID     int64
	users      = make(map[int64]User)
	userStates = make(map[int64]string)
	tempLogin  = make(map[int64]string)
)

func nameHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Println("Привет:", name)
	_, err := bot.Send(&tele.Chat{ID: userID}, name)
	if err != nil {
		http.Error(w, "Не удалось отправить сообщение", http.StatusInternalServerError)
		log.Println("Ошибка отправки:", err)
		return
	}
}

func main() {
	pref := tele.Settings{
		Token:  "8394736122:AAHMPRuZ2dqi_RDCtsWm_qqQALBRuOJn3q8",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	var err error
	bot, err = tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/start", func(c tele.Context) error {
		userID = c.Sender().ID

		if _, exists := users[userID]; exists {
			return c.Send("Вы уже зарегестрированы")
		}

		userStates[userID] = "Ожидание ввода логина"
		return c.Send("Введите логин")
	})

	bot.Handle(tele.OnText, func(c tele.Context) error {
		userID := c.Sender().ID
		text := c.Text()

		switch userStates[userID] {

		case "Ожидание ввода логина":
			tempLogin[userID] = text
			userStates[userID] = "wait_password"
			return c.Send("Теперь введите пароль:")

		case "wait_password":
			users[userID] = User{
				Login:    tempLogin[userID],
				Password: text,
			}
			delete(userStates, userID)
			delete(tempLogin, userID)

			return c.Send("Регистрация завершена")

		default:
			return c.Send("Напишите /start для регистрации")
		}
	})

	go bot.Start()

	http.HandleFunc("/name", nameHandler)
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Произошла ошибка", err.Error())
	}
}
