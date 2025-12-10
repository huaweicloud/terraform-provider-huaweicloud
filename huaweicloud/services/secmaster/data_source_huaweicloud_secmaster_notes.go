package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/notes/search
func DataSourceNotes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotesRead,

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
			"sort_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"from_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"to_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"war_room_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildNotesDataSchema(),
			},
		},
	}
}

func buildNotesDataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildNotesDataDataSchema(),
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"marked_note": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"note_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildNotesDataUserSchema(),
			},
			"war_room_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Field `content` is a JSON format string value.
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildNotesDataDataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildNotesDataUserSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildNotesBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"sort_by":     utils.ValueIgnoreEmpty(d.Get("sort_by")),
		"order":       utils.ValueIgnoreEmpty(d.Get("order")),
		"from_date":   utils.ValueIgnoreEmpty(d.Get("from_date")),
		"to_date":     utils.ValueIgnoreEmpty(d.Get("to_date")),
		"war_room_id": utils.ValueIgnoreEmpty(d.Get("war_room_id")),
		"limit":       100,
	}
}

func dataSourceNotesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/notes/search"
		offset      = 0
		allResult   = make([]interface{}, 0)
		requestBody = utils.RemoveNil(buildNotesBodyParams(d))
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestBody["offset"] = offset
		requestOpt.JSONBody = requestBody
		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster notes: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		data := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		if len(data) == 0 {
			break
		}

		allResult = append(allResult, data...)
		offset += len(data)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenNotes(allResult)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenNotes(allResult []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(allResult))
	for _, v := range allResult {
		rst = append(rst, map[string]interface{}{
			"create_time":  utils.PathSearch("create_time", v, nil),
			"update_time":  utils.PathSearch("update_time", v, nil),
			"data":         flattenNotesData(utils.PathSearch("data", v, nil)),
			"id":           utils.PathSearch("id", v, nil),
			"is_deleted":   utils.PathSearch("is_deleted", v, nil),
			"marked_note":  utils.PathSearch("marked_note", v, nil),
			"note_type":    utils.PathSearch("note_type", v, nil),
			"project_id":   utils.PathSearch("project_id", v, nil),
			"type":         utils.PathSearch("type", v, nil),
			"user":         flattenNotesUser(utils.PathSearch("user", v, nil)),
			"war_room_id":  utils.PathSearch("war_room_id", v, nil),
			"workspace_id": utils.PathSearch("workspace_id", v, nil),
			"content":      utils.JsonToString(utils.PathSearch("content", v, nil)),
		})
	}
	return rst
}

func flattenNotesData(respBody interface{}) []interface{} {
	rstMap := map[string]interface{}{
		"content": utils.PathSearch("content", respBody, nil),
	}

	return []interface{}{rstMap}
}

func flattenNotesUser(respBody interface{}) []interface{} {
	rstMap := map[string]interface{}{
		"id":   utils.PathSearch("id", respBody, nil),
		"name": utils.PathSearch("name", respBody, nil),
	}

	return []interface{}{rstMap}
}
