package swrenterprise

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var enterpriseRetentionPolicyExecuteNonUpdatableParams = []string{
	"instance_id", "namespace_name", "policy_id", "dry_run",
}

// @API SWR POST /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/retention/policies/{policy_id}/executions
func ResourceSwrEnterpriseRetentionPolicyExecute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterpriseRetentionPolicyExecuteCreate,
		UpdateContext: resourceSwrEnterpriseRetentionPolicyExecuteUpdate,
		ReadContext:   resourceSwrEnterpriseRetentionPolicyExecuteRead,
		DeleteContext: resourceSwrEnterpriseRetentionPolicyExecuteDelete,

		CustomizeDiff: config.FlexibleForceNew(enterpriseRetentionPolicyExecuteNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the enterprise instance ID.`,
			},
			"namespace_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the namespace name.`,
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the policy ID.`,
			},
			"dry_run": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Specifies whether to dry run.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"execution_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Specifies the execution ID.`,
			},
		},
	}
}

func resourceSwrEnterpriseRetentionPolicyExecuteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	namespaceName := d.Get("namespace_name").(string)
	policyId := d.Get("policy_id").(string)

	createHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/retention/policies/{policy_id}/executions"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{namespace_name}", namespaceName)
	createPath = strings.ReplaceAll(createPath, "{policy_id}", policyId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"dry_run": d.Get("dry_run"),
		},
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error executing SWR iamge signature policy: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := int(utils.PathSearch("id", createRespBody, float64(-1)).(float64))
	if id == -1 {
		return diag.Errorf("unable to find SWR instance execution ID from the API response")
	}

	d.SetId(instanceId + "/" + namespaceName + "/" + strconv.Itoa(id))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("execution_id", id),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSwrEnterpriseRetentionPolicyExecuteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrEnterpriseRetentionPolicyExecuteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrEnterpriseRetentionPolicyExecuteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting SWR enterprise retention policy execute resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
