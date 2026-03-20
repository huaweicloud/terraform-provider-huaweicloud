package mrs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API MRS GET /v2/{project_id}/metadata/version/{version_name}/available-flavor
func DataSourceAvailableFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAvailableFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the flavors are located.`,
			},
			"version_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The version name of the cluster.`,
			},
			"available_flavors": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of available flavors.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"az_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The availability zone code.`,
						},
						"az_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The availability zone name.`,
						},
						"master": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of flavors supported by master nodes.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"flavor_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The flavor name.`,
									},
								},
							},
						},
						"core": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of flavors supported by core nodes.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"flavor_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The flavor name.`,
									},
								},
							},
						},
						"task": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of flavors supported by task nodes.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"flavor_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The flavor name.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func listAvailableFlavors(client *golangsdk.ServiceClient, versionName string) (interface{}, error) {
	httpUrl := "v2/{project_id}/metadata/version/{version_name}/available-flavor"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{version_name}", versionName)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func dataSourceAvailableFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("mrs", region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	respBody, err := listAvailableFlavors(client, d.Get("version_name").(string))
	if err != nil {
		return diag.Errorf("error retrieving cluster available flavors: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("version_name", utils.PathSearch("version_name", respBody, nil)),
		d.Set("available_flavors", flattenAvailableFlavors(utils.PathSearch("available_flavors", respBody,
			make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAvailableFlavors(flavors []interface{}) []interface{} {
	if len(flavors) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(flavors))
	for _, v := range flavors {
		rst = append(rst, map[string]interface{}{
			"az_code": utils.PathSearch("az_code", v, nil),
			"az_name": utils.PathSearch("az_name", v, nil),
			"master":  flattenNodeAvailableFlavors(utils.PathSearch("master", v, make([]interface{}, 0)).([]interface{})),
			"core":    flattenNodeAvailableFlavors(utils.PathSearch("core", v, make([]interface{}, 0)).([]interface{})),
			"task":    flattenNodeAvailableFlavors(utils.PathSearch("task", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenNodeAvailableFlavors(flavors []interface{}) []interface{} {
	if len(flavors) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(flavors))
	for _, v := range flavors {
		rst = append(rst, map[string]interface{}{
			"flavor_name": utils.PathSearch("flavor_name", v, nil),
		})
	}

	return rst
}
