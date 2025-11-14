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

var enterpriseReplicationPolicyExecuteNonUpdatableParams = []string{
	"instance_id", "policy_id",
}

// @API SWR POST /v2/{project_id}/instances/{instance_id}/replication/executions
func ResourceSwrEnterpriseReplicationPolicyExecute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterpriseReplicationPolicyExecuteCreate,
		UpdateContext: resourceSwrEnterpriseReplicationPolicyExecuteUpdate,
		ReadContext:   resourceSwrEnterpriseReplicationPolicyExecuteRead,
		DeleteContext: resourceSwrEnterpriseReplicationPolicyExecuteDelete,

		CustomizeDiff: config.FlexibleForceNew(enterpriseReplicationPolicyExecuteNonUpdatableParams),

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
			"policy_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the policy ID.`,
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
				Description: `Indicates the execution ID.`,
			},
		},
	}
}

func resourceSwrEnterpriseReplicationPolicyExecuteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createHttpUrl := "v2/{project_id}/instances/{instance_id}/replication/executions"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"policy_id": d.Get("policy_id"),
		},
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error executing SWR iamge replication policy: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := int(utils.PathSearch("id", createRespBody, float64(-1)).(float64))
	if id == -1 {
		return diag.Errorf("unable to find SWR instance replication policy execution ID from the API response")
	}

	d.SetId(instanceId + "/" + strconv.Itoa(id))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("execution_id", id),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSwrEnterpriseReplicationPolicyExecuteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrEnterpriseReplicationPolicyExecuteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrEnterpriseReplicationPolicyExecuteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting SWR enterprise replication policy execute resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
