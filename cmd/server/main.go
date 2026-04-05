package main

import (
    "context"
    "flag"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"

    "booking/internal/config"
    "booking/internal/database"
    "booking/internal/handler"
    "booking/internal/middleware"
    "booking/internal/repository"
    "booking/internal/service"
    "booking/internal/utils"

    "github.com/gorilla/mux"
    httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Room Booking Service API
// @version         1.0
// @description     Сервис бронирования переговорок
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
    seedFlag := flag.Bool("seed", false, "seed test data")
    flag.Parse()

    cfg := config.Load()

    db, err := database.NewDB(cfg)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    if err := database.Migrate(db); err != nil {
        log.Fatal("Failed to migrate:", err)
    }

    if *seedFlag {
        if err := database.Seed(db); err != nil {
            log.Fatal("Failed to seed:", err)
        }
        log.Println("Seeding completed successfully")
        return
    }

    if err := database.EnsureDummyUsers(db); err != nil {
        log.Printf("Warning: failed to ensure dummy users: %v", err)
    }

    utils.InitJWT(cfg.JWTSecret)

    userRepo := repository.NewUserRepository(db)
    roomRepo := repository.NewRoomRepository(db)
    scheduleRepo := repository.NewScheduleRepository(db)
    slotRepo := repository.NewSlotRepository(db)
    bookingRepo := repository.NewBookingRepository(db)

    authService := service.NewAuthService(userRepo, cfg.JWTSecret)
    roomService := service.NewRoomService(roomRepo)
    scheduleService := service.NewScheduleService(scheduleRepo, roomRepo)
    slotService := service.NewSlotService(slotRepo, scheduleRepo)
    bookingService := service.NewBookingService(bookingRepo, slotRepo, userRepo)

    authHandler := handler.NewAuthHandler(authService)
    roomHandler := handler.NewRoomHandler(roomService)
    scheduleHandler := handler.NewScheduleHandler(scheduleService)
    slotHandler := handler.NewSlotHandler(slotService)
    bookingHandler := handler.NewBookingHandler(bookingService)
    infoHandler := handler.NewInfoHandler()

    router := mux.NewRouter()
    router.Use(middleware.Logging)

    // Swagger UI
    router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

    // Public routes
    router.HandleFunc("/_info", infoHandler.Info).Methods("GET")
    router.HandleFunc("/dummyLogin", authHandler.DummyLogin).Methods("POST")
    router.HandleFunc("/register", authHandler.Register).Methods("POST")
    router.HandleFunc("/login", authHandler.Login).Methods("POST")

    // Protected routes
    api := router.PathPrefix("/").Subrouter()
    api.Use(middleware.Auth(cfg.JWTSecret))

    api.HandleFunc("/rooms/list", roomHandler.ListRooms).Methods("GET")
    api.HandleFunc("/rooms/create", roomHandler.CreateRoom).Methods("POST")
    api.HandleFunc("/rooms/{roomId}/schedule/create", scheduleHandler.CreateSchedule).Methods("POST")
    api.HandleFunc("/rooms/{roomId}/slots/list", slotHandler.ListSlots).Methods("GET")
    api.HandleFunc("/bookings/create", bookingHandler.CreateBooking).Methods("POST")
    api.HandleFunc("/bookings/list", bookingHandler.ListAllBookings).Methods("GET")
    api.HandleFunc("/bookings/my", bookingHandler.MyBookings).Methods("GET")
    api.HandleFunc("/bookings/{bookingId}/cancel", bookingHandler.CancelBooking).Methods("POST")

    srv := &http.Server{
        Addr:         cfg.ServerPort,
        Handler:      router,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    go func() {
        log.Printf("Server starting on %s", cfg.ServerPort)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("Server failed:", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server shutdown error:", err)
    }
}