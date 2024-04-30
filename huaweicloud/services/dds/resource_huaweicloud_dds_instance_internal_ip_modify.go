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

// @API DDS POST /v3/{project_id}/instances/{instance_id}/modify-internal-ip
// @API DDS GET /v3/{project_id}/instances
// @API DDS GET /v3/{project_id}/jobs
func ResourceDDSInstanceModifyIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDDSInstanceModifyIPCreate,
		ReadContext:   resourceDDSInstanceModifyIPRead,
		UpdateContext: resourceDDSInstanceModifyIPUpdate,
		DeleteContext: resourceDDSInstanceModifyIPDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
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
			"new_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceDDSInstanceModifyIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DDS client: %s ", err)
	}
	instId := d.Get("instance_id").(string)
	nodeId := d.Get("node_id").(string)

	// modify instance internal IP
	err = modifyInstanceInternalIP(ctx, client, d.Timeout(schema.TimeoutCreate), instId, nodeId, d.Get("new_ip").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instId + "/" + nodeId)

	return resourceDDSInstanceModifyIPRead(ctx, d, meta)
}

func resourceDDSInstanceModifyIPRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
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
	jsonPaths := fmt.Sprintf("instances|[0].groups[*].nodes[?id=='%s'][]|[0].private_ip", nodeID)
	privateIP := utils.PathSearch(jsonPaths, getInstanceInfoRespBody, "")
	if privateIP.(string) == "" {
		return diag.Errorf("error getting private ip of node(%s)", nodeID)
	}

	mErr := multierror.Append(
		d.Set("new_ip", privateIP),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting fields: %s", err)
	}

	return nil
}

func resourceDDSInstanceModifyIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DDS client: %s ", err)
	}
	instID := d.Get("instance_id").(string)
	nodeID := d.Get("node_id").(string)
	newIP := d.Get("new_ip").(string)

	if d.HasChange("new_ip") {
		err = modifyInstanceInternalIP(ctx, client, d.Timeout(schema.TimeoutUpdate), instID, nodeID, newIP)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDDSInstanceModifyIPRead(ctx, d, meta)
}

func resourceDDSInstanceModifyIPDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting modify internal IP resource is not supported. The resource is only removed from the state," +
		" the instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func modifyInstanceInternalIP(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	instID, nodeID, newIP string) error {
	modifyHttpUrl := "v3/{project_id}/instances/{instance_id}/modify-internal-ip"
	modifyPath := client.Endpoint + modifyHttpUrl
	modifyPath = strings.ReplaceAll(modifyPath, "{project_id}", client.ProjectID)
	modifyPath = strings.ReplaceAll(modifyPath, "{instance_id}", instID)
	modifyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"node_id": nodeID,
			"new_ip":  newIP,
		},
	}

	// retry
	retryFunc := func() (interface{}, bool, error) {
		resp, err := client.Request("POST", modifyPath, &modifyOpt)
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	modifyResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instID),
		WaitTarget:   []string{"normal"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error modifying internal ip for instance(%s) node(%s): %s", instID, nodeID, err)
	}

	// get job ID
	modifyRespBody, err := utils.FlattenResponse(modifyResp.(*http.Response))
	if err != nil {
		return fmt.Errorf("error flatten response: %s", err)
	}
	jobID := utils.PathSearch("job_id", modifyRespBody, "")
	if jobID.(string) == "" {
		return fmt.Errorf("error modifying internal ip for instance(%s) node(%s): %s", instID, nodeID, "job_id not found")
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

func resourceDDSInstanceNodeImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, must be <instance_id>/<node_id>")
	}

	d.Set("instance_id", parts[0])
	d.Set("node_id", parts[1])

	return []*schema.ResourceData{d}, nil
}
