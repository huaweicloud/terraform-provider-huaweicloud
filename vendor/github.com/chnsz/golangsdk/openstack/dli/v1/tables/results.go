package tables

type CommonResp struct {
	IsSuccess bool   `json:"is_success"`
	Message   string `json:"message"`
}

type DeleteResp struct {
	CommonResp
	JobMode string `json:"job_mode"`
}

type ListResp struct {
	IsSuccess  bool         `json:"is_success"`
	Message    string       `json:"message"`
	TableCount int          `json:"table_count"`
	Tables     []Table4List `json:"tables"`
}

// If with-detail is set to false in the URI, only values of tables-related parameters data_location, table_name,
// and table_type are returned.
type Table4List struct {
	CreateTime       int    `json:"create_time"`
	DataType         string `json:"data_type"`
	DataLocation     string `json:"data_location"`
	LastAccessTime   int    `json:"last_access_time"`
	Location         string `json:"location"`
	Owner            string `json:"owner"`
	TableName        string `json:"table_name"`
	TableSize        int    `json:"table_size"`
	PartitionColumns string `json:"partition_columns"`
	PageSize         int    `json:"page-size"`
	CurrentPage      int    `json:"current-page"`
	// Type of a table.
	// EXTERNAL: Indicates an OBS table.
	// MANAGED: Indicates a DLI table.
	// VIEW: Indicates a view.
	TableType string `json:"table_type"`
}

type Table struct {
	IsSuccess         bool                     `json:"is_success"`
	Message           string                   `json:"message"`
	ColumnCount       int                      `json:"column_count"`
	Columns           []Column                 `json:"columns"`
	TableType         string                   `json:"table_type"`
	DataType          string                   `json:"data_type"`
	DataLocation      string                   `json:"data_location"`
	StorageProperties []map[string]interface{} `json:"storage_properties"`
	TableComment      string                   `json:"table_comment"`
	CreateTableSql    string                   `json:"create_table_sql"`
}

type Column struct {
	ColumnName        string `json:"column_name"`
	Type              string `json:"type"`
	Description       string `json:"description"`
	IsPartitionColumn bool   `json:"is_partition_column"`
}

type PartitionsResp struct {
	IsSuccess  bool           `json:"is_success"`
	Message    string         `json:"message"`
	Partitions PartitionsInfo `json:"partitions"`
}

type PartitionsInfo struct {
	TotalCount     int         `json:"total_count"`
	PartitionInfos []Partition `json:"partition_infos"`
}

type Partition struct {
	PartitionName  string   `json:"partition_name"`
	CreateTime     int      `json:"create_time"`
	LastAccessTime int      `json:"last_access_time"`
	Locations      []string `json:"locations"`
	LastDdlTime    int      `json:"last_ddl_time"`
	NumRows        int      `json:"num_rows"`
	NumFiles       int      `json:"num_files"`
	TotalSize      int      `json:"total_size"`
}
