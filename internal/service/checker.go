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

func (s *CheckerService) CacheDomain(domain string, newAnalysis bool) (*models.CheckResult, error) {

	if s.cache[domain] != nil {
		fmt.Printf("chace interno del cli\n")
		return s.cache[domain], nil
	}

	fmt.Printf("sin cache")
	host, err := s.CheckDomain(domain, newAnalysis)
	s.cache[domain] = host
	return host, err

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
