package lts

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceLtsStruct() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLtsStructTemplateCreate,
		ReadContext:   resourceLtsStructTemplateRead,
		DeleteContext: resourceLtsStructTemplateDelete,
		UpdateContext: resourceLtsStructTemplateUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: ltsStructResourceImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tokenizer": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"demo_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_analysis": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"content": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"field_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_defined_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"index": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"tag_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"content": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"is_analysis": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"demo_log": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLtsStructTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := httpclient_go.NewHttpClientGo(cfg, "lts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("err creating Client； %s", err)
	}
	var url string
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	region := cfg.GetRegion(d)
	opts := entity.StructTemplateRequest{
		LogGroupId:  d.Get("log_group_id").(string),
		LogStreamId: d.Get("log_stream_id").(string),
	}
	if d.Get("template_type").(string) == "custom" {
		url = "v2/" + cfg.GetProjectID(region) + "/lts/struct/template"
		opts.ToDemoFieldsInfo()
		opts.ParseType = "split"
		opts.Tokenizer = " "
		opts.Content = "127.0.0.1 10.142.203.101 8080 [18/Aug/2021:15:14:33 +0800] GET /apm HTTP/1.1 404 86 6"
	} else {
		url = "v3/" + cfg.GetProjectID(region) + "/lts/struct/template"
		opts.TemplateId = d.Get("template_id").(string)
		opts.TemplateType = d.Get("template_type").(string)
		opts.TemplateName = d.Get("template_name").(string)
	}
	client.WithMethod(httpclient_go.MethodPost).WithUrl(url).WithHeader(header).WithBody(opts)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error request creating StructTemplate fields %s: %s", opts.LogGroupId, err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error convert data %s , %s", string(body), err)
	}
	if response.StatusCode == 201 || response.StatusCode == 200 {
		return resourceLtsStructTemplateRead(ctx, d, meta)
	}
	return diag.Errorf("error creating StructTemplate fields %s: %s", opts.LogGroupId, string(body))
}

func resourceLtsStructTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := httpclient_go.NewHttpClientGo(cfg, "lts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	region := cfg.GetRegion(d)
	client.WithMethod(httpclient_go.MethodGet).WithUrl("v2/" +
		cfg.GetProjectID(region) + "/lts/struct/template?logGroupId=" +
		d.Get("log_group_id").(string) + "&logStreamId=" + d.Get("log_stream_id").(string)).
		WithHeader(header)
	resp, err := client.Do()
	body, diags := client.CheckDeletedDiag(d, err, resp, "error StructTemplate read instance")
	if body == nil {
		return diags
	}
	body = body[1 : len(body)-1]
	body2 := strings.Replace(string(body), `\\\`, "**", -1)
	body3 := strings.Replace(body2, `\`, "", -1)
	body4 := strings.Replace(body3, "**", `\`, -1)
	rlt := &entity.ShowStructTemplateResponse{}
	err = json.Unmarshal([]byte(body4), rlt)
	if err != nil {
		return diag.Errorf("error unmarshal body on entity.ShowStructTemplateResponse")
	}
	d.SetId(rlt.Id)
	mErr := multierror.Append(nil,
		d.Set("demo_log", rlt.DemoLog),
		d.Set("log_group_id", rlt.LogGroupId),
		d.Set("log_stream_id", rlt.LogStreamId),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting LtsStructTemplate fields: %s", err)
	}
	return nil
}

func resourceLtsStructTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := httpclient_go.NewHttpClientGo(cfg, "lts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	region := cfg.GetRegion(d)
	structTemplateDeleteRequest := entity.DeleteStructTemplateReqBody{
		Id: d.Id(),
	}
	client.WithMethod(httpclient_go.MethodDelete).WithUrl("v2/" + cfg.GetProjectID(region) + "/lts/struct/template").
		WithHeader(header).WithBody(structTemplateDeleteRequest)
	resp, err := client.Do()
	if err != nil {
		return diag.Errorf("error delete StructTemplate %s: %s", d.Id(), err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return diag.Errorf("error delete StructTemplate %s: %s", d.Id(), string(body))
	}
	if resp.StatusCode == 200 {
		return nil
	}
	return diag.Errorf("error delete StructTemplate %s:  %s", d.Id(), string(body))
}

func resourceLtsStructTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := httpclient_go.NewHttpClientGo(cfg, "lts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	region := cfg.GetRegion(d)
	structTemplateRequest := entity.StructTemplateRequest{
		LogGroupId:   d.Get("log_group_id").(string),
		LogStreamId:  d.Get("log_stream_id").(string),
		TemplateId:   d.Get("template_id").(string),
		TemplateType: d.Get("template_type").(string),
		TemplateName: d.Get("template_name").(string),
	}
	client.WithMethod(httpclient_go.MethodPut).WithUrl("v3/" + cfg.GetProjectID(region) + "/lts/struct/template").
		WithHeader(header).WithBody(structTemplateRequest)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error update StructTemplate fields %s: %s", structTemplateRequest.LogGroupId, err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error convert data %s , %s", string(body), err)
	}
	if response.StatusCode == 201 {
		rlt := &entity.DeleteStructTemplateReqBody{}
		err = json.Unmarshal(body, rlt)
		if err != nil {
			return diag.Errorf("error unmarshal body on entity.DeleteStructTemplateReqBody")
		}
		d.SetId(rlt.Id)
		return nil
	}
	return diag.Errorf("error update StructTemplate fields %s: %s", structTemplateRequest.LogGroupId, err)
}

func ltsStructResourceImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <id>/<log_group_id>/<log_stream_id>")
	}

	d.SetId(parts[0])
	mErr := multierror.Append(nil,
		d.Set("log_group_id", parts[1]),
		d.Set("log_stream_id", parts[2]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
