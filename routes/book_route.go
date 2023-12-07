package routes

import (
	"book-catalogue/controllers"

	"github.com/gofiber/fiber/v2"
)

func BookRoute(app *fiber.App) {
	app.Get("/books", controllers.GetAllBooks)
	app.Post("/book", controllers.CreateBook)
	app.Post("/book/dummy", controllers.AddDataDummy)
	app.Get("/book/:bookId", controllers.GetABook)
	app.Put("/book/:bookId", controllers.EditABook)
	app.Delete("/book/:bookId", controllers.DeleteABook)
}
