package css

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

// @API CSS GET /v1.0/{project_id}/datastore/{datastore_id}/flavors
func DataSourceDatastoreFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatastoreFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"datastore_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"datastore_version_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datastore_id_str": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dbname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavors": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ram": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"typename": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"diskrange": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cond_operation_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cond_operation_az": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"localdisk": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"flavor_type_cn": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"flavor_type_en": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"edge": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"str_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_allow_https": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"model_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"models": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"datastore_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"datastore_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_text_model": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"model_version_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"desc": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"language": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"arch_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildDatastoreFlavorsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("datastore_version_id"); ok {
		queryParams = fmt.Sprintf("%s?datastore_version_id=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceDatastoreFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1.0/{project_id}/datastore/{datastore_id}/flavors"
		datastoreId = d.Get("datastore_id").(string)
	)

	client, err := cfg.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{datastore_id}", datastoreId)
	getPath += buildDatastoreFlavorsQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving the flavors: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("datastore_id_str", utils.PathSearch("id", getRespBody, nil)),
		d.Set("dbname", utils.PathSearch("dbname", getRespBody, nil)),
		d.Set("versions", flattenDatastoreFlavorsVersions(
			utils.PathSearch("versions", getRespBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("model_list", flattenDatastoreFlavorsModelList(
			utils.PathSearch("modelList", getRespBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDatastoreFlavorsVersions(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"id":   utils.PathSearch("id", v, nil),
			"name": utils.PathSearch("name", v, nil),
			"flavors": flattenDatastoreFlavors(
				utils.PathSearch("flavors", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenDatastoreFlavors(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"cpu":                   utils.PathSearch("cpu", v, nil),
			"ram":                   utils.PathSearch("ram", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"region":                utils.PathSearch("region", v, nil),
			"typename":              utils.PathSearch("typename", v, nil),
			"diskrange":             utils.PathSearch("diskrange", v, nil),
			"cond_operation_status": utils.PathSearch("condOperationStatus", v, nil),
			"cond_operation_az":     utils.PathSearch("condOperationAz", v, nil),
			"localdisk":             utils.PathSearch("localdisk", v, nil),
			"flavor_type_cn":        utils.PathSearch("flavorTypeCn", v, nil),
			"flavor_type_en":        utils.PathSearch("flavorTypeEn", v, nil),
			"edge":                  utils.PathSearch("edge", v, nil),
			"str_id":                utils.PathSearch("str_id", v, nil),
			"is_allow_https":        utils.PathSearch("isAllowHttps", v, nil),
		})
	}

	return result
}

func flattenDatastoreFlavorsModelList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"total_size": utils.PathSearch("totalSize", v, nil),
			"models": flattenDatastoreModels(
				utils.PathSearch("models", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenDatastoreModels(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"datastore_type":    utils.PathSearch("datastore_type", v, nil),
			"datastore_version": utils.PathSearch("datastore_version", v, nil),
			"is_text_model":     utils.PathSearch("is_text_model", v, nil),
			"model_version_id":  utils.PathSearch("model_version_id", v, nil),
			"desc":              utils.PathSearch("desc", v, nil),
			"language":          utils.PathSearch("language", v, nil),
			"arch_type":         utils.PathSearch("arch_type", v, nil),
		})
	}

	return result
}
