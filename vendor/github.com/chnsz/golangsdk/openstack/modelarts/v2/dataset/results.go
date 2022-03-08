package dataset

import "github.com/chnsz/golangsdk/pagination"

const (
	StatusInit         = "INIT"
	StatusCreating     = "CREATING"
	StatusStarting     = "STARTING"
	StatusStopping     = "STOPPING"
	StatusDeleting     = "DELETING"
	StatusRunning      = "RUNNING"
	StatusStopped      = "STOPPED"
	StatusSnapshotting = "SNAPSHOTTING"
	StatusCreateFailed = "CREATE_FAILED"
	StatusStartFailed  = "START_FAILED"
	StatusDeleteFailed = "DELETE_FAILED"
	StatusError        = "ERROR"
	StatusDeleted      = "DELETED"
	StatusFrozen       = "FROZEN"
)

type CreateResp struct {
	DatasetId    string `json:"dataset_id"`
	ErrorCode    string `json:"error_code"`
	ErrorMsg     string `json:"error_msg"`
	ImportTaskId string `json:"import_task_id"`
}

type Dataset struct {
	AnnotatedSampleCount    int                    `json:"annotated_sample_count"`
	AnnotatedSubSampleCount int                    `json:"annotated_sub_sample_count"`
	ContentLabeling         bool                   `json:"content_labeling"`
	CreateTime              int                    `json:"create_time"`
	CurrentVersionId        string                 `json:"current_version_id"`
	CurrentVersionName      string                 `json:"current_version_name"`
	DataFormat              string                 `json:"data_format"`
	DataSources             []DataSource           `json:"data_sources"`
	DataStatistics          map[string]interface{} `json:"data_statistics"`
	DataUpdateTime          int                    `json:"data_update_time"`
	DatasetFormat           int                    `json:"dataset_format"`
	DatasetId               string                 `json:"dataset_id"`
	DatasetName             string                 `json:"dataset_name"`
	DatasetTags             []string               `json:"dataset_tags"`
	DatasetType             int                    `json:"dataset_type"`
	DatasetVersionCount     int                    `json:"dataset_version_count"`
	DeletedSampleCount      int                    `json:"deleted_sample_count"`
	DeletionStats           map[string]int         `json:"deletion_stats"`
	Description             string                 `json:"description"`
	EnterpriseProjectId     string                 `json:"enterprise_project_id"`
	ExistRunningTask        bool                   `json:"exist_running_task"`
	ExistWorkforceTask      bool                   `json:"exist_workforce_task"`
	FeatureSupports         []string               `json:"feature_supports"`
	ImportData              bool                   `json:"import_data"`
	ImportTaskId            string                 `json:"import_task_id"`
	InnerAnnotationPath     string                 `json:"inner_annotation_path"`
	InnerDataPath           string                 `json:"inner_data_path"`
	InnerLogPath            string                 `json:"inner_log_path"`
	InnerTaskPath           string                 `json:"inner_task_path"`
	InnerTempPath           string                 `json:"inner_temp_path"`
	InnerWorkPath           string                 `json:"inner_work_path"`
	LabelTaskCount          int                    `json:"label_task_count"`
	Labels                  []Label                `json:"labels"`
	LoadingSampleCount      int                    `json:"loading_sample_count"`
	Managed                 bool                   `json:"managed"`
	NextVersionNum          int                    `json:"next_version_num"`
	RunningTasksId          []string               `json:"running_tasks_id"`
	Schema                  []Field                `json:"schema"`
	Status                  int                    `json:"status"`
	ThirdPath               string                 `json:"third_path"`
	TotalSampleCount        int                    `json:"total_sample_count"`
	TotalSubSampleCount     int                    `json:"total_sub_sample_count"`
	UnconfirmedSampleCount  int                    `json:"unconfirmed_sample_count"`
	UpdateTime              int                    `json:"update_time"`
	Versions                []Version              `json:"versions"`
	WorkPath                string                 `json:"work_path"`
	WorkPathType            int                    `json:"work_path_type"`
	WorkforceDescriptor     []WorkforceDescriptor  `json:"workforce_descriptor"`
	WorkforceTaskCount      int                    `json:"workforce_task_count"`
	WorkspaceId             string                 `json:"workspace_id"`
}

