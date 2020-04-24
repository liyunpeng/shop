package util

import (
	"context"
	"github.com/kataras/golog"
)

var Logger *golog.Logger

var Ctx context.Context
var Cancel context.CancelFunc
