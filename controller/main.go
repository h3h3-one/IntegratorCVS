package main

import (
	"encoding/json"
	"log"
	"net/http"

	"integratorcvs/models"
	"integratorcvs/service"

	"github.com/gorilla/mux"
)

type MqttController struct {
	mqttService service.MqttService
}

func NewMqttController(mqttService service.MqttService) *MqttController {
	return &MqttController{mqttService: mqttService}
}

func (c *MqttController) SendMessage(w http.ResponseWriter, r *http.Request) {
	var message models.CvsModel

	// Декодируем входящее сообщение JSON в модель
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		log.Printf("Ошибка декодирования JSON: %v", err)
		http.Error(w, "Ошибочный запрос: неверный формат JSON", http.StatusBadRequest)
		return
	}

	topic := "Parking/IntegratorCVS" // Топик для публикации
	log.Printf("Получен JSON = %+v", message)

	// Запись полезной нагрузки (payload) в MQTT
	flag := true // Условный флаг, можно изменить логику по необходимости
	log.Printf("Получен POST запрос. TOPIC: %s PAYLOAD: %s CAM_NUMBER: %s", topic, message.Plate.Model, message.Plate.LicensePlate)

	// Публикуем сообщение в указанный топик
	if err := c.mqttService.Publish(topic, message.Plate.Model, message.Plate.LicensePlate, flag); err != nil {
		log.Printf("Ошибка публикации сообщения: %v", err)
		http.Error(w, "Не удалось отправить сообщение", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ с заголовком успешного завершения и кодом 200
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}

func main() {
	r := mux.NewRouter()                                    // Создаем новый маршрутизатор
	broker := "tcp://localhost:8080"                        // Здесь укажите адрес вашего MQTT брокера
	clientID := "dddddddddd"                                // Укажите клиентский ID
	mqttService := service.NewMqttService(broker, clientID) // Инициализация сервиса MQTT с параметрами
	mqttController := NewMqttController(mqttService)        // Создание контроллера

	// Определение маршрута для отправки сообщений
	r.HandleFunc("/send", mqttController.SendMessage).Methods("POST")
	http.Handle("/", r) // Обработка корневого маршрута

	log.Println("Сервер запущен на порту :8080") // Лог сообщения о запуске сервера
	log.Fatal(http.ListenAndServe(":8080", nil)) // Запуск сервера и обработка ошибок
}
