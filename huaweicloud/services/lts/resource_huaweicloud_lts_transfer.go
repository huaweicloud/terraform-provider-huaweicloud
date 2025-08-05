// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product LTS
// ---------------------------------------------------------------

package lts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LTS DELETE /v2/{project_id}/transfers
// @API LTS GET /v2/{project_id}/transfers
// @API LTS POST /v2/{project_id}/transfers
// @API LTS PUT /v2/{project_id}/transfers
func ResourceLtsTransfer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLtsTransferCreate,
		UpdateContext: resourceLtsTransferUpdate,
		ReadContext:   resourceLtsTransferRead,
		DeleteContext: resourceLtsTransferDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"log_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Log group ID.`,
			},
			"log_streams": {
				Type:        schema.TypeList,
				Elem:        ltsTransferLogStreamsSchema(),
				Required:    true,
				ForceNew:    true,
				Description: `The list of log streams.`,
			},
			"log_transfer_info": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     ltsTransferLogTransferInfoSchema(),
				Required: true,
			},
			"log_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Log group name.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the log transfer, in RFC3339 format.`,
			},
		},
	}
}

func ltsTransferLogStreamsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"log_stream_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Log stream ID.`,
			},
			"log_stream_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `Log stream name.`,
			},
		},
	}
	return &sc
}

func ltsTransferLogTransferInfoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"log_transfer_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Log transfer type.`,
			},
			"log_transfer_mode": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Log transfer mode.`,
			},
			"log_storage_format": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Log transfer format.`,
			},
			"log_transfer_status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Log transfer status.`,
			},
			"log_agency_transfer": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     ltsTransferLogAgencySchema(),
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"log_transfer_detail": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     ltsTransferLogDetailSchema(),
				Required: true,
			},
		},
	}
	return &sc
}

func ltsTransferLogAgencySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"agency_domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Delegator account ID.`,
			},
			"agency_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Delegator account name.`,
			},
			"agency_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The agency name created by the delegator.`,
			},
			"agency_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Project ID of the delegator.`,
			},
		},
	}
	return &sc
}

func ltsTransferLogDetailSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"obs_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Length of the transfer interval for an OBS transfer task.`,
			},
			"obs_period_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Unit of the transfer interval for an OBS transfer task.`,
			},
			"obs_bucket_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `OBS bucket name.`,
			},
			"obs_transfer_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `OBS bucket path, which is the log transfer destination.`,
			},
			"obs_dir_prefix_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Custom transfer path of an OBS transfer task.`,
			},
			"obs_prefix_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Transfer file prefix of an OBS transfer task.`,
			},
			"obs_eps_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Enterprise project ID of an OBS transfer task.`,
			},
			"obs_encrypted_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether OBS bucket encryption is enabled.`,
			},
			"obs_encrypted_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `KMS key ID for an OBS transfer task.`,
			},
			"obs_time_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Time zone for an OBS transfer task.`,
			},
			"obs_time_zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `ID of the time zone for an OBS transfer task.`,
			},
			"dis_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `DIS stream ID.`,
			},
			"dis_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `DIS stream name.`,
			},
			"kafka_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Kafka ID.`,
			},
			"kafka_topic": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Kafka topic.`,
			},
			"lts_tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of built-in fields and custom tags to be transferred.`,
			},
			"stream_tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of stream tag fields to be transferred.`,
			},
			"struct_fields": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of structured fields to be transferred.`,
			},
			"invalid_field_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The value of the invalid field fill.`,
			},
			"delivery_tags": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `The list of tag fields will be delivered when transferring.`,
			},
			"cloud_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The cloud project ID.`,
			},
		},
	}
	return &sc
}

func resourceLtsTransferCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createTransfer: create a log transfer task.
	var (
		createTransferHttpUrl = "v2/{project_id}/transfers"
		createTransferProduct = "lts"
	)
	createTransferClient, err := cfg.NewServiceClient(createTransferProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	createTransferPath := createTransferClient.Endpoint + createTransferHttpUrl
	createTransferPath = strings.ReplaceAll(createTransferPath, "{project_id}", createTransferClient.ProjectID)

	createTransferOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createTransferOpt.JSONBody = utils.RemoveNil(buildCreateTransferBodyParams(d, cfg.DomainID, createTransferClient.ProjectID))
	createTransferResp, err := createTransferClient.Request("POST", createTransferPath, &createTransferOpt)
	if err != nil {
		return diag.Errorf("error creating LTS transfer: %s", err)
	}

	createTransferRespBody, err := utils.FlattenResponse(createTransferResp)
	if err != nil {
		return diag.FromErr(err)
	}

	transferId := utils.PathSearch("log_transfer_id", createTransferRespBody, "").(string)
	if transferId == "" {
		return diag.Errorf("unable to find the LTS transfer ID from the API response")
	}
	d.SetId(transferId)

	return resourceLtsTransferRead(ctx, d, meta)
}

func buildCreateTransferBodyParams(d *schema.ResourceData, domainID, projectID string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"log_group_id":      utils.ValueIgnoreEmpty(d.Get("log_group_id")),
		"log_streams":       buildCreateTransferRequestBodyLogStreams(d.Get("log_streams")),
		"log_transfer_info": buildCreateTransferRequestBodyLogTransferInfo(d.Get("log_transfer_info"), domainID, projectID),
	}
	return bodyParams
}

func buildCreateTransferRequestBodyLogStreams(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"log_stream_id":   utils.ValueIgnoreEmpty(raw["log_stream_id"]),
				"log_stream_name": utils.ValueIgnoreEmpty(raw["log_stream_name"]),
			}
		}
		return rst
	}
	return nil
}

func buildCreateTransferRequestBodyLogTransferInfo(rawParams interface{}, domainID, projectID string) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"log_transfer_type":   utils.ValueIgnoreEmpty(raw["log_transfer_type"]),
			"log_transfer_mode":   utils.ValueIgnoreEmpty(raw["log_transfer_mode"]),
			"log_storage_format":  utils.ValueIgnoreEmpty(raw["log_storage_format"]),
			"log_transfer_status": utils.ValueIgnoreEmpty(raw["log_transfer_status"]),
			"log_agency_transfer": buildLogTransferInfoLogAgency(raw["log_agency_transfer"], domainID, projectID),
			"log_transfer_detail": buildLogTransferInfoLogTransferDetail(raw["log_transfer_detail"]),
		}
		return params
	}
	return nil
}

func buildLogTransferInfoLogAgency(rawParams interface{}, domainID, projectID string) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"agency_domain_id":     utils.ValueIgnoreEmpty(raw["agency_domain_id"]),
			"agency_domain_name":   utils.ValueIgnoreEmpty(raw["agency_domain_name"]),
			"agency_name":          utils.ValueIgnoreEmpty(raw["agency_name"]),
			"agency_project_id":    utils.ValueIgnoreEmpty(raw["agency_project_id"]),
			"be_agency_domain_id":  domainID,
			"be_agency_project_id": projectID,
		}
		return params
	}
	return nil
}

func buildLogTransferInfoLogTransferDetail(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"obs_period":           utils.ValueIgnoreEmpty(raw["obs_period"]),
			"obs_period_unit":      utils.ValueIgnoreEmpty(raw["obs_period_unit"]),
			"obs_bucket_name":      utils.ValueIgnoreEmpty(raw["obs_bucket_name"]),
			"obs_transfer_path":    utils.ValueIgnoreEmpty(raw["obs_transfer_path"]),
			"obs_dir_pre_fix_name": utils.ValueIgnoreEmpty(raw["obs_dir_prefix_name"]),
			"obs_prefix_name":      utils.ValueIgnoreEmpty(raw["obs_prefix_name"]),
			"obs_eps_id":           utils.ValueIgnoreEmpty(raw["obs_eps_id"]),
			"obs_encrypted_enable": utils.ValueIgnoreEmpty(raw["obs_encrypted_enable"]),
			"obs_encrypted_id":     utils.ValueIgnoreEmpty(raw["obs_encrypted_id"]),
			"obs_time_zone":        utils.ValueIgnoreEmpty(raw["obs_time_zone"]),
			"obs_time_zone_id":     utils.ValueIgnoreEmpty(raw["obs_time_zone_id"]),
			"dis_id":               utils.ValueIgnoreEmpty(raw["dis_id"]),
			"dis_name":             utils.ValueIgnoreEmpty(raw["dis_name"]),
			"kafka_id":             utils.ValueIgnoreEmpty(raw["kafka_id"]),
			"kafka_topic":          utils.ValueIgnoreEmpty(raw["kafka_topic"]),
			"lts_tags":             raw["lts_tags"].(*schema.Set).List(),
			"stream_tags":          raw["stream_tags"].(*schema.Set).List(),
			"struct_fields":        raw["struct_fields"].(*schema.Set).List(),
			"invalid_field_value":  utils.ValueIgnoreEmpty(raw["invalid_field_value"]),
			"tags":                 utils.ValueIgnoreEmpty(raw["delivery_tags"]),
			"cloud_project_id":     utils.ValueIgnoreEmpty(raw["cloud_project_id"]),
		}
		return params
	}
	return nil
}

func resourceLtsTransferRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getTransfer: Query the log transfer task.
	var (
		getTransferHttpUrl = "v2/{project_id}/transfers"
		getTransferProduct = "lts"
	)
	getTransferClient, err := cfg.NewServiceClient(getTransferProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	getTransferPath := getTransferClient.Endpoint + getTransferHttpUrl
	getTransferPath = strings.ReplaceAll(getTransferPath, "{project_id}", getTransferClient.ProjectID)

	getTransferOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getTransferResp, err := getTransferClient.Request("GET", getTransferPath, &getTransferOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving LTS transfer")
	}

	getTransferRespBody, err := utils.FlattenResponse(getTransferResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("log_transfers[?log_transfer_id =='%s']|[0]", d.Id())
	getTransferRespBody = utils.PathSearch(jsonPath, getTransferRespBody, nil)
	if getTransferRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("log_group_id", utils.PathSearch("log_group_id", getTransferRespBody, nil)),
		d.Set("log_group_name", utils.PathSearch("log_group_name", getTransferRespBody, nil)),
		d.Set("log_streams", flattenGetTransferResponseBodyLogStreams(getTransferRespBody)),
		d.Set("log_transfer_info", flattenGetTransferResponseBodyLogTransferInfo(getTransferRespBody, d)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("log_transfer_info.log_create_time", getTransferRespBody, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetTransferResponseBodyLogStreams(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("log_streams", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"log_stream_id":   utils.PathSearch("log_stream_id", v, nil),
			"log_stream_name": utils.PathSearch("log_stream_name", v, nil),
		})
	}
	return rst
}

func flattenGetTransferResponseBodyLogTransferInfo(resp interface{}, d *schema.ResourceData) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("log_transfer_info", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"log_transfer_type":   utils.PathSearch("log_transfer_type", curJson, nil),
			"log_transfer_mode":   utils.PathSearch("log_transfer_mode", curJson, nil),
			"log_storage_format":  utils.PathSearch("log_storage_format", curJson, nil),
			"log_transfer_status": utils.PathSearch("log_transfer_status", curJson, nil),
			"log_agency_transfer": flattenLogTransferInfoLogAgency(curJson),
			"log_transfer_detail": flattenLogTransferInfoLogTransferDetail(curJson, d),
		},
	}
	return rst
}

func flattenLogTransferInfoLogAgency(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("log_agency_transfer", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"agency_domain_id":   utils.PathSearch("agency_domain_id", curJson, nil),
			"agency_domain_name": utils.PathSearch("agency_domain_name", curJson, nil),
			"agency_name":        utils.PathSearch("agency_name", curJson, nil),
			"agency_project_id":  utils.PathSearch("agency_project_id", curJson, nil),
		},
	}
	return rst
}

