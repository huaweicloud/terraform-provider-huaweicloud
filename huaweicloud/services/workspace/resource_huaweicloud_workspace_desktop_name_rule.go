package workspace

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

// @API Workspace POST /v2/{project_id}/desktop-name-policies
// @API Workspace GET /v2/{project_id}/desktop-name-policies
// @API Workspace PUT /v2/{project_id}/desktop-name-policies/{policy_id}
// @API Workspace POST /v2/{project_id}/desktop-name-policies/batch-delete
func ResourceDesktopNameRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDesktopNameRuleCreate,
		ReadContext:   resourceDesktopNameRuleRead,
		UpdateContext: resourceDesktopNameRuleUpdate,
		DeleteContext: resourceDesktopNameRuleDelete,

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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},
			"digit_number": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"start_number": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"single_domain_user_increment": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"is_default_policy": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_contain_user": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceDesktopNameRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/desktop-name-policies"
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildDesktopNameRuleBodyParams(d)),
	}
	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating desktop name rule: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	policyId := utils.PathSearch("policy_id", respBody, "")
	d.SetId(policyId.(string))

	return resourceDesktopNameRuleRead(ctx, d, meta)
}

func buildDesktopNameRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"policy_name":            d.Get("name"),
		"name_prefix":            d.Get("name_prefix"),
		"digit_number":           d.Get("digit_number"),
		"start_number":           d.Get("start_number"),
		"single_domain_user_inc": d.Get("single_domain_user_increment"),
		"is_default_policy":      d.Get("is_default_policy"),
	}
}

func resourceDesktopNameRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	rule, err := GetDesktopNameRule(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "desktop name rule")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("policy_name", rule, nil)),
		d.Set("is_contain_user", utils.PathSearch("is_contain_user", rule, false)),
		d.Set("name_prefix", utils.PathSearch("name_prefix", rule, nil)),
		d.Set("digit_number", utils.PathSearch("digit_number", rule, nil)),
		d.Set("start_number", utils.PathSearch("start_number", rule, nil)),
		d.Set("single_domain_user_increment", utils.PathSearch("single_domain_user_inc", rule, nil)),
		d.Set("is_default_policy", utils.PathSearch("is_default_policy", rule, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetDesktopNameRule(client *golangsdk.ServiceClient, policyId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/desktop-name-policies"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += fmt.Sprintf("?policy_id=%v", policyId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	rule := utils.PathSearch("desktop_name_policy_infos|[0]", respBody, nil)
	if rule == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return rule, nil
}

func resourceDesktopNameRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/desktop-name-policies/{policy_id}"
	)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{policy_id}", d.Id())
	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		JSONBody: utils.RemoveNil(buildDesktopNameRuleBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return diag.Errorf("error updating desktop name rule: %s", err)
	}

	return resourceDesktopNameRuleRead(ctx, d, meta)
}

func resourceDesktopNameRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/desktop-name-policies/batch-delete"
	)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDeleteDesktopNameRuleParams(d),
		OkCodes: []int{
			204,
		},
	}

	_, err = client.Request("POST", deletePath, &deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting desktop name rule: %s", err)
	}
	return nil
}

func buildDeleteDesktopNameRuleParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"policy_ids": []string{d.Id()},
	}
}
