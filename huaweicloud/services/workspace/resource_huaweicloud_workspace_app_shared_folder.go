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

// @API Workspace POST /v1/{project_id}/persistent-storages/{storage_id}/actions/create-share-folder
// @API Workspace GET /v1/{project_id}/persistent-storages/actions/list-share-folders
// @API Workspace POST /v1/{project_id}/persistent-storages/{storage_id}/actions/delete-storage-claim
func ResourceAppSharedFolder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppSharedFolderCreate,
		ReadContext:   resourceAppSharedFolderRead,
		DeleteContext: resourceAppSharedFolderDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAppSharedFolderImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the shared folder is located.",
			},
			"storage_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The NAS storage ID to which the shared folder belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the shared folder.",
			},
			"delimiter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The delimiter that the shared folder path used.",
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The path of the shared folder.",
			},
		},
	}
}

func buildAppSharedFolderCreateJsonBody(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"folder_name": d.Get("name").(string),
	}
}

func resourceAppSharedFolderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v1/{project_id}/persistent-storages/{storage_id}/actions/create-share-folder"
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
		JSONBody:         utils.RemoveNil(buildAppSharedFolderCreateJsonBody(d)),
	}
	requestResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error creating shared folder of Workspace APP: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	sharedFolderId := utils.PathSearch("storage_claim_id", respBody, "").(string)
	if sharedFolderId == "" {
		return diag.Errorf("unable to find the shared folder ID from the API response")
	}
	d.SetId(sharedFolderId)

	return resourceAppSharedFolderRead(ctx, d, meta)
}

func listAppSharedFolders(client *golangsdk.ServiceClient, storageId string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/persistent-storages/actions/list-share-folders?storage_id={storage_id}&limit=100"
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
			return nil, fmt.Errorf("error getting list of shared folders: %s", err)
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

func GetAppSharedFolderById(client *golangsdk.ServiceClient, storageId, sharedFolderId string) (interface{}, error) {
	sharedFolders, err := listAppSharedFolders(client, storageId)
	if err != nil {
		return nil, err
	}

	sharedFolder := utils.PathSearch(fmt.Sprintf("[?storage_claim_id=='%s']|[0]", sharedFolderId), sharedFolders, nil)
	if sharedFolder == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte("the shared folder has been removed from the Workspace APP service"),
			},
		}
	}
	return sharedFolder, nil
}

func resourceAppSharedFolderRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		storageId      = d.Get("storage_id").(string)
		sharedFolderId = d.Id()
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	sharedFolder, err := GetAppSharedFolderById(client, storageId, sharedFolderId)
	if err != nil {
		// If NAS storage not exist or the service has been closed, the API request will fail and returns a 404 error.
		return common.CheckDeletedDiag(d, err, "Shared folder of Workspace APP")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("delimiter", utils.PathSearch("delimiter", sharedFolder, nil)),
		d.Set("path", utils.PathSearch("folder_path", sharedFolder, nil)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("unable to setting resource fields of the shared folder: %s", err)
	}
	return nil
}

func resourceAppSharedFolderDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		httpUrl        = "v1/{project_id}/persistent-storages/{storage_id}/actions/delete-storage-claim"
		storageId      = d.Get("storage_id").(string)
		sharedFolderId = d.Id()
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{storage_id}", storageId)
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"items": []interface{}{sharedFolderId},
		},
	}
	_, err = client.Request("POST", deletePath, &deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting shared folder (%s) of Workspace APP", sharedFolderId))
	}
	return nil
}

func retrieveAppNasStorageId(client *golangsdk.ServiceClient, importId string) (string, error) {
	re := regexp.MustCompile(`^\d{18}$`)
	if re.MatchString(importId) {
		// The imported ID is NAS storage ID.
		return importId, nil
	}

	// The imported ID is NAS storage name or other meaningless characters, which are all queried as names.
	storages, err := listAppNasStorages(client)
	if err != nil {
		return "", err
	}
	storageId := utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0].id", importId), storages, "").(string)
	if storageId == "" {
		return "", fmt.Errorf("unable to find the NAS storage ID using its name (%s)", storageId)
	}
	return storageId, nil
}

func retrieveAppSharedFolderId(client *golangsdk.ServiceClient, storageId, importId string) (string, error) {
	re := regexp.MustCompile(`^\d{18}$`)
	if re.MatchString(importId) {
		// The imported ID is shared folder ID.
		return importId, nil
	}

	// The imported ID is shared folder name.
	sharedFolders, err := listAppSharedFolders(client, storageId)
	if err != nil {
		return "", err
	}

	delimiter := utils.PathSearch("[0].delimiter", sharedFolders, "").(string)
	sharedFolderId := utils.PathSearch(fmt.Sprintf("[?contains(folder_path, '%s%s%s')]|[0].storage_claim_id",
		delimiter, importId, delimiter), sharedFolders, "").(string)
	if sharedFolderId == "" {
		return "", fmt.Errorf("unable to find the shared folder ID using its name (%s)", sharedFolderId)
	}
	return sharedFolderId, nil
}

func resourceAppSharedFolderImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var (
		storageId, sharedFolderId string
		importId                  = d.Id()
		cfg                       = meta.(*config.Config)
		region                    = cfg.GetRegion(d)
	)

	parts := strings.Split(importId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<storage_name>/<name>' or '<storage_id>/<id>', but got '%s'", importId)
	}

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
	}

	storageId, err = retrieveAppNasStorageId(client, parts[0])
	if err != nil {
		return nil, err
	}

	sharedFolderId, err = retrieveAppSharedFolderId(client, storageId, parts[1])
	if err != nil {
		return nil, err
	}

	d.SetId(sharedFolderId)
	return []*schema.ResourceData{d}, d.Set("storage_id", storageId)
}
