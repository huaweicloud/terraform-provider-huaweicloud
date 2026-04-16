package rfs

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

// @API RFS GET /v1/private-providers
func DataSourcePrivateProviders() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateProvidersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"providers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_agency_urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_agency_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildListPrivateProvidersQueryParams(d *schema.ResourceData, marker string) string {
	rst := ""

	if v, ok := d.GetOk("sort_key"); ok {
		rst += fmt.Sprintf("&sort_key=%s", v.(string))
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		rst += fmt.Sprintf("&sort_dir=%s", v.(string))
	}

	if marker != "" {
		rst += fmt.Sprintf("&marker=%s", marker)
	}

	if rst != "" {
		rst = "?" + rst[1:]
	}

	return rst
}

func dataSourcePrivateProvidersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v1/private-providers"
		allProviders = make([]interface{}, 0)
		nextMarker   string
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	reqUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request UUID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Client-Request-Id": reqUUID},
		KeepResponseBody: true,
	}

	for {
		requestPathWithQueryParams := requestPath + buildListPrivateProvidersQueryParams(d, nextMarker)
		resp, err := client.Request("GET", requestPathWithQueryParams, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving RFS private providers: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		providers := utils.PathSearch("providers", respBody, make([]interface{}, 0)).([]interface{})
		if len(providers) == 0 {
			break
		}

		allProviders = append(allProviders, providers...)

		nextMarker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	d.SetId(reqUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("providers", flattenPrivateProviders(allProviders)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPrivateProviders(providers []interface{}) []interface{} {
	if len(providers) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(providers))
	for _, v := range providers {
		rst = append(rst, map[string]interface{}{
			"provider_id":          utils.PathSearch("provider_id", v, nil),
			"provider_name":        utils.PathSearch("provider_name", v, nil),
			"provider_description": utils.PathSearch("provider_description", v, nil),
			"provider_source":      utils.PathSearch("provider_source", v, nil),
			"provider_agency_urn":  utils.PathSearch("provider_agency_urn", v, nil),
			"provider_agency_name": utils.PathSearch("provider_agency_name", v, nil),
			"create_time":          utils.PathSearch("create_time", v, nil),
			"update_time":          utils.PathSearch("update_time", v, nil),
		})
	}

	return rst
}
