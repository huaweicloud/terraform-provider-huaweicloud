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

// @API Secmaster GET /v1/{project_id}/workspaces/{workspace_id}/siem/pipes
func DataSourceSecmasterPipes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecmasterPipesRead,

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
			"dataspace_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pipe_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pipe_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// There's a problem with the API; passing this field `pipe_type` will cause an error.
			"pipe_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"dataspace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dataspace_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pipe_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pipe_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pipe_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"shards": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"storage_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildPipesQueryParams(d *schema.ResourceData, offset int) string {
	rst := ""

	if v, ok := d.GetOk("dataspace_id"); ok {
		rst += fmt.Sprintf("&dataspace_id=%v", v)
	}

	if v, ok := d.GetOk("pipe_id"); ok {
		rst += fmt.Sprintf("&pipe_id=%v", v)
	}

	if v, ok := d.GetOk("pipe_name"); ok {
		rst += fmt.Sprintf("&pipe_name=%v", v)
	}

	if v, ok := d.GetOk("pipe_type"); ok {
		rst += fmt.Sprintf("&pipe_type=%v", v)
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

func dataSourceSecmasterPipesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/siem/pipes"
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
		requestCurrentPath := requestPath + buildPipesQueryParams(d, offset)
		resp, err := client.Request("GET", requestCurrentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster pipes: %s", err)
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
		d.Set("records", flattenPipesAttribute(allData)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPipesAttribute(allData []interface{}) []interface{} {
	if len(allData) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(allData))
	for _, v := range allData {
		rst = append(rst, map[string]interface{}{
			"create_by":      utils.PathSearch("create_by", v, nil),
			"create_time":    utils.PathSearch("create_time", v, nil),
			"dataspace_id":   utils.PathSearch("dataspace_id", v, nil),
			"dataspace_name": utils.PathSearch("dataspace_name", v, nil),
			"description":    utils.PathSearch("description", v, nil),
			"domain_id":      utils.PathSearch("domain_id", v, nil),
			"pipe_id":        utils.PathSearch("pipe_id", v, nil),
			"pipe_name":      utils.PathSearch("pipe_name", v, nil),
			"pipe_type":      utils.PathSearch("pipe_type", v, nil),
			"project_id":     utils.PathSearch("project_id", v, nil),
			"shards":         utils.PathSearch("shards", v, nil),
			"storage_period": utils.PathSearch("storage_period", v, nil),
			"update_by":      utils.PathSearch("update_by", v, nil),
			"update_time":    utils.PathSearch("update_time", v, nil),
		})
	}

	return rst
}
