/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OSIndexPolicySpec defines the desired state of OSIndexPolicy.
type OSIndexPolicySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// PolicyID is the unique identifier for the Opensearch Index ISM policy
	PolicyID string `json:"policy_id,omitempty"`
	// Target Opensearch
	OpensearhConnection OpensearhConnection `json:"opensearch_connection,omitempty"`
	// IndexPolicy defines the ISM policy for the index
	Policy OpensearchIndexPolicy `json:"policy,omitempty"`
}

type OpensearhConnection struct {
	// URL of the Opensearch instance
	URL string `json:"url,omitempty"`
	// Username for authentication
	Username string `json:"username,omitempty"`
	// Password for authentication
	Password string `json:"password,omitempty"`
}

// OpensearchIndexPolicy define the desired state of Opensearch Index ISM policy
type OpensearchIndexPolicy struct {
	Description string `json:"description,omitempty"`
	// LastUpdatedTime   time.Time         `json:"last_updated_time,omitempty"`
	ErrorNotification map[string]string `json:"error_notification,omitempty"`
	DefaultState      string            `json:"default_state,omitempty"`
	States            []*State          `json:"states,omitempty"`
	// ISMTemplate defines the template for the index
	ISMTemplate *ISMTemplate `json:"ism_template,omitempty"`
}

type ISMTemplate struct {
	IndexPatterns []string `json:"index_patterns,omitempty"`
	Priority      int      `json:"priority,omitempty"`
}

// State defines a state in the ISM policy
// It includes the name of the state, the actions to be performed in this state,
// and the transitions to other states
type State struct {
	Name        string        `json:"name,omitempty"`
	Actions     []*Action     `json:"actions,omitempty"`
	Transitions []*Transition `json:"transitions,omitempty"`
}

// Transition defines the transition from one state to another
// It includes the state name to transition to and the conditions that must be met for the transition to occur
type Transition struct {
	// StateName is the name of the state to transition to
	StateName string `json:"state_name,omitempty"`
	// Conditions are the conditions that must be met for the transition to occur
	Conditions map[string]string `json:"conditions,omitempty"`
}

type Action struct {
	// DeleteAction defines the action to delete the index
	Delete *DeleteAction `json:"delete,omitempty"`
	// ForceMergeAction defines the action to force merge the index
	ForceMerge *ForceMergeAction `json:"force_merge,omitempty"`
	// ReadOnlyAction defines the action to make the index read-only
	ReadOnly *ReadOnlyAction `json:"read_only,omitempty"`
	// RollOverAction defines the action to roll over the index
	RollOver *RollOverAction `json:"rollover,omitempty"`
	// SnapshotAction defines the action to take a snapshot of the index
	Snapshot *SnapshotAction `json:"snapshot,omitempty"`
	// ReadWriteAction defines the action to make the index read-write
	ReadWrite *ReadWriteAction `json:"read_write,omitempty"`
	// ReplicaCountAction defines the action to set the number of replicas for the index
	ReplicaCount *ReplicaCountAction `json:"replica_count,omitempty"`
	// ShrinkAction defines the action to shrink the index
	Shrink *ShrinkAction `json:"shrink,omitempty"`
	// CloseAction defines the action to close the index
	Close *CloseAction `json:"close,omitempty"`
	// OpenAction defines the action to open the index
	Open *OpenAction `json:"open,omitempty"`
	// NotificationAction defines the action to notify about the index state
	Notification *NotifyAction `json:"notification,omitempty"`
	// ConvertIndexToRemoteAction defines the action to convert the index to removed state
	ConvertIndexToRemote *ConvertIndexToRemoteAction `json:"convert_index_to_remote,omitempty"`
	// IndexPriorityAction defines the action to set the index priority
	IndexPriority *IndexPriorityAction `json:"index_priority,omitempty"`
	// AllocationAction defines the action to set the index allocation
	Allocation *AllocationAction `json:"allocation,omitempty"`
	// RollupAction defines the action to roll up the index
	Rollup *RollupAction `json:"rollup,omitempty"`
	// StopReplicationAction defines the action to stop replication of the index
	StopReplication *StopReplicationAction `json:"stop_replication,omitempty"`
}

type DeleteAction struct {
}

type ForceMergeAction struct {
	ForceMerge *ForceMerge `json:"force_merge,omitempty"`
}

type ReadOnlyAction struct {
}

type RollOverAction struct {
	MinSize             string `json:"min_size,omitempty"`
	MinPrimaryShardSize string `json:"min_primary_shard_size,omitempty"`
	MinDocCount         int    `json:"min_doc_count,omitempty"`
	MinIndexAge         string `json:"min_index_age,omitempty"`
	CopyAlias           bool   `json:"copy_alias,omitempty"`
}

type SnapshotAction struct {
	Repository string `json:"repository,omitempty"`
	Snapshot   string `json:"snapshot,omitempty"`
}

type ReadWriteAction struct {
}

type ReplicaCountAction struct {
	NumberOfReplicas int `json:"number_of_replicas,omitempty"`
}

type ShrinkAction struct {
}

type CloseAction struct {
}

type OpenAction struct {
}
type NotifyAction struct {
}
type ConvertIndexToRemoteAction struct {
}
type IndexPriorityAction struct {
}
type AllocationAction struct {
}
type RollupAction struct {
}
type StopReplicationAction struct {
}

type ForceMerge struct {
	// MaxNumSegments is the maximum number of segments to merge into
	MaxNumSegments       int    `json:"max_num_segments,omitempty"`
	WaitForCompletion    bool   `json:"wait_for_completion,omitempty"`
	TaskExecutionTimeout string `json:"task_execution_timeout,omitempty"`
}

// OSIndexPolicyStatus defines the observed state of OSIndexPolicy.
type OSIndexPolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// OSIndexPolicy is the Schema for the osindexpolicies API.
type OSIndexPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OSIndexPolicySpec   `json:"spec,omitempty"`
	Status OSIndexPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OSIndexPolicyList contains a list of OSIndexPolicy.
type OSIndexPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OSIndexPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OSIndexPolicy{}, &OSIndexPolicyList{})
}
