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

var volumesBatchExpandNonUpdatableParams = []string{
	"volumes",
	"volumes.*.id",
	"volumes.*.new_size",
	"is_auto_pay",
}

// @API EVS POST /v5/{project_id}/volumes/batch-extend
func ResourceVolumesBatchExpand() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumesBatchExpandCreate,
		ReadContext:   resourceVolumesBatchExpandRead,
		UpdateContext: resourceVolumesBatchExpandUpdate,
		DeleteContext: resourceVolumesBatchExpandDelete,

		CustomizeDiff: config.FlexibleForceNew(volumesBatchExpandNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"volumes": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"new_size": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"is_auto_pay": {
				Type:     schema.TypeBool,
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

func buildVolumesBatchExpandBodyParams(d *schema.ResourceData) map[string]interface{} {
	volumes := d.Get("volumes").([]interface{})
	volumesList := make([]map[string]interface{}, len(volumes))

	for i, v := range volumes {
		volume := v.(map[string]interface{})
		volumesList[i] = map[string]interface{}{
			"id":       volume["id"].(string),
			"new_size": volume["new_size"].(int),
		}
	}

	bodyParams := map[string]interface{}{
		"bss_param": map[string]interface{}{
			"is_auto_pay": utils.ValueIgnoreEmpty(d.Get("is_auto_pay")),
		},
		"volumes": volumesList,
	}

	return bodyParams
}

func resourceVolumesBatchExpandCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/volumes/batch-extend"
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildVolumesBatchExpandBodyParams(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch expanding EVS volumes: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return nil
}

func resourceVolumesBatchExpandRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceVolumesBatchExpandUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceVolumesBatchExpandDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to expand volumes.
Deleting this resource will not reset the expanded volumes, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
