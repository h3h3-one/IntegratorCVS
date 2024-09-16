package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MqttServiceImpl представляет реализацию сервиса MQTT с клиентом и конфигурацией.
type MqttServiceImpl struct {
	mqttClient mqtt.Client // MQTT клиент для подключения к брокеру
	mqttConfig string      // Путь к файлу конфигурации
}

// MQTTClientModel представляет структуру конфигурации клиента MQTT.
type MQTTClientModel struct {
	MqttClientIp   string `json:"mqttClientIp"`   // IP-адрес MQTT клиента
	MqttClientPort int    `json:"mqttClientPort"` // Порт MQTT клиента
	MqttUsername   string `json:"mqttUsername"`   // Имя пользователя для подключения к MQTT
	MqttPassword   string `json:"mqttPassword"`   // Пароль для подключения к MQTT
}

// IntegratorCVSModel представляет данные, отправляемые в MQTT как полезная нагрузка.
type IntegratorCVSModel struct {
	Payload   string `json:"payload"`   // Полезная нагрузка
	CamNumber string `json:"camNumber"` // Номер камеры
}

// Publish публикует сообщение на указанный топик MQTT.
func (service *MqttServiceImpl) Publish(topic string, payload string, camNumber string, flag bool) {
	var mqttClientModel MQTTClientModel

	// Читаем данные конфигурации из файла
	configData, err := ioutil.ReadFile(service.mqttConfig)
	if err != nil {
		log.Fatalf("Ошибка чтения файла конфигурации: %v", err)
	}

	// Десериализуем JSON в структуру MQTTClientModel
	err = json.Unmarshal(configData, &mqttClientModel)
	if err != nil {
		log.Fatalf("Ошибка десериализации данных конфигурации: %v", err)
	}

	log.Printf("Создание подключения клиента... HOST_NAME = %s, PORT = %d, USERNAME = %s",
		mqttClientModel.MqttClientIp,
		mqttClientModel.MqttClientPort,
		mqttClientModel.MqttUsername)

	// Проверяем, установлен ли клиент MQTT
	if service.mqttClient == nil {
		// Настройка параметров клиента MQTT
		opts := mqtt.NewClientOptions()
		opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mqttClientModel.MqttClientIp, mqttClientModel.MqttClientPort))
		opts.SetUsername(mqttClientModel.MqttUsername)
		opts.SetPassword(mqttClientModel.MqttPassword)
		opts.SetConnectTimeout(5 * time.Second)
		opts.SetAutoReconnect(true)

		// Создание нового клиента MQTT
		service.mqttClient = mqtt.NewClient(opts)

		// Подключение к MQTT брокеру
		if token := service.mqttClient.Connect(); token.Wait() && token.Error() != nil {
			log.Fatalf("Ошибка подключения клиента: %v", token.Error())
		}

		log.Printf("Клиент успешно подключен к адресу: %s", service.mqttClient.OptionsReader())
	}

	// Если флаг установлен, публикуем сообщение
	if flag {
		log.Printf("Начало публикации. TOPIC: %s, PAYLOAD: %s, CAM_NUMBER: %s", topic, payload, camNumber)
		integratorModel := IntegratorCVSModel{Payload: payload, CamNumber: camNumber}
		jsonData, err := json.Marshal(integratorModel) // Сериализуем модель в JSON
		if err != nil {
			log.Fatalf("Ошибка сериализации JSON: %v", err)
		}

		// Публикуем сообщение на указанный топик
		if token := service.mqttClient.Publish(topic, 0, false, jsonData); token.Wait() && token.Error() != nil {
			log.Fatalf("Ошибка публикации сообщения: %v", token.Error())
		}
		log.Printf("Публикация успешна. Опубликовано = %s", string(jsonData))
	}
}

// isNewFile проверяет, существует ли файл по указанному пути, и создает его, если он не существует.
func (service *MqttServiceImpl) isNewFile(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) { // Проверка существования файла
		file, err := os.Create(filePath) // Создаем файл
		if err != nil {
			log.Fatalf("Ошибка создания файла конфигурации: %v", err)
		}
		defer file.Close() // Закрываем файл после использования

		// Создание пустого экземпляра модели MQTTClientModel и сериализация в JSON
		mqttClientModel := MQTTClientModel{}
		jsonData, err := json.MarshalIndent(mqttClientModel, "", "  ")
		if err != nil {
			log.Fatalf("Ошибка сериализации JSON: %v", err)
		}

		file.Write(jsonData) // Записываем данные JSON в файл
		log.Printf("Файл конфигурации успешно создан. Пожалуйста, перезапустите программу. ПУТЬ: %s", filePath)
		os.Exit(1) // Выход из программы после создания конфигурационного файла
	}
}

// CreateFileConfig создает файл конфигурации, если он не существует.
func (service *MqttServiceImpl) CreateFileConfig() {
	service.mqttConfig = "IntegratorConfig.json" // Определяем путь к конфигурационному файлу
	service.isNewFile(service.mqttConfig)        // Проверяем и создаем файл конфигурации, если необходимо
}
