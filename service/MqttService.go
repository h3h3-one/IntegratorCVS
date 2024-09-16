package service

import (
	"errors"
	"log"
	"regexp"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttService interface {
	Publish(topic string, payload string, camNumber string, flag bool) error
}

type mqttServiceImpl struct {
	client mqtt.Client // MQTT клиент
}

func NewMqttService(broker string, clientID string) MqttService {
	// Настройка параметров клиента MQTT.
	opts := mqtt.NewClientOptions().AddBroker(broker) // Добавление адреса брокера
	opts.SetClientID(clientID)                        // Установка идентификатора клиента
	opts.SetCleanSession(true)                        // Установка чистой сессии (без сохранения состояния)

	// Создание нового клиента MQTT
	client := mqtt.NewClient(opts)

	// Подключение к брокеру MQTT
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Ошибка подключения к брокеру MQTT: %v", token.Error()) // Логирование ошибки, если не удалось подключиться
	}

	log.Println("Успешно подключено к брокеру MQTT") // Логирование успешного подключения

	// Возвращает новый экземпляр mqttServiceImpl с подключенным клиентом
	return &mqttServiceImpl{client: client}
}

func (m *mqttServiceImpl) Publish(topic string, payload string, camNumber string, flag bool) error {
	// Проверка, что топик не пуст
	if topic == "" {
		return errors.New("топик не может быть пустым") // Возвращает ошибку, если топик пустой
	}

	// Проверка полезной нагрузки на допустимые символы и длину
	if !isValidPayload(payload) {
		return errors.New("полезная нагрузка содержит недопустимые символы") // Возвращает ошибку, если полезная нагрузка недопустима
	}

	// Проверка номера камеры на допустимые символы и длину
	if !isValidCamNumber(camNumber) {
		return errors.New("номер камеры содержит недопустимые символы") // Возвращает ошибку, если номер камеры недопустим
	}

	// Публикация сообщения, если флаг установлен в true
	if flag {
		token := m.client.Publish(topic, 0, false, payload) // Публикация сообщения на указанный топик
		token.Wait()                                        // Ожидание завершения публикации
		if err := token.Error(); err != nil {
			return err // Возвращает ошибку, если произошла ошибка при публикации
		}
		log.Printf("Сообщение успешно опубликовано: ТОПИК: %s, ПОЛЕЗНАЯ НАГРУЗКА: %s", topic, payload) // Логирование успешной публикации
	} else {
		log.Println("Публикация отклонена: флаг не установлен") // Логирование, если публикация отклонена
	}

	return nil // Возвращает nil, если все прошло успешно
}

func isValidPayload(payload string) bool {
	return len(payload) >= 1 && len(payload) <= 10 && regexp.MustCompile(`^[ABCDEFHKMOPTXY\d]+$`).MatchString(payload)
}

func isValidCamNumber(camNumber string) bool {
	return len(camNumber) <= 10 && regexp.MustCompile(`^[A-Za-z\d]+$`).MatchString(camNumber)
}
