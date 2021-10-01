package cmd

import (
	"bufio"
	"io"
	"os"

	"github.com/jordyv/reconstore/internal/entities"
	"github.com/sirupsen/logrus"
)

type TagSubdomainsCmd struct{}

func (c *TagSubdomainsCmd) ShouldHandle() bool {
	return tagSubdomainsCmd.Used
}

func (c *TagSubdomainsCmd) Handle() {
	s, _ := os.Stdin.Stat()
	if s.Size() == 0 {
		logrus.Info("No input received")
		return
	}

	buf := bufio.NewReader(os.Stdin)
	for {
		l, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			logrus.Fatal(err.Error())
		}

		var subdomain entities.Subdomain
		if r := db.Where("domain = ?", string(l)).First(&subdomain); r.Error == nil {

			tagEntities := make([]*entities.Tag, 0)
			for _, t := range tags {
				if t != "" {
					tag := entities.Tag{Name: t}
					db.Where("name = ?", t).FirstOrCreate(&tag)
					tagEntities = append(tagEntities, &tag)
				}
			}
			subdomain.Tags = append(subdomain.Tags, tagEntities...)
			db.Save(&subdomain)
		}
	}
}
