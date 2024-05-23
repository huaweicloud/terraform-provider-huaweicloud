package apig

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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apps
func DataSourceApplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the data source.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of the dedicated instance to which the applications belong.",
			},
			"application_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the application to be queried.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the application to be queried.",
			},
			"app_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the key of the application to be queried.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the creator of the application to be queried.",
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the application.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the application.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The status of the application.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the application.",
						},
						"app_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The key of the application.",
						},
						"app_secret": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The secret of the application.",
						},
						"app_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the application.",
						},
						"bind_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of bound APIs.",
						},
						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creator of the application.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the application.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The latest update time of the application.",
						},
					},
				},
			},
		},
	}
}

func buildListApplicationsParams(d *schema.ResourceData) string {
	res := ""
	if applicationId, ok := d.GetOk("application_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, applicationId)
	}
	if name, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, name)
	}
	if appKey, ok := d.GetOk("app_key"); ok {
		res = fmt.Sprintf("%s&app_key=%v", res, appKey)
	}
	return res
}

func queryApplications(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/apps?limit=500"
		instanceId = d.Get("instance_id").(string)
		offset     = 0
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)

	queryParams := buildListApplicationsParams(d)
	listPath += queryParams

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving applications under specified "+
				"dedicated instance (%s): %s", instanceId, err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		applications := utils.PathSearch("apps", respBody, make([]interface{}, 0)).([]interface{})
		if len(applications) < 1 {
			break
		}
		result = append(result, applications...)
		offset += len(applications)
	}
	return result, nil
}

func dataSourceApplicationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}
	applications, err := queryApplications(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("applications", filterApplications(flattenApplications(applications), d)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func filterApplications(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))

	for _, v := range all {
		if param, ok := d.GetOk("created_by"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("created_by", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func flattenApplications(applications []interface{}) []interface{} {
	if len(applications) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(applications))
	for _, policy := range applications {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", policy, nil),
			"name":        utils.PathSearch("name", policy, nil),
			"status":      utils.PathSearch("status", policy, nil),
			"description": utils.PathSearch("remark", policy, nil),
			"app_key":     utils.PathSearch("app_key", policy, nil),
			"app_secret":  utils.PathSearch("app_secret", policy, nil),
			"app_type":    utils.PathSearch("app_type", policy, nil),
			"bind_num":    utils.PathSearch("bind_num", policy, nil),
			"created_by":  utils.PathSearch("creator", policy, nil),
			"created_at":  utils.PathSearch("register_time", policy, nil),
			"updated_at":  utils.PathSearch("update_time", policy, nil),
		})
	}
	return result
}
