package eip

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

// @API EIP POST /v3/{project_id}/geip/vpc-igws
// @API EIP DELETE /v3/{project_id}/geip/vpc-igws/{vpc_igw_id}
// @API EIP GET /v3/{project_id}/geip/vpc-igws/{vpc_igw_id}
// @API EIP PUT /v3/{project_id}/geip/vpc-igws/{vpc_igw_id}
// @API VPC GET /v1/{project_id}/routetables
// @API VPC GET /v1/{project_id}/routetables/{routetable_id}
func ResourceVPCInternetGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPCInternetGatewayCreate,
		ReadContext:   resourceVPCInternetGatewayRead,
		UpdateContext: resourceVPCInternetGatewayUpdate,
		DeleteContext: resourceVPCInternetGatewayDelete,

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
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"add_route": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVPCInternetGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NetworkingV3Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC V3 client: %s", err)
	}

	createIGWHttpUrl := "v3/{project_id}/geip/vpc-igws"
	createIGWPath := client.Endpoint + createIGWHttpUrl
	createIGWPath = strings.ReplaceAll(createIGWPath, "{project_id}", client.ProjectID)

	createIGWOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"vpc_igw": utils.RemoveNil(buildCreateIGWBodyParams(d)),
		},
	}

	createIGWResp, err := client.Request("POST", createIGWPath, &createIGWOpt)
	if err != nil {
		return diag.FromErr(err)
	}
	createIGWRespBody, err := utils.FlattenResponse(createIGWResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("vpc_igw.id", createIGWRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find VPC IGW ID from the API response")
	}
	d.SetId(id)

	return resourceVPCInternetGatewayRead(ctx, d, meta)
}

func buildCreateIGWBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"vpc_id":      d.Get("vpc_id"),
		"network_id":  utils.ValueIgnoreEmpty(d.Get("subnet_id")),
		"add_route":   utils.ValueIgnoreEmpty(d.Get("add_route")),
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"enable_ipv6": utils.ValueIgnoreEmpty(d.Get("enable_ipv6")),
	}
	return bodyParams
}

func resourceVPCInternetGatewayRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NetworkingV3Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC V3 client: %s", err)
	}

	vpcClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC V1 client: %s", err)
	}

	// get IGW info
	getIGWHttpUrl := "v3/{project_id}/geip/vpc-igws/{vpc_igw_id}"
	getIGWPath := client.Endpoint + getIGWHttpUrl
	getIGWPath = strings.ReplaceAll(getIGWPath, "{project_id}", client.ProjectID)
	getIGWPath = strings.ReplaceAll(getIGWPath, "{vpc_igw_id}", d.Id())

	getIGWOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getIGWResp, err := client.Request("GET", getIGWPath, &getIGWOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VPC internet gateway")
	}
	getIGWRespBody, err := utils.FlattenResponse(getIGWResp)
	if err != nil {
		return diag.FromErr(err)
	}

	igw := utils.PathSearch("vpc_igw", getIGWRespBody, nil)
	if igw == nil {
		return diag.Errorf("unable to find VPC IGW from the API response")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("vpc_id", utils.PathSearch("vpc_id", igw, nil)),
		d.Set("subnet_id", utils.PathSearch("network_id", igw, nil)),
		d.Set("name", utils.PathSearch("name", igw, nil)),
		d.Set("enable_ipv6", utils.PathSearch("enable_ipv6", igw, false)),
		d.Set("created_at", utils.PathSearch("created_at", igw, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", igw, nil)),
	)

	// GET IGW API do not return add_route, it must get from default route table of VPC
	addRoute, err := getAddRoute(vpcClient, d.Get("vpc_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	mErr = multierror.Append(mErr,
		d.Set("add_route", addRoute),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting VPC internet gateway fields: %s", err)
	}
	return nil
}

func getAddRoute(client *golangsdk.ServiceClient, vpcID string) (bool, error) {
	url := "v1/{project_id}/routetables"
	path := client.Endpoint + url
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// call List API to filter ID of the default route table
	listRouteTablesPath := path + "?vpc_id=" + vpcID
	listRouteTablesResp, err := client.Request("GET", listRouteTablesPath, &opt)
	if err != nil {
		return false, fmt.Errorf("error getting route tables: %s", err)
	}
	listRouteTablesRespBody, err := utils.FlattenResponse(listRouteTablesResp)
	if err != nil {
		return false, fmt.Errorf("error flatten route tables: %s", err)
	}

	routeTableID := utils.PathSearch("routetables[?default==`true`]|[0].id", listRouteTablesRespBody, "").(string)
	if routeTableID == "" {
		return false, fmt.Errorf("unable to find default route table ID from the API response")
	}

	// call Get API to retrieve more details about the default route table
	getRouteTablePath := path + "/" + routeTableID
	getRouteTableResp, err := client.Request("GET", getRouteTablePath, &opt)
	if err != nil {
		return false, fmt.Errorf("error getting default route table: %s", err)
	}
	getRouteTableRespBody, err := utils.FlattenResponse(getRouteTableResp)
	if err != nil {
		return false, fmt.Errorf("error flatten default route table: %s", err)
	}

	route := utils.PathSearch("routetable.routes[?type=='igw']", getRouteTableRespBody, make([]interface{}, 0))
	return len(route.([]interface{})) == 1, nil
}

func resourceVPCInternetGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NetworkingV3Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC V3 client: %s", err)
	}

	updateChanges := []string{
		"enable_ipv6",
		"name",
	}

	if d.HasChanges(updateChanges...) {
		// Precheck `enable_ipv6`, it is not allow change true to false
		if d.HasChange("enable_ipv6") && !d.Get("enable_ipv6").(bool) {
			return diag.Errorf("error updating enable_ipv6: not allow change into false")
		}

		updateIGWHttpUrl := "v3/{project_id}/geip/vpc-igws/{vpc_igw_id}"
		updateIGWPath := client.Endpoint + updateIGWHttpUrl
		updateIGWPath = strings.ReplaceAll(updateIGWPath, "{project_id}", client.ProjectID)
		updateIGWPath = strings.ReplaceAll(updateIGWPath, "{vpc_igw_id}", d.Id())

		updateIGWOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"vpc_igw": utils.RemoveNil(map[string]interface{}{
					"enable_ipv6": utils.ValueIgnoreEmpty(d.Get("enable_ipv6")),
					"name":        utils.ValueIgnoreEmpty(d.Get("name")),
				}),
			},
		}

		_, err = client.Request("PUT", updateIGWPath, &updateIGWOpt)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceVPCInternetGatewayRead(ctx, d, meta)
}

func resourceVPCInternetGatewayDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NetworkingV3Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC V3 client: %s", err)
	}

	deleteIGWHttpUrl := "v3/{project_id}/geip/vpc-igws/{vpc_igw_id}"
	deleteIGWPath := client.Endpoint + deleteIGWHttpUrl
	deleteIGWPath = strings.ReplaceAll(deleteIGWPath, "{project_id}", client.ProjectID)
	deleteIGWPath = strings.ReplaceAll(deleteIGWPath, "{vpc_igw_id}", d.Id())

	deleteIGWOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deleteIGWPath, &deleteIGWOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
