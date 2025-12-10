package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	logsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cloudflare_audit_logs_processed_total",
		Help: "The total number of processed events",
	})
	logsProcessedErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cloudflare_audit_logs_processed_errors_total",
		Help: "The total number of errors during processing events",
	})
	dnsRecords = make(map[string]AccountAuditLogsResourceDNSResponse)
	lastProcessed = time.Now().Add(time.Duration(-5) * time.Minute).UTC().Format(time.RFC3339)
)

func initDNSRecord(apiKey, orgId string, zones []string) error {
	var failed bool
	client := &http.Client{}
	for _, zone_id := range zones {
		req, err := http.NewRequest("GET", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?per_page=5000000", zone_id), nil)
		if err != nil {
			log.Fatalf("Error creating request to Cloudflare API: %v", err)
		}
		req.Header = http.Header{
			"Authorization": {fmt.Sprintf("Bearer %s", apiKey)},
		}
		res, err := client.Do(req)
		if err != nil {
			log.Printf("Error executing request to Cloudflare API: %v\n", err)
			failed = true
			continue
		}
		var listDNSRecords ListDNSRecords
		err = json.NewDecoder(res.Body).Decode(&listDNSRecords)
		if err != nil {
			log.Printf("Error decoding response from Cloudflare API: %v\n", err)
			failed = true
			continue
		}
		if !listDNSRecords.Success {
			log.Printf("Error accessing Cloudflare API: %+v\n", listDNSRecords.Errors)
			failed = true
			continue
		}
		for _, record := range listDNSRecords.Result {
			dnsRecords[record.ID] = record
		}
	}
	if failed {
		return errors.New("Some of the zones aren't populated")
	}
	return nil
}
// Get audit logs and process them until no more records are returned
func getAuditLogs(apiKey, orgId, slackWebhook string, zones []string) error {
	client := &http.Client{}
	before := time.Now().UTC().Format(time.RFC3339)
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/logs/audit?limit=1000&direction=asc&action_type.not=view&resource_product=dns_records&since=%s&before=%s",
		orgId, lastProcessed, before)
	for _, z := range zones {
		url = fmt.Sprintf("%s&zone_id=%s", url, z)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request to Cloudflare API: %v", err)
	}
	req.Header = http.Header{
		"Authorization": {fmt.Sprintf("Bearer %s", apiKey)},
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	var auditLogs AccountAuditLogs
	err = json.NewDecoder(res.Body).Decode(&auditLogs)
	if err != nil {
		return err
	}
	if !auditLogs.Success {
		return fmt.Errorf("ERROR: %+v\n", auditLogs)
	}
	for _, err := range auditLogs.Errors {
		log.Printf("ERROR: %s\n", err.Error())
	}
	for _, auditLog := range auditLogs.Result {
		err = processAuditLog(auditLog, slackWebhook)
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			logsProcessedErrors.Inc()
			continue
		} else {
			logsProcessed.Inc()
		}
		lastProcessed = auditLog.Action.Time
	}
	// one milliseconds *after* the last shown record
	t, err := time.Parse(time.RFC3339, lastProcessed)
	if err != nil {
		return err
	}
	lastProcessed = t.Add(1*time.Second).Format(time.RFC3339)
	return nil
}

func processAuditLog(auditLog AccountAuditLogsResult, slackWebhook string) error {
	// fmt.Printf("Result ID: %s\n", log.ID)
	if auditLog.Resource.Product != "dns_records" {
		return fmt.Errorf("Not supported audit record product %s", auditLog.Resource.Product)
	}
	var response AccountAuditLogsResourceDNSResponse
	json.Unmarshal(auditLog.Resource.Response, &response)
	if response.ID == "" {
		// cloudflare can return empty response for recently created records
		return fmt.Errorf("Invalid response: %v", auditLog.Resource.Response)
	}
	if auditLog.Action.Result == "success" {
		if auditLog.Action.Type == "create" {
			dnsRecords[response.ID] = response
		}
		if auditLog.Action.Type == "delete" {
			value, ok := dnsRecords[response.ID]
			if ok {
				response = value  // to output what was in the deleted record
			}
		}
	}
	if auditLog.Action.Type == "create" || auditLog.Action.Type == "delete" || auditLog.Action.Type == "update" {
		return slackMessage(response, auditLog, slackWebhook)
	}
	return nil
}

func slackMessage(response AccountAuditLogsResourceDNSResponse, auditLog AccountAuditLogsResult, slackWebhook string) error {
	// Build readable Slack message
	text := fmt.Sprintf(
                "*Cloudflare %s (%s)*\n"+
			"*Time:*        %s\n"+
			"*Actor:*       %s (%s)\n"+
			"*Zone:*        %s\n"+
			"*DNS:*         %s\t%s\t%s\t(%s)\n"+
			"*IsProxied:*  %t\n",
                auditLog.Action.Description, auditLog.Action.Result,
                auditLog.Action.Time,
                auditLog.Actor.Email, auditLog.Actor.Type,
		auditLog.Zone.Name,
                response.Type, response.Name, response.Content, response.ID,
                response.Proxied,
	)
	msg := map[string]string{"text": text}
	payload, _ := json.Marshal(msg)
	_, err := http.Post(slackWebhook, "application/json", bytes.NewBuffer(payload))
	return err
}
