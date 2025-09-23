package evs

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v3TransferAccepterNonUpdatableParams = []string{"transfer_id", "auth_key"}

// @API EVS POST /v3/{project_id}/os-volume-transfer/{transfer_id}/accept
func ResourceV3VolumeTransferAccepter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3VolumeTransferAccepterCreate,
		ReadContext:   resourceV3VolumeTransferAccepterRead,
		UpdateContext: resourceV3VolumeTransferAccepterUpdate,
		DeleteContext: resourceV3VolumeTransferAccepterDelete,

		CustomizeDiff: config.FlexibleForceNew(v3TransferAccepterNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"transfer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auth_key": {
				Type:     schema.TypeString,
				Required: true,
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

func buildV3VolumeTransferAccepterBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"accept": map[string]interface{}{
			"auth_key": d.Get("auth_key"),
		},
	}
}

func resourceV3VolumeTransferAccepterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/os-volume-transfer/{transfer_id}/accept"
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{transfer_id}", d.Get("transfer_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildV3VolumeTransferAccepterBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating EVS v3 volume transfer accepter: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceId := utils.PathSearch("transfer.id", respBody, "").(string)
	if resourceId == "" {
		return diag.Errorf("error creating EVS v3 volume transfer accepter: ID is not found in API response")
	}

	d.SetId(resourceId)

	return resourceV3VolumeTransferAccepterRead(ctx, d, meta)
}

func resourceV3VolumeTransferAccepterRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceV3VolumeTransferAccepterUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceV3VolumeTransferAccepterDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is a one-time action resource.
	return nil
}
