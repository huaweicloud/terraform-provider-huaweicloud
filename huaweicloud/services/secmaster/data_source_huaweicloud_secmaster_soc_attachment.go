package secmaster

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/attachment/{attach_id}
func DataSourceSocAttachment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSocAttachmentRead,

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
			"attach_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_folder": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modifier_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modifier_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_deleted": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSocAttachmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/attachment/{attach_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{attach_id}", d.Get("attach_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster soc attachment: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataResp := utils.PathSearch("data", respBody, nil)

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenSocAttachmentData(dataResp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSocAttachmentData(dataResp interface{}) []interface{} {
	if dataResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":            utils.PathSearch("id", dataResp, nil),
			"file_name":     utils.PathSearch("file_name", dataResp, nil),
			"file_folder":   utils.PathSearch("file_folder", dataResp, nil),
			"workspace_id":  utils.PathSearch("workspace_id", dataResp, nil),
			"creator_id":    utils.PathSearch("creator_id", dataResp, nil),
			"creator_name":  utils.PathSearch("creator_name", dataResp, nil),
			"create_time":   utils.PathSearch("create_time", dataResp, nil),
			"update_time":   utils.PathSearch("update_time", dataResp, nil),
			"modifier_id":   utils.PathSearch("modifier_id", dataResp, nil),
			"modifier_name": utils.PathSearch("modifier_name", dataResp, nil),
			"is_deleted":    utils.PathSearch("is_deleted", dataResp, nil),
			"storage_type":  utils.PathSearch("storage_type", dataResp, nil),
			"storage_url":   utils.PathSearch("storage_url", dataResp, nil),
		},
	}
}
