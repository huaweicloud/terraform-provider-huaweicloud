package tms

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TMS POST /v2/tms/associated-resources/rules/batch-create
// @API TMS GET /v2/tms/associated-resources/rules
// @API TMS PUT /v2/tms/associated-resources/rules/batch-update
// @API TMS DELETE /v2/tms/associated-resources/rules/{setting_name}
func ResourceTmsAssociatedResourceRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTmsAssociatedResourceRuleCreate,
		ReadContext:   resourceTmsAssociatedResourceRuleRead,
		UpdateContext: resourceTmsAssociatedResourceRuleUpdate,
		DeleteContext: resourceTmsAssociatedResourceRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceTmsAssociatedResourceRuleImportState,
		},

		Schema: map[string]*schema.Schema{
			"setting_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the name of the rule setting.",
			},
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the region ID where the rule takes effect.",
			},
			"tag_keys": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the tag keys associated with the rule.",
			},
			"existing_resource_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "disable",
				ValidateFunc: validation.StringInSlice([]string{"enable", "disable"}, false),
				Description:  "Specifies whether the rule takes effect for existing resources.",
			},
			"auto_delete_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "disable",
				ValidateFunc: validation.StringInSlice([]string{"enable", "disable"}, false),
				Description:  "Specifies whether to automatically delete tags when the association is removed.",
			},
		},
	}
}

func resourceTmsAssociatedResourceRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("tms", region)
	if err != nil {
		return diag.Errorf("error creating TMS client: %s", err)
	}

	createPath := client.Endpoint + "v2/tms/associated-resources/rules/batch-create"
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateAssociatedResourceRulesBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating TMS associated resource rule: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkAssociatedResourceRuleErrors(createRespBody); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("setting_name").(string), d.Get("region_id").(string)))

	return resourceTmsAssociatedResourceRuleRead(ctx, d, meta)
}

func buildCreateAssociatedResourceRulesBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"rules": []map[string]interface{}{
			buildAssociatedResourceRuleBodyParams(d),
		},
	}
}

func resourceTmsAssociatedResourceRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("tms", region)
	if err != nil {
		return diag.Errorf("error creating TMS client: %s", err)
	}

	listBasePath := client.Endpoint + "v2/tms/associated-resources/rules"
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	var rule interface{}
	marker := ""
	for {
		listPath := listBasePath + buildGetAssociatedResourceRulePath(marker, d)
		listResp, err := client.Request("GET", listPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving TMS associated resource rules: %s", err)
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		rule = utils.PathSearch(fmt.Sprintf("rules[?setting_name=='%s' && region_id=='%s']|[0]",
			d.Get("setting_name").(string), d.Get("region_id").(string)), listRespBody, nil)
		if rule != nil {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", listRespBody, "").(string)
		if marker == "" {
			break
		}
	}
	if rule == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method: "GET",
				URL:    "/v2/tms/associated-resources/rules",
				Body:   []byte(fmt.Sprintf("the TMS associated resource rule (%s) does not exist", d.Id())),
			},
		}, "error retrieving TMS associated resource rule")
	}

	mErr := multierror.Append(
		d.Set("setting_name", utils.PathSearch("setting_name", rule, nil)),
		d.Set("region_id", utils.PathSearch("region_id", rule, nil)),
		d.Set("tag_keys", utils.PathSearch("tag_keys", rule, nil)),
		d.Set("existing_resource_status",
			normalizeAssociatedResourceStatus(utils.PathSearch("existing_resource_status", rule, "").(string))),
		d.Set("auto_delete_status",
			normalizeAssociatedResourceStatus(utils.PathSearch("auto_delete_status", rule, "").(string))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetAssociatedResourceRulePath(marker string, d *schema.ResourceData) string {
	path := fmt.Sprintf("?limit=100&setting_name=%s&region_id=%s",
		d.Get("setting_name").(string), d.Get("region_id").(string))
	if marker != "" {
		path = fmt.Sprintf("%s&marker=%s", path, marker)
	}
	return path
}

func resourceTmsAssociatedResourceRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("tms", region)
	if err != nil {
		return diag.Errorf("error creating TMS client: %s", err)
	}

	updatePath := client.Endpoint + "v2/tms/associated-resources/rules/batch-update"
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateAssociatedResourceRulesBodyParams(d))

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating TMS associated resource rule: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkAssociatedResourceRuleErrors(updateRespBody); err != nil {
		return diag.FromErr(err)
	}

	return resourceTmsAssociatedResourceRuleRead(ctx, d, meta)
}

func buildUpdateAssociatedResourceRulesBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"rules": []map[string]interface{}{
			utils.RemoveNil(buildAssociatedResourceRuleBodyParams(d)),
		},
	}
}

func buildAssociatedResourceRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"setting_name":             d.Get("setting_name"),
		"tag_keys":                 d.Get("tag_keys"),
		"existing_resource_status": d.Get("existing_resource_status"),
		"auto_delete_status":       d.Get("auto_delete_status"),
		"region_id":                d.Get("region_id"),
	}
}

func checkAssociatedResourceRuleErrors(respBody interface{}) error {
	errors := utils.PathSearch("errors", respBody, make([]interface{}, 0)).([]interface{})
	if len(errors) == 0 {
		return nil
	}

	errMsgs := make([]string, 0, len(errors))
	for _, e := range errors {
		code := utils.PathSearch("code", e, "")
		message := utils.PathSearch("message", e, "")
		errMsgs = append(errMsgs, fmt.Sprintf("code: %v, message: %v", code, message))
	}

	return fmt.Errorf("failed to operate the TMS associated resource rule, errors: %s", strings.Join(errMsgs, "; "))
}

func resourceTmsAssociatedResourceRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("tms", region)
	if err != nil {
		return diag.Errorf("error creating TMS client: %s", err)
	}

	deletePath := client.Endpoint + "v2/tms/associated-resources/rules/{setting_name}"
	deletePath = strings.ReplaceAll(deletePath, "{setting_name}", d.Get("setting_name").(string))
	deletePath += buildDeleteAssociatedResourceRuleQueryParams(d)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting TMS associated resource rule: %s", err)
	}

	return nil
}

func buildDeleteAssociatedResourceRuleQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?region_id=%s", d.Get("region_id").(string))
}

func resourceTmsAssociatedResourceRuleImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <setting_name>/<region_id>")
	}

	d.SetId(parts[0])
	mErr := multierror.Append(
		d.Set("setting_name", parts[0]),
		d.Set("region_id", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
