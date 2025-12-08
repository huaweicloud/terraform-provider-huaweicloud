package nat

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

// @API NAT GET /v3/{project_id}/private-nat/transit-subnets
func DataSourceNatPrivateTransitSubnets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNatPrivateTransitSubnetsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the resource IDs for querying instances.`,
			},
			"names": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the resource names for querying instances.`,
			},
			"descriptions": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the resource descriptions for querying instances.`,
			},
			"virsubnet_project_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the resource subnet project ids for querying instances.`,
			},
			"vpc_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the resource vpc ids for querying instances.`,
			},
			"virsubnet_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the resource subnet ids for querying instances.`,
			},
			"status": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the resource status for querying instances.`,
			},
			"transit_subnets": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virsubnet_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virsubnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_count": {
							Type:     schema.TypeInt,
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
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func buildNatPrivateTransitSubnetsQueryParams(d *schema.ResourceData, marker string) string {
	params := "?limit=2000"

	id := d.Get("ids").([]interface{})
	name := d.Get("names").([]interface{})
	description := d.Get("descriptions").([]interface{})
	virsubnetProjectId := d.Get("virsubnet_project_ids").([]interface{})
	vpcId := d.Get("vpc_ids").([]interface{})
	virsubnetId := d.Get("virsubnet_ids").([]interface{})
	status := d.Get("status").([]interface{})

	if len(id) > 0 {
		for _, v := range id {
			params += fmt.Sprintf("&id=%s", v)
		}
	}
	if len(name) > 0 {
		for _, v := range name {
			params += fmt.Sprintf("&name=%s", v)
		}
	}
	if len(description) > 0 {
		for _, v := range description {
			params += fmt.Sprintf("&description=%s", v)
		}
	}
	if len(virsubnetProjectId) > 0 {
		for _, v := range virsubnetProjectId {
			params += fmt.Sprintf("&virsubnet_project_id=%s", v)
		}
	}
	if len(vpcId) > 0 {
		for _, v := range vpcId {
			params += fmt.Sprintf("&vpc_id=%s", v)
		}
	}
	if len(virsubnetId) > 0 {
		for _, v := range virsubnetId {
			params += fmt.Sprintf("&virsubnet_id=%s", v)
		}
	}
	if len(status) > 0 {
		for _, v := range status {
			params += fmt.Sprintf("&status=%s", v)
		}
	}
	if v, ok := d.GetOk("page_reverse"); ok {
		params += fmt.Sprintf("&page_reverse=%s", v)
	}
	if marker != "" {
		params += fmt.Sprintf("&marker=%s", marker)
	}
	return params
}

func flattenNatPrivateTransitSubnetsResponseBody(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":                   utils.PathSearch("id", v, nil),
			"name":                 utils.PathSearch("name", v, nil),
			"description":          utils.PathSearch("description", v, nil),
			"virsubnet_project_id": utils.PathSearch("virsubnet_project_id", v, nil),
			"project_id":           utils.PathSearch("project_id", v, nil),
			"vpc_id":               utils.PathSearch("vpc_id", v, nil),
			"virsubnet_id":         utils.PathSearch("virsubnet_id", v, nil),
			"cidr":                 utils.PathSearch("cidr", v, nil),
			"type":                 utils.PathSearch("type", v, nil),
			"status":               utils.PathSearch("status", v, nil),
			"ip_count":             utils.PathSearch("ip_count", v, nil),
			"created_at":           utils.PathSearch("created_at", v, nil),
			"updated_at":           utils.PathSearch("updated_at", v, nil),
			"tags":                 utils.FlattenTagsToMap(utils.PathSearch("tags", v, make([]interface{}, 0))),
		})
	}
	return rst
}

func dataSourceNatPrivateTransitSubnetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v3/{project_id}/private-nat/transit-subnets"
		product = "nat"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating NAT client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var allNatPrivateTransitSubnets []interface{}
	var marker string
	var getTemplatesPath string

	for {
		getTemplatesPath = requestPath + buildNatPrivateTransitSubnetsQueryParams(d, marker)
		resp, err := client.Request("GET", getTemplatesPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error querying Nat Private Transit Subnets: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		transitSubnets := utils.PathSearch("transit_subnets", respBody, []interface{}{}).([]interface{})
		if len(transitSubnets) == 0 {
			break
		}
		allNatPrivateTransitSubnets = append(allNatPrivateTransitSubnets, transitSubnets...)
		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("transit_subnets", flattenNatPrivateTransitSubnetsResponseBody(allNatPrivateTransitSubnets)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
