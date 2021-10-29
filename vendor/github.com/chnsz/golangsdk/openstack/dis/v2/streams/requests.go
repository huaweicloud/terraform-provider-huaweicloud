package streams

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

const (
	StreamTypeCommon   = "COMMON"
	StreamTypeAdvanced = "ADVANCED"

	DataTypeBlob = "BLOB"
	DataTypeJson = "JSON"
	DataTypeCsv  = "CSV"

	CompressionTypeSnappy = "snappy"
	CompressionTypeGzip   = "gzip"
	CompressionTypeZip    = "zip"

	StatusCreating    = "CREATING"
	StatusRunning     = "RUNNING"
	StatusTerminating = "TERMINATING"
	StatusTerminated  = "TERMINATED"
	StatusFrozen      = "FROZEN"
)

type CreateOpts struct {
	// Each DIS stream has a unique name. A stream name is 1 to 64 characters in length. Only letters, digits,
	// hyphens (-), and underscores (_) are allowed.
	StreamName string `json:"stream_name" required:"true"`
	// Quantity of the partitions into which data records in the newly created DIS stream will be distributed.
	// Partitions are the base throughput unit of a DIS stream.
	// The value range varies depending on the value of stream_type.
	// If stream_type is not specified or set to COMMON, the value of partition_count is an int from 1 to 50.
	// If the tenant has created N common partitions, the maximum value of partition_count is 50-N.
	// If stream_type is set to ADVANCED, the value of partition_count is an int from 1 to 10.
	// If the tenant has created N advanced partitions, the maximum value of partition_count is 10-N.
	PartitionCount int `json:"partition_count" required:"true"`
	// Stream type. Possible values:
	// COMMON: a common stream. The bandwidth is 1 MB/s.
	// ADVANCED: an advanced stream. The bandwidth is 5 MB/s.
	// Default value: COMMON
	StreamType string `json:"stream_type,omitempty"`
	// Source data type.
	// Possible values:
	// BLOB: a collection of binary data stored as a single entity in a database management system.
	// JSON: an open-standard file format that uses human-readable text to transmit data objects consisting of attribute–value pairs and array data types.
	// CSV: a simple text format for storing tabular data in a plain text file.
	// Default value: BLOB
	DataType string `json:"data_type,omitempty"`
	// Period of time for which data is retained in the DIS stream.
	// Value range: 24 to 72
	// Unit: hour
	// Default value: 24
	// If this parameter is left unspecified, the default value will be used.
	DataDuration int `json:"data_duration,omitempty"`
	// Whether to enable automatic Scaling
	// true: Turn on automatic scale out.
	// false: Turn off automatic scale in.
	// Disabled by default.
	// Default value: false
	AutoScaleEnabled *bool `json:"auto_scale_enabled,omitempty"`
	// When auto scaling is enabled, the minimum number of slices for auto scaling.
	// Minimum value: 1
	AutoScaleMinPartitionCount *int `json:"auto_scale_min_partition_count,omitempty"`
	// When auto scaling is enabled, the maximum number of slices for auto scaling.
	AutoScaleMaxPartitionCount *int `json:"auto_scale_max_partition_count,omitempty"`
	// Source data structure that defines JOSN and CSV formats. It is described in the syntax of Avro.
	// For details about Avro, see http://avro.apache.org/docs/current/#schemas.
	// NOTE:
	// This parameter is mandatory when Dump File Format is Parquet or CarbonData.
	DataSchema string `json:"data_schema,omitempty"`
	// Related attributes of CSV format data, such as delimiter
	CsvProperties *CsvProperty `json:"csv_properties,omitempty"`
	// Data compression type, currently supports:
	// snappy
	// gzip
	// zip
	// No compression by default
	CompressionFormat string             `json:"compression_format,omitempty"`
	Tags              []tags.ResourceTag `json:"tags,omitempty"`
	//the key must be:"_sys_enterprise_project_id"
	SysTags []tags.ResourceTag `json:"sys_tags,omitempty"`
}

type CsvProperty struct {
	Delimiter string `json:"delimiter"`
}

