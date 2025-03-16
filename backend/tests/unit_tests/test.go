package unit_tests

import (
	"context"
	"golbugames/backend/internal/database"
	"log"
)

func CreateUserTest(ctx context.Context) error {
	errTest := database.AddUser(ctx, "test", "test123")
	if errTest != nil {
		log.Println("Error when using function <database.AddUser>", errTest)
		return errTest
	}
	return nil
}

func DeleteUserTest(ctx context.Context) error {
	errTest := database.DeleteUser(ctx, 1)
	if errTest != nil {
		log.Println("Error when using function <database.AddUser>", errTest)
		return errTest
	}
	return nil
}

// Tests de stockage et de récupération des grilles

//  Tests de transformation de type pour le stockage
