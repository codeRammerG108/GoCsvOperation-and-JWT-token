package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt"

	getCSVRoutes "github.com/codeRammerG108/goCSVAssignment/Routes"
	DBconnection "github.com/codeRammerG108/goCSVAssignment/db"
)

var DB *sql.DB
var sampleSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Message struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}

func main() {
	var err error
	fmt.Println("Go-Lang Assignment on CSV")
	DB, err := DBconnection.DBinit()
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("\nConnection Established")
	}
	defer DB.Close()

	// Routes
	app := fiber.New()
	app.Use(cors.New())

	// Just like a function Call
	getCSVRoutes.SetupCSVRoutes(app, DB)

	// Route to generate JWT
	app.Post("/login", handleLogin)

	// Protected route requiring JWT
	app.Get("/protected", verifyJWTMiddleware, handleProtected)

	// Start the Fiber app
	err = app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}

func handleLogin(c *fiber.Ctx) error {
	tokenString, err := generateJWT()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Message{
			Status: "error",
			Info:   "Failed to generate JWT",
		})
	}

	return c.JSON(fiber.Map{
		"token": tokenString,
	})
}

func handleProtected(c *fiber.Ctx) error {
	// Access the user claim stored in the context
	userClaim, ok := c.Locals("user").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(Message{
			Status: "error",
			Info:   "User claim not found in context",
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Access granted to protected route for user: %s", userClaim),
	})
}

func generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256) // Use HS256 for simplicity

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
	claims["authorized"] = true
	claims["user"] = "username"

	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyJWTMiddleware(c *fiber.Ctx) error {
	// Get the JWT token from the Authorization header
	tokenString := c.Get("Authorization")

	// Verify the JWT
	claims, err := verifyJWT(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(Message{
			Status: "error",
			Info:   "Unauthorized",
		})
	}

	// Extract user claim from MapClaims
	userClaim, ok := (*claims)["user"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(Message{
			Status: "error",
			Info:   "Invalid user claim",
		})
	}

	// Store the user claim in the context for further use
	c.Locals("user", userClaim)

	return c.Next()
}

func verifyJWT(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return sampleSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &claims, nil
}
