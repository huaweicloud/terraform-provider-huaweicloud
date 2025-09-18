package swrenterprise

import (
	"context"
	"fmt"
	"log"
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

var enterpriseInstanceNonUpdatableParams = []string{
	"name", "spec", "vpc_id", "subnet_id", "enterprise_project_id",
	"obs_encrypt", "encrypt_type", "obs_bucket_name", "description",
}

// @API SWR POST /v2/{project_id}/instances
// @API SWR DELETE /v2/{project_id}/instances/{instance_id}
// @API SWR GET /v2/{project_id}/instances/{instance_id}
// @API SWR GET /v2/{project_id}/jobs/{job_id}
// @API SWR GET /v2/{project_id}/instances/{instance_id}/configurations
// @API SWR PUT /v2/{project_id}/instances/{instance_id}/configurations
// @API SWR GET /v2/{project_id}/instances/{instance_id}/statistics
// @API SWR POST /v2/{project_id}/instances/{instance_id}/endpoint-policy
// @API SWR PUT /v2/{project_id}/instances/{instance_id}/endpoint-policy
// @API SWR GET /v2/{project_id}/instances/{instance_id}/endpoint-policy
// @API SWR POST /v2/{project_id}/{resource_type}/{resource_id}/tags/create
// @API SWR DELETE /v2/{project_id}/{resource_type}/{resource_id}/tags/delete
// @API SWR GET /v2/{project_id}/{resource_type}/{resource_id}/tags
func ResourceSwrEnterpriseInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterpriseInstanceCreate,
		UpdateContext: resourceSwrEnterpriseInstanceUpdate,
		ReadContext:   resourceSwrEnterpriseInstanceRead,
		DeleteContext: resourceSwrEnterpriseInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(enterpriseInstanceNonUpdatableParams),

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
				Description: `Specifies the name of the instance.`,
			},
			"spec": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the specification of the instance.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the VPC ID .`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the subnet ID.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			"obs_encrypt": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the OBS bucket is encrypted.`,
			},
			"encrypt_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the encrypt type.`,
			},
			"obs_bucket_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the OBS bucket name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description.`,
			},
			"anonymous_access": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to enable anonymous access.`,
			},
			"public_network_access_control_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Enable", "Disable"}, false),
				Description:  `Specifies the public network access control status.`,
			},
			"public_network_access_white_ip_list": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Specifies the public network access white IP list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the IP address or CIDR block.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the description.`,
						},
					},
				},
			},
			"tags": common.TagsSchema(`Specifies the key/value pairs to associate with the instance.`),
			"delete_obs": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to delete OBS bucket when deleting instance.`,
			},
			"delete_dns": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to delete DNS resources when deleting instance.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance version.`,
			},
			"charge_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the charge mode of instance.`,
			},
			"access_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the access address of instance.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the last update time.`,
			},
			"expires_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the expired time.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance status.`,
			},
			"user_def_obs": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user specifies the OBS bucket.`,
			},
			"vpc_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the VPC name.`,
			},
			"vpc_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the range of available subnets for the VPC.`,
			},
			"subnet_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the subnet name.`,
			},
			"subnet_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the range of available subnets for the subnet.`,
			},
		},
	}
}

func resourceSwrEnterpriseInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	createHttpUrl := "v2/{project_id}/instances"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildSwrEnterpriseInstanceBodyParams(d, client.ProjectID)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SWR instance: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("instance_id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find SWR instance ID from the API response")
	}

	d.SetId(id)

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID from the API response")
	}

	err = waitForJobComplete(ctx, client, d.Timeout(schema.TimeoutCreate), jobId)
	if err != nil {
		return diag.FromErr(err)
	}

	// `anonymous_access` default to false
	if d.Get("anonymous_access").(bool) {
		if err := updateSwrEnterpriseInstanceAnonymousAccess(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.Get("public_network_access_control_status").(string) == "Enable" {
		if err := updateSwrEnterpriseInstancePublicNetworkAccessControl(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("public_network_access_white_ip_list"); ok {
		if err := updateSwrEnterpriseInstancePublicNetworkAccessWhiteIpList(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSwrEnterpriseInstanceRead(ctx, d, meta)
}

func buildSwrEnterpriseInstanceBodyParams(d *schema.ResourceData, projectId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		// temporary support postPaid only
		"charge_mode":           "postPaid",
		"name":                  d.Get("name"),
		"spec":                  d.Get("spec"),
		"vpc_id":                d.Get("vpc_id"),
		"subnet_id":             d.Get("subnet_id"),
		"enterprise_project_id": d.Get("enterprise_project_id"),
		"project_id":            projectId,
		"obs_encrypt":           utils.ValueIgnoreEmpty(d.Get("obs_encrypt")),
		"encrypt_type":          utils.ValueIgnoreEmpty(d.Get("encrypt_type")),
		"obs_bucket_name":       utils.ValueIgnoreEmpty(d.Get("obs_bucket_name")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"resource_tags":         utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
	}

	return bodyParams
}

func waitForJobComplete(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration, jobId string) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      getJobStatusRefreshFunc(client, jobId),
		Timeout:      timeout,
		Delay:        120 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for SWR job (%s) to complete: %s", jobId, err)
	}

	return nil
}

func getJobStatusRefreshFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getJobStatusHttpUrl := "v2/{project_id}/jobs/{job_id}"
		getJobStatusPath := client.Endpoint + getJobStatusHttpUrl
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{project_id}", client.ProjectID)
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{job_id}", jobId)

		getJobStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		}
		getJobStatusResp, err := client.Request("GET", getJobStatusPath, &getJobStatusOpt)
		if err != nil {
			return nil, "ERROR", err
		}

		getJobStatusRespBody, err := utils.FlattenResponse(getJobStatusResp)
		if err != nil {
			return nil, "ERROR", err
		}

		status := utils.PathSearch("status", getJobStatusRespBody, "")
		if status == "Success" {
			return getJobStatusRespBody, "SUCCESS", nil
		}
		if status == "Failed" {
			return getJobStatusRespBody, "FAILED", nil
		}

		return getJobStatusRespBody, "PENDING", nil
	}
}

func resourceSwrEnterpriseInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	getHttpUrl := "v2/{project_id}/instances/{instance_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SWR instance")
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("spec", utils.PathSearch("spec", getRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", getRespBody, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", getRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", getRespBody, nil)),
		d.Set("obs_encrypt", utils.PathSearch("obs_encrypt", getRespBody, nil)),
		d.Set("encrypt_type", utils.PathSearch("encrypt_type", getRespBody, nil)),
		d.Set("obs_bucket_name", utils.PathSearch("obs_bucket_name", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("anonymous_access", utils.PathSearch("anonymous_access", getRespBody, nil)),
		d.Set("version", utils.PathSearch("version", getRespBody, nil)),
		d.Set("charge_mode", utils.PathSearch("charge_mode", getRespBody, nil)),
		d.Set("access_address", utils.PathSearch("access_address", getRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_at", getRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", getRespBody, nil)),
		d.Set("expires_at", utils.PathSearch("expires_at", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("user_def_obs", utils.PathSearch("user_def_obs", getRespBody, nil)),
		d.Set("vpc_name", utils.PathSearch("vpc_name", getRespBody, nil)),
		d.Set("vpc_cidr", utils.PathSearch("vpc_cidr", getRespBody, nil)),
		d.Set("subnet_name", utils.PathSearch("subnet_name", getRespBody, nil)),
		d.Set("subnet_cidr", utils.PathSearch("subnet_cidr", getRespBody, nil)),
	)

	configuration, err := getSwrEnterpriseInstanceConfiguration(client, d)
	if err != nil {
		log.Printf("error retrieving SWR instance configuration: %s", err)
	} else {
		mErr = multierror.Append(mErr,
			d.Set("anonymous_access", utils.PathSearch("anonymous_access", configuration, nil)),
		)
	}

	publicAccessControl, err := getSwrEnterpriseInstancepPublicAccessControl(client, d)
	if err != nil {
		log.Printf("error retrieving SWR instance public access control infos: %s", err)
	} else {
		mErr = multierror.Append(mErr,
			d.Set("public_network_access_control_status", utils.PathSearch("status", publicAccessControl, nil)),
			d.Set("public_network_access_white_ip_list", flattenSwrEnterpriseInstancepPublicAccessControlIpList(publicAccessControl)),
		)
	}

	if err := utils.SetResourceTagsToState(d, client, "instances", d.Id()); err != nil {
		mErr = multierror.Append(mErr, err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func getSwrEnterpriseInstanceConfiguration(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v2/{project_id}/instances/{instance_id}/configurations"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Id())
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

func flattenSwrEnterpriseInstancepPublicAccessControlIpList(resp interface{}) []interface{} {
	if rawParams, ok := utils.PathSearch("ip_list", resp, make([]interface{}, 0)).([]interface{}); ok && len(rawParams) > 0 {
		result := make([]interface{}, 0, len(rawParams))
		for _, params := range rawParams {
			m := map[string]interface{}{
				"description": utils.PathSearch("description", params, nil),
				"ip":          utils.PathSearch("ip", params, nil),
			}
			result = append(result, m)
		}

		return result
	}

	return nil
}

func resourceSwrEnterpriseInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	if d.HasChanges("anonymous_access") {
		if err := updateSwrEnterpriseInstanceAnonymousAccess(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		oRaw, nRaw := d.GetChange("tags")
		oMap := oRaw.(map[string]interface{})
		nMap := nRaw.(map[string]interface{})
		// remove old tags
		if len(oMap) > 0 {
			err := deleteSwrEnterpriseInstanceTags(client, d, client.ProjectID, oMap)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		// add new tags
		if len(nMap) > 0 {
			err := addSwrEnterpriseInstanceTags(client, d, client.ProjectID, nMap)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("public_network_access_control_status") {
		n := d.Get("public_network_access_control_status")
		if d.HasChange("public_network_access_white_ip_list") && n == "Disable" {
			if err := updateSwrEnterpriseInstancePublicNetworkAccessWhiteIpList(client, d); err != nil {
				return diag.FromErr(err)
			}
		}

		if err := updateSwrEnterpriseInstancePublicNetworkAccessControl(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}

		if d.HasChange("public_network_access_white_ip_list") && n == "Enable" {
			if err := updateSwrEnterpriseInstancePublicNetworkAccessWhiteIpList(client, d); err != nil {
				return diag.FromErr(err)
			}
		}
	} else if d.HasChange("public_network_access_white_ip_list") {
		if err := updateSwrEnterpriseInstancePublicNetworkAccessWhiteIpList(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSwrEnterpriseInstanceRead(ctx, d, meta)
}

func addSwrEnterpriseInstanceTags(client *golangsdk.ServiceClient, d *schema.ResourceData, projectId string,
	tags map[string]interface{}) error {
	httpUrl := "v2/{project_id}/{resource_type}/{resource_id}/tags/create"
	addPath := client.Endpoint + httpUrl
	addPath = strings.ReplaceAll(addPath, "{project_id}", projectId)
	addPath = strings.ReplaceAll(addPath, "{resource_type}", "instances")
	addPath = strings.ReplaceAll(addPath, "{resource_id}", d.Id())

	addOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		JSONBody: buildUpdateSwrEnterpriseInstanceTagsBodyParams(tags),
	}

	_, err := client.Request("POST", addPath, &addOpt)
	if err != nil {
		return fmt.Errorf("error adding SWR enterprise instance tags: %s", err)
	}

	return nil
}

func deleteSwrEnterpriseInstanceTags(client *golangsdk.ServiceClient, d *schema.ResourceData, projectId string,
	tags map[string]interface{}) error {
	httpUrl := "v2/{project_id}/{resource_type}/{resource_id}/tags/delete"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", projectId)
	deletePath = strings.ReplaceAll(deletePath, "{resource_type}", "instances")
	deletePath = strings.ReplaceAll(deletePath, "{resource_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		JSONBody: buildUpdateSwrEnterpriseInstanceTagsBodyParams(tags),
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return fmt.Errorf("error deleting SWR enterprise instance tags: %s", err)
	}

	return nil
}

func buildUpdateSwrEnterpriseInstanceTagsBodyParams(tags map[string]interface{}) interface{} {
	bodyParams := map[string]interface{}{
		"tags": utils.ExpandResourceTagsMap(tags),
	}
	return bodyParams
}

func updateSwrEnterpriseInstanceAnonymousAccess(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updateHttpUrl := "v2/{project_id}/instances/{instance_id}/configurations"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"anonymous_access": d.Get("anonymous_access"),
		},
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating SWR instance configuration anonymous access: %s", err)
	}

	return nil
}

func updateSwrEnterpriseInstancePublicNetworkAccessControl(ctx context.Context, client *golangsdk.ServiceClient,
	d *schema.ResourceData) error {
	updateHttpUrl := "v2/{project_id}/instances/{instance_id}/endpoint-policy"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"enable": d.Get("public_network_access_control_status").(string) == "Enable",
		},
	}

	_, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating SWR instance public network access control status: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"SUCCESS"},
		Refresh: func() (interface{}, string, error) {
			publicAccessControl, err := getSwrEnterpriseInstancepPublicAccessControl(client, d)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status", publicAccessControl, "")
			if status == d.Get("public_network_access_control_status").(string) {
				return publicAccessControl, "SUCCESS", nil
			}
			if status == "EnableFailed" || status == "DisableFailed" {
				return publicAccessControl, "FAILED", nil
			}

			return publicAccessControl, "PENDING", nil
		},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for updating SWR instance network access control status to be completed: %s", err)
	}

	return nil
}

func getSwrEnterpriseInstancepPublicAccessControl(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v2/{project_id}/instances/{instance_id}/endpoint-policy"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Id())
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

func updateSwrEnterpriseInstancePublicNetworkAccessWhiteIpList(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updateHttpUrl := "v2/{project_id}/instances/{instance_id}/endpoint-policy"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"ip_list": buildSwrEnterpriseInstancePublicNetworkAccessWhiteIpList(d),
		},
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating SWR instance public network access white IP list: %s", err)
	}

	return nil
}

func buildSwrEnterpriseInstancePublicNetworkAccessWhiteIpList(d *schema.ResourceData) interface{} {
	rawParams := d.Get("public_network_access_white_ip_list").(*schema.Set).List()
	if len(rawParams) == 0 {
		return []map[string]interface{}{}
	}

	rst := make([]map[string]interface{}, 0, len(rawParams))
	for _, p := range rawParams {
		if params, ok := p.(map[string]interface{}); ok {
			m := map[string]interface{}{
				"ip":          params["ip"],
				"description": utils.ValueIgnoreEmpty(params["description"]),
			}
			rst = append(rst, m)
		}
	}

	return rst
}

func resourceSwrEnterpriseInstanceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	deleteHttpUrl := "v2/{project_id}/instances/{instance_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"delete_obs": d.Get("delete_obs"),
			"delete_dns": d.Get("delete_dns"),
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SWR instance")
	}

	return nil
}
