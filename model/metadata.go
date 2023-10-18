package model

import "strings"

// Metadata -
type Metadata struct {
	ModelVersion             string               `json:"model_version,omitempty"`
	Ekn                      string               `json:"ekn,omitempty"`
	LastUpdated              string               `json:"last_updated,omitempty"`
	LastUpdateEventOccurred  string               `json:"last_update_event_occurred,omitempty"`
	LastUpdateEventID        string               `json:"last_update_event_id,omitempty"`
	OriginalSourceSystem     string               `json:"original_source_system,omitempty"`
	LastUpdateSourceSystem   string               `json:"last_update_source_system,omitempty"`
	CreatedFromFilename      string               `json:"created_from_filename,omitempty"`
	LastUpdateFromFilename   string               `json:"last_update_from_filename,omitempty"`
	CreatedFromTable         string               `json:"created_from_table,omitempty"`
	LastUpdatedFromTable     string               `json:"last_updated_from_table,omitempty"`
	DuplicateChecked         bool                 `json:"duplicate_checked,omitempty"`
	IsDuplicate              bool                 `json:"is_duplicate,omitempty"`
	DuplicateOf              string               `json:"duplicate_of,omitempty"`
	RequiresMerge            bool                 `json:"requires_merge,omitempty"`
	IsMerged                 bool                 `json:"is_merged,omitempty"`
	MergedEventID            string               `json:"merged_event_id,omitempty"`
	Created                  string               `json:"created,omitempty"`
	CreatedEventID           string               `json:"created_event_id,omitempty"`
	ActionType               string               `json:"action_type,omitempty"`
	IsMigrated               bool                 `json:"is_migrated,omitempty"`
	MigratedFromSystem       string               `json:"migrated_from_system,omitempty"`
	MigratedDate             string               `json:"migrated_date,omitempty"`
	MigratedEventID          string               `json:"migrated_event_id,omitempty"`
	MigratedFromFilename     string               `json:"migrated_from_filename,omitempty"`
	JobTrackerName           string               `json:"job_tracker_name,omitempty"`
	CreatedBy                string               `json:"created_by,omitempty"`
	LastUpdatedBy            string               `json:"last_updated_by,omitempty"`
	LastExportedDate         string               `json:"last_exported_date,omitempty"`
	FirstExportedDate        string               `json:"first_exported_date,omitempty"`
	SourceFileLineNumber     int                  `json:"source_file_line_number,omitempty"`
	SourceProcessName        string               `json:"source_process_name,omitempty"`
	SourceProcessExecutionID string               `json:"source_process_execution_id,omitempty"`
	SourceProcessStartTime   string               `json:"source_process_start_time,omitempty"`
	SourceLastUpdated        string               `json:"source_last_updated,omitempty"`
	SourceLastUpdatedBy      string               `json:"source_last_updated_by,omitempty"`
	ChangedSourceAttrs       []string             `json:"changed_source_attrs,omitempty"`
	AddedSourceAttrs         []string             `json:"added_source_attrs,omitempty"`
	RemovedSourceAttrs       []string             `json:"removed_source_attrs,omitempty"`
	SourceDataBefore         StringNameValuePairs `json:"source_data_before,omitempty"`
	SourceDataAfter          StringNameValuePairs `json:"source_data_after,omitempty"`
	LastChecked              string               `json:"last_checked,omitempty"`
	LastReplicationEventID   string               `json:"last_replication_event_id,omitempty"`
	ReplicationEventIDs      []string             `json:"replication_event_i_ds,omitempty"`
	ChangedAttrCols          string               `json:"changed_attr_cols,omitempty"`
	ChangedAttrVals          string               `json:"changed_attr_vals,omitempty"`
	IsDeleted                StringBool           `json:"is_deleted,omitempty"`
	SourceDataMd5Hash        string               `json:"source_data_md_5_hash,omitempty"`
	CreatedFromEventID       string               `json:"created_from_event_id,omitempty"`
	OriginalTargetStream     string               `json:"original_target_stream,omitempty"`
	ShouldBeParked           StringBool           `json:"should_be_parked,omitempty"`
	IsParked                 StringBool           `json:"is_parked,omitempty"`
	IsTest                   StringBool           `json:"is_test,omitempty"`
}

