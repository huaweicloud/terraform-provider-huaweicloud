package drs

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/drs/v3/jobs"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceDrsJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDrsJobCreate,
		ReadContext:   resourceDrsJobRead,
		UpdateContext: resourceDrsJobUpdate,
		DeleteContext: resourceDrsJobDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^([A-Za-z][A-Za-z0-9-_\.]*)$`),
						"The name consists of 4 to 50 characters, starting with a letter. "+
							"Only letters, digits, underscores (_) and hyphens (-) are allowed."),
					validation.StringLenBetween(4, 50),
				),
			},

			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"migration", "sync", "cloudDataGuard"}, false),
			},

			"engine_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{"mysql", "mongodb", "cloudDataGuard-mysql",
					"gaussdbv5"}, false),
			},

			"direction": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"up", "down", "non-dbs"}, false),
			},

			"source_db": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     dbInfoSchemaResource(),
			},

			"destination_db": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     dbInfoSchemaResource(),
			},

			"destination_db_readnoly": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},

			"net_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "eip",
				ValidateFunc: validation.StringInSlice([]string{"vpn", "vpc", "eip"}, false),
			},

			"migration_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "FULL_INCR_TRANS",
				ValidateFunc: validation.StringInSlice([]string{"FULL_TRANS", "FULL_INCR_TRANS", "INCR_TRANS"}, false),
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[^!<>&'"\\]*$`),
						"The 'description' has special character"),
					validation.StringLenBetween(1, 256),
				),
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"multi_write": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"expired_days": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  14,
			},

			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"migrate_definer": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},

			"limit_speed": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 3,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"speed": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"start_time": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"end_time": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"tags": common.TagsForceNewSchema(),

			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func dbInfoSchemaResource() *schema.Resource {
	nodeResource := schema.Resource{
		Schema: map[string]*schema.Schema{
			"engine_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"mysql", "mongodb", "gaussdbv5"}, false),
			},

			"ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"password": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"ssl_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"ssl_cert_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ssl_cert_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ssl_cert_check_sum": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ssl_cert_password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}

	return &nodeResource
}

func resourceDrsJobCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DrsV3Client(region)
	if err != nil {
		return diag.Errorf("Error creating DRS v3 client, error=%s", err)
	}

	opts, err := buildCreateParamter(d, client.ProjectID, config.GetEnterpriseProjectID(d))
	if err != nil {
		return diag.FromErr(err)
	}

	rst, err := jobs.Create(client, *opts)
	if err != nil {
		return fmtp.DiagErrorf("Error creating DRS job: %s", err)
	}

	jobId := rst.Results[0].Id

	err = waitingforJobStatus(ctx, client, jobId, "create", d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(jobId)

	valid := testConnections(client, jobId, opts.Jobs[0])
	if !valid {
		return fmtp.DiagErrorf("Test db connection of job=%s failed", jobId)
	}

	err = reUpdateJob(client, jobId, opts.Jobs[0], d.Get("migrate_definer").(bool))
	if err != nil {
		return diag.FromErr(err)
	}

	//configTransSpeed
	if v, ok := d.GetOk("limit_speed"); ok {
		configRaw := v.([]interface{})
		speedLimits := make([]jobs.SpeedLimitInfo, len(configRaw))
		for i, v := range configRaw {
			tmp := v.(map[string]interface{})
			speedLimits[i] = jobs.SpeedLimitInfo{
				Speed: tmp["speed"].(string),
				Begin: tmp["begin_time"].(string),
				End:   tmp["end_time"].(string),
			}
		}
		_, err = jobs.LimitSpeed(client, jobs.BatchLimitSpeedReq{
			SpeedLimits: []jobs.LimitSpeedReq{
				{
					JobId:      jobId,
					SpeedLimit: speedLimits,
				},
			},
		})

		if err != nil {
			return fmtp.DiagErrorf("Limit speed of job=%s failed, error: %s", jobId, err)
		}
	}

	err = preCheck(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	startReq := jobs.StartJobReq{
		Jobs: []jobs.StartInfo{
			{
				JobId:     jobId,
				StartTime: d.Get("start_time").(string),
			},
		},
	}
	_, err = jobs.Start(client, startReq)

	if err != nil {
		return fmtp.DiagErrorf("start DRS job failed,error: %s", err)
	}

	err = waitingforJobStatus(ctx, client, jobId, "start", d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceDrsJobRead(ctx, d, meta)
}

func resourceDrsJobRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DrsV3Client(region)
	if err != nil {
		return diag.Errorf("Error creating DRS v3 client, error: %s", err)
	}

	detailResp, err := jobs.Get(client, jobs.QueryJobReq{Jobs: []string{d.Id()}})
	if err != nil {
		return common.CheckDeletedDiag(d, parseDrsJobErrorToError404(err), "Error retrieving DRS job")
	}
	detail := detailResp.Results[0]

	// net_type is not in detail, so query by list
	listResp, err := jobs.List(client, jobs.ListJobsReq{
		CurPage:   1,
		PerPage:   1,
		Name:      d.Id(),
		DbUseType: detail.DbUseType,
	})

	if err != nil {
		return fmtp.DiagErrorf("Query the job list by jobId=%s, error: %s", d.Id(), err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", detail.Name),
		d.Set("type", detail.DbUseType),
		d.Set("engine_type", detail.InstInfo.EngineType),
		d.Set("direction", detail.JobDirection),
		d.Set("net_type", listResp.Jobs[0].NetType),
		d.Set("public_ip", detail.InstInfo.PublicIp),
		d.Set("private_ip", detail.InstInfo.Ip),
		d.Set("destination_db_readnoly", detail.IsTargetReadonly),
		d.Set("migration_type", detail.TaskType),
		d.Set("description", detail.Description),
		d.Set("multi_write", detail.MultiWrite),
		d.Set("created_at", detail.CreateTime),
		d.Set("status", detail.Status),
		setDbInfoToState(d, detail.SourceEndpoint, "source_db"),
		setDbInfoToState(d, detail.TargetEndpoint, "destination_db"),
	)

	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("Error setting DRS job fields: %s", mErr)
	}

	return nil
}

func resourceDrsJobUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DrsV3Client(region)
	if err != nil {
		return diag.Errorf("Error creating DRS v3 client, error: %s", err)
	}

	detailResp, err := jobs.Get(client, jobs.QueryJobReq{Jobs: []string{d.Id()}})
	if err != nil {
		return common.CheckDeletedDiag(d, parseDrsJobErrorToError404(err), "Error retrieving DRS job")
	}
	detail := detailResp.Results[0]

	if utils.StrSliceContains(
		[]string{"RELEASE_RESOURCE_COMPLETE", "RELEASE_RESOURCE_STARTED", "RELEASE_RESOURCE_FAILED"}, detail.Status) {
		return nil
	}

	updateParams := jobs.UpdateReq{
		Jobs: []jobs.UpdateJobReq{
			{
				JobId:       d.Id(),
				Name:        d.Get("name").(string),
				Description: d.Get("description").(string),
			},
		},
	}

	_, err = jobs.Update(client, updateParams)
	if err != nil {
		return fmtp.DiagErrorf("Update job=%s failed,error: %s", d.Id(), err)
	}

	return resourceDrsJobRead(ctx, d, meta)
}

func resourceDrsJobDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DrsV3Client(region)
	if err != nil {
		return diag.Errorf("Error creating DRS v3 client, error: %s", err)
	}

	detailResp, err := jobs.Get(client, jobs.QueryJobReq{Jobs: []string{d.Id()}})
	if err != nil {
		return common.CheckDeletedDiag(d, parseDrsJobErrorToError404(err), "Error retrieving DRS job")
	}

	// force terminate
	if !utils.StrSliceContains([]string{"CREATE_FAILED", "RELEASE_RESOURCE_COMPLETE", "RELEASE_CHILD_TRANSFER_COMPLETE"},
		detailResp.Results[0].Status) {
		if !d.Get("force_destroy").(bool) {
			return fmtp.DiagErrorf("The job=%s cannot be deleted when it is running. If you want to forcibly delete " +
				"the job please set force_destroy to True.")
		}

		dErr := jobs.Delete(client, jobs.BatchDeleteJobReq{
			Jobs: []jobs.DeleteJobReq{
				{
					DeleteType: jobs.DeleteTypeForceTerminate,
					JobId:      d.Id(),
				},
			},
		})

		if dErr.Err != nil {
			return fmtp.DiagErrorf("Terminate DRS job failed. %q: %s", d.Id(), dErr)
		}

		err = waitingforJobStatus(ctx, client, d.Id(), "terminate", d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	//delete
	dErr := jobs.Delete(client, jobs.BatchDeleteJobReq{
		Jobs: []jobs.DeleteJobReq{
			{
				DeleteType: jobs.DeleteTypeDelete,
				JobId:      d.Id(),
			},
		},
	})
	if dErr.Err != nil {
		return fmtp.DiagErrorf("Delete DRS job failed. %q: %s", d.Id(), dErr)
	}

	d.SetId("")

	return nil
}

func waitingforJobStatus(ctx context.Context, client *golangsdk.ServiceClient, id, statusType string,
	timeout time.Duration) error {
	var pending []string
	var target []string

	switch statusType {
	case "create":
		pending = []string{"CREATING"}
		target = []string{"CONFIGURATION"}
	case "start":
		pending = []string{"STARTJOBING", "WAITING_FOR_START"}
		target = []string{"FULL_TRANSFER_STARTED", "FULL_TRANSFER_COMPLETE", "INCRE_TRANSFER_STARTED"}
	case "terminate":
		pending = []string{"RELEASE_RESOURCE_STARTED"}
		target = []string{"RELEASE_RESOURCE_COMPLETE"}
	}

	stateConf := &resource.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (interface{}, string, error) {
			resp, err := jobs.Status(client, jobs.QueryJobReq{Jobs: []string{id}})
			if err != nil {
				return nil, "", err
			}
			if resp.Count == 0 || resp.Results[0].ErrorCode != "" {
				return resp, "failed", fmtp.Errorf("%s: %s", resp.Results[0].ErrorCode, resp.Results[0].ErrorMessage)
			}

			if resp.Results[0].Status == "CREATE_FAILED" || resp.Results[0].Status == "RELEASE_RESOURCE_FAILED" {
				return resp, "failed", fmtp.Errorf("%s", resp.Results[0].Status)
			}

			return resp, resp.Results[0].Status, nil
		},
		Timeout:      timeout,
		PollInterval: 20 * timeout,
		Delay:        20 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.Errorf("Error waiting for DRS job (%s) to be %s: %s", id, statusType, err)
	}
	return nil
}

func buildCreateParamter(d *schema.ResourceData, projectId, enterpriseProjectID string) (*jobs.BatchCreateJobReq, error) {
	jobDirection := d.Get("direction").(string)

	sourceDb, err := buildDbConfigParamter(d, "source_db", projectId)
	if err != nil {
		return nil, err
	}

	targetDb, err := buildDbConfigParamter(d, "destination_db", projectId)
	if err != nil {
		return nil, err
	}

	var subnetId string
	if jobDirection == "up" {
		if targetDb.InstanceId == "" {
			return nil, fmtp.Errorf("destination_db.0.instance_id is required When diretion is down.")
		}
		subnetId = targetDb.SubnetId

	} else {
		if sourceDb.InstanceId == "" {
			return nil, fmtp.Errorf("source_db.0.instance_id is required When diretion is down.")
		}
		subnetId = sourceDb.SubnetId
	}

	var bindEip bool
	if d.Get("net_type").(string) == "eip" {
		bindEip = true
	}

	job := jobs.CreateJobReq{
		Name:             d.Get("name").(string),
		DbUseType:        d.Get("type").(string),
		EngineType:       d.Get("engine_type").(string),
		JobDirection:     jobDirection,
		NetType:          d.Get("net_type").(string),
		BindEip:          utils.Bool(bindEip),
		IsTargetReadonly: utils.Bool(d.Get("destination_db_readnoly").(bool)),
		TaskType:         d.Get("migration_type").(string),
		Description:      d.Get("description").(string),
		MultiWrite:       utils.Bool(d.Get("multi_write").(bool)),
		ExpiredDays:      fmt.Sprint(d.Get("expired_days").(int)),
		NodeType:         "high",
		SourceEndpoint:   *sourceDb,
		TargetEndpoint:   *targetDb,
		SubnetId:         subnetId,
		Tags:             utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		SysTags:          utils.BuildSysTags(enterpriseProjectID),
	}

	return &jobs.BatchCreateJobReq{Jobs: []jobs.CreateJobReq{job}}, nil
}

