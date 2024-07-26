package evs

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EVS POST /v2/{project_id}/os-volume-transfer/{transfer_id}/accept
func ResourceVolumeTransferAccepter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumeTransferAccepterCreate,
		ReadContext:   resourceVolumeTransferAccepterRead,
		DeleteContext: resourceVolumeTransferAccepterDelete,

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
				ForceNew: true,
			},
			"auth_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVolumeTransferAccepterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/os-volume-transfer/{transfer_id}/accept"
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	createTransferAccepterPath := client.Endpoint + httpUrl
	createTransferAccepterPath = strings.ReplaceAll(createTransferAccepterPath, "{project_id}", client.ProjectID)
	createTransferAccepterPath = strings.ReplaceAll(createTransferAccepterPath, "{transfer_id}",
		d.Get("transfer_id").(string))
	createTransferAccepterBodyParams := map[string]interface{}{
		"accept": map[string]interface{}{
			"auth_key": d.Get("auth_key"),
		},
	}
	createTransferAccepterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         createTransferAccepterBodyParams,
	}

	createResponse, err := client.Request("POST", createTransferAccepterPath, &createTransferAccepterOpt)
	if err != nil {
		return diag.Errorf("error creating EVS volume transfer accepter: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	// After successfully calling the API, it will also return the `name` and `volume_id` attributes,
	// but for one-time action resource, they are ignored here.
	resourceId := utils.PathSearch("transfer.id", createRespBody, "").(string)
	if resourceId == "" {
		return diag.Errorf("error creating EVS volume transfer accepter: ID is not found in API response")
	}

	d.SetId(resourceId)

	return resourceVolumeTransferAccepterRead(ctx, d, meta)
}

func resourceVolumeTransferAccepterRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceVolumeTransferAccepterDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is a one-time action resource.
	return nil
}
