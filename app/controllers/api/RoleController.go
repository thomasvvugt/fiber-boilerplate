package api

import (
	"github.com/gofiber/fiber/v2"

	"go-fiber-v2-boilerplate/app/models"
	"go-fiber-v2-boilerplate/database"
)

// Return all roles as JSON
func GetAllRoles(c *fiber.Ctx) error {
	db := database.Instance()
	var Role []models.Role
	if res := db.Find(&Role); res.Error != nil {
		c.SendString("Error occurred while retrieving roles from the database")
		return res.Error
	}
	err := c.JSON(Role)
	if err != nil {
		panic("Error occurred when returning JSON of roles")
	}
	return err
}

// Return a single role as JSON
func GetRole(c *fiber.Ctx) error {
	db := database.Instance()
	Role := new(models.Role)
	id := c.Params("id")
	if res := db.Find(&Role, id); res.Error != nil {
		c.SendString("An error occurred when retrieving the role")
		return res.Error
	}
	if Role.ID == 0 {
		c.SendStatus(404)
		err := c.JSON(fiber.Map{
			"ID": id,
		})
		if err != nil {
			panic("Error occurred when returning JSON of a role")
		}
		return err
	}
	err := c.JSON(Role)
	if err != nil {
		panic("Error occurred when returning JSON of a role")
	}
	return err
}

// Add a single role to the database
func AddRole(c *fiber.Ctx) error {
	db := database.Instance()
	Role := new(models.Role)
	if err := c.BodyParser(Role); err != nil {
		c.SendString("An error occurred when parsing the new role")
		return err
	}
	if res := db.Create(&Role); res.Error != nil {
		c.SendString("An error occurred when storing the new role")
		return res.Error
	}
	err := c.JSON(Role)
	if err != nil {
		panic("Error occurred when returning JSON of a role")
	}
	return err
}

// Edit a single role
func EditRole(c *fiber.Ctx) error {
	db := database.Instance()
	id := c.Params("id")
	EditRole := new(models.Role)
	Role := new(models.Role)
	if err := c.BodyParser(EditRole); err != nil {
		c.SendString("An error occurred when parsing the edited role")
		return err
	}
	if res := db.Find(&Role, id); res.Error != nil {
		c.SendString("An error occurred when retrieving the existing role")
		return res.Error
	}
	// Role does not exist
	if Role.ID == 0 {
		c.SendStatus(404)
		err := c.JSON(fiber.Map{
			"ID": id,
		})
		if err != nil {
			panic("Error occurred when returning JSON of a role")
		}
		return err
	}
	Role.Name = EditRole.Name
	Role.Description = EditRole.Description
	db.Save(&Role)

	err := c.JSON(Role)
	if err != nil {
		panic("Error occurred when returning JSON of a role")
	}
	return err
}

// Delete a single role
func DeleteRole(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.Instance()

	var Role models.Role
	db.Find(&Role, id)
	if res := db.Find(&Role); res.Error != nil {
		c.SendString("An error occurred when finding the role to be deleted")
		return res.Error
	}
	db.Delete(&Role)

	err := c.JSON(fiber.Map{
		"ID": id,
		"Deleted": true,
	})
	if err != nil {
		panic("Error occurred when returning JSON of a role")
	}
	return err
}
