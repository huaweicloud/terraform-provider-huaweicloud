package evs

import (
	"context"
	"errors"
	"fmt"
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

var snapshotMetadataNonUpdatableParams = []string{"snapshot_id"}

// @API EVS POST /v3/{project_id}/snapshots/{snapshot_id}/metadata
// @API EVS GET /v3/{project_id}/snapshots/{snapshot_id}/metadata/{key}
// @API EVS PUT /v3/{project_id}/snapshots/{snapshot_id}/metadata
// @API EVS DELETE /v3/{project_id}/snapshots/{snapshot_id}/metadata/{key}
func ResourceSnapshotMetadata() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSnapshotMetadataCreate,
		ReadContext:   resourceSnapshotMetadataRead,
		UpdateContext: resourceSnapshotMetadataUpdate,
		DeleteContext: resourceSnapshotMetadataDelete,

		CustomizeDiff: config.FlexibleForceNew(snapshotMetadataNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			// The `metadata` field can be set as an empty map.
			"metadata": {
				Type:     schema.TypeMap,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func buildCreateOrUpdateSnapshotMetadataBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"metadata": utils.ExpandToStringMap(d.Get("metadata").(map[string]interface{})),
	}
}

func resourceSnapshotMetadataCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		snapshotID = d.Get("snapshot_id").(string)
		product    = "evs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + "v3/{project_id}/snapshots/{snapshot_id}/metadata"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{snapshot_id}", d.Get("snapshot_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateOrUpdateSnapshotMetadataBodyParams(d),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating EVS snapshot metadata: %s", err)
	}

	d.SetId(snapshotID)

	return resourceSnapshotMetadataRead(ctx, d, meta)
}

func resourceSnapshotMetadataRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		product        = "evs"
		metadataResult = make(map[string]interface{})
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + "v3/{project_id}/snapshots/{snapshot_id}/metadata"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{snapshot_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	metadataInput := d.Get("metadata").(map[string]interface{})
	for key := range metadataInput {
		getPath := fmt.Sprintf("%s/%s", requestPath, key)
		getResp, err := client.Request("GET", getPath, &requestOpt)
		if err != nil {
			// When the `metadata` key does not exist, calling the query API will return a `404` status code,
			// we need to ignore this `404` error.
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(err, &errDefault404) {
				continue
			}

			return diag.Errorf("error retrieving EVS snapshot metadata: %s", err)
		}

		respBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		metaResp := utils.PathSearch("meta", respBody, make(map[string]interface{})).(map[string]interface{})
		for k, v := range metaResp {
			metadataResult[k] = v
		}
	}

	// When all key queries in `metadata` return `404`, it is considered that the resource does not exist.
	// Then execute `checkDeleted` logic.
	if len(metadataInput) > 0 && len(metadataResult) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("metadata", metadataResult),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSnapshotMetadataUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "evs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + "v3/{project_id}/snapshots/{snapshot_id}/metadata"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{snapshot_id}", d.Get("snapshot_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateOrUpdateSnapshotMetadataBodyParams(d),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating EVS snapshot metadata: %s", err)
	}

	return resourceSnapshotMetadataRead(ctx, d, meta)
}

func resourceSnapshotMetadataDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/snapshots/{snapshot_id}/metadata"
		product = "evs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{snapshot_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for key := range d.Get("metadata").(map[string]interface{}) {
		deletePath := fmt.Sprintf("%s/%s", requestPath, key)
		_, err = client.Request("DELETE", deletePath, &requestOpt)
		if err != nil {
			// When the `metadata` key does not exist, calling the delete API will return a `404` status code,
			// we need to ignore the `404` error because it no longer exists.
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(err, &errDefault404) {
				continue
			}

			return diag.Errorf("error deleting EVS snapshot metadata: %s", err)
		}
	}

	return nil
}
