package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/arikkfir/greenstar/common"
	"github.com/arikkfir/greenstar/xlsconverter/pkg"
	"github.com/rs/zerolog"
	"github.com/rueian/rueidis"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type Config struct {
	General common.GeneralConfig
	Redis   RedisConfig `group:"redis" namespace:"redis" env-namespace:"REDIS"`
	WorkDir string      `env:"WORK_DIR" long:"work-dir" description:"Working directory for temporary files" default:"/tmp"`
}

type RedisConfig struct {
	Host     string `env:"HOST" value-name:"HOST" long:"host" description:"Redis host name" default:"localhost"`
	Port     int    `env:"PORT" value-name:"PORT" long:"port" description:"Redis port" default:"6379"`
	PoolSize int    `env:"POOL_SIZE" value-name:"POOL_SIZE" long:"pool-size" description:"Redis connection pool size" default:"3"`
}

func main() {
	config := Config{}
	common.ReadConfig(&config)
	config.General.Apply()

	// Create a context which will be cancelled on termination
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := *zerolog.Ctx(ctx)
	logger.Info().Interface("config", config).Msg("Configured")

	// Create the Redis client
	redisClient, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress:      []string{config.Redis.Host + ":" + strconv.Itoa(config.Redis.Port)},
		ClientName:       "greenstar-xlsconverter",
		BlockingPoolSize: config.Redis.PoolSize,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	defer redisClient.Close()

	// Subscribe to Redis queue
	err = redisClient.Receive(ctx, redisClient.B().Subscribe().Channel(pkg.ChannelName).Build(), func(msg rueidis.PubSubMessage) {
		req := pkg.ConvertXLSFileToXLSXRequest{}
		if err := json.Unmarshal([]byte(msg.Message), &req); err != nil {
			logger.Fatal().Err(err).Msg("Failed to unmarshal message")
		}

		if err := processConversionRequest(ctx, &config, redisClient, &req); err != nil {
			logger.Fatal().Err(err).Msg("Failed to process conversion request")
		}
	})
	if err != nil && !errors.Is(err, rueidis.ErrClosing) {
		logger.Fatal().Err(err).Msg("Failed to subscribe to Redis queue")
	}
}

func processConversionRequest(ctx context.Context, config *Config, redisClient rueidis.Client, req *pkg.ConvertXLSFileToXLSXRequest) error {
	logger := *zerolog.Ctx(ctx)

	var fileName = req.FileName
	if strings.Contains(fileName, "/") {
		_, fileName = path.Split(fileName)
	}
	ext := path.Ext(fileName)
	if ext != ".xls" {
		return fmt.Errorf("unsupported file extension: %s", ext)
	}

	srcFile, err := os.CreateTemp(config.WorkDir, fileName[:len(fileName)-len(ext)]+"-*.xls")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(srcFile.Name())
	defer srcFile.Close()

	if _, err := srcFile.Write(req.Data); err != nil {
		return fmt.Errorf("failed to write file data: %w", err)
	}
	logger = logger.With().Str("xlsFile", srcFile.Name()).Logger()
	logger.Debug().Msg("Wrote XLS data")

	cmd := exec.CommandContext(ctx, "libreoffice", "--convert-to", "xlsx", "--outdir", config.WorkDir, srcFile.Name())
	cmd.Stdout = logger.Level(zerolog.TraceLevel)
	cmd.Stderr = logger.Level(zerolog.WarnLevel)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to convert .xls to .xlsx: %w", err)
	} else {
		logger.Info().Int("exitCode", cmd.ProcessState.ExitCode()).Msg("Conversion succeeded")
	}

	srcFileName := srcFile.Name()
	xlsxFileName := srcFileName[:len(srcFileName)-len(path.Ext(srcFileName))] + ".xlsx"
	dstFile, err := os.Open(xlsxFileName)
	if err != nil {
		return fmt.Errorf("failed to open converted file: %w", err)
	}
	logger = logger.With().Str("xlsxFile", dstFile.Name()).Logger()
	defer os.Remove(dstFile.Name())
	defer dstFile.Close()

	data, err := io.ReadAll(dstFile)
	if err != nil {
		return fmt.Errorf("failed to read converted file: %w", err)
	}

	res := pkg.ConvertXLSFileToXLSXResponse{
		FileName: path.Base(dstFile.Name()),
		Data:     data,
	}
	if payload, err := json.Marshal(&res); err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	} else if resp := redisClient.Do(ctx, redisClient.B().Publish().Channel(req.ReplyTo).Message(string(payload)).Build()); resp.Error() != nil {
		return fmt.Errorf("failed to publish response to '%s': %w", req.ReplyTo, resp.Error())
	} else {
		return nil
	}
}
