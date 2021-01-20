package collector

import (
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

const (
	visitChildMapArea   = "67,64,82,78"
	visitSiblingMapArea = "105,34,120,50"
)

func New(ch chan string) *colly.Collector {
	c := colly.NewCollector(colly.Async(true))

	// Extract OID data and figure out next traversal.
	c.OnHTML("html", func(e *colly.HTMLElement) {
		logger := logrus.WithField("current", e.Request.URL.Path)
		logger.Debug("visited")
		// TODO: Extract OID data.

		var childLink, siblingLink string

		// Figure out next traversal.
		e.ForEachWithBreak("map > area", func(i int, e *colly.HTMLElement) bool {
			switch e.Attr("coords") {
			case visitChildMapArea:
				childLink = e.Attr("href")
			case visitSiblingMapArea:
				siblingLink = e.Attr("href")
			}
			return true
		})

		// Visit child link.
		if childLink != "" {
			logger.WithField("child", childLink).Debug("pushing")
			ch <- childLink
		}
		if siblingLink != "" {
			logger.WithField("sibling", siblingLink).Debug("pushing")
			ch <- siblingLink
		}
	})

	return c
}
