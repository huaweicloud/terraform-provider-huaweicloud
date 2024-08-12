package iam

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

var policyNonUpdatableParams = []string{"name", "policy_document", "path", "description"}

// @API IAM POST /v5/policies
// @API IAM GET /v5/policies/{policy_id}
// @API IAM DELETE /v5/policies/{policy_id}
func ResourceIdentityPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityPolicyCreate,
		ReadContext:   resourceIdentityPolicyRead,
		UpdateContext: resourceIdentityPolicyUpdate,
		DeleteContext: resourceIdentityPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(policyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_document": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachment_count": {
				Type:     schema.TypeInt,
				Computed: true,
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

func resourceIdentityPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam_no_version", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	createPolicyHttpUrl := "v5/policies"
	createPolicyPath := client.Endpoint + createPolicyHttpUrl
	createPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreatePolicyBodyParams(d)),
	}
	createPolicyResp, err := client.Request("POST", createPolicyPath, &createPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating IAM identity policy: %s", err)
	}
	createPolicyRespBody, err := utils.FlattenResponse(createPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("policy.policy_id", createPolicyRespBody, nil)
	if id == nil {
		return diag.Errorf("error creating IAM identity policy: policy_id is not found in API response")
	}
	d.SetId(id.(string))

	return resourceIdentityPolicyRead(ctx, d, meta)
}

func buildCreatePolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"policy_name":     d.Get("name"),
		"policy_document": d.Get("policy_document").(string),
		"path":            utils.ValueIgnoreEmpty(d.Get("path")),
		"description":     utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceIdentityPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam_no_version", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	getPolicyHttpUrl := "v5/policies/{policy_id}"
	getPolicyPath := client.Endpoint + getPolicyHttpUrl
	getPolicyPath = strings.ReplaceAll(getPolicyPath, "{policy_id}", d.Id())
	getPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPolicyResp, err := client.Request("GET", getPolicyPath, &getPolicyOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IAM identity policy")
	}
	getPolicyRespBody, err := utils.FlattenResponse(getPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	policy := utils.PathSearch("policy", getPolicyRespBody, nil)
	if policy == nil {
		return diag.Errorf("error getting IAM trust policy: policy is not found in API response")
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("policy_name", policy, nil)),
		d.Set("path", utils.PathSearch("path", policy, nil)),
		d.Set("description", utils.PathSearch("description", policy, nil)),
		d.Set("urn", utils.PathSearch("urn", policy, nil)),
		d.Set("policy_type", utils.PathSearch("policy_type", policy, nil)),
		d.Set("default_version_id", utils.PathSearch("default_version_id", policy, nil)),
		d.Set("attachment_count", utils.PathSearch("attachment_count", policy, nil)),
		d.Set("created_at", utils.PathSearch("created_at", policy, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", policy, nil)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM trust policy fields: %s", err)
	}

	return nil
}

func resourceIdentityPolicyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam_no_version", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	deletePolicyHttpUrl := "v5/policies/{policy_id}"
	deletePolicyPath := client.Endpoint + deletePolicyHttpUrl
	deletePolicyPath = strings.ReplaceAll(deletePolicyPath, "{policy_id}", d.Id())
	deletePolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePolicyPath, &deletePolicyOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IAM identity policy")
	}

	return nil
}
