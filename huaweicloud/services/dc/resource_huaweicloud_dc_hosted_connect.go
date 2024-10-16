// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product dc
// ---------------------------------------------------------------

package dc

import (
	"context"
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

// ResourceHostedConnect Due to insufficient test conditions, the current resource has not been tested and verified.
// @API DC DELETE /v3/{project_id}/dcaas/hosted-connects/{id}
// @API DC PUT /v3/{project_id}/dcaas/hosted-connects/{id}
// @API DC GET /v3/{project_id}/dcaas/hosted-connects/{id}
// @API DC POST /v3/{project_id}/dcaas/hosted-connects
func ResourceHostedConnect() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostedConnectCreate,
		UpdateContext: resourceHostedConnectUpdate,
		ReadContext:   resourceHostedConnectRead,
		DeleteContext: resourceHostedConnectDelete,
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
				Optional:    true,
				Computed:    true,
				Description: `The name of the hosted connect.`,
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The bandwidth size of the hosted connect in Mbit/s.`,
			},
			"hosting_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the operations connection on which the hosted connect is created.`,
			},
			"vlan": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `The VLAN allocated to the hosted connect.`,
			},
			"resource_tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The tenant ID for whom a hosted connect is to be created.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the hosted connect.`,
			},
			"peer_location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The location of the on-premises facility at the other end of the connection.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the hosted connect.`,
			},
		},
	}
}

func resourceHostedConnectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createHostedConnectHttpUrl = "v3/{project_id}/dcaas/hosted-connects"
		createHostedConnectProduct = "dc"
	)
	createHostedConnectClient, err := cfg.NewServiceClient(createHostedConnectProduct, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	createHostedConnectPath := createHostedConnectClient.Endpoint + createHostedConnectHttpUrl
	createHostedConnectPath = strings.ReplaceAll(createHostedConnectPath, "{project_id}", createHostedConnectClient.ProjectID)

	createHostedConnectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createHostedConnectOpt.JSONBody = utils.RemoveNil(buildCreateHostedConnectBodyParams(d))
	createHostedConnectResp, err := createHostedConnectClient.Request("POST", createHostedConnectPath, &createHostedConnectOpt)
	if err != nil {
		return diag.Errorf("error creating hosted connect: %s", err)
	}

	createHostedConnectRespBody, err := utils.FlattenResponse(createHostedConnectResp)
	if err != nil {
		return diag.FromErr(err)
	}

	connectId := utils.PathSearch("hosted_connect.id", createHostedConnectRespBody, "").(string)
	if connectId == "" {
		return diag.Errorf("unable to find the hosted connect ID from the API response")
	}
	d.SetId(connectId)

	err = hostedConnectWaitingForStateCompleted(ctx, createHostedConnectClient, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the hosted connect (%s) creation to complete: %s", d.Id(), err)
	}
	return resourceHostedConnectRead(ctx, d, meta)
}

func buildCreateHostedConnectBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"hosted_connect": map[string]interface{}{
			"name":               utils.ValueIgnoreEmpty(d.Get("name")),
			"description":        utils.ValueIgnoreEmpty(d.Get("description")),
			"bandwidth":          d.Get("bandwidth"),
			"hosting_id":         d.Get("hosting_id"),
			"vlan":               d.Get("vlan"),
			"resource_tenant_id": d.Get("resource_tenant_id"),
			"peer_location":      utils.ValueIgnoreEmpty(d.Get("peer_location")),
		},
	}
	return bodyParams
}