func buildDbConfigParamter(d *schema.ResourceData, dbType, projectId string) (*jobs.Endpoint, error) {
	configRaw := d.Get(dbType).([]interface{})[0].(map[string]interface{})
	configs := jobs.Endpoint{
		DbType:          configRaw["engine_type"].(string),
		Ip:              configRaw["ip"].(string),
		DbName:          configRaw["name"].(string),
		DbUser:          configRaw["user"].(string),
		DbPassword:      configRaw["password"].(string),
		DbPort:          golangsdk.IntToPointer(configRaw["port"].(int)),
		InstanceId:      configRaw["instance_id"].(string),
		Region:          configRaw["region"].(string),
		SubnetId:        configRaw["subnet_id"].(string),
		ProjectId:       projectId,
		SslCertPassword: configRaw["ssl_cert_password"].(string),
		SslCertCheckSum: configRaw["ssl_cert_check_sum"].(string),
		SslCertKey:      configRaw["ssl_cert_key"].(string),
		SslCertName:     configRaw["ssl_cert_name"].(string),
		SslLink:         utils.Bool(configRaw["ssl_enabled"].(bool)),
	}
	return &configs, nil
}

func parseDrsJobErrorToError404(respErr error) error {
	var apiError jobs.JobDetailResp

	if errCode, ok := respErr.(golangsdk.ErrDefault400); ok {
		pErr := json.Unmarshal(errCode.Body, &apiError)
		if pErr == nil &&
			(apiError.Results[0].ErrorCode == "DRS.M00289" || apiError.Results[0].ErrorCode == "DRS.M05004") {
			return golangsdk.ErrDefault404(errCode)
		}
	}
	return respErr
}

func setDbInfoToState(d *schema.ResourceData, endpoint jobs.Endpoint, fieldName string) error {
	result := make([]interface{}, 1)
	item := map[string]interface{}{
		"engine_type":        endpoint.DbType,
		"ip":                 endpoint.Ip,
		"port":               endpoint.DbPort,
		"password":           endpoint.DbPassword,
		"user":               endpoint.DbUser,
		"instance_id":        endpoint.InstanceId,
		"name":               endpoint.InstanceName,
		"region":             endpoint.Region,
		"subnet_id":          endpoint.SubnetId,
		"ssl_cert_password":  endpoint.SslCertPassword,
		"ssl_cert_check_sum": endpoint.SslCertCheckSum,
		"ssl_cert_key":       endpoint.SslCertKey,
		"ssl_cert_name":      endpoint.SslCertName,
		"ssl_enabled":        endpoint.SslLink,
	}
	result[0] = item
	//lintignore:R001
	return d.Set(fieldName, result)
}

