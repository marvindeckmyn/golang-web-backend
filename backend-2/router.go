package main

func setRouter() {
	footballerHandlers := newFootballerHandlers()
	Get("/footballer", footballerHandlers.getFootballers)
	Get("/footballer/", footballerHandlers.getFootballer)
	Post("/footballer", footballerHandlers.postFootballer)
}