func flattenLogTransferInfoLogTransferDetail(resp interface{}, d *schema.ResourceData) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("log_transfer_detail", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	logTransferDetail := map[string]interface{}{
		"obs_period":           utils.PathSearch("obs_period", curJson, nil),
		"obs_period_unit":      utils.PathSearch("obs_period_unit", curJson, nil),
		"obs_bucket_name":      utils.PathSearch("obs_bucket_name", curJson, nil),
		"obs_transfer_path":    utils.PathSearch("obs_transfer_path", curJson, nil),
		"obs_dir_prefix_name":  utils.PathSearch("obs_dir_pre_fix_name", curJson, nil),
		"obs_prefix_name":      utils.PathSearch("obs_prefix_name", curJson, nil),
		"obs_eps_id":           utils.PathSearch("obs_eps_id", curJson, nil),
		"obs_encrypted_enable": utils.PathSearch("obs_encrypted_enable", curJson, nil),
		"obs_encrypted_id":     utils.PathSearch("obs_encrypted_id", curJson, nil),
		"obs_time_zone":        utils.PathSearch("obs_time_zone", curJson, nil),
		"obs_time_zone_id":     utils.PathSearch("obs_time_zone_id", curJson, nil),
		"dis_id":               utils.PathSearch("dis_id", curJson, nil),
		"dis_name":             utils.PathSearch("dis_name", curJson, nil),
		"kafka_id":             utils.PathSearch("kafka_id", curJson, nil),
		"kafka_topic":          utils.PathSearch("kafka_topic", curJson, nil),
		"lts_tags":             utils.PathSearch("lts_tags", curJson, nil),
		"stream_tags":          utils.PathSearch("stream_tags", curJson, nil),
		"struct_fields":        utils.PathSearch("struct_fields", curJson, nil),
		"delivery_tags":        utils.PathSearch("tags", curJson, nil),
		"cloud_project_id":     utils.PathSearch("cloud_project_id", curJson, nil),
	}

	invalidFieldValue, ok := d.GetOk("log_transfer_info.0.log_transfer_detail.0.invalid_field_value")
	if ok {
		logTransferDetail["invalid_field_value"] = invalidFieldValue
	}

	rst = []interface{}{
		logTransferDetail,
	}
	return rst
}

func resourceLtsTransferUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateTransferChanges := []string{
		"log_transfer_info",
	}

	if d.HasChanges(updateTransferChanges...) {
		var (
			updateTransferHttpUrl = "v2/{project_id}/transfers"
			updateTransferProduct = "lts"
		)
		updateTransferClient, err := cfg.NewServiceClient(updateTransferProduct, region)
		if err != nil {
			return diag.Errorf("error creating LTS client: %s", err)
		}

		updateTransferPath := updateTransferClient.Endpoint + updateTransferHttpUrl
		updateTransferPath = strings.ReplaceAll(updateTransferPath, "{project_id}", updateTransferClient.ProjectID)
		updateTransferPath = strings.ReplaceAll(updateTransferPath, "{id}", d.Id())

		updateTransferOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/json"},
		}

		updateTransferOpt.JSONBody = utils.RemoveNil(buildUpdateTransferBodyParams(d))
		_, err = updateTransferClient.Request("PUT", updateTransferPath, &updateTransferOpt)
		if err != nil {
			return diag.Errorf("error updating LTS transfer: %s", err)
		}
	}
	return resourceLtsTransferRead(ctx, d, meta)
}

func buildUpdateTransferBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"log_transfer_id":   utils.ValueIgnoreEmpty(d.Id()),
		"log_transfer_info": buildUpdateTransferRequestBodyLogTransferInfoUpdate(d.Get("log_transfer_info")),
	}
	return bodyParams
}

func buildUpdateTransferRequestBodyLogTransferInfoUpdate(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"log_storage_format":  utils.ValueIgnoreEmpty(raw["log_storage_format"]),
			"log_transfer_status": utils.ValueIgnoreEmpty(raw["log_transfer_status"]),
			"log_transfer_detail": buildLogTransferInfoLogTransferDetail(raw["log_transfer_detail"]),
		}
		return params
	}
	return nil
}

func resourceLtsTransferDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteTransfer: delete log transfer task
	var (
		deleteTransferHttpUrl = "v2/{project_id}/transfers"
		deleteTransferProduct = "lts"
	)
	deleteTransferClient, err := cfg.NewServiceClient(deleteTransferProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	deleteTransferPath := deleteTransferClient.Endpoint + deleteTransferHttpUrl
	deleteTransferPath = strings.ReplaceAll(deleteTransferPath, "{project_id}", deleteTransferClient.ProjectID)

	deleteTransferPath += fmt.Sprintf("?log_transfer_id=%s", d.Id())

	deleteTransferOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err = deleteTransferClient.Request("DELETE", deleteTransferPath, &deleteTransferOpt)
	if err != nil {
		return diag.Errorf("error deleting LTS transfer: %s", err)
	}

	return nil
}
