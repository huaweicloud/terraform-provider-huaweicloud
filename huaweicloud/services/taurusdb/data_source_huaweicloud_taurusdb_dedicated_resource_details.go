package taurusdb

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TaurusDB GET /v3/{project_id}/dedicated-resource/{dedicated_resource_id}
func DataSourceTaurusDBDedicatedResourceDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBDedicatedResourceDetailsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dedicated_resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"architecture": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dedicated_compute_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dedicatedResourceDetailComputeInfoSchema(),
			},
			"dedicated_storage_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dedicatedResourceDetailStorageInfoSchema(),
			},
		},
	}
}

func dedicatedResourceDetailComputeInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vcpus_total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vcpus_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ram_total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ram_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dedicatedResourceDetailStorageInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceTaurusDBDedicatedResourceDetailsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/dedicated-resource/{dedicated_resource_id}"
	)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{dedicated_resource_id}", d.Get("dedicated_resource_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving TaurusDB dedicated resource detail: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("resource_name", utils.PathSearch("resource_name", getRespBody, nil)),
		d.Set("engine_name", utils.PathSearch("engine_name", getRespBody, nil)),
		d.Set("availability_zone_ids", utils.PathSearch("availability_zone_ids", getRespBody, nil)),
		d.Set("architecture", utils.PathSearch("architecture", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("dedicated_compute_info", flattenDedicatedResourceDetailComputeInfoBody(getRespBody)),
		d.Set("dedicated_storage_info", flattenDedicatedResourceDetailStorageInfoBody(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDedicatedResourceDetailComputeInfoBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("dedicated_compute_info", resp, nil)
	if curJson == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"vcpus_total": utils.PathSearch("vcpus_total", curJson, nil),
			"vcpus_used":  utils.PathSearch("vcpus_used", curJson, nil),
			"ram_total":   utils.PathSearch("ram_total", curJson, nil),
			"ram_used":    utils.PathSearch("ram_used", curJson, nil),
			"spec_code":   utils.PathSearch("spec_code", curJson, nil),
			"host_num":    utils.PathSearch("host_num", curJson, nil),
		},
	}
}

func flattenDedicatedResourceDetailStorageInfoBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("dedicated_storage_info", resp, nil)
	if curJson == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"spec_code": utils.PathSearch("spec_code", curJson, nil),
			"host_num":  utils.PathSearch("host_num", curJson, nil),
		},
	}
}
