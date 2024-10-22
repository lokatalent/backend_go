/*
 * This code was generated by
 * ___ _ _ _ _ _    _ ____    ____ ____ _    ____ ____ _  _ ____ ____ ____ ___ __   __
 *  |  | | | | |    | |  | __ |  | |__| | __ | __ |___ |\ | |___ |__/ |__|  | |  | |__/
 *  |  |_|_| | |___ | |__|    |__| |  | |    |__] |___ | \| |___ |  \ |  |  | |__| |  \
 *
 * Twilio - Taskrouter
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

// TaskrouterV1WorkspaceCumulativeStatistics struct for TaskrouterV1WorkspaceCumulativeStatistics
type TaskrouterV1WorkspaceCumulativeStatistics struct {
	// The SID of the [Account](https://www.twilio.com/docs/iam/api/account) that created the Workspace resource.
	AccountSid *string `json:"account_sid,omitempty"`
	// The average time in seconds between Task creation and acceptance.
	AvgTaskAcceptanceTime int `json:"avg_task_acceptance_time,omitempty"`
	// The beginning of the interval during which these statistics were calculated, in [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.
	StartTime *time.Time `json:"start_time,omitempty"`
	// The end of the interval during which these statistics were calculated, in [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.
	EndTime *time.Time `json:"end_time,omitempty"`
	// The total number of Reservations that were created for Workers.
	ReservationsCreated int `json:"reservations_created,omitempty"`
	// The total number of Reservations accepted by Workers.
	ReservationsAccepted int `json:"reservations_accepted,omitempty"`
	// The total number of Reservations that were rejected.
	ReservationsRejected int `json:"reservations_rejected,omitempty"`
	// The total number of Reservations that were timed out.
	ReservationsTimedOut int `json:"reservations_timed_out,omitempty"`
	// The total number of Reservations that were canceled.
	ReservationsCanceled int `json:"reservations_canceled,omitempty"`
	// The total number of Reservations that were rescinded.
	ReservationsRescinded int `json:"reservations_rescinded,omitempty"`
	// A list of objects that describe the number of Tasks canceled and reservations accepted above and below the thresholds specified in seconds.
	SplitByWaitTime *interface{} `json:"split_by_wait_time,omitempty"`
	// The wait duration statistics (`avg`, `min`, `max`, `total`) for Tasks that were accepted.
	WaitDurationUntilAccepted *interface{} `json:"wait_duration_until_accepted,omitempty"`
	// The wait duration statistics (`avg`, `min`, `max`, `total`) for Tasks that were canceled.
	WaitDurationUntilCanceled *interface{} `json:"wait_duration_until_canceled,omitempty"`
	// The total number of Tasks that were canceled.
	TasksCanceled int `json:"tasks_canceled,omitempty"`
	// The total number of Tasks that were completed.
	TasksCompleted int `json:"tasks_completed,omitempty"`
	// The total number of Tasks created.
	TasksCreated int `json:"tasks_created,omitempty"`
	// The total number of Tasks that were deleted.
	TasksDeleted int `json:"tasks_deleted,omitempty"`
	// The total number of Tasks that were moved from one queue to another.
	TasksMoved int `json:"tasks_moved,omitempty"`
	// The total number of Tasks that were timed out of their Workflows (and deleted).
	TasksTimedOutInWorkflow int `json:"tasks_timed_out_in_workflow,omitempty"`
	// The SID of the Workspace.
	WorkspaceSid *string `json:"workspace_sid,omitempty"`
	// The absolute URL of the Workspace statistics resource.
	Url *string `json:"url,omitempty"`
}