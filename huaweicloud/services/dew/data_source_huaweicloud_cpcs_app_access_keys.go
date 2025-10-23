package dew

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

// @API DEW GET /v1/{project_id}/dew/cpcs/apps/{app_id}/access-keys
func DataSourceCpcsAppAccessKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCpcsAppAccessKeysRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the application ID.`,
			},
			"key_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the access key name filter.`,
			},
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sort attribute.`,
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sort direction.`,
			},
			"access_keys": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the access keys.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The access key ID.`,
						},
						"key_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The access key name.`,
						},
						"access_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The access key AK.`,
						},
						"app_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the application to which the access key belongs.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The access key status.`,
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The creation time of the access key.`,
						},
						"download_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The download time of the access key.`,
						},
						"is_downloaded": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the access key has been downloaded.`,
						},
						"is_imported": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the access key is imported.`,
						},
					},
				},
			},
		},
	}
}

func buildDataSourceCpcsAppAccessKeysQueryParams(d *schema.ResourceData) string {
	rst := ""

	if v, ok := d.GetOk("key_name"); ok {
		rst += fmt.Sprintf("&key_name=%v", v)
	}

	if v, ok := d.GetOk("sort_key"); ok {
		rst += fmt.Sprintf("&sort_key=%v", v)
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		rst += fmt.Sprintf("&sort_dir=%v", v)
	}

	if len(rst) > 0 {
		rst = "?" + rst[1:]
	}

	return rst
}

// When calling the API using `page_num`, the API will report a strange error, so the data source will temporarily ignore paging.
func dataSourceCpcsAppAccessKeysRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/dew/cpcs/apps/{app_id}/access-keys"
		product = "kms"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{app_id}", d.Get("app_id").(string))
	requestPath += buildDataSourceCpcsAppAccessKeysQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DEW CPCS app access keys: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error flattening DEW CPCS app access keys response: %s", err)
	}

	generateId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(generateId)

	accessKeys := utils.PathSearch("result", respBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("access_keys", flattenCpcsAppAccessKeysResponseBody(accessKeys)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCpcsAppAccessKeysResponseBody(accessKeys []interface{}) []interface{} {
	if len(accessKeys) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(accessKeys))
	for _, accessKey := range accessKeys {
		result = append(result, map[string]interface{}{
			"access_key_id": utils.PathSearch("access_key_id", accessKey, nil),
			"key_name":      utils.PathSearch("key_name", accessKey, nil),
			"access_key":    utils.PathSearch("access_key", accessKey, nil),
			"app_name":      utils.PathSearch("app_name", accessKey, nil),
			"status":        utils.PathSearch("status", accessKey, nil),
			"create_time":   utils.PathSearch("create_time", accessKey, nil),
			"download_time": utils.PathSearch("download_time", accessKey, nil),
			"is_downloaded": utils.PathSearch("is_downloaded", accessKey, false),
			"is_imported":   utils.PathSearch("is_imported", accessKey, false),
		})
	}

	return result
}
