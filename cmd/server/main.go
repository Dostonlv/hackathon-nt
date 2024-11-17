package main

import (
	"database/sql"
	"log"

	"github.com/Dostonlv/hackathon-nt/internal/api"
	"github.com/Dostonlv/hackathon-nt/internal/repository/postgres"
	"github.com/Dostonlv/hackathon-nt/internal/service"
	"github.com/Dostonlv/hackathon-nt/internal/utils"
	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func initializeCasbin() (*casbin.Enforcer, error) {
	// Create the adapter
	a := fileadapter.NewAdapter("config/policy.csv")

	// Initialize the enforcer
	enforcer, err := casbin.NewEnforcer("config/model.conf", a)
	if err != nil {
		return nil, err
	}

	// Load the policy from DB
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}

	return enforcer, nil
}

func main() {
	// Database connection
	connStr := "postgres://postgres:postgres@localhost:5433/tender_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize JWT util
	jwtSecret := "secreeet"
	jwtUtil := utils.NewJWTUtil(jwtSecret)

	// Initialize Casbin
	enforcer, err := initializeCasbin()
	if err != nil {
		log.Fatal("Failed to initialize Casbin: ", err)
	}

	// Add some basic policies
	enforcer.AddPolicy("admin", "/api/*", "*")
	enforcer.AddPolicy("user", "/api/users", "GET")
	enforcer.AddPolicy("user", "/api/users/:id", "GET")

	// Save the policies back to storage
	err = enforcer.SavePolicy()
	if err != nil {
		log.Fatal("Failed to save Casbin policies: ", err)
	}

	// Initialize repositories and services
	userRepo := postgres.NewUserRepo(db)
	authService := service.NewAuthService(userRepo, jwtUtil)
	// Redis connection
	redisAddr := "localhost:6379"
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Pass Redis client to NewTenderRepo
	tenderService := service.NewTenderService(postgres.NewTenderRepo(db, redisClient))
	bidService := service.NewBidService(postgres.NewBidRepo(db), postgres.NewTenderRepo(db, redisClient))
	// Setup router with Casbin enforcer
	router := api.SetupRouter(authService, tenderService, bidService, enforcer, jwtSecret)

	// Start the server
	log.Println("Server starting on :8888...")
	if err := router.Run(":8888"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
