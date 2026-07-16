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

// @API DSC GET /v1/{project_id}/sdg/asset/bigdata
func DataSourceDscBigdataAssets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscBigdataAssetsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"assets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     bigdataAssetSchema(),
			},
		},
	}
}

func bigdataAssetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"asset_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorize_fail_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorized": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ds_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ds_authorized": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ds_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ds_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ds_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ds_user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ds_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ins_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ins_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ins_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDscBigdataAssetsQueryParams(d *schema.ResourceData) string {
	queryParams := ""
	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("?type=%v", v)
	}
	return queryParams
}

func dataSourceDscBigdataAssetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/sdg/asset/bigdata"
		product = "dsc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildDscBigdataAssetsQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC bigdata assets: %s", err)
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
		d.Set("assets", flattenBigdataAssets(
			utils.PathSearch("datasources", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBigdataAssets(items []interface{}) []interface{} {
	if len(items) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(items))
	for _, v := range items {
		rst = append(rst, map[string]interface{}{
			"asset_name":            utils.PathSearch("asset_name", v, nil),
			"authorize_fail_reason": utils.PathSearch("authorize_fail_reason", v, nil),
			"authorized":            utils.PathSearch("authorized", v, nil),
			"create_time":           utils.PathSearch("create_time", v, nil),
			"default":               utils.PathSearch("default", v, nil),
			"ds_address":            utils.PathSearch("ds_address", v, nil),
			"ds_authorized":         utils.PathSearch("ds_authorized", v, nil),
			"ds_name":               utils.PathSearch("ds_name", v, nil),
			"ds_port":               utils.PathSearch("ds_port", v, nil),
			"ds_type":               utils.PathSearch("ds_type", v, nil),
			"ds_user":               utils.PathSearch("ds_user", v, nil),
			"ds_version":            utils.PathSearch("ds_version", v, nil),
			"id":                    utils.PathSearch("id", v, nil),
			"ins_id":                utils.PathSearch("ins_id", v, nil),
			"ins_name":              utils.PathSearch("ins_name", v, nil),
			"ins_type":              utils.PathSearch("ins_type", v, nil),
			"region":                utils.PathSearch("region", v, nil),
			"security_group_id":     utils.PathSearch("security_group_id", v, nil),
			"source_type":           utils.PathSearch("source_type", v, nil),
			"subnet_id":             utils.PathSearch("subnet_id", v, nil),
			"vpc_id":                utils.PathSearch("vpc_id", v, nil),
		})
	}

	return rst
}
