package elb

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

// @API ELB POST /v3/{project_id}/elb/security-policies
// @API ELB GET /v3/{project_id}/elb/security-policies/{security_policy_id}
// @API ELB PUT /v3/{project_id}/elb/security-policies/{security_policy_id}
// @API ELB DELETE /v3/{project_id}/elb/security-policies/{security_policy_id}
func ResourceSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityPoliciesV3Create,
		UpdateContext: resourceSecurityPoliciesV3Update,
		ReadContext:   resourceSecurityPoliciesV3Read,
		DeleteContext: resourceSecurityPoliciesV3Delete,
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
			"protocols": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `Specifies the TSL protocol list which the security policy select.`,
			},
			"ciphers": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `Specifies the cipher suite list of the security policy.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the ELB security policy name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the description of the ELB security policy`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the enterprise project ID to which the Enterprise router belongs.`,
			},
			"listeners": {
				Type:        schema.TypeList,
				Elem:        SecurityPoliciesV3ListenerRefSchema(),
				Computed:    true,
				Description: `The listener which the security policy associated with.`,
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

func SecurityPoliciesV3ListenerRefSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSecurityPoliciesV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createSecurityPolicy: create an ELB security policy
	var (
		createSecurityPolicyHttpUrl = "v3/{project_id}/elb/security-policies"
		createSecurityPolicyProduct = "elbv3"
	)
	createSecurityPolicyClient, err := cfg.NewServiceClient(createSecurityPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecurityPolicies Client: %s", err)
	}

	createSecurityPolicyPath := createSecurityPolicyClient.Endpoint + createSecurityPolicyHttpUrl
	createSecurityPolicyPath = strings.ReplaceAll(createSecurityPolicyPath, "{project_id}", createSecurityPolicyClient.ProjectID)

	createSecurityPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createSecurityPolicyOpt.JSONBody = utils.RemoveNil(buildCreateSecurityPolicyBodyParams(d, cfg))
	createSecurityPolicyResp, err := createSecurityPolicyClient.Request("POST", createSecurityPolicyPath, &createSecurityPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating SecurityPolicies: %s", err)
	}

	createSecurityPolicyRespBody, err := utils.FlattenResponse(createSecurityPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	policyId := utils.PathSearch("security_policy.id", createSecurityPolicyRespBody, "").(string)
	if policyId == "" {
		return diag.Errorf("unable to find the security policy ID from the API response")
	}
	d.SetId(policyId)

	return resourceSecurityPoliciesV3Read(ctx, d, meta)
}

func buildCreateSecurityPolicyBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"security_policy": buildCreateSecurityPolicyChildBodyParams(d, cfg),
	}
}

func buildCreateSecurityPolicyChildBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"protocols":             utils.ValueIgnoreEmpty(d.Get("protocols")),
		"ciphers":               utils.ValueIgnoreEmpty(d.Get("ciphers")),
	}
}

func resourceSecurityPoliciesV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteSecurityPolicy: missing operation notes
	var (
		deleteSecurityPolicyHttpUrl = "v3/{project_id}/elb/security-policies/{security_policy_id}"
		deleteSecurityPolicyProduct = "elbv3"
	)
	deleteSecurityPolicyClient, err := cfg.NewServiceClient(deleteSecurityPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecurityPolicies Client: %s", err)
	}

	deleteSecurityPolicyPath := deleteSecurityPolicyClient.Endpoint + deleteSecurityPolicyHttpUrl
	deleteSecurityPolicyPath = strings.ReplaceAll(deleteSecurityPolicyPath, "{project_id}", deleteSecurityPolicyClient.ProjectID)
	deleteSecurityPolicyPath = strings.ReplaceAll(deleteSecurityPolicyPath, "{security_policy_id}", d.Id())

	deleteSecurityPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteSecurityPolicyClient.Request("DELETE", deleteSecurityPolicyPath, &deleteSecurityPolicyOpt)
	if err != nil {
		return diag.Errorf("error deleting SecurityPolicies: %s", err)
	}

	return nil
}

func resourceSecurityPoliciesV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getSecurityPolicy: Query the ELB security policy
	var (
		getSecurityPolicyHttpUrl = "v3/{project_id}/elb/security-policies/{security_policy_id}"
		getSecurityPolicyProduct = "elbv3"
	)
	getSecurityPolicyClient, err := cfg.NewServiceClient(getSecurityPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecurityPolicies Client: %s", err)
	}

	getSecurityPolicyPath := getSecurityPolicyClient.Endpoint + getSecurityPolicyHttpUrl
	getSecurityPolicyPath = strings.ReplaceAll(getSecurityPolicyPath, "{project_id}", getSecurityPolicyClient.ProjectID)
	getSecurityPolicyPath = strings.ReplaceAll(getSecurityPolicyPath, "{security_policy_id}", d.Id())

	getSecurityPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSecurityPolicyResp, err := getSecurityPolicyClient.Request("GET", getSecurityPolicyPath, &getSecurityPolicyOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SecurityPolicies")
	}

	getSecurityPolicyRespBody, err := utils.FlattenResponse(getSecurityPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("security_policy.name", getSecurityPolicyRespBody, nil)),
		d.Set("description", utils.PathSearch("security_policy.description", getSecurityPolicyRespBody, nil)),
		d.Set("protocols", utils.PathSearch("security_policy.protocols", getSecurityPolicyRespBody, nil)),
		d.Set("ciphers", utils.PathSearch("security_policy.ciphers", getSecurityPolicyRespBody, nil)),
		d.Set("listeners", flattenGetSecurityPolicyResponseBodyListenerRef(getSecurityPolicyRespBody)),
		d.Set("created_at", utils.PathSearch("security_policy.created_at", getSecurityPolicyRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("security_policy.updated_at", getSecurityPolicyRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetSecurityPolicyResponseBodyListenerRef(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("security_policy.listeners", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id": utils.PathSearch("id", v, nil),
		})
	}
	return rst
}

func resourceSecurityPoliciesV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateSecurityPolicyHasChanges := []string{
		"name",
		"description",
		"protocols",
		"ciphers",
	}

	if d.HasChanges(updateSecurityPolicyHasChanges...) {
		// updateSecurityPolicy: update the ELB security policy
		var (
			updateSecurityPolicyHttpUrl = "v3/{project_id}/elb/security-policies/{security_policy_id}"
			updateSecurityPolicyProduct = "elbv3"
		)
		updateSecurityPolicyClient, err := cfg.NewServiceClient(updateSecurityPolicyProduct, region)
		if err != nil {
			return diag.Errorf("error creating SecurityPolicies Client: %s", err)
		}

		updateSecurityPolicyPath := updateSecurityPolicyClient.Endpoint + updateSecurityPolicyHttpUrl
		updateSecurityPolicyPath = strings.ReplaceAll(updateSecurityPolicyPath, "{project_id}", updateSecurityPolicyClient.ProjectID)
		updateSecurityPolicyPath = strings.ReplaceAll(updateSecurityPolicyPath, "{security_policy_id}", d.Id())

		updateSecurityPolicyOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateSecurityPolicyOpt.JSONBody = utils.RemoveNil(buildUpdateSecurityPolicyBodyParams(d, cfg))
		_, err = updateSecurityPolicyClient.Request("PUT", updateSecurityPolicyPath, &updateSecurityPolicyOpt)
		if err != nil {
			return diag.Errorf("error updating SecurityPolicies: %s", err)
		}
	}
	return resourceSecurityPoliciesV3Read(ctx, d, meta)
}

func buildUpdateSecurityPolicyBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"security_policy": buildUpdateSecurityPolicyChildBodyParams(d, cfg),
	}
}

func buildUpdateSecurityPolicyChildBodyParams(d *schema.ResourceData, _ *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"protocols":   utils.ValueIgnoreEmpty(d.Get("protocols")),
		"ciphers":     utils.ValueIgnoreEmpty(d.Get("ciphers")),
	}
}
