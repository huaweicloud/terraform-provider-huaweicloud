package vpc

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

var reservationNonUpdatableParams = []string{
	"subnet_id", "ip_version", "cidr", "mask",
}

// @API VPC POST /v3/{project_id}/vpc/virsubnet-cidr-reservations
// @API VPC GET /v3/{project_id}/vpc/virsubnet-cidr-reservations/{reservation_id}
// @API VPC PUT /v3/{project_id}/vpc/virsubnet-cidr-reservations/{reservation_id}
// @API VPC DELETE /v3/{project_id}/vpc/virsubnet-cidr-reservations/{reservation_id}

func ResourceVpcSubnetCidrReservation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubnetCidrReservationCreate,
		ReadContext:   resourceSubnetCidrReservationRead,
		UpdateContext: resourceSubnetCidrReservationUpdate,
		DeleteContext: resourceSubnetCidrReservationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(reservationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the subnet CIDR reservation is located.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the virtual subnet to which the CIDR reservation belongs.`,
			},
			"ip_version": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `IP version of the subnet CIDR reservation.`,
			},
			"cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Reserved CIDR block in CIDR notation.`,
			},
			"mask": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Subnet mask length.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the subnet CIDR reservation.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the subnet CIDR reservation.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the VPC to which the subnet belongs.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The project ID of the subnet CIDR reservation.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the subnet CIDR reservation.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The last update time of the subnet CIDR reservation.`,
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

func buildCreateSubnetCidrReservationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"virsubnet_cidr_reservation": map[string]interface{}{
			"virsubnet_id": d.Get("subnet_id").(string),
			"ip_version":   d.Get("ip_version").(int),
			"cidr":         utils.ValueIgnoreEmpty(d.Get("cidr")),
			"mask":         utils.ValueIgnoreEmpty(d.Get("mask")),
			"name":         utils.ValueIgnoreEmpty(d.Get("name")),
			"description":  utils.ValueIgnoreEmpty(d.Get("description")),
		},
	}
	return bodyParams
}

func resourceSubnetCidrReservationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	httpUrl := "vpc/virsubnet-cidr-reservations"
	createPath := client.ResourceBaseURL() + httpUrl

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}

	createOpts.JSONBody = utils.RemoveNil(buildCreateSubnetCidrReservationBodyParams(d))
	resp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error creating subnet CIDR reservation: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("virsubnet_cidr_reservation.id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating subnet CIDR reservation: ID is not found in API response")
	}

	d.SetId(id)
	return resourceSubnetCidrReservationRead(ctx, d, meta)
}

func resourceSubnetCidrReservationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	httpUrl := fmt.Sprintf("vpc/virsubnet-cidr-reservations/%s", d.Id())
	getPath := client.ResourceBaseURL() + httpUrl

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "subnet CIDR reservation")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	reservation := utils.PathSearch("virsubnet_cidr_reservation", respBody, nil)

	cidr := utils.PathSearch("cidr", reservation, "").(string)
	mask := strings.Split(cidr, "/")[1]

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("subnet_id", utils.PathSearch("virsubnet_id", reservation, nil)),
		d.Set("ip_version", utils.PathSearch("ip_version", reservation, nil)),
		d.Set("cidr", cidr),
		d.Set("mask", utils.StringToInt(&mask)),
		d.Set("name", utils.PathSearch("name", reservation, nil)),
		d.Set("description", utils.PathSearch("description", reservation, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", reservation, nil)),
		d.Set("project_id", utils.PathSearch("project_id", reservation, nil)),
		d.Set("created_at", utils.PathSearch("created_at", reservation, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", reservation, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateSubnetCidrReservationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"virsubnet_cidr_reservation": map[string]interface{}{
			"name":        d.Get("name"),
			"description": d.Get("description"),
		},
	}
	return bodyParams
}

func resourceSubnetCidrReservationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	if d.HasChanges("name", "description") {
		httpUrl := fmt.Sprintf("vpc/virsubnet-cidr-reservations/%s", d.Id())
		updatePath := client.ResourceBaseURL() + httpUrl

		updateOpts := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateOpts.JSONBody = utils.RemoveNil(buildUpdateSubnetCidrReservationBodyParams(d))
		_, err = client.Request("PUT", updatePath, &updateOpts)
		if err != nil {
			return diag.Errorf("error updating subnet CIDR reservation: %s", err)
		}
	}

	return resourceSubnetCidrReservationRead(ctx, d, meta)
}

func resourceSubnetCidrReservationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	httpUrl := fmt.Sprintf("vpc/virsubnet-cidr-reservations/%s", d.Id())
	deletePath := client.ResourceBaseURL() + httpUrl

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting subnet CIDR reservation")
	}

	return nil
}