type GetOpts struct {
	StartPartitionId string `q:"start_partitionId"`
	// maximum number of partitions per request
	//   value range:1~1000。
	//   default:100。
	LimitPartitions int `q:"limit_partitions"`
}

type ListStreamsOpts struct {
	// maximum number of stream per request
	// value range:1~100。
	// default:10。
	Limit int `q:"limit"`
	// Return the stream list from this stream, the returned stream list does not include this stream name.
	// do not pass this field when querying on the first page. When the returned result has_more_streams is true,
	// the next page query is performed, and start_stream_name is the last stream name of the first page query result.
	StartStreamName string `q:"start_stream_name"`
}

type UpdatePartitionOpt struct {
	StreamName string `json:"stream_name" required:"true"`
	// The number of target partitions to be changed.
	// The value is an integer greater than 0.
	// The set value greater than the current number of partitions means expansion,
	// and less than the current number of partitions means shrinking.
	//Notice:
	// The total number of expansion and contraction times for each channel in an hour is up to 5 times,
	// and if the expansion or contraction operation succeeds once within an hour, the expansion or contraction
	// operation is not allowed in the last hour.
	// Minimum value: 0
	TargetPartitionCount int `json:"target_partition_count" required:"true"`
}

type CreatePolicyOpt struct {
	StreamId string `json:"stream_id" required:"true"`
	// Authorized users.
	// If authorized to the specified tenant, the format is: domainName.*;
	// if authorized to the specified sub-user under the tenant, the format is: domainName.userName;
	//Support for adding multiple accounts, separated by ",", for example: domainName1.userName1,domainName2.userName2;
	PrincipalName string `json:"principal_name" required:"true"`
	// Authorized operation type: putRecords,getRecords
	ActionType string `json:"action_type" required:"true"`
	//Authorization impact type: accept,
	Effect string `json:"effect" required:"true"`
}

func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*golangsdk.Result, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r golangsdk.Result
	_, r.Err = c.Post(rootURL(c), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"region":       c.AKSKAuthOptions.Region,
		},
	})
	return &r, r.Err
}

func Get(c *golangsdk.ServiceClient, streamName string, opts GetOpts) (*StreamDetail, error) {
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}

	url := resourceURL(c, streamName)
	url += query.String()

	var rst StreamDetail
	_, err = c.Get(url, &rst, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"region":       c.AKSKAuthOptions.Region,
		},
	})

	if err == nil {
		return &rst, nil
	}
	return nil, err
}

func Delete(c *golangsdk.ServiceClient, streamName string) *golangsdk.ErrResult {
	var r golangsdk.ErrResult
	_, r.Err = c.Delete(resourceURL(c, streamName), &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"region":       c.AKSKAuthOptions.Region,
		},
	})
	return &r
}

func List(c *golangsdk.ServiceClient, opts ListStreamsOpts) (*ListResult, error) {
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}

	url := rootURL(c) + query.String()

	var rst ListResult
	_, err = c.Get(url, &rst, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"region":       c.AKSKAuthOptions.Region,
		},
	})
	if err == nil {
		return &rst, nil
	}
	return nil, err
}

func UpdatePartition(c *golangsdk.ServiceClient, name string, opts UpdatePartitionOpt) (*golangsdk.Result, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r golangsdk.Result
	_, err = c.Put(resourceURL(c, name), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"region":       c.AKSKAuthOptions.Region,
		},
	})
	return &r, err
}

func CreatePolicy(c *golangsdk.ServiceClient, streamName string, opts CreatePolicyOpt) (*golangsdk.Result, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r golangsdk.Result
	_, err = c.Post(policiesURL(c, streamName), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"region":       c.AKSKAuthOptions.Region,
		},
	})
	return &r, err
}

func ListPolicies(c *golangsdk.ServiceClient, streamName string) (*ListPolicyResult, error) {
	var rst ListPolicyResult
	_, err := c.Get(policiesURL(c, streamName), &rst, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"region":       c.AKSKAuthOptions.Region,
		},
	})
	if err == nil {
		return &rst, nil
	}
	return nil, err
}
