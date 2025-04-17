package evs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EVS GET /v2/{project_id}/cloudsnapshots/detail
func DataSourceEvsSnapshots() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceSnapshotsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dedicated_storage_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dedicated_storage_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshots": {
				Type:     schema.TypeList,
				Elem:     snapshotSchema(),
				Computed: true,
			},
		},
	}
}

func snapshotSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
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
			"service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dedicated_storage_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dedicated_storage_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"progress": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func datasourceSnapshotsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getEVSSnapshots: Query EVS snapshots
	var (
		getEVSSnapshotsHttpUrl = "v2/{project_id}/cloudsnapshots/detail"
		getEVSSnapshotsProduct = "evs"
	)
	getEVSSnapshotsClient, err := cfg.NewServiceClient(getEVSSnapshotsProduct, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	getEVSSnapshotsPath := getEVSSnapshotsClient.Endpoint + getEVSSnapshotsHttpUrl
	getEVSSnapshotsPath = strings.ReplaceAll(getEVSSnapshotsPath, "{project_id}",
		getEVSSnapshotsClient.ProjectID)
	getEVSSnapshotsPath += buildEVSSnapshotsQueryParams(d, cfg)

	getEVSSnapshotsResp, err := pagination.ListAllItems(
		getEVSSnapshotsClient,
		"offset",
		getEVSSnapshotsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving snapshots, %s", err)
	}

	listEVSSnapshotsRespJson, err := json.Marshal(getEVSSnapshotsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listEVSSnapshotsRespBody interface{}
	err = json.Unmarshal(listEVSSnapshotsRespJson, &listEVSSnapshotsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}
	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("snapshots", flattenListSnapshotsBody(listEVSSnapshotsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildEVSSnapshotsQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := ""
	epsId := cfg.GetEnterpriseProjectID(d, "all_granted_eps")
	if epsId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsId)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("volume_id"); ok {
		res = fmt.Sprintf("%s&volume_id=%v", res, v)
	}
	if v, ok := d.GetOk("availability_zone"); ok {
		res = fmt.Sprintf("%s&availability_zone=%v", res, v)
	}
	if v, ok := d.GetOk("snapshot_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("dedicated_storage_id"); ok {
		res = fmt.Sprintf("%s&dedicated_storage_id=%v", res, v)
	}
	if v, ok := d.GetOk("dedicated_storage_name"); ok {
		res = fmt.Sprintf("%s&dedicated_storage_name=%v", res, v)
	}
	if v, ok := d.GetOk("service_type"); ok {
		res = fmt.Sprintf("%s&service_type=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenListSnapshotsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("snapshots", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		progress := utils.PathSearch("\"os-extended-snapshot-attributes:progress\"", v, nil)

		rst = append(rst, map[string]interface{}{
			"id":                     utils.PathSearch("id", v, nil),
			"name":                   utils.PathSearch("name", v, nil),
			"status":                 utils.PathSearch("status", v, nil),
			"description":            utils.PathSearch("description", v, nil),
			"size":                   utils.PathSearch("size", v, nil),
			"created_at":             utils.PathSearch("created_at", v, nil),
			"updated_at":             utils.PathSearch("updated_at", v, nil),
			"volume_id":              utils.PathSearch("volume_id", v, nil),
			"service_type":           utils.PathSearch("service_type", v, nil),
			"metadata":               utils.PathSearch("metadata", v, nil),
			"dedicated_storage_id":   utils.PathSearch("dedicated_storage_id", v, nil),
			"dedicated_storage_name": utils.PathSearch("dedicated_storage_name", v, nil),
			"progress":               progress,
		})
	}
	return rst
}
