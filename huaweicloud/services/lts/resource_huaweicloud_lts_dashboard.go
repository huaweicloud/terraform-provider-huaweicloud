package lts

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LTS POST /v2/{project_id}/lts/template-dashboard
// @API LTS GET /v2/{project_id}/dashboards
// @API LTS DELETE /v2/{project_id}/dashboard
func ResourceLtsDashboard() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLtsDashBoardCreate,
		ReadContext:   resourceLtsDashBoardRead,
		DeleteContext: resourceLtsDashBoardDelete,
		UpdateContext: resourceDashBoardUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: ltsResourceImportState,
		},

		Description: "schema: Internal",
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"is_delete_charts": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"title": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_stream_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_title": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"template_type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"last_update_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceLtsDashBoardCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := httpclient_go.NewHttpClientGo(cfg, "lts", region)
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	dashBoardRequest := entity.DashBoardRequest{
		LogGroupId:    d.Get("log_group_id").(string),
		LogGroupName:  d.Get("log_group_name").(string),
		LogStreamId:   d.Get("log_stream_id").(string),
		LogStreamName: d.Get("log_stream_name").(string),
		TemplateTitle: utils.ExpandToStringList(d.Get("template_title").([]interface{})),
		TemplateType:  utils.ExpandToStringList(d.Get("template_type").([]interface{})),
		GroupName:     d.Get("group_name").(string),
	}
	client.WithMethod(httpclient_go.MethodPost).WithUrl("v2/" + cfg.GetProjectID(region) + "/lts/template-dashboard").
		WithHeader(header).WithBody(dashBoardRequest)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error creating LtsDashBoard fields %s: %s", dashBoardRequest.LogGroupId, err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error convert data %s, %s", string(body), err)
	}
	if response.StatusCode == 201 {
		rlt := make([]entity.DashBoard, 0)
		err = json.Unmarshal(body, &rlt)
		if err != nil {
			return diag.Errorf("error convert data %s, %s", string(body), err)
		}
		if len(rlt) == 0 {
			return diag.Errorf("error resource has been created log stream name %s", d.Get("log_stream_name").(string))
		}
		d.SetId(rlt[0].Id)
		return resourceLtsDashBoardRead(ctx, d, meta)
	}
	return diag.Errorf("error creating LtsDashBoard Response %s: %s", dashBoardRequest.LogGroupId, string(body))
}

func resourceLtsDashBoardRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := httpclient_go.NewHttpClientGo(cfg, "lts", region)
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	client.WithMethod(httpclient_go.MethodGet).WithUrl("v2/" + cfg.GetProjectID(region) + "/dashboards?id=" + d.Id()).WithHeader(header)
	response, err := client.Do()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving dashboard")
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	rlt := entity.ReadDashBoardResp{}
	err = json.Unmarshal(body, &rlt)
	d.Set("region", region)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(rlt.Results) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, fmt.Sprintf("error retrieving dashboard %s", d.Id()))
	}
	mErr := multierror.Append(nil,
		d.Set("title", rlt.Results[0].Title),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Lts dashboard fields: %s", err)
	}
	return nil
}

func resourceLtsDashBoardDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := httpclient_go.NewHttpClientGo(cfg, "lts", region)
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	client.WithMethod(httpclient_go.MethodDelete).WithUrl("v2/" + cfg.GetProjectID(region) + "/dashboard?is_delete_charts=" +
		d.Get("is_delete_charts").(string) + "&id=" + d.Id()).WithHeader(header)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error delete LtsDashBoard %s: %s", d.Id(), err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error delete LtsDashBoard %s: %s", d.Id(), err)
	}
	if response.StatusCode == 200 {
		return nil
	}
	return diag.Errorf("error delete LtsDashBoard %s:  %s", d.Id(), string(body))
}

func resourceDashBoardUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := httpclient_go.NewHttpClientGo(cfg, "lts", region)
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	dashBoardRequest := entity.DashBoardRequest{
		LogGroupId:    d.Get("log_group_id").(string),
		LogGroupName:  d.Get("log_group_name").(string),
		LogStreamId:   d.Get("log_stream_id").(string),
		LogStreamName: d.Get("log_stream_name").(string),
		TemplateTitle: utils.ExpandToStringList(d.Get("template_title").([]interface{})),
		TemplateType:  utils.ExpandToStringList(d.Get("template_type").([]interface{})),
		GroupName:     d.Get("group_name").(string),
	}
	client.WithMethod(httpclient_go.MethodPost).WithUrl("v2/" + cfg.GetProjectID(region) + "/lts/template-dashboard").
		WithHeader(header).WithBody(dashBoardRequest)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error update LtsDashBoard fields %s: %s", dashBoardRequest.LogGroupId, err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error convert data %s, %s", string(body), err)
	}
	if response.StatusCode == 200 {
		rlt := make([]entity.DashBoard, 0)
		err = json.Unmarshal(body, &rlt)
		if err != nil {
			return diag.Errorf("error convert data %s, %s", string(body), err)
		}
		d.SetId(rlt[0].Id)
		return nil
	}
	return diag.Errorf("error update LtsDashBoard Response %s: %s", dashBoardRequest.LogGroupId, string(body))
}

func ltsResourceImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 5)
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid format specified for import id, " +
			"must be <id>/<log_group_id>/<log_group_name>/<log_stream_id>/<log_stream_name>")
	}

	d.SetId(parts[0])
	mErr := multierror.Append(nil,
		d.Set("is_delete_charts", "true"),
		d.Set("log_group_id", parts[1]),
		d.Set("log_group_name", parts[2]),
		d.Set("log_stream_id", parts[3]),
		d.Set("log_stream_name", parts[4]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
