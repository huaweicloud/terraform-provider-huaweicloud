package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v1/{project_id}/features
func DataSourceDscFeatures() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscFeaturesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     featureVoSchema(),
			},
		},
	}
}

func featureVoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDscFeaturesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		product     = "dsc"
		httpUrl     = "v1/{project_id}/features"
		limit       = 1000
		offset      = 0
		allFeatures = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	for {
		currentPath := fmt.Sprintf("%s?offset=%d&limit=%d", requestPath, offset, limit)
		requestOpts := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		resp, err := client.Request("GET", currentPath, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving DSC features: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		featuresRaw := utils.PathSearch("features", respBody, make([]interface{}, 0)).([]interface{})
		if len(featuresRaw) == 0 {
			break
		}

		allFeatures = append(allFeatures, featuresRaw...)

		offset += len(featuresRaw)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenFeatures(allFeatures)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenFeatures(features []interface{}) []interface{} {
	if len(features) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(features))
	for _, item := range features {
		result = append(result, map[string]interface{}{
			"description": utils.PathSearch("description", item, nil),
			"enabled":     utils.PathSearch("enabled", item, nil),
			"name":        utils.PathSearch("name", item, nil),
		})
	}

	return result
}
