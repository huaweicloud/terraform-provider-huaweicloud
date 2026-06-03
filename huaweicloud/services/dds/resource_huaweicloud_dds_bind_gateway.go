package dds

import (
	"context"
	"errors"
	"fmt"
	"log"
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

var bindGatewayNonUpdatableParams = []string{
	"instance_id",
	"node_id",
	"nat_gateway_id",
	"public_ip_id",
	"external_service_port",
}

// @API DDS POST /v3/{project_id}/instances/{instance_id}/nodes/{node_id}/public-gateway
// @API DDS DELETE /v3/{project_id}/instances/{instance_id}/nodes/{node_id}/public-gateway
// @API DDS GET /v3/{project_id}/instances
// @API EIP GET /v1/{project_id}/publicips
func ResourceBindGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBindGatewayCreate,
		ReadContext:   resourceBindGatewayRead,
		UpdateContext: resourceBindGatewayUpdate,
		DeleteContext: resourceBindGatewayDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceBindGatewayImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(bindGatewayNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_ip_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"external_service_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildCreateBindGatewayBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"nat_gateway_id":        d.Get("nat_gateway_id"),
		"public_ip_id":          d.Get("public_ip_id"),
		"external_service_port": d.Get("external_service_port"),
	}

	return bodyParams
}

func resourceBindGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/instances/{instance_id}/nodes/{node_id}/public-gateway"
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))
	createPath = strings.ReplaceAll(createPath, "{node_id}", d.Get("node_id").(string))
	createOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
		JSONBody:         buildCreateBindGatewayBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error binding the gateway to the DDS instance node: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	nodeId := utils.PathSearch("node_id", respBody, "").(string)
	if nodeId == "" {
		return diag.Errorf("error binding the gateway to the DDS instance node: unable to find node ID from the API response")
	}

	d.SetId(nodeId)

	return resourceBindGatewayRead(ctx, d, meta)
}

func GetGatewayInfo(client *golangsdk.ServiceClient, instanceId, nodeId string) (interface{}, error) {
	getHttpUrl := "v3/{project_id}/instances?id={instance_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getOpts := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	gatewayInfo := utils.PathSearch(fmt.Sprintf("instances[].groups[].nodes[]|[?id=='%s']|[0]", nodeId), respBody, nil)
	if gatewayInfo == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return gatewayInfo, nil
}

func resourceBindGatewayRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	gatewayInfo, err := GetGatewayInfo(client, instanceId, d.Id())
	if err != nil {
		// When the instance does not exist, the response HTTP status code of the query API is 400
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", instanceNotFoundCodes...),
			"error retrieving the DDS instance node binding gateway information")
	}

	gateway := utils.PathSearch("nat_gateway_id", gatewayInfo, "").(string)
	if gateway == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving the DDS instance node binding gateway information")
	}

	ipAddress := utils.PathSearch("public_ip", gatewayInfo, "").(string)
	var publicIpInfo interface{}
	if ipAddress != "" {
		vpcClient, err := cfg.NewServiceClient("vpc", region)
		if err != nil {
			return diag.Errorf("error creating VPC client: %s", err)
		}

		publicIpInfo, err = getAssociateEipInfo(vpcClient, ipAddress)
		if err != nil {
			log.Printf("[Warn] error retrieving the EIP information")
		}
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("node_id", utils.PathSearch("id", gatewayInfo, nil)),
		d.Set("nat_gateway_id", utils.PathSearch("nat_gateway_id", gatewayInfo, nil)),
		d.Set("external_service_port", utils.PathSearch("external_service_port", gatewayInfo, nil)),
		d.Set("public_ip_id", utils.PathSearch("id", publicIpInfo, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getAssociateEipInfo(client *golangsdk.ServiceClient, ipAddress string) (interface{}, error) {
	httpUrl := "v1/{project_id}/publicips?public_ip_address={public_ip_address}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{public_ip_address}", ipAddress)
	getOpts := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	ipAddressInfo := utils.PathSearch("publicips|[0]", respBody, nil)
	if ipAddressInfo == nil {
		return nil, errors.New("error retrieving EIP information")
	}

	return ipAddressInfo, nil
}

func resourceBindGatewayUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceBindGatewayDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		httpUrl    = "v3/{project_id}/instances/{instance_id}/nodes/{node_id}/public-gateway"
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deletePath = strings.ReplaceAll(deletePath, "{node_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.200602"),
			"error unbinding gateway from DDS instance node")
	}

	return nil
}

func resourceBindGatewayImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
