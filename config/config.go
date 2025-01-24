package config

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func LoadConfig(path string) IConfig {
	envMap, err := godotenv.Read(path)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &config{
		app: &app{
			host: envMap["APP_HOST"],
			port: func() int {
				p, err := strconv.Atoi(envMap["APP_PORT"])
				if err != nil {
					log.Fatalf("Error parsing APP_PORT: %v", err)
				}
				return p
			}(),
			name:    envMap["APP_NAME"],
			version: envMap["APP_VERSION"],
			bodyLimit: func() int {
				l, err := strconv.Atoi(envMap["APP_BODY_LIMIT"])
				if err != nil {
					log.Fatalf("Error parsing APP_BODY_LIMIT: %v", err)
				}
				return l
			}(),
			readTimeout: func() time.Duration {
				t, err := strconv.Atoi(envMap["APP_READ_TIMEOUT"])
				if err != nil {
					log.Fatalf("Error parsing APP_READ_TIMEOUT: %v", err)
				}
				return time.Duration(t) * time.Second
			}(),
			writeTimeout: func() time.Duration {
				t, err := strconv.Atoi(envMap["APP_WRITE_TIMEOUT"])
				if err != nil {
					log.Fatalf("Error parsing APP_WRITE_TIMEOUT: %v", err)
				}
				return time.Duration(t) * time.Second
			}(),
			fileLimit: func() int {
				l, err := strconv.Atoi(envMap["APP_FILE_LIMIT"])
				if err != nil {
					log.Fatalf("Error parsing APP_FILE_LIMIT: %v", err)
				}
				return l
			}(),
			gcpBucket: envMap["APP_GCP_BUCKET"]},
		db: &db{
			host: envMap["DB_HOST"],
			port: func() int {
				p, err := strconv.Atoi(envMap["DB_PORT"])
				if err != nil {
					log.Fatalf("Error parsing DB_PORT: %v", err)
				}
				return p
			}(),
			protocol: envMap["DB_PROTOCOL"],
			username: envMap["DB_USERNAME"],
			password: envMap["DB_PASSWORD"],
			database: envMap["DB_DATABASE"],
			sslMode:  envMap["DB_SSL_MODE"],
			maxConnections: func() int {
				c, err := strconv.Atoi(envMap["DB_MAX_CONNECTIONS"])
				if err != nil {
					log.Fatalf("Error parsing DB_MAX_CONNECTIONS: %v", err)
				}
				return c
			}(),
		},
		jwt: &jwt{
			adminKey:  envMap["JWT_ADMIN_KEY"],
			secretKey: envMap["JWT_SECRET_KEY"],
			apiKey:    envMap["JWT_API_KEY"],
			accessExpAt: func() int {
				t, err := strconv.Atoi(envMap["JWT_ACCESS_EXPIRES"])
				if err != nil {
					log.Fatalf("Error parsing JWT_ACCESS_EXPIRES: %v", err)
				}
				return t
			}(),
			refreshExpAt: func() int {
				t, err := strconv.Atoi(envMap["JWT_REFRESH_EXPIRES"])
				if err != nil {
					log.Fatalf("Error parsing JWT_REFRESH_EXPIRES: %v", err)
				}
				return t
			}(),
		},
	}
}

type IConfig interface {
	APP() IAppConfig
	DB() IDBConfig
	JWT() IJwtConfig
}

type config struct {
	app *app
	db  *db
	jwt *jwt
}

func (c *config) APP() IAppConfig {
	return c.app
}

func (c *config) DB() IDBConfig {
	return c.db
}

func (c *config) JWT() IJwtConfig {
	return c.jwt
}

type IAppConfig interface {
	Url() string
	Name() string
	Version() string
	ReadTimeOut() time.Duration
	WriteTimeOut() time.Duration
	BodyLimit() int
	Filelimit() int
	GCPBucket() string
}

func (a *app) Url() string {
	return fmt.Sprintf("%s:%d", a.host, a.port)
}
func (a *app) Name() string                { return a.name }
func (a *app) Version() string             { return a.version }
func (a *app) ReadTimeOut() time.Duration  { return a.readTimeout }
func (a *app) WriteTimeOut() time.Duration { return a.writeTimeout }
func (a *app) BodyLimit() int              { return a.bodyLimit }
func (a *app) Filelimit() int              { return a.fileLimit }
func (a *app) GCPBucket() string           { return a.gcpBucket }

type app struct {
	host         string
	port         int
	name         string
	version      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	bodyLimit    int
	fileLimit    int
	gcpBucket    string
}

type IDBConfig interface {
	Url() string
	MaxConnections() int
}

func (db *db) Url() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", db.host, db.port, db.username, db.password, db.database, db.sslMode)
}
func (db *db) MaxConnections() int {
	return db.maxConnections
}

type db struct {
	host           string
	port           int
	protocol       string
	username       string
	password       string
	database       string
	sslMode        string
	maxConnections int
}

type IJwtConfig interface {
	SecretKey() []byte
	AdminKey() []byte
	APIKey() []byte
	AccessExpiresAt() int
	RefreshExpiresAt() int
	SetJwtAccessExpires(t int)
	SetJwtRefreshExpires(t int)
}

func (j *jwt) SecretKey() []byte {
	return []byte(j.secretKey)
}

func (j *jwt) AdminKey() []byte {
	return []byte(j.adminKey)
}

func (j *jwt) APIKey() []byte {
	return []byte(j.apiKey)
}

func (j *jwt) AccessExpiresAt() int {
	return j.accessExpAt
}

func (j *jwt) RefreshExpiresAt() int {
	return j.refreshExpAt
}

func (j *jwt) SetJwtAccessExpires(t int) {
	j.accessExpAt = t
}

func (j *jwt) SetJwtRefreshExpires(t int) {
	j.refreshExpAt = t
}

type jwt struct {
	adminKey     string
	secretKey    string
	apiKey       string
	accessExpAt  int //sec
	refreshExpAt int // sec
}
