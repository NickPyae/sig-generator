package models

type ResponseModel struct {
	DeploymentSignature string `json:"deploymentSignature"`
}

type ErrorResponseModel struct {
	Code  int    `json:"code"`
	Error string `json:"message"`
}
