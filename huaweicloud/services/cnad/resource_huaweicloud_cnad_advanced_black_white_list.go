// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CNAD
// ---------------------------------------------------------------

package cnad

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AAD POST /v1/cnad/policies/{policy_id}/ip-list/add
// @API AAD POST /v1/cnad/policies/{policy_id}/ip-list/delete
// @API AAD GET /v1/cnad/policies/{policy_id}
func ResourceBlackWhiteList() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBlackWhiteListCreate,
		UpdateContext: resourceBlackWhiteListUpdate,
		ReadContext:   resourceBlackWhiteListRead,
		DeleteContext: resourceBlackWhiteListDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the policy ID.`,
			},
			"black_ip_list": {
				Type:         schema.TypeSet,
				Elem:         &schema.Schema{Type: schema.TypeString},
				Optional:     true,
				Computed:     true,
				Description:  `Specifies the black IP list.`,
				AtLeastOneOf: []string{"black_ip_list", "white_ip_list"},
			},
			"white_ip_list": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Specifies the white IP list.`,
			},
		},
	}
}

func resourceBlackWhiteListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aad", region)
	if err != nil {
		return diag.Errorf("error creating CNAD Client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	blackIpList := utils.ExpandToStringList(d.Get("black_ip_list").(*schema.Set).List())
	whiteIpList := utils.ExpandToStringList(d.Get("white_ip_list").(*schema.Set).List())

	if err := addPolicyIpList(client, blackIpList, "black", policyID); err != nil {
		return diag.FromErr(err)
	}

	if err := addPolicyIpList(client, whiteIpList, "white", policyID); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(policyID)
	return resourceBlackWhiteListRead(ctx, d, meta)
}

func resourceBlackWhiteListRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var mErr *multierror.Error

	var (
		getBlackWhiteListHttpUrl = "v1/cnad/policies/{policy_id}"
		getBlackWhiteListProduct = "aad"
	)
	getBlackWhiteListClient, err := cfg.NewServiceClient(getBlackWhiteListProduct, region)
	if err != nil {
		return diag.Errorf("error creating CNAD Client: %s", err)
	}

	getBlackWhiteListPath := getBlackWhiteListClient.Endpoint + getBlackWhiteListHttpUrl
	getBlackWhiteListPath = strings.ReplaceAll(getBlackWhiteListPath, "{policy_id}", d.Id())
	getBlackWhiteListOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	getBlackWhiteListResp, err := getBlackWhiteListClient.Request("GET", getBlackWhiteListPath,
		&getBlackWhiteListOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CNAD advanced policy black and white IP list")
	}

	respBody, err := utils.FlattenResponse(getBlackWhiteListResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("black_ip_list", utils.PathSearch("pop_policy.bw_list.black_ip_list",
			respBody, nil)),
		d.Set("white_ip_list", utils.PathSearch("pop_policy.bw_list.white_ip_list",
			respBody, nil)),
		d.Set("policy_id", d.Id()),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBlackWhiteListUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aad", region)
	if err != nil {
		return diag.Errorf("error creating CNAD Client: %s", err)
	}

	if d.HasChanges("black_ip_list") {
		oldRaws, newRaws := d.GetChange("black_ip_list")
		deleteIpList := utils.ExpandToStringList(oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set)).List())
		addIpList := utils.ExpandToStringList(newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set)).List())

		if err := deletePolicyIpList(client, deleteIpList, "black", d.Id()); err != nil {
			return diag.FromErr(err)
		}

		if err := addPolicyIpList(client, addIpList, "black", d.Id()); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("white_ip_list") {
		oldRaws, newRaws := d.GetChange("white_ip_list")
		deleteIpList := utils.ExpandToStringList(oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set)).List())
		addIpList := utils.ExpandToStringList(newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set)).List())

		if err := deletePolicyIpList(client, deleteIpList, "white", d.Id()); err != nil {
			return diag.FromErr(err)
		}

		if err := addPolicyIpList(client, addIpList, "white", d.Id()); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceBlackWhiteListRead(ctx, d, meta)
}

func resourceBlackWhiteListDelete(_ context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aad", region)
	if err != nil {
		return diag.Errorf("error creating CNAD Client: %s", err)
	}

	blackIpList := utils.ExpandToStringList(d.Get("black_ip_list").(*schema.Set).List())
	whiteIpList := utils.ExpandToStringList(d.Get("white_ip_list").(*schema.Set).List())
	if err := deletePolicyIpList(client, blackIpList, "black", d.Id()); err != nil {
		return diag.FromErr(err)
	}

	if err := deletePolicyIpList(client, whiteIpList, "white", d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func addPolicyIpList(client *golangsdk.ServiceClient, ipList []string, ipType, policyID string) error {
	if len(ipList) == 0 {
		return nil
	}
	addBlackWhiteListHttpUrl := "v1/cnad/policies/{policy_id}/ip-list/add"
	addBlackWhiteListPath := client.Endpoint + addBlackWhiteListHttpUrl
	addBlackWhiteListPath = strings.ReplaceAll(addBlackWhiteListPath, "{policy_id}", policyID)
	addBlackWhiteListOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"type":    ipType,
			"ip_list": ipList,
		},
	}

	_, err := client.Request("POST", addBlackWhiteListPath, &addBlackWhiteListOpt)
	if err != nil {
		return fmt.Errorf("error add %s policy IP list, %s", ipType, err)
	}
	return nil
}

func deletePolicyIpList(client *golangsdk.ServiceClient, ipList []string, ipType, policyID string) error {
	if len(ipList) == 0 {
		return nil
	}
	deleteBlackWhiteListHttpUrl := "v1/cnad/policies/{policy_id}/ip-list/delete"
	deleteBlackWhiteListPath := client.Endpoint + deleteBlackWhiteListHttpUrl
	deleteBlackWhiteListPath = strings.ReplaceAll(deleteBlackWhiteListPath, "{policy_id}", policyID)
	deleteBlackWhiteListOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"type":    ipType,
			"ip_list": ipList,
		},
	}

	_, err := client.Request("POST", deleteBlackWhiteListPath, &deleteBlackWhiteListOpt)
	if err != nil {
		return fmt.Errorf("error delete %s policy IP list, %s", ipType, err)
	}
	return nil
}
