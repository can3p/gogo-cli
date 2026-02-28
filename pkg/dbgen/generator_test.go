package dbgen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/can3p/gogo-cli/pkg/config"
	"github.com/can3p/gogo/testcontainers/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	code := m.Run()
	_ = postgres.Cleanup()
	if code != 0 {
		os.Exit(code)
	}
}

func TestGenerateSQLBoilerConfig(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Migrations: "migrations",
			SQLBoiler: config.SQLBoilerConfig{
				PkgName:   "models",
				Output:    "pkg/models",
				Wipe:      true,
				Blacklist: []string{"migrations", "http_sessions"},
			},
		},
	}

	g := NewGenerator("/test/project", cfg)

	connInfo := postgres.ConnectionInfo{
		Host:     "localhost",
		Port:     "5432",
		User:     "testuser",
		Password: "secret",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	configContent, err := g.generateSQLBoilerConfig(connInfo)
	require.NoError(t, err)

	assert.Contains(t, configContent, `pkgname = "models"`)
	assert.Contains(t, configContent, `output = "pkg/models"`)
	assert.Contains(t, configContent, `wipe = true`)
	assert.Contains(t, configContent, `dbname = "testdb"`)
	assert.Contains(t, configContent, `host = "localhost"`)
	assert.Contains(t, configContent, `port = 5432`)
	assert.Contains(t, configContent, `user = "testuser"`)
	assert.Contains(t, configContent, `pass = "secret"`)
	assert.Contains(t, configContent, `blacklist = ["migrations", "http_sessions"]`)
}

func TestConfigLoad(t *testing.T) {
	testdataDir, err := filepath.Abs("testdata")
	require.NoError(t, err)

	cfg, err := config.Load(testdataDir)
	require.NoError(t, err)

	assert.Equal(t, "migrations", cfg.Database.Migrations)
	assert.Equal(t, "models", cfg.Database.SQLBoiler.PkgName)
	assert.Equal(t, "pkg/models", cfg.Database.SQLBoiler.Output)
	assert.True(t, cfg.Database.SQLBoiler.Wipe)
	assert.Contains(t, cfg.Database.SQLBoiler.Blacklist, "migrations")
}

func TestGeneratorRun(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Get absolute path to testdata
	testdataDir, err := filepath.Abs("testdata")
	require.NoError(t, err)

	// Create a temp directory for output
	tmpDir, err := os.MkdirTemp("", "dbgen-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Copy testdata to temp dir
	err = copyDir(testdataDir, tmpDir)
	require.NoError(t, err)

	cfg, err := config.Load(tmpDir)
	require.NoError(t, err)

	// Use local docker image for testing
	g := NewGenerator(tmpDir, cfg)
	g.DockerImage = "gogo-docker-deps:test"

	err = g.Run()
	require.NoError(t, err)

	// Verify generated files exist
	outputDir := filepath.Join(tmpDir, cfg.Database.SQLBoiler.Output)
	assert.DirExists(t, outputDir)

	// Check for expected generated files
	usersFile := filepath.Join(outputDir, "users.go")
	assert.FileExists(t, usersFile)

	postsFile := filepath.Join(outputDir, "posts.go")
	assert.FileExists(t, postsFile)

	// Verify the generated code contains expected content
	usersContent, err := os.ReadFile(usersFile)
	require.NoError(t, err)
	assert.Contains(t, string(usersContent), "package models")
	assert.Contains(t, string(usersContent), "type User struct")
}

func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			data, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}
			if err := os.WriteFile(dstPath, data, 0644); err != nil {
				return err
			}
		}
	}
	return nil
}