type Version struct {
	AddSampleCount               int                    `json:"add_sample_count"`
	AnalysisCachePath            string                 `json:"analysis_cache_path"`
	AnalysisStatus               int                    `json:"analysis_status"`
	AnalysisTaskId               string                 `json:"analysis_task_id"`
	AnnotatedSampleCount         int                    `json:"annotated_sample_count"`
	AnnotatedSubSampleCount      int                    `json:"annotated_sub_sample_count"`
	ClearHardProperty            bool                   `json:"clear_hard_property"`
	Code                         string                 `json:"code"`
	CreateTime                   int                    `json:"create_time"`
	Crop                         bool                   `json:"crop"`
	CropPath                     string                 `json:"crop_path"`
	CropRotateCachePath          string                 `json:"crop_rotate_cache_path"`
	DataAnalysis                 map[string]interface{} `json:"data_analysis"`
	DataPath                     string                 `json:"data_path"`
	DataStatistics               map[string]interface{} `json:"data_statistics"`
	DataValidate                 bool                   `json:"data_validate"`
	DeletedSampleCount           int                    `json:"deleted_sample_count"`
	DeletionStats                map[string]int         `json:"deletion_stats"`
	Description                  string                 `json:"description"`
	ExportImages                 bool                   `json:"export_images"`
	ExtractSerialNumber          bool                   `json:"extract_serial_number"`
	IncludeDatasetData           bool                   `json:"include_dataset_data"`
	IsCurrent                    bool                   `json:"is_current"`
	LabelStats                   []LabelStats           `json:"label_stats"`
	LabelType                    string                 `json:"label_type"`
	ManifestCacheInputPath       string                 `json:"manifest_cache_input_path"`
	ManifestPath                 string                 `json:"manifest_path"`
	Message                      string                 `json:"message"`
	ModifiedSampleCount          int                    `json:"modified_sample_count"`
	PreviousAnnotatedSampleCount int                    `json:"previous_annotated_sample_count"`
	PreviousTotalSampleCount     int                    `json:"previous_total_sample_count"`
	PreviousVersionId            string                 `json:"previous_version_id"`
	ProcessorTaskId              string                 `json:"processor_task_id"`
	ProcessorTaskStatus          int                    `json:"processor_task_status"`
	RemoveSampleUsage            bool                   `json:"remove_sample_usage"`
	Rotate                       bool                   `json:"rotate"`
	RotatePath                   string                 `json:"rotate_path"`
	SampleState                  string                 `json:"sample_state"`
	StartProcessorTask           bool                   `json:"start_processor_task"`
	Status                       int                    `json:"status"`
	Tags                         []string               `json:"tags"`
	TaskType                     int                    `json:"task_type"`
	TotalSampleCount             int                    `json:"total_sample_count"`
	TotalSubSampleCount          int                    `json:"total_sub_sample_count"`
	TrainEvaluateSampleRatio     string                 `json:"train_evaluate_sample_ratio"`
	UpdateTime                   int                    `json:"update_time"`
	VersionFormat                string                 `json:"version_format"`
	VersionId                    string                 `json:"version_id"`
	VersionName                  string                 `json:"version_name"`
	WithColumnHeader             bool                   `json:"with_column_header"`
}

type LabelStats struct {
	Attributes  []LabelAttribute `json:"attributes"`
	Count       int              `json:"count"`
	Name        string           `json:"name"`
	Property    LabelProperty    `json:"property"`
	SampleCount int              `json:"sample_count"`
	Type        int              `json:"type"`
}

type WorkforceDescriptor struct {
	CurrentTaskId                 string   `json:"current_task_id"`
	CurrentTaskName               string   `json:"current_task_name"`
	RejectNum                     int      `json:"reject_num"`
	Repetition                    int      `json:"repetition"`
	IsSynchronizeAutoLabelingData bool     `json:"is_synchronize_auto_labeling_data"`
	IsSynchronizeData             bool     `json:"is_synchronize_data"`
	Workers                       []Worker `json:"workers"`
	WorkforceId                   string   `json:"workforce_id"`
	WorkforceName                 string   `json:"workforce_name"`
}

type ListResp struct {
	Result []Dataset4List `json:"datasets"`
	Total  int            `json:"total_number"`
}

type Dataset4List struct {
	Dataset
	DataUrl string           `json:"data_url"`
	Samples []AnnotationFile `json:"samples"`
}

type AnnotationFile struct {
	CreateTime int               `json:"create_time"`
	DatasetId  string            `json:"dataset_id"`
	Depth      int               `json:"depth"`
	FileName   string            `json:"file_Name"`
	FileId     string            `json:"file_id"`
	FileType   string            `json:"file_type"`
	Height     int               `json:"height"`
	Size       int               `json:"size"`
	Tags       map[string]string `json:"tags"`
	Url        string            `json:"url"`
	Width      int               `json:"width"`
}

type DatasetPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a RouteTablePage struct is empty.
func (b DatasetPage) IsEmpty() (bool, error) {
	arr, err := ExtractDatasets(b)
	return len(arr) == 0, err
}

func ExtractDatasets(r pagination.Page) ([]Dataset4List, error) {
	var s ListResp
	err := (r.(DatasetPage)).ExtractInto(&s)
	return s.Result, err
}
