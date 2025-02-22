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

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/lgsconf/submit
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/lgsconf/confdetail
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/lgsconf/update
// @API CSS DELETE /v1.0/{project_id}/clusters/{cluster_id}/lgsconf/delete
func ResourceLogstashConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogstashConfigurationCreate,
		ReadContext:   resourceLogstashConfigurationRead,
		UpdateContext: resourceLogstashConfigurationUpdate,
		DeleteContext: resourceLogstashConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceLogstashConfigurationImportState,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"conf_content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"setting": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"workers": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"batch_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"batch_delay_ms": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"queue_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"queue_check_point_writes": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"queue_max_bytes_mb": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"sensitive_words": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
	}
}

func resourceLogstashConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	clusterId := d.Get("cluster_id").(string)
	createLogstashConfigHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/lgsconf/submit"
	createLogstashConfigPath := client.Endpoint + createLogstashConfigHttpUrl
	createLogstashConfigPath = strings.ReplaceAll(createLogstashConfigPath, "{project_id}", client.ProjectID)
	createLogstashConfigPath = strings.ReplaceAll(createLogstashConfigPath, "{cluster_id}", clusterId)

	createLogstashConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createLogstashConfigOpt.JSONBody = utils.RemoveNil(buildCreateLogstashCnfParameters(d))

	_, err = client.Request("POST", createLogstashConfigPath, &createLogstashConfigOpt)
	if err != nil {
		return diag.Errorf("error creating CSS logstash cluster configuration: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", clusterId, d.Get("name").(string)))

	checkErr := configFileStatusCheck(ctx, d, client)
	if checkErr != nil {
		return diag.FromErr(checkErr)
	}

	return resourceLogstashConfigurationRead(ctx, d, meta)
}

func buildCreateLogstashCnfParameters(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name").(string),
		"confContent": d.Get("conf_content").(string),
		"setting": map[string]interface{}{
			"queueType":             d.Get("setting.0.queue_type").(string),
			"workers":               utils.IntIgnoreEmpty(d.Get("setting.0.workers").(int)),
			"batchSize":             utils.IntIgnoreEmpty(d.Get("setting.0.batch_size").(int)),
			"batchDelayMs":          utils.IntIgnoreEmpty(d.Get("setting.0.batch_delay_ms").(int)),
			"queueCheckPointWrites": utils.IntIgnoreEmpty(d.Get("setting.0.queue_check_point_writes").(int)),
			"queueMaxBytesMb":       utils.IntIgnoreEmpty(d.Get("setting.0.queue_max_bytes_mb").(int)),
		},
	}

	sensitiveWordsRaw := d.Get("sensitive_words").([]interface{})
	if len(sensitiveWordsRaw) > 0 {
		bodyParams["sensitiveWords"] = sensitiveWordsRaw
	}

	return bodyParams
}

func resourceLogstashConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	var mErr *multierror.Error

	resp, err := GetLogstashConfigDetails(client, d.Get("cluster_id").(string), d.Get("name").(string))
	if err != nil {
		// "CSS.0001": Incorrect parameters. Status code is 400.
		// This error code is a general parameter error identification code.
		// It needs to match the corresponding error message to determine whether to convert it from 400 error to 404 error. e.g.
		// {"errCode": "CSS.0001","externalMessage": "CSS.0001 : Incorrect parameters. (conf not exist)"}
		// Use the string (conf not exist) to confirm that it needs to be converted to a 404 error
		err = common.ConvertExpected400ErrInto404Err(err, "errCode", "CSS.0001")
		// "CSS.0015": The cluster does not exist. Status code is 403.
		// {"errCode": "CSS.0015","externalMessage": "No resources are found or the access is denied."}
		err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
		return common.CheckDeletedDiag(d, err, "error querying CSS logstash cluster configuration")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", resp, nil)),
		d.Set("conf_content", utils.PathSearch("confContent", resp, nil)),
		d.Set("status", utils.PathSearch("status", resp, nil)),
		d.Set("updated_at", utils.PathSearch("updateAt", resp, nil)),
		d.Set("setting", flattenLogstashConfSetting(utils.PathSearch("setting", resp, nil))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenLogstashConfSetting(setting interface{}) []map[string]interface{} {
	if setting == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"workers":                  int(utils.PathSearch("workers", setting, float64(0)).(float64)),
			"batch_size":               int(utils.PathSearch("batchSize", setting, float64(0)).(float64)),
			"batch_delay_ms":           int(utils.PathSearch("batchDelayMs", setting, float64(0)).(float64)),
			"queue_type":               utils.PathSearch("queueType", setting, nil),
			"queue_check_point_writes": int(utils.PathSearch("queueCheckPointWrites", setting, float64(0)).(float64)),
			"queue_max_bytes_mb":       int(utils.PathSearch("queueMaxBytesMb", setting, float64(0)).(float64)),
		},
	}
}

func resourceLogstashConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	clusterId := d.Get("cluster_id").(string)
	updateLogstashConfigHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/lgsconf/update"
	updateLogstashConfigPath := client.Endpoint + updateLogstashConfigHttpUrl
	updateLogstashConfigPath = strings.ReplaceAll(updateLogstashConfigPath, "{project_id}", client.ProjectID)
	updateLogstashConfigPath = strings.ReplaceAll(updateLogstashConfigPath, "{cluster_id}", clusterId)

	updateLogstashConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateLogstashConfigOpt.JSONBody = utils.RemoveNil(buildCreateLogstashCnfParameters(d))
	_, err = client.Request("POST", updateLogstashConfigPath, &updateLogstashConfigOpt)
	if err != nil {
		diag.Errorf("error updating CSS logstash cluster configuration: %s", err)
	}

	checkErr := configFileStatusCheck(ctx, d, client)
	if checkErr != nil {
		return diag.FromErr(checkErr)
	}

	return resourceLogstashConfigurationRead(ctx, d, meta)
}

func resourceLogstashConfigurationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	deleteLogstashConfigHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/lgsconf/delete"
	deleteLogstashConfigPath := client.Endpoint + deleteLogstashConfigHttpUrl
	deleteLogstashConfigPath = strings.ReplaceAll(deleteLogstashConfigPath, "{project_id}", client.ProjectID)
	deleteLogstashConfigPath = strings.ReplaceAll(deleteLogstashConfigPath, "{cluster_id}", d.Get("cluster_id").(string))

	deleteLogstashConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"name": d.Get("name").(string),
		},
	}
	_, err = client.Request("DELETE", deleteLogstashConfigPath, &deleteLogstashConfigOpt)
	if err != nil {
		diag.Errorf("error updating CSS logstash cluster configuration: %s", err)
	}
	if err != nil {
		// 1. "CSS.0001" : Incorrect parameters. Status code is 400.
		// This error code is a general parameter error identification code.
		// It needs to match the corresponding error message to determine whether to convert it from 400 error to 404 error. e.g.
		// {"errCode": "CSS.0001","externalMessage": "CSS.0001 : Incorrect parameters. (conf not exist)"}
		// Use the string (conf not exist) to confirm that it needs to be converted to a 404 error
		err = common.ConvertExpected400ErrInto404Err(err, "errCode", "CSS.0001")
		// "CSS.0015": The cluster does not exist. Status code is 403.
		// {"errCode": "CSS.0015","externalMessage": "No resources are found or the access is denied."}
		err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
		return common.CheckDeletedDiag(d, err, "error deleting CSS logstash cluster configuration")
	}

	return nil
}

func configFileStatusCheck(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"checking"},
		Target:  []string{"available", "unavailable"},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetLogstashConfigDetails(client, d.Get("cluster_id").(string), d.Get("name").(string))
			if err != nil {
				return nil, "failed", err
			}
			return resp, utils.PathSearch("status", resp, "").(string), err
		},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the CSS logstash cluster config to check completed: %s", err)
	}
	return nil
}

func GetLogstashConfigDetails(client *golangsdk.ServiceClient, clusterId, confName string) (interface{}, error) {
	getLogstashConfigDetailsHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/lgsconf/confdetail?name={confName}"
	getLogstashConfigDetailsPath := client.Endpoint + getLogstashConfigDetailsHttpUrl
	getLogstashConfigDetailsPath = strings.ReplaceAll(getLogstashConfigDetailsPath, "{project_id}", client.ProjectID)
	getLogstashConfigDetailsPath = strings.ReplaceAll(getLogstashConfigDetailsPath, "{cluster_id}", clusterId)
	getLogstashConfigDetailsPath = strings.ReplaceAll(getLogstashConfigDetailsPath, "{confName}", confName)

	getLogstashConfigDetailsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getLogstashConfigDetailsResp, err := client.Request("GET", getLogstashConfigDetailsPath, &getLogstashConfigDetailsOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getLogstashConfigDetailsResp)
}

func resourceLogstashConfigurationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, must be <cluster_id>/<name>")
	}

	mErr := multierror.Append(nil,
		d.Set("cluster_id", parts[0]),
		d.Set("name", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
