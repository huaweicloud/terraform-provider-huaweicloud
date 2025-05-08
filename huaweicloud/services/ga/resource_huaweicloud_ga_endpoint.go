// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product GA
// ---------------------------------------------------------------

package ga

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GA POST /v1/endpoint-groups/{endpoint_group_id}/endpoints
// @API GA GET /v1/endpoint-groups/{endpoint_group_id}/endpoints/{endpoint_id}
// @API GA PUT /v1/endpoint-groups/{endpoint_group_id}/endpoints/{endpoint_id}
// @API GA DELETE /v1/endpoint-groups/{endpoint_group_id}/endpoints/{endpoint_id}
func ResourceEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointCreate,
		UpdateContext: resourceEndpointUpdate,
		ReadContext:   resourceEndpointRead,
		DeleteContext: resourceEndpointDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceEndpointImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"endpoint_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the endpoint group to which the endpoint belongs.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the endpoint ID, for example, EIP ID.`,
			},
			// According to the service team's feedback, this parameter will be update to required in the future.
			"ip_address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the IP address of the endpoint.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "EIP",
				ForceNew:    true,
				Description: `Specifies the endpoint type.`,
				ValidateFunc: validation.StringInSlice([]string{
					"EIP",
				}, false),
			},
			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the weight of the endpoint based on which the listener distributes traffic.`,
			},
			"health_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the health check result of the endpoint.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the provisioning status.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies when the endpoint was created.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies when the endpoint was updated.`,
			},
			"frozen_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The frozen details of cloud services or resources.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The status of a cloud service or resource.`,
						},
						"effect": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The status of the resource after being forzen.`,
						},
						"scene": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The service scenario.`,
						},
					},
				},
			},
		},
	}
}

func resourceEndpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf            = meta.(*config.Config)
		region          = conf.GetRegion(d)
		httpUrl         = "v1/endpoint-groups/{endpoint_group_id}/endpoints"
		product         = "ga"
		endpointGroupID = d.Get("endpoint_group_id").(string)
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{endpoint_group_id}", endpointGroupID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateEndpointBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating GA endpoint: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	endpointId := utils.PathSearch("endpoint.id", respBody, "").(string)
	if endpointId == "" {
		return diag.Errorf("error creating GA endpoint: ID is not found in API response")
	}
	d.SetId(endpointId)

	err = createEndpointWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the GA endpoint (%s) creation to complete: %s", d.Id(), err)
	}
	return resourceEndpointRead(ctx, d, meta)
}

func buildCreateEndpointBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"endpoint": map[string]interface{}{
			"ip_address":    utils.ValueIgnoreEmpty(d.Get("ip_address")),
			"resource_id":   utils.ValueIgnoreEmpty(d.Get("resource_id")),
			"resource_type": utils.ValueIgnoreEmpty(d.Get("resource_type")),
			"weight":        utils.ValueIgnoreEmpty(d.Get("weight")),
		},
	}
	return bodyParams
}

func createEndpointWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		httpUrl          = "v1/endpoint-groups/{endpoint_group_id}/endpoints/{endpoint_id}"
		product          = "ga"
		endpointGroupID  = d.Get("endpoint_group_id").(string)
		targetStatus     = []string{"ACTIVE"}
		unexpectedStatus = []string{"ERROR"}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return fmt.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{endpoint_group_id}", endpointGroupID)
	requestPath = strings.ReplaceAll(requestPath, "{endpoint_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Request("GET", requestPath, &requestOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			respBody, err := utils.FlattenResponse(resp)
			if err != nil {
				return nil, "ERROR", err
			}
			status := utils.PathSearch(`endpoint.status`, respBody, "").(string)

			if utils.StrSliceContains(targetStatus, status) {
				return respBody, "COMPLETED", nil
			}

			if utils.StrSliceContains(unexpectedStatus, status) {
				return respBody, status, nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	return err
}

func resourceEndpointRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf            = meta.(*config.Config)
		region          = conf.GetRegion(d)
		mErr            *multierror.Error
		httpUrl         = "v1/endpoint-groups/{endpoint_group_id}/endpoints/{endpoint_id}"
		product         = "ga"
		endpointGroupID = d.Get("endpoint_group_id").(string)
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{endpoint_group_id}", endpointGroupID)
	requestPath = strings.ReplaceAll(requestPath, "{endpoint_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GA endpoint")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("created_at", utils.PathSearch("endpoint.created_at", respBody, nil)),
		d.Set("endpoint_group_id", utils.PathSearch("endpoint.endpoint_group_id", respBody, nil)),
		d.Set("health_state", utils.PathSearch("endpoint.health_state", respBody, nil)),
		d.Set("ip_address", utils.PathSearch("endpoint.ip_address", respBody, nil)),
		d.Set("resource_id", utils.PathSearch("endpoint.resource_id", respBody, nil)),
		d.Set("resource_type", utils.PathSearch("endpoint.resource_type", respBody, nil)),
		d.Set("status", utils.PathSearch("endpoint.status", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("endpoint.updated_at", respBody, nil)),
		d.Set("weight", utils.PathSearch("endpoint.weight", respBody, nil)),
		d.Set("frozen_info", flattenEndpointFrozenInfo(utils.PathSearch("endpoint.frozen_info", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEndpointFrozenInfo(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	frozenInfo := map[string]interface{}{
		"status": utils.PathSearch("status", resp, nil),
		"effect": utils.PathSearch("effect", resp, nil),
		"scene":  utils.PathSearch("scene", resp, []string{}),
	}

	return []map[string]interface{}{frozenInfo}
}

func resourceEndpointUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	if d.HasChange("weight") {
		var (
			httpUrl         = "v1/endpoint-groups/{endpoint_group_id}/endpoints/{endpoint_id}"
			product         = "ga"
			endpointGroupID = d.Get("endpoint_group_id").(string)
		)
		client, err := conf.NewServiceClient(product, region)
		if err != nil {
			return diag.Errorf("error creating GA client: %s", err)
		}

		requestPath := client.Endpoint + httpUrl
		requestPath = strings.ReplaceAll(requestPath, "{endpoint_group_id}", endpointGroupID)
		requestPath = strings.ReplaceAll(requestPath, "{endpoint_id}", d.Id())
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateEndpointBodyParams(d)),
		}

		_, err = client.Request("PUT", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error updating GA endpoint: %s", err)
		}
		err = updateEndpointWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the GA endpoint (%s) update to complete: %s", d.Id(), err)
		}
	}
	return resourceEndpointRead(ctx, d, meta)
}

func buildUpdateEndpointBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"endpoint": map[string]interface{}{
			"weight": utils.ValueIgnoreEmpty(d.Get("weight")),
		},
	}
	return bodyParams
}

func updateEndpointWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		httpUrl          = "v1/endpoint-groups/{endpoint_group_id}/endpoints/{endpoint_id}"
		product          = "ga"
		endpointGroupID  = d.Get("endpoint_group_id").(string)
		targetStatus     = []string{"ACTIVE"}
		unexpectedStatus = []string{"ERROR"}
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return fmt.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{endpoint_group_id}", endpointGroupID)
	requestPath = strings.ReplaceAll(requestPath, "{endpoint_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Request("GET", requestPath, &requestOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			respBody, err := utils.FlattenResponse(resp)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch(`endpoint.status`, respBody, "").(string)
			if utils.StrSliceContains(targetStatus, status) {
				return respBody, "COMPLETED", nil
			}

			if utils.StrSliceContains(unexpectedStatus, status) {
				return respBody, status, nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	return err
}

func resourceEndpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf            = meta.(*config.Config)
		region          = conf.GetRegion(d)
		httpUrl         = "v1/endpoint-groups/{endpoint_group_id}/endpoints/{endpoint_id}"
		product         = "ga"
		endpointGroupID = d.Get("endpoint_group_id").(string)
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{endpoint_group_id}", endpointGroupID)
	requestPath = strings.ReplaceAll(requestPath, "{endpoint_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting GA endpoint: %s", err)
	}

	err = deleteEndpointWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the GA endpoint (%s) deletion to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteEndpointWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		httpUrl          = "v1/endpoint-groups/{endpoint_group_id}/endpoints/{endpoint_id}"
		product          = "ga"
		endpointGroupID  = d.Get("endpoint_group_id").(string)
		unexpectedStatus = []string{"ERROR"}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return fmt.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{endpoint_group_id}", endpointGroupID)
	requestPath = strings.ReplaceAll(requestPath, "{endpoint_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Request("GET", requestPath, &requestOpt)
			if err != nil {
				var errDefault404 golangsdk.ErrDefault404
				if errors.As(err, &errDefault404) {
					// When the error code is `404`, the value of respBody is nil, and a non-null value is returned to
					// avoid continuing the loop check.
					return "Resource Not Found", "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			respBody, err := utils.FlattenResponse(resp)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch(`endpoint.status`, respBody, "").(string)
			if utils.StrSliceContains(unexpectedStatus, status) {
				return respBody, status, nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	return err
}

func resourceEndpointImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <endpoint_group_id>/<id>")
	}

	d.SetId(parts[1])

	return []*schema.ResourceData{d}, d.Set("endpoint_group_id", parts[0])
}
