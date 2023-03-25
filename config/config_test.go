package config

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	// setup
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("JWT_RSA", "pem")

	// cleanup
	defer func() {
		os.Clearenv()
	}()

	// call
	cfg, err := Get()

	// assertions
	if err != nil {
		t.Fatalf("Unable to parse config %s", err)
	}

	if cfg.Database.HOST != os.Getenv("DB_HOST") {
		t.Errorf("Unmatching tags with DB")
	}
	if cfg.Jwt.RSA != os.Getenv("JWT_RSA") {
		t.Errorf("Unmatching tags with jwt")
	}
}
