package dws

import (
	"context"
	"encoding/json"
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

var (
	clusterElbAssociateNonUpdatableParams = []string{
		"cluster_id",
		"elb_id",
	}
	clusterElbAssociateExpected400ErrCodes = []string{
		"DWS.0102",
	}
	clusterElbAssociateRetryableErrCodes = []string{
		"DWS.7023",
	}
)

// @API DWS POST /v2/{project_id}/clusters/{cluster_id}/elbs/{elb_id}
// @API DWS GET /v1.0/{project_id}/clusters/{cluster_id}
// @API DWS DELETE /v2/{project_id}/clusters/{cluster_id}/elbs/{elb_id}
func ResourceClusterElbAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterElbAssociateCreate,
		ReadContext:   resourceClusterElbAssociateRead,
		UpdateContext: resourceClusterElbAssociateUpdate,
		DeleteContext: resourceClusterElbAssociateDelete,

		CustomizeDiff: config.FlexibleForceNew(clusterElbAssociateNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the cluster (to which the ELB associated) is located.`,
			},

			// Required parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster to which the ELB is associated.`,
			},
			"elb_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the ELB to be associated with the DWS cluster.`,
			},

			// Attributes.
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The public IP address of the ELB loadbalancer.`,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The private IP address of the ELB loadbalancer.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					}),
			},
		},
	}
}

func refreshClusterElbAssociateFunc(client *golangsdk.ServiceClient, clusterId, elbId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		elb, err := GetClusterAssociatedElbById(client, clusterId, elbId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "Resource Not Found", "PENDING", nil
			}
			return nil, "ERROR", err
		}

		return elb, "COMPLETED", nil
	}
}

// An error will be returned if the cluster status is not allowed to operation.
// Even if the cluster status is NORMAL, this error may be returned in a very short time after the last operation is completed.
// The error information is as follows:
//
//	{
//	 "externalMessage":"Cluster status is not allowed to operation.",
//	 "errCode":"DWS.7023",
//	 "error_code":"DWS.7023",
//	 "error_msg":"Cluster status is not allowed to operation."
//	}
func isClusterElbAssociateRetryableError(err error) bool {
	if apiErr, ok := err.(golangsdk.ErrDefault403); ok {
		var respBody interface{}
		if jsonErr := json.Unmarshal(apiErr.Body, &respBody); jsonErr != nil {
			log.Printf("[WARN] failed to unmarshal the response body: %s", jsonErr)
			return false
		}
		// AS.2033: "You are not allowed to perform the operation when the AS group is in current [xxx] status."
		// AS.0003: "AS group lock conflict."
		errCode := utils.PathSearch("error_code", respBody, "").(string)
		if utils.StrSliceContains(clusterElbAssociateRetryableErrCodes, errCode) {
			return true
		}
	}
	return false
}

func associateClusterElb(ctx context.Context, client *golangsdk.ServiceClient, clusterId, elbId string, timeout time.Duration) error {
	httpUrl := "v2/{project_id}/clusters/{cluster_id}/elbs/{elb_id}"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", clusterId)
	createPath = strings.ReplaceAll(createPath, "{elb_id}", elbId)

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	err := resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, reqErr := client.Request("POST", createPath, &createOpts)
		if reqErr != nil {
			if isClusterElbAssociateRetryableError(reqErr) {
				// Wait for the update to take effect
				// lintignore:R018
				time.Sleep(10 * time.Second)
				return resource.RetryableError(reqErr)
			}
			return resource.NonRetryableError(reqErr)
		}
		return nil
	})
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshClusterElbAssociateFunc(client, clusterId, elbId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	return err
}

// GetClusterAssociatedElbById queries the cluster ELB information and verifies whether the specified ELB is associated
// with the DWS cluster if the elbId is provided, otherwise returns the current associated ELB.
func GetClusterAssociatedElbById(client *golangsdk.ServiceClient, clusterId string, elbId ...string) (interface{}, error) {
	clusterDetail, err := GetClusterInfoByClusterId(client, clusterId)
	if err != nil {
		return nil, err
	}

	elb := utils.PathSearch("cluster.elb", clusterDetail, nil)
	log.Printf("[Lance] the current associated ELB: %#v", elb)
	if elb == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1.0/{project_id}/clusters/{cluster_id}",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the ELB is not associated with the DWS cluster (%s)", clusterId)),
			},
		}
	}

	currentElbId := utils.PathSearch("id", elb, "").(string)
	if len(elbId) > 0 && elbId[0] != "" && (currentElbId == "" || currentElbId != elbId[0]) {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1.0/{project_id}/clusters/{cluster_id}",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the ELB (%s) is not associated with the DWS cluster (%s)", elbId, clusterId)),
			},
		}
	}

	return elb, nil
}

func resourceClusterElbAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
		elbId     = d.Get("elb_id").(string)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	err = associateClusterElb(ctx, client, clusterId, elbId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error associating ELB (%s) to the DWS cluster (%s): %s", elbId, clusterId, err)
	}

	d.SetId(clusterId)
	return resourceClusterElbAssociateRead(ctx, d, meta)
}

func resourceClusterElbAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Id()
		elbId     = d.Get("elb_id").(string)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	elb, err := GetClusterAssociatedElbById(client, clusterId, elbId)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error retrieving ELB association for the DWS cluster (%s)", clusterId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("cluster_id", clusterId),
		d.Set("elb_id", utils.PathSearch("id", elb, nil)),
		d.Set("public_ip", utils.PathSearch("public_ip", elb, nil)),
		d.Set("private_ip", utils.PathSearch("private_ip", elb, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceClusterElbAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because all parameters are NonUpdatable.
	return resourceClusterElbAssociateRead(ctx, d, meta)
}

func refreshClusterElbDisassociateFunc(client *golangsdk.ServiceClient, clusterId, elbId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		elb, err := GetClusterAssociatedElbById(client, clusterId, elbId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "Resource Not Found", "COMPLETED", nil
			}
			return nil, "ERROR", err
		}
		return elb, "PENDING", nil
	}
}

func disassociateClusterElb(ctx context.Context, client *golangsdk.ServiceClient, clusterId, elbId string, timeout time.Duration) error {
	httpUrl := "v2/{project_id}/clusters/{cluster_id}/elbs/{elb_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{cluster_id}", clusterId)
	deletePath = strings.ReplaceAll(deletePath, "{elb_id}", elbId)

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	err := resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, reqErr := client.Request("DELETE", deletePath, &deleteOpts)
		if reqErr != nil {
			if isClusterElbAssociateRetryableError(reqErr) {
				// Wait for the update to take effect
				// lintignore:R018
				time.Sleep(10 * time.Second)
				return resource.RetryableError(reqErr)
			}
			return resource.NonRetryableError(reqErr)
		}
		return nil
	})
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshClusterElbDisassociateFunc(client, clusterId, elbId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	return err
}

func resourceClusterElbAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Id()
		elbId     = d.Get("elb_id").(string)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	err = disassociateClusterElb(ctx, client, clusterId, elbId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", clusterElbAssociateExpected400ErrCodes...),
			fmt.Sprintf("error disassociating ELB (%s) from the DWS cluster (%s)", elbId, clusterId))
	}

	return nil
}
