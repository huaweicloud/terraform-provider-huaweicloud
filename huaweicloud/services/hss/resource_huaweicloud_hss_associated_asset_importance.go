package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS PUT /v5/{project_id}/host-management/asset-value/associate
func ResourceAssociatedAssetImportance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssociatedAssetImportanceCreate,
		ReadContext:   resourceAssociatedAssetImportanceRead,
		UpdateContext: resourceAssociatedAssetImportanceUpdate,
		DeleteContext: resourceAssociatedAssetImportanceDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"asset_value",
			"host_id_list",
			"enterprise_project_id",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"asset_value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_id_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
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

func buildAssociatedAssetImportanceQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}

	return ""
}

func buildAssociatedAssetImportanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"asset_value":  d.Get("asset_value"),
		"host_id_list": utils.ExpandToStringList(d.Get("host_id_list").([]interface{})),
	}

	return bodyParams
}

func resourceAssociatedAssetImportanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/host-management/asset-value/associate"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAssociatedAssetImportanceQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAssociatedAssetImportanceBodyParams(d),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error operating HSS associated asset importance: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return resourceAssociatedAssetImportanceRead(ctx, d, meta)
}

func resourceAssociatedAssetImportanceRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceAssociatedAssetImportanceUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceAssociatedAssetImportanceDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to operation HSS associated asset importance.
    Deleting this resource will not clear the corresponding request record, but will only remove the resource
    information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