func hostedConnectWaitingForStateCompleted(ctx context.Context, client *golangsdk.ServiceClient, id string, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			var createHostedConnectWaitingHttpUrl = "v3/{project_id}/dcaas/hosted-connects/{id}"

			createHostedConnectWaitingPath := client.Endpoint + createHostedConnectWaitingHttpUrl
			createHostedConnectWaitingPath = strings.ReplaceAll(createHostedConnectWaitingPath, "{project_id}",
				client.ProjectID)
			createHostedConnectWaitingPath = strings.ReplaceAll(createHostedConnectWaitingPath, "{id}", id)

			createHostedConnectWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
				MoreHeaders: map[string]string{"Content-Type": "application/json"},
			}

			createHostedConnectWaitingResp, err := client.Request("GET", createHostedConnectWaitingPath,
				&createHostedConnectWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createHostedConnectWaitingRespBody, err := utils.FlattenResponse(createHostedConnectWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			status := utils.PathSearch(`hosted_connect.status`, createHostedConnectWaitingRespBody, "").(string)

			targetStatus := []string{
				"BUILD",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createHostedConnectWaitingRespBody, "COMPLETED", nil
			}

			pendingStatus := []string{
				"PENDING_CREATE",
				"PENDING_UPDATE",
			}
			if utils.StrSliceContains(pendingStatus, status) {
				return createHostedConnectWaitingRespBody, "PENDING", nil
			}

			return createHostedConnectWaitingRespBody, status, nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceHostedConnectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getHostedConnectHttpUrl = "v3/{project_id}/dcaas/hosted-connects/{id}"
		getHostedConnectProduct = "dc"
	)
	getHostedConnectClient, err := cfg.NewServiceClient(getHostedConnectProduct, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	getHostedConnectPath := getHostedConnectClient.Endpoint + getHostedConnectHttpUrl
	getHostedConnectPath = strings.ReplaceAll(getHostedConnectPath, "{project_id}", getHostedConnectClient.ProjectID)
	getHostedConnectPath = strings.ReplaceAll(getHostedConnectPath, "{id}", d.Id())

	getHostedConnectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getHostedConnectResp, err := getHostedConnectClient.Request("GET", getHostedConnectPath, &getHostedConnectOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving hosted connect")
	}

	getHostedConnectRespBody, err := utils.FlattenResponse(getHostedConnectResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("hosted_connect.name", getHostedConnectRespBody, nil)),
		d.Set("description", utils.PathSearch("hosted_connect.description", getHostedConnectRespBody, nil)),
		d.Set("bandwidth", utils.PathSearch("hosted_connect.bandwidth", getHostedConnectRespBody, nil)),
		d.Set("hosting_id", utils.PathSearch("hosted_connect.hosting_id", getHostedConnectRespBody, nil)),
		d.Set("vlan", utils.PathSearch("hosted_connect.vlan", getHostedConnectRespBody, nil)),
		d.Set("resource_tenant_id", utils.PathSearch("hosted_connect.tenant_id", getHostedConnectRespBody, nil)),
		d.Set("peer_location", utils.PathSearch("hosted_connect.peer_location", getHostedConnectRespBody, nil)),
		d.Set("status", utils.PathSearch("hosted_connect.status", getHostedConnectRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceHostedConnectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateHostedConnectChanges := []string{
		"name",
		"description",
		"bandwidth",
		"peer_location",
	}

	if d.HasChanges(updateHostedConnectChanges...) {
		var (
			updateHostedConnectHttpUrl = "v3/{project_id}/dcaas/hosted-connects/{id}"
			updateHostedConnectProduct = "dc"
		)
		updateHostedConnectClient, err := cfg.NewServiceClient(updateHostedConnectProduct, region)
		if err != nil {
			return diag.Errorf("error creating DC client: %s", err)
		}

		updateHostedConnectPath := updateHostedConnectClient.Endpoint + updateHostedConnectHttpUrl
		updateHostedConnectPath = strings.ReplaceAll(updateHostedConnectPath, "{project_id}", updateHostedConnectClient.ProjectID)
		updateHostedConnectPath = strings.ReplaceAll(updateHostedConnectPath, "{id}", d.Id())

		updateHostedConnectOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/json"},
		}

		updateHostedConnectOpt.JSONBody = utils.RemoveNil(buildUpdateHostedConnectBodyParams(d))
		_, err = updateHostedConnectClient.Request("PUT", updateHostedConnectPath, &updateHostedConnectOpt)
		if err != nil {
			return diag.Errorf("error updating hosted connect: %s", err)
		}
		err = hostedConnectWaitingForStateCompleted(ctx, updateHostedConnectClient, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the hosted connect (%s) update to complete: %s", d.Id(), err)
		}
	}
	return resourceHostedConnectRead(ctx, d, meta)
}

func buildUpdateHostedConnectBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"hosted_connect": map[string]interface{}{
			"name":          utils.ValueIgnoreEmpty(d.Get("name")),
			"description":   utils.ValueIgnoreEmpty(d.Get("description")),
			"bandwidth":     utils.ValueIgnoreEmpty(d.Get("bandwidth")),
			"peer_location": utils.ValueIgnoreEmpty(d.Get("peer_location")),
		},
	}
	return bodyParams
}

func resourceHostedConnectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteHostedConnectHttpUrl = "v3/{project_id}/dcaas/hosted-connects/{id}"
		deleteHostedConnectProduct = "cc"
	)
	deleteHostedConnectClient, err := cfg.NewServiceClient(deleteHostedConnectProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	deleteHostedConnectPath := deleteHostedConnectClient.Endpoint + deleteHostedConnectHttpUrl
	deleteHostedConnectPath = strings.ReplaceAll(deleteHostedConnectPath, "{project_id}", deleteHostedConnectClient.ProjectID)
	deleteHostedConnectPath = strings.ReplaceAll(deleteHostedConnectPath, "{id}", d.Id())

	deleteHostedConnectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err = deleteHostedConnectClient.Request("DELETE", deleteHostedConnectPath, &deleteHostedConnectOpt)
	if err != nil {
		return diag.Errorf("error deleting hosted connect: %s", err)
	}

	err = deleteHostedConnectWaitingForStateCompleted(ctx, deleteHostedConnectClient, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the hosted connect (%s) deletion to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteHostedConnectWaitingForStateCompleted(ctx context.Context, client *golangsdk.ServiceClient, id string, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			var deleteHostedConnectWaitingHttpUrl = "v3/{project_id}/dcaas/hosted-connects/{id}"

			deleteHostedConnectWaitingPath := client.Endpoint + deleteHostedConnectWaitingHttpUrl
			deleteHostedConnectWaitingPath = strings.ReplaceAll(deleteHostedConnectWaitingPath, "{project_id}",
				client.ProjectID)
			deleteHostedConnectWaitingPath = strings.ReplaceAll(deleteHostedConnectWaitingPath, "{id}", id)

			deleteHostedConnectWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
				MoreHeaders: map[string]string{"Content-Type": "application/json"},
			}

			deleteHostedConnectWaitingResp, err := client.Request("GET",
				deleteHostedConnectWaitingPath, &deleteHostedConnectWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return deleteHostedConnectWaitingResp, "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			deleteHostedConnectWaitingRespBody, err := utils.FlattenResponse(deleteHostedConnectWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			status := utils.PathSearch(`hosted_connect.status`, deleteHostedConnectWaitingRespBody, "").(string)

			pendingStatus := []string{
				"PENDING_DELETE",
			}
			if utils.StrSliceContains(pendingStatus, status) {
				return deleteHostedConnectWaitingRespBody, "PENDING", nil
			}

			return deleteHostedConnectWaitingRespBody, status, nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
