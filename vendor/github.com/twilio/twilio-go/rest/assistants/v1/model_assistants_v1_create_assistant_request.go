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

// AssistantsV1CreateAssistantRequest struct for AssistantsV1CreateAssistantRequest
type AssistantsV1CreateAssistantRequest struct {
	CustomerAi AssistantsV1CustomerAi `json:"customer_ai,omitempty"`
	// The name of the assistant.
	Name string `json:"name"`
	// The owner/company of the assistant.
	Owner string `json:"owner,omitempty"`
	// The personality prompt to be used for assistant.
	PersonalityPrompt string                        `json:"personality_prompt,omitempty"`
	SegmentCredential AssistantsV1SegmentCredential `json:"segment_credential,omitempty"`
}