package servicestage

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v3ConfigurationGroupNonUpdatableParams = []string{"name", "description"}

// @API ServiceStage POST /v3/{project_id}/cas/config-groups
// @API ServiceStage GET /v3/{project_id}/cas/config-groups/{config_group_id}
// @API ServiceStage DELETE /v3/{project_id}/cas/config-groups/{config_group_id}
func ResourceV3ConfigurationGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3ConfigurationGroupCreate,
		ReadContext:   resourceV3ConfigurationGroupRead,
		UpdateContext: resourceV3ConfigurationGroupUpdate,
		DeleteContext: resourceV3ConfigurationGroupDelete,

		CustomizeDiff: config.FlexibleForceNew(v3ConfigurationGroupNonUpdatableParams),

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
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the configuration group.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the configuration group.`,
			},
			"creator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the configuration group.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the configuration group, in RFC3339 format.`,
			},
		},
	}
}

func buildV3ConfigurationGroupCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func resourceV3ConfigurationGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v3/{project_id}/cas/config-groups"
	)
	client, err := cfg.NewServiceClient("servicestage", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: utils.RemoveNil(buildV3ConfigurationGroupCreateBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating configuration group: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	configGroupId := utils.PathSearch("id", respBody, "").(string)
	if configGroupId == "" {
		return diag.Errorf("unable to find the configuration group ID from the API response")
	}
	d.SetId(configGroupId)

	return resourceV3ConfigurationGroupRead(ctx, d, meta)
}

// GetV3ConfigurationGroupById is a method used to get configuration group detail bu its ID.
func GetV3ConfigurationGroupById(client *golangsdk.ServiceClient, configGroupId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/cas/config-groups/{config_group_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{config_group_id}", configGroupId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceV3ConfigurationGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		configGroupId = d.Id()
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	respBody, err := GetV3ConfigurationGroupById(client, configGroupId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving configuration group (%s)", configGroupId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("creator", utils.PathSearch("creator", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", respBody,
			float64(0)).(float64))/1000, false)))

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceV3ConfigurationGroupUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV3ConfigurationGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v3/{project_id}/cas/config-groups/{config_group_id}"
		configGroupId = d.Id()
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{config_group_id}", configGroupId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	// The delete API always returns 200 status code whether config group is exist.
	_, err = client.Request("DELETE", deletePath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting configuration group (%s)", configGroupId))
	}
	return nil
}
