package dds

import (
	"context"
	"fmt"
	"net/http"
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

// @API DDS POST /v3/{project_id}/nodes/{node_id}/bind-eip
// @API DDS POST /v3/{project_id}/nodes/{node_id}/unbind-eip
// @API DDS GET /v3/{project_id}/instances
// @API DDS GET /v3/{project_id}/jobs
// @API EIP GET /v1/{project_id}/publicips
func ResourceDDSInstanceBindEIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDDSInstanceBindEIPCreate,
		ReadContext:   resourceDDSInstanceBindEIPRead,
		DeleteContext: resourceDDSInstanceBindEIPDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceDDSInstanceNodeImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDDSInstanceBindEIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.DdsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s ", err)
	}
	vpcClient, err := conf.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	instId := d.Get("instance_id").(string)
	nodeId := d.Get("node_id").(string)
	publicIP := d.Get("public_ip").(string)
	publicID, err := common.GetEipIDbyAddress(vpcClient, publicIP, "all_granted_eps")
	if err != nil {
		return diag.Errorf("error getting EIP ID with public IP %s: %s", publicIP, err)
	}

	// binding eip
	err = bindOrUnbindEIP(ctx, d, client, d.Timeout(schema.TimeoutCreate), publicID, "bind-eip")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instId + "/" + nodeId)

	return resourceDDSInstanceBindEIPRead(ctx, d, meta)
}

func resourceDDSInstanceBindEIPRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.DdsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s ", err)
	}

	instID := d.Get("instance_id").(string)
	getInstanceInfoHttpUrl := "v3/{project_id}/instances?id={instance_id}"
	getInstanceInfoPath := client.Endpoint + getInstanceInfoHttpUrl
	getInstanceInfoPath = strings.ReplaceAll(getInstanceInfoPath, "{project_id}", client.ProjectID)
	getInstanceInfoPath = strings.ReplaceAll(getInstanceInfoPath, "{instance_id}", instID)
	getInstanceInfoOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getInstanceInfoResp, err := client.Request("GET", getInstanceInfoPath, &getInstanceInfoOpt)
	if err != nil {
		return diag.Errorf("error getting instance(%s) info: %s", instID, err)
	}

	getInstanceInfoRespBody, err := utils.FlattenResponse(getInstanceInfoResp)
	if err != nil {
		return diag.Errorf("error flatten response: %s", err)
	}

	nodeID := d.Get("node_id").(string)
	jsonPaths := fmt.Sprintf("instances|[0].groups[*].nodes[?id=='%s'][]|[0].public_ip", nodeID)
	publicIP := utils.PathSearch(jsonPaths, getInstanceInfoRespBody, "")
	if publicIP.(string) == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving public IP")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("public_ip", publicIP),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting fields: %s", err)
	}

	return nil
}

func resourceDDSInstanceBindEIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DDS client: %s ", err)
	}

	// unbinding eip
	err = bindOrUnbindEIP(ctx, d, client, d.Timeout(schema.TimeoutDelete), "", "unbind-eip")
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func bindOrUnbindEIP(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, timeout time.Duration,
	publicID, action string) error {
	instID := d.Get("instance_id").(string)
	nodeID := d.Get("node_id").(string)

	bindOrUnbindEIPHttpUrl := "v3/{project_id}/nodes/{node_id}/{action}"
	bindOrUnbindEIPPath := client.Endpoint + bindOrUnbindEIPHttpUrl
	bindOrUnbindEIPPath = strings.ReplaceAll(bindOrUnbindEIPPath, "{project_id}", client.ProjectID)
	bindOrUnbindEIPPath = strings.ReplaceAll(bindOrUnbindEIPPath, "{node_id}", nodeID)
	bindOrUnbindEIPPath = strings.ReplaceAll(bindOrUnbindEIPPath, "{action}", action)
	bindOrUnbindEIPOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	if action == "bind-eip" {
		bindOrUnbindEIPOpt.JSONBody = map[string]interface{}{
			"public_ip_id": publicID,
			"public_ip":    d.Get("public_ip"),
		}
	}

	// retry
	retryFunc := func() (interface{}, bool, error) {
		resp, err := client.Request("POST", bindOrUnbindEIPPath, &bindOrUnbindEIPOpt)
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	bindOrUnbindEIPResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instID),
		WaitTarget:   []string{"normal"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error %s for instance(%s) node(%s): %s", action, instID, nodeID, err)
	}

	// get job ID
	bindOrUnbindEIPRespBody, err := utils.FlattenResponse(bindOrUnbindEIPResp.(*http.Response))
	if err != nil {
		return fmt.Errorf("error flatten response: %s", err)
	}
	jobID := utils.PathSearch("job_id", bindOrUnbindEIPRespBody, "")
	if jobID.(string) == "" {
		return fmt.Errorf("error %s for instance(%s) node(%s): %s", action, instID, nodeID, "job_id not found")
	}

	// wait for job complete
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed"},
		Refresh:      JobStateRefreshFunc(client, jobID.(string)),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the job (%s) completed: %s ", jobID.(string), err)
	}

	return nil
}
