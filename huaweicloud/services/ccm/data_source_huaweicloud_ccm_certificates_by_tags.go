package ccm

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCM POST /v3/scm/{resource_instances}/action
func DataSourceCertificatesByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceCertificatesByTagsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_instances": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"without_any_tag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     buildCertificatesByTagsTagsSchema(),
			},
			"tags_any": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     buildCertificatesByTagsTagsSchema(),
			},
			"not_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     buildCertificatesByTagsTagsSchema(),
			},
			"not_tags_any": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     buildCertificatesByTagsTagsSchema(),
			},
			"matches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     buildCertificatesByTagsMatchesSchema(),
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildCertificatesByTagsResourcesSchema(),
			},
		},
	}
}

func buildCertificatesByTagsResourcesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildCertificatesByTagsResourcesTagsSchema(),
			},
			"resource_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_detail": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildCertificatesByTagsResourcesResourceDetailSchema(),
			},
		},
	}
}

func buildCertificatesByTagsResourcesResourceDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cert_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert_brand": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"purchase_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"expired_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"wildcard_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sans": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert_des": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signature_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fail_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"partner_order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"push_support": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cert_issued_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"unsubscribe_support": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"origin_cert_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"renewal_cert_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_renew_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"remain_cert_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"auto_deploy_support": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildCertificatesByTagsResourcesTagsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCertificatesByTagsMatchesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildCertificatesByTagsTagsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"values": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildCertificatesByTagsTagsBodyParams(rawArray []interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, raw := range rawArray {
		rawMap, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"key":    rawMap["key"],
			"values": rawMap["values"],
		})
	}

	return rst
}

func buildCertificatesByTagsMatchesBodyParams(rawArray []interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, raw := range rawArray {
		rawMap, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"key":   rawMap["key"],
			"value": rawMap["value"],
		})
	}

	return rst
}

func buildCertificatesByTagsBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"action":          d.Get("action"),
		"without_any_tag": d.Get("without_any_tag"),
		"tags":            buildCertificatesByTagsTagsBodyParams(d.Get("tags").([]interface{})),
		"tags_any":        buildCertificatesByTagsTagsBodyParams(d.Get("tags_any").([]interface{})),
		"not_tags":        buildCertificatesByTagsTagsBodyParams(d.Get("not_tags").([]interface{})),
		"not_tags_any":    buildCertificatesByTagsTagsBodyParams(d.Get("not_tags_any").([]interface{})),
		"matches":         buildCertificatesByTagsMatchesBodyParams(d.Get("matches").([]interface{})),
	}
}

func datasourceCertificatesByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf              = meta.(*config.Config)
		region            = conf.GetRegion(d)
		httpUrl           = "v3/scm/{resource_instances}/action"
		product           = "scm"
		resourceInstances = d.Get("resource_instances").(string)
		action            = d.Get("action").(string)
		offset            = 0
		limit             = 50
		allResults        []interface{}
		totalCount        = 0
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{resource_instances}", resourceInstances)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	bodyParams := buildCertificatesByTagsBodyParams(d)

	for {
		if action == "filter" {
			bodyParams["limit"] = limit
			bodyParams["offset"] = offset
		}

		requestOpt.JSONBody = utils.RemoveNil(bodyParams)
		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving CCM certificates by tags: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		if action == "count" {
			totalCount = int(utils.PathSearch("total_count", respBody, float64(0)).(float64))
			break
		}

		resources := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
		if len(resources) == 0 {
			break
		}

		allResults = append(allResults, resources...)
		offset += len(resources)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("total_count", totalCount),
		d.Set("resources", flattenCertificatesByTagsResources(allResults)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCertificatesByTagsResources(respArray []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		tagsResp := utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})
		resourceDetailResp := utils.PathSearch("resource_detail", v, nil)

		rst = append(rst, map[string]interface{}{
			"resource_id":     utils.PathSearch("resource_id", v, nil),
			"tags":            flattenCertificatesByTagsTagsResources(tagsResp),
			"resource_name":   utils.PathSearch("resource_name", v, nil),
			"resource_detail": flattenCertificatesByTagsResourceDetailResources(resourceDetailResp),
		})
	}

	return rst
}

func flattenCertificatesByTagsTagsResources(respArray []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}

	return rst
}

func flattenCertificatesByTagsResourceDetailResources(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	resourceDetail := map[string]interface{}{
		"cert_id":               utils.PathSearch("cert_id", respBody, nil),
		"cert_name":             utils.PathSearch("cert_name", respBody, nil),
		"domain":                utils.PathSearch("domain", respBody, nil),
		"cert_type":             utils.PathSearch("cert_type", respBody, nil),
		"cert_brand":            utils.PathSearch("cert_brand", respBody, nil),
		"domain_type":           utils.PathSearch("domain_type", respBody, nil),
		"purchase_period":       utils.PathSearch("purchase_period", respBody, nil),
		"expired_time":          utils.PathSearch("expired_time", respBody, nil),
		"order_status":          utils.PathSearch("order_status", respBody, nil),
		"domain_num":            utils.PathSearch("domain_num", respBody, nil),
		"wildcard_number":       utils.PathSearch("wildcard_number", respBody, nil),
		"sans":                  utils.PathSearch("sans", respBody, nil),
		"cert_des":              utils.PathSearch("cert_des", respBody, nil),
		"signature_algorithm":   utils.PathSearch("signature_algorithm", respBody, nil),
		"fail_reason":           utils.PathSearch("fail_reason", respBody, nil),
		"partner_order_id":      utils.PathSearch("partner_order_id", respBody, nil),
		"push_support":          utils.PathSearch("push_support", respBody, nil),
		"cert_issued_time":      utils.PathSearch("cert_issued_time", respBody, nil),
		"resource_id":           utils.PathSearch("resource_id", respBody, nil),
		"unsubscribe_support":   utils.PathSearch("unsubscribe_support", respBody, nil),
		"enterprise_project_id": utils.PathSearch("enterprise_project_id", respBody, nil),
		"origin_cert_id":        utils.PathSearch("origin_cert_id", respBody, nil),
		"renewal_cert_id":       utils.PathSearch("renewal_cert_id", respBody, nil),
		"auto_renew_status":     utils.PathSearch("auto_renew_status", respBody, nil),
		"remain_cert_number":    utils.PathSearch("remain_cert_number", respBody, nil),
		"auto_deploy_support":   utils.PathSearch("auto_deploy_support", respBody, nil),
	}

	return []interface{}{resourceDetail}
}
