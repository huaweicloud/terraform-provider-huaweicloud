package dataset

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type CreateOpts struct {
	DataFormat           string                `json:"data_format,omitempty"`
	DataSources          []DataSource          `json:"data_sources,omitempty"`
	DatasetName          string                `json:"dataset_name" required:"true"`
	DatasetType          int                   `json:"dataset_type,omitempty"`
	Description          string                `json:"description,omitempty"`
	ImportAnnotations    *bool                 `json:"import_annotations,omitempty"`
	ImportData           bool                  `json:"import_data,omitempty"`
	LabelFormat          LabelFormat           `json:"label_format,omitempty"`
	Labels               []Label               `json:"labels,omitempty"`
	Managed              bool                  `json:"managed,omitempty"`
	Schema               []Field               `json:"schema,omitempty"`
	WorkPath             string                `json:"work_path" required:"true"`
	WorkPathType         int                   `json:"work_path_type"`
	WorkforceInformation *WorkforceInformation `json:"workforce_information,omitempty"`
	WorkspaceId          string                `json:"workspace_id,omitempty"`
}

type DataSource struct {
	DataPath         string      `json:"data_path,omitempty"`
	DataType         int         `json:"data_type,omitempty"`
	SchemaMaps       []SchemaMap `json:"schema_maps,omitempty"`
	SourceInfo       SourceInfo  `json:"source_info,omitempty"`
	WithColumnHeader *bool       `json:"with_column_header,omitempty"`
}

type SchemaMap struct {
	DestName string `json:"dest_name,omitempty"`
	SrcName  string `json:"src_name,omitempty"`
}

type SourceInfo struct {
	ClusterId    string `json:"cluster_id,omitempty"`
	ClusterMode  string `json:"cluster_mode,omitempty"`
	ClusterName  string `json:"cluster_name,omitempty"`
	DatabaseName string `json:"database_name,omitempty"`
	Input        string `json:"input,omitempty"`
	Ip           string `json:"ip,omitempty"`
	Port         string `json:"port,omitempty"`
	QueueName    string `json:"queue_name,omitempty"`
	SubnetId     string `json:"subnet_id,omitempty"`
	TableName    string `json:"table_name,omitempty"`
	UserName     string `json:"user_name,omitempty"`
	UserPassword string `json:"user_password,omitempty"`
	VpcId        string `json:"vpc_id,omitempty"`
}

type LabelFormat struct {
	LabelType           string `json:"label_type,omitempty"`
	TextLabelSeparator  string `json:"text_label_separator,omitempty"`
	TextSampleSeparator string `json:"text_sample_separator,omitempty"`
}

type Label struct {
	Attributes []LabelAttribute `json:"attributes,omitempty"`
	Name       string           `json:"name,omitempty"`
	Property   LabelProperty    `json:"property,omitempty"`
	Type       *int             `json:"type,omitempty"`
}

type LabelAttribute struct {
	DefaultValue string              `json:"default_value,omitempty"`
	Id           string              `json:"id,omitempty"`
	Name         string              `json:"name,omitempty"`
	Type         string              `json:"type,omitempty"`
	Values       LabelAttributeValue `json:"values,omitempty"`
}

type LabelAttributeValue struct {
	Id    string `json:"id,omitempty"`
	Value string `json:"value,omitempty"`
}

type LabelProperty struct {
	Color        string `json:"@modelarts:color,omitempty"`
	DefaultShape string `json:"@modelarts:default_shape,omitempty"`
	FromType     string `json:"@modelarts:from_type,omitempty"`
	RenameTo     string `json:"@modelarts:rename_to,omitempty"`
	Shortcut     string `json:"@modelarts:shortcut,omitempty"`
	ToType       string `json:"@modelarts:to_type,omitempty"`
}

type Field struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	SchemaId    int    `json:"schema_id,omitempty"`
	Type        string `json:"type,omitempty"`
}

