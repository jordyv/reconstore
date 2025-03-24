package subdomainimport

import (
	"strings"

	"github.com/jordyv/reconstore/internal/dns"
	"github.com/jordyv/reconstore/internal/entities"
	dnslib "github.com/miekg/dns"

	"gorm.io/gorm"
)

var (
	// Borrowed from https://github.com/projectdiscovery/public-bugbounty-programs/blob/main/pkg/dns/dns.go#L19
	dnsServers = []string{
		"1.1.1.1:53",        // Cloudflare primary
		"1.0.0.1:53",        // Cloudflare secondary
		"8.8.8.8:53",        // Google primary
		"8.8.4.4:53",        // Google secondary
		"9.9.9.9:53",        // Quad9 Primary
		"9.9.9.10:53",       // Quad9 Secondary
		"77.88.8.8:53",      // Yandex Primary
		"77.88.8.1:53",      // Yandex Secondary
		"208.67.222.222:53", // OpenDNS Primary
		"208.67.220.220:53", // OpenDNS Secondary
	}
)

type DNS struct {
	db        *gorm.DB
	dnsClient dns.Client
}

func NewDNS(db *gorm.DB) *DNS {
	return &DNS{
		db:        db,
		dnsClient: dns.NewFallbackClient(dnsServers),
	}
}

func (d *DNS) AfterSave(s *entities.Subdomain) error {
	dnsInfo := entities.DNSInfo{}

	cnameResult, err := d.query(s.Domain, dnslib.TypeCNAME)
	if err != nil {
		return err
	}
	if len(cnameResult.Answer) > 0 {
		dnsInfo.CnameRecord = cnameResult.Answer[0].String()
	}

	aResult, err := d.query(s.Domain, dnslib.TypeA)
	if err != nil {
		return err
	}
	dnsInfo.Status = dnslib.RcodeToString[aResult.Rcode]
	if len(aResult.Answer) > 0 {
		aRecords := []string{}
		for _, a := range aResult.Answer {
			aRecords = append(aRecords, a.String())
		}
		dnsInfo.ARecords = strings.Join(aRecords, ",")
	}

	nsResult, err := d.query(s.Domain, dnslib.TypeNS)
	if err != nil {
		return err
	}
	if len(nsResult.Answer) > 0 {
		nsRecords := []string{}
		for _, a := range nsResult.Answer {
			nsRecords = append(nsRecords, a.String())
		}
		dnsInfo.NSRecords = strings.Join(nsRecords, ",")
	}

	s.DNSInfo = &dnsInfo
	d.db.Save(s)

	return nil
}

func (d *DNS) query(value string, dnsMsgType uint16) (*dnslib.Msg, error) {
	return d.dnsClient.Query(value, dnsMsgType)
}
