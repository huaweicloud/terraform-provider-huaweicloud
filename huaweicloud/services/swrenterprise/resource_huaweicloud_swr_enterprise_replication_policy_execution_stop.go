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

var enterpriseReplicationPolicyExecutionStopNonUpdatableParams = []string{
	"instance_id", "execution_id",
}

// @API SWR PUT /v2/{project_id}/instances/{instance_id}/replication/executions/{execution_id}
func ResourceSwrEnterpriseReplicationPolicyExecutionStop() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterpriseReplicationPolicyExecutionStopCreate,
		UpdateContext: resourceSwrEnterpriseReplicationPolicyExecutionStopUpdate,
		ReadContext:   resourceSwrEnterpriseReplicationPolicyExecutionStopRead,
		DeleteContext: resourceSwrEnterpriseReplicationPolicyExecutionStopDelete,

		CustomizeDiff: config.FlexibleForceNew(enterpriseReplicationPolicyExecutionStopNonUpdatableParams),

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
			"execution_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the execution ID.`,
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

func resourceSwrEnterpriseReplicationPolicyExecutionStopCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	executionId := strconv.Itoa(d.Get("execution_id").(int))
	createHttpUrl := "v2/{project_id}/instances/{instance_id}/replication/executions/{execution_id}"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{execution_id}", executionId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error executing SWR iamge replication policy: %s", err)
	}

	d.SetId(instanceId + "/" + executionId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSwrEnterpriseReplicationPolicyExecutionStopRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrEnterpriseReplicationPolicyExecutionStopUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrEnterpriseReplicationPolicyExecutionStopDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting SWR enterprise replication policy execution stop resource is not supported." +
		" The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
