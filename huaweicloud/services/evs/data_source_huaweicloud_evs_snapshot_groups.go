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

// @API EVS GET /v5/{project_id}/snapshot-groups/detail
func DataSourceEvsSnapshotGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvsSnapshotGroupsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
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
			"tag_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
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
			"server_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshot_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     snapshotGroupSchema(),
			},
		},
	}
}

func snapshotGroupSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
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
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceEvsSnapshotGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v5/{project_id}/snapshot-groups/detail"
		product   = "evs"
		allGroups = make([]interface{}, 0)
		offset    = 0
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	for {
		pagedPath := fmt.Sprintf("%s?limit=1000&offset=%d%s", requestPath, offset, buildEvsSnapshotGroupsQueryParams(d, cfg))
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		resp, err := client.Request("GET", pagedPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error querying EVS snapshot groups: %s", err)
		}
		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}
		groups := utils.PathSearch("snapshot_groups", respBody, make([]interface{}, 0)).([]interface{})
		if len(groups) == 0 {
			break
		}
		allGroups = append(allGroups, groups...)
		offset += len(groups)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("snapshot_groups", flattenEvsSnapshotGroups(allGroups)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildEvsSnapshotGroupsQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	params := ""
	epsId := cfg.GetEnterpriseProjectID(d)
	if epsId != "" {
		params += fmt.Sprintf("&enterprise_project_id=%v", epsId)
	}
	if v, ok := d.GetOk("id"); ok {
		params += fmt.Sprintf("&id=%v", v)
	}
	if v, ok := d.GetOk("name"); ok {
		params += fmt.Sprintf("&name=%v", v)
	}
	if v, ok := d.GetOk("status"); ok {
		params += fmt.Sprintf("&status=%v", v)
	}
	if v, ok := d.GetOk("tag_key"); ok {
		params += fmt.Sprintf("&tag_key=%v", v)
	}
	if v, ok := d.GetOk("tags"); ok {
		params += fmt.Sprintf("&tags=%v", v)
	}
	if v, ok := d.GetOk("sort_key"); ok {
		params += fmt.Sprintf("&sort_key=%v", v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		params += fmt.Sprintf("&sort_dir=%v", v)
	}
	if v, ok := d.GetOk("server_id"); ok {
		params += fmt.Sprintf("&server_id=%v", v)
	}
	return params
}

func flattenEvsSnapshotGroups(groups []interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0, len(groups))
	for _, g := range groups {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", g, nil),
			"created_at":            utils.PathSearch("created_at", g, nil),
			"status":                utils.PathSearch("status", g, nil),
			"updated_at":            utils.PathSearch("updated_at", g, nil),
			"name":                  utils.PathSearch("name", g, nil),
			"description":           utils.PathSearch("description", g, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", g, nil),
			"tags":                  utils.PathSearch("tags", g, map[string]interface{}{}).(map[string]interface{}),
			"server_id":             utils.PathSearch("server_id", g, nil),
		})
	}
	return rst
}
