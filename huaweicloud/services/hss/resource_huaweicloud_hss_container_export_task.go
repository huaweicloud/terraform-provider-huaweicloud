package hss

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var containerExportTaskNonUpdatableParams = []string{
	"export_headers",
	"cluster_container",
	"cluster_type",
	"cluster_name",
	"container_name",
	"pod_name",
	"image_name",
	"status",
	"risky",
	"create_time",
	"create_time.*.start_time",
	"create_time.*.end_time",
	"cpu_limit",
	"memory_limit",
	"enterprise_project_id",
	"export_size",
}

// @API HSS POST /v5/{project_id}/container/export-task
// @API HSS GET /v5/{project_id}/export-task/{task_id}
func ResourceContainerExportTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContainerExportTaskCreate,
		ReadContext:   resourceContainerExportTaskRead,
		UpdateContext: resourceContainerExportTaskUpdate,
		DeleteContext: resourceContainerExportTaskDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(containerExportTaskNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Specifies the region where the resource is located.",
			},
			// Body Params
			"export_headers": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Specifies the headers for the exported container list.",
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			// This field is defined as a Boolean value in the API. Since this field has no special meaning, its type is
			// adjusted to a string.
			"cluster_container": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether the container is in a cluster.",
			},
			"cluster_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the cluster type.",
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the cluster to which the container belongs.",
			},
			"container_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the container to export.",
			},
			"pod_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the Pod to which the container belongs.",
			},
			"image_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the container image.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the container status.",
			},
			// This field is defined as a Boolean value in the API. Since this field has no special meaning, its type is
			// adjusted to a string.
			"risky": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether the container has security risks.",
			},
			"create_time": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Specifies the time range for filtering containers.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the start time for filtering containers.",
						},
						"end_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the end time for filtering containers.",
						},
					},
				},
			},
			"cpu_limit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the CPU limit for the container.",
			},
			"memory_limit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the memory limit for the container.",
			},
			// Query Params
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the enterprise project to which the resource belongs.",
			},
			"export_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the number of containers to export.",
			},
			// Computed Attributes
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the export task.",
			},
			"task_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the export task.",
			},
			"task_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the export task.",
			},
			"file_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the exported file.",
			},
			"file_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the exported file.",
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildCreateContainerExportTaskQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := ""
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("export_size"); ok {
		queryParams = fmt.Sprintf("%s&export_size=%v", queryParams, v)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func buildContainerExportTaskClusterContainerParam(d *schema.ResourceData) interface{} {
	rawValue := d.Get("cluster_container").(string)
	if rawValue == "" {
		return nil
	}

	return utils.StringToBool(rawValue)
}

func buildContainerExportTaskCreateTimeBodyParam(d *schema.ResourceData) interface{} {
	rawArray := d.Get("create_time").([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"start_time": rawMap["start_time"],
		"end_time":   rawMap["end_time"],
	}
}

func buildContainerExportTaskExportHeadersBodyParam(d *schema.ResourceData) interface{} {
	rst := make([]interface{}, 0)
	for _, v := range d.Get("export_headers").([]interface{}) {
		row := utils.ExpandToStringList(v.([]interface{}))
		rst = append(rst, row)
	}
	return rst
}

func buildContainerExportTaskRiskyParam(d *schema.ResourceData) interface{} {
	rawValue := d.Get("risky").(string)
	if rawValue == "" {
		return nil
	}

	return utils.StringToBool(rawValue)
}

func buildCreateContainerExportTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"export_headers":    buildContainerExportTaskExportHeadersBodyParam(d),
		"cluster_container": buildContainerExportTaskClusterContainerParam(d),
		"cluster_type":      utils.ValueIgnoreEmpty(d.Get("cluster_type")),
		"cluster_name":      utils.ValueIgnoreEmpty(d.Get("cluster_name")),
		"container_name":    utils.ValueIgnoreEmpty(d.Get("container_name")),
		"pod_name":          utils.ValueIgnoreEmpty(d.Get("pod_name")),
		"image_name":        utils.ValueIgnoreEmpty(d.Get("image_name")),
		"status":            utils.ValueIgnoreEmpty(d.Get("status")),
		"risky":             buildContainerExportTaskRiskyParam(d),
		"create_time":       buildContainerExportTaskCreateTimeBodyParam(d),
		"cpu_limit":         utils.ValueIgnoreEmpty(d.Get("cpu_limit")),
		"memory_limit":      utils.ValueIgnoreEmpty(d.Get("memory_limit")),
	}
}

func resourceContainerExportTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/container/export-task"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildCreateContainerExportTaskQueryParams(d, epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateContainerExportTaskBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error exporting HSS container task: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskId := utils.PathSearch("task_id", respBody, "").(string)
	if taskId == "" {
		return diag.Errorf("error exporting HSS container task: task ID is not found in API response")
	}

	d.SetId(taskId)

	queryTaskRespBody, err := waitingForExportTaskSuccess(ctx, client, d.Timeout(schema.TimeoutCreate), taskId)
	if err != nil {
		return diag.Errorf("error waiting for the HSS container export task to complete: %s", err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("task_id", taskId),
		d.Set("task_name", utils.PathSearch("task_name", queryTaskRespBody, nil)),
		d.Set("task_status", utils.PathSearch("task_status", queryTaskRespBody, nil)),
		d.Set("file_id", utils.PathSearch("file_id", queryTaskRespBody, nil)),
		d.Set("file_name", utils.PathSearch("file_name", queryTaskRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceContainerExportTaskRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceContainerExportTaskUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceContainerExportTaskDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to export HSS container task. Deleting this resource
	will not clear the corresponding exported record, but will only remove the resource information from the tf state
	file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
