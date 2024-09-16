package models

import (
	"fmt"
	"log"
)

type MQTTClientModel struct {
	MqttClientId   string // Идентификатор клиента MQTT
	MqttClientIp   string // IP-адрес клиента
	MqttClientPort int    // Порт клиента
	MqttUsername   string // Имя пользователя для подключения к брокеру
	MqttPassword   string // Пароль для подключения к брокеру
}

func NewMQTTClientModel() *MQTTClientModel {
	return &MQTTClientModel{
		MqttClientId:   "Integrator",
		MqttClientIp:   "194.87.237.67",
		MqttClientPort: 1883,
		MqttUsername:   "admin",
		MqttPassword:   "333",
	}
}

func NewCustomMQTTClientModel(mqttClientId, mqttClientIp string, mqttClientPort int, mqttUsername, mqttPassword string) *MQTTClientModel {
	return &MQTTClientModel{
		MqttClientId:   mqttClientId,
		MqttClientIp:   mqttClientIp,
		MqttClientPort: mqttClientPort,
		MqttUsername:   mqttUsername,
		MqttPassword:   mqttPassword,
	}
}

func (m *MQTTClientModel) Connect() error {
	log.Printf("Подключение к MQTT брокеру по адресу %s:%d с клиентом %s", m.MqttClientIp, m.MqttClientPort, m.MqttClientId)
	return nil // Предположим, что подключение прошло успешно
}

// Validate проверяет корректность параметров клиента.
func (m *MQTTClientModel) Validate() error {
	if m.MqttClientId == "" {
		return fmt.Errorf("MqttClientId не может быть пустым")
	}
	if m.MqttClientIp == "" {
		return fmt.Errorf("MqttClientIp не может быть пустым")
	}
	if m.MqttClientPort <= 0 || m.MqttClientPort > 65535 {
		return fmt.Errorf("MqttClientPort должен быть в диапазоне 1-65535")
	}
	if m.MqttUsername == "" {
		return fmt.Errorf("MqttUsername не может быть пустым")
	}
	if m.MqttPassword == "" {
		return fmt.Errorf("MqttPassword не может быть пустым")
	}
	return nil
}

func (m *MQTTClientModel) String() string {
	return fmt.Sprintf("MQTTClientModel{MqttClientId: %s, MqttClientIp: %s, MqttClientPort: %d, MqttUsername: %s}",
		m.MqttClientId, m.MqttClientIp, m.MqttClientPort, m.MqttUsername)
}
