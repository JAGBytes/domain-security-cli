package service

import (
	"fmt"
	"time"

	"github.com/JAGBytes/domain-security-cli/internal/client"
	"github.com/JAGBytes/domain-security-cli/internal/models"
)

type CheckerService struct {
	client *client.SSLabsClient
}

func NewCheckerService() *CheckerService {
	return &CheckerService{
		client: client.NewSSLabsClient(),
	}

}

func (s *CheckerService) CheckDomains(doms []string, newAnalysis bool) ([]*models.CheckResult, error) {

	respCh := make(chan *models.CheckResult, len(doms))
	errCh := make(chan error, len(doms))
	results := make([]*models.CheckResult, 0, len(doms))

	for _, dom := range doms {
		go s.CheckDomain(dom, newAnalysis, respCh, errCh)
	}

	for range doms {
		select {
		case result := <-respCh:
			results = append(results, result)
		case err := <-errCh:
			return results, err
		}
	}

	return results, nil

}

func (s *CheckerService) CheckDomain(domain string, newAnalysis bool, respCh chan<- *models.CheckResult, errCh chan<- error) {
	//fmt.Printf("Starting analysis for %s...\n", domain)

	host, err := s.client.Analyze(domain, newAnalysis)

	if err != nil {
		errCh <- fmt.Errorf("failed to analyze domain %s: %w", domain, err)
	}

	if host.Status != "READY" && host.Status != "ERROR" {
		//fmt.Println("making a new analysis...")
		//fmt.Println("state different to READY, making polling")

		host, err = s.client.Poll(domain, 10*time.Minute)
		if err != nil {
			errCh <- fmt.Errorf("polling failed: %w", err)
		}
	}

	if host.Status == "ERROR" {
		errCh <- fmt.Errorf("analysis error: %s", host.StatusMessage)
	}

	//fmt.Println("Analysis complete")

	bestGrade := "N/A"
	for _, ep := range host.Endpoints {
		if ep.Grade != "" && (bestGrade == "N/A" || ep.Grade < bestGrade) {
			bestGrade = ep.Grade
		}
	}

	respCh <- &models.CheckResult{
		Domain:    domain,
		Grade:     bestGrade,
		Endpoints: host.Endpoints,
		Status:    host.Status,
	}
}
