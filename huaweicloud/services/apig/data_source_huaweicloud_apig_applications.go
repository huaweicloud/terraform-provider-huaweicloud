package apig

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
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
				Description: `The region where the applications are located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the applications belong.`,
			},

			// Optional parameters.
			"application_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the application to be queried.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the application to be queried.`,
			},
			"app_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The key of the application to be queried.`,
			},

			// Attributes.
			"applications": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the applications that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the application.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the application.`,
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The status of the application.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the application.`,
						},
						"app_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: `The key of the application.`,
						},
						"app_secret": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: `The secret of the application.`,
						},
						"app_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the application.`,
						},
						"bind_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of bound APIs.`,
						},
						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of the application.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the application, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the application, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

func buildApplicationsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("application_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("app_key"); ok {
		res = fmt.Sprintf("%s&app_key=%v", res, v)
	}

	return res
}

func listApplications(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/apps?limit={limit}"
		instanceId = d.Get("instance_id").(string)
		limit      = 500
		offset     = 0
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildApplicationsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		applications := utils.PathSearch("apps", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, applications...)
		if len(applications) < limit {
			break
		}
		offset += len(applications)
	}
	return result, nil
}

func flattenApplications(applications []interface{}) []interface{} {
	if len(applications) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(applications))
	for _, application := range applications {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", application, nil),
			"name":        utils.PathSearch("name", application, nil),
			"status":      utils.PathSearch("status", application, nil),
			"description": utils.PathSearch("remark", application, nil),
			"app_key":     utils.PathSearch("app_key", application, nil),
			"app_secret":  utils.PathSearch("app_secret", application, nil),
			"app_type":    utils.PathSearch("app_type", application, nil),
			"bind_num":    utils.PathSearch("bind_num", application, nil),
			"created_by":  utils.PathSearch("creator", application, nil),
			"created_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("register_time",
				application, "").(string))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("update_time",
				application, "").(string))/1000, false),
		})
	}
	return result
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

	applications, err := listApplications(client, d)
	if err != nil {
		return diag.Errorf("error querying applications: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("applications", flattenApplications(applications)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
