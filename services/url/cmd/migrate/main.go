//go:build ignore
// +build ignore

// Install atlas
// curl -sSf https://atlasgo.sh | sh

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tuan-dd/go-pkg/common"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/generated/ent/migrate"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/global"
)

const (
	migrationDir = "./migrations"
)

var (
	defaultDBURL = ""
	backupDBURL  = ""
)

var flags = flag.NewFlagSet("migrate", flag.ExitOnError)

type AppConfig struct{}

func main() {
	flags.Usage = usage
	_ = flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) == 0 {
		flags.Usage()
		return
	}
	configDb, err := common.LoadConfig[global.Config]("")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	defaultDBURL = fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s", configDb.SQL.RDBMS, configDb.SQL.Username, configDb.SQL.Password, configDb.SQL.Host, configDb.SQL.Port, configDb.SQL.DBname)

	backupDBURL = os.Getenv("BK_DATABASE_URL")

	switch args[0] {
	case "ent_add":
		addEntMigrate(args)
	case "add":
		migrateAdd(args)
	case "up":
		migrateApply()
	case "down":
		migrateDown(args)
	case "status":
		migrateStatus()
	case "hash":
		migrateHash()
	default:
		log.Fatalf("Unknown command: %s", args[0])
	}
}

const usageCommands = `
Commands:
    up       - Apply all available migrations
    down     - Roll back to a specific migration version
    status   - Show migration status
    hash     - Update schema checksum
    ent_add  - Generate migration file from schema
    add      - Create an empty migration file
`

func usage() {
	fmt.Println("Usage: migrate COMMAND")
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

func runCommand(command string, args ...string) {
	cmd := exec.Command("atlas", append([]string{"migrate", command}, args...)...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	fmt.Println("Executing:", cmd.String())
	if err := cmd.Run(); err != nil {
		log.Fatalf("Command failed: %v", err)
		os.Exit(1)
	}
}

func ensureMigrationDir() {
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		log.Fatalf("Failed to create migration directory: %v", err)
	}
}

func addEntMigrate(args []string) {
	if len(args) < 2 || args[1] == "" {
		log.Fatal("Migration name is required")
	}

	ensureMigrationDir()

	localDir, err := atlas.NewLocalDir(migrationDir)
	if err != nil {
		log.Fatalf("Failed to create atlas migration directory: %v", err)
	}

	opts := []schema.MigrateOption{
		schema.WithDir(localDir),
		schema.WithMigrationMode(schema.ModeInspect),
		schema.WithDialect(dialect.MySQL),
		schema.WithFormatter(atlas.DefaultFormatter),
		schema.WithIndent("  "),
		schema.WithDropIndex(true),
		schema.WithDropColumn(true),
	}

	ctx := context.Background()
	if err := migrate.NamedDiff(ctx, defaultDBURL, args[1], opts...); err != nil {
		log.Fatalf("Migration file generation failed: %v", err)
	}
}

func migrateAdd(args []string) {
	if len(args) < 2 || args[1] == "" {
		log.Fatal("Migration name is required")
	}
	runCommand("new", "--dir", "file://"+migrationDir, args[1])
}

func migrateApply() {
	runCommand("apply", "--dir", "file://"+migrationDir, "--url", defaultDBURL)
}

func migrateDown(args []string) {
	if len(args) < 2 || args[1] == "" {
		log.Fatal("Version is required")
	}
	runCommand(
		"down",
		"--to-version", args[1],
		"--url", defaultDBURL,
		"--dir", "file://"+migrationDir,
		"--dev-url", backupDBURL,
	)
}

func migrateHash() {
	runCommand("hash", "--dir", "file://"+migrationDir)
}

func migrateStatus() {
	runCommand("status", "--url", defaultDBURL, "--dir", "file://"+migrationDir)
}
