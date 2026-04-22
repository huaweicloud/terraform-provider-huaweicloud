package dataarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v1/{project_id}/service/authorize/apps/{app_id}
func DataSourceDataServiceAppAuthorizedApis() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataServiceAppAuthorizedApisRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the authorized APIs are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the APP belongs.`,
			},
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the APP used to query authorized APIs.`,
			},

			// Attributes.
			"apis": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataServiceAppAuthorizedApiElemSchema(),
				Description: `The list of APIs authorized to the APP.`,
			},
		},
	}
}

func dataServiceAppAuthorizedApiElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the API.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the API.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the API.`,
			},
			"approval_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The approval time, in RFC3339 format.`,
			},
			"manager": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the API reviewer.`,
			},
			"deadline": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The deadline for using the API, in RFC3339 format.`,
			},
			"relationship_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The relationship between the authorized API and the APP.`,
			},
			"static_params": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataServiceAppAuthorizedApiStaticParamElemSchema(),
				Description: `The configuration of the static parameters.`,
			},
		},
	}
}

func dataServiceAppAuthorizedApiStaticParamElemSchema() *schema.Resource {
	return &schema.Resource{
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
}

func queryDataServiceAppAuthorizedApis(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/service/authorize/apps/{app_id}?limit={limit}"
		offset  = 0
		limit   = 100
		result  = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{app_id}", d.Get("app_id").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"dlm-type":     "EXCLUSIVE",
			"workspace":    d.Get("workspace_id").(string),
		},
	}

	for {
		listPathWithOffset := listPathWithLimit + fmt.Sprintf("&offset=%s", strconv.Itoa(offset))
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		apis := utils.PathSearch("records", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, apis...)
		if len(apis) < limit {
			break
		}
		offset += len(apis)
	}

	return result, nil
}

func flattenDataServiceAppAuthorizedApiStaticParams(staticParams []interface{}) []interface{} {
	result := make([]interface{}, 0, len(staticParams))

	for _, staticParam := range staticParams {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("para_name", staticParam, nil),
			"value": utils.PathSearch("para_value", staticParam, nil),
		})
	}

	return result
}

func flattenDataServiceAppAuthorizedApis(apis []interface{}) []interface{} {
	result := make([]interface{}, 0, len(apis))

	for _, api := range apis {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("id", api, nil),
			"name":              utils.PathSearch("name", api, nil),
			"description":       utils.PathSearch("description", api, nil),
			"approval_time":     utils.FormatTimeStampRFC3339(int64(utils.PathSearch("approval_time", api, float64(0)).(float64))/1000, false),
			"manager":           utils.PathSearch("manager", api, nil),
			"deadline":          utils.FormatTimeStampRFC3339(int64(utils.PathSearch("deadline", api, float64(0)).(float64))/1000, false),
			"relationship_type": utils.PathSearch("relationship_type", api, nil),
			"static_params": flattenDataServiceAppAuthorizedApiStaticParams(
				utils.PathSearch("static_params", api, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func dataSourceDataServiceAppAuthorizedApisRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	authorizedApis, err := queryDataServiceAppAuthorizedApis(client, d)
	if err != nil {
		return diag.Errorf("error getting authorized APIs for APP (%s): %s", d.Get("app_id").(string), err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("apis", flattenDataServiceAppAuthorizedApis(authorizedApis)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
