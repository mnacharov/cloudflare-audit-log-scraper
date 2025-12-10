package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"strconv"
	"strings"

	"github.com/go-co-op/gocron"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func runCronJobs(apiKey, orgId, slackWebhook string, zones []string, lookupInteval int) {
	s := gocron.NewScheduler(time.UTC)

	s.Every(lookupInteval).Second().Do(func() {
		err := getAuditLogs(apiKey, orgId, slackWebhook, zones)
		if err != nil {
			log.Fatal(err)
		}
	})

	s.StartBlocking()
}

func main() {
	apiKey := os.Getenv("CLOUDFLARE_API_KEY")
	orgId := os.Getenv("CLOUDFLARE_ORGANIZATION_ID")
	interval := os.Getenv("CLOUDFLARE_LOOKUP_INTERVAL")
	slackWebhook := os.Getenv("SLACK_WEBHOOK")
	zones := os.Getenv("CLOUDFLARE_ZONE_IDS")

	if apiKey == "" {
		log.Fatal("Must specify CLOUDFLARE_API_KEY")
	}
	if orgId == "" {
		log.Fatal("Must specify CLOUDFLARE_ORGANIZATION_ID")
	}
	if slackWebhook == "" {
		log.Fatal("Must specify SLACK_WEBHOOK")
	}
	if orgId == "" {
		log.Fatal("Must specify CLOUDFLARE_ORGANIZATION_ID")
	}


	
	// Convert user supplied look back value to an integer otherwise set a default
	var lookupInterval int
	if interval != "" {
		var err error
		lookupInterval, err = strconv.Atoi(interval)
		if err != nil {
			log.Fatal("CLOUDFLARE_LOOKUP_INTERVAL must be an integer", err)
		}
	} else {
		lookupInterval = 300
	}

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":2112", nil)
	}()
	if len(zones) > 0 {
		err := initDNSRecord(apiKey, orgId, strings.Split(zones, ","))
		if err != nil {
			log.Printf("ERROR: failed to list DNS records (deleted records will have only IDs) %s\n", err)
		}
	} else {
		log.Println("WARN: CLOUDFLARE_ZONES not specified, deleted records will have only IDs")
	}

	runCronJobs(apiKey, orgId, slackWebhook, strings.Split(zones, ","), lookupInterval)
}
