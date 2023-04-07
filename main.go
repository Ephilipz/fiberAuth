package main

import (
	"fmt"
	"log"

	"github.com/Ephilipz/fiberAuth/config"
	"github.com/Ephilipz/fiberAuth/database"
	repo_gorm "github.com/Ephilipz/fiberAuth/repository/gorm"
	"github.com/Ephilipz/fiberAuth/route"
	"github.com/Ephilipz/fiberAuth/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("Error initializing config %s\n", err)
	}
	app := fiber.New()
	initApp(app, cfg)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.App.PORT)))
}

func initApp(app *fiber.App, cfg *config.Config) {
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

	if err := seedInitialData(userService, roleService); err != nil {
		log.Fatalf("Error seeding database")
	}

	app.Get("/metrics", monitor.New())
}

func seedInitialData(userService service.User, roleService service.Role) error {
	err := service.InitRoles(roleService)
	if err != nil {
		return err
	}
	return service.InitUsers(userService)
}
