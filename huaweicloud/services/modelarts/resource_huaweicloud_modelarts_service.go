// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ModelArts
// ---------------------------------------------------------------

package modelarts

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts POST /v1/{project_id}/services
// @API ModelArts DELETE /v1/{project_id}/services/{id}
// @API ModelArts GET /v1/{project_id}/services/{id}
// @API ModelArts PUT /v1/{project_id}/services/{id}
func ResourceModelartsService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceModelartsServiceCreate,
		UpdateContext: resourceModelartsServiceUpdate,
		ReadContext:   resourceModelartsServiceRead,
		DeleteContext: resourceModelartsServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Service name, which consists of 1 to 64 characters.`,
			},
			"infer_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Inference mode.`,
			},
			"config": {
				Type:     schema.TypeList,
				Elem:     modelartsServiceConfigSchema(),
				Required: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `ID of the workspace to which a service belongs.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the service.`,
			},
			"pool_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the new dedicated resource pool.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The VPC ID to which a real-time service instance is deployed.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The subnet ID.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The security group ID.`,
			},
			"schedule": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     modelartsServiceScheduleSchema(),
				Optional: true,
				Computed: true,
			},
			"additional_properties": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     modelartsServiceAdditionalPropertySchema(),
				Optional: true,
				Computed: true,
			},
			"change_status_to": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Which status you want to change the service to, the valid value can be **running** or **stopped**.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `User to which a service belongs`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Service status.`,
			},
			"access_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Access address of an inference request.`,
			},
			"bind_access_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Request address of a custom domain name.`,
			},
			"invocation_times": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Total number of service calls.`,
			},
			"failed_times": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Number of failed service calls.`,
			},
			"is_shared": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether a service is subscribed.`,
			},
			"shared_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Number of subscribed services.`,
			},
			"debug_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Online debugging address of a real-time service.`,
			},
			"is_free": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether a free-of-charge flavor is used.`,
			},
		},
	}
}

func modelartsServiceScheduleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"duration": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Value mapping a time unit.`,
			},
			"time_unit": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Scheduling time unit. Possible values are DAYS, HOURS, and MINUTES.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Scheduling type. Only the value **stop** is supported.`,
			},
		},
	}
	return &sc
}

func modelartsServiceConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"custom_spec": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     modelartsServiceConfigCustomSpecSchema(),
				Optional: true,
				Computed: true,
			},
			"envs": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Environment variable key-value pair required for running a model.`,
			},
			"specification": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Resource flavors.`,
			},
			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Weight of traffic allocated to a model.`,
			},
			"model_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Model ID, which can be obtained by calling the API for obtaining a model list.`,
			},
			"src_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `OBS path to the input data of a batch job.`,
			},
			"req_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Inference API called in a batch task, which is the RESTful API exposed in the model image.`,
			},
			"mapping_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Mapping type of the input data. Mandatory for batch services.`,
			},
			"pool_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `ID of a dedicated resource pool. Optional for real-time services.`,
			},
			"nodes": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Edge node ID array. Mandatory for edge services.`,
			},
			"mapping_rule": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Mapping between input parameters and CSV data. Optional for batch services.`,
			},
			"src_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Data source type, which can be ManifestFile. Mandatory for batch services.`,
			},
			"dest_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `OBS path to the output data of a batch job. Mandatory for batch services.`,
			},
			"instance_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Number of instances deployed for a model.`,
			},
		},
	}
	return &sc
}

func modelartsServiceConfigCustomSpecSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"memory": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Memory in MB, which must be an integer.`,
			},
			"cpu": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Computed:    true,
				Description: `Number of CPU cores, which can be a decimal. The value cannot be smaller than 0.01.`,
			},
			"gpu_p4": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Computed:    true,
				Description: `Number of GPU cores, which can be a decimal.`,
			},
			"ascend_a310": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Number of Ascend chips. Either this parameter or **gpu_p4** is configured.`,
			},
		},
	}
	return &sc
}

func modelartsServiceAdditionalPropertySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"smn_notification": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     modelartsServiceAdditionalPropertySmnNotificationSchema(),
				Optional: true,
				Computed: true,
			},
			"log_report_channels": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     modelartsServiceAdditionalPropertyLogReportChannelSchema(),
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func modelartsServiceAdditionalPropertySmnNotificationSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"topic_urn": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `URN of an SMN topic.`,
			},
			"events": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
				Computed:    true,
				Description: `Event ID.`,
			},
		},
	}
	return &sc
}

func modelartsServiceAdditionalPropertyLogReportChannelSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The type of log report channel. The valid value is **LTS**.`,
			},
		},
	}
	return &sc
}

func resourceModelartsServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createService: create a ModelArts service.
	var (
		createServiceHttpUrl = "v1/{project_id}/services"
		createServiceProduct = "modelarts"
	)
	createServiceClient, err := cfg.NewServiceClient(createServiceProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts Client: %s", err)
	}

	createServicePath := createServiceClient.Endpoint + createServiceHttpUrl
	createServicePath = strings.ReplaceAll(createServicePath, "{project_id}", createServiceClient.ProjectID)

	createServiceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	createServiceOpt.JSONBody = utils.RemoveNil(buildCreateServiceBodyParams(d))
	createServiceResp, err := createServiceClient.Request("POST", createServicePath, &createServiceOpt)
	if err != nil {
		return diag.Errorf("error creating ModelartsService: %s", err)
	}

	createServiceRespBody, err := utils.FlattenResponse(createServiceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	serviceId := utils.PathSearch("service_id", createServiceRespBody, "").(string)
	if serviceId == "" {
		return diag.Errorf("unable to find the ModelArts service ID from the API response")
	}
	d.SetId(serviceId)

	err = serviceWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Create of ModelartsService (%s) to complete: %s", d.Id(), err)
	}
	return resourceModelartsServiceRead(ctx, d, meta)
}

func buildCreateServiceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"workspace_id":          utils.ValueIgnoreEmpty(d.Get("workspace_id")),
		"service_name":          d.Get("name"),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"infer_type":            d.Get("infer_type"),
		"pool_name":             utils.ValueIgnoreEmpty(d.Get("pool_name")),
		"vpc_id":                utils.ValueIgnoreEmpty(d.Get("vpc_id")),
		"subnet_network_id":     utils.ValueIgnoreEmpty(d.Get("subnet_id")),
		"security_group_id":     utils.ValueIgnoreEmpty(d.Get("security_group_id")),
		"schedule":              buildServiceRequestBodySchedule(d.Get("schedule")),
		"config":                buildServiceRequestBodyConfig(d.Get("config")),
		"additional_properties": buildServiceRequestBodyAdditionalProperty(d.Get("additional_properties")),
	}
	return bodyParams
}

func buildServiceRequestBodySchedule(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"duration":  utils.ValueIgnoreEmpty(raw["duration"]),
				"time_unit": utils.ValueIgnoreEmpty(raw["time_unit"]),
				"type":      utils.ValueIgnoreEmpty(raw["type"]),
			}
		}
		return rst
	}
	return nil
}

func buildServiceRequestBodyConfig(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"custom_spec":    buildConfigCustomSpec(raw["custom_spec"]),
				"envs":           utils.ValueIgnoreEmpty(raw["envs"]),
				"specification":  utils.ValueIgnoreEmpty(raw["specification"]),
				"weight":         utils.ValueIgnoreEmpty(raw["weight"]),
				"model_id":       utils.ValueIgnoreEmpty(raw["model_id"]),
				"src_path":       utils.ValueIgnoreEmpty(raw["src_path"]),
				"req_uri":        utils.ValueIgnoreEmpty(raw["req_uri"]),
				"mapping_type":   utils.ValueIgnoreEmpty(raw["mapping_type"]),
				"pool_name":      utils.ValueIgnoreEmpty(raw["pool_name"]),
				"nodes":          utils.ValueIgnoreEmpty(raw["nodes"]),
				"mapping_rule":   utils.ValueIgnoreEmpty(raw["mapping_rule"]),
				"src_type":       utils.ValueIgnoreEmpty(raw["src_type"]),
				"dest_path":      utils.ValueIgnoreEmpty(raw["dest_path"]),
				"instance_count": utils.ValueIgnoreEmpty(raw["instance_count"]),
			}
		}
		return rst
	}
	return nil
}

func buildConfigCustomSpec(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"memory":      utils.ValueIgnoreEmpty(raw["memory"]),
			"cpu":         utils.ValueIgnoreEmpty(raw["cpu"]),
			"gpu_p4":      utils.ValueIgnoreEmpty(raw["gpu_p4"]),
			"ascend_a310": utils.ValueIgnoreEmpty(raw["ascend_a310"]),
		}
		return params
	}
	return nil
}

func buildServiceRequestBodyAdditionalProperty(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"smn_notification":    buildAdditionalPropertySmnNotification(raw["smn_notification"]),
			"log_report_channels": buildAdditionalPropertyLogReportChannel(raw["log_report_channels"]),
		}
		return params
	}
	return nil
}

func buildAdditionalPropertySmnNotification(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"topic_urn": utils.ValueIgnoreEmpty(raw["topic_urn"]),
			"events":    utils.ValueIgnoreEmpty(raw["events"]),
		}
		return params
	}
	return nil
}

func buildAdditionalPropertyLogReportChannel(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"type": utils.ValueIgnoreEmpty(raw["type"]),
			}
		}
		return rst
	}
	return nil
}

func serviceWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			var (
				createServiceWaitingHttpUrl = "v1/{project_id}/services/{id}"
				createServiceWaitingProduct = "modelarts"
			)
			createServiceWaitingClient, err := cfg.NewServiceClient(createServiceWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating ModelArts Client: %s", err)
			}

			createServiceWaitingPath := createServiceWaitingClient.Endpoint + createServiceWaitingHttpUrl
			createServiceWaitingPath = strings.ReplaceAll(createServiceWaitingPath, "{project_id}", createServiceWaitingClient.ProjectID)
			createServiceWaitingPath = strings.ReplaceAll(createServiceWaitingPath, "{id}", d.Id())

			createServiceWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
				MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
			}

			createServiceWaitingResp, err := createServiceWaitingClient.Request("GET", createServiceWaitingPath, &createServiceWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createServiceWaitingRespBody, err := utils.FlattenResponse(createServiceWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			status := utils.PathSearch(`status`, createServiceWaitingRespBody, "").(string)

			targetStatus := []string{
				"running",
				"finished",
				"stopped",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createServiceWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"failed",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return createServiceWaitingRespBody, status, nil
			}

			return createServiceWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceModelartsServiceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getService: Query the ModelArts service.
	var (
		getServiceHttpUrl = "v1/{project_id}/services/{id}"
		getServiceProduct = "modelarts"
	)
	getServiceClient, err := cfg.NewServiceClient(getServiceProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts Client: %s", err)
	}

	getServicePath := getServiceClient.Endpoint + getServiceHttpUrl
	getServicePath = strings.ReplaceAll(getServicePath, "{project_id}", getServiceClient.ProjectID)
	getServicePath = strings.ReplaceAll(getServicePath, "{id}", d.Id())

	getServiceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	getServiceResp, err := getServiceClient.Request("GET", getServicePath, &getServiceOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ModelartsService")
	}

	getServiceRespBody, err := utils.FlattenResponse(getServiceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("service_name", getServiceRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getServiceRespBody, nil)),
		d.Set("owner", utils.PathSearch("owner", getServiceRespBody, nil)),
		d.Set("infer_type", utils.PathSearch("infer_type", getServiceRespBody, nil)),
		d.Set("workspace_id", utils.PathSearch("workspace_id", getServiceRespBody, nil)),
		d.Set("pool_name", utils.PathSearch("pool_name", getServiceRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", getServiceRespBody, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_network_id", getServiceRespBody, nil)),
		d.Set("security_group_id", utils.PathSearch("security_group_id", getServiceRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getServiceRespBody, nil)),
		d.Set("config", flattenGetServiceResponseBodyConfig(getServiceRespBody)),
		d.Set("access_address", utils.PathSearch("access_address", getServiceRespBody, nil)),
		d.Set("bind_access_address", utils.PathSearch("bind_access_address", getServiceRespBody, nil)),
		d.Set("invocation_times", utils.PathSearch("invocation_times", getServiceRespBody, nil)),
		d.Set("failed_times", utils.PathSearch("failed_times", getServiceRespBody, nil)),
		d.Set("is_shared", utils.PathSearch("is_shared", getServiceRespBody, nil)),
		d.Set("shared_count", utils.PathSearch("shared_count", getServiceRespBody, nil)),
		d.Set("schedule", flattenGetServiceResponseBodySchedule(getServiceRespBody)),
		d.Set("debug_url", utils.PathSearch("debug_url", getServiceRespBody, nil)),
		d.Set("is_free", utils.PathSearch("is_free", getServiceRespBody, nil)),
		d.Set("additional_properties", flattenServiceResponseAdditionalProperty(getServiceRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetServiceResponseBodyConfig(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("config", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"custom_spec":    flattenConfigCustomSpec(v),
			"envs":           utils.PathSearch("envs", v, nil),
			"specification":  utils.PathSearch("specification", v, nil),
			"weight":         utils.PathSearch("weight", v, nil),
			"model_id":       utils.PathSearch("model_id", v, nil),
			"src_path":       utils.PathSearch("src_path", v, nil),
			"req_uri":        utils.PathSearch("req_uri", v, nil),
			"mapping_type":   utils.PathSearch("mapping_type", v, nil),
			"pool_name":      utils.PathSearch("pool_name", v, nil),
			"nodes":          utils.PathSearch("nodes", v, nil),
			"mapping_rule":   utils.PathSearch("mapping_rule", v, nil),
			"src_type":       utils.PathSearch("src_type", v, nil),
			"dest_path":      utils.PathSearch("dest_path", v, nil),
			"instance_count": utils.PathSearch("instance_count", v, nil),
		})
	}
	return rst
}

func flattenConfigCustomSpec(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("custom_spec", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"memory":      utils.PathSearch("memory", curJson, nil),
			"cpu":         utils.PathSearch("cpu", curJson, nil),
			"gpu_p4":      utils.PathSearch("gpu_p4", curJson, nil),
			"ascend_a310": utils.PathSearch("ascend_a310", curJson, nil),
		},
	}
	return rst
}

func flattenGetServiceResponseBodySchedule(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("schedule", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"duration":  utils.PathSearch("duration", v, nil),
			"time_unit": utils.PathSearch("time_unit", v, nil),
			"type":      utils.PathSearch("type", v, nil),
		})
	}
	return rst
}

func flattenServiceResponseAdditionalProperty(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("additional_properties", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"smn_notification":    flattenAdditionalPropertySmnNotification(curJson),
			"log_report_channels": flattenAdditionalPropertyLogReportChannels(curJson),
		},
	}
	return rst
}

func flattenAdditionalPropertySmnNotification(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("smn_notification", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"topic_urn": utils.PathSearch("topic_urn", curJson, nil),
			"events":    utils.PathSearch("events", curJson, nil),
		},
	}
	return rst
}

func flattenAdditionalPropertyLogReportChannels(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("log_report_channels", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"type": utils.PathSearch("type", v, nil),
		})
	}
	return rst
}

func resourceModelartsServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateService: update the ModelArts service.
	var (
		updateServiceHttpUrl = "v1/{project_id}/services/{id}"
		updateServiceProduct = "modelarts"
	)
	updateServiceClient, err := cfg.NewServiceClient(updateServiceProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts Client: %s", err)
	}

	updateServicePath := updateServiceClient.Endpoint + updateServiceHttpUrl
	updateServicePath = strings.ReplaceAll(updateServicePath, "{project_id}", updateServiceClient.ProjectID)
	updateServicePath = strings.ReplaceAll(updateServicePath, "{id}", d.Id())

	updateServiceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	updateServiceChanges := []string{
		"description",
		"schedule",
		"config",
		"additional_properties",
	}

	if d.HasChanges(updateServiceChanges...) {
		updateServiceOpt.JSONBody = utils.RemoveNil(buildUpdateServiceBodyParams(d))
		_, err = updateServiceClient.Request("PUT", updateServicePath, &updateServiceOpt)
		if err != nil {
			return diag.Errorf("error updating ModelartsService: %s", err)
		}
		err = serviceWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the update of ModelartsService (%s) to complete: %s", d.Id(), err)
		}
	}

	if d.HasChange("change_status_to") {
		updateServiceOpt.JSONBody = utils.RemoveNil(buildChangeStateParams(d))
		_, err = updateServiceClient.Request("PUT", updateServicePath, &updateServiceOpt)
		if err != nil {
			return diag.Errorf("error updating ModelartsService: %s", err)
		}
		err = serviceWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the update of ModelartsService (%s) to complete: %s", d.Id(), err)
		}
	}
	return resourceModelartsServiceRead(ctx, d, meta)
}

func buildUpdateServiceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"schedule":              buildServiceRequestBodySchedule(d.Get("schedule")),
		"config":                buildServiceRequestBodyConfig(d.Get("config")),
		"additional_properties": buildServiceRequestBodyAdditionalProperty(d.Get("additional_properties")),
	}
	return bodyParams
}

func buildChangeStateParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"status": utils.ValueIgnoreEmpty(d.Get("change_status_to")),
	}
	return bodyParams
}

func resourceModelartsServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteService: delete ModelArts service
	var (
		deleteServiceHttpUrl = "v1/{project_id}/services/{id}"
		deleteServiceProduct = "modelarts"
	)
	deleteServiceClient, err := cfg.NewServiceClient(deleteServiceProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts Client: %s", err)
	}

	deleteServicePath := deleteServiceClient.Endpoint + deleteServiceHttpUrl
	deleteServicePath = strings.ReplaceAll(deleteServicePath, "{project_id}", deleteServiceClient.ProjectID)
	deleteServicePath = strings.ReplaceAll(deleteServicePath, "{id}", d.Id())

	deleteServiceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	_, err = deleteServiceClient.Request("DELETE", deleteServicePath, &deleteServiceOpt)
	if err != nil {
		return diag.Errorf("error deleting ModelartsService: %s", err)
	}

	err = deleteServiceWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the delete of ModelartsService (%s) to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteServiceWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			var (
				deleteServiceWaitingHttpUrl = "v1/{project_id}/services/{id}"
				deleteServiceWaitingProduct = "modelarts"
			)
			deleteServiceWaitingClient, err := cfg.NewServiceClient(deleteServiceWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating ModelArts Client: %s", err)
			}

			deleteServiceWaitingPath := deleteServiceWaitingClient.Endpoint + deleteServiceWaitingHttpUrl
			deleteServiceWaitingPath = strings.ReplaceAll(deleteServiceWaitingPath, "{project_id}", deleteServiceWaitingClient.ProjectID)
			deleteServiceWaitingPath = strings.ReplaceAll(deleteServiceWaitingPath, "{id}", d.Id())

			deleteServiceWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
				MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
			}

			deleteServiceWaitingResp, err := deleteServiceWaitingClient.Request("GET", deleteServiceWaitingPath, &deleteServiceWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
					return "Resource Not Found", "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			return deleteServiceWaitingResp, "PENDING", nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
