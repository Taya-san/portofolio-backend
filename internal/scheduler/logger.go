package scheduler

import (
	"github.com/Taya-san/portofolio-backend/pkg/logger"
	"github.com/go-co-op/gocron/v2"
)

type GocronLogger struct {
	L logger.Interface
}

var _ gocron.Logger = GocronLogger{}

func (g GocronLogger) Debug(msg string, args ...any) { g.L.Debug(msg, args) }
func (g GocronLogger) Error(msg string, args ...any) { g.L.Error(msg, args) }
func (g GocronLogger) Info(msg string, args ...any)  { g.L.Info(msg, args) }
func (g GocronLogger) Warn(msg string, args ...any)  { g.L.Warn(msg, args) }
