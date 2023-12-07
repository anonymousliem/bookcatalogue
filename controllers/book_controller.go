package controllers

import (
	"book-catalogue/configs"
	"book-catalogue/models"
	"book-catalogue/responses"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var bookCollection *mongo.Collection = configs.GetCollection(configs.DB, "books")
var validate = validator.New()

func GetAllBooks(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var books []models.Book
	defer cancel()

	results, err := bookCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleBook models.Book
		if err = results.Decode(&singleBook); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		books = append(books, singleBook)
	}

	return c.Status(http.StatusOK).JSON(
		responses.BookResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": books}},
	)
}

func CreateBook(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var book models.Book
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&book); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.BookResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&book); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.BookResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newBook := models.Book{
		Id:     primitive.NewObjectID(),
		Name:   book.Name,
		Author: book.Author,
	}

	result, err := bookCollection.InsertOne(ctx, newBook)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.BookResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetABook(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	bookId := c.Params("bookId")
	var book models.Book
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(bookId)

	err := bookCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&book)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.BookResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": book}})
}

func EditABook(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	bookId := c.Params("bookId")
	var book models.Book
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(bookId)

	//validate the request body
	if err := c.BodyParser(&book); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.BookResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&book); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.BookResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"name": book.Name, "author": book.Author}

	result, err := bookCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated book details
	var updatedBook models.Book
	if result.MatchedCount == 1 {
		err := bookCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedBook)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.BookResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedBook}})
}

func DeleteABook(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	bookId := c.Params("bookId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(bookId)

	result, err := bookCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.BookResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "book with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.BookResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "book successfully deleted!"}},
	)
}

func AddDataDummy(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var dummyBooks []models.Book

	// Generate 100 dummy books
	for i := 0; i < 100; i++ {
		book := models.Book{
			Id:     primitive.NewObjectID(),
			Name:   fmt.Sprintf("Book %d", i+1),
			Author: fmt.Sprintf("Author %d", i+1),
		}
		dummyBooks = append(dummyBooks, book)
	}

	// Convert dummyBooks to a slice of interface{}
	var booksAsInterface []interface{}
	for _, b := range dummyBooks {
		booksAsInterface = append(booksAsInterface, b)
	}

	// Insert the dummy books into the collection
	result, err := bookCollection.InsertMany(ctx, booksAsInterface)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(responses.BookResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": result.InsertedIDs}})
}
