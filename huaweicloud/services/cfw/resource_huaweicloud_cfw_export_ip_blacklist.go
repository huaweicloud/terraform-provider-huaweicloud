package cfw

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var exportIpBlacklistNonUpdatableParams = []string{"fw_instance_id", "name"}

// @API CFW POST /v1/{project_id}/ptf/ip-blacklist/export
func ResourceExportIpBlacklist() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceExportIpBlacklistCreate,
		ReadContext:   resourceExportIpBlacklistRead,
		UpdateContext: resourceExportIpBlacklistUpdate,
		DeleteContext: resourceExportIpBlacklistDelete,

		CustomizeDiff: config.FlexibleForceNew(exportIpBlacklistNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"fw_instance_id": {
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
			"data": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildExportIpBlacklistQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?fw_instance_id=%s&name=%s", d.Get("fw_instance_id"), d.Get("name"))
}

func resourceExportIpBlacklistCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/ptf/ip-blacklist/export"
	)

	client, err := cfg.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildExportIpBlacklistQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error exporting CFW IP blacklist: %s", err)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.Errorf("error reading export response body: %s", err)
	}

	if err := d.Set("data", string(bodyBytes)); err != nil {
		return diag.Errorf("error setting exported IP blacklist data field: %s", err)
	}

	d.SetId(d.Get("fw_instance_id").(string))

	return nil
}

func resourceExportIpBlacklistRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceExportIpBlacklistUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceExportIpBlacklistDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to export the IP blacklist. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
