package aad

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var blackWhiteListNonUpdatableParams = []string{"instance_id", "type"}

// @API AAD POST /v1/{project_id}/aad/external/bwlist
// @API AAD DELETE /v1/{project_id}/aad/external/bwlist
// @API AAD GET /v2/aad/policies/ddos/blackwhite-list
func ResourceBlackWhiteList() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBlackWhiteListCreate,
		ReadContext:   resourceBlackWhiteListRead,
		UpdateContext: resourceBlackWhiteListUpdate,
		DeleteContext: resourceBlackWhiteListDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBlackWhiteListImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(blackWhiteListNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the AAD instance ID.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the rule type.`,
			},
			"ips": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `Specifies the IP list.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func createBlackWhiteList(client *golangsdk.ServiceClient, d *schema.ResourceData, ips []interface{}) error {
	if len(ips) == 0 {
		return nil
	}

	requestPath := client.Endpoint + "v1/{project_id}/aad/external/bwlist"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: map[string]interface{}{
			"instance_id": d.Get("instance_id"),
			"type":        d.Get("type"),
			"ips":         ips,
		},
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func resourceBlackWhiteListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "aad"
		instanceID = d.Get("instance_id").(string)
		ips        = d.Get("ips").(*schema.Set).List()
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	if err := createBlackWhiteList(client, d, ips); err != nil {
		return diag.Errorf("error creating AAD black white list: %s", err)
	}

	d.SetId(instanceID)

	return resourceBlackWhiteListRead(ctx, d, meta)
}

func ReadBlackWhiteList(client *golangsdk.ServiceClient, typeParam, instanceID string) ([]interface{}, error) {
	requestPath := client.Endpoint + "v2/aad/policies/ddos/blackwhite-list"
	requestPath = fmt.Sprintf("%s?type=%s&instance_id=%s", requestPath, typeParam, instanceID)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	ips := utils.PathSearch("ips[].ip", respBody, make([]interface{}, 0)).([]interface{})
	if len(ips) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return ips, nil
}

func resourceBlackWhiteListRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "aad"
		typeParam  = d.Get("type").(string)
		instanceID = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	ips, err := ReadBlackWhiteList(client, typeParam, instanceID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AAD black white list")
	}

	return diag.FromErr(d.Set("ips", ips))
}

func deleteBlackWhiteList(client *golangsdk.ServiceClient, d *schema.ResourceData, ips []interface{}) error {
	if len(ips) == 0 {
		return nil
	}

	requestPath := client.Endpoint + "v1/{project_id}/aad/external/bwlist"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: map[string]interface{}{
			"instance_id": d.Get("instance_id"),
			"type":        d.Get("type"),
			"ips":         ips,
		},
	}

	_, err := client.Request("DELETE", requestPath, &requestOpt)
	return err
}

func resourceBlackWhiteListUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	if d.HasChange("ips") {
		oldRaws, newRaws := d.GetChange("ips")
		deleteIpList := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set)).List()
		addIpList := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set)).List()

		if err := deleteBlackWhiteList(client, d, deleteIpList); err != nil {
			return diag.Errorf("error deleting AAD black white list in update operation: %s", err)
		}

		if err := createBlackWhiteList(client, d, addIpList); err != nil {
			return diag.Errorf("error creating AAD black white list in update operation: %s", err)
		}
	}

	return resourceBlackWhiteListRead(ctx, d, meta)
}

func resourceBlackWhiteListDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		ips     = d.Get("ips").(*schema.Set).List()
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	if err := deleteBlackWhiteList(client, d, ips); err != nil {
		return diag.Errorf("error deleting AAD black white list: %s", err)
	}

	return nil
}

func resourceBlackWhiteListImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <instance_id>/<type>")
	}

	d.SetId(parts[0])
	mErr := multierror.Append(
		d.Set("instance_id", parts[0]),
		d.Set("type", parts[1]),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import AAD black white list, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
