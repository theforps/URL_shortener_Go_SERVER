package models

type Request struct {
	BaseUrl    string `json:"base_url" validate:"required,url"`
	DayLife    int    `json:"day_life" validate:"required,min=1,max=7"`
	CodeLength int    `json:"code_length" validate:"required,min=6,max=12"`
}

type Response[T any] struct {
	StatusCode  int    `json:"status_code"`
	Description string `json:"description"`
	Data        *T     `json:"data"`
}

type URL struct {
	UniqCode    string `json:"uniq_code"`
	Views       int    `json:"views"`
	FinallyDate string `json:"finally_date"`
}

type URLStats struct {
	Views       int    `json:"views"`
	FinallyDate string `json:"finally_date"`
}
