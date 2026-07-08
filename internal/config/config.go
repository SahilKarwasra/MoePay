package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoUri          string
	MongoDB           string
	ServerPort        string
	JWTSecret         string
	JWTDB             string
	RedisAddr         string
	RedisPass         string
	RedisDB           int
	JWTPrivateKeyPath string
	JWTPublicKeyPath  string
	SignupAdminSecret string
	EnableHTTPS       bool
	HTTPSPort         string
	SSLCertPath       string
	SSLKeyPath        string
}

func Load() (Config, error) {

	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("Failed to load .env")
	}

	mongoUri, err := extracEnv("MONGO_URI")
	if err != nil {
		return Config{}, err
	}

	mongoDB, err := extracEnv("DATABASE_NAME")
	if err != nil {
		return Config{}, err
	}

	jwtSecret, err := extracEnv("JWT_SECRET")
	if err != nil {
		return Config{}, err
	}

	jwtDB, err := extracEnv("JWT_DB")
	if err != nil {
		return Config{}, err
	}

	port, err := extracEnv("PORT")
	if err != nil {
		return Config{}, err
	}

	redisAddr, err := extracEnv("REDIS_ADDR")
	if err != nil {
		return Config{}, err
	}

	// redisPass, err := extracEnv("REDIS_PASS")
	// if err != nil {
	// 	return Config{}, err
	// }

	redisDB, err := getEnvInt("REDIS_DB", 0)
	if err != nil {
		return Config{}, err
	}

	// ApiKey configs removed

	jwtPrivateKeyPath := os.Getenv("JWT_PRIVATE_KEY_PATH")
	if jwtPrivateKeyPath == "" {
		jwtPrivateKeyPath = "private_key.pem"
	}

	jwtPublicKeyPath := os.Getenv("JWT_PUBLIC_KEY_PATH")
	if jwtPublicKeyPath == "" {
		jwtPublicKeyPath = "public_key.pem"
	}

	signupAdminSecret := os.Getenv("SIGNUP_ADMIN_SECRET")

	enableHTTPS := true
	if val := os.Getenv("ENABLE_HTTPS"); val == "false" {
		enableHTTPS = false
	}

	httpsPort := os.Getenv("HTTPS_PORT")
	if httpsPort == "" {
		httpsPort = "8443"
	}

	sslCertPath := os.Getenv("SSL_CERT_PATH")
	if sslCertPath == "" {
		sslCertPath = "cert.pem"
	}

	sslKeyPath := os.Getenv("SSL_KEY_PATH")
	if sslKeyPath == "" {
		sslKeyPath = "key.pem"
	}

	return Config{
		MongoUri:          mongoUri,
		MongoDB:           mongoDB,
		ServerPort:        port,
		JWTSecret:         jwtSecret,
		JWTDB:             jwtDB,
		RedisAddr:         redisAddr,
		RedisPass:         "",
		RedisDB:           redisDB,
		JWTPrivateKeyPath: jwtPrivateKeyPath,
		JWTPublicKeyPath:  jwtPublicKeyPath,
		SignupAdminSecret: signupAdminSecret,
		EnableHTTPS:       enableHTTPS,
		HTTPSPort:         httpsPort,
		SSLCertPath:       sslCertPath,
		SSLKeyPath:        sslKeyPath,
	}, nil

}

func extracEnv(key string) (string, error) {
	val := strings.TrimSpace(os.Getenv(key))

	if val == "" {
		return "", fmt.Errorf("missing req env: %s", key)
	}

	return val, nil
}

func getEnvInt(key string, fallback int) (int, error) {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i, nil
		}
	}
	return fallback, nil
}
