package models

type Deployment struct {
	Services Services `json:"services"`
}

type Services struct {
	Location Location `json:"location"`
}

type Location struct {
	Image string `json:"image"`
}
