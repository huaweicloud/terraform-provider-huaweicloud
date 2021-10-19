package tables

import (
	"fmt"

	"github.com/chnsz/golangsdk"
)

const (
	TableTypeOBS  = "OBS"
	TableTypeDLI  = "DLI"
	TableTypeVIEW = "VIEW"
)

type CreateTableOpts struct {
	TableName string `json:"table_name" required:"true"`
	// Data storage location. OBS tables,DLI tables and view are available.
	DataLocation    string       `json:"data_location" required:"true"`
	Description     string       `json:"description,omitempty"`
	Columns         []ColumnOpts `json:"columns" required:"true"`
	SelectStatement string       `json:"select_statement,omitempty"`
	// Type of the data to be added to the OBS table. The options: Parquet, ORC, CSV, JSON, Carbon, and Avro.
	DataType string `json:"data_type,omitempty"`
	// Storage path of data in the new OBS table, which must be a path on OBS and must begin with obs. start with s3a
	DataPath string `json:"data_path,omitempty"`
	// Whether the table header is included in the OBS table data. Only data in CSV files has this attribute.
	WithColumnHeader *bool `json:"with_column_header,omitempty"`
	// User-defined data delimiter. Only data in CSV files has this attribute.
	Delimiter string `json:"delimiter,omitempty"`
	// User-defined reference character. Double quotation marks ("\") are used by default. Only data in CSV files
	// has this attribute.
	QuoteChar string `json:"quote_char,omitempty"`
	// User-defined escape character. Backslashes (\\) are used by default. Only data in CSV files has this attribute.
	EscapeChar string `json:"escape_char,omitempty"`
	// User-defined date type. yyyy-MM-dd is used by default. Only data in CSV and JSON files has this attribute.
	DateFormat string `json:"date_format,omitempty"`
	// User-defined timestamp type. yyyy-MM-dd HH:mm:ss is used by default. Only data in CSV and JSON files has
	// this attribute.
	TimestampFormat string `json:"timestamp_format,omitempty"`
}

type ColumnOpts struct {
	ColumnName        string `json:"column_name" required:"true"`
	Type              string `json:"type" required:"true"`
	Description       string `json:"description,omitempty"`
	IsPartitionColumn *bool  `json:"is_partition_column,omitempty"`
}

type ListOpts struct {
	Keyword          string `q:"keyword"`
	WithDetail       *bool  `q:"with-detail"`
	PageSize         *int   `q:"page-size"`
	CurrentPage      *int   `q:"current-page"`
	WithPriv         *bool  `q:"with-priv"`
	TableType        string `q:"table-type"`
	DatasourceType   string `q:"table-type"`
	WithoutTablemeta *bool  `q:"table-type"`
}

type PartitionsOpts struct {
	// default:100
	Limit  int `q:"limit"`
	Offset int `q:"offset"`
}

type UpdateOwnerOpts struct {
	NewOwner string `json:"new_owner" required:"true"`
}

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

func Create(c *golangsdk.ServiceClient, databaseName string, opts CreateTableOpts) (*CommonResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst CommonResp
	_, err = c.Post(createURL(c, databaseName), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func Delete(c *golangsdk.ServiceClient, databaseName string, tableName string, asyncFlag bool) (*DeleteResp, error) {
	url := deleteURL(c, databaseName, tableName)
	url += fmt.Sprintf("%s%t", "?async=", asyncFlag)
	var rst DeleteResp

	_, err := c.DeleteWithResponse(url, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func List(c *golangsdk.ServiceClient, databaseName string, opts ListOpts) (*ListResp, error) {
	url := listURL(c, databaseName)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst ListResp
	_, err = c.Get(url, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func Get(c *golangsdk.ServiceClient, databaseName string, tableName string) (*Table, error) {
	var rst Table

	_, err := c.Get(getURL(c, databaseName, tableName), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func Partitions(c *golangsdk.ServiceClient, databaseName string, tableName string, opts PartitionsOpts) (*PartitionsResp, error) {
	url := partitionsURL(c, databaseName, tableName)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst PartitionsResp
	_, err = c.Get(url, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func UpdateOwner(c *golangsdk.ServiceClient, databaseName string, tableName string, opts UpdateOwnerOpts) (*CommonResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst CommonResp
	_, err = c.Put(updateOwnerURL(c, databaseName, tableName), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}
