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

// @API RAM POST /v1/shared-resources/search-distinct-resource
func DataSourceDistinctSharedResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDistinctSharedResourcesRead,
		Schema: map[string]*schema.Schema{
			"resource_owner": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"principal": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_urns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"distinct_shared_resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     distinctSharedResourcesSchema(),
			},
		},
	}
}

func distinctSharedResourcesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
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

func dataSourceDistinctSharedResourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/shared-resources/search-distinct-resource"
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
			JSONBody:         utils.RemoveNil(buildDistinctSharedResourcesBody(d, limit, marker)),
		}

		resp, err := client.Request("POST", requestPath, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving RAM distinct shared resources: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		resourcesResp := utils.PathSearch(
			"distinct_shared_resources", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, resourcesResp...)
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
		d.Set("distinct_shared_resources", flattenDistinctSharedResources(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDistinctSharedResourcesBody(d *schema.ResourceData, limit int, marker string) map[string]interface{} {
	params := map[string]interface{}{
		"resource_owner": d.Get("resource_owner"),
		"limit":          limit,
	}

	if marker != "" {
		params["marker"] = marker
	}

	if v, ok := d.GetOk("resource_ids"); ok {
		params["resource_ids"] = utils.ExpandToStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("principal"); ok {
		params["principal"] = v
	}

	if v, ok := d.GetOk("resource_region"); ok {
		params["resource_region"] = v
	}

	if v, ok := d.GetOk("resource_urns"); ok {
		params["resource_urns"] = utils.ExpandToStringList(v.([]interface{}))
	}

	return params
}

func flattenDistinctSharedResources(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"resource_urn":  utils.PathSearch("resource_urn", v, nil),
			"resource_type": utils.PathSearch("resource_type", v, nil),
			"updated_at":    utils.PathSearch("updated_at", v, nil),
		})
	}

	return rst
}
