package evs

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

// @API EVS GET /v5/{project_id}/snapshots/detail
func DataSourceEvsv5Snapshots() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvsv5SnapshotsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshot_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshot_chain_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshot_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshots": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     evsv5SnapshotSchema(),
			},
		},
	}
}

func evsv5SnapshotSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cmk_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instant_access": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"retention_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instant_access_retention_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"incremental": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"snapshot_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"progress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encrypt_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_chains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     evsv5SnapshotChainSchema(),
			},
			"snapshot_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func evsv5SnapshotChainSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"capacity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceEvsv5SnapshotsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/snapshots/detail"
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	var allSnapshots []interface{}
	offset := 0
	for {
		pagedPath := fmt.Sprintf("%s?limit=1000&offset=%d%s", requestPath, offset, buildEvsv5SnapshotsQueryParams(d, cfg))
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		resp, err := client.Request("GET", pagedPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error querying EVS v5 snapshots: %s", err)
		}
		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}
		snapshots := utils.PathSearch("snapshots", respBody, []interface{}{}).([]interface{})
		if len(snapshots) == 0 {
			break
		}
		allSnapshots = append(allSnapshots, snapshots...)
		offset += len(snapshots)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("snapshots", flattenEvsv5Snapshots(allSnapshots)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildEvsv5SnapshotsQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	params := ""
	epsId := cfg.GetEnterpriseProjectID(d, "all_granted_eps")
	params += fmt.Sprintf("&enterprise_project_id=%v", epsId)
	if v, ok := d.GetOk("volume_id"); ok {
		params += fmt.Sprintf("&volume_id=%v", v)
	}
	if v, ok := d.GetOk("availability_zone"); ok {
		params += fmt.Sprintf("&availability_zone=%v", v)
	}
	if v, ok := d.GetOk("name"); ok {
		params += fmt.Sprintf("&name=%v", v)
	}
	if v, ok := d.GetOk("status"); ok {
		params += fmt.Sprintf("&status=%v", v)
	}
	if v, ok := d.GetOk("sort_key"); ok {
		params += fmt.Sprintf("&sort_key=%v", v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		params += fmt.Sprintf("&sort_dir=%v", v)
	}
	if v, ok := d.GetOk("id"); ok {
		params += fmt.Sprintf("&id=%v", v)
	}
	if v, ok := d.GetOk("ids"); ok {
		params += fmt.Sprintf("&ids=%v", v)
	}
	if v, ok := d.GetOk("snapshot_type"); ok {
		params += fmt.Sprintf("&snapshot_type=%v", v)
	}
	if v, ok := d.GetOk("tag_key"); ok {
		params += fmt.Sprintf("&tag_key=%v", v)
	}
	if v, ok := d.GetOk("tags"); ok {
		params += fmt.Sprintf("&tags=%v", v)
	}
	if v, ok := d.GetOk("snapshot_chain_id"); ok {
		params += fmt.Sprintf("&snapshot_chain_id=%v", v)
	}
	if v, ok := d.GetOk("snapshot_group_id"); ok {
		params += fmt.Sprintf("&snapshot_group_id=%v", v)
	}
	return params
}

func flattenEvsv5Snapshots(snapshots []interface{}) []map[string]interface{} {
	if len(snapshots) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(snapshots))
	for _, s := range snapshots {
		rst = append(rst, map[string]interface{}{
			"id":                          utils.PathSearch("id", s, nil),
			"name":                        utils.PathSearch("name", s, nil),
			"description":                 utils.PathSearch("description", s, nil),
			"created_at":                  utils.PathSearch("created_at", s, nil),
			"updated_at":                  utils.PathSearch("updated_at", s, nil),
			"volume_id":                   utils.PathSearch("volume_id", s, nil),
			"size":                        utils.PathSearch("size", s, nil),
			"status":                      utils.PathSearch("status", s, nil),
			"enterprise_project_id":       utils.PathSearch("enterprise_project_id", s, nil),
			"encrypted":                   utils.PathSearch("encrypted", s, nil),
			"cmk_id":                      utils.PathSearch("cmk_id", s, nil),
			"category":                    utils.PathSearch("category", s, nil),
			"availability_zone":           utils.PathSearch("availability_zone", s, nil),
			"tags":                        utils.PathSearch("tags", s, map[string]interface{}{}).(map[string]interface{}),
			"instant_access":              utils.PathSearch("instant_access", s, nil),
			"retention_at":                utils.PathSearch("retention_at", s, nil),
			"instant_access_retention_at": utils.PathSearch("instant_access_retention_at", s, nil),
			"incremental":                 utils.PathSearch("incremental", s, nil),
			"snapshot_type":               utils.PathSearch("snapshot_type", s, nil),
			"progress":                    utils.PathSearch("progress", s, nil),
			"encrypt_algorithm":           utils.PathSearch("encrypt_algorithm", s, nil),
			"snapshot_chains":             flattenEvsv5SnapshotChains(utils.PathSearch("snapshot_chains", s, []interface{}{}).([]interface{})),
			"snapshot_group_id":           utils.PathSearch("snapshot_group_id", s, nil),
		})
	}
	return rst
}

func flattenEvsv5SnapshotChains(chains []interface{}) []map[string]interface{} {
	if len(chains) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(chains))
	for _, c := range chains {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", c, nil),
			"availability_zone": utils.PathSearch("availability_zone", c, nil),
			"snapshot_count":    utils.PathSearch("snapshot_count", c, nil),
			"capacity":          utils.PathSearch("capacity", c, nil),
			"volume_id":         utils.PathSearch("volume_id", c, nil),
			"category":          utils.PathSearch("category", c, nil),
			"created_at":        utils.PathSearch("created_at", c, nil),
			"updated_at":        utils.PathSearch("updated_at", c, nil),
		})
	}
	return rst
}
