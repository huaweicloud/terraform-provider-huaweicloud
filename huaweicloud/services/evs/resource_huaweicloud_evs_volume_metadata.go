package evs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EVS PUT /v3/{project_id}/volumes/{volume_id}/metadata
// @API EVS POST /v3/{project_id}/volumes/{volume_id}/metadata
// @API EVS GET /v3/{project_id}/volumes/{volume_id}/metadata/{key}
// @API EVS DELETE  /v3/{project_id}/volumes/{volume_id}/metadata/{key}
func ResourceEVSVolumedMeta() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEVSVolumedMetaCreate,
		UpdateContext: resourceEVSVolumedMetaUpdate,
		ReadContext:   resourceEVSVolumedMetaRead,
		DeleteContext: resourceEVSVolumedMetaDelete,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"volume_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The disk ID.`,
			},
			"metadata": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: `The snapshot metadata, which is made up of key-value pairs.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The disk ID.`,
			},
		},
	}
}

func resourceEVSVolumedMetaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		JSONBody:         utils.RemoveNil(requestBody),
	}
	_, err = client.Request("POST", requestPath, &createVolumeMetadataOpt)
	if err != nil {
		return diag.Errorf("error creating EVS volume metadata: %s", err)
	}

	d.SetId(volumeID)

	return resourceEVSVolumedMetaRead(ctx, d, meta)
}

func resourceEVSVolumedMetaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if len(metadataInput) > 0 {
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
				return diag.Errorf("error retrieving metadata for key %s: %s", key, err)
			}

			respBody, err := utils.FlattenResponse(resp)
			if err != nil {
				return diag.Errorf("error formatting metadata for key %s: %s", key, err)
			}
			path := fmt.Sprintf("meta.%s", key)
			metadataResult[key] = utils.PathSearch(path, respBody, nil)
		}
		if len(metadataResult) == 0 {
			return diag.Errorf("all metadata keys are not found")
		}
	}

	mErr = multierror.Append(
		mErr,
		d.Set("metadata", metadataResult),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEVSVolumedMetaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.HasChange("metadata") {
		return nil
	}
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
	log.Printf("[DEBUG] metadata requestBody is %v", requestBody)
	updateVolumeMetadataOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody:         requestBody,
	}
	log.Printf("[DEBUG] updateVolumeMetadataOpt is %v", updateVolumeMetadataOpt)

	_, err = client.Request("PUT", requestPath, &updateVolumeMetadataOpt)
	if err != nil {
		return diag.Errorf("error updating EVS volume metadata: %s", err)
	}
	return resourceEVSVolumedMetaRead(ctx, d, meta)
}

func resourceEVSVolumedMetaDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			return diag.Errorf("error deleting metadata for key %s: %s", key, err)
		}
	}
	return nil
}
