package hooks

import (
	"sync"

	"github.com/jordyv/reconstore/internal/entities"
	"github.com/jordyv/reconstore/internal/hooks/subdomainimport"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Hook interface {
	AfterSave(*entities.Subdomain) error
}

var (
	afterSaveHooks []Hook
)

func Init(db *gorm.DB) {
	afterSaveHooks = []Hook{
		subdomainimport.NewDNS(db),
		subdomainimport.NewHTTP(db),
	}
}

func TriggerAfterSubdomainSave(s *entities.Subdomain) {
	wg := sync.WaitGroup{}

	for _, h := range afterSaveHooks {
		wg.Add(1)
		go func(hook Hook) {
			if err := hook.AfterSave(s); err != nil {
				logrus.WithError(err).Error("error in after save hook")
			}
			wg.Done()
		}(h)
	}
	wg.Wait()
}