type WorkforceInformation struct {
	DataSyncType                *int             `json:"data_sync_type,omitempty"`
	Repetition                  int              `json:"repetition,omitempty"`
	SynchronizeAutoLabelingData *bool            `json:"synchronize_auto_labeling_data,omitempty"`
	SynchronizeData             *bool            `json:"synchronize_data,omitempty"`
	TaskId                      string           `json:"task_id,omitempty"`
	TaskName                    string           `json:"task_name" required:"true"`
	WorkforcesConfig            WorkforcesConfig `json:"workforces_config,omitempty"`
}

type WorkforcesConfig struct {
	Agency     string            `json:"agency,omitempty"`
	Workforces []WorkforceConfig `json:"workforces,omitempty"`
}

type WorkforceConfig struct {
	Workers       []Worker `json:"workers,omitempty"`
	WorkforceId   string   `json:"workforce_id,omitempty"`
	WorkforceName string   `json:"workforce_name,omitempty"`
}

type Worker struct {
	CreateTime  *int   `json:"create_time,omitempty"`
	Description string `json:"description,omitempty"`
	Email       string `json:"email,omitempty"`
	Role        *int   `json:"role,omitempty"`
	Status      *int   `json:"status,omitempty"`
	UpdateTime  *int   `json:"update_time,omitempty"`
	WorkerId    string `json:"worker_id,omitempty"`
	WorkforceId string `json:"workforce_id,omitempty"`
}

type UpdateOpts struct {
	AddLabels        []Label `json:"add_labels,omitempty"`
	CurrentVersionId string  `json:"current_version_id,omitempty"`
	DatasetName      string  `json:"dataset_name,omitempty"`
	DeleteLabels     []Label `json:"delete_labels,omitempty"`
	Description      *string `json:"description,omitempty"`
	UpdateLabels     []Label `json:"update_labels,omitempty"`
}

type ListOpts struct {
	CheckRunningTask   bool   `q:"check_running_task,omitempty"`
	ContainVersions    *bool  `q:"contain_versions,omitempty"`
	DatasetType        *int   `q:"dataset_type,omitempty"`
	FilePreview        bool   `q:"file_preview,omitempty"`
	Limit              *int   `q:"limit,omitempty"`
	Offset             int    `q:"offset,omitempty"`
	Order              string `q:"order,omitempty"`
	RunningTaskType    *int   `q:"running_task_type,omitempty"`
	SearchContent      string `q:"search_content,omitempty"`
	SortBy             string `q:"sort_by,omitempty"`
	SupportExport      bool   `q:"support_export,omitempty"`
	TrainEvaluateRatio string `q:"train_evaluate_ratio,omitempty"`
	VersionFormat      *int   `q:"version_format,omitempty"`
	WithLabels         bool   `q:"with_labels,omitempty"`
	WorkspaceId        string `q:"workspace_id,omitempty"`
}

type GetOpts struct {
	CheckRunningTask bool `q:"check_running_task,omitempty"`
	RunningTaskType  *int `q:"running_task_type,omitempty"`
}

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*CreateResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst CreateResp
	_, err = c.Post(createURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func Get(c *golangsdk.ServiceClient, id string, opts GetOpts) (*Dataset, error) {
	var rst Dataset

	url := getURL(c, id)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	_, err = c.Get(url, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Update(c *golangsdk.ServiceClient, id string, opts UpdateOpts) *golangsdk.ErrResult {
	var rst golangsdk.ErrResult
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		rst.Err = err
		return &rst
	}

	_, rst.Err = c.Put(updateURL(c, id), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst
}

func Delete(c *golangsdk.ServiceClient, id string) *golangsdk.ErrResult {
	url := deleteURL(c, id)
	var rst golangsdk.ErrResult
	_, rst.Err = c.Delete(url, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst
}

func List(c *golangsdk.ServiceClient, opts ListOpts) (*pagination.Pager, error) {
	url := listURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	page := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := DatasetPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	})

	return &page, nil
}
