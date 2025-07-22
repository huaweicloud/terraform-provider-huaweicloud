package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS DELETE /v5/{project_id}/event/white-list/alarm
var eventAlarmWhiteListDeleteNonUpdatableParams = []string{"enterprise_project_id", "data_list",
	"data_list.*.event_type", "data_list.*.hash", "data_list.*.description", "data_list.*.delete_white_rule",
	"data_list.*.white_field", "data_list.*.judge_type", "data_list.*.field_value", "data_list.*.file_hash",
	"data_list.*.file_path", "restore_alarm", "delete_all", "event_type"}

func ResourceEventAlarmWhiteListDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventAlarmWhiteListDeleteCreate,
		ReadContext:   resourceEventAlarmWhiteListDeleteRead,
		UpdateContext: resourceEventAlarmWhiteListDeleteUpdate,
		DeleteContext: resourceEventAlarmWhiteListDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(eventAlarmWhiteListDeleteNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			// Currently, this parameter is optional in the API documentation, but it is actually required.
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Currently, this parameter is optional in the API documentation, but it is actually required.
						"event_type": {
							Type:     schema.TypeInt,
							Required: true,
						},
						// Currently, this parameter is optional in the API documentation, but it is actually required.
						"hash": {
							Type:     schema.TypeString,
							Required: true,
						},
						// Currently, this parameter is optional in the API documentation, but it is actually required.
						"description": {
							Type:     schema.TypeString,
							Required: true,
						},
						"delete_white_rule": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"white_field": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"judge_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"field_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"file_hash": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"file_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"restore_alarm": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"delete_all": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"event_type": {
				Type:     schema.TypeInt,
				Optional: true,
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

func buildBoolParam(val interface{}) interface{} {
	b, ok := val.(bool)
	if ok && b {
		return true
	}

	return nil
}

func buildEventAlarmWhiteListDeleteBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"restore_alarm": buildBoolParam(d.Get("restore_alarm")),
		"delete_all":    buildBoolParam(d.Get("delete_all")),
		"event_type":    utils.ValueIgnoreEmpty(d.Get("event_type")),
	}

	rawDataList := d.Get("data_list").([]interface{})
	dataListParams := make([]map[string]interface{}, 0, len(rawDataList))

	for _, rawData := range rawDataList {
		dataMap := map[string]interface{}{
			"event_type":        utils.ValueIgnoreEmpty(utils.PathSearch("event_type", rawData, nil)),
			"hash":              utils.ValueIgnoreEmpty(utils.PathSearch("hash", rawData, nil)),
			"description":       utils.ValueIgnoreEmpty(utils.PathSearch("description", rawData, nil)),
			"white_field":       utils.ValueIgnoreEmpty(utils.PathSearch("white_field", rawData, nil)),
			"judge_type":        utils.ValueIgnoreEmpty(utils.PathSearch("judge_type", rawData, nil)),
			"field_value":       utils.ValueIgnoreEmpty(utils.PathSearch("field_value", rawData, nil)),
			"file_hash":         utils.ValueIgnoreEmpty(utils.PathSearch("file_hash", rawData, nil)),
			"file_path":         utils.ValueIgnoreEmpty(utils.PathSearch("file_path", rawData, nil)),
			"delete_white_rule": buildBoolParam(utils.PathSearch("delete_white_rule", rawData, nil)),
		}

		dataListParams = append(dataListParams, dataMap)
	}

	bodyParams["data_list"] = dataListParams

	return bodyParams
}

func resourceEventAlarmWhiteListDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	requestPath := client.Endpoint + "v5/{project_id}/event/white-list/alarm"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	if epsId != "" {
		requestPath += fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildEventAlarmWhiteListDeleteBodyParams(d)),
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting HSS event alarm white list: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return nil
}

func resourceEventAlarmWhiteListDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceEventAlarmWhiteListDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceEventAlarmWhiteListDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to delete HSS alarm white list. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information from the
    tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
