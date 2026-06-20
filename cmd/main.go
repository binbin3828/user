// @title           User Service API
// @version         1.0
// @description     User management and geo-based friend discovery REST API
// @termsOfService  https://github.com/binbin3828/user

// @contact.name   Developer
// @contact.url    https://github.com/binbin3828/user

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer " followed by your JWT token
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
	"user/pkg/mailer"
	"user/pkg/ratelimit"
	"user/pkg/redis"
	"user/service"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		runMigrate()
		return
	}
	runServer()
}

func runMigrate() {
	zapLog := logger.NewZapLogger()

	db, err := dbconn.NewMysql(zapLog)
	if err != nil {
		log.Fatalf("mysql init failed: %v", err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Friends{},
		&model.FriendRequest{},
		&model.Blacklist{},
		&model.PasswordResetToken{},
	); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
	zapLog.Sync()
	zapLog.Infof("migration completed successfully")
}

func runServer() {
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

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("get sql.DB failed: %v", err)
	}
	service.RegisterDBMetrics(sqlDB)

	userDao := dao.NewUserDao(db, zapLog)
	friendsDao := dao.NewFriendsDao(db, zapLog)
	friendRequestDao := dao.NewFriendRequestDao(db, zapLog)
	blacklistDao := dao.NewBlacklistDao(db, zapLog)
	passwordResetDao := dao.NewPasswordResetDao(db, zapLog)
	mailerInst := &mailer.DevMailer{}

	var rl ratelimit.Limiter
	if err := redis.Init(); err != nil {
		zapLog.Warnf("redis init failed, using memory limiter: %v", err)
		rl = ratelimit.NewMemoryLimiter(10, time.Minute)
	} else if redis.Enabled() {
		rl = ratelimit.NewRedisLimiterFromClient(redis.Client(), 10, time.Minute)
		zapLog.Infof("rate limiter: redis")
	} else {
		rl = ratelimit.NewMemoryLimiter(10, time.Minute)
		zapLog.Infof("rate limiter: memory (not suitable for multi-instance)")
	}

	svc := service.NewService(zapLog, userDao, friendsDao, friendRequestDao, blacklistDao, passwordResetDao, mailerInst, rl)
	router := service.NewRouter(svc)

	tlsEnabled, _ := config.Get("config.tls.enabled").(bool)
	certFile, _ := config.Get("config.tls.certFile").(string)
	keyFile, _ := config.Get("config.tls.keyFile").(string)

	tlsConfigured := tlsEnabled && certFile != "" && keyFile != ""

	var httpsServer *http.Server
	var httpServer *http.Server

	if tlsConfigured {
		httpsServer = &http.Server{
			Addr:    ":443",
			Handler: router,
		}

		redirectHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			target := "https://" + r.Host + r.URL.String()
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		})
		httpServer = &http.Server{
			Addr:    ":8080",
			Handler: redirectHandler,
		}

		go func() {
			zapLog.Infof("HTTPS server starting on :443")
			if err := httpsServer.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTPS server error: %v", err)
			}
		}()
		go func() {
			zapLog.Infof("HTTP redirect server starting on :8080 -> :443")
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTP redirect server error: %v", err)
			}
		}()
	} else {
		zapLog.Warnf("TLS not configured (enabled=%v, cert=%v, key=%v), serving HTTP on :8080", tlsEnabled, certFile != "", keyFile != "")
		httpServer = &http.Server{
			Addr:    ":8080",
			Handler: router,
		}
		go func() {
			zapLog.Infof("HTTP server starting on :8080")
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTP server error: %v", err)
			}
		}()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	zapLog.Infof("received signal %v, shutting down...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if httpsServer != nil {
		if err := httpsServer.Shutdown(ctx); err != nil {
			zapLog.Errorf("HTTPS server shutdown error: %v", err)
		}
	}
	if httpServer != nil {
		if err := httpServer.Shutdown(ctx); err != nil {
			zapLog.Errorf("HTTP server shutdown error: %v", err)
		}
	}

	sqlDB, err = db.DB()
	if err == nil {
		sqlDB.Close()
	}

	zapLog.Sync()
}
