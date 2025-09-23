package workspace

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace POST /v1/{project_id}/persistent-storages/{storage_id}/actions/assign-folder
// @API Workspace GET /v1/{project_id}/persistent-storages/actions/list-attachments
// @API Workspace POST /v1/{project_id}/persistent-storages/{storage_id}/actions/delete-user-attachment
func ResourceAppPersonalFolders() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppPersonalFoldersCreate,
		ReadContext:   resourceAppPersonalFoldersRead,
		DeleteContext: resourceAppPersonalFoldersDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAppPersonalFolderImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the personal folders are located.",
			},
			"storage_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The NAS storage ID to which the personal folders belong.",
			},
			"assignments": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_statement_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The ID of the storage permission policy.",
						},
						"attach": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The object name of personal folder assignment.",
						},
						"attach_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "The type of personal folder assignment.",
						},
					},
				},
				Description: "The assignment configuration of personal folders.",
			},
		},
	}
}

func buildAppPersonalFoldersCreateJsonBody(d *schema.ResourceData) map[string]interface{} {
	assignments := d.Get("assignments").(*schema.Set)
	items := make([]interface{}, 0, assignments.Len())

	for _, assignment := range assignments.List() {
		items = append(items, map[string]interface{}{
			"policy_statement_id": utils.PathSearch("policy_statement_id", assignment, nil),
			"attach":              utils.PathSearch("attach", assignment, nil),
			"attach_type":         utils.ValueIgnoreEmpty(utils.PathSearch("attach_type", assignment, nil)),
		})
	}

	return map[string]interface{}{
		"items": items,
	}
}

func resourceAppPersonalFoldersCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v1/{project_id}/persistent-storages/{storage_id}/actions/assign-folder"
		storageId = d.Get("storage_id").(string)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{storage_id}", storageId)
	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildAppPersonalFoldersCreateJsonBody(d)),
	}
	_, err = client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error assigning personal folders of Workspace APP: %s", err)
	}

	d.SetId(storageId)

	return resourceAppPersonalFoldersRead(ctx, d, meta)
}

func ListAppPersonalFolders(client *golangsdk.ServiceClient, storageId string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/persistent-storages/actions/list-attachments?claim_mode=USER&storage_id={storage_id}&limit=100"
		listOpt = golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		offset = 0
		result = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{storage_id}", storageId)

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, fmt.Errorf("error getting list of personal folders: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		items := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		if len(items) < 1 {
			break
		}
		result = append(result, items...)
		offset += len(items)
	}
	return result, nil
}

func flattenAppPersonalFolderAssignments(assignments []interface{}) []interface{} {
	result := make([]interface{}, 0, len(assignments))

	for _, assignment := range assignments {
		result = append(result, map[string]interface{}{
			"policy_statement_id": utils.PathSearch("policy_statement.policy_statement_id", assignment, nil),
			"attach":              utils.PathSearch("attachment.attach", assignment, nil),
			"attach_type":         utils.PathSearch("attachment.attach_type", assignment, nil),
		})
	}

	return result
}

func resourceAppPersonalFoldersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		storageId = d.Id()
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	assignments, err := ListAppPersonalFolders(client, storageId)
	if err != nil {
		// If NAS storage not exist or the service has been closed, the API request will fail and returns a 404 error.
		return common.CheckDeletedDiag(d, err, "Personal folders of Workspace APP")
	}
	if len(assignments) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte("All personal folders have been removed from the Workspace APP service"),
			},
		}, "")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("assignments", flattenAppPersonalFolderAssignments(assignments)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("unable to setting resource fields of the personal folders: %s", err)
	}
	return nil
}

func buildAppPersonalFoldersDeleteJsonBody(client *golangsdk.ServiceClient, storageId string) (map[string]interface{}, error) {
	// The assignments is empey if resource is tainted.
	assignments, err := ListAppPersonalFolders(client, storageId)
	if err != nil {
		return nil, err
	}

	items := utils.PathSearch("items[*].attachment.attach", assignments, make([]interface{}, 0)).([]interface{})
	if len(items) < 1 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte("all personal folders removed from the NAS storage for Workspace APP"),
			},
		}
	}
	return map[string]interface{}{
		"items": items,
	}, nil
}

func resourceAppPersonalFoldersDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v1/{project_id}/persistent-storages/{storage_id}/actions/delete-user-attachment"
		storageId = d.Get("storage_id").(string)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	deleteItems, err := buildAppPersonalFoldersDeleteJsonBody(client, storageId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Personal folders query before delete")
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{storage_id}", storageId)
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(deleteItems),
	}
	_, err = client.Request("POST", deletePath, &deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting personal folders of Workspace APP: %s", err)
	}
	return nil
}

func resourceAppPersonalFolderImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var (
		storageId string
		importId  = d.Id()
		re        = regexp.MustCompile(`^\d{18}$`)
	)

	if re.MatchString(d.Id()) {
		// The imported ID is NAS storage ID.
		storageId = importId
	} else {
		// The imported ID is NAS storage name or other meaningless characters, which are all queried as names.
		cfg := meta.(*config.Config)
		region := cfg.GetRegion(d)

		client, err := cfg.NewServiceClient("appstream", region)
		if err != nil {
			return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
		}

		storages, err := listAppNasStorages(client)
		if err != nil {
			return nil, err
		}
		storageId = utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0].id", importId), storages, "").(string)
		if storageId == "" {
			return nil, fmt.Errorf("unable to find the NAS storage using its name (%s): %s", importId, err)
		}
	}

	d.SetId(storageId)
	return []*schema.ResourceData{d}, d.Set("storage_id", storageId)
}
