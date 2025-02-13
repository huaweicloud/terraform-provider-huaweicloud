package fgs

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

// @API FunctionGraph GET /v2/{project_id}/fgs/applications
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
			"application_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The application ID used to query specified application.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The application name used to query specified application.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status of the application to be queried.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the application to be queried.`,
			},
			"applications": {
				Type:        schema.TypeList,
				Elem:        applicationSchema(),
				Computed:    true,
				Description: `All applications that match the filter parameters.`,
			},
		},
	}
}

func applicationSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of application.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of application.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of application.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of application.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the application, in RFC3339 format.`,
			},
		},
	}
	return &sc
}

func getApplications(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/fgs/applications"
		marker  float64
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithMarker := fmt.Sprintf("%s?marker=%v", listPath, marker)
		requestResp, err := client.Request("GET", listPathWithMarker, &listOpt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		applications := utils.PathSearch("applications", respBody, make([]interface{}, 0)).([]interface{})
		if len(applications) < 1 {
			break
		}
		result = append(result, applications...)
		// The attribute page_info.next_marker is offline.
		// In this API, marker has the same meaning as offset.
		nextMarker := utils.PathSearch("next_marker", respBody, float64(0)).(float64)
		if nextMarker == marker || nextMarker == 0 {
			// Make sure the next marker value is correct, not the previous marker or zero (in the last page).
			break
		}
		marker = nextMarker
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
			"id":     utils.PathSearch("id", application, nil),
			"name":   utils.PathSearch("name", application, nil),
			"status": utils.PathSearch("status", application, nil),
			"updated_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("last_modified_time",
				application, float64(0)).(float64))/1000, false),
			"description": utils.PathSearch("description", application, nil),
		})
	}
	return result
}

func filterApplications(d *schema.ResourceData, functions []interface{}) []interface{} {
	result := functions

	if appId, ok := d.GetOk("application_id"); ok && len(result) > 0 {
		result = utils.PathSearch(fmt.Sprintf("[?id=='%v']", appId), result, make([]interface{}, 0)).([]interface{})
	}

	if name, ok := d.GetOk("name"); ok && len(result) > 0 {
		result = utils.PathSearch(fmt.Sprintf("[?name=='%v']", name), result, make([]interface{}, 0)).([]interface{})
	}

	if status, ok := d.GetOk("status"); ok && len(result) > 0 {
		result = utils.PathSearch(fmt.Sprintf("[?status=='%v']", status), result, make([]interface{}, 0)).([]interface{})
	}

	if desc, ok := d.GetOk("description"); ok && len(result) > 0 {
		result = utils.PathSearch(fmt.Sprintf("[?description=='%v']", desc), result, make([]interface{}, 0)).([]interface{})
	}

	return result
}

func dataSourceApplicationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	applications, err := getApplications(client)
	if err != nil {
		return diag.Errorf("error querying the applications: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("applications", flattenApplications(filterApplications(d, applications))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
