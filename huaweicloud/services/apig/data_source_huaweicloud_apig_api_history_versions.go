package apig

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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apis/publish/{api_id}
func DataSourceApiHistoryVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApiHistoryVersionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the API history versions are located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the API belongs.`,
			},
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the API.`,
			},

			// Optional parameters.
			"env_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the environment.`,
			},
			"env_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the environment.`,
			},

			// Attributes.
			"api_versions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the API history versions that matched filter parameters.`,
				Elem:        apiHistoryVersionSchema(),
			},
		},
	}
}

func apiHistoryVersionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the API history version.`,
			},
			"number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version number of the API.`,
			},
			"api_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the API.`,
			},
			"env_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the published environment.`,
			},
			"env_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the published environment.`,
			},
			"remark": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The publish description.`,
			},
			"publish_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The publish time of the version, in RFC3339 format.`,
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The status of the API version.`,
			},
		},
	}
}

func buildApiHistoryVersionsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("env_id"); ok {
		res = fmt.Sprintf("%s&env_id=%v", res, v)
	}
	if v, ok := d.GetOk("env_name"); ok {
		res = fmt.Sprintf("%s&env_name=%v", res, v)
	}
	return res
}

func listApiHistoryVersions(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		result     = make([]interface{}, 0)
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/apis/publish/{api_id}?limit={limit}"
		instanceId = d.Get("instance_id").(string)
		apiId      = d.Get("api_id").(string)
		limit      = 100
		offset     = 0
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{api_id}", apiId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildApiHistoryVersionsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
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

		apiVersions := utils.PathSearch("api_versions", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, apiVersions...)
		if len(apiVersions) < limit {
			break
		}
		offset += len(apiVersions)
	}

	return result, nil
}

func flattenApiHistoryVersions(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("version_id", item, nil),
			"number":       utils.PathSearch("version_no", item, nil),
			"api_id":       utils.PathSearch("api_id", item, nil),
			"env_id":       utils.PathSearch("env_id", item, nil),
			"env_name":     utils.PathSearch("env_name", item, nil),
			"remark":       utils.PathSearch("remark", item, nil),
			"publish_time": utils.PathSearch("publish_time", item, nil),
			"status":       utils.PathSearch("status", item, nil),
		})
	}

	return result
}

func dataSourceApiHistoryVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		apiId  = d.Get("api_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	resp, err := listApiHistoryVersions(client, d)
	if err != nil {
		return diag.Errorf("error querying history versions for specified API (%s): %s", apiId, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("api_versions", flattenApiHistoryVersions(resp)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
