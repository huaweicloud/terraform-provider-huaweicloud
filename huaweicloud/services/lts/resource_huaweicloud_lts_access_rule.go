package lts

import (
	"context"
	"encoding/json"
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

// @API LTS POST /v2/{project_id}/lts/aom-mapping
// @API LTS GET /v2/{project_id}/lts/aom-mapping/{rule_id}
// @API LTS PUT /v2/{project_id}/lts/aom-mapping
// @API LTS DELETE /v2/{project_id}/lts/aom-mapping
func ResourceAomMappingRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAomMappingRuleCreate,
		ReadContext:   resourceAomMappingRuleRead,
		DeleteContext: resourceAomMappingRuleDelete,
		UpdateContext: resourceAomMappingRuleUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Description: "schema: Internal",
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rule_name": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: SuppressCaseDiffs,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name_space": {
				Type:     schema.TypeString,
				Required: true,
			},
			"container_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployments": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"files": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"log_stream_info": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_log_group_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"target_log_group_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"target_log_stream_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"target_log_stream_name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildLogStreamOpts(rawRule []interface{}) entity.AomMappingLogStreamInfo {
	s := rawRule[0].(map[string]interface{})
	rst := entity.AomMappingLogStreamInfo{
		TargetLogGroupId:    s["target_log_group_id"].(string),
		TargetLogGroupName:  s["target_log_group_name"].(string),
		TargetLogStreamId:   s["target_log_stream_id"].(string),
		TargetLogStreamName: s["target_log_stream_name"].(string),
	}
	return rst
}

func buildFileOpts(rawRules []interface{}) []entity.AomMappingfilesInfo {
	file := make([]entity.AomMappingfilesInfo, len(rawRules))
	for i, v := range rawRules {
		rawRule := v.(map[string]interface{})
		file[i].FileName = rawRule["file_name"].(string)
		file[i].LogStreamInfo = buildLogStreamOpts(rawRule["log_stream_info"].([]interface{}))
	}
	return file
}

func resourceAomMappingRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := httpclient_go.NewHttpClientGo(cfg, "lts", region)
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}
	aomMappingRequestInfo := entity.AomMappingRequestInfo{
		ProjectId: cfg.GetProjectID(region),
		RuleName:  d.Get("rule_name").(string),
		RuleInfo: entity.AomMappingRuleInfo{
			ClusterId:   d.Get("cluster_id").(string),
			ClusterName: d.Get("cluster_name").(string),
			Namespace:   d.Get("name_space").(string),
			Deployments: utils.ExpandToStringList(d.Get("deployments").([]interface{})),
			Files:       buildFileOpts(d.Get("files").([]interface{})),
		},
	}
	client.WithMethod(httpclient_go.MethodPost).WithUrl("v2/" + cfg.GetProjectID(region) + "/lts/aom-mapping" + "?isBatch=false").
		WithBody(aomMappingRequestInfo)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error creating AomMappingRule %s: %s", aomMappingRequestInfo.RuleName, err)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error convert data %s , %s", string(body), err)
	}
	if response.StatusCode == 201 {
		rlt := make([]entity.AomMappingRuleResp, 0)
		err = json.Unmarshal(body, &rlt)
		if err != nil {
			return diag.Errorf("error convert data %s, %s", string(body), err)
		}
		d.SetId(rlt[0].RuleId)
		return resourceAomMappingRuleRead(ctx, d, meta)
	}
	return diag.Errorf("error AomMappingRule Response %s : %s", aomMappingRequestInfo.RuleName, string(body))
}

func resourceAomMappingRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := httpclient_go.NewHttpClientGo(cfg, "lts", region)
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	client.WithMethod(httpclient_go.MethodGet).WithUrl("v2/" + cfg.GetProjectID(region) + "/lts/aom-mapping/" + d.Id()).WithHeader(header)
	response, err := client.Do()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AOM mapping rule")
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	rlt := make([]entity.AomMappingRequestInfo, 0)
	err = json.Unmarshal(body, &rlt)
	if err != nil {
		return diag.Errorf("error retrieving AomMappingRule %s", d.Id())
	}

	if len(rlt) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error getting AOM mapping rule")
	}

	mErr := multierror.Append(nil,
		d.Set("rule_name", rlt[0].RuleName),
		d.Set("cluster_id", rlt[0].RuleInfo.ClusterId),
		d.Set("cluster_name", rlt[0].RuleInfo.ClusterName),
		d.Set("container_name", rlt[0].RuleInfo.ContainerName),
		d.Set("deployments", rlt[0].RuleInfo.Deployments),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting AomMappingRule fields: %s", err)
	}
	return nil
}

func resourceAomMappingRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := httpclient_go.NewHttpClientGo(cfg, "lts", region)
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	client.WithMethod(httpclient_go.MethodDelete).WithUrl("v2/" + cfg.GetProjectID(region) + "/lts/aom-mapping?id=" + d.Id()).WithHeader(header)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error delete AomMappingRule %s: %s", d.Id(), err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error delete AomMappingRule %s: %s", d.Id(), err)
	}
	if response.StatusCode == 200 {
		return nil
	}
	return diag.Errorf("error delete AomMappingRule %s:  %s", d.Id(), string(body))
}

func resourceAomMappingRuleUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := httpclient_go.NewHttpClientGo(cfg, "lts", region)
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}
	opts := entity.AomMappingRequestInfo{
		ProjectId: cfg.GetProjectID(region),
		RuleId:    d.Id(),
		RuleName:  d.Get("rule_name").(string),
		RuleInfo: entity.AomMappingRuleInfo{
			ClusterId:   d.Get("cluster_id").(string),
			ClusterName: d.Get("cluster_name").(string),
			Namespace:   d.Get("name_space").(string),
			Deployments: utils.ExpandToStringList(d.Get("deployments").([]interface{})),
			Files:       buildFileOpts(d.Get("files").([]interface{})),
		},
	}
	client.WithMethod(httpclient_go.MethodPut).WithUrl("v2/" + cfg.GetProjectID(region) + "/lts/aom-mapping").WithBody(opts)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error update AomMappingRule fields %s: %s", opts.RuleName, err)
	}
	d.SetId(opts.RuleId)

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error update AomMappingRule %s: %s", string(body), err)
	}
	if response.StatusCode == 200 {
		return nil
	}
	return diag.Errorf("error update AomMappingRule %s:  %s", opts.RuleName, string(body))
}

func SuppressCaseDiffs(_, old, new string, _ *schema.ResourceData) bool {
	return strings.Split(old, "_")[0] == strings.Split(new, "_")[0]
}
