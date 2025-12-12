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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"distinct_shared_resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     distinctSharedResourceSchema(),
			},
		},
	}
}

func distinctSharedResourceSchema() *schema.Resource {
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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listDistinctSharedResourcesHttpUrl = "v1/shared-resources/search-distinct-resource"
		listDistinctSharedResourcesProduct = "ram"
	)
	ramClient, err := cfg.NewServiceClient(listDistinctSharedResourcesProduct, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	listDistinctSharedResourcesHttpPath := ramClient.Endpoint + listDistinctSharedResourcesHttpUrl

	var distinctSharedResources []interface{}
	var marker string
	var limit = 200

	for {
		listDistinctSharedResourcesHttpOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildDistinctSharedResourcesBody(d, limit, marker),
		}
		listDistinctSharedResourcesHttpResp, err := ramClient.Request(
			"POST",
			listDistinctSharedResourcesHttpPath,
			&listDistinctSharedResourcesHttpOpt,
		)
		if err != nil {
			return diag.FromErr(err)
		}
		listDistinctSharedResourcesRespBody, err := utils.FlattenResponse(listDistinctSharedResourcesHttpResp)
		if err != nil {
			return diag.FromErr(err)
		}

		onePageSearchDistinctResource := flattenDistinctSharedResourcesResp(listDistinctSharedResourcesRespBody)
		distinctSharedResources = append(distinctSharedResources, onePageSearchDistinctResource...)
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
		d.Set("distinct_shared_resources", distinctSharedResources),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDistinctSharedResourcesBody(d *schema.ResourceData, limit int, marker string) map[string]interface{} {
	params := make(map[string]interface{})

	params["limit"] = limit
	if marker != "" {
		params["marker"] = marker
	}

	if v, ok := d.GetOk("resource_ids"); ok {
		params["resource_ids"] = v
	}

	if v, ok := d.GetOk("principal"); ok {
		params["principal"] = v
	}

	if v, ok := d.GetOk("resource_region"); ok {
		params["resource_region"] = v
	}

	if v, ok := d.GetOk("resource_urns"); ok {
		params["resource_urns"] = v
	}

	params["resource_owner"] = d.Get("resource_owner")

	return params
}

func flattenDistinctSharedResourcesResp(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("distinct_shared_resources", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"resource_urn":  utils.PathSearch("resource_urn", v, nil),
			"resource_type": utils.PathSearch("resource_type", v, nil),
			"updated_at":    utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}
