package api

import (
	"fiber-boilerplate/app/models"
	"fiber-boilerplate/database"

	"github.com/gofiber/fiber/v2"
)

// Return all roles as JSON
func GetAllRoles(db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var Role []models.Role
		if response := db.Find(&Role); response.Error != nil {
			panic("Error occurred while retrieving roles from the database: " + response.Error.Error())
		}
		err := ctx.JSON(Role)
		if err != nil {
			panic("Error occurred when returning JSON of roles: " + err.Error())
		}
		return err
	}
}

// Return a single role as JSON
func GetRole(db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		Role := new(models.Role)
		id := ctx.Params("id")
		if response := db.Find(&Role, id); response.Error != nil {
			panic("An error occurred when retrieving the role: " + response.Error.Error())
		}
		if Role.ID == 0 {
			// Send status not found
			err := ctx.SendStatus(fiber.StatusNotFound)
			if err != nil {
				panic("Cannot return status not found: " + err.Error())
			}
			// Set ID
			err = ctx.JSON(fiber.Map{
				"ID": id,
			})
			if err != nil {
				panic("Error occurred when returning JSON of a role: " + err.Error())
			}
			return err
		}
		err := ctx.JSON(Role)
		if err != nil {
			panic("Error occurred when returning JSON of a role: " + err.Error())
		}
		return err
	}
}

// Add a single role to the database
func AddRole(db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		Role := new(models.Role)
		if err := ctx.BodyParser(Role); err != nil {
			panic("An error occurred when parsing the new role: " + err.Error())
		}
		if response := db.Create(&Role); response.Error != nil {
			panic("An error occurred when storing the new role: " + response.Error.Error())
		}
		err := ctx.JSON(Role)
		if err != nil {
			panic("Error occurred when returning JSON of a role: " + err.Error())
		}
		return err
	}
}

// Edit a single role
func EditRole(db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		EditRole := new(models.Role)
		Role := new(models.Role)
		if err := ctx.BodyParser(EditRole); err != nil {
			panic("An error occurred when parsing the edited role: " + err.Error())
		}
		if response := db.Find(&Role, id); response.Error != nil {
			panic("An error occurred when retrieving the existing role: " + response.Error.Error())
		}
		// Role does not exist
		if Role.ID == 0 {
			err := ctx.SendStatus(fiber.StatusNotFound)
			if err != nil {
				panic("Cannot return status not found: " + err.Error())
			}
			err = ctx.JSON(fiber.Map{
				"ID": id,
			})
			if err != nil {
				panic("Error occurred when returning JSON of a role: " + err.Error())
			}
			return err
		}
		Role.Name = EditRole.Name
		Role.Description = EditRole.Description
		db.Save(&Role)

		err := ctx.JSON(Role)
		if err != nil {
			panic("Error occurred when returning JSON of a role: " + err.Error())
		}
		return err
	}
}

// Delete a single role
func DeleteRole(db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		var Role models.Role
		db.Find(&Role, id)
		if response := db.Find(&Role); response.Error != nil {
			panic("An error occurred when finding the role to be deleted: " + response.Error.Error())
		}
		db.Delete(&Role)

		err := ctx.JSON(fiber.Map{
			"ID": id,
			"Deleted": true,
		})
		if err != nil {
			panic("Error occurred when returning JSON of a role: " + err.Error())
		}
		return err
	}
}
