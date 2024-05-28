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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes
func DataSourceAppcodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppcodesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the application belongs.`,
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the application to be queried.`,
			},
			"appcodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All APPCODEs of the specified application.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the APPCODE.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The APPCODE value (content).`,
						},
						"application_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the application.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the APPCODE, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

func queryAppcodes(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes"
		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
		offset     = 0
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{app_id}", appId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s?limit=100&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving APPCODEs under specified application (%s): %s", appId, err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		appcodes := utils.PathSearch("app_codes", respBody, make([]interface{}, 0)).([]interface{})
		if len(appcodes) < 1 {
			break
		}
		result = append(result, appcodes...)
		offset += len(appcodes)
	}
	return result, nil
}

func dataSourceAppcodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}
	signatures, err := queryAppcodes(client, d)
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
		d.Set("appcodes", flattenAppcodes(signatures)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAppcodes(appcodes []interface{}) []interface{} {
	if len(appcodes) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(appcodes))
	for _, appcode := range appcodes {
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("id", appcode, nil),
			"value":          utils.PathSearch("app_code", appcode, nil),
			"application_id": utils.PathSearch("app_id", appcode, nil),
			"created_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("create_time",
				appcode, "").(string))/1000, false),
		})
	}
	return result
}
