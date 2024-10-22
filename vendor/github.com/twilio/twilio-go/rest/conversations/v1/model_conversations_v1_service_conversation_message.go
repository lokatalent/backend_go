/*
 * This code was generated by
 * ___ _ _ _ _ _    _ ____    ____ ____ _    ____ ____ _  _ ____ ____ ____ ___ __   __
 *  |  | | | | |    | |  | __ |  | |__| | __ | __ |___ |\ | |___ |__/ |__|  | |  | |__/
 *  |  |_|_| | |___ | |__|    |__| |  | |    |__] |___ | \| |___ |  \ |  |  | |__| |  \
 *
 * Twilio - Conversations
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

// ConversationsV1ServiceConversationMessage struct for ConversationsV1ServiceConversationMessage
type ConversationsV1ServiceConversationMessage struct {
	// The unique ID of the [Account](https://www.twilio.com/docs/iam/api/account) responsible for this message.
	AccountSid *string `json:"account_sid,omitempty"`
	// The SID of the [Conversation Service](https://www.twilio.com/docs/conversations/api/service-resource) the Participant resource is associated with.
	ChatServiceSid *string `json:"chat_service_sid,omitempty"`
	// The unique ID of the [Conversation](https://www.twilio.com/docs/conversations/api/conversation-resource) for this message.
	ConversationSid *string `json:"conversation_sid,omitempty"`
	// A 34 character string that uniquely identifies this resource.
	Sid *string `json:"sid,omitempty"`
	// The index of the message within the [Conversation](https://www.twilio.com/docs/conversations/api/conversation-resource).
	Index int `json:"index,omitempty"`
	// The channel specific identifier of the message's author. Defaults to `system`.
	Author *string `json:"author,omitempty"`
	// The content of the message, can be up to 1,600 characters long.
	Body *string `json:"body,omitempty"`
	// An array of objects that describe the Message's media, if the message contains media. Each object contains these fields: `content_type` with the MIME type of the media, `filename` with the name of the media, `sid` with the SID of the Media resource, and `size` with the media object's file size in bytes. If the Message has no media, this value is `null`.
	Media *[]interface{} `json:"media,omitempty"`
	// A string metadata field you can use to store any data you wish. The string value must contain structurally valid JSON if specified.  **Note** that if the attributes are not set \"{}\" will be returned.
	Attributes *string `json:"attributes,omitempty"`
	// The unique ID of messages's author participant. Null in case of `system` sent message.
	ParticipantSid *string `json:"participant_sid,omitempty"`
	// The date that this resource was created.
	DateCreated *time.Time `json:"date_created,omitempty"`
	// The date that this resource was last updated. `null` if the message has not been edited.
	DateUpdated *time.Time `json:"date_updated,omitempty"`
	// An object that contains the summary of delivery statuses for the message to non-chat participants.
	Delivery *interface{} `json:"delivery,omitempty"`
	// An absolute API resource URL for this message.
	Url *string `json:"url,omitempty"`
	// Contains an absolute API resource URL to access the delivery & read receipts of this message.
	Links *map[string]interface{} `json:"links,omitempty"`
	// The unique ID of the multi-channel [Rich Content](https://www.twilio.com/docs/content) template.
	ContentSid *string `json:"content_sid,omitempty"`
}
