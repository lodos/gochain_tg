package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

const (
	maxMessageLength = 4096
	licenseServer    = "https://license.gochain.space/license/check"
)

type MessageRequest struct {
	ServerKey string `json:"server_key"`
	Message   string `json:"message"`
	BotApiKey string `json:"bot_api_key"`
	BotChatId string `json:"bot_chat_id"`
}

func main() {
	app := fiber.New()

	// Middleware для CORS
	app.Use(cors.New())

	// Ручка для открытия страницы
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("front/index.html")
	})

	// Ручка для добавления сообщения в канал
	app.Post("/add", func(c *fiber.Ctx) error {
		// Получаем параметры запроса
		var req MessageRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON format",
			})
		}

		// Проверяем длину текста
		if len(req.Message) > maxMessageLength {
			return c.Status(http.StatusExpectationFailed).JSON(fiber.Map{
				"error": "Message length exceeds the maximum allowed",
			})
		}

		// Проверяем лицензию
		goLicense, err := CheckLicense(req.ServerKey)
		if err != nil {
			log.Println("Error checking license:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"details": err.Error(),
			})
		}

		// Если лицензия не прошла
		if !goLicense {
			return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
				"error":   "Invalid License",
				"details": map[string]interface{}{"Error": map[string]interface{}{"error": fiber.StatusExpectationFailed, "message": "Invalid License"}},
			})
		}

		// Ваш токен бота
		token := req.BotApiKey

		// ID вашего канала (узнать его можно у бота @LyusherBot)
		channelID := req.BotChatId

		log.Println("token: ", token)
		log.Println("channelID: ", channelID)

		// Создаем нового бота с указанным токеном
		bot, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			log.Println("Error creating bot:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		}

		// Устанавливаем уровень логирования
		bot.Debug = true

		log.Printf("Авторизован как %s", bot.Self.UserName)

		// Отправляем сообщение в канал Telegram
		msg := tgbotapi.NewMessageToChannel(channelID, req.Message)
		msg.ParseMode = "markdown" // Устанавливаем parse_mode в "HTML"

		_, err = bot.Send(msg)
		if err != nil {
			log.Println("Error sending message:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"details": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Message sent to channel",
		})
	})

	// Запускаем Fiber
	log.Fatal(app.Listen(":3300"))
}

// Функция для проверки лицензии
func CheckLicense(serverKey string) (bool, error) {
	url := fmt.Sprintf("%s/%s", licenseServer, serverKey)

	response, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing response body:", err)
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	fmt.Println("Response Body:", string(body)) // Добавим вывод тела ответа

	var responseData map[string]interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		return false, err
	}

	statusCode := response.StatusCode

	if statusCode == http.StatusOK {
		return true, nil
	} else {
		return false, fmt.Errorf("non-OK status code: %d, Error: %v", statusCode, string(body))
	}
}
