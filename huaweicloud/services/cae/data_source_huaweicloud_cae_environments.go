package cae

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CAE GET /v1/{project_id}/cae/environments
func DataSourceEnvironments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnvironmentRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"environment_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the environment to be queried.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the environment to be queried.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the status of the environment to be queried.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the enterprise project to which the environments belong.",
			},

			// attributes
			"environments": {
				Type:        schema.TypeList,
				Elem:        dataSchemaEnvironments(),
				Computed:    true,
				Description: "The list of the environments.",
			},
		},
	}
}

func dataSchemaEnvironments() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the environment.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the environment.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the environment.",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "The additional attributes of the environment.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the environment.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the environment.",
			},
		},
	}
	return &sc
}

func dataSourceEnvironmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listEnvironments: Query the list of CAE environments
	var (
		listEnvironmentsHttpUrl = "v1/{project_id}/cae/environments"
		listEnvironmentsProduct = "cae"
	)
	listEnvironmentsClient, err := cfg.NewServiceClient(listEnvironmentsProduct, region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	listEnvironmentsPath := listEnvironmentsClient.Endpoint + listEnvironmentsHttpUrl
	listEnvironmentsPath = strings.ReplaceAll(listEnvironmentsPath, "{project_id}", listEnvironmentsClient.ProjectID)

	listEnvironmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listEnvironmentsResp, err := listEnvironmentsClient.Request("GET", listEnvironmentsPath, &listEnvironmentOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CAE environments")
	}

	listEnvironmentsRespBody, err := utils.FlattenResponse(listEnvironmentsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("environments", filterListEnvironmentBody(
			flattenListEnvironmentsBody(listEnvironmentsRespBody), d, cfg.GetEnterpriseProjectID(d))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListEnvironmentsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("items", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"status":      utils.PathSearch("status", v, nil),
			"annotations": utils.PathSearch("annotations", v, nil),
			"created_at":  utils.PathSearch("created_at", v, nil),
			"updated_at":  utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}

func filterListEnvironmentBody(all []interface{}, d *schema.ResourceData, enterpriseProjectId string) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if environmentId, ok := d.GetOk("environment_id"); ok && environmentId.(string) != utils.PathSearch("id", v, "").(string) {
			continue
		}
		if name, ok := d.GetOk("name"); ok && name.(string) != utils.PathSearch("name", v, "").(string) {
			continue
		}

		if status, ok := d.GetOk("status"); ok && status.(string) != utils.PathSearch("status", v, "").(string) {
			continue
		}
		if enterpriseProjectId != "" && enterpriseProjectId != utils.PathSearch("annotations.enterprise_project_id", v, "").(string) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}
