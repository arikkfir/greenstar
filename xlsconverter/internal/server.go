package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arikkfir/greenstar/xlsconverter/pkg"
	xlsxPkg "github.com/arikkfir/greenstar/xlsxprocessor/pkg"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"io"
	"os"
	"os/exec"
	"path"
)

type Server struct {
	Config Config
	Redis  *redis.Client
}

func NewServer(config Config, Redis *redis.Client) *Server {
	return &Server{
		Config: config,
		Redis:  Redis,
	}
}

func (s *Server) Run(ctx context.Context) error {
	logger := *zerolog.Ctx(ctx)
	subscriber := s.Redis.Subscribe(ctx, pkg.InputChannelName)
	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to receive message from Redis")
		}

		req := pkg.ConvertXLSFileToXLSXRequest{}
		if err := json.Unmarshal([]byte(msg.Payload), &req); err != nil {
			logger.Fatal().Err(err).Msg("Failed to unmarshal message")
		}

		if err := s.processConversionRequest(ctx, &req); err != nil {
			logger.Fatal().Err(err).Msg("Failed to process conversion request")
		}
	}
}

func (s *Server) processConversionRequest(ctx context.Context, req *pkg.ConvertXLSFileToXLSXRequest) error {
	logger := *zerolog.Ctx(ctx)

	dir, fileName := path.Split(req.FileName)
	ext := path.Ext(fileName)
	if ext != ".xls" {
		return fmt.Errorf("unsupported file extension: %s", ext)
	}

	srcFile, err := os.CreateTemp(s.Config.WorkDir, fileName[:len(fileName)-len(ext)]+"-*.xls")
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

	cmd := exec.CommandContext(ctx, "libreoffice", "--convert-to", "xlsx", "--outdir", dir, srcFile.Name())
	cmd.Stdout = logger.Level(zerolog.TraceLevel)
	cmd.Stderr = logger.Level(zerolog.WarnLevel)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to convert .xls to .xlsx: %w", err)
	}

	dstFile, err := os.Open(path.Join(s.Config.WorkDir, fileName[:len(fileName)-len(ext)]+".xlsx"))
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

	res := xlsxPkg.ProcessXLSXRequest{
		FileName: path.Base(dstFile.Name()),
		Data:     data,
	}
	if payload, err := json.Marshal(&res); err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	} else if err := s.Redis.Publish(ctx, pkg.TargetChannelName, payload).Err(); err != nil {
		return fmt.Errorf("failed to publish response: %w", err)
	} else {
		return nil
	}
}
