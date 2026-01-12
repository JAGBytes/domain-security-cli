package models

type Host struct {
	Host          string     `json:"host"`
	Port          int        `json:"port"`
	Status        string     `json:"status"`
	StatusMessage string     `json:"statusMessage"`
	Endpoints     []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	IPAddress         string `json:"ipAddress"`
	ServerName        string `json:"serverName"`
	StatusMessage     string `json:"statusMessage"`
	Grade             string `json:"grade"`
	GradeTrustIgnored string `json:"gradeTrustIgnored"`
	HasWarnings       bool   `json:"hasWarnings"`
	IsExceptional     bool   `json:"isExceptional"`
	Progress          int    `json:"progress"`
	Duration          int    `json:"duration"`
	Delegation        int    `json:"delegation"`
}

type CheckResult struct {
    Domain    string     `json:"domain"`
    Grade     string     `json:"grade"`
    Endpoints []Endpoint `json:"endpoints"`
    Status    string     `json:"status"`
}
