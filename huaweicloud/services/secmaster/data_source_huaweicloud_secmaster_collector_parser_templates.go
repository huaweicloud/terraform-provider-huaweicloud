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

// @API Secmaster GET /v1/{project_id}/workspaces/{workspace_id}/collector/logstash/parsers/templates
func DataSourceSecmasterCollectorParserTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecmasterCollectorParserTemplatesRead,

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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parser_id": {
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

func buildCollectorParserTemplatesQueryParams(d *schema.ResourceData, offset int) string {
	rst := ""

	if v, ok := d.GetOk("title"); ok {
		rst += fmt.Sprintf("&title=%v", v)
	}

	if v, ok := d.GetOk("description"); ok {
		rst += fmt.Sprintf("&description=%v", v)
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		rst += fmt.Sprintf("&sort_dir=%v", v)
	}

	if v, ok := d.GetOk("sort_key"); ok {
		rst += fmt.Sprintf("&sort_key=%v", v)
	}

	if offset > 0 {
		rst += fmt.Sprintf("&offset=%d", offset)
	}

	if rst != "" {
		rst = "?" + rst[1:]
	}

	return rst
}

func dataSourceSecmasterCollectorParserTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/collector/logstash/parsers/templates"
		offset  = 0
		allData = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestCurrentPath := requestPath + buildCollectorParserTemplatesQueryParams(d, offset)
		resp, err := client.Request("GET", requestCurrentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster collector parser templates: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		records := utils.PathSearch("records", respBody, make([]interface{}, 0)).([]interface{})
		if len(records) == 0 {
			break
		}

		allData = append(allData, records...)
		offset += len(records)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenCollectorParserTemplates(allData)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCollectorParserTemplates(allData []interface{}) []interface{} {
	if len(allData) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(allData))
	for _, v := range allData {
		rst = append(rst, map[string]interface{}{
			"description": utils.PathSearch("description", v, nil),
			"parser_id":   utils.PathSearch("parser_id", v, nil),
			"title":       utils.PathSearch("title", v, nil),
		})
	}

	return rst
}
