package subdomainimport

import (
	"fmt"
	"strings"

	"github.com/jordyv/reconstore/internal/entities"
	"github.com/miekg/dns"
	"gorm.io/gorm"
)

var (
	dnsServer = "8.8.8.8:53"
)

type DNS struct {
	db        *gorm.DB
	dnsClient *dns.Client
}

func NewDNS(db *gorm.DB) *DNS {
	return &DNS{
		db:        db,
		dnsClient: new(dns.Client),
	}
}

func (d *DNS) AfterSave(s *entities.Subdomain) error {
	dnsInfo := entities.DNSInfo{}

	cnameResult, err := d.query(s.Domain, dns.TypeCNAME)
	if err != nil {
		return err
	}
	if len(cnameResult.Answer) > 0 {
		dnsInfo.CnameRecord = cnameResult.Answer[0].String()
	}

	aResult, err := d.query(s.Domain, dns.TypeA)
	if err != nil {
		return err
	}
	dnsInfo.Status = dns.RcodeToString[aResult.Rcode]
	if len(aResult.Answer) > 0 {
		aRecords := []string{}
		for _, a := range aResult.Answer {
			aRecords = append(aRecords, a.String())
		}
		dnsInfo.ARecords = strings.Join(aRecords, ",")
	}

	nsResult, err := d.query(s.Domain, dns.TypeNS)
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

	s.DNSInfo = dnsInfo
	d.db.Save(s)

	return nil
}

func (d *DNS) query(value string, dnsMsgType uint16) (*dns.Msg, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(value), uint16(dnsMsgType))
	r, _, err := d.dnsClient.Exchange(m, dnsServer)
	if err != nil {
		return nil, fmt.Errorf("could not get %s record - %w", dns.TypeToString[dnsMsgType], err)
	}
	return r, nil
}
