package secmaster

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

// @API SecMaster GET /v1/{project_id}/configurations/dictionaries
func DataSourceConfigurationDictionaries() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConfigurationDictionariesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"success_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dictionariesSchema(),
			},
			"failed_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dictionariesSchema(),
			},
		},
	}
}

func dictionariesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dict_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dict_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dict_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dict_val": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dict_pkey": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dict_pcode": {
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
			"publish_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extension_field": {
				Type:     schema.TypeString, // JSON string
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"language": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceConfigurationDictionariesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/configurations/dictionaries"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving dictionaries: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	successList := utils.PathSearch("success_list", getRespBody, make([]interface{}, 0)).([]interface{})
	failedList := utils.PathSearch("failed_list", getRespBody, make([]interface{}, 0)).([]interface{})

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("success_list", flattenConfigurationDictionaries(successList)),
		d.Set("failed_list", flattenConfigurationDictionaries(failedList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenConfigurationDictionaries(dictionariesResp []interface{}) []interface{} {
	if len(dictionariesResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dictionariesResp))
	for _, v := range dictionariesResp {
		rst = append(rst, map[string]interface{}{
			"id":              utils.PathSearch("id", v, nil),
			"version":         utils.PathSearch("version", v, nil),
			"dict_id":         utils.PathSearch("dict_id", v, nil),
			"dict_key":        utils.PathSearch("dict_key", v, nil),
			"dict_code":       utils.PathSearch("dict_code", v, nil),
			"dict_val":        utils.PathSearch("dict_val", v, nil),
			"dict_pkey":       utils.PathSearch("dict_pkey", v, nil),
			"dict_pcode":      utils.PathSearch("dict_pcode", v, nil),
			"create_time":     utils.PathSearch("create_time", v, nil),
			"update_time":     utils.PathSearch("update_time", v, nil),
			"publish_time":    utils.PathSearch("publish_time", v, nil),
			"scope":           utils.PathSearch("scope", v, nil),
			"description":     utils.PathSearch("description", v, nil),
			"extension_field": utils.JsonToString(utils.PathSearch("extension_field", v, nil)),
			"project_id":      utils.PathSearch("project_id", v, nil),
			"language":        utils.PathSearch("language", v, nil),
		})
	}

	return rst
}
