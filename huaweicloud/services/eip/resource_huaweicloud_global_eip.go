package eip

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

// @API EIP POST /v3/{domain_id}/global-eips
// @API EIP DELETE /v3/{domain_id}/global-eips/{id}
// @API EIP GET /v3/{domain_id}/global-eips/{id}
// @API EIP PUT /v3/{domain_id}/global-eips/{id}
// @API EIP POST /v3/global-eip/{resource_id}/tags/delete
// @API EIP POST /v3/global-eip/{resource_id}/tags/create
func ResourceGlobalEIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlobalEIPCreate,
		ReadContext:   resourceGlobalEIPRead,
		UpdateContext: resourceGlobalEIPUpdate,
		DeleteContext: resourceGlobalEIPDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"access_site": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"geip_pool_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"internet_bandwidth_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			"isp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"frozen": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"frozen_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"polluted": {
				Type:     schema.TypeBool,
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"global_connection_bandwidth_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"global_connection_bandwidth_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"associate_instance_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"associate_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"associate_instance_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGlobalEIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	createGEIPHttpUrl := "v3/{domain_id}/global-eips"
	createGEIPPath := client.Endpoint + createGEIPHttpUrl
	createGEIPPath = strings.ReplaceAll(createGEIPPath, "{domain_id}", cfg.DomainID)

	createGEIPOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"global_eip": utils.RemoveNil(buildCreateGEIPBodyParams(d, cfg.GetEnterpriseProjectID(d))),
		},
	}

	createGEIPResp, err := client.Request("POST", createGEIPPath, &createGEIPOpt)
	if err != nil {
		return diag.FromErr(err)
	}
	createGEIPRespBody, err := utils.FlattenResponse(createGEIPResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("global_eip.id", createGEIPRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find global EIP ID from the API response")
	}
	d.SetId(id)

	err = waitForGEIPComplete(ctx, d.Timeout(schema.TimeoutCreate), d.Id(), cfg.DomainID, client)
	if err != nil {
		return diag.Errorf("Error creating global EIP: %s", err)
	}

	return resourceGlobalEIPRead(ctx, d, meta)
}

func waitForGEIPComplete(ctx context.Context, timeout time.Duration, id string, domainID string,
	client *golangsdk.ServiceClient) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      geipStatusRefreshFunc(id, domainID, client),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func buildCreateGEIPBodyParams(d *schema.ResourceData, epsID string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"access_site":           d.Get("access_site"),
		"geip_pool_name":        d.Get("geip_pool_name"),
		"internet_bandwidth_id": d.Get("internet_bandwidth_id"),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(epsID),
		"tags":                  utils.ValueIgnoreEmpty(utils.ExpandResourceTags(d.Get("tags").(map[string]interface{}))),
	}
	return bodyParams
}

func geipStatusRefreshFunc(id string, domainID string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getGEIPHttpUrl := "v3/{domain_id}/global-eips/{id}"
		getGEIPPath := client.Endpoint + getGEIPHttpUrl
		getGEIPPath = strings.ReplaceAll(getGEIPPath, "{domain_id}", domainID)
		getGEIPPath = strings.ReplaceAll(getGEIPPath, "{id}", id)
		getGEIPOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getGEIPResp, err := client.Request("GET", getGEIPPath, &getGEIPOpt)
		if err != nil {
			return nil, "ERROR", err
		}
		getGEIPRespBody, err := utils.FlattenResponse(getGEIPResp)
		if err != nil {
			return nil, "ERROR", err
		}
		status := utils.PathSearch("global_eip.status", getGEIPRespBody, "").(string)
		if status == "" {
			return nil, "ERROR", fmt.Errorf("unable to find global EIP status from the API response")
		}
		if status == "inuse" || status == "idle" {
			return getGEIPRespBody, "SUCCESS", nil
		}
		return getGEIPRespBody, "PENDING", nil
	}
}

func resourceGlobalEIPRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	getGEIPHttpUrl := "v3/{domain_id}/global-eips/{id}"
	getGEIPPath := client.Endpoint + getGEIPHttpUrl
	getGEIPPath = strings.ReplaceAll(getGEIPPath, "{domain_id}", cfg.DomainID)
	getGEIPPath = strings.ReplaceAll(getGEIPPath, "{id}", d.Id())

	getGEIPOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getGEIPResp, err := client.Request("GET", getGEIPPath, &getGEIPOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving global EIP")
	}
	getGEIPRespBody, err := utils.FlattenResponse(getGEIPResp)
	if err != nil {
		return diag.FromErr(err)
	}

	geip := utils.PathSearch("global_eip", getGEIPRespBody, nil)
	if geip == nil {
		return diag.Errorf("unable to find global EIP from the API response")
	}

	mErr := multierror.Append(nil,
		d.Set("access_site", utils.PathSearch("access_site", geip, nil)),
		d.Set("geip_pool_name", utils.PathSearch("geip_pool_name", geip, nil)),
		d.Set("internet_bandwidth_id", utils.PathSearch("internet_bandwidth_info.id", geip, nil)),
		d.Set("description", utils.PathSearch("description", geip, nil)),
		d.Set("name", utils.PathSearch("name", geip, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", geip, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", geip, nil))),
		d.Set("isp", utils.PathSearch("isp", geip, nil)),
		d.Set("polluted", utils.PathSearch("polluted", geip, false)),
		d.Set("ip_version", utils.PathSearch("ip_version", geip, float64(0))),
		d.Set("ip_address", utils.PathSearch("ip_address", geip, nil)),
		d.Set("frozen", utils.PathSearch("freezen", geip, false)),
		d.Set("frozen_info", utils.PathSearch("freezen_info", geip, nil)),
		d.Set("created_at", utils.PathSearch("created_at", geip, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", geip, nil)),
		d.Set("status", utils.PathSearch("status", geip, nil)),
		d.Set("global_connection_bandwidth_id", utils.PathSearch("global_connection_bandwidth_info.gcb_id", geip, nil)),
		d.Set("global_connection_bandwidth_type", utils.PathSearch("global_connection_bandwidth_info.gcb_type", geip, nil)),
		d.Set("associate_instance_region", utils.PathSearch("associate_instance_info.region", geip, nil)),
		d.Set("associate_instance_id", utils.PathSearch("associate_instance_info.instance_id", geip, nil)),
		d.Set("associate_instance_type", utils.PathSearch("associate_instance_info.instance_type", geip, nil)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting global EIP fields: %s", err)
	}
	return nil
}

func resourceGlobalEIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	updateChanges := []string{
		"description",
		"name",
	}

	if d.HasChanges(updateChanges...) {
		updateGEIPHttpUrl := "v3/{domain_id}/global-eips/{id}"
		updateGEIPPath := client.Endpoint + updateGEIPHttpUrl
		updateGEIPPath = strings.ReplaceAll(updateGEIPPath, "{domain_id}", cfg.DomainID)
		updateGEIPPath = strings.ReplaceAll(updateGEIPPath, "{id}", d.Id())

		updateGEIPOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"global_eip": utils.RemoveNil(map[string]interface{}{
					"description": d.Get("description"),
					"name":        utils.ValueIgnoreEmpty(d.Get("name")),
				}),
			},
		}

		_, err = client.Request("PUT", updateGEIPPath, &updateGEIPOpt)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		tagErr := updateTags(client, d, "global-eip", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of global EIP (%s): %s", d.Id(), tagErr)
		}
	}

	err = waitForGEIPComplete(ctx, d.Timeout(schema.TimeoutCreate), d.Id(), cfg.DomainID, client)
	if err != nil {
		return diag.Errorf("Error updating global EIP: %s", err)
	}

	return resourceGlobalEIPRead(ctx, d, meta)
}

func resourceGlobalEIPDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	deleteGEIPHttpUrl := "v3/{domain_id}/global-eips/{id}"
	deleteGEIPPath := client.Endpoint + deleteGEIPHttpUrl
	deleteGEIPPath = strings.ReplaceAll(deleteGEIPPath, "{domain_id}", cfg.DomainID)
	deleteGEIPPath = strings.ReplaceAll(deleteGEIPPath, "{id}", d.Id())

	deleteGEIPOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deleteGEIPPath, &deleteGEIPOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
