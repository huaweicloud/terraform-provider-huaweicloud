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
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AAD POST /v1/cnad/packages/{package_id}/protected-ips
// @API AAD GET /v1/cnad/protected-ips
func ResourceProtectedObject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProtectedObjectCreateOrUpdate,
		UpdateContext: resourceProtectedObjectCreateOrUpdate,
		ReadContext:   resourceProtectedObjectRead,
		DeleteContext: resourceProtectedObjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the instance ID in which to bind protected objects.`,
			},
			"protected_objects": {
				Type:        schema.TypeList,
				Elem:        protectedObjectSchema(),
				Required:    true,
				Description: `Specifies the advanced protected objects.`,
			},
		},
	}
}

func protectedObjectSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the protected object ID.`,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the protected object IP.`,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Specifies the protected object type. Valid values are **VPN**, **NAT**, **VIP**, **CCI**,
**EIP**, **ELB**, **REROUTING_IP**, **SubEni** and **NetInterFace**.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the protected object.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance ID which the protected object belongs to.`,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance name which the protected object belongs to.`,
			},
			"instance_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance version which the protected object belongs to.`,
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The region which the protected object belongs to.`,
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The status of the protected object.`,
			},
			"block_threshold": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The blocking threshold of the protected object.`,
			},
			"clean_threshold": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The cleaning threshold of the protected object.`,
			},
			"policy_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The policy name which the protected object binding.`,
			},
		},
	}
	return &sc
}

func resourceProtectedObjectCreateOrUpdate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	var (
		cfg                    = meta.(*config.Config)
		region                 = cfg.GetRegion(d)
		protectedObjectHttpUrl = "v1/cnad/packages/{package_id}/protected-ips"
		protectedObjectProduct = "aad"
		instanceID             = d.Get("instance_id").(string)
	)
	protectedObjectClient, err := cfg.NewServiceClient(protectedObjectProduct, region)
	if err != nil {
		return diag.Errorf("error creating CNAD Client: %s", err)
	}

	protectedObjectPath := protectedObjectClient.Endpoint + protectedObjectHttpUrl
	protectedObjectPath = strings.ReplaceAll(protectedObjectPath, "{package_id}", instanceID)

	protectedObjectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: buildProtectedObjectBodyParams(d),
	}

	_, err = protectedObjectClient.Request("POST", protectedObjectPath, &protectedObjectOpt)
	if err != nil {
		return diag.Errorf("error binding CNAD advanced protected object: %s", err)
	}

	d.SetId(instanceID)
	return resourceProtectedObjectRead(ctx, d, meta)
}

func buildProtectedObjectBodyParams(d *schema.ResourceData) map[string]interface{} {
	var protectedObjects []map[string]interface{}
	if rawArray, ok := d.Get("protected_objects").([]interface{}); ok {
		protectedObjects = make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			protectedObjects[i] = map[string]interface{}{
				"id":   utils.ValueIgnoreEmpty(raw["id"]),
				"ip":   utils.ValueIgnoreEmpty(raw["ip_address"]),
				"type": utils.ValueIgnoreEmpty(raw["type"]),
			}
		}
	}

	return map[string]interface{}{
		"protected_ip_list": protectedObjects,
	}
}

func resourceProtectedObjectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                       = meta.(*config.Config)
		region                    = cfg.GetRegion(d)
		mErr                      *multierror.Error
		getProtectedObjectHttpUrl = "v1/cnad/protected-ips"
		getProtectedObjectProduct = "aad"
	)
	getProtectedObjectClient, err := cfg.NewServiceClient(getProtectedObjectProduct, region)
	if err != nil {
		return diag.Errorf("error creating CNAD Client: %s", err)
	}

	getProtectedObjectPath := getProtectedObjectClient.Endpoint + getProtectedObjectHttpUrl
	getProtectedObjectPath += buildGetProtectedObjectQueryParams(d)

	getProtectedObjectsResp, err := pagination.ListAllItems(
		getProtectedObjectClient,
		"offset",
		getProtectedObjectPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		// here is no special error code
		return diag.Errorf("error retrieving protected objects, %s", err)
	}

	getProtectedObjectsRespJson, err := json.Marshal(getProtectedObjectsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getProtectedObjectsRespBody interface{}
	if err := json.Unmarshal(getProtectedObjectsRespJson, &getProtectedObjectsRespBody); err != nil {
		return diag.FromErr(err)
	}

	protectedObjects := flattenGetProtectedObjects(getProtectedObjectsRespBody)
	if len(protectedObjects) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving advanced protected objects")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("instance_id", d.Id()),
		d.Set("protected_objects", protectedObjects),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetProtectedObjects(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("items", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, len(curArray))
	for i, v := range curArray {
		rst[i] = map[string]interface{}{
			"id":               utils.PathSearch("id", v, nil),
			"ip_address":       utils.PathSearch("ip", v, nil),
			"type":             utils.PathSearch("type", v, nil),
			"name":             utils.PathSearch("name", v, nil),
			"instance_id":      utils.PathSearch("package_id", v, nil),
			"instance_name":    utils.PathSearch("package_name", v, nil),
			"instance_version": utils.PathSearch("package_version", v, nil),
			"region":           utils.PathSearch("region", v, nil),
			"status":           utils.PathSearch("status", v, nil),
			"block_threshold":  utils.PathSearch("block_threshold", v, nil),
			"clean_threshold":  utils.PathSearch("clean_threshold", v, nil),
			"policy_name":      utils.PathSearch("policy_name", v, nil),
		}
	}
	return rst
}

func buildGetProtectedObjectQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?package_id=%v", d.Id())
}

func resourceProtectedObjectDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                    = meta.(*config.Config)
		region                 = cfg.GetRegion(d)
		protectedObjectHttpUrl = "v1/cnad/packages/{package_id}/protected-ips"
		protectedObjectProduct = "aad"
	)
	protectedObjectClient, err := cfg.NewServiceClient(protectedObjectProduct, region)
	if err != nil {
		return diag.Errorf("error creating CNAD Client: %s", err)
	}

	protectedObjectPath := protectedObjectClient.Endpoint + protectedObjectHttpUrl
	protectedObjectPath = strings.ReplaceAll(protectedObjectPath, "{package_id}", d.Id())

	protectedObjectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: buildDeleteProtectedObjectBodyParams(),
	}

	_, err = protectedObjectClient.Request("POST", protectedObjectPath, &protectedObjectOpt)
	if err != nil {
		return diag.Errorf("error unbinding CNAD advanced protected object: %s", err)
	}
	return nil
}

func buildDeleteProtectedObjectBodyParams() map[string]interface{} {
	return map[string]interface{}{
		"protected_ip_list": []map[string]interface{}{},
	}
}
