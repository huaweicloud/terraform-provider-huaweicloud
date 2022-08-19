package rds

import (
	"context"
	"strconv"

	"github.com/chnsz/golangsdk/openstack/rds/v3/flavors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func DataSourceRdsFlavor() *schema.Resource {
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
				ValidateFunc: validation.StringInSlice([]string{
					"MySQL", "PostgreSQL", "SQLServer",
				}, true),
			},
			"db_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ha", "single", "replica",
				}, false),
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
					},
				},
			},
		},
	}
}

func dataSourceRdsFlavorRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud rds client: %s", err)
	}

	dbType := d.Get("db_type").(string)
	listOpts := flavors.DbFlavorsOpts{Versionname: d.Get("db_version").(string)}

	pages, err := flavors.List(client, listOpts, dbType).AllPages()
	if err != nil {
		return fmtp.DiagErrorf(err.Error())
	}

	flavorsResp, err := flavors.ExtractDbFlavors(pages)
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve RDS flavors: %s", err)
	}

	mode := d.Get("instance_mode").(string)
	az := d.Get("availability_zone").(string)
	version := d.Get("db_version").(string)
	groupType := d.Get("group_type").(string)

	var vcpus string
	if v, ok := d.GetOk("vcpus"); ok {
		vcpus = strconv.Itoa(v.(int))
	}

	filter := map[string]interface{}{
		"Vcpus":        vcpus,
		"Instancemode": mode,
		"GroupType":    groupType,
	}
	if mem, ok := d.GetOk("memory"); ok {
		filter["Ram"] = mem.(int)
	}

	filterFlavors, err := utils.FilterSliceWithField(flavorsResp.Flavorslist, filter)

	if err != nil {
		return fmtp.DiagErrorf("filter RDS flavors failed: %s", err)
	}

	var resultFlavors []interface{}
	var ids []string

	for _, item := range filterFlavors {
		flavor := item.(flavors.Flavors)

		// filter availability_zones
		var azList []string
		for k, v := range flavor.Azstatus {
			if v == "normal" && (az == "" || (az != "" && az == k)) {
				azList = append(azList, k)
			}
		}

		// filter db_versions
		var versionList []string
		for _, v := range flavor.VersionName {
			if version == "" || (version != "" && version == v) {
				versionList = append(versionList, v)
			}
		}

		if len(azList) > 0 && len(versionList) > 0 {
			resultFlavors = append(resultFlavors, flattenRdsFlavor(flavor, azList, versionList))
		}

		ids = append(ids, flavor.ID)

	}

	logp.Printf("[DEBUG]RDS flavors api return:%d, after filter: %d, %v", len(flavorsResp.Flavorslist), len(resultFlavors), resultFlavors)

	mErr := d.Set("flavors", resultFlavors)
	if mErr != nil {
		return fmtp.DiagErrorf("set flavors err:%s", mErr)
	}

	d.SetId(hashcode.Strings(ids))

	return nil
}

func flattenRdsFlavor(flavor flavors.Flavors, azList, versionList []string) map[string]interface{} {
	vcpus, _ := strconv.Atoi(flavor.Vcpus)

	return map[string]interface{}{
		"id":                 flavor.ID,
		"name":               flavor.Speccode,
		"vcpus":              vcpus,
		"memory":             flavor.Ram,
		"group_type":         flavor.GroupType,
		"instance_mode":      flavor.Instancemode,
		"mode":               flavor.Instancemode,
		"availability_zones": azList,
		"db_versions":        versionList,
	}
}
