package css

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

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/lgsconf/start
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/lgsconf/hot-start
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/lgsconf/listpipelines
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/lgsconf/stop
func ResourceLogstashPipeline() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogstashPipelineCreate,
		ReadContext:   resourceLogstashPipelineRead,
		UpdateContext: resourceLogstashPipelineUpdate,
		DeleteContext: resourceLogstashPipelineDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"keep_alive": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"pipelines": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"keep_alive": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"events": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"in": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"filtered": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"out": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceLogstashPipelineCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterId := d.Get("cluster_id").(string)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	startPipelineHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/lgsconf/start"
	startPipelinePath := cssV1Client.Endpoint + startPipelineHttpUrl
	startPipelinePath = strings.ReplaceAll(startPipelinePath, "{project_id}", cssV1Client.ProjectID)
	startPipelinePath = strings.ReplaceAll(startPipelinePath, "{cluster_id}", clusterId)

	startPipelineOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	startPipelineOpt.JSONBody = map[string]interface{}{
		"names":      utils.ExpandToStringList(d.Get("names").(*schema.Set).List()),
		"keep_alive": d.Get("keep_alive").(bool),
	}

	_, err = cssV1Client.Request("POST", startPipelinePath, &startPipelineOpt)
	if err != nil {
		return diag.Errorf("error creating CSS logstash cluster pipeline: %s", err)
	}

	d.SetId(clusterId)

	checkErr := pipelineStatusCheck(ctx, d, cssV1Client, d.Timeout(schema.TimeoutCreate), "WORKING")
	if checkErr != nil {
		return diag.FromErr(checkErr)
	}

	return resourceLogstashPipelineRead(ctx, d, meta)
}

func resourceLogstashPipelineRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	pipelineList, err := getPipelines(d, cssV1Client)
	if err != nil {
		err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
		return common.CheckDeletedDiag(d, err, "error get CSS logstash cluster pipeline")
	}
	if len(pipelineList) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error get CSS logstash cluster pipeline")
	}
	names, pipelines := flattenPipeline(pipelineList)
	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("cluster_id", d.Id()), // it needs to set when import this resource.
		d.Set("names", names),
		d.Set("keep_alive", pipelines[0]["keep_alive"].(bool)),
		d.Set("pipelines", pipelines),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceLogstashPipelineUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	oldRaws, newRaws := d.GetChange("names")
	hotStopNames := utils.ExpandToStringListBySet(oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set)))
	hotStartNames := utils.ExpandToStringListBySet(newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set)))
	if len(hotStopNames) > 0 {
		return diag.Errorf("cannot hot stop a pipeline")
	}
	if len(hotStartNames) > 0 {
		err := addPipeline(ctx, d, cssV1Client, hotStartNames)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceLogstashPipelineRead(ctx, d, meta)
}

func resourceLogstashPipelineDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	stopPipelineHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/lgsconf/stop"
	stopPipelinePath := cssV1Client.Endpoint + stopPipelineHttpUrl
	stopPipelinePath = strings.ReplaceAll(stopPipelinePath, "{project_id}", cssV1Client.ProjectID)
	stopPipelinePath = strings.ReplaceAll(stopPipelinePath, "{cluster_id}", d.Get("cluster_id").(string))

	stopPipelineOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = cssV1Client.Request("POST", stopPipelinePath, &stopPipelineOpt)
	if err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
		// "CSS.5090": In this status, the current operation is not allowed. (stop logstash failed, no pipeline is working.)
		err = common.ConvertExpected400ErrInto404Err(err, "errCode", "CSS.5090")
		return common.CheckDeletedDiag(d, err, "error deleting CSS logstash cluster pipeline")
	}

	checkErr := pipelineStatusCheck(ctx, d, cssV1Client, d.Timeout(schema.TimeoutUpdate), "DELETED")
	if checkErr != nil {
		return diag.FromErr(checkErr)
	}

	return nil
}

