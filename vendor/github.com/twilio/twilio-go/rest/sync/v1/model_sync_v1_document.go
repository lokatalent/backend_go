/*
 * This code was generated by
 * ___ _ _ _ _ _    _ ____    ____ ____ _    ____ ____ _  _ ____ ____ ____ ___ __   __
 *  |  | | | | |    | |  | __ |  | |__| | __ | __ |___ |\ | |___ |__/ |__|  | |  | |__/
 *  |  |_|_| | |___ | |__|    |__| |  | |    |__] |___ | \| |___ |  \ |  |  | |__| |  \
 *
 * Twilio - Sync
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

// SyncV1Document struct for SyncV1Document
type SyncV1Document struct {
	// The unique string that we created to identify the Document resource.
	Sid *string `json:"sid,omitempty"`
	// An application-defined string that uniquely identifies the resource. It can be used in place of the resource's `sid` in the URL to address the resource and can be up to 320 characters long.
	UniqueName *string `json:"unique_name,omitempty"`
	// The SID of the [Account](https://www.twilio.com/docs/iam/api/account) that created the Document resource.
	AccountSid *string `json:"account_sid,omitempty"`
	// The SID of the [Sync Service](https://www.twilio.com/docs/sync/api/service) the resource is associated with.
	ServiceSid *string `json:"service_sid,omitempty"`
	// The absolute URL of the Document resource.
	Url *string `json:"url,omitempty"`
	// The URLs of resources related to the Sync Document.
	Links *map[string]interface{} `json:"links,omitempty"`
	// The current revision of the Sync Document, represented as a string. The `revision` property is used with conditional updates to ensure data consistency.
	Revision *string `json:"revision,omitempty"`
	// An arbitrary, schema-less object that the Sync Document stores. Can be up to 16 KiB in length.
	Data *interface{} `json:"data,omitempty"`
	// The date and time in GMT when the Sync Document expires and will be deleted, specified in [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format. If the Sync Document does not expire, this value is `null`. The Document resource might not be deleted immediately after it expires.
	DateExpires *time.Time `json:"date_expires,omitempty"`
	// The date and time in GMT when the resource was created specified in [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.
	DateCreated *time.Time `json:"date_created,omitempty"`
	// The date and time in GMT when the resource was last updated specified in [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.
	DateUpdated *time.Time `json:"date_updated,omitempty"`
	// The identity of the Sync Document's creator. If the Sync Document is created from the client SDK, the value matches the Access Token's `identity` field. If the Sync Document was created from the REST API, the value is `system`.
	CreatedBy *string `json:"created_by,omitempty"`
}
