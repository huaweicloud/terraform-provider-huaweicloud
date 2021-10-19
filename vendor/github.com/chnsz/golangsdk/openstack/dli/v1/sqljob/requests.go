package sqljob

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

const (

	// Maximum number of records to be written into a single file. If the value is zero or negative, there is no limit.
	ConfigSparkSqlFilesMaxRecordsPerFile = "spark.sql.files.maxRecordsPerFile"
	// Maximum size of the table that displays all working nodes when a connection is executed.
	// You can set this parameter to -1 to disable the display.
	// NOTE:
	// only the configuration unit metastore table that runs the ANALYZE TABLE COMPUTE statistics noscan command and
	// the file-based data source table that directly calculates statistics based on data files are supported.
	ConfigSparkSqlAutoBroadcastJoinThreshold = "spark.sql.autoBroadcastJoinThreshold"
	// Default number of partitions used to filter data for join or aggregation.
	ConfigSparkSqlShufflePartitions = "spark.sql.shuffle.partitions"
	// In dynamic mode, Spark does not delete the previous partitions and only overwrites the partitions without
	// data during execution.
	ConfigSparkSqlDynamicPartitionOverwriteEnabled = "spark.sql.dynamicPartitionOverwrite.enabled"
	// Maximum number of bytes to be packed into a single partition when a file is read.
	ConfigSparkSqlMaxPartitionBytes = "spark.sql.files.maxPartitionBytes"
	// Path of bad records.
	ConfigSparkSqlBadRecordsPath = "spark.sql.badRecordsPath"
	// Indicates whether DDL and DCL statements are executed asynchronously. The value true indicates that
	// asynchronous execution is enabled.
	ConfigDliSqlasyncEnabled = "dli.sql.sqlasync.enabled"
	// Sets the job running timeout interval. If the timeout interval expires, the job is canceled. Unit: ms.
	ConfigDliSqljobTimeout = "dli.sql.job.timeout"
)

type SqlJobOpts struct {
	// SQL statement that you want to execute.
	Sql string `json:"sql" required:"true"`
	// Database where the SQL is executed. This parameter does not need to be configured during database creation.
	Currentdb string `json:"currentdb,omitempty"`
	QueueName string `json:"queue_name,omitempty"`
	// You can set the configuration parameters for the SQL job in the form of Key/Value
	Conf []string           `json:"conf,omitempty"`
	Tags []tags.ResourceTag `json:"tags,omitempty"`
}

// ListJobsOpts
type ListJobsOpts struct {
	// Maximum number of jobs displayed on each page. The value range is as follows: [1, 100]. The default value is 50.
	PageSize *int `q:"page-size"`
	// Current page number. The default value is 1.
	CurrentPage *int `q:"current-page"`
	// Queries the jobs executed later than the time. The time is a UNIX timestamp in milliseconds.
	Start *int `q:"start"`
	// Queries the jobs executed earlier than the time. The time is a UNIX timestamp in milliseconds.
	End *int `q:"end"`
	// Type of a job to be queried. Job types include:DDL、DCL、IMPORT、EXPORT、QUERY、INSERT、DATA_MIGRATION、UPDATE、
	// DELETE、RESTART_QUEUE、SCALE_QUEUE, To query all types of jobs, enter ALL
	JobType   string `q:"job-type"`
	JobStatus string `q:"job-status"`
	JobId     string `q:"job-id"`
	DbName    string `q:"db_name"`
	TableName string `q:"table_name"`
	// Specifies queue_name as the filter to query jobs running on the specified queue.
	QueueName string `q:"queue_name"`
	// Specifies the SQL segment as the filter. It is case insensitive.
	SqlPattern string `q:"sql_pattern"`
	// Specifies the job sorting mode. The default value is start_time_desc (job submission time in descending order).
	// Four sorting modes are supported: duration_desc (job running duration in descending order),
	// duration_asc (job running duration in ascending order),
	// start_time_desc (job submission time in descending order),
	// and start_time_asc (job submission time in ascending order).
	Order      string `q:"order"`
	EngineType string `q:"engine-type"`
}