func addPipeline(ctx context.Context, d *schema.ResourceData, cssV1Client *golangsdk.ServiceClient,
	hotStartNames []string) error {
	hotStartPipelineHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/lgsconf/hot-start"
	hotStartPipelinePath := cssV1Client.Endpoint + hotStartPipelineHttpUrl
	hotStartPipelinePath = strings.ReplaceAll(hotStartPipelinePath, "{project_id}", cssV1Client.ProjectID)
	hotStartPipelinePath = strings.ReplaceAll(hotStartPipelinePath, "{cluster_id}", d.Get("cluster_id").(string))

	hotStartPipelineOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for _, name := range hotStartNames {
		hotStartPipelineOpt.JSONBody = map[string]interface{}{
			"name":       name,
			"keep_alive": d.Get("keep_alive").(bool),
		}

		_, err := cssV1Client.Request("POST", hotStartPipelinePath, &hotStartPipelineOpt)
		if err != nil {
			return fmt.Errorf("error hot start CSS logstash cluster pipeline: %s", err)
		}

		checkErr := pipelineStatusCheck(ctx, d, cssV1Client, d.Timeout(schema.TimeoutUpdate), "WORKING")
		if checkErr != nil {
			return checkErr
		}
	}

	return nil
}

func pipelineStatusCheck(ctx context.Context, d *schema.ResourceData,
	cssV1Client *golangsdk.ServiceClient, duration time.Duration, waitingTarget string) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{waitingTarget},
		Refresh:      pipelineStateRefreshFunc(d, cssV1Client),
		Timeout:      duration,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the CSS logstash cluster pipeline to completed: %s", err)
	}
	return nil
}

func pipelineStateRefreshFunc(d *schema.ResourceData, cssV1Client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pipelines, err := getPipelines(d, cssV1Client)
		if err != nil {
			return pipelines, "ERROR", err
		}
		if len(pipelines) == 0 {
			return pipelines, "DELETED", nil
		}
		startingPipelines := utils.PathSearch(
			"[?status!='working']", pipelines, make([]interface{}, 0)).([]interface{})
		if len(startingPipelines) == 0 {
			return pipelines, "WORKING", nil
		}
		return pipelines, "PENDING", nil
	}
}

func getPipelines(d *schema.ResourceData, cssV1Client *golangsdk.ServiceClient) ([]interface{}, error) {
	getPipelinesHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/lgsconf/listpipelines"
	getPipelinesPath := cssV1Client.Endpoint + getPipelinesHttpUrl
	getPipelinesPath = strings.ReplaceAll(getPipelinesPath, "{project_id}", cssV1Client.ProjectID)
	getPipelinesPath = strings.ReplaceAll(getPipelinesPath, "{cluster_id}", d.Id())

	getPipelineOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getPipelineResp, err := cssV1Client.Request("GET", getPipelinesPath, &getPipelineOpt)
	if err != nil {
		return nil, err
	}
	getPipelineRespBody, err := utils.FlattenResponse(getPipelineResp)
	if err != nil {
		return nil, err
	}

	pipelines := utils.PathSearch(
		"pipelines|[?status!='stopped']", getPipelineRespBody, make([]interface{}, 0)).([]interface{})

	return pipelines, nil
}

func flattenPipeline(pipelineList []interface{}) ([]string, []map[string]interface{}) {
	if len(pipelineList) == 0 {
		return nil, nil
	}
	var names []string
	rst := make([]map[string]interface{}, 0, len(pipelineList))
	for _, v := range pipelineList {
		names = append(names, utils.PathSearch("name", v, "").(string))
		rst = append(rst, map[string]interface{}{
			"name":       utils.PathSearch("name", v, nil),
			"keep_alive": utils.PathSearch("keepAlive", v, false),
			"events":     flattenPipelineEvent(utils.PathSearch("events", v, nil)),
			"status":     utils.PathSearch("status", v, nil),
			"updated_at": utils.PathSearch("updateAt", v, nil),
		})
	}
	return names, rst
}

func flattenPipelineEvent(events interface{}) []map[string]interface{} {
	if events == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"in":       int(utils.PathSearch("in", events, float64(0)).(float64)),
			"filtered": int(utils.PathSearch("filtered", events, float64(0)).(float64)),
			"out":      int(utils.PathSearch("out", events, float64(0)).(float64)),
		},
	}
	return rst
}
