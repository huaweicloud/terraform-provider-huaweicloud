package aom

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableParams = []string{
	"installer_agent_id", "public_net_flag", "proxy_region_id",
	"inner_ip", "port", "account", "password", "os_type", "vpc_id", "coc_cmdb_id",
	"icagent_install_flag", "icagent_install_version", "access_key", "secret_key",
}

// @API AOM POST /v1/{project_id}/uniagent-console/mainview/batch-import
// @API AOM POST /v1/{project_id}/uniagent-console/upgrade/batch-upgrade
func ResourceUniAgent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUniAgentCreate,
		ReadContext:   resourceUniAgentRead,
		UpdateContext: resourceUniAgentUpdate,
		DeleteContext: resourceUniAgentDelete,

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(nonUpdatableParams),
			func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
				if oldValue, newValue := d.GetChange("version"); oldValue != newValue && oldValue != "" {
					agentID := d.Get("agent_id").(string)
					if agentID == "" {
						return errors.New("only support to update `version` when `agent_id` is not empty")
					}
				}
				return nil
			},
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"installer_agent_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the installer agent ID.`,
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the UniAgent version to be installed.`,
			},
			"public_net_flag": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Specifies the whether to use public network access.`,
			},
			"proxy_region_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the proxy region ID.`,
			},
			// host info
			"inner_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the IP of the host where the UniAgent will be installed.`,
			},
			"port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the login port of the host where the UniAgent will be installed.`,
			},
			"account": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the login account of the host where the UniAgent will be installed.`,
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `Specifies the login password of the host where the UniAgent will be installed.`,
			},
			"os_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the OS type of the host where the UniAgent will be installed.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the VPC ID of the host where the UniAgent will be installed.`,
			},
			"agent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the agent ID of the host where the UniAgent will be installed.`,
			},
			"coc_cmdb_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the COC CMDB ID of the host where the UniAgent will be installed.`,
			},
			// icagent
			"icagent_install_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to install ICAgent.`,
			},
			"icagent_install_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ICAgent version to be installed.`,
			},
			"access_key": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"secret_key"},
				Description:  `Specifies the access key of the IAM account where the host ICAgent is not installed.`,
			},
			"secret_key": {
				Type:         schema.TypeString,
				Sensitive:    true,
				Optional:     true,
				RequiredWith: []string{"access_key"},
				Description:  `Specifies the secret key of the IAM account where the host ICAgent is not installed.`,
			},
		},
	}
}

func resourceUniAgentCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createHttpUrl := "v1/{project_id}/uniagent-console/mainview/batch-import"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"region": region,
		},
		JSONBody: utils.RemoveNil(buildCreateUniAgentBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating AOM UniAgent: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening response: %s", err)
	}

	if !utils.PathSearch("state", createRespBody, false).(bool) {
		return diag.Errorf("error creating AOM UniAgent: %v", createRespBody)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID")
	}
	d.SetId(id)

	return nil
}

func buildCreateUniAgentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"installer_agent_id":        d.Get("installer_agent_id"),
		"version":                   d.Get("version"),
		"public_net_flag":           d.Get("public_net_flag"),
		"proxy_region_id":           d.Get("proxy_region_id"),
		"agent_import_param_list":   buildCreateUniAgentBodyParamsAgentImportParamList(d),
		"icagent_install_flag":      utils.ValueIgnoreEmpty(d.Get("icagent_install_flag")),
		"plugin_install_base_param": buildCreateUniAgentBodyParamsPluginInstallBaseParam(d),
	}

	return bodyParams
}

func buildCreateUniAgentBodyParamsAgentImportParamList(d *schema.ResourceData) []map[string]interface{} {
	bodyParams := map[string]interface{}{
		"inner_ip":    d.Get("inner_ip"),
		"port":        d.Get("port"),
		"account":     d.Get("account"),
		"password":    d.Get("password"),
		"os_type":     d.Get("os_type"),
		"agent_id":    utils.ValueIgnoreEmpty(d.Get("agent_id")),
		"vpc_id":      utils.ValueIgnoreEmpty(d.Get("vpc_id")),
		"coc_cmdb_id": utils.ValueIgnoreEmpty(d.Get("coc_cmdb_id")),
	}

	return []map[string]interface{}{bodyParams}
}

func buildCreateUniAgentBodyParamsPluginInstallBaseParam(d *schema.ResourceData) interface{} {
	bodyParams := utils.RemoveNil(map[string]interface{}{
		"install_version": utils.ValueIgnoreEmpty(d.Get("icagent_install_version")),
		"domain_ak":       utils.ValueIgnoreEmpty(d.Get("access_key")),
		"domain_sk":       utils.ValueIgnoreEmpty(d.Get("secret_key")),
	})

	return bodyParams
}

func resourceUniAgentRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceUniAgentUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	if d.HasChange("version") {
		updateHttpUrl := "v1/{project_id}/uniagent-console/upgrade/batch-upgrade"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildUpdateUniAgentBodyParams(d),
		}

		updateResp, err := client.Request("POST", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating UniAgent: %s", err)
		}
		updateRespBody, err := utils.FlattenResponse(updateResp)
		if err != nil {
			return diag.Errorf("error flattening response: %s", err)
		}

		if !utils.PathSearch("state", updateRespBody, false).(bool) {
			return diag.Errorf("error updating UniAgent: %v", updateResp)
		}
	}

	return nil
}

func buildUpdateUniAgentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"version":    d.Get("version"),
		"agent_list": buildCUpdateUniAgentBodyParamsAgentList(d),
	}

	return bodyParams
}

func buildCUpdateUniAgentBodyParamsAgentList(d *schema.ResourceData) []map[string]interface{} {
	bodyParams := map[string]interface{}{
		"inner_ip": d.Get("inner_ip"),
		"agent_id": d.Get("agent_id"),
	}

	return []map[string]interface{}{bodyParams}
}

func resourceUniAgentDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting AOM UniAgent resource is not supported. The UniAgent resource is only removed from the state," +
		" the UniAgent remains in the host."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
