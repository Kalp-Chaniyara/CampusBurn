package controller

import (
	"campusburn-backend/dbConnection"
	"campusburn-backend/middleware"
	"campusburn-backend/model"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

/*
THE FLOW OF REGISTERING EVERY USER ON THE APPLICATION:
------------------------------------------------------
1. Get all the data coming from the client side through body parser.
2. Do sanity check on the data, if it is accurate or not.
3. Check if this user already exists in the db or not.
4. If exists, return with a 409 http code.
5. If not, then first hash the password using a middleware just before creating the user and saving it into the db.
6. Handle all the errors gracefully to send specific JSON messages back to the frontend.
*/

func RegisterUser(c *fiber.Ctx) error {

	var user model.User

	// CHECKING IF THE INCOMING DATA FROM THE REQUEST IS OKAY OR NOT
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// FINDING IF USER ALREADY EXISTS OR NOT
	result := dbConnection.DB.Where("Email = ?", user.Email).Limit(1).Find(&model.User{})

	// IF USER ALREADY EXISTS, THEN SENDING ERROR JSON RESPONSE
	if result.Error == nil && result.RowsAffected > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"Error": "User already exists",
		})
	}

	// TODO: I have to write a function (probably middleware) to hash the password entered by the user to make sure that only the hashed password goes into the database.

	//Hashing the password before saving it into the database
	hashedPassword, err := middleware.HashPassword(user.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Cannot save the password in the database",
		})
	}

	user.Password = hashedPassword

	// SAVING THE USER INTO THE DATABASE
	createdUser := dbConnection.DB.Create(&user)

	// IF USER NOT CREATED, THEN SENDING A ERROR RESPONSE
	if createdUser.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User not created",
		})
	}

	fmt.Println("The ID of the created user is: ", user.ID)
	fmt.Println("Number of rows affected: ", createdUser.RowsAffected)

	// IF USER CREATED, THEN SENDING SUCCESS RESPONSE
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})

}