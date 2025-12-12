package ram

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RAM POST /v1/shared-principals/search-distinct-principal
func DataSourceDistinctSharedPrincipals() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDistinctSharedPrincipalsRead,
		Schema: map[string]*schema.Schema{
			"resource_owner": {
				Type:     schema.TypeString,
				Required: true,
			},
			"principals": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resource_urn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"distinct_shared_principals": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     distinctSharedPrincipalSchema(),
			},
		},
	}
}

func distinctSharedPrincipalSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourceDistinctSharedPrincipalsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listDistinctSharedPrincipalsHttpUrl = "v1/shared-principals/search-distinct-principal"
		listDistinctSharedPrincipalsProduct = "ram"
	)
	listDistinctSharedPrincipalsClient, err := cfg.NewServiceClient(listDistinctSharedPrincipalsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	listDistinctSharedPrincipalsHttpPath := listDistinctSharedPrincipalsClient.Endpoint + listDistinctSharedPrincipalsHttpUrl

	var distinctSharedPrincipals []interface{}
	var marker string
	var limit = 200

	for {
		listDistinctSharedPrincipalsHttpOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildDistinctSharedPrincipalsBody(d, limit, marker),
		}
		listDistinctSharedPrincipalsHttpResp, err := listDistinctSharedPrincipalsClient.Request(
			"POST",
			listDistinctSharedPrincipalsHttpPath,
			&listDistinctSharedPrincipalsHttpOpt,
		)
		if err != nil {
			return diag.Errorf("error retrieving RAM distinct shared principles, %s", err)
		}
		listDistinctSharedPrincipalsRespBody, err := utils.FlattenResponse(listDistinctSharedPrincipalsHttpResp)
		if err != nil {
			return diag.FromErr(err)
		}

		onePageSearchDistinctResource := flattenDistinctSharedPrincipalsResp(listDistinctSharedPrincipalsRespBody)
		distinctSharedPrincipals = append(distinctSharedPrincipals, onePageSearchDistinctResource...)
		marker = utils.PathSearch("page_info.next_marker", onePageSearchDistinctResource, "").(string)
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
		d.Set("distinct_shared_principals", distinctSharedPrincipals),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDistinctSharedPrincipalsBody(d *schema.ResourceData, limit int, marker string) map[string]interface{} {
	params := make(map[string]interface{})

	params["limit"] = limit
	if marker != "" {
		params["marker"] = marker
	}

	if v, ok := d.GetOk("principals"); ok {
		params["principals"] = v
	}

	if v, ok := d.GetOk("resource_urn"); ok {
		params["resource_urn"] = v
	}

	if v, ok := d.GetOk("resource_owner"); ok {
		params["resource_owner"] = v
	}

	return params
}

func flattenDistinctSharedPrincipalsResp(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("distinct_shared_principals", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"updated_at": utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}
