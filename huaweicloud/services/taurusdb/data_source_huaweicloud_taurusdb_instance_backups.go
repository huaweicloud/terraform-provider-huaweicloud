package taurusdb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/backups
func DataSourceTaurusDBInstanceBackups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBInstanceBackupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"order_field": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"order_rule"},
			},
			"order_rule": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"order_field"},
			},
			"filter_field": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"filter_content"},
			},
			"filter_content": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"filter_field"},
			},
			"backups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     instanceBackupsSchema(),
			},
		},
	}
}

func instanceBackupsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size_unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"use_detail": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTaurusDBInstanceBackupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/instances/{instance_id}/backups?offset=0&limit=100"
		product = "gaussdb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath += buildGetInstanceBackupsQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving TaurusDB instance (%s) backups: %s", instanceId, err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("backups", flattenGetInstanceBackupsBody(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetInstanceBackupsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("order_field"); ok {
		res = fmt.Sprintf("%s&order_field=%v", res, v)
	}
	if v, ok := d.GetOk("order_rule"); ok {
		res = fmt.Sprintf("%s&order_rule=%v", res, v)
	}
	if v, ok := d.GetOk("filter_field"); ok {
		res = fmt.Sprintf("%s&filter_field=%v", res, v)
	}
	if v, ok := d.GetOk("filter_content"); ok {
		res = fmt.Sprintf("%s&filter_content=%v", res, v)
	}
	return res
}

func flattenGetInstanceBackupsBody(resp interface{}) []interface{} {
	backupsJson := utils.PathSearch("backups", resp, make([]interface{}, 0))
	backupsArray := backupsJson.([]interface{})
	if len(backupsArray) < 1 {
		return nil
	}
	rst := make([]interface{}, 0, len(backupsArray))

	for _, v := range backupsArray {
		rst = append(rst, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"description":   utils.PathSearch("description", v, nil),
			"instance_id":   utils.PathSearch("instance_id", v, nil),
			"instance_name": utils.PathSearch("instance_name", v, nil),
			"name":          utils.PathSearch("name", v, nil),
			"size":          utils.PathSearch("size", v, nil),
			"size_unit":     utils.PathSearch("size_unit", v, nil),
			"status":        utils.PathSearch("status", v, nil),
			"created":       utils.PathSearch("created", v, nil),
			"updated":       utils.PathSearch("updated", v, nil),
			"backup_type":   utils.PathSearch("backup_type", v, nil),
			"backup_level":  utils.PathSearch("backup_level", v, nil),
			"backup_method": utils.PathSearch("backup_method", v, nil),
			"use_detail":    utils.PathSearch("use_detail", v, nil),
			"time_zone":     utils.PathSearch("time_zone", v, nil),
		})
	}
	return rst
}
