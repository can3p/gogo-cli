package dbgen

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/can3p/gogo-cli/pkg/config"
	"github.com/can3p/gogo/testcontainers/postgres"
)

const (
	DefaultDockerImage = "ghcr.io/can3p/gogo-docker-deps:latest"
)

// Generator handles sqlboiler code generation
type Generator struct {
	ProjectDir  string
	Config      *config.Config
	DockerImage string
}

// NewGenerator creates a new Generator instance
func NewGenerator(projectDir string, cfg *config.Config) *Generator {
	return &Generator{
		ProjectDir:  projectDir,
		Config:      cfg,
		DockerImage: DefaultDockerImage,
	}
}

// Run executes the code generation process:
// 1. Starts a postgres container
// 2. Applies migrations
// 3. Runs sqlboiler via docker to generate code
func (g *Generator) Run() error {
	migrationsDir := filepath.Join(g.ProjectDir, g.Config.Database.Migrations)

	// Check migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory not found: %s", migrationsDir)
	}

	fmt.Println("Starting PostgreSQL container...")
	testDB, err := postgres.NewTestDB(postgres.Options{
		MigrationsDir: migrationsDir,
	})
	if err != nil {
		return fmt.Errorf("failed to start test database: %w", err)
	}
	defer func() {
		fmt.Println("Cleaning up database...")
		_ = testDB.Close()
	}()

	fmt.Println("Database ready, generating sqlboiler config...")

	// Create temporary sqlboiler config
	configContent, err := g.generateSQLBoilerConfig(testDB.ConnInfo)
	if err != nil {
		return fmt.Errorf("failed to generate sqlboiler config: %w", err)
	}

	// Write config to temp file in project directory (so docker can access it)
	configPath := filepath.Join(g.ProjectDir, ".sqlboiler.tmp.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write sqlboiler config: %w", err)
	}
	defer os.Remove(configPath)

	fmt.Println("Running sqlboiler...")

	// Run sqlboiler via docker
	if err := g.runSQLBoiler(configPath); err != nil {
		return fmt.Errorf("sqlboiler failed: %w", err)
	}

	fmt.Println("Code generation complete!")
	return nil
}

const sqlboilerConfigTemplate = `pkgname = "{{.PkgName}}"
output = "{{.Output}}"
wipe = {{.Wipe}}

[psql]
dbname = "{{.DBName}}"
host = "{{.Host}}"
port = {{.Port}}
user = "{{.User}}"
pass = "{{.Password}}"
sslmode = "{{.SSLMode}}"
{{- if .Blacklist}}
blacklist = [{{range $i, $v := .Blacklist}}{{if $i}}, {{end}}"{{$v}}"{{end}}]
{{- end}}
`

type sqlboilerConfigData struct {
	PkgName   string
	Output    string
	Wipe      bool
	DBName    string
	Host      string
	Port      string
	User      string
	Password  string
	SSLMode   string
	Blacklist []string
}

func (g *Generator) generateSQLBoilerConfig(connInfo postgres.ConnectionInfo) (string, error) {
	tmpl, err := template.New("sqlboiler").Parse(sqlboilerConfigTemplate)
	if err != nil {
		return "", err
	}

	data := sqlboilerConfigData{
		PkgName:   g.Config.Database.SQLBoiler.PkgName,
		Output:    g.Config.Database.SQLBoiler.Output,
		Wipe:      g.Config.Database.SQLBoiler.Wipe,
		DBName:    connInfo.DBName,
		Host:      connInfo.Host,
		Port:      connInfo.Port,
		User:      connInfo.User,
		Password:  connInfo.Password,
		SSLMode:   connInfo.SSLMode,
		Blacklist: g.Config.Database.SQLBoiler.Blacklist,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (g *Generator) runSQLBoiler(configPath string) error {
	// Get relative config path for docker
	relConfigPath, err := filepath.Rel(g.ProjectDir, configPath)
	if err != nil {
		relConfigPath = filepath.Base(configPath)
	}

	args := []string{
		"run", "--rm",
		"-v", fmt.Sprintf("%s:/workspace", g.ProjectDir),
		"--network", "host",
		g.DockerImage,
		"--add-panic-variants",
		"--no-hooks",
		"--no-tests",
		"--add-enum-types",
		"-c", relConfigPath,
		"psql",
	}

	cmd := exec.Command("docker", args...)
	cmd.Dir = g.ProjectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
