package waf

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

var batchDeleteAlarmNotificationsNonUpdatableParams = []string{
	"enterprise_project_id",
	"alert_notice_configs",
	"alert_notice_configs.*.id",
}

// @API WAF POST /v2/{project_id}/waf/alert/batch-delete
func ResourceWafBatchDeleteAlarmNotifications() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafBatchDeleteAlarmNotificationsCreate,
		ReadContext:   resourceWafBatchDeleteAlarmNotificationsRead,
		UpdateContext: resourceWafBatchDeleteAlarmNotificationsUpdate,
		DeleteContext: resourceWafBatchDeleteAlarmNotificationsDelete,

		CustomizeDiff: config.FlexibleForceNew(batchDeleteAlarmNotificationsNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alert_notice_configs": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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

func buildBatchDeleteAlarmNotificationsBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawArray := d.Get("alert_notice_configs").([]interface{})
	noticeConfigs := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		noticeConfigs = append(noticeConfigs, map[string]interface{}{
			"id": rawMap["id"],
		})
	}

	bodyParams := map[string]interface{}{
		"alert_notice_configs": noticeConfigs,
	}

	return bodyParams
}

func resourceWafBatchDeleteAlarmNotificationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/waf/alert/batch-delete"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath += fmt.Sprintf("?enterpriseProjectId=%s", epsId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
			"X-Language":   "en-us",
		},
		JSONBody: buildBatchDeleteAlarmNotificationsBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error batch deleting WAF alarm notifications: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	return resourceWafBatchDeleteAlarmNotificationsRead(ctx, d, meta)
}

func resourceWafBatchDeleteAlarmNotificationsRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is an action resource.
	return nil
}

func resourceWafBatchDeleteAlarmNotificationsUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is an action resource.
	return nil
}

func resourceWafBatchDeleteAlarmNotificationsDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch delete alarm notifications. Deleting this resource
    will not clear the corresponding request record, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
