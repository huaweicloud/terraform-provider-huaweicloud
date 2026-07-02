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

// @API DSC GET /v1/{project_id}/asset-center/database/instances
func DataSourceDatabaseInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     databaseInstancesSchema(),
			},
		},
	}
}

func databaseInstancesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
			"ins_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ins_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bind_database": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_external": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildDatabaseInstancesQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?instance_type=%s", d.Get("instance_type").(string))
}

func dataSourceDatabaseInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/asset-center/database/instances"
		product = "dsc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildDatabaseInstancesQueryParams(d)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	resp, err := client.Request("GET", listPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC database instances: %s", err)
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
		d.Set("instances", flattenDatabaseInstances(utils.PathSearch(
			"instance_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDatabaseInstances(instances []interface{}) []interface{} {
	if len(instances) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(instances))
	for _, v := range instances {
		rst = append(rst, utils.RemoveNil(map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"ins_id":            utils.PathSearch("ins_id", v, nil),
			"ins_name":          utils.PathSearch("ins_name", v, nil),
			"ins_status":        utils.PathSearch("ins_status", v, nil),
			"ins_type":          utils.PathSearch("ins_type", v, nil),
			"db_type":           utils.PathSearch("db_type", v, nil),
			"version":           utils.PathSearch("version", v, nil),
			"address":           utils.PathSearch("address", v, nil),
			"port":              utils.PathSearch("port", v, nil),
			"bind_database":     utils.PathSearch("bind_database", v, nil),
			"is_external":       utils.PathSearch("is_external", v, nil),
			"project_id":        utils.PathSearch("project_id", v, nil),
			"vpc_id":            utils.PathSearch("vpc_id", v, nil),
			"subnet_id":         utils.PathSearch("subnet_id", v, nil),
			"security_group_id": utils.PathSearch("security_group_id", v, nil),
			"create_time":       utils.PathSearch("create_time", v, nil),
		}))
	}

	return rst
}
