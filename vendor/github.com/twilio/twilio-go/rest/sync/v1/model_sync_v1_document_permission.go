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

// SyncV1DocumentPermission struct for SyncV1DocumentPermission
type SyncV1DocumentPermission struct {
	// The SID of the [Account](https://www.twilio.com/docs/iam/api/account) that created the Document Permission resource.
	AccountSid *string `json:"account_sid,omitempty"`
	// The SID of the [Sync Service](https://www.twilio.com/docs/sync/api/service) the resource is associated with.
	ServiceSid *string `json:"service_sid,omitempty"`
	// The SID of the Sync Document to which the Document Permission applies.
	DocumentSid *string `json:"document_sid,omitempty"`
	// The application-defined string that uniquely identifies the resource's User within the Service to an FPA token.
	Identity *string `json:"identity,omitempty"`
	// Whether the identity can read the Sync Document.
	Read *bool `json:"read,omitempty"`
	// Whether the identity can update the Sync Document.
	Write *bool `json:"write,omitempty"`
	// Whether the identity can delete the Sync Document.
	Manage *bool `json:"manage,omitempty"`
	// The absolute URL of the Sync Document Permission resource.
	Url *string `json:"url,omitempty"`
}
