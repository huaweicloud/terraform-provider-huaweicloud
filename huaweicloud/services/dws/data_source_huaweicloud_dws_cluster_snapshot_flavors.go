package dws

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS GET /v2/{project_id}/snapshots/{snapshot_id}/flavors
func DataSourceClusterSnapshotFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterSnapshotFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the cluster flavors are located.`,
			},

			// Required parameters.
			"snapshot_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the snapshot to be queried.`,
			},

			// Attributes.
			"flavors": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        snapshotFlavorSchema(),
				Description: `The list of cluster flavors for the snapshot.`,
			},
		},
	}
}

func snapshotFlavorSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor ID.`,
			},
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor code.`,
			},
			"classify": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor type.`,
			},
			"scenario": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor scenario.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor version.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor status.`,
			},
			"default_capacity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The default capacity of the flavor.`,
			},
			"duplicate": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of replicas used by the flavor.`,
			},
			"default_node": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The default number of nodes.`,
			},
			"min_node": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The minimum number of nodes.`,
			},
			"max_node": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum number of nodes.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The underlying flavor ID.`,
			},
			"flavor_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The underlying flavor code.`,
			},
			"volume_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of disks.`,
			},
			"attribute": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        snapshotFlavorAttributeSchema(),
				Description: `The list of extended information.`,
			},
			"product_version_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        snapshotFlavorProductVersionSchema(),
				Description: `The list of product versions supported by the flavor.`,
			},
			"volume_used": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        snapshotFlavorVolumeUsedSchema(),
				Description: `The disk usage information of the snapshot source cluster.`,
			},
		},
	}
}

func snapshotFlavorAttributeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The extended information code.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The extended information value.`,
			},
		},
	}
}

func snapshotFlavorProductVersionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"min_cn": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The minimum number of CN nodes supported by this version.`,
			},
			"max_cn": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum number of CN nodes supported by this version.`,
			},
			"version_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of this version.`,
			},
			"datastore_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The datastore version name.`,
			},
		},
	}
}

func snapshotFlavorVolumeUsedSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The disk type.`,
			},
			"volume_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of disks.`,
			},
			"capacity": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The available storage capacity of a single node, in GB.`,
			},
			"volume_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The physical storage capacity of a single disk, in GB.`,
			},
		},
	}
}

func flattenSnapshotFlavorAttributes(items []interface{}) []interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"code":  utils.PathSearch("code", item, nil),
			"value": utils.PathSearch("value", item, nil),
		})
	}
	return result
}

func flattenSnapshotFlavorProductVersions(items []interface{}) []interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"min_cn":            utils.PathSearch("min_cn", item, nil),
			"max_cn":            utils.PathSearch("max_cn", item, nil),
			"version_type":      utils.PathSearch("version_type", item, nil),
			"datastore_version": utils.PathSearch("datastore_version", item, nil),
		})
	}
	return result
}

func flattenSnapshotFlavorVolumeUsed(volumeUsed map[string]interface{}) []interface{} {
	if len(volumeUsed) < 1 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"volume_type": utils.PathSearch("volume_type", volumeUsed, nil),
			"volume_num":  utils.PathSearch("volume_num", volumeUsed, nil),
			"capacity":    utils.PathSearch("capacity", volumeUsed, nil),
			"volume_size": utils.PathSearch("volume_size", volumeUsed, nil),
		},
	}
}

func flattenSnapshotFlavors(items []interface{}) []interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"id":               utils.PathSearch("id", item, nil),
			"code":             utils.PathSearch("code", item, nil),
			"classify":         utils.PathSearch("classify", item, nil),
			"scenario":         utils.PathSearch("scenario", item, nil),
			"version":          utils.PathSearch("version", item, nil),
			"status":           utils.PathSearch("status", item, nil),
			"default_capacity": utils.PathSearch("default_capacity", item, nil),
			"duplicate":        utils.PathSearch("duplicate", item, nil),
			"default_node":     utils.PathSearch("default_node", item, nil),
			"min_node":         utils.PathSearch("min_node", item, nil),
			"max_node":         utils.PathSearch("max_node", item, nil),
			"flavor_id":        utils.PathSearch("flavor_id", item, nil),
			"flavor_code":      utils.PathSearch("flavor_code", item, nil),
			"volume_num":       utils.PathSearch("volume_num", item, nil),
			"attribute": flattenSnapshotFlavorAttributes(utils.PathSearch("attribute", item,
				make([]interface{}, 0)).([]interface{})),
			"product_version_list": flattenSnapshotFlavorProductVersions(utils.PathSearch("product_version_list", item,
				make([]interface{}, 0)).([]interface{})),
			"volume_used": flattenSnapshotFlavorVolumeUsed(utils.PathSearch("volume_used", item,
				make(map[string]interface{})).(map[string]interface{})),
		})
	}
	return result
}

func dataSourceClusterSnapshotFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		snapshotId = d.Get("snapshot_id").(string)
		httpUrl    = "v2/{project_id}/snapshots/{snapshot_id}/flavors"
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{snapshot_id}", snapshotId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return diag.Errorf("error querying the snapshot flavors by snapshot ID (%s): %s", snapshotId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error flattening response of the snapshot flavors by snapshot ID (%s): %s", snapshotId, err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("flavors", flattenSnapshotFlavors(
			utils.PathSearch("flavors", respBody, make([]interface{}, 0)).([]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
