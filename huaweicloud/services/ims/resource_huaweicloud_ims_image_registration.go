package ims

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

var imageRegistrationNonUpdatableParams = []string{"image_id", "image_url"}

// @API IMS PUT /v1/cloudimages/{image_id}/upload
// @API IMS GET /v1/{project_id}/jobs/{job_id}
// @API IMS GET /v2/cloudimages
// @API IMS DELETE /v2/images/{image_id}
func ResourceImageRegistration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImageRegistrationCreate,
		ReadContext:   resourceImageRegistrationRead,
		UpdateContext: resourceImageRegistrationUpdate,
		DeleteContext: resourceImageRegistrationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(imageRegistrationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			// Attributes.
			"file": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"self": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"visibility": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protected": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"container_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_ram": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_ram": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__os_bit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__os_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__isregistered": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__platform": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_disk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"virtual_env_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__image_source_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__imagetype": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__originalimagename": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__backup_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__productcode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__image_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__data_origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__lazyloading": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__image_displayname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__os_feature_list": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_kvm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_xen": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_largememory": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_diskintensive": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_highperformance": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_xen_gpu_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_kvm_gpu_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_xen_hana": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_kvm_infiniband": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__system_support_market": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"__is_offshelved": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__root_origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__sequence_num": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_fc_inject": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hw_firmware_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hw_vif_multiqueue_enabled": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_arm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_agent_list": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__system__cmkid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__account_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_amd": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_kvm_hi1822_hisriov": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_kvm_hi1822_hivirtionet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_shutdown_timeout": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateImageRegistrationBodyParam(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"image_url": d.Get("image_url"),
	}
}

func resourceImageRegistrationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ims"
		httpUrl = "v1/cloudimages/{image_id}/upload"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{image_id}", d.Get("image_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildCreateImageRegistrationBodyParam(d),
	}

	createResp, err := client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error registering IMS image: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error registering IMS image: job ID is not found in API response")
	}

	imageId, err := waitForImageRegistrationJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for IMS image registration to complete: %s", err)
	}

	d.SetId(imageId)

	return resourceImageRegistrationRead(ctx, d, meta)
}

func waitForImageRegistrationJobCompleted(ctx context.Context, client *golangsdk.ServiceClient, jobId string,
	timeout time.Duration) (string, error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"COMPLETED"},
		Refresh:    imageRegistrationJobStatusRefreshFunc(jobId, client),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "", fmt.Errorf("error waiting for IMS image registration job (%s) to succeed: %s", jobId, err)
	}

	return getImageIdByJob(client, jobId)
}

func imageRegistrationJobStatusRefreshFunc(jobId string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getPath := client.Endpoint + "v1/{project_id}/jobs/{job_id}"
		getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
		getPath = strings.ReplaceAll(getPath, "{job_id}", fmt.Sprintf("%v", jobId))
		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return getResp, "ERROR", fmt.Errorf("error retrieving IMS image registration job: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return getRespBody, "ERROR", err
		}

		status := utils.PathSearch("status", getRespBody, "").(string)
		if status == "SUCCESS" {
			return "SUCCESS", "COMPLETED", nil
		}

		if status == "FAIL" {
			return getRespBody, "COMPLETED", errors.New("the image registration job execution has failed")
		}

		if status == "" {
			return getRespBody, "ERROR", errors.New("status field is not found in API response")
		}

		return getRespBody, "PENDING", nil
	}
}

func getImageIdByJob(client *golangsdk.ServiceClient, jobId string) (string, error) {
	getPath := client.Endpoint + "v1/{project_id}/jobs/{job_id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{job_id}", fmt.Sprintf("%v", jobId))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return "", fmt.Errorf("error retrieving IMS image registration job: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return "", err
	}

	imageId := utils.PathSearch("entities.image_id", getRespBody, "").(string)
	if imageId == "" {
		return "", errors.New("the image ID is not found in API response")
	}

	return imageId, nil
}

func resourceImageRegistrationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ims"
		httpUrl = "v2/cloudimages"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath += fmt.Sprintf("?id=%s", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving IMS image: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	image := utils.PathSearch("images[0]", getRespBody, nil)
	// If the list API return empty, then process `CheckDeleted` logic.
	if image == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving IMS image")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("image_id", utils.PathSearch("id", image, nil)),
		d.Set("file", utils.PathSearch("file", image, nil)),
		d.Set("self", utils.PathSearch("self", image, nil)),
		d.Set("schema", utils.PathSearch("schema", image, nil)),
		d.Set("status", utils.PathSearch("status", image, nil)),
		d.Set("tags", utils.ExpandToStringList(utils.PathSearch("tags", image, make([]interface{}, 0)).([]interface{}))),
		d.Set("visibility", utils.PathSearch("visibility", image, nil)),
		d.Set("name", utils.PathSearch("name", image, nil)),
		d.Set("protected", utils.PathSearch("protected", image, nil)),
		d.Set("container_format", utils.PathSearch("container_format", image, nil)),
		d.Set("min_ram", utils.PathSearch("min_ram", image, nil)),
		d.Set("max_ram", utils.PathSearch("max_ram", image, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", image, nil)),
		d.Set("__os_bit", utils.PathSearch("__os_bit", image, nil)),
		d.Set("__os_version", utils.PathSearch("__os_version", image, nil)),
		d.Set("__description", utils.PathSearch("__description", image, nil)),
		d.Set("disk_format", utils.PathSearch("disk_format", image, nil)),
		d.Set("__isregistered", utils.PathSearch("__isregistered", image, nil)),
		d.Set("__platform", utils.PathSearch("__platform", image, nil)),
		d.Set("__os_type", utils.PathSearch("__os_type", image, nil)),
		d.Set("min_disk", utils.PathSearch("min_disk", image, nil)),
		d.Set("virtual_env_type", utils.PathSearch("virtual_env_type", image, nil)),
		d.Set("__image_source_type", utils.PathSearch("__image_source_type", image, nil)),
		d.Set("__imagetype", utils.PathSearch("__imagetype", image, nil)),
		d.Set("created_at", utils.PathSearch("created_at", image, nil)),
		d.Set("__originalimagename", utils.PathSearch("__originalimagename", image, nil)),
		d.Set("__backup_id", utils.PathSearch("__backup_id", image, nil)),
		d.Set("__productcode", utils.PathSearch("__productcode", image, nil)),
		d.Set("__image_size", utils.PathSearch("__image_size", image, nil)),
		d.Set("__data_origin", utils.PathSearch("__data_origin", image, nil)),
		d.Set("__lazyloading", utils.PathSearch("__lazyloading", image, nil)),
		d.Set("active_at", utils.PathSearch("active_at", image, nil)),
		d.Set("__image_displayname", utils.PathSearch("__image_displayname", image, nil)),
		d.Set("__os_feature_list", utils.PathSearch("__os_feature_list", image, nil)),
		d.Set("__support_kvm", utils.PathSearch("__support_kvm", image, nil)),
		d.Set("__support_xen", utils.PathSearch("__support_xen", image, nil)),
		d.Set("__support_largememory", utils.PathSearch("__support_largememory", image, nil)),
		d.Set("__support_diskintensive", utils.PathSearch("__support_diskintensive", image, nil)),
		d.Set("__support_highperformance", utils.PathSearch("__support_highperformance", image, nil)),
		d.Set("__support_xen_gpu_type", utils.PathSearch("__support_xen_gpu_type", image, nil)),
		d.Set("__support_kvm_gpu_type", utils.PathSearch("__support_kvm_gpu_type", image, nil)),
		d.Set("__support_xen_hana", utils.PathSearch("__support_xen_hana", image, nil)),
		d.Set("__support_kvm_infiniband", utils.PathSearch("__support_kvm_infiniband", image, nil)),
		d.Set("__system_support_market", utils.PathSearch("__system_support_market", image, nil)),
		d.Set("__is_offshelved", utils.PathSearch("__is_offshelved", image, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", image, nil)),
		d.Set("__root_origin", utils.PathSearch("__root_origin", image, nil)),
		d.Set("__sequence_num", utils.PathSearch("__sequence_num", image, nil)),
		d.Set("__support_fc_inject", utils.PathSearch("__support_fc_inject", image, nil)),
		d.Set("hw_firmware_type", utils.PathSearch("hw_firmware_type", image, nil)),
		d.Set("hw_vif_multiqueue_enabled", utils.PathSearch("hw_vif_multiqueue_enabled", image, nil)),
		d.Set("__support_arm", utils.PathSearch("__support_arm", image, nil)),
		d.Set("__support_agent_list", utils.PathSearch("__support_agent_list", image, nil)),
		d.Set("__system__cmkid", utils.PathSearch("__system__cmkid", image, nil)),
		d.Set("__account_code", utils.PathSearch("__account_code", image, nil)),
		d.Set("__support_amd", utils.PathSearch("__support_amd", image, nil)),
		d.Set("__support_kvm_hi1822_hisriov", utils.PathSearch("__support_kvm_hi1822_hisriov", image, nil)),
		d.Set("__support_kvm_hi1822_hivirtionet", utils.PathSearch("__support_kvm_hi1822_hivirtionet", image, nil)),
		d.Set("os_shutdown_timeout", utils.PathSearch("os_shutdown_timeout", image, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceImageRegistrationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceImageRegistrationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ims"
		httpUrl = "v2/images/{image_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{image_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting IMS image: %s", err)
	}

	// Because the delete API always return `204` status code,
	// so we need to call the list query API to check if the image has been successfully deleted.
	err = waitForImageDeleted(ctx, client, d)
	if err != nil {
		return diag.Errorf("error waiting for IMS image to be deleted: %s", err)
	}

	return nil
}

func waitForImageDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			getPath := client.Endpoint + "v2/cloudimages"
			getPath += fmt.Sprintf("?id=%s", d.Id())
			getOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}

			getResp, err := client.Request("GET", getPath, &getOpt)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error retrieving IMS image: %s", err)
			}

			getRespBody, err := utils.FlattenResponse(getResp)
			if err != nil {
				return nil, "ERROR", err
			}

			image := utils.PathSearch("images[0]", getRespBody, nil)
			if image == nil {
				return "SUCCESS", "COMPLETED", nil
			}

			return image, "PENDING", nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}
