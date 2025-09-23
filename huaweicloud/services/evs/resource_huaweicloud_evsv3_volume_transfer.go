package evs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v3TransferNonUpdatableParams = []string{"volume_id", "name"}

// @API EVS POST /v3/{project_id}/os-volume-transfer
// @API EVS GET /v3/{project_id}/os-volume-transfer/{transfer_id}
// @API EVS DELETE /v3/{project_id}/os-volume-transfer/{transfer_id}
func ResourceV3VolumeTransfer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3VolumeTransferCreate,
		ReadContext:   resourceV3VolumeTransferRead,
		UpdateContext: resourceV3VolumeTransferUpdate,
		DeleteContext: resourceV3VolumeTransferDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(v3TransferNonUpdatableParams),

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
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     v3TransferLinksComputeSchema(),
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

func v3TransferLinksComputeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"href": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rel": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateV3TransferBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"transfer": map[string]interface{}{
			"name":      d.Get("name"),
			"volume_id": d.Get("volume_id"),
		},
	}
}

func resourceV3VolumeTransferCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/os-volume-transfer"
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
		JSONBody:         buildCreateV3TransferBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating EVS v3 volume transfer: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	transferId := utils.PathSearch("transfer.id", respBody, "").(string)
	if transferId == "" {
		return diag.Errorf("error creating EVS v3 volume transfer: ID is not found in API response")
	}

	d.SetId(transferId)
	// The `auth_key` field is an important attribute and will only be returned when calling the create API.
	// So it is necessary to make a no null judgment and set value here.
	authKey := utils.PathSearch("transfer.auth_key", respBody, "").(string)
	if authKey == "" {
		return diag.Errorf("error creating EVS v3 volume transfer: auth_key is not found in API response")
	}

	if err := d.Set("auth_key", authKey); err != nil {
		return diag.Errorf("error setting attribute `auth_key` after creating EVS v3 volume transfer: %s", err)
	}

	return resourceV3VolumeTransferRead(ctx, d, meta)
}

func flattenV3TransferLinksAttribute(respBody interface{}) []interface{} {
	links := utils.PathSearch("transfer.links", respBody, make([]interface{}, 0)).([]interface{})
	if len(links) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(links))
	for _, v := range links {
		rst = append(rst, map[string]interface{}{
			"href": utils.PathSearch("href", v, nil),
			"rel":  utils.PathSearch("rel", v, nil),
		})
	}

	return rst
}

func resourceV3VolumeTransferRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr    *multierror.Error
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/os-volume-transfer/{transfer_id}"
		product = "evs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{transfer_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		// When the resource does not exist, calling the query API will return a `404` status code.
		return common.CheckDeletedDiag(d, err, "error retrieving EVS v3 volume transfer")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("volume_id", utils.PathSearch("transfer.volume_id", respBody, nil)),
		d.Set("name", utils.PathSearch("transfer.name", respBody, nil)),
		d.Set("links", flattenV3TransferLinksAttribute(respBody)),
		d.Set("created_at", utils.PathSearch("transfer.created_at", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceV3VolumeTransferUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV3VolumeTransferDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/os-volume-transfer/{transfer_id}"
		product = "evs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{transfer_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		// When the resource does not exist, calling the delete API will return a `404` status code.
		return common.CheckDeletedDiag(d, err, "error deleting EVS v3 volume transfer")
	}

	return nil
}
