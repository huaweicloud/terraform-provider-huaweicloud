// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product GaussDB
// ---------------------------------------------------------------

package geminidb

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

// @API GaussDBforNoSQL POST /v3/{project_id}/instances/{instance_id}/nodes/{node_id}/public-ip
// @API GaussDBforNoSQL GET /v3/{project_id}/instances
func ResourceGaussRedisEipAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussRedisEipAssociateCreate,
		ReadContext:   resourceGaussRedisEipAssociateRead,
		DeleteContext: resourceGaussRedisEipAssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

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
				ForceNew:    true,
				Description: `Specifies the ID of a GaussDB Redis instance.`,
			},
			"node_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of a GaussDB Redis node.`,
			},
			"public_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Indicates the EIP address to associate.`,
			},
		},
	}
}

func resourceGaussRedisEipAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createGaussRedisEipAssociate: create GaussDB Redis node EIP associate
	var (
		createGaussRedisEipAssociateHttpUrl = "v3/{project_id}/instances/{instance_id}/nodes/{node_id}/public-ip"
		createGaussRedisEipAssociateProduct = "geminidb"
	)
	createGaussRedisEipAssociateClient, err := cfg.NewServiceClient(createGaussRedisEipAssociateProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB for Redis client: %s", err)
	}

	vpcClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	nodeID := d.Get("node_id").(string)
	publicIP := d.Get("public_ip").(string)
	epsID := "all_granted_eps"
	publicID, err := common.GetEipIDbyAddress(vpcClient, publicIP, epsID)
	if err != nil {
		return diag.Errorf("unable to get ID of public IP %s: %s", publicIP, err)
	}

	createGaussRedisEipAssociatePath := createGaussRedisEipAssociateClient.Endpoint + createGaussRedisEipAssociateHttpUrl
	createGaussRedisEipAssociatePath = strings.ReplaceAll(createGaussRedisEipAssociatePath, "{project_id}",
		createGaussRedisEipAssociateClient.ProjectID)
	createGaussRedisEipAssociatePath = strings.ReplaceAll(createGaussRedisEipAssociatePath, "{instance_id}",
		instanceID)
	createGaussRedisEipAssociatePath = strings.ReplaceAll(createGaussRedisEipAssociatePath, "{node_id}", nodeID)

	createGaussRedisEipAssociateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createGaussRedisEipAssociateOpt.JSONBody = utils.RemoveNil(buildGaussRedisEipAssociateBodyParams("BIND", publicIP,
		publicID))

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err = createGaussRedisEipAssociateClient.Request("POST", createGaussRedisEipAssociatePath,
			&createGaussRedisEipAssociateOpt)
		isRetry, err := handleOperationError(err)
		if isRetry {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.Errorf("error creating EipAssociate: %s", err)
	}

	d.SetId(instanceID + "/" + nodeID)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"BIND_EIP"},
		Target:  []string{"available"},
		Refresh: GaussRedisInstanceUpdateRefreshFunc(createGaussRedisEipAssociateClient, instanceID,
			[]string{"BIND_EIP"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to become ready: %s", instanceID, err)
	}

	return resourceGaussRedisEipAssociateRead(ctx, d, meta)
}

func resourceGaussRedisEipAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getGaussRedisEipAssociate: Query GaussDB Redis node EIP associate
	var (
		getGaussRedisEipAssociateHttpUrl = "v3/{project_id}/instances"
		getGaussRedisEipAssociateProduct = "geminidb"
	)
	getGaussRedisEipAssociateClient, err := cfg.NewServiceClient(getGaussRedisEipAssociateProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB for Redis Client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <instance_id>/<node_id>")
	}
	instanceID := parts[0]
	nodeID := parts[1]

	getGaussRedisEipAssociatePath := getGaussRedisEipAssociateClient.Endpoint + getGaussRedisEipAssociateHttpUrl
	getGaussRedisEipAssociatePath = strings.ReplaceAll(getGaussRedisEipAssociatePath, "{project_id}",
		getGaussRedisEipAssociateClient.ProjectID)

	getGaussRedisEipAssociateQueryParams := buildGetGaussRedisEipAssociateQueryParams(instanceID)
	getGaussRedisEipAssociatePath += getGaussRedisEipAssociateQueryParams

	getGaussRedisEipAssociateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getGaussRedisEipAssociateResp, err := getGaussRedisEipAssociateClient.Request("GET",
		getGaussRedisEipAssociatePath, &getGaussRedisEipAssociateOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving EipAssociate")
	}

	getGaussRedisEipAssociateRespBody, err := utils.FlattenResponse(getGaussRedisEipAssociateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	publicIP := utils.PathSearch(fmt.Sprintf("instances[?id=='%s']|[0].groups[0].nodes[?id=='%s']|[0].public_ip",
		instanceID, nodeID), getGaussRedisEipAssociateRespBody, "")
	if publicIP == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("node_id", nodeID),
		d.Set("public_ip", publicIP),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetGaussRedisEipAssociateQueryParams(instanceID string) string {
	return fmt.Sprintf("?id=%s", instanceID)
}

func resourceGaussRedisEipAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteGaussRedisEipAssociate: Delete GaussDB Redis node EIP associate
	var (
		deleteGaussRedisEipAssociateHttpUrl = "v3/{project_id}/instances/{instance_id}/nodes/{node_id}/public-ip"
		deleteGaussRedisEipAssociateProduct = "geminidb"
	)
	deleteGaussRedisEipAssociateClient, err := cfg.NewServiceClient(deleteGaussRedisEipAssociateProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB for Redis client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	deleteGaussRedisEipAssociatePath := deleteGaussRedisEipAssociateClient.Endpoint + deleteGaussRedisEipAssociateHttpUrl
	deleteGaussRedisEipAssociatePath = strings.ReplaceAll(deleteGaussRedisEipAssociatePath, "{project_id}",
		deleteGaussRedisEipAssociateClient.ProjectID)
	deleteGaussRedisEipAssociatePath = strings.ReplaceAll(deleteGaussRedisEipAssociatePath, "{instance_id}",
		instanceID)
	deleteGaussRedisEipAssociatePath = strings.ReplaceAll(deleteGaussRedisEipAssociatePath, "{node_id}",
		fmt.Sprintf("%v", d.Get("node_id")))

	deleteGaussRedisEipAssociateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	publicIP := d.Get("public_ip").(string)

	deleteGaussRedisEipAssociateOpt.JSONBody = utils.RemoveNil(buildGaussRedisEipAssociateBodyParams("UNBIND",
		publicIP, nil))

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err = deleteGaussRedisEipAssociateClient.Request("POST", deleteGaussRedisEipAssociatePath,
			&deleteGaussRedisEipAssociateOpt)
		isRetry, err := handleOperationError(err)
		if isRetry {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.Errorf("error deleting EipAssociate: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"UNBIND_EIP"},
		Target:  []string{"available"},
		Refresh: GaussRedisInstanceUpdateRefreshFunc(deleteGaussRedisEipAssociateClient, instanceID,
			[]string{"UNBIND_EIP"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to become ready: %s", instanceID, err)
	}
	return diag.FromErr(err)
}

func buildGaussRedisEipAssociateBodyParams(action, publicIP, publicID interface{}) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action":       action,
		"public_ip":    publicIP,
		"public_ip_id": publicID,
	}
	return bodyParams
}
