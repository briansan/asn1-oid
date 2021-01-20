package main

import (
	"context"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"

	"asn1-oid/internal/collector"
)

const (
	urlBufSize = 10000
	colBufSize = 100
	sampleRate = 100
	startURL   = "http://oid-info.com/get/0"
	timeout    = 10 * time.Second
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	// Make URL channel.
	urlCh := make(chan string, 1000)
	defer close(urlCh)

	// Make Collector channel.
	collectorPool := make([]*colly.Collector, colBufSize)
	for i := 0; i < colBufSize; i++ {
		collectorPool[i] = collector.New(urlCh)
	}
	colCh := make(chan *colly.Collector, urlBufSize)
	defer close(colCh)
	for _, col := range collectorPool {
		colCh <- col
	}

	// Throw initial url and run Collectors on URL channel.
	urlCh <- startURL
	ctx, more := context.Background(), true

	start := time.Now()
	for i := 0; more; i++ {
		time.Sleep(100 * time.Millisecond)
		if i%sampleRate == sampleRate-1 {
			now := time.Now()
			rate := now.Sub(start) / sampleRate
			logrus.WithFields(logrus.Fields{
				"i": i, "rate": rate.String()}).Info("healthcheck")
			start = now
		}

		select {
		case <-ctx.Done():
			logrus.Info("context canceled")
			more = false
		case next := <-urlCh:
			c := <-colCh
			c.Visit(next)
			colCh <- c
		case <-time.After(timeout):
			logrus.Info("no more urls")
			more = false
		}
	}

	logrus.Info("Done")
}
