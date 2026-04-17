package rfs

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

// @API RFS GET /v1/stack-sets/{stack_set_name}/operations/{stack_set_operation_id}/metadata
func DataSourceStackSetOperationMetadata() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceStackSetOperationMetadataRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"stack_set_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stack_set_operation_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stack_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Field `call_identity` has no response value.
			"call_identity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"administration_agency_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"administration_agency_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"managed_agency_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_targets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceDeploymentTargetsSchema(),
			},
			"operation_preferences": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataOperationPreferencesSchema(),
			},
		},
	}
}

func dataOperationPreferencesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"region_concurrency_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_order": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"failure_tolerance_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"failure_tolerance_percentage": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_concurrent_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_concurrent_percentage": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"failure_tolerance_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDeploymentTargetsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"domain_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"domain_ids_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organizational_unit_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"domain_id_filter_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildStackSetOperationMetadataQueryParams(d *schema.ResourceData) string {
	rst := ""

	if v, ok := d.GetOk("stack_set_id"); ok {
		rst += fmt.Sprintf("&stack_set_id=%s", v.(string))
	}

	if v, ok := d.GetOk("call_identity"); ok {
		rst += fmt.Sprintf("&call_identity=%s", v.(string))
	}

	if rst != "" {
		rst = "?" + rst[1:]
	}

	return rst
}

func dataSourceStackSetOperationMetadataRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/stack-sets/{stack_set_name}/operations/{stack_set_operation_id}/metadata"
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{stack_set_name}", d.Get("stack_set_name").(string))
	requestPath = strings.ReplaceAll(requestPath, "{stack_set_operation_id}", d.Get("stack_set_operation_id").(string))
	requestPath += buildStackSetOperationMetadataQueryParams(d)

	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Client-Request-Id": uuid,
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving RFS stack set operation metadata: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uuid)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("stack_set_name", d.Get("stack_set_name")),
		d.Set("stack_set_operation_id", d.Get("stack_set_operation_id")),
		d.Set("stack_set_id", utils.PathSearch("stack_set_id", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("status_message", utils.PathSearch("status_message", respBody, nil)),
		d.Set("action", utils.PathSearch("action", respBody, nil)),
		d.Set("administration_agency_name", utils.PathSearch("administration_agency_name", respBody, nil)),
		d.Set("administration_agency_urn", utils.PathSearch("administration_agency_urn", respBody, nil)),
		d.Set("managed_agency_name", utils.PathSearch("managed_agency_name", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", respBody, nil)),
		d.Set("deployment_targets", flattenDeploymentTargets(utils.PathSearch("deployment_targets", respBody, nil))),
		d.Set("operation_preferences", flattenOperationPreferences(utils.PathSearch("operation_preferences", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDeploymentTargets(deploymentTargets interface{}) []interface{} {
	if deploymentTargets == nil {
		return nil
	}

	targets := map[string]interface{}{
		"regions":                 utils.PathSearch("regions", deploymentTargets, nil),
		"domain_ids":              utils.PathSearch("domain_ids", deploymentTargets, nil),
		"domain_ids_uri":          utils.PathSearch("domain_ids_uri", deploymentTargets, nil),
		"organizational_unit_ids": utils.PathSearch("organizational_unit_ids", deploymentTargets, nil),
		"domain_id_filter_type":   utils.PathSearch("domain_id_filter_type", deploymentTargets, nil),
	}

	return []interface{}{
		targets,
	}
}

func flattenOperationPreferences(operationPreferences interface{}) []interface{} {
	if operationPreferences == nil {
		return nil
	}

	preferences := map[string]interface{}{
		"region_concurrency_type":      utils.PathSearch("region_concurrency_type", operationPreferences, nil),
		"region_order":                 utils.PathSearch("region_order", operationPreferences, nil),
		"failure_tolerance_count":      utils.PathSearch("failure_tolerance_count", operationPreferences, nil),
		"failure_tolerance_percentage": utils.PathSearch("failure_tolerance_percentage", operationPreferences, nil),
		"max_concurrent_count":         utils.PathSearch("max_concurrent_count", operationPreferences, nil),
		"max_concurrent_percentage":    utils.PathSearch("max_concurrent_percentage", operationPreferences, nil),
		"failure_tolerance_mode":       utils.PathSearch("failure_tolerance_mode", operationPreferences, nil),
	}

	return []interface{}{
		preferences,
	}
}
