package eip

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP POST /v3/{domain_id}/global-eips/{global_eip_id}/attach-internet-bandwidth
// @API EIP POST /v3/{domain_id}/global-eips/{global_eip_id}/detach-internet-bandwidth
// @API EIP GET /v3/{domain_id}/global-eips/{global_eip_id}
func ResourceInternetBandwidthAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInternetBandwidthAssociateCreate,
		ReadContext:   resourceInternetBandwidthAssociateRead,
		UpdateContext: resourceInternetBandwidthAssociateUpdate,
		DeleteContext: resourceInternetBandwidthAssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"global_eip_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internet_bandwidth_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internet_bandwidth_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     globalEipInternetBandwidthInfoSchema(),
			},
		},
	}
}

func globalEipInternetBandwidthInfoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func buildAttachInternetBandwidthBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"global_eip": map[string]interface{}{
			"internet_bandwidth_id": d.Get("internet_bandwidth_id").(string),
		},
	}

	return bodyParams
}

func resourceInternetBandwidthAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v3/{domain_id}/global-eips/{global_eip_id}/attach-internet-bandwidth"
		globalEipId = d.Get("global_eip_id").(string)
	)

	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", cfg.DomainID)
	requestPath = strings.ReplaceAll(requestPath, "{global_eip_id}", globalEipId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAttachInternetBandwidthBodyParams(d),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating global EIP internet bandwidth associate: %s", err)
	}

	d.SetId(globalEipId)

	return resourceInternetBandwidthAssociateRead(ctx, d, meta)
}

func GetGlobalEIPWithInternetBandwidth(client *golangsdk.ServiceClient, domainID, globalEipId string) (interface{}, error) {
	getGEIPHttpUrl := "v3/{domain_id}/global-eips/{global_eip_id}"
	getGEIPPath := client.Endpoint + getGEIPHttpUrl
	getGEIPPath = strings.ReplaceAll(getGEIPPath, "{domain_id}", domainID)
	getGEIPPath = strings.ReplaceAll(getGEIPPath, "{global_eip_id}", globalEipId)
	getGEIPOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getGEIPResp, err := client.Request("GET", getGEIPPath, &getGEIPOpt)
	if err != nil {
		return nil, err
	}

	getGEIPRespBody, err := utils.FlattenResponse(getGEIPResp)
	if err != nil {
		return nil, err
	}

	target := utils.PathSearch("global_eip.internet_bandwidth_info", getGEIPRespBody, nil)
	if target == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return utils.PathSearch("global_eip", getGEIPRespBody, nil), nil
}

func resourceInternetBandwidthAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	globalEIPRaw, err := GetGlobalEIPWithInternetBandwidth(client, cfg.DomainID, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving global EIP internet bandwidth info")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("global_eip_id", utils.PathSearch("id", globalEIPRaw, nil)),
		d.Set("internet_bandwidth_info", flattenInternetBandwidthInfo(globalEIPRaw)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInternetBandwidthInfo(geip interface{}) []interface{} {
	if geip == nil {
		return nil
	}

	internetBandwidthInfo := utils.PathSearch("internet_bandwidth_info", geip, nil)
	if internetBandwidthInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":   utils.PathSearch("id", internetBandwidthInfo, nil),
			"size": utils.PathSearch("size", internetBandwidthInfo, nil),
		},
	}
}

func resourceInternetBandwidthAssociateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInternetBandwidthAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	deleteGEIPHttpUrl := "v3/{domain_id}/global-eips/{global_eip_id}/detach-internet-bandwidth"
	deleteGEIPPath := client.Endpoint + deleteGEIPHttpUrl
	deleteGEIPPath = strings.ReplaceAll(deleteGEIPPath, "{domain_id}", cfg.DomainID)
	deleteGEIPPath = strings.ReplaceAll(deleteGEIPPath, "{global_eip_id}", d.Id())

	deleteGEIPOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("POST", deleteGEIPPath, &deleteGEIPOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "GEIP.5001"),
			"error deleting internet bandwidth associate")
	}

	return nil
}
