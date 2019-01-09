package main

import (
	"github.com/go-kit/kit/log"
	"time"
)

const TimeFormat  = "2006-01-02T15:04:05Z07:00"

type loggingMiddleware struct {
	logger log.Logger
	next   ReportingService
}

func (mw loggingMiddleware) reporting (req BusinessesRequest) (output Business, err error) {
	defer func(begin time.Time) {
		callTime := time.Now().Format(TimeFormat)
		mw.logger.Log(
			"callTime", callTime,
			"method", "reporting",
			"Limit", req.Limit,
			"Offset", req.Offset,
			"BusinessID", req.BusinessID,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.reporting(req)
	return
}

