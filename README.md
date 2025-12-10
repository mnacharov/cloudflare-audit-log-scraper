# cloudflare-audit-log-scraper

Scrapes DNS Changes in Audit Logs From Cloudflare and sends notifications to Slack.

# Configuration Options

`CLOUDFLARE_API_KEY`
* Description: An API key associated with the CLOUDFLARE_API_EMAIL
* Default: null
* Required: true

`CLOUDFLARE_ORGANIZATION_ID`
* Description: The organization for which you walt to collect audit logs
* Default: null
* Required: true

`SLACK_WEBHOOK`
* Description: Slack webhook for sending the message
* Default: null
* Required: true

`CLOUDFLARE_LOOKUP_INTERVAL`
* Description: How often execute cloudflare access log api in seconds
* Default: 300
* Required: false

`CLOUDFLARE_ZONE_IDS`
* Description: List of zone IDs separated by comma. Monitor all zones by default
* Default: null
* Required: false
