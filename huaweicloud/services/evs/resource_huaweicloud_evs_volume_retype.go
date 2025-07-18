package evs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var volumeRetypeNonUpdatableParams = []string{
	"volume_id",
	"is_auto_pay",
	"new_type",
	"iops",
	"throughput",
}

// @API EVS POST /v2/{project_id}/volumes/{volume_id}/retype
func ResourceVolumeRetype() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumeRetypeCreate,
		ReadContext:   resourceVolumeRetypeRead,
		UpdateContext: resourceVolumeRetypeUpdate,
		DeleteContext: resourceVolumeRetypeDelete,

		CustomizeDiff: config.FlexibleForceNew(volumeRetypeNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"new_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_auto_pay": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"iops": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"throughput": {
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

func buildVolumeRetypeBodyParams(d *schema.ResourceData) map[string]interface{} {
	osRetype := map[string]interface{}{}
	if v, ok := d.GetOk("new_type"); ok {
		osRetype["new_type"] = v
	}
	if v, ok := d.GetOk("iops"); ok {
		osRetype["iops"] = v
	}
	if v, ok := d.GetOk("throughput"); ok {
		osRetype["throughput"] = v
	}
	bodyParams := map[string]interface{}{
		"os-retype": osRetype,
	}
	bssParam := map[string]interface{}{}
	if v, ok := d.GetOk("is_auto_pay"); ok {
		bssParam["isAutoPay"] = v
	}
	if len(bssParam) > 0 {
		bodyParams["bssParam"] = bssParam
	}

	return bodyParams
}

func resourceVolumeRetypeCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "evs"
		httpUrl  = "v2/{project_id}/volumes/{volume_id}/retype"
		volumeID = d.Get("volume_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	url := client.Endpoint + httpUrl
	url = strings.ReplaceAll(url, "{project_id}", client.ProjectID)
	url = strings.ReplaceAll(url, "{volume_id}", volumeID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildVolumeRetypeBodyParams(d)),
	}

	_, err = client.Request("POST", url, &opt)
	if err != nil {
		return diag.Errorf("failed to retype volume: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return nil
}

func resourceVolumeRetypeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceVolumeRetypeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceVolumeRetypeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to retype the volume.  Deleting this resource will
not reset the volume type, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
