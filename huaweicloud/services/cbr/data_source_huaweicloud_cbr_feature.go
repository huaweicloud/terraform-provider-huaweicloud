package cbr

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CBR GET /v3/{project_id}/cbr-features/{feature_key}
func DataSourceCbrFeature() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCbrFeatureRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the region where the CBR feature is located.`,
			},
			"feature_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the key of the feature to query.`,
			},
			"feature_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the specified feature in JSON format.`,
			},
		},
	}
}

func dataSourceCbrFeatureRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/cbr-features/{feature_key}"
		product = "cbr"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{feature_key}", d.Get("feature_key").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CBR feature: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	jsonString, err := json.Marshal(respBody)
	if err != nil {
		return diag.Errorf("error marshaling CBR feature: %s", err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("feature_value", string(jsonString)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
