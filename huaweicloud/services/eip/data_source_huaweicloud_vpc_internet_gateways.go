package eip

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

const pageLimit = 10

// @API EIP GET /v3/{project_id}/geip/vpc-igws
func DataSourceVPCInternetGateways() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVPCInternetGatewaysRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"igw_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"igw_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_igws": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_ipv6": {
							Type:     schema.TypeBool,
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
				},
			},
		},
	}
}

func dataSourceVPCInternetGatewaysRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NetworkingV3Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC V3 client: %s", err)
	}

	getInternetGatewaysHttpUrl := "v3/{project_id}/geip/vpc-igws"
	getInternetGatewaysPath := client.Endpoint + getInternetGatewaysHttpUrl
	getInternetGatewaysPath = strings.ReplaceAll(getInternetGatewaysPath, "{project_id}", client.ProjectID)
	getInternetGatewaysOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getInternetGatewaysPath += fmt.Sprintf("?limit=%v", pageLimit)
	if igwName, ok := d.GetOk("igw_name"); ok {
		getInternetGatewaysPath += fmt.Sprintf("&name=%s", igwName)
	}
	if igwID, ok := d.GetOk("igw_id"); ok {
		getInternetGatewaysPath += fmt.Sprintf("&id=%s", igwID)
	}
	if vpcID, ok := d.GetOk("vpc_id"); ok {
		getInternetGatewaysPath += fmt.Sprintf("&vpc_id=%s", vpcID)
	}
	currentTotal := 0

	results := make([]map[string]interface{}, 0)
	for {
		// Although the offset is not in HuaweiCloud API docs, but it works actually.
		currentPath := getInternetGatewaysPath + fmt.Sprintf("&offset=%d", currentTotal)
		getInternetGatewaysResp, err := client.Request("GET", currentPath, &getInternetGatewaysOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		getInternetGatewaysRespBody, err := utils.FlattenResponse(getInternetGatewaysResp)
		if err != nil {
			return diag.FromErr(err)
		}

		igws := utils.PathSearch("vpc_igws", getInternetGatewaysRespBody, make([]interface{}, 0)).([]interface{})
		for _, igw := range igws {
			results = append(results, map[string]interface{}{
				"id":          utils.PathSearch("id", igw, nil),
				"name":        utils.PathSearch("name", igw, nil),
				"vpc_id":      utils.PathSearch("vpc_id", igw, nil),
				"subnet_id":   utils.PathSearch("network_id", igw, nil),
				"enable_ipv6": utils.PathSearch("enable_ipv6", igw, false),
				"created_at":  utils.PathSearch("created_at", igw, nil),
				"updated_at":  utils.PathSearch("updated_at", igw, nil),
			})
		}

		// `current_count` means the number of `igws` in this page, and the limit of page is `10`.
		currentCount := utils.PathSearch("page_info.current_count", getInternetGatewaysRespBody, float64(0))
		if currentCount.(float64) < pageLimit {
			break
		}
		currentTotal += len(igws)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("vpc_igws", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
