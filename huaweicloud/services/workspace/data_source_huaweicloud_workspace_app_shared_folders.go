package workspace

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v1/{project_id}/persistent-storages/actions/list-share-folders
func DataSourceAppSharedFolders() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppSharedFoldersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the WKS storage is located.`,
			},

			// Required parameters.
			"storage_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The WKS storage ID.`,
			},

			// Optional parameters.
			"storage_claim_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The WKS storage directory claim ID.`,
			},
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The shared folder path for query.`,
			},

			// Attributes.
			"shared_folders": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        appSharedFoldersSchema(),
				Description: `The list of the shared folders that matched filter parameters.`,
			},
		},
	}
}

func appSharedFoldersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"storage_claim_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The WKS storage directory claim ID.`,
			},
			"folder_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The storage object path.`,
			},
			"delimiter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The path delimiter.`,
			},
			"claim_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The storage claim type.`,
			},
			"count": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: `The number of associated users and user groups of the shared folder.`,
			},
		},
	}
}

func buildSharedFoldersQueryParams(d *schema.ResourceData) string {
	res := ""

	res = fmt.Sprintf("%s&storage_id=%v", res, d.Get("storage_id"))

	if v, ok := d.GetOk("storage_claim_id"); ok {
		res = fmt.Sprintf("%s&storage_claim_id=%v", res, v)
	}
	if v, ok := d.GetOk("path"); ok {
		res = fmt.Sprintf("%s&path=%v", res, v)
	}

	return res
}

func listSharedFolders(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/persistent-storages/actions/list-share-folders?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	listPathWithLimit += buildSharedFoldersQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPathWithLimit + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		items := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, items...)
		if len(items) < limit {
			break
		}

		offset += len(items)
	}

	return result, nil
}

func flattenSharedFolders(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"storage_claim_id": utils.PathSearch("storage_claim_id", item, nil),
			"folder_path":      utils.PathSearch("folder_path", item, nil),
			"delimiter":        utils.PathSearch("delimiter", item, nil),
			"claim_mode":       utils.PathSearch("claim_mode", item, nil),
			"count":            utils.PathSearch("count", item, make(map[string]interface{})),
		})
	}

	return result
}

func dataSourceAppSharedFoldersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	resp, err := listSharedFolders(client, d)
	if err != nil {
		return diag.Errorf("error querying shared folders: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("shared_folders", flattenSharedFolders(resp)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
