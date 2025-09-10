package dew

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

// @API DEW GET /v2/{project_id}/kms/regions
func DataSourceKMSKeyRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKMSKeyRegionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"regions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of regions supported by cross-regional keys.",
			},
		},
	}
}

func dataSourceKMSKeyRegionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		mErr   *multierror.Error
		offset = 0
		result = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("kms", region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	httpUrl := "v2/{project_id}/kms/regions?limit=50"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		getListPath := fmt.Sprintf("%s&offset=%v", listPath, offset)
		requestResp, err := client.Request("GET", getListPath, &opt)
		if err != nil {
			return diag.Errorf("error retrieving KMS key regions: %s", err)
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}

		regions := utils.PathSearch("regions", respBody, make([]interface{}, 0)).([]interface{})
		if len(regions) == 0 {
			break
		}
		result = append(result, regions...)
		offset += len(regions)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("regions", result),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
