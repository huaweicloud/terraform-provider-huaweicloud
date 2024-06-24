package rms

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API Config POST /v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/policy-states/run-evaluation
func ResourcePolicyAssignmentEvaluate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyAssignmentEvaluateCreate,
		ReadContext:   resourcePolicyAssignmentEvaluateRead,
		DeleteContext: resourcePolicyAssignmentEvaluateDelete,

		Schema: map[string]*schema.Schema{
			"policy_assignment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the policy assignment to evaluate.`,
			},
		},
	}
}

func resourcePolicyAssignmentEvaluateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createPolicyAssignmentEvaluate: Create a RMS policy assignment evaluate.
	var (
		createPolicyAssignmentEvaluateHttpUrl = "v1/resource-manager/domains/{domain_id}/policy-assignments/" +
			"{policy_assignment_id}/policy-states/run-evaluation"
		createPolicyAssignmentEvaluateProduct = "rms"
	)
	createPolicyAssignmentEvaluateClient, err := cfg.NewServiceClient(createPolicyAssignmentEvaluateProduct, region)
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	createPolicyAssignmentEvaluatePath := createPolicyAssignmentEvaluateClient.Endpoint + createPolicyAssignmentEvaluateHttpUrl
	createPolicyAssignmentEvaluatePath = strings.ReplaceAll(createPolicyAssignmentEvaluatePath, "{domain_id}", cfg.DomainID)
	createPolicyAssignmentEvaluatePath = strings.ReplaceAll(
		createPolicyAssignmentEvaluatePath, "{policy_assignment_id}", d.Get("policy_assignment_id").(string))

	createPolicyAssignmentEvaluateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = createPolicyAssignmentEvaluateClient.Request("POST", createPolicyAssignmentEvaluatePath,
		&createPolicyAssignmentEvaluateOpt)
	if err != nil {
		return diag.Errorf("error creating RMS policy assignment evaluate: %s", err)
	}

	d.SetId(d.Get("policy_assignment_id").(string))

	return resourcePolicyAssignmentEvaluateRead(ctx, d, meta)
}

func resourcePolicyAssignmentEvaluateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePolicyAssignmentEvaluateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting policy assignment evaluate is not supported. The policy assignment evaluate is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
