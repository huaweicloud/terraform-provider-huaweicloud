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

// The paging parameters mentioned in the API documentation are meaningless here, and the paging function is ignored.

// @API CBR GET /v3/{project_id}/cbr-features
func DataSourceCbrFeatures() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCbrFeaturesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the region where the CBR features are located.`,
			},
			"feature_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The feature values in JSON format.`,
			},
		},
	}
}

func dataSourceCbrFeaturesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/cbr-features"
		product = "cbr"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CBR features: %s", err)
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

	// The response data is a Map structure, so it needs to be converted to a JSON string.
	jsonString, err := json.Marshal(respBody)
	if err != nil {
		return diag.Errorf("error marshaling CBR features: %s", err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("feature_value", string(jsonString)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
