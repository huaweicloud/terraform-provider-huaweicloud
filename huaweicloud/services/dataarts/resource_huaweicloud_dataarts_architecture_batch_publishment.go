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
// @API DataArtsStudio POST /v2/{project_id}/design/approvals/batch-offline
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
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
				Description: utils.SchemaDesc(
					"Specifies whether to automatically review.",
					utils.SchemaDescInput{
						Internal: true,
					},
				),
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

func resourceArchitectureBatchPublishmentCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dataarts", cfg.GetRegion(d))
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

func resourceArchitectureBatchPublishmentRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceArchitectureBatchPublishmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dataarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}
	if err = batchOfflineResource(client, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func batchOfflineResource(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v2/{project_id}/design/approvals/batch-offline"
	offlinePath := client.Endpoint + httpUrl
	offlinePath = strings.ReplaceAll(offlinePath, "{project_id}", client.ProjectID)

	offlineOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"workspace":    d.Get("workspace_id").(string),
			"Content-Type": "application/json;charset=UTF-8",
		},
		JSONBody: buildActionArchitectureResourceBodyParams(d, false),
	}

	resp, err := client.Request("POST", offlinePath, &offlineOpt)
	if err != nil {
		return fmt.Errorf("error offlining DataArts Architecture resource: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	errMsg := getActionErrMsgs(respBody)
	if errMsg != "" {
		return fmt.Errorf("error offlining some resources: %s", errMsg)
	}
	return nil
}

func getActionErrMsgs(respBody interface{}) string {
	status := utils.PathSearch("data.value.operation_status", respBody, "").(string)
	if status != "FAILED" {
		return ""
	}

	errMsgs := make([]string, 0)
	for _, v := range utils.PathSearch("data.value.groups", respBody, make([]interface{}, 0)).([]interface{}) {
		if utils.PathSearch("operation_status", v, nil) != "FAILED" {
			continue
		}
		bizId := utils.PathSearch("biz_id", v, "").(string)
		failedMsg := utils.PathSearch("remark", v, "").(string)
		errMsgs = append(errMsgs, fmt.Sprintf("%s | %s;", bizId, failedMsg))
	}

	if len(errMsgs) > 0 {
		return strings.Join(errMsgs, "\n")
	}
	return ""
}
