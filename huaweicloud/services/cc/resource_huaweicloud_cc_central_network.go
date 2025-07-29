// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CC
// ---------------------------------------------------------------

package cc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CC POST /v3/{domain_id}/gcn/central-networks
// @API CC GET /v3/{domain_id}/gcn/central-networks/{central_network_id}
// @API CC DELETE /v3/{domain_id}/gcn/central-networks/{central_network_id}
// @API CC PUT /v3/{domain_id}/gcn/central-networks/{central_network_id}
func ResourceCentralNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCentralNetworkCreate,
		UpdateContext: resourceCentralNetworkUpdate,
		ReadContext:   resourceCentralNetworkRead,
		DeleteContext: resourceCentralNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the central network.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the central network.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The enterprise project ID to which the central network belongs.`,
			},
			"tags": common.TagsSchema(),
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The state of Central network.`,
			},
		},
	}
}

func resourceCentralNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createCentralNetwork: create a central network
	var (
		createCentralNetworkHttpUrl = "v3/{domain_id}/gcn/central-networks"
		createCentralNetworkProduct = "cc"
	)
	createCentralNetworkClient, err := cfg.NewServiceClient(createCentralNetworkProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	createCentralNetworkPath := createCentralNetworkClient.Endpoint + createCentralNetworkHttpUrl
	createCentralNetworkPath = strings.ReplaceAll(createCentralNetworkPath, "{domain_id}", cfg.DomainID)

	createCentralNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createCentralNetworkOpt.JSONBody = utils.RemoveNil(buildCreateCentralNetworkBodyParams(d, cfg))
	createCentralNetworkResp, err := createCentralNetworkClient.Request("POST", createCentralNetworkPath, &createCentralNetworkOpt)
	if err != nil {
		return diag.Errorf("error creating central network: %s", err)
	}

	createCentralNetworkRespBody, err := utils.FlattenResponse(createCentralNetworkResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("central_network.id", createCentralNetworkRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating central network: ID is not found in API response")
	}
	d.SetId(id)

	err = centralNetworkWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the central network (%s) creation to complete: %s", d.Id(), err)
	}
	return resourceCentralNetworkRead(ctx, d, meta)
}

func buildCreateCentralNetworkBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"central_network": map[string]interface{}{
			"name":                  d.Get("name"),
			"description":           utils.ValueIgnoreEmpty(d.Get("description")),
			"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
			"tags":                  utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		},
	}
	return bodyParams
}

func centralNetworkWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			// createCentralNetworkWaiting: missing operation notes
			var (
				createCentralNetworkWaitingHttpUrl = "v3/{domain_id}/gcn/central-networks/{id}"
				createCentralNetworkWaitingProduct = "cc"
			)
			createCentralNetworkWaitingClient, err := cfg.NewServiceClient(createCentralNetworkWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating CC client: %s", err)
			}

			createCentralNetworkWaitingPath := createCentralNetworkWaitingClient.Endpoint + createCentralNetworkWaitingHttpUrl
			createCentralNetworkWaitingPath = strings.ReplaceAll(createCentralNetworkWaitingPath, "{domain_id}", cfg.DomainID)
			createCentralNetworkWaitingPath = strings.ReplaceAll(createCentralNetworkWaitingPath, "{id}", d.Id())

			createCentralNetworkWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
				MoreHeaders: map[string]string{"Content-Type": "application/json"},
			}

			createCentralNetworkWaitingResp, err := createCentralNetworkWaitingClient.Request("GET",
				createCentralNetworkWaitingPath, &createCentralNetworkWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createCentralNetworkWaitingRespBody, err := utils.FlattenResponse(createCentralNetworkWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`central_network.state`, createCentralNetworkWaitingRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parsing %s from response body", `central_network.state`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			if utils.StrSliceContains([]string{"FAILED", "DELETED"}, status) {
				return createCentralNetworkWaitingRespBody, "", fmt.Errorf("unexpected status '%s'", status)
			}

			if status == "AVAILABLE" {
				return createCentralNetworkWaitingRespBody, "COMPLETED", nil
			}

			return createCentralNetworkWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceCentralNetworkRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getCentralNetwork: Query the central network
	var (
		getCentralNetworkHttpUrl = "v3/{domain_id}/gcn/central-networks/{id}"
		getCentralNetworkProduct = "cc"
	)
	getCentralNetworkClient, err := cfg.NewServiceClient(getCentralNetworkProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	getCentralNetworkPath := getCentralNetworkClient.Endpoint + getCentralNetworkHttpUrl
	getCentralNetworkPath = strings.ReplaceAll(getCentralNetworkPath, "{domain_id}", cfg.DomainID)
	getCentralNetworkPath = strings.ReplaceAll(getCentralNetworkPath, "{id}", d.Id())

	getCentralNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getCentralNetworkResp, err := getCentralNetworkClient.Request("GET", getCentralNetworkPath, &getCentralNetworkOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CentralNetwork")
	}

	getCentralNetworkRespBody, err := utils.FlattenResponse(getCentralNetworkResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("central_network.name", getCentralNetworkRespBody, nil)),
		d.Set("description", utils.PathSearch("central_network.description", getCentralNetworkRespBody, nil)),
		d.Set("state", utils.PathSearch("central_network.state", getCentralNetworkRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("central_network.enterprise_project_id", getCentralNetworkRespBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("central_network.tags", getCentralNetworkRespBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCentralNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateCentralNetworkChanges := []string{
		"name",
		"description",
		"tags",
	}

	if d.HasChanges(updateCentralNetworkChanges...) {
		// updateCentralNetwork: update the central network
		var (
			updateCentralNetworkHttpUrl = "v3/{domain_id}/gcn/central-networks/{id}"
			updateCentralNetworkProduct = "cc"
		)
		updateCentralNetworkClient, err := cfg.NewServiceClient(updateCentralNetworkProduct, region)
		if err != nil {
			return diag.Errorf("error creating CC client: %s", err)
		}

		updateCentralNetworkPath := updateCentralNetworkClient.Endpoint + updateCentralNetworkHttpUrl
		updateCentralNetworkPath = strings.ReplaceAll(updateCentralNetworkPath, "{domain_id}", cfg.DomainID)
		updateCentralNetworkPath = strings.ReplaceAll(updateCentralNetworkPath, "{id}", d.Id())

		updateCentralNetworkOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/json"},
		}

		updateCentralNetworkOpt.JSONBody = utils.RemoveNil(buildUpdateCentralNetworkBodyParams(d))
		_, err = updateCentralNetworkClient.Request("PUT", updateCentralNetworkPath, &updateCentralNetworkOpt)
		if err != nil {
			return diag.Errorf("error updating central network: %s", err)
		}
		err = centralNetworkWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the central network (%s) update to complete: %s", d.Id(), err)
		}
	}
	return resourceCentralNetworkRead(ctx, d, meta)
}

func buildUpdateCentralNetworkBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"central_network": map[string]interface{}{
			"name":        utils.ValueIgnoreEmpty(d.Get("name")),
			"description": utils.ValueIgnoreEmpty(d.Get("description")),
			"tags":        utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		},
	}
	return bodyParams
}

func resourceCentralNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteCentralNetwork: delete the central network
	var (
		deleteCentralNetworkHttpUrl = "v3/{domain_id}/gcn/central-networks/{id}"
		deleteCentralNetworkProduct = "cc"
	)
	deleteCentralNetworkClient, err := cfg.NewServiceClient(deleteCentralNetworkProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	deleteCentralNetworkPath := deleteCentralNetworkClient.Endpoint + deleteCentralNetworkHttpUrl
	deleteCentralNetworkPath = strings.ReplaceAll(deleteCentralNetworkPath, "{domain_id}", cfg.DomainID)
	deleteCentralNetworkPath = strings.ReplaceAll(deleteCentralNetworkPath, "{id}", d.Id())

	deleteCentralNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err = deleteCentralNetworkClient.Request("DELETE", deleteCentralNetworkPath, &deleteCentralNetworkOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting central network")
	}

	err = deleteCentralNetworkWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the central network (%s) deletion to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteCentralNetworkWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			// deleteCentralNetworkWaiting: missing operation notes
			var (
				deleteCentralNetworkWaitingHttpUrl = "v3/{domain_id}/gcn/central-networks/{id}"
				deleteCentralNetworkWaitingProduct = "cc"
			)
			deleteCentralNetworkWaitingClient, err := cfg.NewServiceClient(deleteCentralNetworkWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating CC client: %s", err)
			}

			deleteCentralNetworkWaitingPath := deleteCentralNetworkWaitingClient.Endpoint + deleteCentralNetworkWaitingHttpUrl
			deleteCentralNetworkWaitingPath = strings.ReplaceAll(deleteCentralNetworkWaitingPath, "{domain_id}", cfg.DomainID)
			deleteCentralNetworkWaitingPath = strings.ReplaceAll(deleteCentralNetworkWaitingPath, "{id}", d.Id())

			deleteCentralNetworkWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
				MoreHeaders: map[string]string{"Content-Type": "application/json"},
			}

			deleteCentralNetworkWaitingResp, err := deleteCentralNetworkWaitingClient.Request("GET",
				deleteCentralNetworkWaitingPath, &deleteCentralNetworkWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return deleteCentralNetworkWaitingResp, "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			deleteCentralNetworkWaitingRespBody, err := utils.FlattenResponse(deleteCentralNetworkWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`central_network.state`, deleteCentralNetworkWaitingRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parsing %s from response body", `central_network.state`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			if status == "FAILED" {
				return deleteCentralNetworkWaitingRespBody, "", fmt.Errorf("unexpected status '%s'", status)
			}

			if status == "DELETED" {
				return deleteCentralNetworkWaitingRespBody, "COMPLETED", nil
			}

			return deleteCentralNetworkWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
