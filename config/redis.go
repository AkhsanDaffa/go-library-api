package config

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

// Variabel Global supaya bisa diakses dari mana saja (Handler, Service, dll)
var RedisClient *redis.Client

// InitRedis bertugas menyalakan koneksi
func InitRedis() {
	// Ambil konfigurasi dari .env
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	// Kalau lupa set .env, kita kasih nilai default (Jaga-jaga)
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	// Konfigurasi Client
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword, // "" kalau kosong
		DB:       0,             // Default DB
	})

	// Coba PING ke Redis (Tes Koneksi)
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("❌ Gagal connect ke Redis: %v", err))
	}

	fmt.Println("✅ Berhasil connect ke Redis!")
}
