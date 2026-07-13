package gaussdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var gaussdbIdleSessionsQueryAndKillNonUpdatableParams = []string{
	"instance_id",
	"node_id",
	"component_id",
}

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/kill-free-session
func ResourceIdleSessionsQueryAndKill() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdleSessionsQueryAndKillCreate,
		UpdateContext: resourceIdleSessionsQueryAndKillUpdate,
		ReadContext:   resourceIdleSessionsQueryAndKillRead,
		DeleteContext: resourceIdleSessionsQueryAndKillDelete,

		CustomizeDiff: config.FlexibleForceNew(gaussdbIdleSessionsQueryAndKillNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"component_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"success": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceIdleSessionsQueryAndKillCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createPath := client.Endpoint + "v3/{project_id}/instances/{instance_id}/kill-free-session"
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateIdleSessionsQueryAndKillBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error querying and killing idle sessions: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	nodeId := d.Get("node_id").(string)
	componentId := d.Get("component_id").(string)
	resourceId := fmt.Sprintf("%s/%s/%s", instanceId, nodeId, componentId)
	d.SetId(resourceId)

	mErr := multierror.Append(nil,
		d.Set("success", utils.PathSearch("success", createRespBody, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceIdleSessionsQueryAndKillRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdleSessionsQueryAndKillUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdleSessionsQueryAndKillDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting kill idle sessions resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func buildCreateIdleSessionsQueryAndKillBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"node_id":      d.Get("node_id"),
		"component_id": d.Get("component_id"),
	}
}
