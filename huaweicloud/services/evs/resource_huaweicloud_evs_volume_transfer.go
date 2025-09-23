package evs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EVS POST /v2/{project_id}/os-volume-transfer
// @API EVS GET /v2/{project_id}/os-volume-transfer/{transfer_id}
// @API EVS DELETE /v2/{project_id}/os-volume-transfer/{transfer_id}
func ResourceVolumeTransfer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumeTransferCreate,
		ReadContext:   resourceVolumeTransferRead,
		DeleteContext: resourceVolumeTransferDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auth_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVolumeTransferCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                         = meta.(*config.Config)
		region                      = cfg.GetRegion(d)
		createVolumeTransferHttpUrl = "v2/{project_id}/os-volume-transfer"
		product                     = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	createVolumeTransferPath := client.Endpoint + createVolumeTransferHttpUrl
	createVolumeTransferPath = strings.ReplaceAll(createVolumeTransferPath, "{project_id}", client.ProjectID)
	createVolumeTransferBodyParams := map[string]interface{}{
		"transfer": map[string]interface{}{
			"name":      d.Get("name"),
			"volume_id": d.Get("volume_id"),
		},
	}
	createVolumeTransferOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         createVolumeTransferBodyParams,
	}

	createResponse, err := client.Request("POST", createVolumeTransferPath, &createVolumeTransferOpt)
	if err != nil {
		return diag.Errorf("error creating EVS volume transfer: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	transferId := utils.PathSearch("transfer.id", createRespBody, "").(string)
	if transferId == "" {
		return diag.Errorf("error creating EVS volume transfer: ID is not found in API response")
	}

	d.SetId(transferId)
	// The `auth_key` field is an important attribute and will only be returned when calling the create API.
	// So it is necessary to make a no null judgment and set value here.
	authKey := utils.PathSearch("transfer.auth_key", createRespBody, "").(string)
	if authKey == "" {
		return diag.Errorf("error creating EVS volume transfer: auth_key is not found in API response")
	}

	if err := d.Set("auth_key", authKey); err != nil {
		return diag.Errorf("error setting attribute `auth_key` after creating EVS volume transfer: %s", err)
	}

	return resourceVolumeTransferRead(ctx, d, meta)
}

func resourceVolumeTransferRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr                     *multierror.Error
		cfg                      = meta.(*config.Config)
		region                   = cfg.GetRegion(d)
		getVolumeTransferHttpUrl = "v2/{project_id}/os-volume-transfer/{transfer_id}"
		product                  = "evs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	getVolumeTransferPath := client.Endpoint + getVolumeTransferHttpUrl
	getVolumeTransferPath = strings.ReplaceAll(getVolumeTransferPath, "{project_id}", client.ProjectID)
	getVolumeTransferPath = strings.ReplaceAll(getVolumeTransferPath, "{transfer_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResponse, err := client.Request("GET", getVolumeTransferPath, &getOpt)
	if err != nil {
		// When the resource does not exist, calling the query API will return a `404` status code.
		return common.CheckDeletedDiag(d, err, "error retrieving EVS volume transfer")
	}

	getRespBody, err := utils.FlattenResponse(getResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	transfer := utils.PathSearch("transfer", getRespBody, nil)
	createdAtStr := utils.PathSearch("created_at", transfer, "").(string)
	// The `created_at` field return format is **2024-07-23T11:19:06.469300**,
	// convert it into the standard RFC3339 format here.
	createdAtTimestamp := utils.ConvertTimeStrToNanoTimestamp(createdAtStr, "2006-01-02T15:04:05.999999")
	createdAt := utils.FormatTimeStampRFC3339(createdAtTimestamp/1000, false)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("volume_id", utils.PathSearch("volume_id", transfer, "").(string)),
		d.Set("name", utils.PathSearch("name", transfer, "").(string)),
		d.Set("created_at", createdAt),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceVolumeTransferDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                         = meta.(*config.Config)
		region                      = cfg.GetRegion(d)
		deleteVolumeTransferHttpUrl = "v2/{project_id}/os-volume-transfer/{transfer_id}"
		product                     = "evs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	deleteVolumeTransferPath := client.Endpoint + deleteVolumeTransferHttpUrl
	deleteVolumeTransferPath = strings.ReplaceAll(deleteVolumeTransferPath, "{project_id}", client.ProjectID)
	deleteVolumeTransferPath = strings.ReplaceAll(deleteVolumeTransferPath, "{transfer_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deleteVolumeTransferPath, &deleteOpt)
	if err != nil {
		// When the resource does not exist, calling the delete API will return a `404` status code.
		return common.CheckDeletedDiag(d, err, "error deleting EVS volume transfer")
	}

	return nil
}
