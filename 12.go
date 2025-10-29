package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
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

const (
	StateWaitLogin      = "wait_login"
	StateWaitPassword   = "wait_password"
	StateLoginWaitLogin = "login_wait_login"
	StateLoginWaitPass  = "login_wait_password"
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

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func main() {
	erld := godotenv.Load()
	if erld != nil {
		log.Fatal("Error loading .env file")
	}
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
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
			return c.Send("Вы уже зарегестрированы, для авторизации используйте /login")
		}
		return c.Send("Здравствуйте, для регистрации используйте /register")
	})

	bot.Handle("/register", func(c tele.Context) error {
		userID := c.Sender().ID

		if _, exists := users[userID]; exists {
			return c.Send("Вы уже зарегистрированы")
		}

		userStates[userID] = StateWaitLogin
		return c.Send("Введите логин:")
	})

	bot.Handle("/login", func(c tele.Context) error {
		userID := c.Sender().ID
		if _, exists := users[userID]; !exists {
			return c.Send("Вы ещё не зарегестрированы")
		}
		userStates[userID] = StateLoginWaitLogin
		return c.Send("Введите ваш логин")
	})

	bot.Handle(tele.OnText, func(c tele.Context) error {
		userID := c.Sender().ID
		text := c.Text()

		switch userStates[userID] {

		case StateWaitLogin:
			tempLogin[userID] = text
			userStates[userID] = "wait_password"
			return c.Send("Теперь введите пароль:")

		case StateWaitPassword:
			users[userID] = User{
				Login:    tempLogin[userID],
				Password: hashPassword(text),
			}

			delete(userStates, userID)
			delete(tempLogin, userID)

			return c.Send("Регистрация завершена")

		case StateLoginWaitLogin:
			if users[userID].Login != text {
				return c.Send("Неверный логин. Попробуйте /login снова")
			}
			userStates[userID] = StateLoginWaitPass
			return c.Send("Введите пароль:")

		case StateLoginWaitPass:
			if users[userID].Password != hashPassword(text) {
				return c.Send("Неверный пароль. Попробуйте /login снова")
			}

			delete(userStates, userID)
			return c.Send("Вы успешно авторизованы")

		default:
			return c.Send("Напишите /register для регистрации")
		}
	})

	go bot.Start()

	http.HandleFunc("/name", nameHandler)
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Произошла ошибка", err.Error())
	}
}
