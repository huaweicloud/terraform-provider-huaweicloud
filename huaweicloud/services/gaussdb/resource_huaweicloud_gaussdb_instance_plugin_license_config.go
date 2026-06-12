package gaussdb

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var gaussdbPluginLicenseNonUpdatableParams = []string{"instance_id", "license_str"}

// @API GaussDB PUT /v3/{project_id}/instances/{instance_id}/kernel-plugin-license
func ResourceGaussDbPluginLicense() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginLicenseCreate,
		ReadContext:   resourcePluginLicenseRead,
		UpdateContext: resourcePluginLicenseUpdate,
		DeleteContext: resourcePluginLicenseDelete,

		CustomizeDiff: config.FlexibleForceNew(gaussdbPluginLicenseNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"license_str": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"enable_force_new": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "schema: Internal",
			},
		},
	}
}

func buildPluginLicenseBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"license_str": d.Get("license_str"),
	}
}

func resourcePluginLicenseCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/kernel-plugin-license"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", d.Get("instance_id").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildPluginLicenseBodyParams(d),
	}

	_, err = client.Request("PUT", requestPath, &opt)
	if err != nil {
		return diag.Errorf("error configuring GaussDB plugin license: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(resourceId)

	return nil
}

func resourcePluginLicenseRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePluginLicenseUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePluginLicenseDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB instance plugin license resource is not supported. The license resource is only removed " +
		"from the state, the GaussDB instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
