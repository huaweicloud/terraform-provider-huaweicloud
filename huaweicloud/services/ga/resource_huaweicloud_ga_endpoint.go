// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product GA
// ---------------------------------------------------------------

package ga

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

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
		},
	}
}

func resourceEndpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	// createEndpoint: Create a GA Endpoint.
	var (
		createEndpointHttpUrl = "v1/endpoint-groups/{endpoint_group_id}/endpoints"
		createEndpointProduct = "ga"
	)
	createEndpointClient, err := conf.NewServiceClient(createEndpointProduct, region)
	if err != nil {
		return diag.Errorf("error creating Endpoint Client: %s", err)
	}

	createEndpointPath := createEndpointClient.Endpoint + createEndpointHttpUrl
	createEndpointPath = strings.ReplaceAll(createEndpointPath, "{endpoint_group_id}", fmt.Sprintf("%v", d.Get("endpoint_group_id")))

	createEndpointOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createEndpointOpt.JSONBody = utils.RemoveNil(buildCreateEndpointBodyParams(d))
	createEndpointResp, err := createEndpointClient.Request("POST", createEndpointPath, &createEndpointOpt)
	if err != nil {
		return diag.Errorf("error creating Endpoint: %s", err)
	}

	createEndpointRespBody, err := utils.FlattenResponse(createEndpointResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("endpoint.id", createEndpointRespBody)
	if err != nil {
		return diag.Errorf("error creating Endpoint: ID is not found in API response")
	}
	d.SetId(id.(string))

	err = createEndpointWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Create of Endpoint (%s) to complete: %s", d.Id(), err)
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
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// createEndpointWaiting: missing operation notes
			var (
				createEndpointWaitingHttpUrl = "v1/endpoint-groups/{endpoint_group_id}/endpoints/{endpoint_id}"
				createEndpointWaitingProduct = "ga"
			)
			createEndpointWaitingClient, err := config.NewServiceClient(createEndpointWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating Endpoint Client: %s", err)
			}

			createEndpointWaitingPath := createEndpointWaitingClient.Endpoint + createEndpointWaitingHttpUrl
			createEndpointWaitingPath = strings.ReplaceAll(createEndpointWaitingPath,
				"{endpoint_group_id}", fmt.Sprintf("%v", d.Get("endpoint_group_id")))
			createEndpointWaitingPath = strings.ReplaceAll(createEndpointWaitingPath, "{endpoint_id}", d.Id())

			createEndpointWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			createEndpointWaitingResp, err := createEndpointWaitingClient.Request("GET", createEndpointWaitingPath, &createEndpointWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createEndpointWaitingRespBody, err := utils.FlattenResponse(createEndpointWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`endpoint.status`, createEndpointWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `endpoint.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"ACTIVE",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createEndpointWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return createEndpointWaitingRespBody, status, nil
			}

			return createEndpointWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceEndpointRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	// getEndpoint: Query the GA Endpoint detail
	var (
		getEndpointHttpUrl = "v1/endpoint-groups/{endpoint_group_id}/endpoints/{endpoint_id}"
		getEndpointProduct = "ga"
	)
	getEndpointClient, err := conf.NewServiceClient(getEndpointProduct, region)
	if err != nil {
		return diag.Errorf("error creating Endpoint Client: %s", err)
	}

	getEndpointPath := getEndpointClient.Endpoint + getEndpointHttpUrl
	getEndpointPath = strings.ReplaceAll(getEndpointPath, "{endpoint_group_id}", fmt.Sprintf("%v", d.Get("endpoint_group_id")))
	getEndpointPath = strings.ReplaceAll(getEndpointPath, "{endpoint_id}", d.Id())

	getEndpointOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getEndpointResp, err := getEndpointClient.Request("GET", getEndpointPath, &getEndpointOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Endpoint")
	}

	getEndpointRespBody, err := utils.FlattenResponse(getEndpointResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("created_at", utils.PathSearch("endpoint.created_at", getEndpointRespBody, nil)),
		d.Set("endpoint_group_id", utils.PathSearch("endpoint.endpoint_group_id", getEndpointRespBody, nil)),
		d.Set("health_state", utils.PathSearch("endpoint.health_state", getEndpointRespBody, nil)),
		d.Set("ip_address", utils.PathSearch("endpoint.ip_address", getEndpointRespBody, nil)),
		d.Set("resource_id", utils.PathSearch("endpoint.resource_id", getEndpointRespBody, nil)),
		d.Set("resource_type", utils.PathSearch("endpoint.resource_type", getEndpointRespBody, nil)),
		d.Set("status", utils.PathSearch("endpoint.status", getEndpointRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("endpoint.updated_at", getEndpointRespBody, nil)),
		d.Set("weight", utils.PathSearch("endpoint.weight", getEndpointRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEndpointUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	updateEndpointhasChanges := []string{
		"weight",
	}

	if d.HasChanges(updateEndpointhasChanges...) {
		// updateEndpoint: Update the configuration of GA Endpoint
		var (
			updateEndpointHttpUrl = "v1/endpoint-groups/{endpoint_group_id}/endpoints/{endpoint_id}"
			updateEndpointProduct = "ga"
		)
		updateEndpointClient, err := conf.NewServiceClient(updateEndpointProduct, region)
		if err != nil {
			return diag.Errorf("error creating Endpoint Client: %s", err)
		}

		updateEndpointPath := updateEndpointClient.Endpoint + updateEndpointHttpUrl
		updateEndpointPath = strings.ReplaceAll(updateEndpointPath, "{endpoint_group_id}", fmt.Sprintf("%v", d.Get("endpoint_group_id")))
		updateEndpointPath = strings.ReplaceAll(updateEndpointPath, "{endpoint_id}", d.Id())

		updateEndpointOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateEndpointOpt.JSONBody = utils.RemoveNil(buildUpdateEndpointBodyParams(d))
		_, err = updateEndpointClient.Request("PUT", updateEndpointPath, &updateEndpointOpt)
		if err != nil {
			return diag.Errorf("error updating Endpoint: %s", err)
		}
		err = updateEndpointWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the Update of Endpoint (%s) to complete: %s", d.Id(), err)
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
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// updateEndpointWaiting: missing operation notes
			var (
				updateEndpointWaitingHttpUrl = "v1/endpoint-groups/{endpoint_group_id}/endpoints/{endpoint_id}"
				updateEndpointWaitingProduct = "ga"
			)
			updateEndpointWaitingClient, err := config.NewServiceClient(updateEndpointWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating Endpoint Client: %s", err)
			}

			updateEndpointWaitingPath := updateEndpointWaitingClient.Endpoint + updateEndpointWaitingHttpUrl
			updateEndpointWaitingPath = strings.ReplaceAll(updateEndpointWaitingPath,
				"{endpoint_group_id}", fmt.Sprintf("%v", d.Get("endpoint_group_id")))
			updateEndpointWaitingPath = strings.ReplaceAll(updateEndpointWaitingPath, "{endpoint_id}", d.Id())

			updateEndpointWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			updateEndpointWaitingResp, err := updateEndpointWaitingClient.Request("GET", updateEndpointWaitingPath, &updateEndpointWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			updateEndpointWaitingRespBody, err := utils.FlattenResponse(updateEndpointWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`endpoint.status`, updateEndpointWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `endpoint.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"ACTIVE",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return updateEndpointWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return updateEndpointWaitingRespBody, status, nil
			}

			return updateEndpointWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceEndpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	// deleteEndpoint: Delete an existing GA Endpoint
	var (
		deleteEndpointHttpUrl = "v1/endpoint-groups/{endpoint_group_id}/endpoints/{endpoint_id}"
		deleteEndpointProduct = "ga"
	)
	deleteEndpointClient, err := conf.NewServiceClient(deleteEndpointProduct, region)
	if err != nil {
		return diag.Errorf("error creating Endpoint Client: %s", err)
	}

	deleteEndpointPath := deleteEndpointClient.Endpoint + deleteEndpointHttpUrl
	deleteEndpointPath = strings.ReplaceAll(deleteEndpointPath, "{endpoint_group_id}", fmt.Sprintf("%v", d.Get("endpoint_group_id")))
	deleteEndpointPath = strings.ReplaceAll(deleteEndpointPath, "{endpoint_id}", d.Id())

	deleteEndpointOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteEndpointClient.Request("DELETE", deleteEndpointPath, &deleteEndpointOpt)
	if err != nil {
		return diag.Errorf("error deleting Endpoint: %s", err)
	}

	err = deleteEndpointWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the Delete of Endpoint (%s) to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteEndpointWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// deleteEndpointWaiting: missing operation notes
			var (
				deleteEndpointWaitingHttpUrl = "v1/endpoint-groups/{endpoint_group_id}/endpoints/{endpoint_id}"
				deleteEndpointWaitingProduct = "ga"
			)
			deleteEndpointWaitingClient, err := config.NewServiceClient(deleteEndpointWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating Endpoint Client: %s", err)
			}

			deleteEndpointWaitingPath := deleteEndpointWaitingClient.Endpoint + deleteEndpointWaitingHttpUrl
			deleteEndpointWaitingPath = strings.ReplaceAll(deleteEndpointWaitingPath,
				"{endpoint_group_id}", fmt.Sprintf("%v", d.Get("endpoint_group_id")))
			deleteEndpointWaitingPath = strings.ReplaceAll(deleteEndpointWaitingPath, "{endpoint_id}", d.Id())

			deleteEndpointWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			deleteEndpointWaitingResp, err := deleteEndpointWaitingClient.Request("GET", deleteEndpointWaitingPath, &deleteEndpointWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
					return "Resource Not Found", "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			deleteEndpointWaitingRespBody, err := utils.FlattenResponse(deleteEndpointWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`endpoint.status`, deleteEndpointWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `endpoint.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return deleteEndpointWaitingRespBody, status, nil
			}

			return deleteEndpointWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceEndpointImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <endpoint_group_id>/<id>")
	}

	d.Set("endpoint_group_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
