package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "lab10-web-go-python/go-service/chat"
    "lab10-web-go-python/go-service/handlers"
    "lab10-web-go-python/go-service/middleware"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.New()
    r.Use(gin.Recovery())
    r.Use(middleware.Logger())

    // HTTP API (задание 1 средней)
    r.GET("/ping", handlers.Ping)
    r.GET("/hello/:name", handlers.Hello)
    r.POST("/echo", handlers.Echo)

    // WebSocket чат (задание 4 повышенной)
    hub := chat.NewHub()
    go hub.Run()
    r.GET("/ws", func(c *gin.Context) {
        chat.ServeWS(hub, c.Writer, c.Request)
    })

    // Graceful shutdown (задание 7 средней)
    srv := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }

    go func() {
        log.Println("Go HTTP server listening on :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down HTTP server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Forced to shutdown:", err)
    }
    log.Println("HTTP server exited")
}
