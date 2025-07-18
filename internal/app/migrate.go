package app

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	_defaultMaxAttempts = 20
	_defaultMaxTimeout  = time.Second * 2
)

func migrateDB(URL string) error {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return fmt.Errorf("parsing db url: %v", err)
	}

	q := parsedURL.Query()
	q.Set("sslmode", "disable")
	parsedURL.RawQuery = q.Encode()
	URL = parsedURL.String()

	var m *migrate.Migrate
	for attempt := _defaultMaxAttempts; attempt > 0; attempt-- {
		m, err = migrate.New("file://migrations", URL)
		if err == nil {
			break
		}

		slog.Debug("trying to connect to db", slog.Int("attempt", attempt))

		time.Sleep(_defaultMaxTimeout)
	}

	if err != nil {
		return fmt.Errorf("connecting to db: %v", err)
	}

	err = m.Up()
	defer m.Close()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("applying up migration: %v", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		slog.Info("no new change")
		return nil
	}

	slog.Info("sucsess")
	return nil
}
