package workspace

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace POST /v2/{project_id}/statistics/notify-rules
// @API Workspace GET /v2/{project_id}/statistics/notify-rules
// @API Workspace PUT /v2/{project_id}/statistics/notify-rules/{rule_id}
// @API Workspace DELETE /v2/{project_id}/statistics/notify-rules/{rule_id}
func ResourceNotificationRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotificationRuleCreate,
		ReadContext:   resourceNotificationRuleRead,
		UpdateContext: resourceNotificationRuleUpdate,
		DeleteContext: resourceNotificationRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the notification rule is located.`,
			},

			// Required parameters.
			"metric_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the metric. Currently only supports "desktop_idle_duration".`,
			},
			"comparison_operator": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The comparison operator for the metric value and threshold.`,
			},
			"enable": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Whether to enable the rule.`,
			},
			"notify_object": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The notification object, which is the SMN topic URN.`,
			},

			// Optional parameters.
			"threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The threshold (in days) for the rule configuration.`,
			},
			"interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The interval time (in days) for the next notification after triggering. Default is once per day.`,
			},
		},
	}
}

func buildNotificationRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metric_name":         d.Get("metric_name"),
		"comparison_operator": d.Get("comparison_operator"),
		"enable":              d.Get("enable"),
		"notify_object":       d.Get("notify_object"),
		"threshold":           utils.ValueIgnoreEmpty(d.Get("threshold")),
		"interval":            utils.ValueIgnoreEmpty(d.Get("interval")),
	}

	return utils.RemoveNil(bodyParams)
}

func resourceNotificationRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/statistics/notify-rules"
	)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildNotificationRuleBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Workspace notification rule: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleId := utils.PathSearch("rule_id", respBody, "").(string)
	if ruleId == "" {
		return diag.Errorf("unable to find rule ID from API response")
	}
	d.SetId(ruleId)

	return resourceNotificationRuleRead(ctx, d, meta)
}

func listNotificationRules(client *golangsdk.ServiceClient, ruleId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/statistics/notify-rules?rule_id={rule_id}&limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{rule_id}", ruleId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%v", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &getOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		items := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, items...)
		if len(items) < limit {
			break
		}
		offset += len(items)
	}

	return result, nil
}

func GetNotificationRuleById(client *golangsdk.ServiceClient, ruleId string) (interface{}, error) {
	rules, err := listNotificationRules(client, ruleId)
	if err != nil {
		return nil, err
	}

	notificationRule := utils.PathSearch(fmt.Sprintf("[?rule_id=='%s']|[0]", ruleId), rules, nil)
	if notificationRule == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/statistics/notify-rules",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the notification rule with ID '%s' has been removed", ruleId)),
			},
		}
	}
	return notificationRule, nil
}

func resourceNotificationRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	respBody, err := GetNotificationRuleById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Workspace notification rule")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("metric_name", utils.PathSearch("metric_name", respBody, nil)),
		d.Set("comparison_operator", utils.PathSearch("comparison_operator", respBody, nil)),
		d.Set("enable", utils.PathSearch("enable", respBody, nil)),
		d.Set("notify_object", utils.PathSearch("notify_object", respBody, nil)),
		d.Set("threshold", utils.PathSearch("threshold", respBody, nil)),
		d.Set("interval", utils.PathSearch("interval", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceNotificationRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/statistics/notify-rules/{rule_id}"
	)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{rule_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildNotificationRuleBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating Workspace notification rule: %s", err)
	}

	return resourceNotificationRuleRead(ctx, d, meta)
}

func resourceNotificationRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/statistics/notify-rules/{rule_id}"
	)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{rule_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// When deleting a deleted rule, the API does not return an error, so CheckDeletedDiag is not used here.
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting Workspace notification rule: %s", err)
	}

	return nil
}
