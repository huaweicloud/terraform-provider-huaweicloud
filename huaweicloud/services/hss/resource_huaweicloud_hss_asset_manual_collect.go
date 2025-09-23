package hss

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

var assetManualCollectNonUpdatableParams = []string{"type", "host_id", "enterprise_project_id"}

// @API HSS POST /v5/{project_id}/asset/manual-collect/{type}
func ResourceAssetManualCollect() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssetManualCollectCreate,
		ReadContext:   resourceAssetManualCollectRead,
		UpdateContext: resourceAssetManualCollectUpdate,
		DeleteContext: resourceAssetManualCollectDelete,

		CustomizeDiff: config.FlexibleForceNew(assetManualCollectNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_id": {
				Type:     schema.TypeString,
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

func buildAssetManualCollectQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?host_id=%v", d.Get("host_id"))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func resourceAssetManualCollectCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		assectType = d.Get("type").(string)
		hostId     = d.Get("host_id").(string)
		epsId      = cfg.GetEnterpriseProjectID(d)
		httpUrl    = "v5/{project_id}/asset/manual-collect/{type}"
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{type}", assectType)
	requestPath += buildAssetManualCollectQueryParams(d, epsId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error collecting the host asset fingerprints: %s", err)
	}

	d.SetId(hostId)

	return nil
}

func resourceAssetManualCollectRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceAssetManualCollectUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceAssetManualCollectDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to collect asset fingerprint of a host. Deleting this resource
    will not clear the corresponding request record, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
