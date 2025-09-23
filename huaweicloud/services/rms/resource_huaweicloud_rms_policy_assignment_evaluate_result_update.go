package rms

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var policyAssignmentEvaluateResultUpdateNonUpdatableParams = []string{"policy_assignment_id", "trigger_type",
	"compliance_state", "evaluation_time", "evaluation_hash", "policy_resource", "policy_resource.*.resource_id",
	"policy_resource.*.resource_name", "policy_resource.*.resource_provider", "policy_resource.*.resource_type",
	"policy_resource.*.region_id", "policy_resource.*.domain_id", "policy_assignment_name",
}

// @API Config PUT /v1/resource-manager/domains/{domain_id}/policy-states
func ResourcePolicyAssignmentEvaluateResultUpdate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyAssignmentEvaluateResultUpdateCreate,
		ReadContext:   resourcePolicyAssignmentEvaluateResultUpdateRead,
		UpdateContext: resourcePolicyAssignmentEvaluateResultUpdateUpdate,
		DeleteContext: resourcePolicyAssignmentEvaluateResultUpdateDelete,

		CustomizeDiff: config.FlexibleForceNew(policyAssignmentEvaluateResultUpdateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"policy_assignment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"trigger_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"compliance_state": {
				Type:     schema.TypeString,
				Required: true,
			},
			"evaluation_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"evaluation_hash": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_resource": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     policyAssignmentEvaluateResultUpdatePolicyResource(),
			},
			"policy_assignment_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func policyAssignmentEvaluateResultUpdatePolicyResource() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_provider": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	return &sc
}

func resourcePolicyAssignmentEvaluateResultUpdateCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/resource-manager/domains/{domain_id}/policy-states"
		product = "rms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Config client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{domain_id}", cfg.DomainID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Config policy assignment evaluate result update: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId)

	return nil
}

func resourcePolicyAssignmentEvaluateResultUpdateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePolicyAssignmentEvaluateResultUpdateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePolicyAssignmentEvaluateResultUpdateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting Config policy assignment evaluate result update is not supported. The resource is only removed" +
		" from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
