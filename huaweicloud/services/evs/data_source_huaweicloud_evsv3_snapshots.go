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

// @API EVS GET /v3/{project_id}/snapshots/detail
func DataSourceV3Snapshots() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceV3SnapshotsRead,

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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshots": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
						"metadata": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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
					},
				},
			},
		},
	}
}

func buildV3SnapshotsGetPathWithOffset(queryParams string, offset int) string {
	if offset == 0 {
		// Ignore the offset of the first page.
		return queryParams
	}

	if queryParams == "" {
		// No query conditions were added.
		return fmt.Sprintf("?offset=%d", offset)
	}

	return fmt.Sprintf("%s&offset=%d", queryParams, offset)
}

func buildV3SnapshotsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("volume_id"); ok {
		res = fmt.Sprintf("%s&volume_id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}

func datasourceV3SnapshotsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		mErr         *multierror.Error
		product      = "evs"
		getHttpUrl   = "v3/{project_id}/snapshots/detail"
		queryParams  = buildV3SnapshotsQueryParams(d)
		offset       = 0
		allSnapshots []interface{}
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		getPathWithOffset := getPath + buildV3SnapshotsGetPathWithOffset(queryParams, offset)
		resp, err := client.Request("GET", getPathWithOffset, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving EVS v3 snapshots: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		snapshotsResp := utils.PathSearch("snapshots", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(snapshotsResp) == 0 {
			break
		}

		allSnapshots = append(allSnapshots, snapshotsResp...)
		offset += len(snapshotsResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("snapshots", flattenV3Snapshots(allSnapshots)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenV3Snapshots(allSnapshots []interface{}) []interface{} {
	if len(allSnapshots) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(allSnapshots))
	for _, v := range allSnapshots {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"created_at":  utils.PathSearch("created_at", v, nil),
			"updated_at":  utils.PathSearch("updated_at", v, nil),
			"metadata":    utils.PathSearch("metadata", v, nil),
			"volume_id":   utils.PathSearch("volume_id", v, nil),
			"size":        utils.PathSearch("size", v, nil),
			"status":      utils.PathSearch("status", v, nil),
		})
	}

	return rst
}
