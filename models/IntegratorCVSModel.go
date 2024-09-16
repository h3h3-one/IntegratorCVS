package models

import (
	"fmt"
	"time"
)

type IntegratorCVSModel struct {
	GRZ       string    // Государственный регистрационный знак
	CamNumber string    // Номер камеры
	Timestamp time.Time // Время создания записи
	Location  string    // Местоположение (можно добавить)
	VehicleID string    // Идентификатор транспортного средства (можно добавить)
}

func NewIntegratorCVSModel(grz string, camNumber string) *IntegratorCVSModel {
	return &IntegratorCVSModel{
		GRZ:       grz,
		CamNumber: camNumber,
		Timestamp: time.Now(), // Устанавливаем текущее время
	}
}

func NewIntegratorCVSModelEmpty() *IntegratorCVSModel {
	return &IntegratorCVSModel{
		Timestamp: time.Now(), // Устанавливаем текущее время
	}
}

func (m *IntegratorCVSModel) String() string {
	return fmt.Sprintf("IntegratorCVSModel{GRZ: %s, CamNumber: %s, Timestamp: %s, Location: %s, VehicleID: %s}",
		m.GRZ, m.CamNumber, m.Timestamp.Format(time.RFC3339), m.Location, m.VehicleID)
}

func (m *IntegratorCVSModel) IsEmpty() bool {
	return m.GRZ == "" && m.CamNumber == "" && m.Location == "" && m.VehicleID == ""
}
