package dew

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW POST /v1/{project_id}/dew/cpcs/associate-apps
// @API DEW POST /v1/{project_id}/dew/cpcs/disassociate-apps
// @API DEW GET /v1/{project_id}/dew/cpcs/associations
func ResourceCpcsAppClusterAssociation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCpcsAppClusterAssociationCreate,
		ReadContext:   resourceCpcsAppClusterAssociationRead,
		UpdateContext: resourceCpcsAppClusterAssociationUpdate,
		DeleteContext: resourceCpcsAppClusterAssociationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCpcsAppClusterAssociationImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"app_id",
			"cluster_id",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the application.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the cluster.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the cluster.`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the application.`,
			},
			"vpc_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the VPC.`,
			},
			"subnet_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the subnet.`,
			},
			"cluster_server_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the cluster server.`,
			},
			"vpcep_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The address of the VPC endpoint.`,
			},
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The update time of the association, UNIX timestamp in milliseconds.`,
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The creation time of the association, UNIX timestamp in milliseconds.`,
			},
		},
	}
}

func waitingForCpcsAppClusterAssociationSuccess(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	appId := d.Get("app_id").(string)
	clusterId := d.Get("cluster_id").(string)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			associationDetail, err := QueryCpcsAppClusterAssociation(client, appId, clusterId)
			if err != nil {
				return nil, "ERROR", err
			}

			// Although the API documentation does not explicitly mention the existence of the `bind_status` field,
			// we still need to use this field to guide the binding process.
			bindStatus := utils.PathSearch("bind_status", associationDetail, "").(string)
			if bindStatus == "BIND_SUCCESS" {
				return associationDetail, "COMPLETED", nil
			}

			return associationDetail, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceCpcsAppClusterAssociationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/dew/cpcs/associate-apps"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"app_id":     d.Get("app_id").(string),
			"cluster_id": d.Get("cluster_id").(string),
		},
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating the association between DEW CPCS application and cluster: %s", err)
	}

	generateId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(generateId)

	if err := waitingForCpcsAppClusterAssociationSuccess(ctx, client, d, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for DEW CPCS application cluster association to be created: %s", err)
	}

	return resourceCpcsAppClusterAssociationRead(ctx, d, meta)
}

func QueryCpcsAppClusterAssociation(client *golangsdk.ServiceClient, appId, clusterId string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/dew/cpcs/associations"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += fmt.Sprintf("?app_id=%s&cluster_id=%s", appId, clusterId)
	requestOpt := golangsdk.RequestOpts{
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

	associationDetail := utils.PathSearch("result|[0]", respBody, nil)
	if associationDetail == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return associationDetail, nil
}

func resourceCpcsAppClusterAssociationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		product   = "kms"
		appId     = d.Get("app_id").(string)
		clusterId = d.Get("cluster_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	associationDetail, err := QueryCpcsAppClusterAssociation(client, appId, clusterId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DEW CPCS application cluster association")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("cluster_id", utils.PathSearch("cluster_id", associationDetail, nil)),
		d.Set("cluster_name", utils.PathSearch("cluster_name", associationDetail, nil)),
		d.Set("app_id", utils.PathSearch("app_id", associationDetail, nil)),
		d.Set("app_name", utils.PathSearch("app_name", associationDetail, nil)),
		d.Set("vpc_name", utils.PathSearch("vpc_name", associationDetail, nil)),
		d.Set("subnet_name", utils.PathSearch("subnet_name", associationDetail, nil)),
		d.Set("cluster_server_type", utils.PathSearch("cluster_server_type", associationDetail, nil)),
		d.Set("vpcep_address", utils.PathSearch("vpcep_address", associationDetail, nil)),
		d.Set("update_time", utils.PathSearch("update_time", associationDetail, nil)),
		d.Set("create_time", utils.PathSearch("create_time", associationDetail, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCpcsAppClusterAssociationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func waitingForCpcsAppClusterAssociationDelete(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	appId := d.Get("app_id").(string)
	clusterId := d.Get("cluster_id").(string)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			associationDetail, err := QueryCpcsAppClusterAssociation(client, appId, clusterId)
			if err != nil {
				var errDefault404 golangsdk.ErrDefault404
				if errors.As(err, &errDefault404) {
					return "deleted", "COMPLETED", nil
				}
				return nil, "ERROR", err
			}

			return associationDetail, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceCpcsAppClusterAssociationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/dew/cpcs/disassociate-apps"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"app_ids":    []string{d.Get("app_id").(string)},
			"cluster_id": d.Get("cluster_id").(string),
		},
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting DEW CPCS application cluster association: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	if err := waitingForCpcsAppClusterAssociationDelete(context.Background(), client, d, timeout); err != nil {
		return diag.Errorf("error waiting for DEW CPCS application cluster association to be deleted: %s", err)
	}

	return nil
}

func resourceCpcsAppClusterAssociationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <app_id>/<cluster_id>, but got %s", d.Id())
	}

	mErr := multierror.Append(
		d.Set("app_id", parts[0]),
		d.Set("cluster_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
