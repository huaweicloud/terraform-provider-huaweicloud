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
func ResourceArchitectureBatchPublishment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArchitectureBatchPublishmentCreate,
		ReadContext:   resourceArchitectureBatchPublishmentRead,
		DeleteContext: resourceArchitectureBatchPublishmentDelete,

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
				Description: "Specifies the list of the business information.",
				Elem: &schema.Resource{
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
				},
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
				Description: "Whether to automatically review.",
			},
		},
	}
}

func resourceArchitectureBatchPublishmentCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/design/approvals/batch-publish"
		product = "dataarts"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"workspace":    d.Get("workspace_id").(string),
			"Content-Type": "application/json;charset=UTF-8",
		},
		JSONBody: utils.RemoveNil(buildCreateArchitectureBatchPublishmentBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error publishing DataArts Architecture resource: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	groupId := utils.PathSearch("data.value.group_id", createRespBody, nil)
	if groupId == nil {
		return diag.Errorf("error creating DataArts Architecture batch publishment: ID is not found in API response")
	}

	publishStatus := utils.PathSearch("data.value.operation_status", createRespBody, nil)
	errMsgs := make([]string, 0)
	if publishStatus == "FAILED" {
		for _, v := range utils.PathSearch("data.value.groups", createRespBody, make([]interface{}, 0)).([]interface{}) {
			if utils.PathSearch("operation_status", v, nil) != "FAILED" {
				continue
			}
			bizId := utils.PathSearch("biz_id", v, "").(string)
			failedMsg := utils.PathSearch("remark", v, "").(string)
			errMsgs = append(errMsgs, fmt.Sprintf("%s | %s;", bizId, failedMsg))
		}
	}

	if len(errMsgs) > 0 {
		return diag.Errorf("error publishing some resources: %s", strings.Join(errMsgs, "\n"))
	}

	d.SetId(groupId.(string))
	return nil
}

func buildCreateArchitectureBatchPublishmentBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"approver_user_id":   d.Get("approver_user_id"),
		"approver_user_name": d.Get("approver_user_name"),
		"biz_infos":          buildBusinessInfos(d.Get("biz_infos").([]interface{})),
		"fast_approval":      utils.ValueIgnoreEmpty(d.Get("fast_approval")),
	}
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

func resourceArchitectureBatchPublishmentRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceArchitectureBatchPublishmentDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for publshing resource. Deleting this resource will
	not change the status of the currently published resource, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
