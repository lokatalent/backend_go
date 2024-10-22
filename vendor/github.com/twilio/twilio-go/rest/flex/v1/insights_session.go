/*
 * This code was generated by
 * ___ _ _ _ _ _    _ ____    ____ ____ _    ____ ____ _  _ ____ ____ ____ ___ __   __
 *  |  | | | | |    | |  | __ |  | |__| | __ | __ |___ |\ | |___ |__/ |__|  | |  | |__/
 *  |  |_|_| | |___ | |__|    |__| |  | |    |__] |___ | \| |___ |  \ |  |  | |__| |  \
 *
 * Twilio - Flex
 * This is the public Twilio REST API.
 *
 * NOTE: This class is auto generated by OpenAPI Generator.
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

package openapi

import (
	"encoding/json"
	"net/url"
)

// Optional parameters for the method 'CreateInsightsSession'
type CreateInsightsSessionParams struct {
	// The Authorization HTTP request header
	Authorization *string `json:"Authorization,omitempty"`
}

func (params *CreateInsightsSessionParams) SetAuthorization(Authorization string) *CreateInsightsSessionParams {
	params.Authorization = &Authorization
	return params
}

// To obtain session details for fetching reports and dashboards
func (c *ApiService) CreateInsightsSession(params *CreateInsightsSessionParams) (*FlexV1InsightsSession, error) {
	path := "/v1/Insights/Session"

	data := url.Values{}
	headers := map[string]interface{}{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	if params != nil && params.Authorization != nil {
		headers["Authorization"] = *params.Authorization
	}
	resp, err := c.requestHandler.Post(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &FlexV1InsightsSession{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}
