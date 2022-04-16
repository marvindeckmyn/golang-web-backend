package main

func setRouter() {
	Get("/footballer", getFootballers)
	Get("/footballer/", getFootballer)
	Post("/footballer", postFootballer)
	Put("/footballer/", updateFootballer)
	Delete("/footballer/", deleteFootballer)
}
