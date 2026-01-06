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

// @API Workspace POST /v1/{project_id}/persistent-storages
// @API Workspace GET /v1/{project_id}/persistent-storages
// @API Workspace DELETE /v1/{project_id}/persistent-storages/{storage_id}
func ResourceAppNasStorage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppNasStorageCreate,
		ReadContext:   resourceAppNasStorageRead,
		DeleteContext: resourceAppNasStorageDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAppNasStorageImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the NAS storage is located.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the NAS storage.",
			},
			"storage_metadata": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: `The metadata of the corresponding storage.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_handle": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The storage name.",
						},
						"storage_class": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The storage type.",
						},
						"export_location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The storage access URL.",
						},
					},
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildAppNasStorageMetadata(metadataConfigs []interface{}) map[string]interface{} {
	if len(metadataConfigs) < 1 {
		return nil
	}

	metadata := metadataConfigs[0]
	return map[string]interface{}{
		"storage_handle": utils.PathSearch("storage_handle", metadata, nil),
		"storage_class":  utils.PathSearch("storage_class", metadata, nil),
	}
}

func buildAppNasStorageCreateOpts(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":             d.Get("name").(string),
		"storage_metadata": buildAppNasStorageMetadata(d.Get("storage_metadata").([]interface{})),
	}
}

func resourceAppNasStorageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/persistent-storages"
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildAppNasStorageCreateOpts(d)),
	}
	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating NAS storage of Workspace APP: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	storageId := utils.PathSearch("id", respBody, "").(string)
	if storageId == "" {
		return diag.Errorf("unable to find the storage ID from the API response")
	}
	d.SetId(storageId)

	return resourceAppNasStorageRead(ctx, d, meta)
}

func flattenStorageMetadata(metadata interface{}) []map[string]interface{} {
	if metadata != nil {
		return []map[string]interface{}{
			{
				"storage_handle":  utils.PathSearch("storage_handle", metadata, nil),
				"storage_class":   utils.PathSearch("storage_class", metadata, nil),
				"export_location": utils.PathSearch("export_location", metadata, nil),
			},
		}
	}
	return nil
}

func GetAppNasStorageById(client *golangsdk.ServiceClient, storageId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/persistent-storages?storage_id={storage_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{storage_id}", storageId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	storageConfig := utils.PathSearch("items|[0]", respBody, nil)
	if storageConfig == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/persistent-storages",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the NAS storage (%s) has been removed from the Workspace APP service", storageId)),
			},
		}
	}
	return storageConfig, nil
}

func resourceAppNasStorageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		storageId = d.Id()
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	storageConfig, err := GetAppNasStorageById(client, storageId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "NAS Storage of Workspace APP")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", storageConfig, nil)),
		d.Set("storage_metadata", flattenStorageMetadata(utils.PathSearch("storage_metadata", storageConfig, nil))),
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("create_time",
			storageConfig, "").(string))/1000, false)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("unable to setting resource fields of the NAS storage: %s", err)
	}
	return nil
}

func resourceAppNasStorageDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v1/{project_id}/persistent-storages/{storage_id}"
		storageId = d.Id()
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{storage_id}", storageId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// Although the deletion result of the main region shows that the interface returns a 200 status code when
		// deleting a non-existent NAS storage, in order to avoid the possible return of a 404 status code in the
		// future, the CheckDeleted design is retained here.
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting NAS storage (%s)", storageId))
	}
	return nil
}

func listAppNasStorages(client *golangsdk.ServiceClient, queryPath ...string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/persistent-storages?limit=100"
		listOpt = golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		offset = 0
		result = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	if len(queryPath) > 0 {
		listPath += queryPath[0]
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, fmt.Errorf("error getting list of NAS storages: %s", err)
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

func resourceAppNasStorageImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
	return []*schema.ResourceData{d}, nil
}
