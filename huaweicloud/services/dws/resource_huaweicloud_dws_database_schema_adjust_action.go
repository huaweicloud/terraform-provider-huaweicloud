package dws

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var dwsDatabaseSchemaAdjustActionNonUpdatableParams = []string{
	"cluster_id",
	"database",
	"schema",
	"perm_space",
}

// @API DWS PUT /v2/{project_id}/clusters/{cluster_id}/databases/{database_name}/schemas
func ResourceDatabaseSchemaAdjustAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseSchemaAdjustActionCreate,
		ReadContext:   resourceDatabaseSchemaAdjustActionRead,
		UpdateContext: resourceDatabaseSchemaAdjustActionUpdate,
		DeleteContext: resourceDatabaseSchemaAdjustActionDelete,

		CustomizeDiff: config.FlexibleForceNew(dwsDatabaseSchemaAdjustActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the DWS cluster is located.`,
			},

			// Required parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the DWS cluster.`,
			},
			"database": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the database.`,
			},
			"schema": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the schema.`,
			},
			"perm_space": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The space threshold of the schema.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildDatabaseSchemaAdjustActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"schema_name": d.Get("schema"),
		"perm_space":  d.Get("perm_space"),
	}
}

func adjustDatabaseSchemaSpace(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl      = "v2/{project_id}/clusters/{cluster_id}/databases/{database}/schemas"
		clusterId    = d.Get("cluster_id").(string)
		databaseName = d.Get("database").(string)
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{cluster_id}", clusterId)
	updatePath = strings.ReplaceAll(updatePath, "{database}", databaseName)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildDatabaseSchemaAdjustActionBodyParams(d),
	}

	resp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	retCode := int(utils.PathSearch("ret_code", respBody, float64(-1)).(float64))
	if retCode != 0 {
		return fmt.Errorf("unexpected ret_code (%d), response body: %s", retCode, respBody)
	}

	return nil
}

func resourceDatabaseSchemaAdjustActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	err = adjustDatabaseSchemaSpace(client, d)
	if err != nil {
		return diag.Errorf("error adjusting database schema space: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	return resourceDatabaseSchemaAdjustActionRead(ctx, d, meta)
}

func resourceDatabaseSchemaAdjustActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDatabaseSchemaAdjustActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDatabaseSchemaAdjustActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for updating the schema space limit. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
