package main

import (
	"log"

	"github.com/Ephilipz/fiberAuth/config"
	"github.com/Ephilipz/fiberAuth/database"
	"github.com/Ephilipz/fiberAuth/repository/gorm"
	"github.com/Ephilipz/fiberAuth/route"
	"github.com/Ephilipz/fiberAuth/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()

	app := fiber.New()
	initApp(app)

	log.Fatal(app.Listen(":3000"))
}

func initApp(app *fiber.App) {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("Error initializing config %s\n", err)
	}

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Error initializing db %s\n", err)
	}

	initServices(app, cfg, db)
}

func initServices(app *fiber.App, cfg *config.Config, db *gorm.DB) {
	userRepo := repo_gorm.NewUserGormRepo(db)
	userService := service.NewUserService(userRepo)

	roleRepo := repo_gorm.NewRoleGormRepo(db)
	roleService := service.NewroleService(roleRepo)

	jwtService, err := service.NewJWTService(cfg.Jwt.RSA)
	if err != nil {
		log.Fatalf("Error initializing jwt service %s\n", err)
	}

	route.Auth(app.Group("/auth"), userService, jwtService)

	// protected routes with jwt
	protected := route.ProtectedMiddlware(jwtService)
	route.User(app.Group("/user", protected), userService)
	route.Role(app.Group("/role", protected), roleService)

	app.Get("/metrics", monitor.New())
}