type CheckSQLGramarOpts struct {
	// SQL statement that you want to execute.
	Sql string `json:"sql" required:"true"`
	// Database where the SQL statement is executed.
	// NOTE:
	// If the SQL statement contains db_name, for example, select * from db1.t1, you do not need to set this parameter.
	// If the SQL statement does not contain db_name, the semantics check will fail when you do not set this parameter
	// or set this parameter to an incorrect value.
	Currentdb string `json:"currentdb"`
}

type ExportQueryResultOpts struct {
	// Path for storing the exported data. Currently, data can be stored only on OBS.
	// The OBS path cannot contain folders, for example, the path folder in the sample request.
	DataPath string `json:"data_path" required:"true"`
	// Compression format of exported data. Currently, gzip, bzip2, and deflate are supported.
	// The default value is none, indicating that data is not compressed.
	Compress string `json:"compress,omitempty"`
	// Storage format of exported data. Currently, only CSV and JSON are supported.
	DataType string `json:"data_type" required:"true"`
	// Name of the queue that is specified to execute a task. If no queue is specified, the default queue is used.
	QueueName string `json:"queue_name,omitempty"`
	// Export mode. The parameter value can be ErrorIfExists or Overwrite.
	// If export_mode is not specified, this parameter is set to ErrorIfExists by default.
	// ErrorIfExists: Ensure that the specified export directory does not exist.
	// If the specified export directory exists, an error is reported and the export operation cannot be performed.
	// Overwrite: If you add new files to a specific directory, existing files will be deleted.
	ExportMode string `json:"export_mode,omitempty"`
	// Whether to export column names when exporting CSV and JSON data.
	// If this parameter is set to true, the column names are exported.
	// If this parameter is set to false, the column names are not exported.
	// If this parameter is left blank, the default value false is used.
	WithColumnHeader *bool `json:"with_column_header,omitempty"`
	// Number of data records to be exported. The default value is 0, indicating that all data records are exported.
	LimitNum *int `json:"limit_num,omitempty"`
}

type ImportDataOpts struct {
	// Path to the data to be imported. Currently, only OBS data can be imported.
	DataPath string `json:"data_path" required:"true"`
	// Type of the data to be imported. Currently, data types of CSV, Parquet, ORC, JSON, and Avro are supported.
	// NOTE:
	// Data in Avro format generated by Hive tables cannot be imported.
	DataType string `json:"data_type" required:"true"`
	// Name of the database where the table to which data is imported resides.
	DatabaseName string `json:"database_name" required:"true"`
	// Name of the table to which data is imported.
	TableName string `json:"table_name" required:"true"`
	// Whether the first line of the imported data contains column names, that is, headers. The default value is false, indicating that column names are not contained. This parameter can be specified when CSV data is imported.
	WithColumnHeader *bool `json:"with_column_header,omitempty"`
	// User-defined data delimiter. The default value is a comma (,). This parameter can be specified when CSV data is imported.
	Delimiter string `json:"delimiter,omitempty"`
	// User-defined quotation character. The default value is double quotation marks ("). This parameter can be specified when CSV data is imported.
	QuoteChar string `json:"quote_char,omitempty"`
	// User-defined escape character. The default value is a backslash (\). This parameter can be specified when CSV data is imported.
	EscapeChar string `json:"escape_char,omitempty"`
	// Specified date format. The default value is yyyy-MM-dd. For details about the characters involved in the date format, see Table 3. This parameter can be specified when data in the CSV or JSON format is imported.
	DateFormat string `json:"date_format,omitempty"`
	// Bad records storage directory during job execution. After configuring this item, the bad records is not imported into the target table.
	BadRecordsPath string `json:"bad_records_path,omitempty"`
	// Specified time format. The default value is yyyy-MM-dd HH:mm:ss. For definitions about characters in the time format, see Table 3. This parameter can be specified when data in the CSV or JSON format is imported.
	TimestampFormat string `json:"timestamp_format,omitempty"`
	// Name of the queue that is specified to execute a task. If no queue is specified, the default queue is used.
	QueueName string `json:"queue_name,omitempty"`
	// Whether to overwrite data. The default value is false, indicating appending write. If the value is true, it indicates overwriting.
	Overwrite *bool `json:"overwrite,omitempty"`
	// Partition to which data is to be imported.
	// If this parameter is not set, the entire table data is dynamically imported. The imported data must contain the data in the partition column.
	// If this parameter is set and all partition information is configured during data import, data is imported to the specified partition. The imported data cannot contain data in the partition column.
	// If not all partition information is configured during data import, the imported data must contain all non-specified partition data. Otherwise, abnormal values such as null exist in the partition field column of non-specified data after data import.
	PartitionSpec map[string]string `json:"partition_spec,omitempty"`
	// User-defined parameter that applies to the job. Currently, dli.sql.dynamicPartitionOverwrite.enabled can be set to false by default. If it is set to true, data in a specified partition is overwritten. If it is set to false, data in the entire DataSource table is dynamically overwritten.
	// NOTE:
	// For dynamic overwrite of Hive partition tables, only the involved partition data can be overwritten. The entire table data cannot be overwritten.
	Conf []string `json:"conf,omitempty"`
}

