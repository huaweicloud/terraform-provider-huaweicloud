package rds

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS GET /v3/{project_id}/flavors/{database_name}
func DataSourceRdsFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsFlavorRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"db_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vcpus": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"group_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_flexus": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"flavors": {
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
						"vcpus": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "use instance_mode instead",
						},
						"instance_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"db_versions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"az_status": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceRdsFlavorRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/flavors/{database_name}"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{database_name}", d.Get("db_type").(string))

	getFlavorsQueryParams := buildGetFlavorsQueryParams(d)
	getPath += getFlavorsQueryParams

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving RDS flavors: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("flavors", flattenGetFlavors(getRespBody, d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetFlavorsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("db_version"); ok {
		res = fmt.Sprintf("%s&version_name=%v", res, v)
	}
	if v, ok := d.GetOk("is_flexus"); ok {
		res = fmt.Sprintf("%s&is_flexus=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenGetFlavors(resp interface{}, d *schema.ResourceData) []interface{} {
	if resp == nil {
		return nil
	}

	instanceMode, instanceModeOk := d.GetOk("instance_mode")
	availabilityZone, availabilityZoneOk := d.GetOk("availability_zone")
	groupType, groupTypeOk := d.GetOk("group_type")
	vcpus, vcpusOk := d.GetOk("vcpus")
	memory, memoryOk := d.GetOk("memory")

	curJson := utils.PathSearch("flavors", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0)
	for _, v := range curArray {
		instanceModeRaw := utils.PathSearch("instance_mode", v, nil)
		if instanceModeOk && instanceMode != instanceModeRaw {
			continue
		}
		groupTypeRaw := utils.PathSearch("group_type", v, nil)
		if groupTypeOk && groupType != groupTypeRaw {
			continue
		}
		vcpusRaw := utils.PathSearch("vcpus", v, nil)
		if vcpusOk && fmt.Sprint(vcpus) != fmt.Sprint(vcpusRaw) {
			continue
		}
		memoryRaw := utils.PathSearch("ram", v, nil)
		if memoryOk && fmt.Sprint(memory) != fmt.Sprint(memoryRaw) {
			continue
		}
		azStatusRaw := utils.PathSearch("az_status", v, make(map[string]interface{})).(map[string]interface{})
		azList := make([]string, 0)
		for az, status := range azStatusRaw {
			if status == "normal" && (!availabilityZoneOk || availabilityZone == az) {
				azList = append(azList, az)
			}
		}

		vcpusValue, _ := strconv.Atoi(vcpusRaw.(string))
		rst = append(rst, map[string]interface{}{
			"id":                 utils.PathSearch("id", v, nil),
			"name":               utils.PathSearch("spec_code", v, nil),
			"vcpus":              vcpusValue,
			"memory":             memoryRaw,
			"group_type":         groupTypeRaw,
			"mode":               instanceModeRaw,
			"instance_mode":      instanceModeRaw,
			"availability_zones": azList,
			"db_versions":        utils.PathSearch("version_name", v, nil),
			"az_status":          utils.PathSearch("az_status", v, nil),
		})
	}
	return rst
}
