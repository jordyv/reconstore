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

	var newCount = 0
	var updateCount = 0

	wg := sync.WaitGroup{}
	for i := 0; i < saveThreads; i++ {
		wg.Add(1)

		go func() {
			for line := range lines {
				domain := strings.TrimSpace(strings.ToLower(line))

				var countExisting int64
				db.Model(&entities.Subdomain{}).Where("domain = ?", domain).Count(&countExisting)

				var newSubdomain bool
				var s entities.Subdomain
				if countExisting == 0 {
					newSubdomain = true
					s = entities.Subdomain{
						Domain: domain,
					}
				} else if countExisting > 0 && update {
					newSubdomain = false
					db.Model(&entities.Subdomain{}).Joins("HTTPInfo").Joins("DNSInfo").Where("domain = ?", domain).Find(&s)
				} else {
					continue
				}
				if program.ID != 0 {
					s.Program = &program
				}

				if newSubdomain {
					db.Create(&s)
					newCount++
				} else {
					db.Save(&s)
					updateCount++
				}

				hooks.TriggerAfterSubdomainSave(&s)

				logrus.Infof("Saved %s", domain)
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

	logrus.Infof("Stored %d new subdomains, updated %d subdomains", newCount, updateCount)
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
