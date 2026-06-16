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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes
func DataSourceAppcodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppcodesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the APPCODEs are located.`,
			},

			// Required parameters.
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

			// Attributes.
			"appcodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        appcodeSchema(),
				Description: `The list of the APPCODEs of the specified application.`,
			},
		},
	}
}

func appcodeSchema() *schema.Resource {
	return &schema.Resource{
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
	}
}

func listAppcodes(client *golangsdk.ServiceClient, instanceId, appId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{app_id}", appId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

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

		appcodes := utils.PathSearch("app_codes", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, appcodes...)
		if len(appcodes) < limit {
			break
		}
		offset += len(appcodes)
	}

	return result, nil
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

func dataSourceAppcodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	resp, err := listAppcodes(client, instanceId, appId)
	if err != nil {
		return diag.Errorf("error querying APPCODEs: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("appcodes", flattenAppcodes(resp)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
