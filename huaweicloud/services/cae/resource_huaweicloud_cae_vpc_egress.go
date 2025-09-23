package cae

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

var vpcEgressResourceNotFoundCodes = []string{
	"CAE.01500404", // The environment not found, and status code is 400.
	"CAE.01500000", // The resource not found, and status code is 500.
}

// @API CAE POST /v1/{project_id}/cae/vpc-egress
// @API CAE GET /v1/{project_id}/cae/vpc-egress
// @API CAE DELETE /v1/{project_id}/cae/vpc-egress/{vpc_egress_id}
func ResourceVpcEgress() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcEgressCreate,
		ReadContext:   resourceVpcEgressRead,
		DeleteContext: resourceVpcEgressDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceVpcEgressImportState,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the CAE environment.`,
			},
			"route_table_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the route table corresponding to the subnet to which the CAE environment belongs.`,
			},
			"cidr": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The destination CIDR of the routing table corresponding to the subnet to which the CAE environment belongs.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The ID of the enterprise project to which the VPC egress belongs.`,
			},
		},
	}
}

func buildCreateVpcEgressBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"api_version": "v1",
		"kind":        "VpcEgress",
		"spec": map[string]interface{}{
			"cidrs": []map[string]interface{}{
				{
					"route_table_id": d.Get("route_table_id"),
					"cidr":           d.Get("cidr"),
				},
			},
		},
	}
}

func resourceVpcEgressCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v1/{project_id}/cae/vpc-egress"
		environmentId = d.Get("environment_id").(string)
	)
	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(environmentId, cfg.GetEnterpriseProjectID(d)),
		JSONBody:         buildCreateVpcEgressBodyParams(d),
	}
	createResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error creating CAE environment (%s) access VPC: %s", environmentId, err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	vpcEgressId := utils.PathSearch("spec.cidrs|[0].id", createRespBody, "").(string)
	if vpcEgressId == "" {
		return diag.Errorf("unable to find the egress ID of the CAE environment access VPC from the API response")
	}

	d.SetId(vpcEgressId)
	return resourceVpcEgressRead(ctx, d, meta)
}

func GetVpcEgressById(client *golangsdk.ServiceClient, environmentId, vpcEgressId, epsId string) (interface{}, error) {
	vpcEgresses, err := getVpcEgresses(client, environmentId, epsId)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", vpcEgressResourceNotFoundCodes...)
	}

	vpcEgress := utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0]", vpcEgressId), vpcEgresses, nil)
	if vpcEgress == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return vpcEgress, nil
}

func getVpcEgresses(client *golangsdk.ServiceClient, environmentId string, epsId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/cae/vpc-egress"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(environmentId, epsId),
	}
	resp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("spec.cidrs", respBody, nil), nil
}

func resourceVpcEgressRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	vpcEgress, err := GetVpcEgressById(client, d.Get("environment_id").(string), d.Id(), cfg.GetEnterpriseProjectID(d))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CAE environment access VPC configuration")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("route_table_id", utils.PathSearch("route_table_id", vpcEgress, nil)),
		d.Set("cidr", utils.PathSearch("cidr", vpcEgress, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceVpcEgressDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/cae/vpc-egress/{vpc_egress_id}"
		vpcEgressId = d.Id()
	)
	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{vpc_egress_id}", vpcEgressId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(d.Get("environment_id").(string), cfg.GetEnterpriseProjectID(d)),
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(
			common.ConvertExpected500ErrInto404Err(err, "error_code", vpcEgressResourceNotFoundCodes...),
			"error_code",
			vpcEgressResourceNotFoundCodes...), fmt.Sprintf("error deleting CAE environment access VPC (%s)", vpcEgressId))
	}

	return nil
}

func resourceVpcEgressImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var (
		cfg        = meta.(*config.Config)
		importedId = d.Id()
		// The resource ID cannot be obtained on the console. so we need to import by the <environment_id>,<route_table_id>,<cidr>.
		// The cidr contains a slash (/), so use commas (,) to separate the imported IDs.
		parts = strings.Split(importedId, ",")
	)

	if len(parts) != 3 && len(parts) != 4 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<environment_id>,<route_table_id>,<cidr>', or "+
			"'<environment_id>,<route_table_id>,<cidr>,<enterprise_project_id>' but got '%s'",
			importedId)
	}

	var (
		environmentId = parts[0]
		routeTableId  = parts[1]
		cidr          = parts[2]
	)

	mErr := multierror.Append(nil,
		d.Set("environment_id", environmentId),
		d.Set("route_table_id", routeTableId),
		d.Set("cidr", cidr),
	)

	if len(parts) == 4 {
		mErr = multierror.Append(mErr, d.Set("enterprise_project_id", parts[3]))
	}

	if mErr.ErrorOrNil() != nil {
		return nil, mErr
	}

	client, err := cfg.NewServiceClient("cae", cfg.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating CAE client: %s", err)
	}

	vpcEgresses, err := getVpcEgresses(client, environmentId, cfg.GetEnterpriseProjectID(d))
	if err != nil {
		return nil, fmt.Errorf("error retrieving the list of the CAE environment access VPC configurations: %s", err)
	}

	vpcEgressId := utils.PathSearch(fmt.Sprintf("[?route_table_id=='%s'&&cidr=='%s']|[0].id", routeTableId, cidr), vpcEgresses, "").(string)
	if vpcEgressId == "" {
		return nil, fmt.Errorf("unable to find the VPC egress ID of the route table (%s) and CIDR (%s) from API response",
			routeTableId, cidr)
	}

	d.SetId(vpcEgressId)
	return []*schema.ResourceData{d}, nil
}
