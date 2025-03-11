package pkg

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger     *slog.Logger
	loggerOnce sync.Once
)

type LoggerConfig struct {
	FileName         string `koanf:"file_name"`
	Directory        string `koanf:"directory"`
	UseLocalTime     bool   `koanf:"use_local_time"`
	FileMaxSizeInMB  int    `koanf:"file_max_size_in_mb"`
	FileMaxAgeInDays int    `koanf:"file_max_age_in_days"`
}

func InitLogger(config LoggerConfig) error {
	var err error
	loggerOnce.Do(func() {
		cErr := os.Mkdir(config.Directory, 0644)
		if cErr != nil && !os.IsExist(cErr) {
			err = cErr
			return
		}
		workingDir, gErr := os.Getwd()
		if err != nil {
			err = gErr
			return
		}
		fileWriter := &lumberjack.Logger{
			Filename:  filepath.Join(filepath.Join(workingDir, config.Directory), config.FileName),
			LocalTime: config.UseLocalTime,
			MaxSize:   config.FileMaxSizeInMB,
			MaxAge:    config.FileMaxAgeInDays,
		}
		Logger = slog.New(
			slog.NewJSONHandler(io.MultiWriter(fileWriter, os.Stdout), &slog.HandlerOptions{}),
		)
	})
	if err != nil {
		return err
	}
	return nil
}
