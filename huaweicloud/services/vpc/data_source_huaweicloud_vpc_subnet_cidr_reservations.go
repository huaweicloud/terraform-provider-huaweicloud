package vpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPC GET /v3/{project_id}/vpc/virsubnet-cidr-reservations
func DataSourceVpcSubnetCidrReservations() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceVpcSubnetCidrReservationsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"reservation_id": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `Specifies the subnet CIDR reservation resource IDs. `,
			},
			"subnet_id": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `Specifies the IDs of the subnets containing the CIDR reservations.`,
			},
			"cidr": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `Specifies the CIDRs of the subnet reservations.`,
			},
			"ip_version": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
				Description: `Specifies the IP versions of the subnets.`,
			},
			"name": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `Specifies the names of the subnet CIDR reservations.`,
			},
			"description": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `Specifies the descriptions of the subnet CIDR reservations.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			"reservations": {
				Type:        schema.TypeList,
				Elem:        reservationsReservationSchema(),
				Computed:    true,
				Description: `Indicates the list of VPC subnet CIDR reservations.`,
			},
		},
	}
}

func reservationsReservationSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the unique identifier of the subnet CIDR reservation.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the subnet containing the CIDR reservation.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the VPC containing the subnet.`,
			},
			"ip_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the IP version of the subnet CIDR reservation.`,
			},
			"cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the CIDR of the subnet reservation.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the subnet CIDR reservation.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the description of the subnet CIDR reservation.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the project ID to which the CIDR reservation belongs.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time of the CIDR reservation.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the last update time of the CIDR reservation.`,
			},
		},
	}
	return &sc
}

func resourceVpcSubnetCidrReservationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getReservations: Query the List of VPC subnet CIDR reservations
	var (
		getReservationsHttpUrl = "v3/{project_id}/vpc/virsubnet-cidr-reservations"
		getReservationsProduct = "vpc"
	)
	getReservationsClient, err := cfg.NewServiceClient(getReservationsProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	getReservationsBasePath := getReservationsClient.Endpoint + getReservationsHttpUrl
	getReservationsBasePath = strings.ReplaceAll(getReservationsBasePath, "{project_id}", getReservationsClient.ProjectID)

	getReservationsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var reservations []interface{}
	var marker string
	var getReservationsPath string
	for {
		getReservationsPath = getReservationsBasePath + buildGetReservationsQueryParams(d, cfg, marker)
		getReservationsResp, err := getReservationsClient.Request("GET", getReservationsPath, &getReservationsOpt)
		if err != nil {
			return diag.Errorf("error retrieving VPC subnet CIDR reservations: %s", err)
		}

		getReservationsRespBody, err := utils.FlattenResponse(getReservationsResp)
		if err != nil {
			return diag.FromErr(err)
		}
		reservations = append(reservations, flattenGetReservationsResponseBodyReservation(getReservationsRespBody)...)
		marker = utils.PathSearch("page_info.next_marker", getReservationsRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("reservations", reservations),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetReservationsResponseBodyReservation(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("virsubnet_cidr_reservations", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))

	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"subnet_id":   utils.PathSearch("virsubnet_id", v, nil),
			"vpc_id":      utils.PathSearch("vpc_id", v, nil),
			"ip_version":  utils.PathSearch("ip_version", v, nil),
			"cidr":        utils.PathSearch("cidr", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"project_id":  utils.PathSearch("project_id", v, nil),
			"created_at":  utils.PathSearch("created_at", v, nil),
			"updated_at":  utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}

func buildGetReservationsQueryParams(d *schema.ResourceData, cfg *config.Config, marker string) string {
	res := "?limit=2000" // Maximum limit as per API documentation

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	// Add array parameters
	if v, ok := d.GetOk("reservation_id"); ok {
		for _, id := range v.([]interface{}) {
			res = fmt.Sprintf("%s&id=%v", res, id)
		}
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		for _, virsubnetId := range v.([]interface{}) {
			res = fmt.Sprintf("%s&virsubnet_id=%v", res, virsubnetId)
		}
	}

	if v, ok := d.GetOk("cidr"); ok {
		for _, cidr := range v.([]interface{}) {
			res = fmt.Sprintf("%s&cidr=%v", res, cidr)
		}
	}

	if v, ok := d.GetOk("ip_version"); ok {
		for _, ipVersion := range v.([]interface{}) {
			res = fmt.Sprintf("%s&ip_version=%v", res, ipVersion)
		}
	}

	if v, ok := d.GetOk("name"); ok {
		for _, name := range v.([]interface{}) {
			res = fmt.Sprintf("%s&name=%v", res, name)
		}
	}

	if v, ok := d.GetOk("description"); ok {
		for _, description := range v.([]interface{}) {
			res = fmt.Sprintf("%s&description=%v", res, description)
		}
	}

	res = fmt.Sprintf("%s&enterprise_project_id=%v", res, cfg.GetEnterpriseProjectID(d, "all_granted_eps"))

	return res
}
