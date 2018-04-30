package event

// FilterCompleted is the structure of each event consumed.
type FilterCompleted struct {
	FilterID  string `avro:"filter_output_id"`
	DatasetID string `avro:"dataset_id"`
	Edition   string `avro:"edition"`
	Version   string `avro:"version"`
	Email     string `avro:"email"`
}
