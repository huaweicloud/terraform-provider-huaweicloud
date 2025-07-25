package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var ignoreFailedPCCNonUpdatableParams = []string{
	"action", "host_ids", "operate_all", "enterprise_project_id",
}

// @API HSS POST /v5/{project_id}/baseline/password-complexity/action
func ResourceIgnoreFailedPCC() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIgnoreFailedPCCCreate,
		ReadContext:   resourceIgnoreFailedPCCRead,
		UpdateContext: resourceIgnoreFailedPCCUpdate,
		DeleteContext: resourceIgnoreFailedPCCDelete,

		CustomizeDiff: config.FlexibleForceNew(ignoreFailedPCCNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Specifies the region where the resource is located. If omitted, the provider-level region will be used.",
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the action type to perform.",
			},
			"operate_all": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether to perform the action on all hosts.",
			},
			"host_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the list of host IDs to perform the action on.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the enterprise project ID.",
			},
		},
	}
}

func buildIgnoreFailedPCCQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	rst := fmt.Sprintf("?action=%s", d.Get("action").(string))

	if epsID := cfg.GetEnterpriseProjectID(d); epsID != "" {
		return fmt.Sprintf("%s&enterprise_project_id=%s", rst, epsID)
	}

	return rst
}

func buildIgnoreFailedPCCBodyParams(d *schema.ResourceData) map[string]interface{} {
	rst := map[string]interface{}{
		// Customize its default value to false.
		"operate_all": d.Get("operate_all"),
	}

	if rawArray, ok := d.Get("host_ids").([]interface{}); ok && len(rawArray) > 0 {
		rst["host_ids"] = utils.ExpandToStringList(rawArray)
	}

	return rst
}

func resourceIgnoreFailedPCCCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		action  = d.Get("action").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/baseline/password-complexity/action"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildIgnoreFailedPCCQueryParams(d, cfg)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildIgnoreFailedPCCBodyParams(d),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error operating (%s) action to HSS ignore failed password complexity check: %s", action, err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(id)

	return resourceIgnoreFailedPCCRead(ctx, d, meta)
}

func resourceIgnoreFailedPCCRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceIgnoreFailedPCCUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceIgnoreFailedPCCDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to ignore failed password complexity check. Deleting
	this resource will not clear the corresponding ignore records, but will only remove the resource information from
	the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
