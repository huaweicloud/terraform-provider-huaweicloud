package eip

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var eipBatchAttachShareBandwidthNonUpdatableParams = []string{
	"publicips",
	"publicips.*.bandwidth_id",
	"publicips.*.publicip_id",
}

// @API EIP POST /v3/{project_id}/eip/publicips/attach-share-bandwidth
func ResourceEipBatchAttachShareBandwidth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEipBatchAttachShareBandwidthCreate,
		ReadContext:   resourceEipBatchAttachShareBandwidthRead,
		UpdateContext: resourceEipBatchAttachShareBandwidthUpdate,
		DeleteContext: resourceEipBatchAttachShareBandwidthDelete,

		CustomizeDiff: config.FlexibleForceNew(eipBatchAttachShareBandwidthNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"publicips": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bandwidth_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"publicip_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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

func resourceEipBatchAttachShareBandwidthCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/eip/publicips/attach-share-bandwidth"
		product = "vpc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VPC EIP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
		JSONBody:         buildBatchAttachShareBandwidthBodyParams(d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error attaching EIPs to share bandwidth: %s", err)
	}

	_, err = utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildBatchAttachShareBandwidthBodyParams(d *schema.ResourceData) map[string]interface{} {
	publicipsRaw := d.Get("publicips").([]interface{})

	publicips := make([]map[string]interface{}, 0, len(publicipsRaw))
	for _, item := range publicipsRaw {
		publicip := item.(map[string]interface{})
		publicips = append(publicips, map[string]interface{}{
			"bandwidth_id": publicip["bandwidth_id"],
			"publicip_id":  publicip["publicip_id"],
		})
	}

	bodyParams := map[string]interface{}{
		"publicips": publicips,
	}

	return bodyParams
}

func resourceEipBatchAttachShareBandwidthRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Read()' method because resource is a one-time action resource.
	return nil
}

func resourceEipBatchAttachShareBandwidthUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Update()' method because resource is a one-time action resource.
	return nil
}

func resourceEipBatchAttachShareBandwidthDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to attach EIPs to shared bandwidth. 
Deleting this resource will not detach the EIPs from the shared bandwidth, but will only remove 
the resource information from the Terraform state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
