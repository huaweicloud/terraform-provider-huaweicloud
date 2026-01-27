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
				Elem:     &schema.Schema{Type: schema.TypeString},
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
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/shared-principals/search-distinct-principal"
		product = "ram"
		result  []interface{}
		marker  string
		limit   = 2000
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	for {
		requestPath := client.Endpoint + httpUrl
		requestOpts := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildDistinctSharedPrincipalsBodyParams(d, limit, marker)),
		}

		resp, err := client.Request("POST", requestPath, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving RAM distinct shared principles: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		principalsResp := utils.PathSearch(
			"distinct_shared_principals", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, principalsResp...)
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
		d.Set("distinct_shared_principals", flattenDistinctSharedPrincipals(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDistinctSharedPrincipalsBodyParams(d *schema.ResourceData, limit int, marker string) map[string]interface{} {
	params := map[string]interface{}{
		"resource_owner": d.Get("resource_owner"),
		"limit":          limit,
	}

	if marker != "" {
		params["marker"] = marker
	}

	if v, ok := d.GetOk("principals"); ok {
		params["principals"] = utils.ExpandToStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("resource_urn"); ok {
		params["resource_urn"] = v
	}

	return params
}

func flattenDistinctSharedPrincipals(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"updated_at": utils.PathSearch("updated_at", v, nil),
		})
	}

	return rst
}
