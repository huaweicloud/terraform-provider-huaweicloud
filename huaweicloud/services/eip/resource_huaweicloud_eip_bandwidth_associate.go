package eip

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var eipBandwidthAssociateNonUpdatableParams = []string{
	"publicip_id",
	"bandwidth_id",
	"bandwidth_charge_mode",
	"bandwidth_size",
	"bandwidth_name",
}

// @API EIP POST /v3/{project_id}/eip/publicips/{publicip_id}/attach-share-bandwidth
// @API EIP GET /v3/{project_id}/eip/publicips/{publicip_id}
// @API EIP POST /v3/{project_id}/eip/publicips/{publicip_id}/detach-share-bandwidth
func ResourceEipBandwidthAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEipBandwidthAssociateCreate,
		ReadContext:   resourceEipBandwidthAssociateRead,
		UpdateContext: resourceEipBandwidthAssociateUpdate,
		DeleteContext: resourceEipBandwidthAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceEipBandwidthAssociateImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(eipBandwidthAssociateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"publicip_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bandwidth_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bandwidth_charge_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Required: true}),
			},
			"bandwidth_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Required: true}),
			},
			"bandwidth_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"public_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ipv6_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publicip_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildAttachShareBandwidthBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"publicip": map[string]interface{}{
			"bandwidth_id": d.Get("bandwidth_id"),
		},
	}
	return bodyParams
}

func buildDetachShareBandwidthBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"publicip": map[string]interface{}{
			"bandwidth": map[string]interface{}{
				"name":        d.Get("bandwidth_name"),
				"size":        d.Get("bandwidth_size"),
				"charge_mode": d.Get("bandwidth_charge_mode"),
			},
		},
	}

	return utils.RemoveNil(bodyParams)
}

func resourceEipBandwidthAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v3/{project_id}/eip/publicips/{publicip_id}/attach-share-bandwidth"
		product    = "vpc"
		publicipId = d.Get("publicip_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VPC EIP client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{publicip_id}", publicipId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildAttachShareBandwidthBodyParams(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating VPC EIP bandwidth associate: %s", err)
	}

	resourceId := fmt.Sprintf("%s/%s", publicipId, d.Get("bandwidth_id").(string))
	d.SetId(resourceId)

	return resourceEipBandwidthAssociateRead(ctx, d, meta)
}

func GetEipBandwidthAssociate(client *golangsdk.ServiceClient, publicipId, expectedBandwidthId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/eip/publicips/{publicip_id}"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{publicip_id}", publicipId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	publicip := utils.PathSearch("publicip", respBody, nil)
	if publicip == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	bandwidthId := utils.PathSearch("bandwidth.id", publicip, "").(string)
	shareType := utils.PathSearch("bandwidth.share_type", publicip, "").(string)

	// The EIP must be bound to the expected bandwidth ID.
	// The bandwidth type must be "WHOLE" (indicating shared bandwidth).
	// If any condition is not met, it will be treated as a non-existent resource (404).
	if bandwidthId != expectedBandwidthId || shareType != "WHOLE" {
		return nil, golangsdk.ErrDefault404{}
	}

	return publicip, nil
}

func resourceEipBandwidthAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		product     = "vpc"
		publicipId  = d.Get("publicip_id").(string)
		bandwidthId = d.Get("bandwidth_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VPC EIP client: %s", err)
	}

	publicip, err := GetEipBandwidthAssociate(client, publicipId, bandwidthId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VPC EIP bandwidth associate")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("public_ip_address", utils.PathSearch("public_ip_address", publicip, nil)),
		d.Set("public_ipv6_address", utils.PathSearch("public_ipv6_address", publicip, nil)),
		d.Set("publicip_type", utils.PathSearch("type", publicip, nil)),
		d.Set("ip_version", utils.PathSearch("ip_version", publicip, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEipBandwidthAssociateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Update()' method because resource is a one-time action resource.
	return nil
}

func resourceEipBandwidthAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v3/{project_id}/eip/publicips/{publicip_id}/detach-share-bandwidth"
		product    = "vpc"
		publicipId = d.Get("publicip_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VPC EIP client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{publicip_id}", publicipId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildDetachShareBandwidthBodyParams(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "EIP.7902"),
			"error deleting VPC EIP bandwidth associate")
	}

	return nil
}

func resourceEipBandwidthAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID,"+
			" want '<publicip_id>/<bandwidth_id>', but got '%s'", importedId)
	}

	mErr := multierror.Append(nil,
		d.Set("publicip_id", parts[0]),
		d.Set("bandwidth_id", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
