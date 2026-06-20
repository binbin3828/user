package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user/dao"
	"user/model"
	"user/pkg/config"
	"user/pkg/dbconn"
	"user/pkg/logger"
	"user/service"
)

func main() {
	zapLog := logger.NewZapLogger()

	tp, err := service.InitTracerProvider("user-service")
	if err != nil {
		log.Fatalf("tracer provider init failed: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tp.Shutdown(ctx)
	}()

	db, err := dbconn.NewMysql(zapLog)
	if err != nil {
		log.Fatalf("mysql init failed: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Friends{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	userDao := dao.NewUserDao(db, zapLog)
	friendsDao := dao.NewFriendsDao(db, zapLog)

	svc := service.NewService(zapLog, userDao, friendsDao)
	router := service.NewRouter(svc)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		zapLog.Infof("server starting on :8080")

		tlsEnabled, _ := config.Get("config.tls.enabled").(bool)
		if tlsEnabled {
			certFile, _ := config.Get("config.tls.certFile").(string)
			keyFile, _ := config.Get("config.tls.keyFile").(string)
			if err := server.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
				log.Fatalf("server error: %v", err)
			}
		} else {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("server error: %v", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	zapLog.Infof("received signal %v, shutting down...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		zapLog.Errorf("server shutdown error: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	zapLog.Sync()
}
