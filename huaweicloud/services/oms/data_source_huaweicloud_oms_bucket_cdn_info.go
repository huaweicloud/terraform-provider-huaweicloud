package oms

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

// @API OMS POST /v2/{project_id}/objectstorage/buckets/cdn-info
func DataSourceBucketCdnInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBucketCdnInfoRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cloud_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ak": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"sk": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"source_cdn": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Required: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"authentication_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"connection_string": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"prefix": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keys": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"is_same_cloud_type": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_download_available": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"checked_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_etag_matching": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_object_existing": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildBucketCdnInfoParams(d *schema.ResourceData, region string) map[string]interface{} {
	return map[string]interface{}{
		"cloud_type":        d.Get("cloud_type"),
		"bucket":            d.Get("bucket"),
		"ak":                d.Get("ak"),
		"sk":                d.Get("sk"),
		"region":            region,
		"source_cdn":        buildSourceCdnBodyParams(d.Get("source_cdn").([]interface{})),
		"connection_string": utils.ValueIgnoreEmpty(d.Get("connection_string")),
		"app_id":            utils.ValueIgnoreEmpty(d.Get("app_id")),
		"prefix":            buildPrefixBodyParams(d.Get("prefix").([]interface{})),
	}
}

func buildSourceCdnBodyParams(sourceCdnInfo []interface{}) map[string]interface{} {
	if len(sourceCdnInfo) == 0 {
		return nil
	}

	bodyParamsInfo, ok := sourceCdnInfo[0].(map[string]interface{})
	if !ok {
		return nil
	}

	params := map[string]interface{}{
		"authentication_type": utils.PathSearch("authentication_type", bodyParamsInfo, nil),
		"domain":              utils.PathSearch("domain", bodyParamsInfo, nil),
		"protocol":            utils.PathSearch("protocol", bodyParamsInfo, nil),
		"authentication_key":  utils.ValueIgnoreEmpty(utils.PathSearch("authentication_key", bodyParamsInfo, nil)),
	}

	return params
}

func buildPrefixBodyParams(prefixInfo []interface{}) map[string]interface{} {
	if len(prefixInfo) == 0 {
		return nil
	}

	rawParams, ok := prefixInfo[0].(map[string]interface{})
	if !ok {
		return nil
	}

	bodyParams := map[string]interface{}{
		"keys": utils.ExpandToStringList(rawParams["keys"].([]interface{})),
	}

	return bodyParams
}

func dataSourceBucketCdnInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/objectstorage/buckets/cdn-info"
	)

	client, err := cfg.NewServiceClient("oms", region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildBucketCdnInfoParams(d, region)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving the bucket CDN information: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("is_same_cloud_type", utils.PathSearch("is_same_cloud_type", respBody, nil)),
		d.Set("is_download_available", utils.PathSearch("is_download_available", respBody, nil)),
		d.Set("checked_keys", flattenBucketCdnInfo(
			utils.PathSearch("checked_keys", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBucketCdnInfo(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"key":                utils.PathSearch("key", v, nil),
			"is_etag_matching":   utils.PathSearch("is_etag_matching", v, nil),
			"is_object_existing": utils.PathSearch("is_object_existing", v, nil),
		})
	}

	return result
}
