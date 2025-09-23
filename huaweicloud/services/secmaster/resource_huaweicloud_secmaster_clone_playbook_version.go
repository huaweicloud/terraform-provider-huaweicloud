package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var clonePlaybookAndVersionNonUpdatableParams = []string{"workspace_id", "version_id", "name"}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/clone
func ResourceClonePlaybookAndVersion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClonePlaybookAndVersionCreate,
		ReadContext:   resourceClonePlaybookAndVersionRead,
		UpdateContext: resourceClonePlaybookAndVersionUpdate,
		DeleteContext: resourceClonePlaybookAndVersionDelete,

		CustomizeDiff: config.FlexibleForceNew(clonePlaybookAndVersionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
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

func buildClonePlaybookAndVersionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name": d.Get("name"),
	}

	return bodyParams
}

func resourceClonePlaybookAndVersionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/clone"
		workspaceId   = d.Get("workspace_id").(string)
		versionId     = d.Get("version_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)
	createPath = strings.ReplaceAll(createPath, "{version_id}", versionId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildClonePlaybookAndVersionBodyParams(d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error cloning playbook and playbook version: %s", err)
	}

	respBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	playbookId := utils.PathSearch("data.id", respBody, "").(string)
	if playbookId == "" {
		return diag.Errorf("unable to find the playbook ID from the API response")
	}

	d.SetId(playbookId)

	return nil
}

func resourceClonePlaybookAndVersionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceClonePlaybookAndVersionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceClonePlaybookAndVersionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
