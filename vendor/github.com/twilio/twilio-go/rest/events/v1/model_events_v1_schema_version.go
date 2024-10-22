/*
 * This code was generated by
 * ___ _ _ _ _ _    _ ____    ____ ____ _    ____ ____ _  _ ____ ____ ____ ___ __   __
 *  |  | | | | |    | |  | __ |  | |__| | __ | __ |___ |\ | |___ |__/ |__|  | |  | |__/
 *  |  |_|_| | |___ | |__|    |__| |  | |    |__] |___ | \| |___ |  \ |  |  | |__| |  \
 *
 * Twilio - Events
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

// EventsV1SchemaVersion struct for EventsV1SchemaVersion
type EventsV1SchemaVersion struct {
	// The unique identifier of the schema. Each schema can have multiple versions, that share the same id.
	Id *string `json:"id,omitempty"`
	// The version of this schema.
	SchemaVersion int `json:"schema_version,omitempty"`
	// The date the schema version was created, given in ISO 8601 format.
	DateCreated *time.Time `json:"date_created,omitempty"`
	// The URL of this resource.
	Url *string `json:"url,omitempty"`
	Raw *string `json:"raw,omitempty"`
}
