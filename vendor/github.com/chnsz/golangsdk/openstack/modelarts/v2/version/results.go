package version

import (
	"github.com/chnsz/golangsdk/openstack/modelarts/v2/dataset"
	"github.com/chnsz/golangsdk/pagination"
)

type CreateResp struct {
	VersionId string `json:"version_id"`
}

type DatasetVersion struct {
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
	LabelStats                   []dataset.Label        `json:"label_stats"`
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

type ListDatasetVersions struct {
	TotalNumber int              `json:"total_number"`
	Versions    []DatasetVersion `json:"versions"`
}

type DatasetVersionPage struct {
	pagination.OffsetPageBase
}

func (b DatasetVersionPage) IsEmpty() (bool, error) {
	arr, err := ExtractDatasetVersions(b)
	return len(arr) == 0, err
}

func ExtractDatasetVersions(r pagination.Page) ([]DatasetVersion, error) {
	var s ListDatasetVersions
	err := (r.(DatasetVersionPage)).ExtractInto(&s)
	return s.Versions, err
}
