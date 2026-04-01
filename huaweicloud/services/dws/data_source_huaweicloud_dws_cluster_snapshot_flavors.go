package dws

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
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
				Description: `The region where the snapshot flavors are located.`,
			},
			"snapshot_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the snapshot.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the snapshot flavor.`,
			},
			"available_zones": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The availability zone code for restoration.`,
			},
			"fine_grained_restore": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether it is fine-grained backup restoration.`,
			},

			// Attributes
			"flavors": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of snapshot flavors that matched filter parameters.`,
				Elem:        flavorSchema(),
			},
		},
	}
}

func flavorSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the flavor.`,
			},
			"classify": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The classify of the flavor.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of the flavor.`,
			},
			"default_node": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The default node number of the flavor.`,
			},
			"max_node": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum node number of the flavor.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor ID of the flavor.`,
			},
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The code of the flavor.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the flavor.`,
			},
			"attributes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The attributes of the flavor.`,
				Elem:        attributeSchema(),
			},
			"min_node": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The minimum node number of the flavor.`,
			},
			"flavor_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor code of the flavor.`,
			},
			"product_versions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The product version list of the flavor.`,
				Elem:        productVersionSchema(),
			},
			"volume_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The volume number of the flavor.`,
			},
			"default_capacity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The default capacity of the flavor.`,
			},
			"scenario": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The scenario of the flavor.`,
			},
			"duplicate": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The duplicate of the flavor.`,
			},
			"volume_used": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The volume used information of the flavor.`,
				Elem:        volumeUsedSchema(),
			},
		},
	}
}

func attributeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The code of the attribute.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the attribute.`,
			},
		},
	}
}

func productVersionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datastore_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The datastore version of the product version.`,
			},
			"min_cn": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The minimum CN number of the product version.`,
			},
			"max_cn": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum CN number of the product version.`,
			},
			"version_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version type of the product version.`,
			},
		},
	}
}

func volumeUsedSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The volume type of the volume used.`,
			},
			"volume_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The volume number of the volume used.`,
			},
			"capacity": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The capacity of the volume used.`,
			},
			"volume_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The volume size of the volume used.`,
			},
		},
	}
}

func buildClusterSnapshotFlavorsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("type"); ok {
		res += fmt.Sprintf("&type=%v", v)
	}
	if v, ok := d.GetOk("available_zones"); ok {
		res += fmt.Sprintf("&az_code=%v", v)
	}
	if v, ok := d.GetOk("fine_grained_restore"); ok {
		res += fmt.Sprintf("&fine_grained_restore=%v", v)
	}

	if len(res) > 1 {
		return "?" + res[1:]
	}
	return res
}

func dataSourceClusterSnapshotFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	snapshotID := d.Get("snapshot_id").(string)
	httpUrl := "v2/{project_id}/snapshots/{snapshot_id}/flavors"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{snapshot_id}", snapshotID)
	listPath += buildClusterSnapshotFlavorsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return diag.Errorf("error retrieving snapshot flavors: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	flavors := utils.PathSearch("flavors", respBody, make([]interface{}, 0)).([]interface{})

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("flavors", flattenClusterSnapshotFlavors(flavors)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenClusterSnapshotFlavors(flavors []interface{}) []interface{} {
	if len(flavors) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(flavors))
	for _, flavor := range flavors {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("id", flavor, nil),
			"classify":     utils.PathSearch("classify", flavor, nil),
			"version":      utils.PathSearch("version", flavor, nil),
			"default_node": utils.PathSearch("default_node", flavor, nil),
			"max_node":     utils.PathSearch("max_node", flavor, nil),
			"flavor_id":    utils.PathSearch("flavor_id", flavor, nil),
			"code":         utils.PathSearch("code", flavor, nil),
			"status":       utils.PathSearch("status", flavor, nil),
			"attributes": flattenClusterSnapshotAttributes(
				utils.PathSearch("attribute", flavor, make([]interface{}, 0)).([]interface{})),
			"min_node":    utils.PathSearch("min_node", flavor, nil),
			"flavor_code": utils.PathSearch("flavor_code", flavor, nil),
			"product_versions": flattenClusterSnapshotProductVersions(
				utils.PathSearch("product_version_list", flavor, make([]interface{}, 0)).([]interface{})),
			"volume_num":       utils.PathSearch("volume_num", flavor, nil),
			"default_capacity": utils.PathSearch("default_capacity", flavor, nil),
			"scenario":         utils.PathSearch("scenario", flavor, nil),
			"duplicate":        utils.PathSearch("duplicate", flavor, nil),
			"volume_used":      flattenClusterSnapshotVolumeUsed(utils.PathSearch("volume_used", flavor, nil)),
		})
	}
	return result
}

func flattenClusterSnapshotAttributes(attributes []interface{}) []interface{} {
	if len(attributes) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(attributes))
	for _, attr := range attributes {
		result = append(result, map[string]interface{}{
			"code":  utils.PathSearch("code", attr, nil),
			"value": utils.PathSearch("value", attr, nil),
		})
	}
	return result
}

func flattenClusterSnapshotProductVersions(productVersions []interface{}) []interface{} {
	if len(productVersions) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(productVersions))
	for _, pv := range productVersions {
		result = append(result, map[string]interface{}{
			"datastore_version": utils.PathSearch("datastore_version", pv, nil),
			"min_cn":            utils.PathSearch("min_cn", pv, nil),
			"max_cn":            utils.PathSearch("max_cn", pv, nil),
			"version_type":      utils.PathSearch("version_type", pv, nil),
		})
	}
	return result
}

func flattenClusterSnapshotVolumeUsed(volumeUsed interface{}) []interface{} {
	if volumeUsed == nil {
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
