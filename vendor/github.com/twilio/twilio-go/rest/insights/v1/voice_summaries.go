/*
 * This code was generated by
 * ___ _ _ _ _ _    _ ____    ____ ____ _    ____ ____ _  _ ____ ____ ____ ___ __   __
 *  |  | | | | |    | |  | __ |  | |__| | __ | __ |___ |\ | |___ |__/ |__|  | |  | |__/
 *  |  |_|_| | |___ | |__|    |__| |  | |    |__] |___ | \| |___ |  \ |  |  | |__| |  \
 *
 * Twilio - Insights
 * This is the public Twilio REST API.
 *
 * NOTE: This class is auto generated by OpenAPI Generator.
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

package openapi

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/twilio/twilio-go/client"
)

// Optional parameters for the method 'ListCallSummaries'
type ListCallSummariesParams struct {
	// A calling party. Could be an E.164 number, a SIP URI, or a Twilio Client registered name.
	From *string `json:"From,omitempty"`
	// A called party. Could be an E.164 number, a SIP URI, or a Twilio Client registered name.
	To *string `json:"To,omitempty"`
	// An origination carrier.
	FromCarrier *string `json:"FromCarrier,omitempty"`
	// A destination carrier.
	ToCarrier *string `json:"ToCarrier,omitempty"`
	// A source country code based on phone number in From.
	FromCountryCode *string `json:"FromCountryCode,omitempty"`
	// A destination country code. Based on phone number in To.
	ToCountryCode *string `json:"ToCountryCode,omitempty"`
	// A boolean flag indicating whether or not the caller was verified using SHAKEN/STIR.One of 'true' or 'false'.
	VerifiedCaller *bool `json:"VerifiedCaller,omitempty"`
	// A boolean flag indicating the presence of one or more [Voice Insights Call Tags](https://www.twilio.com/docs/voice/voice-insights/api/call/details-call-tags).
	HasTag *bool `json:"HasTag,omitempty"`
	// A Start time of the calls. xm (x minutes), xh (x hours), xd (x days), 1w, 30m, 3d, 4w or datetime-ISO. Defaults to 4h.
	StartTime *string `json:"StartTime,omitempty"`
	// An End Time of the calls. xm (x minutes), xh (x hours), xd (x days), 1w, 30m, 3d, 4w or datetime-ISO. Defaults to 0m.
	EndTime *string `json:"EndTime,omitempty"`
	// A Call Type of the calls. One of `carrier`, `sip`, `trunking` or `client`.
	CallType *string `json:"CallType,omitempty"`
	// A Call State of the calls. One of `ringing`, `completed`, `busy`, `fail`, `noanswer`, `canceled`, `answered`, `undialed`.
	CallState *string `json:"CallState,omitempty"`
	// A Direction of the calls. One of `outbound_api`, `outbound_dial`, `inbound`, `trunking_originating`, `trunking_terminating`.
	Direction *string `json:"Direction,omitempty"`
	// A Processing State of the Call Summaries. One of `completed`, `partial` or `all`.
	ProcessingState *string `json:"ProcessingState,omitempty"`
	// A Sort By criterion for the returned list of Call Summaries. One of `start_time` or `end_time`.
	SortBy *string `json:"SortBy,omitempty"`
	// A unique SID identifier of a Subaccount.
	Subaccount *string `json:"Subaccount,omitempty"`
	// A boolean flag indicating an abnormal session where the last SIP response was not 200 OK.
	AbnormalSession *bool `json:"AbnormalSession,omitempty"`
	// An Answered By value for the calls based on `Answering Machine Detection (AMD)`. One of `unknown`, `machine_start`, `machine_end_beep`, `machine_end_silence`, `machine_end_other`, `human` or `fax`.
	AnsweredBy *string `json:"AnsweredBy,omitempty"`
	// Either machine or human.
	AnsweredByAnnotation *string `json:"AnsweredByAnnotation,omitempty"`
	// A Connectivity Issue with the calls. One of `no_connectivity_issue`, `invalid_number`, `caller_id`, `dropped_call`, or `number_reachability`.
	ConnectivityIssueAnnotation *string `json:"ConnectivityIssueAnnotation,omitempty"`
	// A subjective Quality Issue with the calls. One of `no_quality_issue`, `low_volume`, `choppy_robotic`, `echo`, `dtmf`, `latency`, `owa`, `static_noise`.
	QualityIssueAnnotation *string `json:"QualityIssueAnnotation,omitempty"`
	// A boolean flag indicating spam calls.
	SpamAnnotation *bool `json:"SpamAnnotation,omitempty"`
	// A Call Score of the calls. Use a range of 1-5 to indicate the call experience score, with the following mapping as a reference for the rated call [5: Excellent, 4: Good, 3 : Fair, 2 : Poor, 1: Bad].
	CallScoreAnnotation *string `json:"CallScoreAnnotation,omitempty"`
	// A boolean flag indicating whether or not the calls were branded using Twilio Branded Calls. One of 'true' or 'false'
	BrandedEnabled *bool `json:"BrandedEnabled,omitempty"`
	// A boolean flag indicating whether or not the phone number had voice integrity enabled.One of 'true' or 'false'
	VoiceIntegrityEnabled *bool `json:"VoiceIntegrityEnabled,omitempty"`
	// A unique SID identifier of the Branded Call.
	BrandedBundleSid *string `json:"BrandedBundleSid,omitempty"`
	// A unique SID identifier of the Voice Integrity Profile.
	VoiceIntegrityBundleSid *string `json:"VoiceIntegrityBundleSid,omitempty"`
	// A Voice Integrity Use Case . Is of type enum. One of 'abandoned_cart', 'appointment_reminders', 'appointment_scheduling', 'asset_management', 'automated_support', 'call_tracking', 'click_to_call', 'contact_tracing', 'contactless_delivery', 'customer_support', 'dating/social', 'delivery_notifications', 'distance_learning', 'emergency_notifications', 'employee_notifications', 'exam_proctoring', 'field_notifications', 'first_responder', 'fraud_alerts', 'group_messaging', 'identify_&_verification', 'intelligent_routing', 'lead_alerts', 'lead_distribution', 'lead_generation', 'lead_management', 'lead_nurturing', 'marketing_events', 'mass_alerts', 'meetings/collaboration', 'order_notifications', 'outbound_dialer', 'pharmacy', 'phone_system', 'purchase_confirmation', 'remote_appointments', 'rewards_program', 'self-service', 'service_alerts', 'shift_management', 'survey/research', 'telehealth', 'telemarketing', 'therapy_(individual+group)'.
	VoiceIntegrityUseCase *string `json:"VoiceIntegrityUseCase,omitempty"`
	// A Business Identity of the calls. Is of type enum. One of 'direct_customer', 'isv_reseller_or_partner'.
	BusinessProfileIdentity *string `json:"BusinessProfileIdentity,omitempty"`
	// A Business Industry of the calls. Is of type enum. One of 'automotive', 'agriculture', 'banking', 'consumer', 'construction', 'education', 'engineering', 'energy', 'oil_and_gas', 'fast_moving_consumer_goods', 'financial', 'fintech', 'food_and_beverage', 'government', 'healthcare', 'hospitality', 'insurance', 'legal', 'manufacturing', 'media', 'online', 'professional_services', 'raw_materials', 'real_estate', 'religion', 'retail', 'jewelry', 'technology', 'telecommunications', 'transportation', 'travel', 'electronics', 'not_for_profit'
	BusinessProfileIndustry *string `json:"BusinessProfileIndustry,omitempty"`
	// A unique SID identifier of the Business Profile.
	BusinessProfileBundleSid *string `json:"BusinessProfileBundleSid,omitempty"`
	// A Business Profile Type of the calls. Is of type enum. One of 'primary', 'secondary'.
	BusinessProfileType *string `json:"BusinessProfileType,omitempty"`
	// How many resources to return in each list page. The default is 50, and the maximum is 1000.
	PageSize *int `json:"PageSize,omitempty"`
	// Max number of records to return.
	Limit *int `json:"limit,omitempty"`
}

func (params *ListCallSummariesParams) SetFrom(From string) *ListCallSummariesParams {
	params.From = &From
	return params
}
func (params *ListCallSummariesParams) SetTo(To string) *ListCallSummariesParams {
	params.To = &To
	return params
}
func (params *ListCallSummariesParams) SetFromCarrier(FromCarrier string) *ListCallSummariesParams {
	params.FromCarrier = &FromCarrier
	return params
}
func (params *ListCallSummariesParams) SetToCarrier(ToCarrier string) *ListCallSummariesParams {
	params.ToCarrier = &ToCarrier
	return params
}
func (params *ListCallSummariesParams) SetFromCountryCode(FromCountryCode string) *ListCallSummariesParams {
	params.FromCountryCode = &FromCountryCode
	return params
}
func (params *ListCallSummariesParams) SetToCountryCode(ToCountryCode string) *ListCallSummariesParams {
	params.ToCountryCode = &ToCountryCode
	return params
}
func (params *ListCallSummariesParams) SetVerifiedCaller(VerifiedCaller bool) *ListCallSummariesParams {
	params.VerifiedCaller = &VerifiedCaller
	return params
}
func (params *ListCallSummariesParams) SetHasTag(HasTag bool) *ListCallSummariesParams {
	params.HasTag = &HasTag
	return params
}
func (params *ListCallSummariesParams) SetStartTime(StartTime string) *ListCallSummariesParams {
	params.StartTime = &StartTime
	return params
}
func (params *ListCallSummariesParams) SetEndTime(EndTime string) *ListCallSummariesParams {
	params.EndTime = &EndTime
	return params
}
func (params *ListCallSummariesParams) SetCallType(CallType string) *ListCallSummariesParams {
	params.CallType = &CallType
	return params
}
func (params *ListCallSummariesParams) SetCallState(CallState string) *ListCallSummariesParams {
	params.CallState = &CallState
	return params
}
func (params *ListCallSummariesParams) SetDirection(Direction string) *ListCallSummariesParams {
	params.Direction = &Direction
	return params
}
func (params *ListCallSummariesParams) SetProcessingState(ProcessingState string) *ListCallSummariesParams {
	params.ProcessingState = &ProcessingState
	return params
}
func (params *ListCallSummariesParams) SetSortBy(SortBy string) *ListCallSummariesParams {
	params.SortBy = &SortBy
	return params
}
func (params *ListCallSummariesParams) SetSubaccount(Subaccount string) *ListCallSummariesParams {
	params.Subaccount = &Subaccount
	return params
}
func (params *ListCallSummariesParams) SetAbnormalSession(AbnormalSession bool) *ListCallSummariesParams {
	params.AbnormalSession = &AbnormalSession
	return params
}
func (params *ListCallSummariesParams) SetAnsweredBy(AnsweredBy string) *ListCallSummariesParams {
	params.AnsweredBy = &AnsweredBy
	return params
}
func (params *ListCallSummariesParams) SetAnsweredByAnnotation(AnsweredByAnnotation string) *ListCallSummariesParams {
	params.AnsweredByAnnotation = &AnsweredByAnnotation
	return params
}
func (params *ListCallSummariesParams) SetConnectivityIssueAnnotation(ConnectivityIssueAnnotation string) *ListCallSummariesParams {
	params.ConnectivityIssueAnnotation = &ConnectivityIssueAnnotation
	return params
}
func (params *ListCallSummariesParams) SetQualityIssueAnnotation(QualityIssueAnnotation string) *ListCallSummariesParams {
	params.QualityIssueAnnotation = &QualityIssueAnnotation
	return params
}
func (params *ListCallSummariesParams) SetSpamAnnotation(SpamAnnotation bool) *ListCallSummariesParams {
	params.SpamAnnotation = &SpamAnnotation
	return params
}
func (params *ListCallSummariesParams) SetCallScoreAnnotation(CallScoreAnnotation string) *ListCallSummariesParams {
	params.CallScoreAnnotation = &CallScoreAnnotation
	return params
}
func (params *ListCallSummariesParams) SetBrandedEnabled(BrandedEnabled bool) *ListCallSummariesParams {
	params.BrandedEnabled = &BrandedEnabled
	return params
}
func (params *ListCallSummariesParams) SetVoiceIntegrityEnabled(VoiceIntegrityEnabled bool) *ListCallSummariesParams {
	params.VoiceIntegrityEnabled = &VoiceIntegrityEnabled
	return params
}
func (params *ListCallSummariesParams) SetBrandedBundleSid(BrandedBundleSid string) *ListCallSummariesParams {
	params.BrandedBundleSid = &BrandedBundleSid
	return params
}
func (params *ListCallSummariesParams) SetVoiceIntegrityBundleSid(VoiceIntegrityBundleSid string) *ListCallSummariesParams {
	params.VoiceIntegrityBundleSid = &VoiceIntegrityBundleSid
	return params
}
func (params *ListCallSummariesParams) SetVoiceIntegrityUseCase(VoiceIntegrityUseCase string) *ListCallSummariesParams {
	params.VoiceIntegrityUseCase = &VoiceIntegrityUseCase
	return params
}
func (params *ListCallSummariesParams) SetBusinessProfileIdentity(BusinessProfileIdentity string) *ListCallSummariesParams {
	params.BusinessProfileIdentity = &BusinessProfileIdentity
	return params
}
func (params *ListCallSummariesParams) SetBusinessProfileIndustry(BusinessProfileIndustry string) *ListCallSummariesParams {
	params.BusinessProfileIndustry = &BusinessProfileIndustry
	return params
}
func (params *ListCallSummariesParams) SetBusinessProfileBundleSid(BusinessProfileBundleSid string) *ListCallSummariesParams {
	params.BusinessProfileBundleSid = &BusinessProfileBundleSid
	return params
}
func (params *ListCallSummariesParams) SetBusinessProfileType(BusinessProfileType string) *ListCallSummariesParams {
	params.BusinessProfileType = &BusinessProfileType
	return params
}
func (params *ListCallSummariesParams) SetPageSize(PageSize int) *ListCallSummariesParams {
	params.PageSize = &PageSize
	return params
}
func (params *ListCallSummariesParams) SetLimit(Limit int) *ListCallSummariesParams {
	params.Limit = &Limit
	return params
}

// Retrieve a single page of CallSummaries records from the API. Request is executed immediately.
func (c *ApiService) PageCallSummaries(params *ListCallSummariesParams, pageToken, pageNumber string) (*ListCallSummariesResponse, error) {
	path := "/v1/Voice/Summaries"

	data := url.Values{}
	headers := map[string]interface{}{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	if params != nil && params.From != nil {
		data.Set("From", *params.From)
	}
	if params != nil && params.To != nil {
		data.Set("To", *params.To)
	}
	if params != nil && params.FromCarrier != nil {
		data.Set("FromCarrier", *params.FromCarrier)
	}
	if params != nil && params.ToCarrier != nil {
		data.Set("ToCarrier", *params.ToCarrier)
	}
	if params != nil && params.FromCountryCode != nil {
		data.Set("FromCountryCode", *params.FromCountryCode)
	}
	if params != nil && params.ToCountryCode != nil {
		data.Set("ToCountryCode", *params.ToCountryCode)
	}
	if params != nil && params.VerifiedCaller != nil {
		data.Set("VerifiedCaller", fmt.Sprint(*params.VerifiedCaller))
	}
	if params != nil && params.HasTag != nil {
		data.Set("HasTag", fmt.Sprint(*params.HasTag))
	}
	if params != nil && params.StartTime != nil {
		data.Set("StartTime", *params.StartTime)
	}
	if params != nil && params.EndTime != nil {
		data.Set("EndTime", *params.EndTime)
	}
	if params != nil && params.CallType != nil {
		data.Set("CallType", *params.CallType)
	}
	if params != nil && params.CallState != nil {
		data.Set("CallState", *params.CallState)
	}
	if params != nil && params.Direction != nil {
		data.Set("Direction", *params.Direction)
	}
	if params != nil && params.ProcessingState != nil {
		data.Set("ProcessingState", *params.ProcessingState)
	}
	if params != nil && params.SortBy != nil {
		data.Set("SortBy", *params.SortBy)
	}
	if params != nil && params.Subaccount != nil {
		data.Set("Subaccount", *params.Subaccount)
	}
	if params != nil && params.AbnormalSession != nil {
		data.Set("AbnormalSession", fmt.Sprint(*params.AbnormalSession))
	}
	if params != nil && params.AnsweredBy != nil {
		data.Set("AnsweredBy", *params.AnsweredBy)
	}
	if params != nil && params.AnsweredByAnnotation != nil {
		data.Set("AnsweredByAnnotation", *params.AnsweredByAnnotation)
	}
	if params != nil && params.ConnectivityIssueAnnotation != nil {
		data.Set("ConnectivityIssueAnnotation", *params.ConnectivityIssueAnnotation)
	}
	if params != nil && params.QualityIssueAnnotation != nil {
		data.Set("QualityIssueAnnotation", *params.QualityIssueAnnotation)
	}
	if params != nil && params.SpamAnnotation != nil {
		data.Set("SpamAnnotation", fmt.Sprint(*params.SpamAnnotation))
	}
	if params != nil && params.CallScoreAnnotation != nil {
		data.Set("CallScoreAnnotation", *params.CallScoreAnnotation)
	}
	if params != nil && params.BrandedEnabled != nil {
		data.Set("BrandedEnabled", fmt.Sprint(*params.BrandedEnabled))
	}
	if params != nil && params.VoiceIntegrityEnabled != nil {
		data.Set("VoiceIntegrityEnabled", fmt.Sprint(*params.VoiceIntegrityEnabled))
	}
	if params != nil && params.BrandedBundleSid != nil {
		data.Set("BrandedBundleSid", *params.BrandedBundleSid)
	}
	if params != nil && params.VoiceIntegrityBundleSid != nil {
		data.Set("VoiceIntegrityBundleSid", *params.VoiceIntegrityBundleSid)
	}
	if params != nil && params.VoiceIntegrityUseCase != nil {
		data.Set("VoiceIntegrityUseCase", *params.VoiceIntegrityUseCase)
	}
	if params != nil && params.BusinessProfileIdentity != nil {
		data.Set("BusinessProfileIdentity", *params.BusinessProfileIdentity)
	}
	if params != nil && params.BusinessProfileIndustry != nil {
		data.Set("BusinessProfileIndustry", *params.BusinessProfileIndustry)
	}
	if params != nil && params.BusinessProfileBundleSid != nil {
		data.Set("BusinessProfileBundleSid", *params.BusinessProfileBundleSid)
	}
	if params != nil && params.BusinessProfileType != nil {
		data.Set("BusinessProfileType", *params.BusinessProfileType)
	}
	if params != nil && params.PageSize != nil {
		data.Set("PageSize", fmt.Sprint(*params.PageSize))
	}

	if pageToken != "" {
		data.Set("PageToken", pageToken)
	}
	if pageNumber != "" {
		data.Set("Page", pageNumber)
	}

	resp, err := c.requestHandler.Get(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &ListCallSummariesResponse{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}

// Lists CallSummaries records from the API as a list. Unlike stream, this operation is eager and loads 'limit' records into memory before returning.
func (c *ApiService) ListCallSummaries(params *ListCallSummariesParams) ([]InsightsV1CallSummaries, error) {
	response, errors := c.StreamCallSummaries(params)

	records := make([]InsightsV1CallSummaries, 0)
	for record := range response {
		records = append(records, record)
	}

	if err := <-errors; err != nil {
		return nil, err
	}

	return records, nil
}

// Streams CallSummaries records from the API as a channel stream. This operation lazily loads records as efficiently as possible until the limit is reached.
func (c *ApiService) StreamCallSummaries(params *ListCallSummariesParams) (chan InsightsV1CallSummaries, chan error) {
	if params == nil {
		params = &ListCallSummariesParams{}
	}
	params.SetPageSize(client.ReadLimits(params.PageSize, params.Limit))

	recordChannel := make(chan InsightsV1CallSummaries, 1)
	errorChannel := make(chan error, 1)

	response, err := c.PageCallSummaries(params, "", "")
	if err != nil {
		errorChannel <- err
		close(recordChannel)
		close(errorChannel)
	} else {
		go c.streamCallSummaries(response, params, recordChannel, errorChannel)
	}

	return recordChannel, errorChannel
}

func (c *ApiService) streamCallSummaries(response *ListCallSummariesResponse, params *ListCallSummariesParams, recordChannel chan InsightsV1CallSummaries, errorChannel chan error) {
	curRecord := 1

	for response != nil {
		responseRecords := response.CallSummaries
		for item := range responseRecords {
			recordChannel <- responseRecords[item]
			curRecord += 1
			if params.Limit != nil && *params.Limit < curRecord {
				close(recordChannel)
				close(errorChannel)
				return
			}
		}

		record, err := client.GetNext(c.baseURL, response, c.getNextListCallSummariesResponse)
		if err != nil {
			errorChannel <- err
			break
		} else if record == nil {
			break
		}

		response = record.(*ListCallSummariesResponse)
	}

	close(recordChannel)
	close(errorChannel)
}

func (c *ApiService) getNextListCallSummariesResponse(nextPageUrl string) (interface{}, error) {
	if nextPageUrl == "" {
		return nil, nil
	}
	resp, err := c.requestHandler.Get(nextPageUrl, nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &ListCallSummariesResponse{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}
	return ps, nil
}
