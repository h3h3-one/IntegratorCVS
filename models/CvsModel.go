package models

type CvsModel struct {
	MessageId string        `json:"messageId"`
	Plate     CvsModelPlate `json:"plate"`
}

type CvsModelPlate struct {
	LicensePlate string `json:"licensePlate"` // Номерной знак
	Color        string `json:"color"`        // Цвет автомобиля
	Model        string `json:"model"`        // Модель автомобиля
	Manufacturer string `json:"manufacturer"` // Производитель
}
