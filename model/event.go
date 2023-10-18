package model

type EventHeader struct {
	EventName     string `json:"event_name,omitempty"`
	EventCategory string `json:"event_category,omitempty"`
	EventSource   string `json:"event_source,omitempty"` // Should match an Account.OperatorId
}

// Event -
type Event struct {
	EventHeader
	Metadata                 *Metadata         `json:"metadata,omitempty"`
	ID                       string            `json:"id,omitempty"`
	Reference                string            `json:"reference,omitempty"`
	BackOfficeUserID         string            `json:"back_office_user_id,omitempty"` // for cases where events are raised administratively
	EventOccurredTime        string            `json:"event_occurred_time,omitempty"`
	EventSourceSurrogateFor  string            `json:"event_source_surrogate_for,omitempty"` // the source that EventSource is acting for
	EventSourceFileName      string            `json:"event_source_file_name,omitempty"`
	EventSourceTableName     string            `json:"event_source_table_name,omitempty"`
	EventReceivedTime        string            `json:"event_received_time,omitempty"`
	EventProcessedTime       string            `json:"event_processed_time,omitempty"`
	EventPersisted           bool              `json:"event_persisted,omitempty"`
	RecordUpdated            bool              `json:"record_updated,omitempty"`
	RequiresRecordUpdate     bool              `json:"requires_record_update,omitempty"`
	SensitiveFieldsEncrypted bool              `json:"sensitive_fields_encrypted,omitempty"`
	ExternalEventID          string            `json:"external_event_id,omitempty"`
	Data                     map[string]string `json:"data,omitempty"`
	SensitiveDataFieldNames  []string          `json:"sensitive_data_field_names,omitempty"`
	MessageContext           *MessageContext   `json:"message_context,omitempty"`
	SourceMessageContext     *MessageContext   `json:"source_message_context,omitempty"`
	ChangeLog                []*LogMessage     `json:"change_log,omitempty"`
	TriggerEventID           string            `json:"trigger_event_id,omitempty"`
	EventContext             *EventContext     `json:"event_context,omitempty"`
	RetryCount               int               `json:"retry_count,omitempty"`
	// MultipartEventTrigger    *MultipartEvent
	ContainsEntityID    string                 `json:"contains_entity_id,omitempty"`
	DerivedEventDetails []*DerivedEventDetails `json:"derived_event_details,omitempty"`
	DoNotPersistMe      bool                   `json:"do_not_persist_me,omitempty"`
	// TODO - Emailer                  *email.SesEmailer      `json:"emailer,omitempty" bigquery:"-"`
}

// LogMessage -
type LogMessage struct {
	ID                     string                  `json:"id,omitempty"`
	DocType                string                  `json:"doc_type,omitempty"`
	LogType                string                  `json:"log_type,omitempty"`
	Timestamp              string                  `json:"timestamp,omitempty"`
	Metadata               *Metadata               `json:"metadata,omitempty"`
	MessageContext         *MessageContext         `json:"message_context,omitempty"`
	SourceMessageContext   *MessageContext         `json:"source_message_context,omitempty"`
	ChangeLog              *ReferenceDataChangeLog `json:"change_log,omitempty"`
	ProcessingLog          *MessageProcessingLog   `json:"processing_log,omitempty"`
	SensitiveDataAccessLog *SensitiveDataAccessLog `json:"sensitive_data_access_log,omitempty"`
}

// ReferenceDataChangeLog -
type ReferenceDataChangeLog struct {
	CustomerID           string      `json:"customer_id,omitempty"`
	EventID              string      `json:"event_id,omitempty"`
	EventName            string      `json:"event_name,omitempty"`
	PartyID              string      `json:"party_id,omitempty"`
	CampaignID           string      `json:"campaign_id,omitempty"`
	ProductID            string      `json:"product_id,omitempty"`
	ReferenceDataDocType string      `json:"reference_data_doc_type,omitempty"`
	AttributeChanged     string      `json:"attribute_changed,omitempty"`
	OriginalValue        interface{} `json:"original_value,omitempty" bigquery:"-"`
	NewValue             interface{} `json:"new_value,omitempty" bigquery:"-"`
}

// SensitiveDataAccessLog -
type SensitiveDataAccessLog struct {
	AccessedIDs    []string `json:"accessed_ids,omitempty"`
	AccessedByUser string   `json:"accessed_by_user,omitempty"`
	ClientID       string   `json:"client_id,omitempty"`
	RequestID      string   `json:"request_id,omitempty"`
}

