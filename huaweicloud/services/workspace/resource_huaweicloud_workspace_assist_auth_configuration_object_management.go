package workspace

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var assistAuthConfigurationObjectManagementParamKeys = []string{"objects"}

// @API Workspace POST /v2/{project_id}/assist-auth-config/apply-objects
// @API Workspace GET /v2/{project_id}/assist-auth-config/apply-objects
func ResourceAssistAuthConfigurationObjectManagement() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssistAuthConfigurationObjectManagementCreate,
		ReadContext:   resourceAssistAuthConfigurationObjectManagementRead,
		UpdateContext: resourceAssistAuthConfigurationObjectManagementUpdate,
		DeleteContext: resourceAssistAuthConfigurationObjectManagementDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the apply objects of assist auth configuration are located.`,
			},

			// Required parameters.
			"objects": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the binding object.`,
						},
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The ID of the binding object.`,
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the binding object.`,
						},
					},
				},
				DiffSuppressFunc: utils.SuppressObjectSliceDiffs(),
				Description:      `The list of objects to be managed with the assist auth configuration.`,
			},

			// Internal parameters.
			"objects_origin": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The type of the binding object.`,
						},
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The ID of the binding object.`,
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The name of the binding object.`,
						},
					},
				},
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
the new value next time the change is made. The corresponding parameter name is 'objects'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildAssistAuthConfigurationAppliedObjects(objects []interface{}) []map[string]interface{} {
	if len(objects) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(objects))
	for _, obj := range objects {
		result = append(result, map[string]interface{}{
			"object_type": utils.PathSearch("type", obj, nil),
			"object_id":   utils.PathSearch("id", obj, nil),
			"object_name": utils.PathSearch("name", obj, nil),
		})
	}

	return result
}

func addAssistAuthConfigurationAppliedObjects(client *golangsdk.ServiceClient, objects []interface{}) error {
	httpUrl := "v2/{project_id}/assist-auth-config/apply-objects"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"add": buildAssistAuthConfigurationAppliedObjects(objects),
		},
		OkCodes: []int{204},
	}

	_, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return err
	}

	return nil
}

func deleteAssistAuthConfigurationAppliedObjects(client *golangsdk.ServiceClient, objects []interface{}) error {
	httpUrl := "v2/{project_id}/assist-auth-config/apply-objects"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"delete": buildAssistAuthConfigurationAppliedObjects(objects),
		},
		OkCodes: []int{204},
	}

	_, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return err
	}

	return nil
}

func resourceAssistAuthConfigurationObjectManagementCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	err = addAssistAuthConfigurationAppliedObjects(client, d.Get("objects").([]interface{}))
	if err != nil {
		return diag.Errorf("error adding assist auth configuration applied objects: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshSliceParamOriginValues(d, assistAuthConfigurationObjectManagementParamKeys)
	if err != nil {
		// Don't fail the creation if origin refresh fails
		log.Printf("[WARN] Unable to refresh the origin values: %s", err)
	}

	return resourceAssistAuthConfigurationObjectManagementRead(ctx, d, meta)
}

func orderAssistAuthConfigurationAppliedObjects(appliedObjects, originObjects []interface{}) []interface{} {
	if len(originObjects) < 1 || len(appliedObjects) < 1 {
		return appliedObjects
	}

	sortedObjects := make([]interface{}, 0, len(appliedObjects))
	objectsCopy := appliedObjects
	for _, originObject := range originObjects {
		for index, appliedObj := range objectsCopy {
			if utils.PathSearch("type", originObject, nil) != utils.PathSearch("object_type", appliedObj, nil) ||
				utils.PathSearch("id", originObject, nil) != utils.PathSearch("object_id", appliedObj, nil) ||
				utils.PathSearch("name", originObject, nil) != utils.PathSearch("object_name", appliedObj, nil) {
				continue
			}
			// Add the found applied object to the sorted objects list.
			sortedObjects = append(sortedObjects, objectsCopy[index])
			// Remove the processed applied object from the original array.
			objectsCopy = append(objectsCopy[:index], objectsCopy[index+1:]...)
		}
	}

	// Add any remaining unsorted applied objects to the end of the sorted list.
	sortedObjects = append(sortedObjects, objectsCopy...)
	return sortedObjects
}

func filterAssistAuthConfigurationAppliedObjects(appliedObjects, originObjects []interface{}) []interface{} {
	filteredObjects := make([]interface{}, 0, len(appliedObjects))

	for _, originObj := range originObjects {
		for _, obj := range appliedObjects {
			if utils.PathSearch("object_type", obj, nil) == utils.PathSearch("type", originObj, nil) &&
				utils.PathSearch("object_id", obj, nil) == utils.PathSearch("id", originObj, nil) &&
				utils.PathSearch("object_name", obj, nil) == utils.PathSearch("name", originObj, nil) {
				filteredObjects = append(filteredObjects, originObj)
				break
			}
		}
	}

	return filteredObjects
}

func ListAndReorderAssistAuthConfigurationAppliedObjects(client *golangsdk.ServiceClient, originObjects []interface{}) ([]interface{}, error) {
	appliedObjects, err := listAssistAuthConfigurationAppliedObjects(client)
	if err != nil {
		return nil, err
	}
	if len(appliedObjects) < 1 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/assist-auth-config/apply-objects",
				RequestId: "NONE",
				Body:      []byte(`all applied objects have been removed`),
			},
		}
	}

	if len(originObjects) > 0 {
		filteredObjects := filterAssistAuthConfigurationAppliedObjects(appliedObjects, originObjects)
		if len(filteredObjects) < 1 {
			return nil, golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Method:    "GET",
					URL:       "/v2/{project_id}/assist-auth-config/apply-objects",
					RequestId: "NONE",
					Body:      []byte(`all applied objects (managed by provider) have been removed`),
				},
			}
		}

		// Do not reorder the filtered objects, just reorder the objects of the remote console.
		return orderAssistAuthConfigurationAppliedObjects(appliedObjects, originObjects), nil
	}
	return appliedObjects, nil
}

func flattenAssistAuthConfigurationAppliedObjects(objects []interface{}) []interface{} {
	result := make([]interface{}, 0, len(objects))

	for _, obj := range objects {
		result = append(result, map[string]interface{}{
			"type": utils.PathSearch("object_type", obj, nil),
			"id":   utils.PathSearch("object_id", obj, nil),
			"name": utils.PathSearch("object_name", obj, nil),
		})
	}

	return result
}

func resourceAssistAuthConfigurationObjectManagementRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	objects, err := ListAndReorderAssistAuthConfigurationAppliedObjects(client, d.Get("objects_origin").([]interface{}))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error listing assist auth configuration applied objects")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("objects", flattenAssistAuthConfigurationAppliedObjects(objects)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getUpdateAssistAuthConfigurationAppliedObjects(d *schema.ResourceData) (newObjects, rmObjects []interface{}) {
	var (
		consoleObjects, scriptObjects = d.GetChange("objects")
		consoleObjectsList            = consoleObjects.([]interface{})
		scriptObjectsList             = scriptObjects.([]interface{})
		originObjectsList             = d.Get("objects_origin").([]interface{})
		isObjMatched                  bool
	)

	newObjects = make([]interface{}, 0, len(scriptObjectsList))
	rmObjects = make([]interface{}, 0, len(originObjectsList))

	for _, obj := range scriptObjectsList {
		isObjMatched = false
		for _, consoleObj := range consoleObjectsList {
			if utils.PathSearch("type", obj, nil) == utils.PathSearch("type", consoleObj, nil) &&
				utils.PathSearch("id", obj, nil) == utils.PathSearch("id", consoleObj, nil) &&
				utils.PathSearch("name", obj, nil) == utils.PathSearch("name", consoleObj, nil) {
				isObjMatched = true
				break
			}
		}
		if !isObjMatched {
			newObjects = append(newObjects, obj)
		}
	}

	for _, obj := range originObjectsList {
		isObjMatched = false
		for _, scriptObj := range scriptObjectsList {
			if utils.PathSearch("type", obj, nil) == utils.PathSearch("type", scriptObj, nil) &&
				utils.PathSearch("id", obj, nil) == utils.PathSearch("id", scriptObj, nil) &&
				utils.PathSearch("name", obj, nil) == utils.PathSearch("name", scriptObj, nil) {
				isObjMatched = true
				break
			}
		}
		if !isObjMatched {
			rmObjects = append(rmObjects, obj)
		}
	}

	return
}

func resourceAssistAuthConfigurationObjectManagementUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	newObjects, rmObjects := getUpdateAssistAuthConfigurationAppliedObjects(d)
	if len(newObjects) > 0 {
		err = addAssistAuthConfigurationAppliedObjects(client, newObjects)
		if err != nil {
			return diag.Errorf("error adding assist auth configuration applied objects: %s", err)
		}
	}

	if len(rmObjects) > 0 {
		err = deleteAssistAuthConfigurationAppliedObjects(client, rmObjects)
		if err != nil {
			return diag.Errorf("error deleting assist auth configuration applied objects: %s", err)
		}
	}

	// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshSliceParamOriginValues(d, assistAuthConfigurationObjectManagementParamKeys)
	if err != nil {
		// Don't fail the update if origin refresh fails
		log.Printf("[WARN] Unable to refresh the origin values: %s", err)
	}
	return resourceAssistAuthConfigurationObjectManagementRead(ctx, d, meta)
}

func resourceAssistAuthConfigurationObjectManagementDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	err = deleteAssistAuthConfigurationAppliedObjects(client, d.Get("objects").([]interface{}))
	if err != nil {
		return diag.Errorf("error deleting assist auth configuration applied objects: %s", err)
	}

	return nil
}
