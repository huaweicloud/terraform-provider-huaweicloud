package aom

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var uniAgentBatchInstallNonUpdatableParams = []string{"agent_import_param_list", "proxy_region_id",
	"installer_agent_id", "version", "public_net_flag", "icagent_install_flag", "plugin_install_base_param"}

// @API AOM POST /v1/{project_id}/uniagent-console/mainview/batch-import
func ResourceUniAgentBatchInstall() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUniAgentBatchInstallCreate,
		ReadContext:   resourceUniAgentBatchInstallRead,
		UpdateContext: resourceUniAgentBatchInstallUpdate,
		DeleteContext: resourceUniAgentBatchInstallDelete,

		CustomizeDiff: config.FlexibleForceNew(uniAgentBatchInstallNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the target machines to be operated are located.`,
			},

			// Required parameters.
			"agent_import_param_list": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        agentImportParamListSchema(),
				Description: `The list of machine parameters for installing UniAgent.`,
			},
			"proxy_region_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The proxy region ID.`,
			},
			"installer_agent_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The agent ID of the installation machine.`,
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The version number of UniAgent to be installed.`,
			},
			"public_net_flag": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Whether to use public network access.`,
			},

			// Optional parameters.
			"icagent_install_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to install ICAgent plugin.`,
			},
			"plugin_install_base_param": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"install_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The specified ICAgent version to install.`,
						},
					},
				},
				Description: `The basic information for plugin installation.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func agentImportParamListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `The login password of the target machine.`,
			},
			"inner_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The IP address of the target machine.`,
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The login port of the target machine.`,
			},
			"account": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The SSH account of the target machine.`,
			},
			"os_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The operating system type of the target machine.`,
			},
			"agent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The unique value of the agent.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The VPC ID of the target machine.`,
			},
			"coc_cmdb_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The external unique identifier for COC usage.`,
			},
		},
	}
}

func buildPluginInstallBaseParam(pluginParams []interface{}) map[string]interface{} {
	if len(pluginParams) < 1 {
		return nil
	}

	return map[string]interface{}{
		"install_version": utils.ValueIgnoreEmpty(utils.PathSearch("install_version", pluginParams[0], nil)),
	}
}

func buildAgentImportParamList(agentParams []interface{}) []map[string]interface{} {
	if len(agentParams) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(agentParams))
	for _, agentParam := range agentParams {
		result = append(result, map[string]interface{}{
			"account":     utils.PathSearch("account", agentParam, nil),
			"password":    utils.PathSearch("password", agentParam, nil),
			"inner_ip":    utils.PathSearch("inner_ip", agentParam, nil),
			"port":        utils.PathSearch("port", agentParam, nil),
			"os_type":     utils.PathSearch("os_type", agentParam, nil),
			"agent_id":    utils.ValueIgnoreEmpty(utils.PathSearch("agent_id", agentParam, nil)),
			"vpc_id":      utils.ValueIgnoreEmpty(utils.PathSearch("vpc_id", agentParam, nil)),
			"coc_cmdb_id": utils.ValueIgnoreEmpty(utils.PathSearch("coc_cmdb_id", agentParam, nil)),
		})
	}
	return result
}

func buildUniAgentBatchInstallBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"agent_import_param_list":   buildAgentImportParamList(d.Get("agent_import_param_list").([]interface{})),
		"proxy_region_id":           d.Get("proxy_region_id"),
		"installer_agent_id":        d.Get("installer_agent_id"),
		"version":                   d.Get("version"),
		"public_net_flag":           d.Get("public_net_flag"),
		"icagent_install_flag":      utils.ValueIgnoreEmpty(d.Get("icagent_install_flag")),
		"plugin_install_base_param": buildPluginInstallBaseParam(d.Get("plugin_install_base_param").([]interface{})),
	}
}

func parseUniAgentInstallTaskDispatchError(respBody interface{}) error {
	state := utils.PathSearch("state", respBody, false).(bool)
	failNum := int64(utils.PathSearch("fail_num", respBody, 0).(float64))
	errorMsg := utils.PathSearch("error_msg", respBody, "").(string)

	if !state || failNum > 0 || errorMsg != "" {
		return fmt.Errorf("UniAgent install task failed, install failed number: %v, error message: %s", failNum, errorMsg)
	}
	return nil
}

func resourceUniAgentBatchInstallCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/uniagent-console/mainview/batch-import"
	)

	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"region":       region,
		},
		JSONBody: utils.RemoveNil(buildUniAgentBatchInstallBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating UniAgent batch install task: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	err = parseUniAgentInstallTaskDispatchError(respBody)
	if err != nil {
		return diag.Errorf("error dispatching UniAgent batch install task: %s", err)
	}

	return resourceUniAgentBatchInstallRead(ctx, d, meta)
}

func resourceUniAgentBatchInstallRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceUniAgentBatchInstallUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceUniAgentBatchInstallDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for batch installing AOM UniAgents. Deleting this resource
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
