package drs

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

// @API DRS GET /v5/{project_id}/link-flavors
func DataSourceDrsLinkFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsLinkFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the data source. If omitted, the provider-level region will be used.",
			},
			"engine_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the engine type of the link.",
			},
			"job_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the job type of the link.",
			},
			"job_direction": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the job direction of the link.",
			},
			"net_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the network type of the link.",
			},
			"node_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the node type of the link.",
			},
			"is_multi_write": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether the disaster recovery is dual-master.",
			},
			"vm_flavor": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VM flavor information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The flavor ID.",
						},
						"cloud_service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cloud service type.",
						},
						"spec_type_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The spec type code.",
						},
						"spec_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The spec code.",
						},
					},
				},
			},
			"volume_flavor": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The volume flavor information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The flavor ID.",
						},
						"cloud_service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cloud service type.",
						},
						"spec_type_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The spec type code.",
						},
						"spec_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The spec code.",
						},
					},
				},
			},
			"flow_flavor": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The flow flavor information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The flavor ID.",
						},
						"cloud_service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cloud service type.",
						},
						"spec_type_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The spec type code.",
						},
						"spec_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The spec code.",
						},
					},
				},
			},
		},
	}
}

func dataSourceDrsLinkFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v5/{project_id}/link-flavors"
		product = "drs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildListLinkFlavorsQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS link flavors: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("vm_flavor", flattenLinkFlavorInfo(getRespBody, "vm_flavor")),
		d.Set("volume_flavor", flattenLinkFlavorInfo(getRespBody, "volume_flavor")),
		d.Set("flow_flavor", flattenLinkFlavorInfo(getRespBody, "flow_flavor")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListLinkFlavorsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("engine_type"); ok {
		res = fmt.Sprintf("%s&engine_type=%v", res, v)
	}
	if v, ok := d.GetOk("job_type"); ok {
		res = fmt.Sprintf("%s&job_type=%v", res, v)
	}
	if v, ok := d.GetOk("job_direction"); ok {
		res = fmt.Sprintf("%s&job_direction=%v", res, v)
	}
	if v, ok := d.GetOk("net_type"); ok {
		res = fmt.Sprintf("%s&net_type=%v", res, v)
	}
	if v, ok := d.GetOk("node_type"); ok {
		res = fmt.Sprintf("%s&node_type=%v", res, v)
	}
	if v, ok := d.GetOk("is_multi_write"); ok {
		res = fmt.Sprintf("%s&is_multi_write=%v", res, v.(bool))
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenLinkFlavorInfo(resp interface{}, field string) []interface{} {
	curJson := utils.PathSearch(field, resp, nil)
	if curJson == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"id":                 utils.PathSearch("id", curJson, nil),
			"cloud_service_type": utils.PathSearch("cloud_service_type", curJson, nil),
			"spec_type_code":     utils.PathSearch("spec_type_code", curJson, nil),
			"spec_code":          utils.PathSearch("spec_code", curJson, nil),
		},
	}
}
