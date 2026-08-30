package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"pyra/pkg/migrate"

	"github.com/joho/godotenv"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	// FIX: load depending on env
	if err := godotenv.Load("/Users/sagg0t/devel/pyra/.env.test"); err != nil {
		panic(err)
	}
	config := migrate.NewConfig()

	dbConfig := migrate.NewDBConfig("pgx")
	dbConfig.Attrs.Add("sslmode", fetchEnv("DB_SSLMODE", "disable"))

	var command string
	if len(os.Args) == 1 {
		command = "apply"
	} else {
		command = os.Args[1]
	}

	engine, err := migrate.NewEngine(config, dbConfig)
	if err != nil {
		panic(err)
	}

	var commandError error = nil
	switch command {
	case "add":
		commandError = addMigration(engine)
	case "apply":
		var count uint64 = 0
		if len(os.Args) > 2 {
			argCount, err := parseCount(os.Args[2])
			if err != nil {
				fmt.Println(err)
				commandError = err
				break
			}

			count = argCount
		}

		commandError = applyMigration(engine, count)
	case "rollback":
		var count uint64 = 1
		if len(os.Args) > 2 {
			argCount, err := parseCount(os.Args[3])
			if err != nil {
				commandError = err
				break
			}

			count = argCount
		}

		commandError = rollbackMigration(engine, count)
	case "status":
	case "version":
		commandError = migrationVersion(engine)
	case "help":
		commandError = showHelp()
	default:
		commandError = fmt.Errorf("unknown command %q", command)
	}

	if commandError != nil {
		fmt.Println(fmt.Errorf("migrate: %w", commandError))
	}
}

func fetchEnv(name, fallback string) string {
	val, ok := os.LookupEnv(name)
	if !ok {
		return fallback
	}

	return val
}

func addMigration(engine *migrate.Engine) error {
	if len(os.Args) < 3 {
		return errors.New("missing migration name")
	}

	name := os.Args[2]

	fileNames, err := engine.CreateMigration(name)
	if err != nil {
		return fmt.Errorf("failed to create migration files - %w", err)
	}

	fmt.Printf("created: %s\n", fileNames[0])
	fmt.Printf("created: %s\n", fileNames[1])

	return nil
}

func applyMigration(engine *migrate.Engine, n uint64) error {
	migs, err := engine.Apply(n)
	if err != nil {
		return err
	}

	for _, m := range migs {
		slog.Info("successfully applied a migration", "version", m.Version)
	}

	return nil
}

func rollbackMigration(engine *migrate.Engine, n uint64) error {
	migs, err := engine.Rollback(n)
	if err != nil {
		return fmt.Errorf("rollback failure - %w", err)
	}

	for _, m := range migs {
		fmt.Printf("rollback success: %s", m.Version)
	}

	return nil
}

func migrationVersion(engine *migrate.Engine) error {
	mig, err := engine.CurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to retireve the current migration version - %w", err)
	}

	if mig.Version == "0" {
		fmt.Println("no DB migrations have been applied so far")
	} else {
		fmt.Println(mig.Version)
	}

	return nil
}

func parseCount(sCount string) (uint64, error) {
	count, err := strconv.ParseUint(sCount, 10, strconv.IntSize)
	if err != nil {
		return 0, fmt.Errorf("failed to read migrations N - %w", err)
	}

	return count, nil
}

var helpMsg = "Migrate:\n" +
	"\tadd NAME\tcreate migration files with the given NAME\n" +
	"\tapply [N]\tapply N migrations, defaults to 0 (all pending migrations)\n" +
	"\trollback [N]\trollback N migrations, defaults to 1\n" +
	"\tstatus\t\tshow migration status\n" +
	"\tversion\t\tshow last applied migration version\n" +
	"\thelp\t\tshow this message\n"

func showHelp() error {
	fmt.Print(helpMsg)
	return nil
}