// MessageProcessingLog -
type MessageProcessingLog struct {
	ProcessedEventID   string `json:"processed_event_id,omitempty"`
	ProcessedEventName string `json:"processed_event_name,omitempty"`
	ProcessorID        string `json:"processor_id,omitempty"`
	Status             string `json:"status,omitempty"`
	Error              string `json:"error,omitempty"`
}

type DerivedEventDetails struct {
	ID                string `json:"id,omitempty"`
	EventName         string `json:"event_name,omitempty"`
	Category          string `json:"category,omitempty"`
	EventOccurredTime string `json:"event_occurred_time,omitempty"`
}

// EventContext -
type EventContext struct {
	Location *Location `json:"location,omitempty" bigquery:",nullable"`
	Device   *Device   `json:"device,omitempty" bigquery:",nullable"`
}

// Location -
type Location struct {
	LocationID       string           `json:"location_id,omitempty"`
	LocationType     string           `json:"location_type,omitempty"`
	LocationName     string           `json:"location_name,omitempty"`
	Source           string           `json:"source,omitempty"`
	Channel          *Channel         `json:"channel,omitempty"`
	Territory        string           `json:"territory,omitempty"`
	Timestamp        string           `json:"timestamp,omitempty"`
	ExternalID       *ExternalID      `json:"external_id,omitempty"`
	DerivedFrom      string           `json:"derived_from,omitempty"`
	Lat              string           `json:"lat,omitempty"`
	Long             string           `json:"long,omitempty"`
	Flat             string           `json:"flat,omitempty"`
	HouseNameNumber  string           `json:"house_name_number,omitempty"`
	AddressLine1     string           `json:"address_line_1,omitempty"`
	AddressLine2     string           `json:"address_line_2,omitempty"`
	TownCity         string           `json:"town_city,omitempty"`
	Country          string           `json:"country,omitempty"`
	CountryCode      string           `json:"country_code,omitempty"`
	Area             string           `json:"area,omitempty"`
	State            string           `json:"state,omitempty"`
	Postcode         string           `json:"postcode,omitempty"`
	PostalOutCode    string           `json:"postal_out_code,omitempty"`
	Description      string           `json:"description,omitempty"`
	IPAddress        string           `json:"ip_address,omitempty"`
	GeoJSON          string           `json:"geo_json,omitempty"`
	GeoScheme        string           `json:"geo_scheme,omitempty"`
	CustomAttributes []*NameValuePair `json:"custom_attributes,omitempty"`
}

// Channel data structure
type Channel struct {
	ChannelCode        string        `json:"channel_code,omitempty"`
	ChannelName        string        `json:"channel_name,omitempty"`
	ChannelInstanceID  string        `json:"channel_instance_id,omitempty"` // e.g. Retail Store ID, Website ID
	ChannelCountry     string        `json:"channel_country,omitempty"`
	ChannelExternalIDs []*ExternalID `json:"channel_external_i_ds,omitempty"`
	ChannelSourceID    string        `json:"channel_source_id,omitempty"`
	ChannelSourceCode  string        `json:"channel_source_code,omitempty"`
}

// ExternalID data structure
type ExternalID struct {
	System string `json:"system,omitempty"`
	ID     string `json:"id,omitempty"`
}

// Device -
type Device struct {
	DeviceType     string `json:"device_type,omitempty"`
	ID             string `json:"id,omitempty"`
	Os             string `json:"os,omitempty"`
	OsVersion      string `json:"os_version,omitempty"`
	Model          string `json:"model,omitempty"`
	Brand          string `json:"brand,omitempty"`
	Browser        string `json:"browser,omitempty"`
	BrowserType    string `json:"browser_type,omitempty"`
	BrowserVersion string `json:"browser_version,omitempty"`
	UserAgent      string `json:"user_agent,omitempty"`
	Source         string `json:"source,omitempty"`
	DateAdded      string `json:"date_added,omitempty"`
	DateUpdated    string `json:"date_updated,omitempty"`
	BrowserWidth   string `json:"browser_width,omitempty"`
	BrowserHeight  string `json:"browser_height,omitempty"`
}

// MessageContext data structure
type MessageContext struct {
	TopicName string `json:"topic_name,omitempty"`
	Partition int32  `json:"partition,omitempty"`
	Offset    int64  `json:"offset,omitempty"`
}

// AddDerivedEventDetailsIfNotExist appends derived details if they don't exist
func (e *Event) AddDerivedEventDetailsIfNotExist(ev *DerivedEventDetails) bool {
	for _, ded := range e.DerivedEventDetails {
		if ev.EventName == ded.EventName {
			return false
		}
	}
	e.DerivedEventDetails = append(e.DerivedEventDetails, ev)
	return true
}
