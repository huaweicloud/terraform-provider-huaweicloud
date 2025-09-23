package swrenterprise

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

var enterprisePrivateNetworkAccessControlNonUpdatableParams = []string{
	"instance_id", "vpc_id", "subnet_id", "project_id", "description",
}

// @API SWR POST /v2/{project_id}/instances/{instance_id}/internal-endpoints
// @API SWR GET /v2/{project_id}/instances/{instance_id}/internal-endpoints/{internal_endpoints_id}
// @API SWR DELETE /v2/{project_id}/instances/{instance_id}/internal-endpoints/{internal_endpoints_id}
func ResourceSwrEnterprisePrivateNetworkAccessControl() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterprisePrivateNetworkAccessControlCreate,
		UpdateContext: resourceSwrEnterprisePrivateNetworkAccessControlUpdate,
		ReadContext:   resourceSwrEnterprisePrivateNetworkAccessControlRead,
		DeleteContext: resourceSwrEnterprisePrivateNetworkAccessControlDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSwrEnterprisePrivateNetworkAccessControlImportStateFunc,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(enterprisePrivateNetworkAccessControlNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the enterprise instance ID.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the VPC ID.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the subnet ID.`,
			},
			// VPC should be in same region with the enterprise instance, but can belong to different project with the instance,
			// use same project_id of region in default
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the project ID to which the VPC belongs.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"vpcep_endpoint_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the VPCEP endpoint ID.`,
			},
			"project_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the project name to which the VPC belongs.`,
			},
			"vpc_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the VPC name.`,
			},
			"vpc_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the VPC cidr.`,
			},
			"subnet_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the subnet name.`,
			},
			"subnet_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the subnet cidr.`,
			},
			"endpoint_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the endpoint IP.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the endpoint status.`,
			},
			"status_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status text`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time of the private network access control.`,
			},
		},
	}
}

func resourceSwrEnterprisePrivateNetworkAccessControlCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	createHttpUrl := "v2/{project_id}/instances/{instance_id}/internal-endpoints"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateSwrEnterprisePrivateNetworkAccessControlBodyParams(d, client.ProjectID)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SWR enterprise instance private network access control: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find ID from the API response")
	}

	d.SetId(id)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      getPrivateNetworkAccessControlStatusRefreshFunc(client, d, false),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for private network access control (%s) status to be completed: %s", d.Id(), err)
	}

	return resourceSwrEnterprisePrivateNetworkAccessControlRead(ctx, d, meta)
}

func buildCreateSwrEnterprisePrivateNetworkAccessControlBodyParams(d *schema.ResourceData, projectId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"vpc_id":      d.Get("vpc_id"),
		"subnet_id":   d.Get("subnet_id"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"project_id":  projectId,
	}

	if v, ok := d.GetOk("project_id"); ok {
		bodyParams["project_id"] = v
	}

	return bodyParams
}

func getPrivateNetworkAccessControlStatusRefreshFunc(client *golangsdk.ServiceClient, d *schema.ResourceData,
	isDelete bool) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		rst, err := getSwrEnterprisePrivateNetworkAccessControl(client, d)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && isDelete {
				return "Resource Not Found", "DELETED", nil
			}
			return nil, "ERROR", err
		}

		status := utils.PathSearch("status", rst, "").(string)
		if status == "Running" {
			return rst, "SUCCESS", nil
		}
		if status == "CreateError" || status == "DeleteError" {
			return nil, "ERROR", errors.New("status unnormal")
		}

		return rst, "PENDING", nil
	}
}

func resourceSwrEnterprisePrivateNetworkAccessControlRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	rst, err := getSwrEnterprisePrivateNetworkAccessControl(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SWR enterprise instance private network access control")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("vpc_id", utils.PathSearch("vpc_id", rst, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", rst, nil)),
		d.Set("project_id", utils.PathSearch("project_id", rst, nil)),
		d.Set("description", utils.PathSearch("description", rst, nil)),
		d.Set("vpcep_endpoint_id", utils.PathSearch("vpcep_endpoint_id", rst, nil)),
		d.Set("project_name", utils.PathSearch("project_name", rst, nil)),
		d.Set("vpc_name", utils.PathSearch("vpc_name", rst, nil)),
		d.Set("vpc_cidr", utils.PathSearch("vpc_cidr", rst, nil)),
		d.Set("subnet_name", utils.PathSearch("subnet_name", rst, nil)),
		d.Set("subnet_cidr", utils.PathSearch("subnet_cidr", rst, nil)),
		d.Set("endpoint_ip", utils.PathSearch("endpoint_ip", rst, nil)),
		d.Set("status", utils.PathSearch("status", rst, nil)),
		d.Set("status_text", utils.PathSearch("status_text", rst, nil)),
		d.Set("created_at", utils.PathSearch("created_at", rst, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getSwrEnterprisePrivateNetworkAccessControl(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	getHttpUrl := "v2/{project_id}/instances/{instance_id}/internal-endpoints/{internal_endpoints_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{internal_endpoints_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func resourceSwrEnterprisePrivateNetworkAccessControlUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrEnterprisePrivateNetworkAccessControlDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	deleteHttpUrl := "v2/{project_id}/instances/{instance_id}/internal-endpoints/{internal_endpoints_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{internal_endpoints_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SWR enterprise instance private network access control")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      getPrivateNetworkAccessControlStatusRefreshFunc(client, d, true),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for private network access control (%s) status to be deleted: %s", d.Id(), err)
	}

	return nil
}

func resourceSwrEnterprisePrivateNetworkAccessControlImportStateFunc(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s'", d.Id())
	}

	d.SetId(parts[1])
	if err := d.Set("instance_id", parts[0]); err != nil {
		return nil, fmt.Errorf("error saving instance ID: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}
