package iotda

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

// @API IoTDA GET /v5/iot/{project_id}/device-authorizers
func DataSourceCustomAuthentications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomAuthenticationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"authorizer_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"authorizers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authorizer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"authorizer_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"func_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"func_urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"signing_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"default_authorizer": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cache_enable": {
							Type:     schema.TypeBool,
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

func buildCustomAuthenticationsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("authorizer_name"); ok {
		queryParams = fmt.Sprintf("%s&authorizer_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceCustomAuthenticationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		isDerived      = WithDerivedAuth(cfg, region)
		httpUrl        = "v5/iot/{project_id}/device-authorizers?limit=50"
		allAuthorizers []interface{}
		offset         = 0
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildCustomAuthenticationsQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		listResp, err := client.Request("GET", requestPathWithOffset, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving IoTDA custom authentications: %s", err)
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		authorizers := utils.PathSearch("authorizers", listRespBody, make([]interface{}, 0)).([]interface{})
		if len(authorizers) == 0 {
			break
		}

		allAuthorizers = append(allAuthorizers, authorizers...)
		offset += len(authorizers)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("authorizers", flattenListCustomAuthentications(allAuthorizers)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListCustomAuthentications(authorizers []interface{}) []interface{} {
	if len(authorizers) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(authorizers))
	for _, v := range authorizers {
		result = append(result, map[string]interface{}{
			"authorizer_id":      utils.PathSearch("authorizer_id", v, nil),
			"authorizer_name":    utils.PathSearch("authorizer_name", v, nil),
			"func_name":          utils.PathSearch("func_name", v, nil),
			"func_urn":           utils.PathSearch("func_urn", v, nil),
			"signing_enable":     utils.PathSearch("signing_enable", v, nil),
			"default_authorizer": utils.PathSearch("default_authorizer", v, nil),
			"status":             utils.PathSearch("status", v, nil),
			"cache_enable":       utils.PathSearch("cache_enable", v, nil),
			"create_time":        utils.PathSearch("create_time", v, nil),
			"update_time":        utils.PathSearch("update_time", v, nil),
		})
	}
	return result
}
