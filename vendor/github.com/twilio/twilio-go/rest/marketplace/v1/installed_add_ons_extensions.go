/*
 * This code was generated by
 * ___ _ _ _ _ _    _ ____    ____ ____ _    ____ ____ _  _ ____ ____ ____ ___ __   __
 *  |  | | | | |    | |  | __ |  | |__| | __ | __ |___ |\ | |___ |__/ |__|  | |  | |__/
 *  |  |_|_| | |___ | |__|    |__| |  | |    |__] |___ | \| |___ |  \ |  |  | |__| |  \
 *
 * Twilio - Marketplace
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
	"strings"

	"github.com/twilio/twilio-go/client"
)

// Fetch an instance of an Extension for the Installed Add-on.
func (c *ApiService) FetchInstalledAddOnExtension(InstalledAddOnSid string, Sid string) (*MarketplaceV1InstalledAddOnExtension, error) {
	path := "/v1/InstalledAddOns/{InstalledAddOnSid}/Extensions/{Sid}"
	path = strings.Replace(path, "{"+"InstalledAddOnSid"+"}", InstalledAddOnSid, -1)
	path = strings.Replace(path, "{"+"Sid"+"}", Sid, -1)

	data := url.Values{}
	headers := map[string]interface{}{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	resp, err := c.requestHandler.Get(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &MarketplaceV1InstalledAddOnExtension{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}

// Optional parameters for the method 'ListInstalledAddOnExtension'
type ListInstalledAddOnExtensionParams struct {
	// How many resources to return in each list page. The default is 50, and the maximum is 1000.
	PageSize *int `json:"PageSize,omitempty"`
	// Max number of records to return.
	Limit *int `json:"limit,omitempty"`
}

func (params *ListInstalledAddOnExtensionParams) SetPageSize(PageSize int) *ListInstalledAddOnExtensionParams {
	params.PageSize = &PageSize
	return params
}
func (params *ListInstalledAddOnExtensionParams) SetLimit(Limit int) *ListInstalledAddOnExtensionParams {
	params.Limit = &Limit
	return params
}

// Retrieve a single page of InstalledAddOnExtension records from the API. Request is executed immediately.
func (c *ApiService) PageInstalledAddOnExtension(InstalledAddOnSid string, params *ListInstalledAddOnExtensionParams, pageToken, pageNumber string) (*ListInstalledAddOnExtensionResponse, error) {
	path := "/v1/InstalledAddOns/{InstalledAddOnSid}/Extensions"

	path = strings.Replace(path, "{"+"InstalledAddOnSid"+"}", InstalledAddOnSid, -1)

	data := url.Values{}
	headers := map[string]interface{}{
		"Content-Type": "application/x-www-form-urlencoded",
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

	ps := &ListInstalledAddOnExtensionResponse{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}

// Lists InstalledAddOnExtension records from the API as a list. Unlike stream, this operation is eager and loads 'limit' records into memory before returning.
func (c *ApiService) ListInstalledAddOnExtension(InstalledAddOnSid string, params *ListInstalledAddOnExtensionParams) ([]MarketplaceV1InstalledAddOnExtension, error) {
	response, errors := c.StreamInstalledAddOnExtension(InstalledAddOnSid, params)

	records := make([]MarketplaceV1InstalledAddOnExtension, 0)
	for record := range response {
		records = append(records, record)
	}

	if err := <-errors; err != nil {
		return nil, err
	}

	return records, nil
}

// Streams InstalledAddOnExtension records from the API as a channel stream. This operation lazily loads records as efficiently as possible until the limit is reached.
func (c *ApiService) StreamInstalledAddOnExtension(InstalledAddOnSid string, params *ListInstalledAddOnExtensionParams) (chan MarketplaceV1InstalledAddOnExtension, chan error) {
	if params == nil {
		params = &ListInstalledAddOnExtensionParams{}
	}
	params.SetPageSize(client.ReadLimits(params.PageSize, params.Limit))

	recordChannel := make(chan MarketplaceV1InstalledAddOnExtension, 1)
	errorChannel := make(chan error, 1)

	response, err := c.PageInstalledAddOnExtension(InstalledAddOnSid, params, "", "")
	if err != nil {
		errorChannel <- err
		close(recordChannel)
		close(errorChannel)
	} else {
		go c.streamInstalledAddOnExtension(response, params, recordChannel, errorChannel)
	}

	return recordChannel, errorChannel
}

func (c *ApiService) streamInstalledAddOnExtension(response *ListInstalledAddOnExtensionResponse, params *ListInstalledAddOnExtensionParams, recordChannel chan MarketplaceV1InstalledAddOnExtension, errorChannel chan error) {
	curRecord := 1

	for response != nil {
		responseRecords := response.Extensions
		for item := range responseRecords {
			recordChannel <- responseRecords[item]
			curRecord += 1
			if params.Limit != nil && *params.Limit < curRecord {
				close(recordChannel)
				close(errorChannel)
				return
			}
		}

		record, err := client.GetNext(c.baseURL, response, c.getNextListInstalledAddOnExtensionResponse)
		if err != nil {
			errorChannel <- err
			break
		} else if record == nil {
			break
		}

		response = record.(*ListInstalledAddOnExtensionResponse)
	}

	close(recordChannel)
	close(errorChannel)
}

func (c *ApiService) getNextListInstalledAddOnExtensionResponse(nextPageUrl string) (interface{}, error) {
	if nextPageUrl == "" {
		return nil, nil
	}
	resp, err := c.requestHandler.Get(nextPageUrl, nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &ListInstalledAddOnExtensionResponse{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}
	return ps, nil
}

// Optional parameters for the method 'UpdateInstalledAddOnExtension'
type UpdateInstalledAddOnExtensionParams struct {
	// Whether the Extension should be invoked.
	Enabled *bool `json:"Enabled,omitempty"`
}

func (params *UpdateInstalledAddOnExtensionParams) SetEnabled(Enabled bool) *UpdateInstalledAddOnExtensionParams {
	params.Enabled = &Enabled
	return params
}

// Update an Extension for an Add-on installation.
func (c *ApiService) UpdateInstalledAddOnExtension(InstalledAddOnSid string, Sid string, params *UpdateInstalledAddOnExtensionParams) (*MarketplaceV1InstalledAddOnExtension, error) {
	path := "/v1/InstalledAddOns/{InstalledAddOnSid}/Extensions/{Sid}"
	path = strings.Replace(path, "{"+"InstalledAddOnSid"+"}", InstalledAddOnSid, -1)
	path = strings.Replace(path, "{"+"Sid"+"}", Sid, -1)

	data := url.Values{}
	headers := map[string]interface{}{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	if params != nil && params.Enabled != nil {
		data.Set("Enabled", fmt.Sprint(*params.Enabled))
	}

	resp, err := c.requestHandler.Post(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &MarketplaceV1InstalledAddOnExtension{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}
