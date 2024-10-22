/*
 * This code was generated by
 * ___ _ _ _ _ _    _ ____    ____ ____ _    ____ ____ _  _ ____ ____ ____ ___ __   __
 *  |  | | | | |    | |  | __ |  | |__| | __ | __ |___ |\ | |___ |__/ |__|  | |  | |__/
 *  |  |_|_| | |___ | |__|    |__| |  | |    |__] |___ | \| |___ |  \ |  |  | |__| |  \
 *
 * Twilio - Assistants
 * This is the public Twilio REST API.
 *
 * NOTE: This class is auto generated by OpenAPI Generator.
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

package openapi

import (
	"time"
)

// AssistantsV1KnowledgeStatus struct for AssistantsV1KnowledgeStatus
type AssistantsV1KnowledgeStatus struct {
	// The SID of the [Account](https://www.twilio.com/docs/iam/api/account) that created the Knowledge resource.
	AccountSid string `json:"account_sid,omitempty"`
	// The status of processing the knowledge source ('QUEUED', 'PROCESSING', 'COMPLETED', 'FAILED')
	Status string `json:"status"`
	// The last status of processing the knowledge source ('QUEUED', 'PROCESSING', 'COMPLETED', 'FAILED')
	LastStatus string `json:"last_status,omitempty"`
	// The date and time in GMT when the Knowledge was last updated specified in [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.
	DateUpdated time.Time `json:"date_updated,omitempty"`
}
