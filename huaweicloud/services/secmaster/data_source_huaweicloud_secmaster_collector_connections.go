package secmaster

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

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/collector/connections
func DataSourceCollectorConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCollectorConnectionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"connection_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"title": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
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
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel_refer_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"info": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"module_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildConnectionQueryParams(d *schema.ResourceData, offset int) string {
	queryParams := ""

	if v, ok := d.GetOk("connection_type"); ok {
		queryParams = fmt.Sprintf("%s&connection_type=%v", queryParams, v)
	}

	if v, ok := d.GetOk("title"); ok {
		queryParams = fmt.Sprintf("%s&title=%v", queryParams, v)
	}

	if v, ok := d.GetOk("description"); ok {
		queryParams = fmt.Sprintf("%s&description=%v", queryParams, v)
	}

	if v, ok := d.GetOk("sort_key"); ok {
		queryParams = fmt.Sprintf("%s&sort_key=%v", queryParams, v)
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		queryParams = fmt.Sprintf("%s&sort_dir=%v", queryParams, v)
	}

	if offset > 0 {
		queryParams = fmt.Sprintf("%s&offset=%v", queryParams, offset)
	}

	if len(queryParams) > 0 {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceCollectorConnectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/collector/connections"
		offset      = 0
		allResult   = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := getPath + buildConnectionQueryParams(d, offset)
		getResp, err := client.Request("GET", requestPathWithOffset, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving collector connections: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		records := utils.PathSearch("records", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(records) == 0 {
			break
		}

		allResult = append(allResult, records...)
		offset += len(records)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenCollectorRecords(allResult)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCollectorRecords(allResult []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(allResult))
	for _, v := range allResult {
		rst = append(rst, map[string]interface{}{
			"channel_refer_count": utils.PathSearch("channel_refer_count", v, nil),
			"connection_id":       utils.PathSearch("connection_id", v, nil),
			"connection_type":     utils.PathSearch("connection_type", v, nil),
			"description":         utils.PathSearch("description", v, nil),
			"info":                utils.PathSearch("info", v, nil),
			"module_id":           utils.PathSearch("module_id", v, nil),
			"template_title":      utils.PathSearch("template_title", v, nil),
			"title":               utils.PathSearch("title", v, nil),
		})
	}
	return rst
}
