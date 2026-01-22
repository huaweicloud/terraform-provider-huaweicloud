package eps

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EPS GET /v1.0/associated-resources/{resource_id}
func DataSourceAssociatedResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssociatedResourcesRead,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"associated_resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildAssociatedResourcesSchema(),
			},
			"errors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildErrorsSchema(),
			},
			// Deprecated
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: utils.SchemaDesc("No project ID configuration required.", utils.SchemaDescInput{Deprecated: true}),
			},
			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: utils.SchemaDesc("No region ID configuration required.", utils.SchemaDescInput{Deprecated: true}),
			},
		},
	}
}

func buildErrorsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildAssociatedResourcesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildQueryAssociatedResourcesProjectId(cfg *config.Config, region string, d *schema.ResourceData) string {
	if projectId := d.Get("project_id").(string); projectId != "" {
		return projectId
	}
	return cfg.GetProjectID(region)
}

func buildQueryAssociatedResourcesRegionId(region string, d *schema.ResourceData) string {
	if regionId := d.Get("region_id").(string); regionId != "" {
		return regionId
	}
	return region
}

func dataSourceAssociatedResourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v1.0/associated-resources/{resource_id}"
		projectId    = buildQueryAssociatedResourcesProjectId(cfg, region, d)
		regionId     = buildQueryAssociatedResourcesRegionId(region, d)
		resourceType = d.Get("resource_type").(string)
		resourceId   = d.Get("resource_id").(string)
	)
	client, err := cfg.NewServiceClient("eps", region)
	if err != nil {
		return diag.Errorf("error creating EPS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{resource_id}", resourceId)
	requestPath += fmt.Sprintf("?project_id=%s&region_id=%s&resource_type=%s", projectId, regionId, resourceType)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving EPS associated resources: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceId)
	mErr := multierror.Append(
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("type", utils.PathSearch("type", respBody, nil)),
		d.Set("associated_resources", flattenListAssociatedResources(respBody)),
		d.Set("errors", flattenListErrors(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListAssociatedResources(respBody interface{}) []interface{} {
	respArray := utils.PathSearch("associated_resources", respBody, make([]interface{}, 0)).([]interface{})
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"name":          utils.PathSearch("name", v, nil),
			"eip":           utils.PathSearch("eip", v, nil),
			"resource_type": utils.PathSearch("resource_type", v, nil),
		})
	}
	return rst
}

func flattenListErrors(resp interface{}) []interface{} {
	respArray := utils.PathSearch("errors", resp, make([]interface{}, 0)).([]interface{})
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"project_id":    utils.PathSearch("project_id", v, nil),
			"resource_type": utils.PathSearch("resource_type", v, nil),
			"resource_id":   utils.PathSearch("resource_id", v, nil),
			"error_code":    utils.PathSearch("error_code", v, nil),
			"error_msg":     utils.PathSearch("error_msg", v, nil),
		})
	}
	return rst
}
