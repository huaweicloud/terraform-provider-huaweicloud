package eip

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP GET /v3/{domain_id}/global-eip-segments/support-masks
func DataSourceGlobalEipSegmentSupportMasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGlobalEipSegmentSupportMasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// This field is of boolean type in the API documentation,
			// here it is modified to string to meet different scenarios.
			"page_reverse": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// The `sort_key` field in the API documentation is of type list, which is unnecessary.
			// Here, it is modified to type string.
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The `sort_dir` field in the API documentation is of type list, which is unnecessary.
			// Here, it is modified to type string.
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Named as `id` in the API documentation, now changed to `mask_ids`.
			"mask_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_version": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"mask": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"support_masks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mask": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildGlobalEipSegmentSupportMasksQueryParams(d *schema.ResourceData, marker string) string {
	queryParams := "?limit=2000"

	if raw, ok := d.GetOk("page_reverse"); ok {
		pageReverse, err := strconv.ParseBool(raw.(string))
		if err != nil {
			log.Printf("[ERROR] error parsing 'page_reverse' field to Boolean: %s", err)
		}
		queryParams += fmt.Sprintf("&page_reverse=%v", pageReverse)
	}
	if raw, ok := d.GetOk("fields"); ok {
		for _, f := range raw.([]interface{}) {
			queryParams += fmt.Sprintf("&fields=%s", f)
		}
	}
	if raw, ok := d.GetOk("sort_key"); ok {
		queryParams += fmt.Sprintf("&sort_key=%s", raw)
	}
	if raw, ok := d.GetOk("sort_dir"); ok {
		queryParams += fmt.Sprintf("&sort_dir=%s", raw)
	}
	if raw, ok := d.GetOk("mask_ids"); ok {
		for _, id := range raw.([]interface{}) {
			queryParams += fmt.Sprintf("&id=%s", id)
		}
	}
	if raw, ok := d.GetOk("ip_version"); ok {
		for _, v := range raw.([]interface{}) {
			queryParams += fmt.Sprintf("&ip_version=%s", v)
		}
	}
	if v, ok := d.GetOk("mask"); ok {
		queryParams += fmt.Sprintf("&mask=%s", v)
	}
	if marker != "" {
		queryParams += fmt.Sprintf("&marker=%s", marker)
	}

	return queryParams
}

func dataSourceGlobalEipSegmentSupportMasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "geip"
		httpUrl    = "v3/{domain_id}/global-eip-segments/support-masks"
		result     = make([]interface{}, 0)
		nextMarker string
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", cfg.DomainID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithMarker := requestPath + buildGlobalEipSegmentSupportMasksQueryParams(d, nextMarker)
		resp, err := client.Request("GET", requestPathWithMarker, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving global EIP segment support masks: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		masksResp := utils.PathSearch("support_masks", respBody, make([]interface{}, 0)).([]interface{})
		if len(masksResp) == 0 {
			break
		}

		result = append(result, masksResp...)

		nextMarker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("support_masks", flattenGlobalEipSegmentSupportMasks(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGlobalEipSegmentSupportMasks(masks []interface{}) []interface{} {
	if len(masks) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(masks))
	for _, m := range masks {
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("id", m, nil),
			"ip_version": utils.PathSearch("ip_version", m, nil),
			"mask":       utils.PathSearch("mask", m, nil),
			"created_at": utils.PathSearch("created_at", m, nil),
			"updated_at": utils.PathSearch("updated_at", m, nil),
		})
	}

	return rst
}
