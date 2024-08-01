package css

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

	cssv1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/css/v1"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/css/v1/model"

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
	cssV1Client, err := conf.HcCssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	createCnfOpts := buildCreateLogstashCnfParameters(d)
	clusterId := d.Get("cluster_id").(string)

	_, err = cssV1Client.CreateCnf(&model.CreateCnfRequest{
		ClusterId: clusterId,
		Body:      createCnfOpts,
	})
	if err != nil {
		return diag.Errorf("error creating CSS logstash cluster configuration: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", clusterId, d.Get("name").(string)))

	checkErr := configFileStatusCheck(ctx, d, cssV1Client)
	if checkErr != nil {
		return diag.FromErr(checkErr)
	}

	return resourceLogstashConfigurationRead(ctx, d, meta)
}

func buildCreateLogstashCnfParameters(d *schema.ResourceData) *model.CreateCnfReq {
	createOpts := &model.CreateCnfReq{
		Name:        d.Get("name").(string),
		ConfContent: d.Get("conf_content").(string),
		Setting: &model.Setting{
			QueueType:             d.Get("setting.0.queue_type").(string),
			Workers:               utils.Int32IgnoreEmpty(int32(d.Get("setting.0.workers").(int))),
			BatchSize:             utils.Int32IgnoreEmpty(int32(d.Get("setting.0.batch_size").(int))),
			BatchDelayMs:          utils.Int32IgnoreEmpty(int32(d.Get("setting.0.batch_delay_ms").(int))),
			QueueCheckPointWrites: utils.Int32IgnoreEmpty(int32(d.Get("setting.0.queue_check_point_writes").(int))),
			QueueMaxBytesMb:       utils.Int32IgnoreEmpty(int32(d.Get("setting.0.queue_max_bytes_mb").(int))),
		},
	}

	sensitiveWordsRaw := d.Get("sensitive_words").([]interface{})
	if len(sensitiveWordsRaw) > 0 {
		createOpts.SensitiveWords = utils.ExpandToStringListPointer(sensitiveWordsRaw)
	}

	return createOpts
}

func resourceLogstashConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssV1Client, err := conf.HcCssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	var mErr *multierror.Error

	logstashConfDetail, err := cssV1Client.ShowGetConfDetail(&model.ShowGetConfDetailRequest{
		ClusterId: d.Get("cluster_id").(string),
		Name:      d.Get("name").(string),
	})
	if err != nil {
		// "CSS.0001": Incorrect parameters. Status code is 400.
		// This error code is a general parameter error identification code.
		// It needs to match the corresponding error message to determine whether to convert it from 400 error to 404 error. e.g.
		// {"errCode": "CSS.0001","externalMessage": "CSS.0001 : Incorrect parameters. (conf not exist)"}
		// Use the string (conf not exist) to confirm that it needs to be converted to a 404 error
		err = ConvertExpectedHwSdkErrInto404Err(err, http.StatusBadRequest, "CSS.0001", "conf not exist")
		// "CSS.0015": The cluster does not exist. Status code is 403.
		// {"errCode": "CSS.0015","externalMessage": "No resources are found or the access is denied."}
		err = ConvertExpectedHwSdkErrInto404Err(err, http.StatusForbidden, "CSS.0015", "")
		return common.CheckDeletedDiag(d, err, "error querying CSS logstash cluster configuration")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", logstashConfDetail.Name),
		d.Set("conf_content", logstashConfDetail.ConfContent),
		d.Set("setting", flattenLogstashConfSetting(logstashConfDetail)),
		d.Set("status", logstashConfDetail.Status),
		d.Set("updated_at", logstashConfDetail.UpdateAt),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenLogstashConfSetting(logstashConfDetail *model.ShowGetConfDetailResponse) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	result = append(result, map[string]interface{}{
		"workers":                  logstashConfDetail.Setting.Workers,
		"batch_size":               logstashConfDetail.Setting.BatchSize,
		"batch_delay_ms":           logstashConfDetail.Setting.BatchDelayMs,
		"queue_type":               logstashConfDetail.Setting.QueueType,
		"queue_check_point_writes": logstashConfDetail.Setting.QueueCheckPointWrites,
		"queue_max_bytes_mb":       logstashConfDetail.Setting.QueueMaxBytesMb,
	})

	return result
}

func resourceLogstashConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssV1Client, err := conf.HcCssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	updateCnfOpts := buildCreateLogstashCnfParameters(d)

	_, err = cssV1Client.UpdateCnf(&model.UpdateCnfRequest{
		ClusterId: d.Get("cluster_id").(string),
		Body:      updateCnfOpts,
	})
	if err != nil {
		diag.Errorf("error updating CSS logstash cluster configuration: %s", err)
	}

	checkErr := configFileStatusCheck(ctx, d, cssV1Client)
	if checkErr != nil {
		return diag.FromErr(checkErr)
	}

	return resourceLogstashConfigurationRead(ctx, d, meta)
}

func resourceLogstashConfigurationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssV1Client, err := conf.HcCssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	_, err = cssV1Client.DeleteConf(&model.DeleteConfRequest{
		ClusterId: d.Get("cluster_id").(string),
		Body: &model.DeleteConfReq{
			Name: d.Get("name").(string),
		},
	})
	if err != nil {
		// 1. "CSS.0001" : Incorrect parameters. Status code is 400.
		// This error code is a general parameter error identification code.
		// It needs to match the corresponding error message to determine whether to convert it from 400 error to 404 error. e.g.
		// {"errCode": "CSS.0001","externalMessage": "CSS.0001 : Incorrect parameters. (conf not exist)"}
		// Use the string (conf not exist) to confirm that it needs to be converted to a 404 error
		err = ConvertExpectedHwSdkErrInto404Err(err, http.StatusBadRequest, "CSS.0001", "conf not exist")
		// 2. "CSS.0015": The cluster does not exist. Status code is 403.
		// {"errCode": "CSS.0015","externalMessage": "No resources are found or the access is denied."}
		err = ConvertExpectedHwSdkErrInto404Err(err, http.StatusForbidden, "CSS.0015", "")
		return common.CheckDeletedDiag(d, err, "error deleting CSS logstash cluster configuration")
	}

	return nil
}

func configFileStatusCheck(ctx context.Context, d *schema.ResourceData, cssV1Client *cssv1.CssClient) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"checking"},
		Target:  []string{"available", "unavailable"},
		Refresh: func() (interface{}, string, error) {
			resp, err := cssV1Client.ShowGetConfDetail(&model.ShowGetConfDetailRequest{
				ClusterId: d.Get("cluster_id").(string),
				Name:      d.Get("name").(string),
			})
			if err != nil {
				return nil, "failed", err
			}
			return resp, *resp.Status, err
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
