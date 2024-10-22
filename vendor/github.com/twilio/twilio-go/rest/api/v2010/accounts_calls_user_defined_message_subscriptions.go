/*
 * This code was generated by
 * ___ _ _ _ _ _    _ ____    ____ ____ _    ____ ____ _  _ ____ ____ ____ ___ __   __
 *  |  | | | | |    | |  | __ |  | |__| | __ | __ |___ |\ | |___ |__/ |__|  | |  | |__/
 *  |  |_|_| | |___ | |__|    |__| |  | |    |__] |___ | \| |___ |  \ |  |  | |__| |  \
 *
 * Twilio - Api
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
	"strings"
)

// Optional parameters for the method 'CreateUserDefinedMessageSubscription'
type CreateUserDefinedMessageSubscriptionParams struct {
	// The SID of the [Account](https://www.twilio.com/docs/iam/api/account) that subscribed to the User Defined Messages.
	PathAccountSid *string `json:"PathAccountSid,omitempty"`
	// The URL we should call using the `method` to send user defined events to your application. URLs must contain a valid hostname (underscores are not permitted).
	Callback *string `json:"Callback,omitempty"`
	// A unique string value to identify API call. This should be a unique string value per API call and can be a randomly generated.
	IdempotencyKey *string `json:"IdempotencyKey,omitempty"`
	// The HTTP method Twilio will use when requesting the above `Url`. Either `GET` or `POST`. Default is `POST`.
	Method *string `json:"Method,omitempty"`
}

func (params *CreateUserDefinedMessageSubscriptionParams) SetPathAccountSid(PathAccountSid string) *CreateUserDefinedMessageSubscriptionParams {
	params.PathAccountSid = &PathAccountSid
	return params
}
func (params *CreateUserDefinedMessageSubscriptionParams) SetCallback(Callback string) *CreateUserDefinedMessageSubscriptionParams {
	params.Callback = &Callback
	return params
}
func (params *CreateUserDefinedMessageSubscriptionParams) SetIdempotencyKey(IdempotencyKey string) *CreateUserDefinedMessageSubscriptionParams {
	params.IdempotencyKey = &IdempotencyKey
	return params
}
func (params *CreateUserDefinedMessageSubscriptionParams) SetMethod(Method string) *CreateUserDefinedMessageSubscriptionParams {
	params.Method = &Method
	return params
}

// Subscribe to User Defined Messages for a given Call SID.
func (c *ApiService) CreateUserDefinedMessageSubscription(CallSid string, params *CreateUserDefinedMessageSubscriptionParams) (*ApiV2010UserDefinedMessageSubscription, error) {
	path := "/2010-04-01/Accounts/{AccountSid}/Calls/{CallSid}/UserDefinedMessageSubscriptions.json"
	if params != nil && params.PathAccountSid != nil {
		path = strings.Replace(path, "{"+"AccountSid"+"}", *params.PathAccountSid, -1)
	} else {
		path = strings.Replace(path, "{"+"AccountSid"+"}", c.requestHandler.Client.AccountSid(), -1)
	}
	path = strings.Replace(path, "{"+"CallSid"+"}", CallSid, -1)

	data := url.Values{}
	headers := map[string]interface{}{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	if params != nil && params.Callback != nil {
		data.Set("Callback", *params.Callback)
	}
	if params != nil && params.IdempotencyKey != nil {
		data.Set("IdempotencyKey", *params.IdempotencyKey)
	}
	if params != nil && params.Method != nil {
		data.Set("Method", *params.Method)
	}

	resp, err := c.requestHandler.Post(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &ApiV2010UserDefinedMessageSubscription{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}

// Optional parameters for the method 'DeleteUserDefinedMessageSubscription'
type DeleteUserDefinedMessageSubscriptionParams struct {
	// The SID of the [Account](https://www.twilio.com/docs/iam/api/account) that subscribed to the User Defined Messages.
	PathAccountSid *string `json:"PathAccountSid,omitempty"`
}

func (params *DeleteUserDefinedMessageSubscriptionParams) SetPathAccountSid(PathAccountSid string) *DeleteUserDefinedMessageSubscriptionParams {
	params.PathAccountSid = &PathAccountSid
	return params
}

// Delete a specific User Defined Message Subscription.
func (c *ApiService) DeleteUserDefinedMessageSubscription(CallSid string, Sid string, params *DeleteUserDefinedMessageSubscriptionParams) error {
	path := "/2010-04-01/Accounts/{AccountSid}/Calls/{CallSid}/UserDefinedMessageSubscriptions/{Sid}.json"
	if params != nil && params.PathAccountSid != nil {
		path = strings.Replace(path, "{"+"AccountSid"+"}", *params.PathAccountSid, -1)
	} else {
		path = strings.Replace(path, "{"+"AccountSid"+"}", c.requestHandler.Client.AccountSid(), -1)
	}
	path = strings.Replace(path, "{"+"CallSid"+"}", CallSid, -1)
	path = strings.Replace(path, "{"+"Sid"+"}", Sid, -1)

	data := url.Values{}
	headers := map[string]interface{}{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	resp, err := c.requestHandler.Delete(c.baseURL+path, data, headers)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
