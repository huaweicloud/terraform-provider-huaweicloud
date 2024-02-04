// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CNAD
// ---------------------------------------------------------------

package cnad

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AAD GET /v1/cnad/protected-ips
func DataSourceProtectedObjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceProtectedObjectsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the CNAD advanced instance ID.`,
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the CNAD advanced policy ID.`,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the CNAD advanced protected object IP.`,
			},
			"protected_object_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the protected object ID.`,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Specifies the type of the protected object. Valid values are **VPN**, **NAT**, **VIP**,
**CCI**, **EIP**, **ELB**, **REROUTING_IP**, **SubEni** and **NetInterFace**`,
			},
			"is_unbound": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether query protected objects which policies unbound.`,
			},
			"protected_objects": {
				Type:        schema.TypeList,
				Elem:        dataSourceProtectedObjectsSchema(),
				Computed:    true,
				Description: `Indicates the list of the advanced protected objects.`,
			},
		},
	}
}

func dataSourceProtectedObjectsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the protected object.`,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the IP of the protected object.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the type of the protected object.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the protected object.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance ID of the protected object.`,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance name of the protected object.`,
			},
			"instance_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance version of the protected object.`,
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the region to which the protected object belongs.`,
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the status of the protected object.`,
			},
			"block_threshold": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the block threshold of the protected object.`,
			},
			"clean_threshold": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the clean threshold of the protected object.`,
			},
			"policy_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the policy name of the protected object.`,
			},
		},
	}
	return &sc
}

func resourceProtectedObjectsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                        = meta.(*config.Config)
		region                     = cfg.GetRegion(d)
		mErr                       *multierror.Error
		getProtectedObjectsHttpUrl = "v1/cnad/protected-ips"
		getProtectedObjectsProduct = "aad"
	)
	getProtectedObjectsClient, err := cfg.NewServiceClient(getProtectedObjectsProduct, region)
	if err != nil {
		return diag.Errorf("error creating CNAD Client: %s", err)
	}

	getProtectedObjectsPath := getProtectedObjectsClient.Endpoint + getProtectedObjectsHttpUrl
	getProtectedObjectsPath += buildGetProtectedObjectsQueryParams(d)

	getProtectedObjectsResp, err := pagination.ListAllItems(
		getProtectedObjectsClient,
		"offset",
		getProtectedObjectsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving protected objects, %s", err)
	}

	getProtectedObjectsRespJson, err := json.Marshal(getProtectedObjectsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getProtectedObjectsRespBody interface{}
	err = json.Unmarshal(getProtectedObjectsRespJson, &getProtectedObjectsRespBody)
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
		d.Set("protected_objects", flattenGetProtectedObjectsResponseBody(getProtectedObjectsRespBody, d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetProtectedObjectsResponseBody(resp interface{}, d *schema.ResourceData) []interface{} {
	if resp == nil {
		return nil
	}
	objectID := d.Get("protected_object_id").(string)
	objectType := d.Get("type").(string)
	isUnbound := d.Get("is_unbound").(bool)

	curJson := utils.PathSearch("items", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		idResp := utils.PathSearch("id", v, "")
		if objectID != "" && objectID != idResp {
			continue
		}

		typeResp := utils.PathSearch("type", v, "")
		if objectType != "" && objectType != typeResp {
			continue
		}

		// policy_name not empty mean the protected object has been bound to a policy already
		policyNameResp := utils.PathSearch("policy_name", v, "")
		if isUnbound && policyNameResp != "" {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"id":               idResp,
			"type":             typeResp,
			"policy_name":      policyNameResp,
			"ip_address":       utils.PathSearch("ip", v, nil),
			"name":             utils.PathSearch("name", v, nil),
			"instance_id":      utils.PathSearch("package_id", v, nil),
			"instance_name":    utils.PathSearch("package_name", v, nil),
			"instance_version": utils.PathSearch("package_version", v, nil),
			"region":           utils.PathSearch("region", v, nil),
			"status":           utils.PathSearch("status", v, nil),
			"block_threshold":  utils.PathSearch("block_threshold", v, nil),
			"clean_threshold":  utils.PathSearch("clean_threshold", v, nil),
		})
	}
	return rst
}

func buildGetProtectedObjectsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&package_id=%v", res, v)
	}

	if v, ok := d.GetOk("policy_id"); ok {
		res = fmt.Sprintf("%s&policy_id=%v", res, v)
	}

	if v, ok := d.GetOk("ip_address"); ok {
		res = fmt.Sprintf("%s&ip=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
