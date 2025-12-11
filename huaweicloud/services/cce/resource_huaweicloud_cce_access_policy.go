package cce

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE POST /api/v3/access-policies
// @API CCE GET /api/v3/access-policies/{policy_id}
// @API CCE PUT /api/v3/access-policies/{policy_id}
// @API CCE DELETE /api/v3/access-policies/{policy_id}
func ResourceAccessPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccessPolicyCreate,
		ReadContext:   resourceAccessPolicyRead,
		UpdateContext: resourceAccessPolicyUpdate,
		DeleteContext: resourceAccessPolicyDelete,
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
			"clusters": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"access_scope": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespaces": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"policy_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"principal": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ids": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildAccessPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"kind":        "AccessPolicy",
		"apiVersion":  "v3",
		"name":        d.Get("name"),
		"clusters":    d.Get("clusters"),
		"accessScope": buildAccessPolicyScopeParams(d),
		"policyType":  d.Get("policy_type"),
		"principal":   buildAccessPolicyPrincipalParams(d),
	}

	return bodyParams
}

func buildAccessPolicyScopeParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"namespaces": utils.PathSearch("[0].namespaces", d.Get("access_scope"), nil),
	}

	return bodyParams
}

func buildAccessPolicyPrincipalParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type": utils.PathSearch("[0].type", d.Get("principal"), nil),
		"ids":  utils.PathSearch("[0].ids", d.Get("principal"), nil),
	}

	return bodyParams
}

func resourceAccessPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createAccessPolicyHttpUrl = "api/v3/access-policies"
		createAccessPolicyProduct = "cce"
	)
	createAccessPolicyClient, err := cfg.NewServiceClient(createAccessPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	createAccessPolicyPath := createAccessPolicyClient.Endpoint + createAccessPolicyHttpUrl

	createAccessPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createAccessPolicyOpt.JSONBody = utils.RemoveNil(buildAccessPolicyBodyParams(d))
	createAccessPolicyResp, err := createAccessPolicyClient.Request("POST", createAccessPolicyPath, &createAccessPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating CCE access policy: %s", err)
	}

	createAccessPolicyRespBody, err := utils.FlattenResponse(createAccessPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("policyId", createAccessPolicyRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CCE access policy: ID is not found in API response")
	}
	d.SetId(id)

	return resourceAccessPolicyRead(ctx, d, meta)
}

func resourceAccessPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getAccessPolicyHttpUrl = "api/v3/access-policies/{policy_id}"
		getAccessPolicyProduct = "cce"
	)
	getAccessPolicyClient, err := cfg.NewServiceClient(getAccessPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	getAccessPolicyPath := getAccessPolicyClient.Endpoint + getAccessPolicyHttpUrl
	getAccessPolicyPath = strings.ReplaceAll(getAccessPolicyPath, "{policy_id}", d.Id())

	getAccessPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getAccessPolicyResp, err := getAccessPolicyClient.Request("GET", getAccessPolicyPath, &getAccessPolicyOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CCE access policy")
	}

	getAccessPolicyRespBody, err := utils.FlattenResponse(getAccessPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("name", utils.PathSearch("name", getAccessPolicyRespBody, nil)),
		d.Set("clusters", utils.PathSearch("clusters", getAccessPolicyRespBody, nil)),
		d.Set("access_scope", flattenAccessPolicyScope(getAccessPolicyRespBody)),
		d.Set("policy_type", utils.PathSearch("policyType", getAccessPolicyRespBody, nil)),
		d.Set("principal", flattenAccessPolicyPrincipal(getAccessPolicyRespBody)),
		d.Set("created_at", utils.PathSearch("createTime", getAccessPolicyRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updateTime", getAccessPolicyRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAccessPolicyScope(respBody interface{}) []map[string]interface{} {
	namespaces := utils.PathSearch("accessScope.namespaces", respBody, make([]interface{}, 0))
	return []map[string]interface{}{
		{
			"namespaces": namespaces,
		},
	}
}

func flattenAccessPolicyPrincipal(respBody interface{}) []map[string]interface{} {
	principalType := utils.PathSearch("principal.type", respBody, "").(string)
	ids := utils.PathSearch("principal.ids", respBody, make([]interface{}, 0))
	return []map[string]interface{}{
		{
			"type": principalType,
			"ids":  ids,
		},
	}
}

func resourceAccessPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateAccessPolicyProduct = "cce"
		updateAccessPolicyHttpUrl = "api/v3/access-policies/{policy_id}"
	)

	updateAccessPolicyClient, err := cfg.NewServiceClient(updateAccessPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	updateAccessPolicyPath := updateAccessPolicyClient.Endpoint + updateAccessPolicyHttpUrl
	updateAccessPolicyPath = strings.ReplaceAll(updateAccessPolicyPath, "{policy_id}", d.Id())

	updateAccessPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildAccessPolicyBodyParams(d)),
	}

	_, err = updateAccessPolicyClient.Request("PUT", updateAccessPolicyPath, &updateAccessPolicyOpt)
	if err != nil {
		return diag.Errorf("error updating CCE access policy: %s", err)
	}

	return resourceAccessPolicyRead(ctx, d, meta)
}

func resourceAccessPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteAccessPolicyHttpUrl = "api/v3/access-policies/{policy_id}"
		deleteAccessPolicyProduct = "cce"
	)
	deleteAccessPolicyClient, err := cfg.NewServiceClient(deleteAccessPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	deleteAccessPolicyPath := deleteAccessPolicyClient.Endpoint + deleteAccessPolicyHttpUrl
	deleteAccessPolicyPath = strings.ReplaceAll(deleteAccessPolicyPath, "{policy_id}", d.Id())

	deleteAccessPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = deleteAccessPolicyClient.Request("DELETE", deleteAccessPolicyPath, &deleteAccessPolicyOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CCE access policy")
	}

	return nil
}
