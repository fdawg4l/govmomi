// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package settings

import "github.com/vmware/govmomi/vapi/rest"

// TaskSummary contains a brief summary about a task.
type TaskSummary struct {

	// The ID of the task.
	Id string `json:"id"`
}

type Progress string

const (
	Running   Progress = "RUNNING"
	Blocked   Progress = "BLOCKED"
	Succeeded Progress = "SUCCEEDED"
	Failed    Progress = "FAILED"
)

// TaskInfo contains information about a task and its subtasks of which it
// consists.
type TaskInfo struct {

	// Progress of the operation.
	Progress Progress

	// Information about the subtasks that this task contains. This field will
	// be unset if the task has no subtasks. This field will be unset if the
	// task is a bulk operation.
	Subtasks map[string]TaskInfo `json:"subtasks"`

	/**
	 * Summary about the subtasks for each entity if the current task is a bulk
	 * operation. The map is null or empty if the current task is not a
	 * bulk operation.
	 */
	SubtaskSummaries map[string][]TaskSummary `json:"subtask_summaries"`

	// Notifications to the user.
	Notifications rest.Notifications `json:"notifications"`

	// Task result.
	Result Progress

	/**
	 * Time when the task was last updated
	 */
	LastUpdateTime uint64 `json:"last_update_time"`
}
