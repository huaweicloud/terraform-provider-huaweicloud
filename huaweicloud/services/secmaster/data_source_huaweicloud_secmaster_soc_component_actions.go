package secmaster

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Secmaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/components/{component_id}/action
func DataSourceSocComponentActions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSocComponentActionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"component_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: schemaSocComponentActionData(),
				},
			},
		},
	}
}

func schemaSocComponentActionData() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"action_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"action_desc": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"action_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"create_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"creator_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"can_update": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"action_version_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"action_version_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"action_version_number": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"action_enable": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceSocComponentActionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/components/{component_id}/action"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{component_id}", d.Get("component_id").(string))
	requestPath += fmt.Sprintf("?enabled=%v", d.Get("enabled").(bool))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&limit=200&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster soc component actions: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		result = append(result, dataResp...)
		offset += len(dataResp)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenSocComponentActionsData(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSocComponentActionsData(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"action_name":           utils.PathSearch("action_name", v, nil),
			"action_desc":           utils.PathSearch("action_desc", v, nil),
			"action_type":           utils.PathSearch("action_type", v, nil),
			"create_time":           utils.PathSearch("create_time", v, nil),
			"creator_name":          utils.PathSearch("creator_name", v, nil),
			"can_update":            utils.PathSearch("can_update", v, nil),
			"action_version_id":     utils.PathSearch("action_version_id", v, nil),
			"action_version_name":   utils.PathSearch("action_version_name", v, nil),
			"action_version_number": utils.PathSearch("action_version_number", v, nil),
			"action_enable":         utils.PathSearch("action_enable", v, nil),
		})
	}

	return rst
}
