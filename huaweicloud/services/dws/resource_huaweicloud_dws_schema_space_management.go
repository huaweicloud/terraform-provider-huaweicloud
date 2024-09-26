package dws

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API APIG PUT /v2/{project_id}/clusters/{cluster_id}/databases/{database_name}/schemas
func ResourceSchemaSpaceManagement() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSchemaSpaceManagementCreate,
		ReadContext:   resourceSchemaSpaceManagementRead,
		DeleteContext: resourceSchemaSpaceManagementDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the DWS cluster ID.",
			},
			"database_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the database name to which the schema space management belongs.",
			},
			"schema_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the name of the schema.",
			},
			"space_limit": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies space limit of the schema.",
			},
		},
	}
}

func resourceSchemaSpaceManagementCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		httpUrl      = "v2/{project_id}/clusters/{cluster_id}/databases/{database_name}/schemas"
		clusterId    = d.Get("cluster_id").(string)
		databaseName = d.Get("database_name").(string)
	)
	client, err := cfg.NewServiceClient("dws", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{cluster_id}", clusterId)
	path = strings.ReplaceAll(path, "{database_name}", databaseName)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"schema_name": d.Get("schema_name"),
			"perm_space":  d.Get("space_limit"),
		},
	}

	_, err = client.Request("PUT", path, &createOpt)
	if err != nil {
		return diag.Errorf("error modifying schema space limit for database (%s): %s", databaseName, err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	return resourceSchemaSpaceManagementRead(ctx, d, meta)
}

func resourceSchemaSpaceManagementRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSchemaSpaceManagementDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for modifying schema space limit. Deleting this resource will
	not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
