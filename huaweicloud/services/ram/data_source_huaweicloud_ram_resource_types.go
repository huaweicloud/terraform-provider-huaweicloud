package ram

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RAM GET /v1/resource-types
func DataSourceResourceTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceTypesRead,
		Schema: map[string]*schema.Schema{
			"resource_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceTypesSchema(),
			},
		},
	}
}

func resourceTypesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourceResourceTypesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/resource-types"
		product = "ram"
		marker  = ""
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithMarker := requestPath + buildListResourceTypesQueryParams(marker)
		resp, err := client.Request("GET", requestPathWithMarker, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving RAM resource types: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		typesResp := utils.PathSearch("resource_types", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, typesResp...)
		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("resource_types", flattenResourceTypes(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListResourceTypesQueryParams(marker string) string {
	// the default value of limit is `2000`
	queryParams := "?limit=2000"

	if marker != "" {
		queryParams = fmt.Sprintf("%s&marker=%v", queryParams, marker)
	}

	return queryParams
}

func flattenResourceTypes(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"region_id":     utils.PathSearch("region_id", v, nil),
			"resource_type": utils.PathSearch("resource_type", v, nil),
		})
	}

	return rst
}
