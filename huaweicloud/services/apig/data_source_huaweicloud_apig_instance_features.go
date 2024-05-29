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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/features
func DataSourceInstanceFeatures() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceFeaturesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specified the ID of the dedicated instance to which the features belong.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specified the name of the feature.`,
			},
			"features": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All instance features that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the feature.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the feature.`,
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the feature is enabled.`,
						},
						"config": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The detailed configuration of the instance feature.`,
						},
						// The format is "yyyy-MM-ddTHH:mm:ss+08:00".
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the feature, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

func queryInstanceFeature(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/features?limit=500"
		instanceId = d.Get("instance_id").(string)
		offset     = 0
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving features under specified dedicated instance (%s): %s", instanceId, err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		features := utils.PathSearch("features", respBody, make([]interface{}, 0)).([]interface{})
		if len(features) < 1 {
			break
		}
		result = append(result, features...)
		offset += len(features)
	}
	return result, nil
}

func dataSourceInstanceFeaturesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}
	features, err := queryInstanceFeature(client, d)
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
		d.Set("features", filterInstanceFeatures(flattenInstanceFeatures(features), d)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func filterInstanceFeatures(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("name"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("name", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func flattenInstanceFeatures(features []interface{}) []interface{} {
	if len(features) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(features))
	for _, feature := range features {
		updateTime := utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("update_time", feature, "2006-01-02T15:04:05Z").(string))
		result = append(result, map[string]interface{}{
			"id":      utils.PathSearch("id", feature, nil),
			"name":    utils.PathSearch("name", feature, nil),
			"enabled": utils.PathSearch("enable", feature, nil),
			"config":  utils.PathSearch("config", feature, nil),
			// If this feature has not been configured, the time format is "0001-01-01T00:00:00Z",
			// the corresponding timestamp is a negative, and this format is uniformly processed as an empty string.
			"updated_at": utils.FormatTimeStampRFC3339(updateTime/1000, false),
		})
	}
	return result
}
