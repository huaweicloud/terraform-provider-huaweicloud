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

var gaussdbSpecifiedSessionsQueryAndKillNonUpdatableParams = []string{
	"instance_id",
	"node_id",
	"component_id",
	"session_ids",
}

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/kill-session
func ResourceGaussdbSpecifiedSessionsQueryAndKill() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussdbSpecifiedSessionsQueryAndKillCreate,
		UpdateContext: resourceGaussdbSpecifiedSessionsQueryAndKillUpdate,
		ReadContext:   resourceGaussdbSpecifiedSessionsQueryAndKillRead,
		DeleteContext: resourceGaussdbSpecifiedSessionsQueryAndKillDelete,

		CustomizeDiff: config.FlexibleForceNew(gaussdbSpecifiedSessionsQueryAndKillNonUpdatableParams),

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
			"session_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func resourceGaussdbSpecifiedSessionsQueryAndKillCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	httpUrl := "v3/{project_id}/instances/{instance_id}/kill-session"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussdbSpecifiedSessionsQueryAndKillBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error querying and killing GaussDB specified sessions: %s", err)
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
		d.Set("session_ids", utils.PathSearch("session_ids", createRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCreateGaussdbSpecifiedSessionsQueryAndKillBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"node_id":      d.Get("node_id"),
		"component_id": d.Get("component_id"),
		"session_ids":  d.Get("session_ids"),
		"instance_id":  d.Get("instance_id"),
	}
	return bodyParams
}

func resourceGaussdbSpecifiedSessionsQueryAndKillRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussdbSpecifiedSessionsQueryAndKillUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussdbSpecifiedSessionsQueryAndKillDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "This resource is a one-time action resource for killing GaussDB specified sessions. " +
		"Deleting this resource will not clear the corresponding request record, " +
		"but will only remove the resource information from the tfstate file."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
