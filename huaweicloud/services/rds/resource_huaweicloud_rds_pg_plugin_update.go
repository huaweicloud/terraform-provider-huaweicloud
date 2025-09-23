package rds

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var pgPluginUpdateNonUpdatableParams = []string{"instance_id", "database_name", "extension_name"}

// @API RDS PUT /v3/{project_id}/instances/{instance_id}/extensions
func ResourceRdsPgPluginUpdate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsPgPluginUpdateCreate,
		ReadContext:   resourceRdsPgPluginUpdateRead,
		UpdateContext: resourceRdsPgPluginUpdateUpdate,
		DeleteContext: resourceRdsPgPluginUpdateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(pgPluginUpdateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the RDS instance.`,
			},
			"database_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the database name.`,
			},
			"extension_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the extension name.`,
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

func resourceRdsPgPluginUpdateCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/extensions"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateRdsPluginUpdateBodyParams(d))

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating RDS PostgreSQL plugin update: %s", err)
	}

	id := fmt.Sprintf("%s/%s/%s", d.Get("instance_id").(string), d.Get("database_name").(string),
		d.Get("extension_name").(string))
	d.SetId(id)

	return nil
}

func buildCreateRdsPluginUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"database_name":  d.Get("database_name"),
		"extension_name": d.Get("extension_name"),
	}
	return bodyParams
}

func resourceRdsPgPluginUpdateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsPgPluginUpdateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsPgPluginUpdateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting RDS PostgreSQL plugin update resource is not supported. The resource is only removed from the" +
		"state, the instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
