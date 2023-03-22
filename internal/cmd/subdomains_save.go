package cmd

import (
	"bufio"
	"os"
	"strings"
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

	var program entities.Program
	db.Where("slug = ?", saveProgramSlug).Find(&program)

	lines := make(chan string)

	var count = 0

	wg := sync.WaitGroup{}
	for i := 0; i < saveThreads; i++ {
		wg.Add(1)

		go func() {
			for line := range lines {
				domain := strings.TrimSpace(strings.ToLower(line))

				var countExisting int64
				db.Model(&entities.Subdomain{}).Where("domain = ?", domain).Count(&countExisting)
				if countExisting == 0 {
					s := &entities.Subdomain{
						Domain: domain,
					}
					if program.ID != 0 {
						s.Program = &program
					}
					db.Create(s)

					hooks.TriggerAfterSubdomainSave(s)

					logrus.Infof("Saved %s", domain)
					count++
				}
			}
			wg.Done()
		}()
	}

	buf := bufio.NewScanner(os.Stdin)
	for buf.Scan() {
		lines <- buf.Text()
	}

	close(lines)
	wg.Wait()

	logrus.Infof("Stored %d new subdomains", count)
}

func (c *SaveSubdomainsCmd) validate() {
	if saveProgramSlug != "" {
		var count int64
		db.Model(&entities.Program{}).Where("slug = ?", saveProgramSlug).Count(&count)

		if count == 0 {
			logrus.Fatalf("No program with slug '%s' found", saveProgramSlug)
		}
	}
}
