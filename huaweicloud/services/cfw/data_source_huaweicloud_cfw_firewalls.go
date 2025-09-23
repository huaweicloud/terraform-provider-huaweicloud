// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CFW
// ---------------------------------------------------------------

package cfw

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CFW GET /v1/{project_id}/firewall/exist
func DataSourceFirewalls() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceFirewallsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fw_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `|-
                    The firewall instance ID.`,
			},
			"service_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: `|-
                    Service type
                      0. North-south firewall
                      1. East-west firewall`,
			},
			"records": {
				Type:        schema.TypeList,
				Elem:        firewallsGetFirewallInstanceResponseRecordSchema(),
				Computed:    true,
				Description: `The firewall instance records.`,
			},
		},
	}
}

func firewallsGetFirewallInstanceResponseRecordSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"charge_mode": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Billing mode. The value can be 0 (yearly/monthly) or 1 (pay-per-use).`,
			},
			"engine_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Engine type`,
			},
			"feature_toggle": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeBool},
				Computed:    true,
				Description: `Whether to enable the feature. The options are true (yes) and false (no).`,
			},
			"flavor": {
				Type:        schema.TypeList,
				Elem:        firewallsGetFirewallInstanceResponseRecordFlavorSchema(),
				Computed:    true,
				Description: `The flavor of the firewall.`,
			},
			"fw_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Firewall ID`,
			},
			"ha_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Cluster type`,
			},
			"is_old_firewall_instance": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the engine is an old engine. The options are true (yes) and false (no).`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Firewall name`,
			},
			"protect_objects": {
				Type:        schema.TypeList,
				Elem:        firewallsGetFirewallInstanceResponseRecordProtectObjectVOSchema(),
				Computed:    true,
				Description: `Project list`,
			},
			"resources": {
				Type:        schema.TypeList,
				Elem:        firewallsGetFirewallInstanceResponseRecordFirewallInstanceResourceSchema(),
				Computed:    true,
				Description: `Firewall instance resources`,
			},
			"service_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Service type`,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
				//nolint:revive
				Description: `Firewall status list. The options are as follows: -1: waiting for payment; 0: creating; 1: deleting; 2: running; 3: upgrading; 4: deletion completed; 5: freezing; 6: creation failed; 7: deletion failed; 8: freezing failed; 9: storage in progress; 10: storage failed; 11: upgrade failed`,
			},
			"support_ipv6": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether IPv6 is supported. The options are true (yes) and false (no).`,
			},
		},
	}
	return &sc
}

func firewallsGetFirewallInstanceResponseRecordFlavorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Bandwidth`,
			},
			"eip_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Number of EIPs`,
			},
			"log_storage": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Log storage`,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
				//nolint:revive
				Description: `Firewall version. The value can be 0 (standard edition), 1 (professional edition), 2 (platinum edition), or 3 (basic edition).`,
			},
			"vpc_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Number of VPCs`,
			},
		},
	}
	return &sc
}

func firewallsGetFirewallInstanceResponseRecordProtectObjectVOSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"object_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Protected object ID`,
			},
			"object_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Protected object name`,
			},
			"type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Project type. The options are as follows: 0: north-south; 1: east-west.`,
			},
		},
	}
	return &sc
}

func firewallsGetFirewallInstanceResponseRecordFirewallInstanceResourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"cloud_service_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Service type, which is used by CBC. The value is hws.service.type.cfw.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Resource ID`,
			},
			"resource_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Resource quantity`,
			},
			"resource_size_measure_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Resource unit name`,
			},
			"resource_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Inventory unit code`,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
				//nolint:revive
				Description: `Resource type. The options are as follows:1. CFW: hws.resource.type.cfw 2. EIP:hws.resource.type.cfw.exp.eip 3. Bandwidth: hws.resource.type.cfw.exp.bandwidth 4. VPC: hws.resource.type.cfw.exp.vpc 5. Log storage: hws.resource.type.cfw.exp.logaudit`,
			},
		},
	}
	return &sc
}

func resourceFirewallsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	// listFirewalls: Query the List of CFW firewalls
	var (
		listFirewallsHttpUrl = "v1/{project_id}/firewall/exist"
		listFirewallsProduct = "cfw"
	)
	listFirewallsClient, err := conf.NewServiceClient(listFirewallsProduct, region)
	if err != nil {
		return diag.Errorf("error creating Firewalls Client: %s", err)
	}

	listFirewallsPath := listFirewallsClient.Endpoint + listFirewallsHttpUrl
	listFirewallsPath = strings.ReplaceAll(listFirewallsPath, "{project_id}", listFirewallsClient.ProjectID)

	listFirewallsqueryParams := buildListFirewallsQueryParams(d)
	listFirewallsPath += listFirewallsqueryParams

	listFirewallsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	listFirewallsResp, err := listFirewallsClient.Request("GET", listFirewallsPath, &listFirewallsOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Firewalls")
	}

	listFirewallsRespBody, err := utils.FlattenResponse(listFirewallsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("records", flattenListFirewallsBodyGetFirewallInstanceResponseRecord(listFirewallsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListFirewallsBodyGetFirewallInstanceResponseRecord(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("data.records", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"charge_mode":              utils.PathSearch("charge_mode", v, nil),
			"engine_type":              utils.PathSearch("engine_type", v, nil),
			"feature_toggle":           utils.PathSearch("feature_toggle", v, nil),
			"flavor":                   flattenGetFirewallInstanceResponseRecordFlavor(v),
			"fw_instance_id":           utils.PathSearch("fw_instance_id", v, nil),
			"ha_type":                  utils.PathSearch("ha_type", v, nil),
			"is_old_firewall_instance": utils.PathSearch("is_old_firewall_instance", v, nil),
			"name":                     utils.PathSearch("name", v, nil),
			"protect_objects":          flattenGetFirewallInstanceResponseRecordProtectObjects(v),
			"resources":                flattenGetFirewallInstanceResponseRecordResources(v),
			"service_type":             utils.PathSearch("service_type", v, nil),
			"status":                   utils.PathSearch("status", v, nil),
			"support_ipv6":             utils.PathSearch("support_ipv6", v, nil),
		})
	}
	return rst
}

func flattenGetFirewallInstanceResponseRecordFlavor(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("flavor", resp, nil)
	if curJson == nil {
		log.Printf("[ERROR] error parsing flavor from response= %#v", resp)
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"bandwidth":   utils.PathSearch("bandwidth", curJson, nil),
			"eip_count":   utils.PathSearch("eip_count", curJson, nil),
			"log_storage": utils.PathSearch("log_storage", curJson, nil),
			"version":     utils.PathSearch("version", curJson, nil),
			"vpc_count":   utils.PathSearch("vpc_count", curJson, nil),
		},
	}
	return rst
}

func flattenGetFirewallInstanceResponseRecordProtectObjects(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("protect_objects", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"object_id":   utils.PathSearch("object_id", v, nil),
			"object_name": utils.PathSearch("object_name", v, nil),
			"type":        utils.PathSearch("type", v, nil),
		})
	}
	return rst
}

func flattenGetFirewallInstanceResponseRecordResources(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("resources", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"cloud_service_type":       utils.PathSearch("cloud_service_type", v, nil),
			"resource_id":              utils.PathSearch("resource_id", v, nil),
			"resource_size":            utils.PathSearch("resource_size", v, nil),
			"resource_size_measure_id": utils.PathSearch("resource_size_measure_id", v, nil),
			"resource_spec_code":       utils.PathSearch("resource_spec_code", v, nil),
			"resource_type":            utils.PathSearch("resource_type", v, nil),
		})
	}
	return rst
}

func buildListFirewallsQueryParams(d *schema.ResourceData) string {
	res := "?offset=0&limit=10"
	res = fmt.Sprintf("%s&service_type=%v", res, d.Get("service_type"))

	if v, ok := d.GetOk("fw_instance_id"); ok {
		res = fmt.Sprintf("%s&fw_instance_id=%v", res, v)
	}

	return res
}
