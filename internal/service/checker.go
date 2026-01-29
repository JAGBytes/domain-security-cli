package service

import (
	"fmt"
	"time"

	"github.com/JAGBytes/domain-security-cli/internal/client"
	"github.com/JAGBytes/domain-security-cli/internal/models"
)

type CheckerService struct {
	client *client.SSLabsClient
	cache  map[string]*models.CheckResult
}

func NewCheckerService() *CheckerService {
	return &CheckerService{
		client: client.NewSSLabsClient(),
		cache:  make(map[string]*models.CheckResult),
	}

}


func (s *CheckerService) CheckDomains(doms []string, newAnalysis bool) ([]*models.CheckResult, []error) {

	sem := make(chan struct{}, 5)
	respCh := make(chan *models.CheckResult, len(doms))
	errCh := make(chan error, len(doms))
	results := make([]*models.CheckResult, 0, len(doms))
	errors := make([]error, 0)

	for _, dom := range doms {
		go func(domain string) {
			sem <- struct{}{}
			defer func() { <-sem }()

			res, err := s.CheckDomain(domain, newAnalysis)
			if err != nil {
				errCh <- err
			} else {
				respCh <- res
			}

		}(dom)

	}

	for range doms {
		select {
		case result := <-respCh:
			results = append(results, result)
		case err := <-errCh:
			errors = append(errors, err)
		}
	}

	return results, errors


}

func (s *CheckerService) CheckDomain(domain string, newAnalysis bool) (*models.CheckResult, error) {
	//fmt.Printf("Starting analysis for %s...\n", domain)

	host, err := s.client.Analyze(domain, newAnalysis)

	if err != nil {
		return nil, fmt.Errorf("failed to analyze domain %s: %w", domain, err)
	}

	if host.Status != "READY" && host.Status != "ERROR" {
		//fmt.Println("making a new analysis...")
		//fmt.Println("state different to READY, making polling")

		host, err = s.client.Poll(domain, 10*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("polling failed: %w", err)
		}
	}

	if host.Status == "ERROR" {
		return nil, fmt.Errorf("analysis error: %s", host.StatusMessage)
	}

	//fmt.Println("Analysis complete")

	bestGrade := "N/A"
	for _, ep := range host.Endpoints {
		if ep.Grade != "" && (bestGrade == "N/A" || ep.Grade < bestGrade) {
			bestGrade = ep.Grade
		}
	}

	return &models.CheckResult{
		Domain:    domain,
		Grade:     bestGrade,
		Endpoints: host.Endpoints,
		Status:    host.Status,
	}, nil
}
