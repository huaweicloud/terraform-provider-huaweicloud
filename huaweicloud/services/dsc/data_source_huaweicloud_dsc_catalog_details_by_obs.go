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

// @API DSC GET /v1/{project_id}/metadata/catalog/column-details/obs-dim
func DataSourceCatalogDetailsByObs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCatalogDetailsByObsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     catalogDetailsObsDimSchema(),
			},
		},
	}
}

func catalogDetailsObsDimSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"classification_tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_level_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCatalogDetailsByObsQueryParams(d *schema.ResourceData) string {
	if v, ok := d.GetOk("type_id"); ok {
		return fmt.Sprintf("?type_id=%v", v)
	}

	return ""
}

func dataSourceCatalogDetailsByObsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/metadata/catalog/column-details/obs-dim"
		product = "dsc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	currentPath := requestPath + buildCatalogDetailsByObsQueryParams(d)
	resp, err := client.Request("GET", currentPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC catalog details by OBS: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	results := utils.PathSearch("results", respBody, make([]interface{}, 0)).([]interface{})

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("results", flattenCatalogDetailsObsDimList(results)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCatalogDetailsObsDimList(results []interface{}) []interface{} {
	if len(results) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(results))
	for _, v := range results {
		rst = append(rst, map[string]interface{}{
			"bucket_name": utils.PathSearch("bucket_name", v, nil),
			"classification_tags": utils.ExpandToStringList(
				utils.PathSearch("classification_tags", v, make([]interface{}, 0)).([]interface{})),
			"security_level_name": utils.PathSearch("security_level_name", v, nil),
		})
	}

	return rst
}
