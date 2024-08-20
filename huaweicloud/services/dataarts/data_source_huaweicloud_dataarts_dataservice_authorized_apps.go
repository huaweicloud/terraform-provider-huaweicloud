package dataarts

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

// @API DataArtsStudio GET /v1/{project_id}/service/authorize/apis/{api_id}
func DataSourceDataServiceAuthorizedApps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataServiceAuthorizedAppsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the API is located.`,
			},

			// Parameters in request header
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the API belongs.`,
			},
			"dlm_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of DLM engine.`,
			},

			// Argument
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the API used to authorize the APPs.`,
			},

			// Attributes
			"apps": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataServiceAuthorizedAppElemSchema(),
				Description: `All APPs authorized by API.`,
			},
		},
	}
}

func dataServiceAuthorizedAppElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the application that has authorization.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the application that has authorization.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance ID to which the authorized API belongs.`,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance name to which the authorized API belongs.`,
			},
			"expired_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The expiration time, in RFC3339 format.`,
			},
			"approved_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The approve time, in RFC3339 format.`,
			},
			"relationship_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The relationship between the authorized API and the authorized APP list.`,
			},
			"static_params": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataServiceAuthorizedAppStaticParamElemSchema(),
				Description: `The configuration of the static parameters.`,
			},
		},
	}
	return &sc
}

func dataServiceAuthorizedAppStaticParamElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the static parameter.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the static parameter.`,
			},
		},
	}
	return &sc
}

func queryDataServiceAuthorizedApps(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/service/authorize/apis/{api_id}?limit=100"
		offset  = 0
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{api_id}", d.Get("api_id").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    d.Get("workspace_id").(string),
			"dlm-type":     d.Get("dlm_type").(string),
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		apps := utils.PathSearch("records", respBody, make([]interface{}, 0)).([]interface{})
		if len(apps) < 1 {
			break
		}
		result = append(result, apps...)
		offset += len(apps)
	}

	return result, nil
}

func flattenDataServiceAuthorizedApps(apps []interface{}) []interface{} {
	result := make([]interface{}, 0, len(apps))

	for _, app := range apps {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("app_id", app, nil),
			"name":              utils.PathSearch("app_name", app, nil),
			"instance_id":       utils.PathSearch("instance_id", app, nil),
			"instance_name":     utils.PathSearch("instance_name", app, nil),
			"expired_at":        utils.FormatTimeStampRFC3339(int64(utils.PathSearch("api_using_time", app, float64(0)).(float64))/1000, false),
			"approved_at":       utils.FormatTimeStampRFC3339(int64(utils.PathSearch("approve_time", app, float64(0)).(float64))/1000, false),
			"relationship_type": utils.PathSearch("relationship_type", app, nil),
			"static_params":     flattenDataServiceAppStaticParams(utils.PathSearch("static_params", app, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenDataServiceAppStaticParams(staticParams []interface{}) []interface{} {
	result := make([]interface{}, 0, len(staticParams))

	for _, staticParam := range staticParams {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("param_name", staticParam, nil),
			"value": utils.PathSearch("param_value", staticParam, nil),
		})
	}

	return result
}

func dataSourceDataServiceAuthorizedAppsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	authorizedApps, err := queryDataServiceAuthorizedApps(client, d)
	if err != nil {
		return diag.Errorf("error getting Data Service API (%s) for DataArts Studio", d.Id())
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("apps", flattenDataServiceAuthorizedApps(authorizedApps)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
