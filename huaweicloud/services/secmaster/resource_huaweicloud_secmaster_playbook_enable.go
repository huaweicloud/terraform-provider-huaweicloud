package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableParamsPlaybookEnable = []string{"workspace_id", "playbook_name", "playbook_id", "active_version_id"}

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/{playbook_id}/versions
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{id}
func ResourcePlaybookEnable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlaybookEnableCreate,
		UpdateContext: resourcePlaybookEnableRead,
		ReadContext:   resourcePlaybookEnableUpdate,
		DeleteContext: resourcePlaybookEnableDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsPlaybookEnable),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the workspace to which the playbook version belongs.`,
			},
			"playbook_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the playbook ID.`,
			},
			"playbook_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the playbook name.`,
			},
			"active_version_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the actived playbook version ID.`,
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

func resourcePlaybookEnableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	bodyParams := map[string]interface{}{
		"name":              d.Get("playbook_name"),
		"enabled":           true,
		"active_version_id": d.Get("active_version_id"),
	}

	playbookId := d.Get("playbook_id").(string)
	// updatePlaybookVersion: Update the configuration of SecMaster playbook
	err = updatePlaybook(client, d.Get("workspace_id").(string), playbookId, bodyParams)
	if err != nil {
		return diag.Errorf("error enabling SecMaster playbook: %s", err)
	}

	d.SetId(playbookId)

	return resourcePlaybookEnableRead(ctx, d, meta)
}

func resourcePlaybookEnableRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	playbook, err := GetPlaybook(client, d.Get("workspace_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "code", "SecMaster.20010001"), "error retrieving SecMaster playbook")
	}
	enabled := utils.PathSearch("enabled", playbook, false).(bool)
	if !enabled {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("playbook_id", utils.PathSearch("id", playbook, nil)),
		d.Set("playbook_name", utils.PathSearch("name", playbook, nil)),
		d.Set("active_version_id", utils.PathSearch("version_id", playbook, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePlaybookEnableUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePlaybookEnableDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	bodyParams := map[string]interface{}{
		"name":    d.Get("playbook_name"),
		"enabled": false,
	}

	// updatePlaybookVersion: Update the configuration of SecMaster playbook
	err = updatePlaybook(client, d.Get("workspace_id").(string), d.Get("playbook_id").(string), bodyParams)
	if err != nil {
		return diag.Errorf("error disabling SecMaster playbook: %s", err)
	}

	return nil
}

func updatePlaybook(client *golangsdk.ServiceClient, workspaceId, id string, bodyParam interface{}) error {
	updatePlaybookHttpUrl := "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/{id}"
	updatePlaybookPath := client.Endpoint + updatePlaybookHttpUrl
	updatePlaybookPath = strings.ReplaceAll(updatePlaybookPath, "{project_id}", client.ProjectID)
	updatePlaybookPath = strings.ReplaceAll(updatePlaybookPath, "{workspace_id}", workspaceId)
	updatePlaybookPath = strings.ReplaceAll(updatePlaybookPath, "{id}", id)

	updatePlaybookOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         bodyParam,
	}

	_, err := client.Request("PUT", updatePlaybookPath, &updatePlaybookOpt)

	return err
}