func testConnections(client *golangsdk.ServiceClient, jobId string, opts jobs.CreateJobReq) (valid bool) {
	reqParams := jobs.TestConnectionsReq{
		Jobs: []jobs.TestEndPoint{
			{
				JobId:        jobId,
				NetType:      opts.NetType,
				EndPointType: "so",
				ProjectId:    client.ProjectID,
				Region:       opts.SourceEndpoint.Region,
				VpcId:        opts.SourceEndpoint.VpcId,
				SubnetId:     opts.SourceEndpoint.SubnetId,
				DbType:       opts.SourceEndpoint.DbType,
				Ip:           opts.SourceEndpoint.Ip,
				DbUser:       opts.SourceEndpoint.DbUser,
				DbPassword:   opts.SourceEndpoint.DbPassword,
				DbPort:       opts.SourceEndpoint.DbPort,
				SslLink:      opts.SourceEndpoint.SslLink,
				InstId:       opts.SourceEndpoint.InstanceId,
			},
			{
				JobId:        jobId,
				NetType:      opts.NetType,
				EndPointType: "ta",
				ProjectId:    client.ProjectID,
				Region:       opts.TargetEndpoint.Region,
				VpcId:        opts.TargetEndpoint.VpcId,
				SubnetId:     opts.TargetEndpoint.SubnetId,
				DbType:       opts.TargetEndpoint.DbType,
				Ip:           opts.TargetEndpoint.Ip,
				DbUser:       opts.TargetEndpoint.DbUser,
				DbPassword:   opts.TargetEndpoint.DbPassword,
				DbPort:       opts.TargetEndpoint.DbPort,
				SslLink:      opts.TargetEndpoint.SslLink,
				InstId:       opts.TargetEndpoint.InstanceId,
			},
		},
	}
	rsp, err := jobs.TestConnections(client, reqParams)
	if err != nil || rsp.Count != 2 {
		logp.Printf("[ERROR]Test connections of job=%s failed,error: %s", jobId, err)
		return false
	}

	valid = rsp.Results[0].Success && rsp.Results[1].Success
	return
}

func reUpdateJob(client *golangsdk.ServiceClient, jobId string, opts jobs.CreateJobReq, migrateDefiner bool) error {
	reqParams := jobs.UpdateReq{
		Jobs: []jobs.UpdateJobReq{
			{
				JobId:            jobId,
				Name:             opts.Name,
				NetType:          opts.NetType,
				EngineType:       opts.EngineType,
				NodeType:         opts.NodeType,
				StoreDbInfo:      true,
				IsRecreate:       utils.Bool(false),
				DbUseType:        opts.DbUseType,
				Description:      opts.Description,
				TaskType:         opts.TaskType,
				JobDirection:     opts.JobDirection,
				IsTargetReadonly: opts.IsTargetReadonly,
				ReplaceDefiner:   &migrateDefiner,
				SourceEndpoint:   &opts.SourceEndpoint,
				TargetEndpoint:   &opts.TargetEndpoint,
			},
		},
	}

	_, err := jobs.Update(client, reqParams)
	if err != nil {
		return fmtp.Errorf("Update job failed,error: %s", err)
	}

	return nil
}

func preCheck(ctx context.Context, client *golangsdk.ServiceClient, jobId string, timeout time.Duration) error {
	_, err := jobs.PreCheckJobs(client, jobs.BatchPrecheckReq{
		Jobs: []jobs.PreCheckInfo{
			{
				JobId:        jobId,
				PrecheckMode: "forStartJob",
			},
		},
	})
	if err != nil {
		return fmtp.Errorf("Start job=%s preCheck failed,error: %s", jobId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"complete"},
		Refresh: func() (interface{}, string, error) {
			resp, err := jobs.CheckResults(client, jobs.QueryPrecheckResultReq{
				Jobs: []string{jobId},
			})
			if err != nil {
				return nil, "", err
			}
			if resp.Count == 0 || resp.Results[0].ErrorCode != "" {
				return resp, "failed", fmtp.Errorf("%s: %s", resp.Results[0].ErrorCode, resp.Results[0].ErrorMsg)
			}

			if resp.Results[0].Process != "100%" {
				return resp, "pending", nil
			}

			if resp.Results[0].TotalPassedRate == "100%" {
				return resp, "complete", nil
			}

			return resp, "failed", fmtp.Errorf("Some preCheck item failed: %s", resp)
		},
		Timeout:      timeout,
		PollInterval: 20 * timeout,
		Delay:        20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.Errorf("Error waiting for DRS job (%s) to be terminate: %s", jobId, err)
	}
	return nil
}
