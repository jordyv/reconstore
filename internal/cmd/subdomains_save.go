package cmd

import (
	"bufio"
	"io"
	"os"
	"sync"

	"github.com/jordyv/reconstore/internal/entities"
	"github.com/jordyv/reconstore/internal/hooks"
	"github.com/sirupsen/logrus"
)

type SaveSubdomainsCmd struct{}

func (c *SaveSubdomainsCmd) ShouldHandle() bool {
	return saveSubdomainsCmd.Used
}

func (c *SaveSubdomainsCmd) Handle() {
	c.validate()

	wg := sync.WaitGroup{}

	var program entities.Program
	db.Where("slug = ?", saveProgramSlug).Find(&program)

	var count = 0
	buf := bufio.NewReader(os.Stdin)
	for {
		l, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			logrus.WithError(err).Fatal("error while reading stdin")
		}
		domain := string(l)
		wg.Add(1)

		go func() {
			var countExisting int64
			db.Model(&entities.Subdomain{}).Where("domain = ?", domain).Count(&countExisting)
			if countExisting == 0 {
				s := &entities.Subdomain{
					Domain:  domain,
					Program: program,
				}
				db.Create(s)

				hooks.TriggerAfterSubdomainSave(s)

				logrus.Infof("Saved %s", domain)
				count++
			}
			wg.Done()
		}()
	}

	wg.Wait()

	logrus.Infof("Stored %d new subdomains", count)
}

func (c *SaveSubdomainsCmd) validate() {
	if saveProgramSlug == "" {
		logrus.Fatal("Program slug is required")
	}

	var count int64
	db.Model(&entities.Program{}).Where("slug = ?", saveProgramSlug).Count(&count)

	if count == 0 {
		logrus.Fatalf("No program with slug '%s' found", saveProgramSlug)
	}
}
