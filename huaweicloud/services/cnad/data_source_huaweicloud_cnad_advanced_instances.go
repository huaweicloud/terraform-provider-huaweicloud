// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CNAD
// ---------------------------------------------------------------

package cnad

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AAD GET /v1/cnad/packages
func DataSourceInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceInstancesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the region in which to query the data source.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the instance id.`,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the instance name.`,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Specifies the instance type. Valid values are **cnad_pro**, **cnad_ip**,
**cnad_ep**, **cnad_full_high**, **cnad_vic** and **cnad_intl_ep**.`,
			},
			"instances": {
				Type:        schema.TypeList,
				Elem:        instanceSchema(),
				Computed:    true,
				Description: `Indicates the list of the Advanced instances`,
			},
		},
	}
}

func instanceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance id.`,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance name.`,
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the region where the instance belongs to.`,
			},
			"instance_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance type of the instance.`,
			},
			"protection_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the protection type of the instance.`,
			},
			"ip_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates ip num of the instance.`,
			},
			"ip_num_now": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the current ip num of the instance.`,
			},
			"protection_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the protection num of the instance, value **9999** means unlimited times.`,
			},
			"protection_num_now": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the current protection num of the instance.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the created time.`,
			},
		},
	}
	return &sc
}

func resourceInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error
	var (
		getAdvancedInstancesHttpUrl = "v1/cnad/packages"
		getAdvancedInstancesProduct = "aad"
	)
	getAdvancedInstancesClient, err := cfg.NewServiceClient(getAdvancedInstancesProduct, region)
	if err != nil {
		return diag.Errorf("error creating CNAD Client: %s", err)
	}

	getAdvancedInstancesPath := getAdvancedInstancesClient.Endpoint + getAdvancedInstancesHttpUrl
	getAdvancedInstancesOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}
	getAdvancedInstancesResp, err := getAdvancedInstancesClient.Request("GET", getAdvancedInstancesPath,
		&getAdvancedInstancesOpt)
	if err != nil {
		return diag.Errorf("error retrieving advanced instances, %s", err)
	}
	getAdvancedInstancesRespBody, err := utils.FlattenResponse(getAdvancedInstancesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("instances", flattenGetInstancesResponseBody(getAdvancedInstancesRespBody, d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetInstancesResponseBody(resp interface{}, d *schema.ResourceData) []interface{} {
	if resp == nil {
		return nil
	}
	instanceID := d.Get("instance_id").(string)
	instanceName := d.Get("instance_name").(string)
	region := d.Get("region").(string)
	instanceType := d.Get("instance_type").(string)

	curJson := utils.PathSearch("items", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		instanceIDResp := utils.PathSearch("package_id", v, "")
		if instanceID != "" && instanceID != instanceIDResp {
			continue
		}

		instanceNameResp := utils.PathSearch("package_name", v, "")
		if instanceName != "" && instanceName != instanceNameResp {
			continue
		}

		regionResp := utils.PathSearch("region_id", v, "")
		if region != "" && region != regionResp {
			continue
		}

		instanceTypeResp := utils.PathSearch("instance_type", v, "")
		if instanceType != "" && instanceType != instanceTypeResp {
			continue
		}

		createdAt := utils.PathSearch("create_time", v, float64(0)).(float64)
		rst = append(rst, map[string]interface{}{
			"instance_id":        instanceIDResp,
			"instance_name":      instanceNameResp,
			"region":             regionResp,
			"instance_type":      instanceTypeResp,
			"protection_type":    utils.PathSearch("protection_type", v, nil),
			"ip_num":             utils.PathSearch("ip_num", v, nil),
			"ip_num_now":         utils.PathSearch("ip_num_now", v, nil),
			"protection_num":     utils.PathSearch("protection_num", v, nil),
			"protection_num_now": utils.PathSearch("protection_num_now", v, nil),
			"created_at":         utils.FormatTimeStampUTC(int64(createdAt) / 1000),
		})
	}
	return rst
}
