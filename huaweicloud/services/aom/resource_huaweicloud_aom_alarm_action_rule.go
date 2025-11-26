// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product AOM
// ---------------------------------------------------------------

package aom

import (
	"context"
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

// @API AOM DELETE /v2/{project_id}/alert/action-rules
// @API AOM POST /v2/{project_id}/alert/action-rules
// @API AOM PUT /v2/{project_id}/alert/action-rules
// @API AOM GET /v2/{project_id}/alert/action-rules/{rule_name}
func ResourceAlarmActionRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmActionRuleCreate,
		UpdateContext: resourceAlarmActionRuleUpdate,
		ReadContext:   resourceAlarmActionRuleRead,
		DeleteContext: resourceAlarmActionRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: utils.SchemaDesc(
					`The IAM user name to which the action rule belongs.`,
					utils.SchemaDescInput{
						Required: true,
					},
				),
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the action rule name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the action rule type.`,
			},
			"smn_topics": {
				Type:        schema.TypeList,
				Elem:        AlarmActionRuleSmnTopicsSchema(),
				Required:    true,
				Description: `Specifies the SMN topic configurations.`,
			},
			"notification_template": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the notification template.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the action rule description.`,
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The last update time.`,
			},
		},
	}
}

func AlarmActionRuleSmnTopicsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"topic_urn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the SMN topic URN.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the SMN topic name.`,
			},
		},
	}
	return &sc
}

func resourceAlarmActionRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createAlarmActionRule: create a Alarm Action Rule.
	var (
		createAlarmActionRuleHttpUrl = "v2/{project_id}/alert/action-rules"
		createAlarmActionRuleProduct = "aom"
	)
	createAlarmActionRuleClient, err := cfg.NewServiceClient(createAlarmActionRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating AOM Client: %s", err)
	}

	createAlarmActionRulePath := createAlarmActionRuleClient.Endpoint + createAlarmActionRuleHttpUrl
	createAlarmActionRulePath = strings.ReplaceAll(createAlarmActionRulePath, "{project_id}", createAlarmActionRuleClient.ProjectID)

	createAlarmActionRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createAlarmActionRuleOpt.JSONBody = utils.RemoveNil(buildAlarmActionRuleBodyParams(d))
	_, err = createAlarmActionRuleClient.Request("POST", createAlarmActionRulePath, &createAlarmActionRuleOpt)
	if err != nil {
		return diag.Errorf("error creating AlarmActionRule: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceAlarmActionRuleRead(ctx, d, meta)
}

func buildAlarmActionRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"user_name":             utils.ValueIgnoreEmpty(d.Get("user_name")),
		"rule_name":             utils.ValueIgnoreEmpty(d.Get("name")),
		"desc":                  utils.ValueIgnoreEmpty(d.Get("description")),
		"type":                  utils.ValueIgnoreEmpty(d.Get("type")),
		"smn_topics":            buildAlarmActionRuleRequestBodySmnTopics(d.Get("smn_topics")),
		"notification_template": utils.ValueIgnoreEmpty(d.Get("notification_template")),
	}
	return bodyParams
}

func buildAlarmActionRuleRequestBodySmnTopics(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"topic_urn": utils.ValueIgnoreEmpty(raw["topic_urn"]),
			}

			if raw["name"] != "" {
				rst[i]["name"] = raw["name"]
			} else {
				parts := strings.Split(raw["topic_urn"].(string), ":")
				rst[i]["name"] = parts[len(parts)-1]
			}
		}
		return rst
	}
	return nil
}

func getAlarmActionRule(client *golangsdk.ServiceClient, ruleName string) (interface{}, error) {
	httpUrl := "v2/{project_id}/alert/action-rules/{rule_name}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{rule_name}", ruleName)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceAlarmActionRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM Client: %s", err)
	}

	respBody, err := getAlarmActionRule(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AlarmActionRule")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("user_name", utils.PathSearch("user_name", respBody, nil)),
		d.Set("name", utils.PathSearch("rule_name", respBody, nil)),
		d.Set("description", utils.PathSearch("desc", respBody, nil)),
		d.Set("type", utils.PathSearch("type", respBody, nil)),
		d.Set("smn_topics", flattenGetAlarmActionRuleResponseBodySmnTopics(respBody)),
		d.Set("notification_template", utils.PathSearch("notification_template", respBody, nil)),
		d.Set("created_at", utils.PathSearch("create_time", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetAlarmActionRuleResponseBodySmnTopics(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("smn_topics", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"topic_urn": utils.PathSearch("topic_urn", v, nil),
			"name":      utils.PathSearch("name", v, nil),
		})
	}
	return rst
}

func resourceAlarmActionRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateAlarmActionRuleChanges := []string{
		"description",
		"type",
		"smn_topics",
		"notification_template",
	}

	if d.HasChanges(updateAlarmActionRuleChanges...) {
		// updateAlarmActionRule: update the Alarm Action Rule
		var (
			updateAlarmActionRuleHttpUrl = "v2/{project_id}/alert/action-rules"
			updateAlarmActionRuleProduct = "aom"
		)
		updateAlarmActionRuleClient, err := cfg.NewServiceClient(updateAlarmActionRuleProduct, region)
		if err != nil {
			return diag.Errorf("error creating AOM Client: %s", err)
		}

		updateAlarmActionRulePath := updateAlarmActionRuleClient.Endpoint + updateAlarmActionRuleHttpUrl
		updateAlarmActionRulePath = strings.ReplaceAll(updateAlarmActionRulePath, "{project_id}", updateAlarmActionRuleClient.ProjectID)

		updateAlarmActionRuleOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}
		updateAlarmActionRuleOpt.JSONBody = utils.RemoveNil(buildAlarmActionRuleBodyParams(d))
		_, err = updateAlarmActionRuleClient.Request("PUT", updateAlarmActionRulePath, &updateAlarmActionRuleOpt)
		if err != nil {
			return diag.Errorf("error updating AlarmActionRule: %s", err)
		}
	}
	return resourceAlarmActionRuleRead(ctx, d, meta)
}

func resourceAlarmActionRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM Client: %s", err)
	}

	var (
		httpUrl = "v2/{project_id}/alert/action-rules"
	)

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	deleteOpt.JSONBody = []string{d.Id()}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting AlarmActionRule: %s", err)
	}

	err = waitForAlarmActionRuleDeleted(ctx, client, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the alarm action rule (%s) to be deleted: %s", d.Id(), err)
	}

	return nil
}

func waitForAlarmActionRuleDeleted(ctx context.Context, client *golangsdk.ServiceClient, ruleName string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshAlarmActionRuleStatus(client, ruleName),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshAlarmActionRuleStatus(client *golangsdk.ServiceClient, ruleName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := getAlarmActionRule(client, ruleName)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "ResourceNotFound", "COMPLETED", nil
			}
			return nil, "ERROR", err
		}

		return respBody, "PENDING", nil
	}
}