// StringBool -
type StringBool string

// Bool data structure
func (s StringBool) Bool() bool {
	switch strings.ToUpper(string(s)) {
	case "T", "TRUE", "Y", "YES", "1":
		return true
	}
	return false
}
func (s StringBool) String() string {
	return string(s)
}

func NewStringBool(b bool) StringBool {
	if b {
		return "Y"
	}
	return "N"
}

// SetLastUpdateEventOccurred setter
func (m *Metadata) SetLastUpdateEventOccurred(s string) {
	m.LastUpdateEventOccurred = s
}

// GetLastUpdateEventOccurred getter
func (m *Metadata) GetLastUpdateEventOccurred() string {
	return m.LastUpdateEventOccurred
}

// SetLastUpdateEventID setter
func (m *Metadata) SetLastUpdateEventID(s string) {
	m.LastUpdateEventID = s
}

// GetLastUpdateEventID getter
func (m *Metadata) GetLastUpdateEventID() string {
	return m.LastUpdateEventID
}

// SetLastUpdated setter
func (m *Metadata) SetLastUpdated(s string) {
	m.LastUpdated = s
}

// GetLastUpdated getter
func (m *Metadata) GetLastUpdated() string {
	return m.LastUpdated
}

// GetLastUpdatedBy getter
func (m *Metadata) GetLastUpdatedBy() string {
	return m.LastUpdatedBy
}

// SetLastUpdatedBy setter
func (m *Metadata) SetLastUpdatedBy(s string) {
	m.LastUpdatedBy = s
}

// SetLastReplicationEventID setter
func (m *Metadata) SetLastReplicationEventID(s string) {
	m.LastReplicationEventID = s
}

// GetLastReplicationEventID getter
func (m *Metadata) GetLastReplicationEventID() string {
	return m.LastReplicationEventID
}

// GetSourceLastUpdated getter
func (m *Metadata) GetSourceLastUpdated() string {
	return m.SourceLastUpdated
}

// SetSourceLastUpdated setter
func (m *Metadata) SetSourceLastUpdated(sld string) {
	m.SourceLastUpdated = sld
}

// GetReplicationIDs getter
func (m *Metadata) GetReplicationIDs() []string {
	if m.ReplicationEventIDs == nil {
		m.ReplicationEventIDs = []string{}
	}
	return m.ReplicationEventIDs
}

// AddReplicationIDIfNotExist setter
func (m *Metadata) AddReplicationIDIfNotExist(id string) bool {
	for _, r := range m.GetReplicationIDs() {
		if r == id {
			return false
		}
	}
	m.ReplicationEventIDs = append(m.ReplicationEventIDs, id)
	return true
}

// SetChangedSourceAttrs setter
func (m *Metadata) SetChangedSourceAttrs(s []string) {
	m.ChangedSourceAttrs = s
}

// GetChangedSourceAttrs getter
func (m *Metadata) GetChangedSourceAttrs() []string {
	return m.ChangedSourceAttrs
}

// GetSourceDataBefore getter
func (m *Metadata) GetSourceDataBefore() StringNameValuePairs {
	return m.SourceDataBefore
}

// GetSourceDataAfter getter
func (m *Metadata) GetSourceDataAfter() StringNameValuePairs {
	return m.SourceDataAfter
}

// StringNameValuePairs contains string name-value pairs in a list
type StringNameValuePairs []*NameValuePair

// NameValuePair data structure
type NameValuePair struct {
	Name string `json:"name,omitempty"`
	Val  string `json:"val,omitempty"`
}
