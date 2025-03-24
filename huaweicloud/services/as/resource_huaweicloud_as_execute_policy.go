package as

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API AS POST /autoscaling-api/v1/{project_id}/scaling_policy/{scaling_policy_id}/action
var nonUpdatableParams = []string{"scaling_policy_id"}

func ResourceExecutePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceExecutePolicyCreate,
		ReadContext:   resourceExecutePolicyRead,
		UpdateContext: resourceExecutePolicyUpdate,
		DeleteContext: resourceExecutePolicyDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scaling_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the AS group policy ID or AS bandwidth policy ID.`,
			},
		},
	}
}

func resourceExecutePolicyCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "autoscaling-api/v1/{project_id}/scaling_policy/{scaling_policy_id}/action"
		policyId = d.Get("scaling_policy_id").(string)
	)

	client, err := cfg.NewServiceClient("autoscaling", region)
	if err != nil {
		return diag.Errorf("error creating AS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{scaling_policy_id}", policyId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"action": "execute",
		},
		OkCodes: []int{
			200, 201, 204,
		},
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error executing the AS policy (%s): %s", policyId, err)
	}

	d.SetId(policyId)
	d.Set("region", region)

	return nil
}

func resourceExecutePolicyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceExecutePolicyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceExecutePolicyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
