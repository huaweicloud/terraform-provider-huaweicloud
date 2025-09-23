package sfsturbo

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/sfs_turbo/v1/shares"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SFSTurbo GET /v1/{project_id}/sfs-turbo/shares/detail
func DataSourceTurbos() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTurbosRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"share_proto": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"share_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"turbos": {
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
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"share_proto": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"share_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enhanced": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"available_capacity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"export_location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"crypt_key_id": {
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
					},
				},
			},
		},
	}
}

func filterTurbosBySize(turbos []interface{}, size int) []interface{} {
	if len(turbos) < 1 {
		return nil
	}
	result := make([]interface{}, 0, len(turbos))
	for _, val := range turbos {
		if turbo, ok := val.(shares.Turbo); ok {
			re := regexp.MustCompile(fmt.Sprintf(`^%d+(\.\d+)?$`, size))
			if re.MatchString(turbo.Size) {
				result = append(result, turbo)
			}
		}
	}
	return result
}

func flattenTurbos(turbos []interface{}) ([]map[string]interface{}, []string) {
	if len(turbos) < 1 {
		return nil, nil
	}
	result := make([]map[string]interface{}, len(turbos))
	ids := make([]string, len(turbos))
	for i, val := range turbos {
		if turbo, ok := val.(shares.Turbo); ok {
			rm := map[string]interface{}{
				"id":                    turbo.ID,
				"name":                  turbo.Name,
				"share_proto":           turbo.ShareProto,
				"share_type":            turbo.ShareType,
				"enterprise_project_id": turbo.EnterpriseProjectId,
				"version":               turbo.Version,
				"availability_zone":     turbo.AvailabilityZone,
				"available_capacity":    turbo.AvailCapacity,
				"export_location":       turbo.ExportLocation,
				"crypt_key_id":          turbo.CryptKeyID,
				"vpc_id":                turbo.VpcID,
				"subnet_id":             turbo.SubnetID,
				"security_group_id":     turbo.SecurityGroupID,
			}

			if turbo.ExpandType == "bandwidth" {
				rm["enhanced"] = true
			} else {
				rm["enhanced"] = false
			}

			// High-precision to low-precision, discarding digits after the dot (.).
			floatSize, _ := strconv.ParseFloat(turbo.Size, 64)
			rm["size"] = int(floatSize)

			result[i] = rm
			ids[i] = turbo.ID
		}
	}
	return result, ids
}

func dataSourceTurbosRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	resp, err := shares.List(client)
	if err != nil {
		return diag.Errorf("error getting SFS turbo list: %s", err)
	}

	filter := map[string]interface{}{
		"Name":                d.Get("name"),
		"ShareProto":          d.Get("share_proto"),
		"ShareType":           d.Get("share_type"),
		"EnterpriseProjectId": d.Get("enterprise_project_id"),
	}

	result, err := utils.FilterSliceWithField(resp, filter)
	if err != nil {
		return diag.Errorf("error filtering SFS turbo list: %s", err)
	}
	if size, ok := d.GetOk("size"); ok {
		result = filterTurbosBySize(result, size.(int))
	}
	log.Printf("[DEBUG] the filter result of STS turbo list is: %s", result)

	turbos, ids := flattenTurbos(result)
	d.SetId(hashcode.Strings(ids))
	if err := d.Set("turbos", turbos); err != nil {
		return diag.Errorf("error setting field of SFS turbo list: %s", err)
	}
	return nil
}
