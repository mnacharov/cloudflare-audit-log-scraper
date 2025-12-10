// https://developers.cloudflare.com/api/resources/accounts/subresources/logs/#(resource)%20accounts.logs.audit
package main

import (
	"encoding/json"
	"fmt"
)


type AccountAuditLogsError struct {
	Message string `json:"message"`
}

func (e *AccountAuditLogsError) Error() string {
	return e.Message
}

type AccountAuditLogsResultInfo struct {
	Count  int    `json:"count"`
	Cursor string `json:"cursor"`
}

type AccountAuditLogsAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`	
}
type AccountAuditLogsAction struct {
	Description string `json:"description"`
	Result      string `json:"result"`
	Time        string `json:"time"`
	Type        string `json:"type"`	
}
type AccountAuditLogsActor struct {
	ID        string `json:"id"`
	Context   string `json:"context"`
	Email     string `json:"email"`
	IPAddress string `json:"ip_address"`
	TokenID   string `json:"token_id"`
	TokenName string `json:"token_name"`
	Type      string `json:"type"`
}
type AccountAuditLogsRaw struct {
	CfRayID    string `json:"cf_ray_id"`
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
	URI        string `json:"uri"`
	UserAgent  string `json:"user_agent"`
}
type AccountAuditLogsResourceDNSRequest struct {
	Name    string   `json:"name"`
	Tags    []string `json:"tags"`
	Type    string   `json:"type"`
	Comment string   `json:"comment"`
	Content string   `json:"content"`
	Proxied bool     `json:"proxied"`
}
type AccountAuditLogsResourceDNSResponse struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Tags    []string `json:"tags"`
	Type    string   `json:"type"`
	Comment string   `json:"comment"`
	Content string   `json:"content"`
	Proxied bool     `json:"proxied"`
}	
type AccountAuditLogsResource struct {
	ID       string          `json:"id"`
	Product  string          `json:"product"`
	Request  json.RawMessage `json:"request"`
	Response json.RawMessage `json:"response"`
        Scope    string          `json:"scope"`
}
type AccountAuditLogsZone struct {
	ID   string `json:"id"`
	Name string `json:"name"`	
}

type AccountAuditLogsResult struct {
	ID       string                   `json:"id"`
	Account  AccountAuditLogsAccount  `json:"account"`
	Action   AccountAuditLogsAction   `json:"action"`
	Actor    AccountAuditLogsActor    `json:"actor"`
	Raw      AccountAuditLogsRaw      `json:"raw"`
	Resource AccountAuditLogsResource `json:"resource"`
	Zone     AccountAuditLogsZone     `json:"zone"`
}

type AccountAuditLogs struct {
	Errors     []AccountAuditLogsError    `json:"errors"`
	Result     []AccountAuditLogsResult   `json:"result"`
	ResultInfo AccountAuditLogsResultInfo `json:"result_info"`
	Success    bool                       `json:"success"`
}

type ListDNSRecordsError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ListDNSRecordsError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// to populate dnsRecords map and be able to output deleted record info
type ListDNSRecords struct {
	Errors     []ListDNSRecordsError                 `json:"errors"`
	Result     []AccountAuditLogsResourceDNSResponse `json:"result"`
	Success    bool                                  `json:"success"`
}
