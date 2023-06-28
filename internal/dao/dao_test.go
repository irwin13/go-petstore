package dao_test

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/irwin13/go-petstore/internal/config"
	"github.com/irwin13/go-petstore/internal/dao/imp"
	"github.com/irwin13/go-petstore/pkg/client"
	"github.com/irwin13/go-petstore/pkg/client/db"
	"github.com/irwin13/go-petstore/pkg/logger"
	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

var appLogger *zap.Logger

var pgxClient client.DbClient

var petDao *imp.PetDaoPgx

func TestMain(m *testing.M) {

	appLogger = logger.InitAppLogger()
	defer appLogger.Sync()

	appConfig, err := config.LoadAppConfig()
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	// check to make sure the db integration test is pointing to localhost or database with username or db name contains integration_test
	if !strings.Contains(appConfig.DbURL, "integration_test") && !strings.Contains(appConfig.DbURL, "localhost") {
		log.Fatalf("Database used is not test database")
		return
	}

	pgxClient = db.NewPgxClient(appConfig)

	err = pgxClient.Start()
	if err != nil {
		log.Fatalf("Failed init DB connection pool : %s", err.Error())
	}

	pgxClient.RunMigration()

	petDao = imp.NewPetDaoPgx(pgxClient)

	exitCode := m.Run()

	pgxClient.Shutdown()

	os.Exit(exitCode)
}

func executeRawSql(sqlList []string) error {

	conn, err := pgxClient.GetConnection()
	if err != nil {
		return err
	}

	pgxConn := conn.(*pgx.ConnPool)

	tx, err := pgxConn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, sql := range sqlList {
		appLogger.Info("execute raw sql", zap.String("sql", sql))
		_, err = tx.Exec(sql)
	}

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
