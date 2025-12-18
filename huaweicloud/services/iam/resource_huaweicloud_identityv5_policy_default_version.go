package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var policyV5DefaultVersionNonUpdatableParams = []string{"policy_id", "version_id"}

// @API IAM POST /v5/policies/{policy_id}/versions/{version_id}/set-default
func ResourceIamV5PolicyDefaultVersion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIamV5PolicyDefaultVersionCreate,
		ReadContext:   resourceIamV5PolicyDefaultVersionRead,
		UpdateContext: resourceIamV5PolicyDefaultVersionUpdate,
		DeleteContext: resourceIamV5PolicyDefaultVersionDelete,

		CustomizeDiff: config.FlexibleForceNew(policyV5DefaultVersionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the identity policy.`,
			},
			"version_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The version ID of the identity policy to be set as default.`,
			},
		},
	}
}

func resourceIamV5PolicyDefaultVersionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		policyId  = d.Get("policy_id").(string)
		versionId = d.Get("version_id").(string)
		httpUrl   = "v5/policies/{policy_id}/versions/{version_id}/set-default"
	)

	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{policy_id}", policyId)
	createPath = strings.ReplaceAll(createPath, "{version_id}", versionId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error setting default version (%s) for identity policy (%s): %s", versionId, policyId, err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomId)

	return nil
}

func resourceIamV5PolicyDefaultVersionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIamV5PolicyDefaultVersionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIamV5PolicyDefaultVersionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to set the default version of an identity policy. Deleting this
	resource will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
