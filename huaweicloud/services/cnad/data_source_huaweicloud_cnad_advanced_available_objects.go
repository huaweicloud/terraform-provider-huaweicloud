// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CNAD
// ---------------------------------------------------------------

package cnad

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AAD GET /v1/cnad/packages/{package_id}/unbound-protected-ips
func DataSourceAvailableProtectedObjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceAvailableProtectedObjectsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the instance id.`,
			},
			"protected_object_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the protected object id.`,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the protected object ip.`,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Specifies the type of the protected object. Valid values are **VPN**, **NAT**, **VIP**,
**CCI**, **EIP**, **ELB**, **REROUTING_IP**, **SubEni** and **NetInterFace**`,
			},
			"protected_objects": {
				Type:        schema.TypeList,
				Elem:        availableProtectedObjectsSchema(),
				Computed:    true,
				Description: `Indicates the list of the advanced available protected object.`,
			},
		},
	}
}

func availableProtectedObjectsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the id of the protected object.`,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ip of the protected object.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the type of the protected object.`,
			},
		},
	}
	return &sc
}

func resourceAvailableProtectedObjectsRead(_ context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getProtectedObjectsHttpUrl = "v1/cnad/packages/{package_id}/unbound-protected-ips"
		getProtectedObjectsProduct = "aad"
	)
	getProtectedObjectsClient, err := cfg.NewServiceClient(getProtectedObjectsProduct, region)
	if err != nil {
		return diag.Errorf("error creating CNAD Client: %s", err)
	}

	getProtectedObjectsPath := getProtectedObjectsClient.Endpoint + getProtectedObjectsHttpUrl
	getProtectedObjectsPath = strings.ReplaceAll(getProtectedObjectsPath, "{package_id}",
		fmt.Sprintf("%v", d.Get("instance_id")))

	getProtectedObjectsResp, err := pagination.ListAllItems(
		getProtectedObjectsClient,
		"offset",
		getProtectedObjectsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving available protected objects, %s", err)
	}

	getProtectedObjectsRespJson, err := json.Marshal(getProtectedObjectsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getProtectedObjectsRespBody interface{}
	if err := json.Unmarshal(getProtectedObjectsRespJson, &getProtectedObjectsRespBody); err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("protected_objects", flattenProtectedObjectsResponseBody(getProtectedObjectsRespBody, d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenProtectedObjectsResponseBody(resp interface{}, d *schema.ResourceData) []interface{} {
	if resp == nil {
		return nil
	}
	objectID := d.Get("protected_object_id").(string)
	objectIp := d.Get("ip_address").(string)
	objectType := d.Get("type").(string)

	curJson := utils.PathSearch("ips", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		idResp := utils.PathSearch("id", v, "")
		if objectID != "" && objectID != idResp {
			continue
		}

		ipResp := utils.PathSearch("ip", v, "")
		if objectIp != "" && objectIp != ipResp {
			continue
		}

		typeResp := utils.PathSearch("type", v, "")
		if objectType != "" && objectType != typeResp {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"id":         idResp,
			"ip_address": ipResp,
			"type":       typeResp,
		})
	}
	return rst
}
