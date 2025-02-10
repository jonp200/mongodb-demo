package main

type Movie struct {
	Title string `json:"title" query:"title" validate:"not_blank"`
}
