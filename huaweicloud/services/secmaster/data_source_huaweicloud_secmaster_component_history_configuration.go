package secmaster

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/components/{component_id}/configurations/versions
func DataSourceComponentHistoryConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComponentHistoryConfigurationRead,

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
			"component_id": {
				Type:     schema.TypeString,
				Required: true,
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
				Elem:     componentHistoryConfigurationRecordsSchema(),
			},
		},
	}
}

func componentHistoryConfigurationRecordsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"configuration_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     componentHistoryConfigurationListSchema(),
			},
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"specification": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func componentHistoryConfigurationListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"param": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildComponentHistoryConfigurationQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("sort_key"); ok {
		queryParams = fmt.Sprintf("%s&sort_key=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		queryParams = fmt.Sprintf("%s&sort_dir=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceComponentHistoryConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/components/{component_id}/configurations/versions"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{component_id}", d.Get("component_id").(string))
	requestPath += buildComponentHistoryConfigurationQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster component history configuration: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		recordsResp := utils.PathSearch("records", respBody, make([]interface{}, 0)).([]interface{})
		if len(recordsResp) == 0 {
			break
		}

		result = append(result, recordsResp...)
		offset += len(recordsResp)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenComponentHistoryConfigurationRecords(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenComponentHistoryConfigurationRecords(recordsResp []interface{}) []interface{} {
	if len(recordsResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(recordsResp))
	for _, v := range recordsResp {
		rst = append(rst, map[string]interface{}{
			"configuration_status": utils.PathSearch("configuration_status", v, nil),
			"list": flattenComponentHistoryConfigurationList(
				utils.PathSearch("list", v, make([]interface{}, 0)).([]interface{})),
			"node_id":       utils.PathSearch("node_id", v, nil),
			"node_name":     utils.PathSearch("node_name", v, nil),
			"node_status":   utils.PathSearch("node_status", v, nil),
			"specification": utils.PathSearch("specification", v, nil),
		})
	}

	return rst
}

func flattenComponentHistoryConfigurationList(listResp []interface{}) []interface{} {
	if len(listResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(listResp))
	for _, v := range listResp {
		rst = append(rst, map[string]interface{}{
			"configuration_id": utils.PathSearch("configuration_id", v, nil),
			"file_name":        utils.PathSearch("file_name", v, nil),
			"file_type":        utils.PathSearch("file_type", v, nil),
			"node_id":          utils.PathSearch("node_id", v, nil),
			"param":            utils.PathSearch("param", v, nil),
			"type":             utils.PathSearch("type", v, nil),
			"version":          utils.PathSearch("version", v, nil),
		})
	}

	return rst
}
