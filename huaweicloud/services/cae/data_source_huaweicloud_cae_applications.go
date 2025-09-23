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

// @API CAE GET /v1/{project_id}/cae/applications
func DataSourceApplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the enterprise project to which the applications belong.",
			},
			// attributes
			"applications": {
				Type:     schema.TypeList,
				Elem:     dataSchemaApplications(),
				Computed: true,
			},
		},
	}
}

func dataSchemaApplications() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
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
	return &sc
}

func dataSourceApplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listApplications: Query the list of CAE applications
	var (
		listApplicationsHttpUrl = "v1/{project_id}/cae/applications"
		listApplicationsProduct = "cae"
	)
	listApplicationsClient, err := cfg.NewServiceClient(listApplicationsProduct, region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	listApplicationsPath := listApplicationsClient.Endpoint + listApplicationsHttpUrl
	listApplicationsPath = strings.ReplaceAll(listApplicationsPath, "{project_id}", listApplicationsClient.ProjectID)
	listApplicationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(d.Get("environment_id").(string), cfg.GetEnterpriseProjectID(d)),
	}

	listApplicationsResp, err := listApplicationsClient.Request("GET", listApplicationsPath, &listApplicationOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CAE applications")
	}

	listApplicationsRespBody, err := utils.FlattenResponse(listApplicationsResp)
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
		d.Set("applications", filterListApplicationBody(
			flattenListApplicationsBody(listApplicationsRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListApplicationsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("items", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"name":       utils.PathSearch("name", v, nil),
			"created_at": utils.PathSearch("created_at", v, nil),
			"updated_at": utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}

func filterListApplicationBody(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if applicationId, ok := d.GetOk("application_id"); ok && applicationId.(string) != utils.PathSearch("id", v, "").(string) {
			continue
		}
		if name, ok := d.GetOk("name"); ok && name.(string) != utils.PathSearch("name", v, "").(string) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}
