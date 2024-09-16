package models

type ModelPlate struct {
	Camera          string `json:"camera"`
	Channel         string `json:"channel"`
	DateTime        string `json:"dateTime"`
	Description     string `json:"description"`
	Direction       string `json:"direction"`
	GroupId         string `json:"groupId"`
	Id              string `json:"id"`
	Image           string `json:"image"`
	InList          string `json:"inList"`
	Passed          string `json:"passed"`
	Plate           string `json:"plate"`
	Quality         string `json:"quality"`
	StayTimeMinutes string `json:"stayTimeMinutes"`
	Type            string `json:"type"`
	Weight          string `json:"weight"`
}
