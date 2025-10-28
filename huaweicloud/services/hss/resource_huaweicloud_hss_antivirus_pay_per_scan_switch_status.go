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

var antivirusPayPerScanSwitchStatusNonUpdatableParams = []string{
	"enabled",
	"enterprise_project_id",
}

// @API HSS PUT /v5/{project_id}/antivirus/pay-per-scan
func ResourceAntivirusPayPerScanSwitchStatus() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAntivirusPayPerScanSwitchStatusCreate,
		ReadContext:   resourceAntivirusPayPerScanSwitchStatusRead,
		UpdateContext: resourceAntivirusPayPerScanSwitchStatusUpdate,
		DeleteContext: resourceAntivirusPayPerScanSwitchStatusDelete,

		CustomizeDiff: config.FlexibleForceNew(antivirusPayPerScanSwitchStatusNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
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

func buildCreateAntivirusPayPerScanSwitchStatusQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func buildCreateAntivirusPayPerScanSwitchStatusBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"enabled": d.Get("enabled"),
	}
}

func resourceAntivirusPayPerScanSwitchStatusCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/antivirus/pay-per-scan"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildCreateAntivirusPayPerScanSwitchStatusQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateAntivirusPayPerScanSwitchStatusBodyParams(d),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error update HSS antivirus pay-per-scan switch status: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return resourceAntivirusPayPerScanSwitchStatusRead(ctx, d, meta)
}

func resourceAntivirusPayPerScanSwitchStatusRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceAntivirusPayPerScanSwitchStatusUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceAntivirusPayPerScanSwitchStatusDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to update HSS antivirus pay-per-scan switch status.
    Deleting this resource will not clear the corresponding request record, but will only remove the resource
    information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
