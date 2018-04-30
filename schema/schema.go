package schema

import (
	"github.com/ONSdigital/go-ns/avro"
)

var filterCompletedEvent = `{
  "type": "record",
  "name": "filter-output-completed",
  "fields": [
    {"name": "filter_output_id", "type": "string", "default": ""},
    {"name": "email", "type": "string", "default": ""},
    {"name": "dataset_id", "type": "string", "default": ""},
    {"name": "edition", "type": "string", "default": ""},
    {"name": "version", "type": "string", "default": ""}
  ]
}`

// FilterCompletedEvent the Avro schema for FilterOutputSubmitted messages.
var FilterCompletedEvent = &avro.Schema{
	Definition: filterCompletedEvent,
}
