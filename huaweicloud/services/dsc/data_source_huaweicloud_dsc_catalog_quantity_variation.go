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

// @API DSC GET /v1/{project_id}/metadata/catalog/quantity-variation
func DataSourceDscCatalogQuantityVariation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscCatalogQuantityVariationRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"label_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The label ID used to filter the data quantity variation.",
			},
			"type_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type ID used to filter the data quantity variation.",
			},
			"sensitive_number_variation": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The sensitive number variation trend.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"time_axis": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The time axis.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"total_number_variation": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The total number variation trend.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func buildDscCatalogQuantityVariationQueryParams(d *schema.ResourceData) string {
	queryParams := ""
	if v, ok := d.GetOk("label_id"); ok {
		queryParams = fmt.Sprintf("%s&label_id=%v", queryParams, v)
	}

	if v, ok := d.GetOk("type_id"); ok {
		queryParams = fmt.Sprintf("%s&type_id=%v", queryParams, v)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceDscCatalogQuantityVariationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/metadata/catalog/quantity-variation"
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

	currentPath := requestPath + buildDscCatalogQuantityVariationQueryParams(d)

	resp, err := client.Request("GET", currentPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC catalog quantity variation: %s", err)
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

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("sensitive_number_variation",
			utils.PathSearch("sensitive_number_variation", respBody, make([]interface{}, 0))),
		d.Set("time_axis",
			utils.PathSearch("time_axis", respBody, make([]interface{}, 0))),
		d.Set("total_number_variation",
			utils.PathSearch("total_number_variation", respBody, make([]interface{}, 0))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
