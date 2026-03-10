package css

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

// @API CSS GET /v1.0/{project_id}/resize-flavors
func DataSourceResizeFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResizeFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datastore_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dbname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"versions": {
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
						"flavors": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"str_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cpu": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ram": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"diskrange": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"typename": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cond_operation_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"localdisk": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"edge": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildResizeFlavorsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceResizeFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v1.0/{project_id}/resize-flavors"
		clusterId = d.Get("cluster_id").(string)
	)

	client, err := cfg.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = fmt.Sprintf("%s?clusterId=%s", getPath, clusterId)
	getPath += buildResizeFlavorsQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving the flavors: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("datastore_id", utils.PathSearch("id", getRespBody, nil)),
		d.Set("dbname", utils.PathSearch("dbname", getRespBody, nil)),
		d.Set("versions", flattenResizeFlavorsVersions(
			utils.PathSearch("versions", getRespBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenResizeFlavorsVersions(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"id":   utils.PathSearch("id", v, nil),
			"name": utils.PathSearch("name", v, nil),
			"flavors": flattenResizeFlavors(
				utils.PathSearch("flavors", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenResizeFlavors(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"str_id":                utils.PathSearch("str_id", v, nil),
			"cpu":                   utils.PathSearch("cpu", v, nil),
			"ram":                   utils.PathSearch("ram", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"region":                utils.PathSearch("region", v, nil),
			"diskrange":             utils.PathSearch("diskrange", v, nil),
			"typename":              utils.PathSearch("typename", v, nil),
			"cond_operation_status": utils.PathSearch("condOperationStatus", v, nil),
			"localdisk":             utils.PathSearch("localdisk", v, nil),
			"edge":                  utils.PathSearch("edge", v, nil),
		})
	}

	return result
}
