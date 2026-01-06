package aom

import (
	"context"
	"errors"
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

var prometheusRecordingRuleNonUpdatableParams = []string{"instance_id"}

// @API AOM POST /v1/{project_id}/{prometheus_instance}/aom/api/v1/rules
// @API AOM PUT /v1/{project_id}/{prometheus_instance}/aom/api/v1/rules/{rule_id}
// @API AOM GET /v1/{project_id}/{prometheus_instance}/aom/api/v1/rules
func ResourceRecordingRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecordingRuleCreate,
		ReadContext:   resourceRecordingRuleRead,
		UpdateContext: resourceRecordingRuleUpdate,
		DeleteContext: resourceRecordingRuleDelete,

		CustomizeDiff: config.FlexibleForceNew(prometheusRecordingRuleNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceRecordingRuleImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the Prometheus instance is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Prometheus instance.`,
			},
			"recording_rule": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The content of the recording rule, in YAML format.`,
			},
		},
	}
}

func buildRecordingRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"recording_rule": d.Get("recording_rule").(string),
	}
}

func createRecordingRule(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/{prometheus_instance}/aom/api/v1/rules"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{prometheus_instance}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildRecordingRuleBodyParams(d),
	}

	_, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return fmt.Errorf("error creating recording rule: %s", err)
	}
	return nil
}

func GetRecordingRuleByInstanceId(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/{prometheus_instance}/aom/api/v1/rules"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{prometheus_instance}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return "", err
	}

	respBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return "", err
	}

	return respBody, nil
}

func resourceRecordingRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	respBody, _ := GetRecordingRuleByInstanceId(client, instanceId)
	ruleId := utils.PathSearch("rule_id", respBody, "").(string)
	if ruleId != "" {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Recording rule already exists",
				Detail: fmt.Sprintf(
					"A recording rule for Prometheus instance %s already exists (ID: %s).\n"+
						"Please import it using:\n"+
						"terraform import huaweicloud_aom_recording_rule.<name> %s/%s",
					instanceId, ruleId,
					instanceId, ruleId,
				),
			},
		}
	}

	err = createRecordingRule(client, d)
	if err != nil {
		return diag.Errorf("error creating Prometheus recording rule: %s", err)
	}

	respBody, err = GetRecordingRuleByInstanceId(client, instanceId)
	if err != nil {
		return diag.Errorf("error querying Prometheus recording rule ID: %s", err)
	}
	ruleId = utils.PathSearch("rule_id", respBody, "").(string)
	if ruleId == "" {
		return diag.Errorf("error setting Prometheus recording rule ID, ruleId is empty.")
	}
	d.SetId(ruleId)

	return resourceRecordingRuleRead(ctx, d, meta)
}

func resourceRecordingRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	respBody, err := GetRecordingRuleByInstanceId(client, instanceId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying Prometheus recording rule")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("recording_rule", utils.PathSearch("recording_rule", respBody, "").(string)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func updateRecordingRule(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/{prometheus_instance}/aom/api/v1/rules/{rule_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{prometheus_instance}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{rule_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildRecordingRuleBodyParams(d),
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating Prometheus recording rule: %s", err)
	}
	return nil
}

func resourceRecordingRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	err = updateRecordingRule(client, d)
	if err != nil {
		return diag.Errorf("error updating Prometheus recoding rule: %s", err)
	}

	return resourceRecordingRuleRead(ctx, d, meta)
}

func resourceRecordingRuleDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This is resource bind with the Prometheus instance resource. Deleting this resource will not clear
 the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourceRecordingRuleImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")

	if len(parts) != 2 && len(parts) != 1 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<rule_id>' or '<instance_id>', but got '%s'",
			importedId)
	}

	if len(parts) == 2 {
		d.SetId(parts[1])
	} else {
		var (
			cfg        = meta.(*config.Config)
			region     = cfg.GetRegion(d)
			instanceId = parts[0]
		)

		client, err := cfg.NewServiceClient("aom", region)
		if err != nil {
			return nil, fmt.Errorf("error creating AOM client: %s", err)
		}

		respBody, err := GetRecordingRuleByInstanceId(client, instanceId)
		if err != nil {
			return nil, fmt.Errorf("error querying Prometheus recording rule ID: %s", err)
		}
		ruleId := utils.PathSearch("rule_id", respBody, "").(string)
		if ruleId == "" {
			return nil, errors.New("error setting Prometheus recording rule ID, ruleId is empty")
		}
		d.SetId(ruleId)
	}

	return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
}
