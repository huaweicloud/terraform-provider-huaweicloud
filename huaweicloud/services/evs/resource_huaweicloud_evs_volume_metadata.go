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

var nonUpdatableVolumeMetadataParams = []string{
	"volume_id",
}

// @API EVS PUT /v3/{project_id}/volumes/{volume_id}/metadata
// @API EVS POST /v3/{project_id}/volumes/{volume_id}/metadata
// @API EVS GET /v3/{project_id}/volumes/{volume_id}/metadata/{key}
// @API EVS DELETE  /v3/{project_id}/volumes/{volume_id}/metadata/{key}
func ResourceVolumeMetadata() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumeMetadataCreate,
		UpdateContext: resourceVolumeMetadataUpdate,
		ReadContext:   resourceVolumeMetadataRead,
		DeleteContext: resourceVolumeMetadataDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableVolumeMetadataParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"volume_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The disk ID.`,
			},
			"metadata": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: `The volume metadata, which is made up of key-value pairs.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
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

func resourceVolumeMetadataCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/volumes/{volume_id}/metadata"
		product = "evs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	volumeID := d.Get("volume_id").(string)

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{volume_id}", volumeID)
	requestBody := map[string]interface{}{
		"metadata": d.Get("metadata").(map[string]interface{}),
	}
	createVolumeMetadataOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody:         requestBody,
	}
	_, err = client.Request("POST", requestPath, &createVolumeMetadataOpt)
	if err != nil {
		return diag.Errorf("error creating EVS volume metadata: %s", err)
	}

	d.SetId(volumeID)

	return resourceVolumeMetadataRead(ctx, d, meta)
}

func resourceVolumeMetadataRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr    *multierror.Error
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/volumes/{volume_id}/metadata/{key}"
		product = "evs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	metadataInput := d.Get("metadata").(map[string]interface{})
	metadataResult := make(map[string]interface{})

	for key := range metadataInput {
		requestPath := client.Endpoint + httpUrl
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
		requestPath = strings.ReplaceAll(requestPath, "{volume_id}", d.Get("volume_id").(string))
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		requestPath = strings.ReplaceAll(requestPath, "{key}", key)

		resp, err := client.Request("GET", requestPath, &requestOpt)
		if err != nil {
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(err, &errDefault404) {
				continue
			}
			return diag.Errorf("error retrieving EVS volume metadata: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}
		path := fmt.Sprintf("meta.%s", key)
		metadataResult[key] = utils.PathSearch(path, respBody, nil)
	}

	if len(metadataInput) > 0 && len(metadataResult) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("metadata", metadataResult),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceVolumeMetadataUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/volumes/{volume_id}/metadata"
		product = "evs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{volume_id}", d.Get("volume_id").(string))
	requestBody := map[string]interface{}{
		"metadata": d.Get("metadata").(map[string]interface{}),
	}

	updateVolumeMetadataOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody:         requestBody,
	}

	_, err = client.Request("PUT", requestPath, &updateVolumeMetadataOpt)
	if err != nil {
		return diag.Errorf("error updating EVS volume metadata: %s", err)
	}
	return resourceVolumeMetadataRead(ctx, d, meta)
}

func resourceVolumeMetadataDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/volumes/{volume_id}/metadata/{key}"
		product = "evs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	metadataInput := d.Get("metadata").(map[string]interface{})

	for key := range metadataInput {
		requestPath := client.Endpoint + httpUrl
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
		requestPath = strings.ReplaceAll(requestPath, "{volume_id}", d.Get("volume_id").(string))
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		requestPath = strings.ReplaceAll(requestPath, "{key}", key)

		_, err := client.Request("DELETE", requestPath, &requestOpt)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				continue
			}
			return diag.Errorf("error deleting EVS volume metadata for key (%s): %s", key, err)
		}
	}
	return nil
}
