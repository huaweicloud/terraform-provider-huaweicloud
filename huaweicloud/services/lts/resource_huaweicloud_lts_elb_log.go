package lts

import (
	"context"
	"encoding/json"
	"io"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
)

// @API ELB POST /v3/{project_id}/elb/logtanks
// @API ELB GET /v3/{project_id}/elb/logtanks/{logtank_id}
// @API ELB PUT /v3/{project_id}/elb/logtanks/{logtank_id}
// @API ELB DELETE /v3/{project_id}/elb/logtanks/{logtank_id}
func ResourceLtsElb() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLtsElbCreate,
		ReadContext:   resourceLtsElbRead,
		DeleteContext: resourceLtsElbDelete,
		UpdateContext: resourceLtsElbUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Description: "schema: Internal",
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_topic_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceLtsElbCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := httpclient_go.NewHttpClientGo(cfg, "elb", region)
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	logTank := entity.CreateLogTankOption{
		LogGroupId:     d.Get("log_group_id").(string),
		LoadBalancerId: d.Get("loadbalancer_id").(string),
		LogTopicId:     d.Get("log_topic_id").(string),
	}
	logTankRequest := entity.CreateLogtankRequestBody{
		Logtank: &logTank,
	}
	client.WithMethod(httpclient_go.MethodPost).WithUrl("v3/" + cfg.GetProjectID(region) + "/elb/logtanks").WithHeader(header).
		WithBody(logTankRequest).WithTransport()
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error creating LogTank fields %s: %s", logTankRequest.Logtank.LogGroupId, err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error convert data %s, %s", string(body), err)
	}
	if response.StatusCode == 201 {
		rlt := &entity.CreateLogtankResponse{}
		err = json.Unmarshal(body, rlt)
		if err != nil {
			return diag.Errorf("error unmarshal body on entity.CreateLogtankResponse")
		}
		d.SetId(rlt.Logtank.ID)
		return resourceLtsElbRead(ctx, d, meta)
	}
	return diag.Errorf("error creating LogTank fields %s: %s", logTankRequest.Logtank.LogGroupId, string(body))
}

func resourceLtsElbRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := httpclient_go.NewHttpClientGo(cfg, "elb", region)
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	client.WithMethod(httpclient_go.MethodGet).WithUrl("v3/" + cfg.GetProjectID(region) +
		"/elb/logtanks/" + d.Id()).WithHeader(header).WithTransport()
	response, err := client.Do()
	body, diags := client.CheckDeletedDiag(d, err, response, "error Elb LogTank read instance")
	if body == nil {
		return diags
	}
	rlt := &entity.CreateLogtankResponse{}
	err = json.Unmarshal(body, rlt)
	if err != nil {
		return diag.Errorf("error retrieving Elb LogTank %s", d.Id())
	}
	mErr := multierror.Append(nil,
		d.Set("loadbalancer_id", rlt.Logtank.LoadBalancerID),
		d.Set("log_group_id", rlt.Logtank.LogGroupID),
		d.Set("log_topic_id", rlt.Logtank.LogTopicID),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Elb LogTank fields: %s", err)
	}
	return nil
}

func resourceLtsElbDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := httpclient_go.NewHttpClientGo(cfg, "elb", region)
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	client.WithMethod(httpclient_go.MethodDelete).WithUrl("v3/" + cfg.GetProjectID(region) +
		"/elb/logtanks/" + d.Id()).WithHeader(header).WithTransport()
	resp, err := client.Do()
	if err != nil {
		return diag.Errorf("error delete LogTank %s: %s", d.Id(), err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.Errorf("error delete LogTank %s: %s", d.Id(), err)
	}
	if resp.StatusCode == 204 {
		return nil
	}
	return diag.Errorf("error delete LogTank %s:  %s", d.Id(), string(body))
}

func resourceLtsElbUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := httpclient_go.NewHttpClientGo(cfg, "elb", region)
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	logTankRequest := entity.CreateLogTankOption{
		LogGroupId: d.Get("log_group_id").(string),
		LogTopicId: d.Get("log_topic_id").(string),
	}
	client.WithMethod(httpclient_go.MethodPut).WithUrl("v3/" + cfg.GetProjectID(region) + "/elb/logtanks/" + d.Id()).
		WithHeader(header).WithBody(logTankRequest).WithTransport()
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error update LogTank fields %s: %s", logTankRequest, err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error update LogTank %s: %s", string(body), err)
	}

	if response.StatusCode == 200 {
		return nil
	}
	return diag.Errorf("error update LogTank %s: %s", d.Id(), string(body))
}
