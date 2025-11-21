package nat

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

// @API NAT POST /v3/{project_id}/private-nat/transit-subnets
// @API NAT GET /v3/{project_id}/private-nat/transit-subnets/{transit_subnet_id}
// @API NAT DELETE /v3/{project_id}/private-nat/transit-subnets/{transit_subnet_id}
// @API NAT PUT /v3/{project_id}/private-nat/transit-subnets/{transit_subnet_id}
// @API NAT POST  /v3/{project_id}/transit-subnets/{resource_id}/tags
func ResourcePrivateTransitSubnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateTransitSubnetCreate,
		ReadContext:   resourcePrivateTransitSubnetRead,
		UpdateContext: resourcePrivateTransitSubnetUpdate,
		DeleteContext: resourcePrivateTransitSubnetDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the transit subnet is located.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The transit subnet ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the name of the private transit subnet.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the description of the private transit subnet.",
			},
			"virsubnet_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the ID of the project to which the transit subnet belongs.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the project ID.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the ID of the VPC to which the transit subnet belongs.",
			},
			"virsubnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the transit subnet to which the transit IP belongs.",
			},
			"cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the CIDR block of the transit subnet.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the transit subnet type.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the transit subnet.",
			},
			"ip_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Specifies the number of IP addresses that has been assigned from the transit subnet.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the transit subnet for private NAT.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the transit subnet for private NAT.",
			},
			"tags": common.TagsSchema(),
		},
	}
}

func buildCreatePrivateTransitSubnetBodyParams(d *schema.ResourceData) map[string]interface{} {
	transitSubnetBodyParams := map[string]interface{}{
		"name":                 d.Get("name"),
		"virsubnet_id":         d.Get("virsubnet_id"),
		"virsubnet_project_id": d.Get("virsubnet_project_id"),
		"description":          utils.ValueIgnoreEmpty(d.Get("description")),
		"tags":                 utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
	}

	return map[string]interface{}{
		"transit_subnet": transitSubnetBodyParams,
	}
}

func resourcePrivateTransitSubnetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/private-nat/transit-subnets"
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreatePrivateTransitSubnetBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating nat private transit subnet: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	transitSubnetId := utils.PathSearch("transit_subnet.id", respBody, "").(string)
	if transitSubnetId == "" {
		return diag.Errorf("error creating nat private transit subnet: ID is not found in API response")
	}

	d.SetId(transitSubnetId)

	return resourcePrivateTransitSubnetRead(ctx, d, meta)
}

func GetTransitSubnet(client *golangsdk.ServiceClient, transitSubnetId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/private-nat/transit-subnets/{transit_subnet_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{transit_subnet_id}", transitSubnetId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func resourcePrivateTransitSubnetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	respBody, err := GetTransitSubnet(client, d.Id())
	if err != nil {
		// If the transit subnet does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving nat private transit subnet")
	}
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("transit_subnet.name", respBody, nil)),
		d.Set("description", utils.PathSearch("transit_subnet.description", respBody, nil)),
		d.Set("virsubnet_project_id", utils.PathSearch("transit_subnet.virsubnet_project_id", respBody, nil)),
		d.Set("project_id", utils.PathSearch("transit_subnet.project_id", respBody, nil)),
		d.Set("vpc_id", utils.PathSearch("transit_subnet.vpc_id", respBody, nil)),
		d.Set("virsubnet_id", utils.PathSearch("transit_subnet.virsubnet_id", respBody, nil)),
		d.Set("cidr", utils.PathSearch("transit_subnet.cidr", respBody, nil)),
		d.Set("type", utils.PathSearch("transit_subnet.type", respBody, nil)),
		d.Set("status", utils.PathSearch("transit_subnet.status", respBody, nil)),
		d.Set("ip_count", utils.PathSearch("transit_subnet.ip_count", respBody, nil)),
		d.Set("created_at", utils.PathSearch("transit_subnet.created_at", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("transit_subnet.updated_at", respBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("transit_subnet.tags", respBody, make([]interface{}, 0)))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateNatPrivateTransitSubnetBodyParams(d *schema.ResourceData) map[string]interface{} {
	subnetBodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": d.Get("description"),
	}

	return map[string]interface{}{
		"transit_subnet": subnetBodyParams,
	}
}

func resourcePrivateTransitSubnetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NatV3Client(region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	if d.HasChanges("name", "description") {
		httpUrl := "v3/{project_id}/private-nat/transit-subnets/{transit_subnet_id}"
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{transit_subnet_id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildUpdateNatPrivateTransitSubnetBodyParams(d),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating nat private transit subnet: %s", err)
		}
	}

	if d.HasChange("tags") {
		err = utils.UpdateResourceTags(client, d, "transit-subnets", d.Id())
		if err != nil {
			return diag.Errorf("error updating tags of the private transit subnet: %s", err)
		}
	}

	return resourcePrivateTransitSubnetRead(ctx, d, meta)
}

func resourcePrivateTransitSubnetDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/private-nat/transit-subnets/{transit_subnet_id}"
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{transit_subnet_id}", d.Id())
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		// If the transit subnet does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting nat private transit subnet")
	}

	return nil
}
