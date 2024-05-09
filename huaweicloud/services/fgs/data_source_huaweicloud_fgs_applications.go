package fgs

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

// @API FunctionGraph GET /v2/{project_id}/fgs/applications
func DataSourceFunctionGraphApplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceFunctionGraphApplicationsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"applications": {
				Type:     schema.TypeList,
				Elem:     applicationSchema(),
				Computed: true,
			},
		},
	}
}

func applicationSchema() *schema.Resource {
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func datasourceFunctionGraphApplicationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                     = meta.(*config.Config)
		region                  = cfg.GetRegion(d)
		mErr                    *multierror.Error
		httpUrl                 = "v2/{project_id}/fgs/applications"
		listApplicationsProduct = "fgs"
	)

	client, err := cfg.NewServiceClient(listApplicationsProduct, region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	listApplicationsPath := client.Endpoint + httpUrl
	listApplicationsPath = strings.ReplaceAll(listApplicationsPath, "{project_id}", client.ProjectID)
	listApplicationsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var applications []interface{}
	var marker string
	for {
		getApplicationsResp, err := client.Request("GET",
			listApplicationsPath, &listApplicationsOpt)

		if err != nil {
			return diag.Errorf("error querying the applications: %s", err)
		}

		getApplicationsRespBody, err := utils.FlattenResponse(getApplicationsResp)
		if err != nil {
			return diag.FromErr(err)
		}
		applications = append(applications, flattenListApplicationsBody(getApplicationsRespBody)...)
		marker = utils.PathSearch("page_info.next_marker", getApplicationsRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("applications", filterListApplicationsBody(applications, d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListApplicationsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("applications", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"status":      utils.PathSearch("status", v, nil),
			"updated_at":  utils.FormatTimeStampRFC3339(int64(utils.PathSearch("last_modified_time", v, float64(0)).(float64))/1000, false),
			"description": utils.PathSearch("description", v, nil),
		})
	}
	return rst
}

func filterListApplicationsBody(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if applicationId, ok := d.GetOk("application_id"); ok && applicationId.(string) != utils.PathSearch("id", v, "").(string) {
			continue
		}

		if name, ok := d.GetOk("name"); ok && name.(string) != utils.PathSearch("name", v, "").(string) {
			continue
		}

		if status, ok := d.GetOk("status"); ok && status.(string) != utils.PathSearch("status", v, "").(string) {
			continue
		}

		if description, ok := d.GetOk("description"); ok && description.(string) != utils.PathSearch("description", v, "").(string) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}