type ExportDataOpts struct {
	// Path for storing the exported data. Currently, data can be stored only on OBS.
	// If export_mode is set to errorifexists, the OBS path cannot contain the specified folder,
	// for example, the test folder in the example request.
	DataPath string `json:"data_path" required:"true"`
	// Type of data to be exported. Currently, only CSV and JSON are supported.
	DataType string `json:"data_type" required:"true"`
	// Name of the database where the table from which data is exported resides.
	DatabaseName string `json:"database_name" required:"true"`
	// Name of the table from which data is exported.
	TableName string `json:"table_name" required:"true"`
	// Compression mode for exported data. Currently, the compression modes gzip, bzip2, and deflate are supported. If you do not want to compress data, enter none.
	Compress string `json:"compress" required:"true"`
	// Name of the queue that is specified to execute a task. If no queue is specified, the default queue is used.
	QueueName string `json:"queue_name,omitempty"`
	// Export mode. The parameter value can be ErrorIfExists or Overwrite. If export_mode is not specified, this parameter is set to ErrorIfExists by default.
	// ErrorIfExists: Ensure that the specified export directory does not exist. If the specified export directory exists, an error is reported and the export operation cannot be performed.
	// Overwrite: If you add new files to a specific directory, existing files will be deleted.
	ExportMode string `json:"export_mode,omitempty"`
	// Whether to export column names when exporting CSV and JSON data.
	// If this parameter is set to true, the column names are exported.
	// If this parameter is set to false, the column names are not exported.
	// If this parameter is left blank, the default value false is used.
	WithColumnHeader *bool `json:"with_column_header,omitempty"`
}

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Submit Job
func Submit(c *golangsdk.ServiceClient, opts SqlJobOpts) (*SubmitJobResult, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst SubmitJobResult
	_, err = c.Post(submitURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Cancel(c *golangsdk.ServiceClient, jobId string) (*CommonResp, error) {
	var rst CommonResp
	_, err := c.DeleteWithResponse(resourceURL(c, jobId), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func List(c *golangsdk.ServiceClient, opts ListJobsOpts) (*ListJobsResp, error) {
	url := listURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst ListJobsResp
	_, err = c.Get(url, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func Status(c *golangsdk.ServiceClient, jobId string) (*JobStatus, error) {
	var rst JobStatus
	_, err := c.Get(queryStatusURL(c, jobId), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Get(c *golangsdk.ServiceClient, jobId string) (*JobDetail, error) {
	var rst JobDetail
	_, err := c.Get(detailURL(c, jobId), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func CheckSQLGramar(c *golangsdk.ServiceClient, opts CheckSQLGramarOpts) (*CheckSqlResult, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst CheckSqlResult
	_, err = c.Post(checkSqlURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func ExportQueryResult(c *golangsdk.ServiceClient, jobId string, opts ExportQueryResultOpts) (*JobResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst JobResp
	_, err = c.Post(exportResultURL(c, jobId), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func Progress(c *golangsdk.ServiceClient, jobId string) (*JobProgress, error) {
	var rst JobProgress
	_, err := c.Get(progressURL(c, jobId), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func ImportData(c *golangsdk.ServiceClient, opts ImportDataOpts) (*JobResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst JobResp
	_, err = c.Post(importTableURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func ExportData(c *golangsdk.ServiceClient, opts ExportDataOpts) (*JobResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst JobResp
	_, err = c.Post(exportTableURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}
