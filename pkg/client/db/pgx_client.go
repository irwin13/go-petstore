package db

import (
	"crypto/tls"
	"database/sql"
	"errors"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/irwin13/go-petstore/internal/config"
	"github.com/irwin13/go-petstore/pkg/logger"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type PgxClient struct {
	logger    *zap.Logger
	appConfig config.AppConfig
}

func NewPgxClient(appConfig config.AppConfig) *PgxClient {
	return &PgxClient{
		logger:    logger.GetAppLogger(),
		appConfig: appConfig,
	}
}

var (
	pgxConnPool *pgx.ConnPool = nil
)

func (c *PgxClient) Start() error {
	c.logger.Info("Start initializing db connection ..")

	if pgxConnPool == nil {
		var err error
		pgxConnPool, err = pgx.NewConnPool(c.dbConfig())
		if err != nil {
			c.logger.Error("Error when connect to database",
				zap.String("error", err.Error()),
			)
			return err
		}

		c.logger.Info("Success start DbClient",
			zap.Int("AvailableConnections", pgxConnPool.Stat().AvailableConnections),
			zap.Int("CurrentConnections", pgxConnPool.Stat().CurrentConnections),
			zap.Int("MaxConnections", pgxConnPool.Stat().MaxConnections),
		)
	}
	return nil
}

func (c *PgxClient) Shutdown() error {
	c.logger.Info("Closing DB connection pool ...")
	pgxConnPool.Close()
	c.logger.Info("Close DB connection pool finished")
	return nil
}

func (c *PgxClient) GetConnection() (interface{}, error) {
	if pgxConnPool == nil {
		return nil, errors.New("connection pool is nil")
	}
	return pgxConnPool, nil
}

func (c *PgxClient) runValidationQuery() int {
	var val int
	rows, err := pgxConnPool.Query(c.appConfig.DbValidationQuery)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&val)
		if err != nil {
			c.logger.Error("Error when connect to database",
				zap.String("error", err.Error()),
			)
		}
		c.logger.Debug("validationQuery", zap.Int("result", val))
	}

	return val
}

func (c *PgxClient) RunMigration() error {
	c.logger.Info("Running DB Migration...")
	var databaseURL = c.dbURL()

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		c.logger.Error("Error db migration connecting to db", zap.String("message", err.Error()))
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		c.logger.Error("Error db migration load driver", zap.String("message", err.Error()))
		return err
	}

	dbMigrationPath := c.appConfig.DbMigrationPath

	m, err := migrate.NewWithDatabaseInstance(
		dbMigrationPath,
		"postgres", driver)
	if err != nil {
		c.logger.Error("Error db migration load migration file",
			zap.String("message", err.Error()),
			zap.String("dbMigrationPath", dbMigrationPath),
		)
		return err
	}

	migrationVersion := c.appConfig.DbMigrationVersion

	err = m.Steps(migrationVersion)
	if err != nil {
		c.logger.Warn("Warn db migration when upgrading version, but we can probably ignore it",
			zap.String("message", err.Error()),
			zap.Int("migrationVersion", migrationVersion),
		)
	}

	c.logger.Info("Success running DB Migration", zap.Int("migrationVersion", migrationVersion))
	return nil
}

func (c *PgxClient) dbURL() string {

	databaseURL := c.appConfig.DbURL

	if c.appConfig.DbUseSSL {
		databaseURL = databaseURL + "?sslmode=require"
	} else {
		databaseURL = databaseURL + "?sslmode=disable"
	}

	c.logger.Debug("DB URL", zap.String("url", databaseURL))

	return databaseURL
}

func (c *PgxClient) dbConfig() pgx.ConnPoolConfig {

	maxConnection := c.appConfig.DbMaxConn
	acquireTimeout := c.appConfig.DbAcquireTimeout

	var databaseURL = c.dbURL()

	c.logger.Info("DB Config",
		zap.Bool("useSsl", c.appConfig.DbUseSSL),
		zap.Int("maxConnection", maxConnection),
		zap.Int("acquireTimeout", acquireTimeout),
	)

	var config pgx.ConnPoolConfig

	parsedURL, err := url.Parse(databaseURL)
	if err != nil {
		c.logger.Error("Invalid db url",
			zap.String("dbURL", databaseURL),
		)
		log.Fatalf("Fatal. invalid dbURL %s", databaseURL)
	}

	pass, _ := parsedURL.User.Password()
	port, _ := strconv.ParseUint(parsedURL.Port(), 10, 16)

	config.Host = parsedURL.Hostname()
	config.Port = uint16(port)
	config.User = parsedURL.User.Username()
	config.Password = pass
	config.Database = strings.TrimPrefix(parsedURL.Path, "/")
	config.MaxConnections = maxConnection
	config.AcquireTimeout = time.Duration(acquireTimeout) * time.Second

	if c.appConfig.DbUseSSL {
		config.ConnConfig.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return config

}
