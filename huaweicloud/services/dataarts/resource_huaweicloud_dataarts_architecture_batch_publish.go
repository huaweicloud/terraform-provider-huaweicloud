package dataarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio POST /v2/{project_id}/design/approvals/batch-publish
func ResourceArchitectureBatchPublish() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArchitectureBatchPublishCreate,
		ReadContext:   resourceArchitectureBatchPublishRead,
		DeleteContext: resourceArchitectureBatchPublishDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of DataArts Studio workspace.",
			},
			"biz_infos": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the list of objects to be published.",
				Elem:        bizInfoSchema(),
			},
			"approver_user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the user ID of the architecture reviewer.",
			},
			"approver_user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the user name of the architecture reviewer.",
			},
			"fast_approval": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies whether to automatically review.",
			},
			"schedule_time": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies scheduling time of the DataArts quality job.",
			},
		},
	}
}

func bizInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"biz_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the object to be published.`,
			},
			"biz_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the type of the object to be published.`,
			},
		},
	}
}

func resourceArchitectureBatchPublishCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dataarts"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	groupId, err := batchPublishResource(client, d, true)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(groupId.(string))
	return nil
}

func batchPublishResource(client *golangsdk.ServiceClient, d *schema.ResourceData, isPublish bool) (interface{}, error) {
	httpUrl := "v2/{project_id}/design/approvals/batch-publish"
	publisPath := client.Endpoint + httpUrl
	publisPath = strings.ReplaceAll(publisPath, "{project_id}", client.ProjectID)

	publisOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"workspace":    d.Get("workspace_id").(string),
			"Content-Type": "application/json;charset=UTF-8",
		},
		JSONBody: buildActionArchitectureResourceBodyParams(d, isPublish),
	}

	resp, err := client.Request("POST", publisPath, &publisOpt)
	if err != nil {
		return nil, fmt.Errorf("error publishing DataArts Architecture resource: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	groupId := utils.PathSearch("data.value.group_id", respBody, nil)
	if groupId == nil {
		return nil, fmt.Errorf("error creating DataArts Architecture batch publishment: ID is not found in API response")
	}

	errMsg := getActionErrMsgs(respBody)
	if errMsg != "" {
		return nil, fmt.Errorf("error publishing some resources: %s", errMsg)
	}
	return groupId, nil
}

func buildActionArchitectureResourceBodyParams(d *schema.ResourceData, isPublish bool) map[string]interface{} {
	params := map[string]interface{}{
		"approver_user_id":   d.Get("approver_user_id"),
		"approver_user_name": d.Get("approver_user_name"),
		"biz_infos":          buildBusinessInfos(d.Get("biz_infos").([]interface{})),
		"fast_approval":      d.Get("fast_approval"),
	}

	if isPublish {
		params["schedule_time"] = utils.ValueIgnoreEmpty(d.Get("schedule_time"))
	}
	return params
}

func buildBusinessInfos(bizInfos []interface{}) []interface{} {
	result := make([]interface{}, len(bizInfos))
	for i, v := range bizInfos {
		result[i] = map[string]interface{}{
			"biz_id":   utils.PathSearch("biz_id", v, ""),
			"biz_type": utils.PathSearch("biz_type", v, ""),
		}
	}
	return result
}

func resourceArchitectureBatchPublishRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceArchitectureBatchPublishDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for publshing resources. Deleting this resource will not clear
	the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
