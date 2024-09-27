package secmaster

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	DataObjectRelationNotFound = "SecMaster.20030005"
)

var nonUpdatableParamsDataObjectRelations = []string{"workspace_id", "data_class", "data_object_id", "related_data_class"}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/{dataclass_type}/{data_object_id}/{related_dataclass_type}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/{dataclass_type}/{data_object_id}/{related_dataclass_type}
// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/{dataclass_type}/{data_object_id}/{related_dataclass_type}/search
func ResourceDataObjectRelations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataObjectRelationsCreate,
		ReadContext:   resourceDataObjectRelationsRead,
		UpdateContext: resourceDataObjectRelationsUpdate,
		DeleteContext: resourceDataObjectRelationsDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDataObjectRelationsImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsDataObjectRelations),

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
				Description: `Specifies the ID of the workspace to which the data object belongs.`,
			},
			"data_class": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the data class to which the data object belongs.`,
			},
			"data_object_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the data object.`,
			},
			"related_data_class": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the data class to which the related data object belongs.`,
			},
			"related_data_object_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the IDs of the data object.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceDataObjectRelationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	relatedDataObjectIds := utils.ExpandToStringList(d.Get("related_data_object_ids").(*schema.Set).List())

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	// createDataObjectRelations: Create the SecMaster data object relations.
	err = bindDataObjectRelations(d, client, relatedDataObjectIds)
	if err != nil {
		return diag.Errorf("error creating data object relations: %s", err)
	}

	workspaceId := d.Get("workspace_id").(string)
	dataClass := d.Get("data_class").(string)
	dataObjectId := d.Get("data_object_id").(string)
	relatedDataClass := d.Get("related_data_class").(string)

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", workspaceId, dataClass, dataObjectId, relatedDataClass))

	return resourceDataObjectRelationsRead(ctx, d, meta)
}

func resourceDataObjectRelationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	// getDataObjectRelations: Query the SecMaster data object relations detail
	relatedDataObjectIds, err := getDataObjectRelations(d, client)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving data object relations")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("related_data_object_ids", relatedDataObjectIds),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDataObjectRelationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	oldRaws, newRaws := d.GetChange("related_data_object_ids")
	unbindIds := utils.ExpandToStringListBySet(oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set)))
	bindIds := utils.ExpandToStringListBySet(newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set)))
	if len(unbindIds) > 0 {
		err := unbindDataObjectRelations(d, client, unbindIds)
		if err != nil {
			return diag.Errorf("error updating unbinding data object relations: %s", err)
		}
	}
	if len(bindIds) > 0 {
		err := bindDataObjectRelations(d, client, bindIds)
		if err != nil {
			return diag.Errorf("error updating binding data object relations: %s", err)
		}
	}

	return resourceDataObjectRelationsRead(ctx, d, meta)
}

func resourceDataObjectRelationsDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	resp, err := getDataObjectRelations(d, client)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving data object relations")
	}

	// deleteDataObjectRelations: delete the SecMaster data object relations.
	err = unbindDataObjectRelations(d, client, utils.ExpandToStringList(resp))
	if err != nil {
		return diag.Errorf("error deleting data object relations: %s", err)
	}

	return nil
}

func bindDataObjectRelations(d *schema.ResourceData, client *golangsdk.ServiceClient, relatedDataObjectIds []string) error {
	bindDataObjectRelationsHttpUrl :=
		"v1/{project_id}/workspaces/{workspace_id}/soc/{dataclass_type}/{data_object_id}/{related_dataclass_type}"
	bindDataObjectRelationsPath := client.Endpoint + bindDataObjectRelationsHttpUrl
	bindDataObjectRelationsPath = strings.ReplaceAll(bindDataObjectRelationsPath, "{project_id}", client.ProjectID)
	bindDataObjectRelationsPath = strings.ReplaceAll(bindDataObjectRelationsPath,
		"{workspace_id}", d.Get("workspace_id").(string))
	bindDataObjectRelationsPath = strings.ReplaceAll(bindDataObjectRelationsPath,
		"{dataclass_type}", d.Get("data_class").(string))
	bindDataObjectRelationsPath = strings.ReplaceAll(bindDataObjectRelationsPath,
		"{data_object_id}", d.Get("data_object_id").(string))
	bindDataObjectRelationsPath = strings.ReplaceAll(bindDataObjectRelationsPath,
		"{related_dataclass_type}", d.Get("related_data_class").(string))

	bindDataObjectRelationsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	bindDataObjectRelationsOpt.JSONBody = map[string]interface{}{
		"ids": relatedDataObjectIds,
	}
	_, err := client.Request("POST", bindDataObjectRelationsPath, &bindDataObjectRelationsOpt)
	if err != nil {
		return err
	}

	return nil
}

func unbindDataObjectRelations(d *schema.ResourceData, client *golangsdk.ServiceClient, unbindIds []string) error {
	unbindDataObjectRelationsHttpUrl := "v1/{project_id}/workspaces/{workspace_id}/soc/{dataclass_type}/{data_object_id}/{related_dataclass_type}"
	unbindDataObjectRelationsPath := client.Endpoint + unbindDataObjectRelationsHttpUrl
	unbindDataObjectRelationsPath = strings.ReplaceAll(unbindDataObjectRelationsPath, "{project_id}", client.ProjectID)
	unbindDataObjectRelationsPath = strings.ReplaceAll(unbindDataObjectRelationsPath,
		"{workspace_id}", d.Get("workspace_id").(string))
	unbindDataObjectRelationsPath = strings.ReplaceAll(unbindDataObjectRelationsPath,
		"{dataclass_type}", d.Get("data_class").(string))
	unbindDataObjectRelationsPath = strings.ReplaceAll(unbindDataObjectRelationsPath,
		"{data_object_id}", d.Get("data_object_id").(string))
	unbindDataObjectRelationsPath = strings.ReplaceAll(unbindDataObjectRelationsPath,
		"{related_dataclass_type}", d.Get("related_data_class").(string))

	unbindDataObjectRelationsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	unbindDataObjectRelationsOpt.JSONBody = map[string]interface{}{
		"ids": unbindIds,
	}
	_, err := client.Request("DELETE", unbindDataObjectRelationsPath, &unbindDataObjectRelationsOpt)
	if err != nil {
		return err
	}

	return nil
}

func getDataObjectRelations(d *schema.ResourceData, client *golangsdk.ServiceClient) ([]interface{}, error) {
	getDataObjectRelationsHttpUrl := "v1/{project_id}/workspaces/{workspace_id}/soc/{dataclass_type}/{data_object_id}/{related_dataclass_type}/search"
	getDataObjectRelationsPath := client.Endpoint + getDataObjectRelationsHttpUrl
	getDataObjectRelationsPath = strings.ReplaceAll(getDataObjectRelationsPath, "{project_id}", client.ProjectID)
	getDataObjectRelationsPath = strings.ReplaceAll(getDataObjectRelationsPath,
		"{workspace_id}", d.Get("workspace_id").(string))
	getDataObjectRelationsPath = strings.ReplaceAll(getDataObjectRelationsPath,
		"{dataclass_type}", d.Get("data_class").(string))
	getDataObjectRelationsPath = strings.ReplaceAll(getDataObjectRelationsPath,
		"{data_object_id}", d.Get("data_object_id").(string))
	getDataObjectRelationsPath = strings.ReplaceAll(getDataObjectRelationsPath,
		"{related_dataclass_type}", d.Get("related_data_class").(string))

	getDataObjectRelationsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getDataObjectRelationsOpt.JSONBody = map[string]interface{}{
		"limit":  1000,
		"offset": 0,
	}

	getDataObjectRelationsResp, err := client.Request("POST", getDataObjectRelationsPath, &getDataObjectRelationsOpt)
	if err != nil {
		// SecMaster.20010001: workspace ID not found
		// SecMaster.20030005: the incident not found
		err = common.ConvertExpected403ErrInto404Err(err, "code", WorkspaceNotFound)
		err = common.ConvertExpected400ErrInto404Err(err, "code", DataObjectRelationNotFound)
		return nil, err
	}

	getDataObjectRelationsRespBody, err := utils.FlattenResponse(getDataObjectRelationsResp)
	if err != nil {
		return nil, err
	}

	relatedDataObjectIds := utils.PathSearch("data[*].data_object.id",
		getDataObjectRelationsRespBody, make([]interface{}, 0)).([]interface{})
	if len(relatedDataObjectIds) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return relatedDataObjectIds, nil
}

func resourceDataObjectRelationsImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 4 {
		return nil, fmt.Errorf(`invalid format specified for import id,
			must be <workspace_id>/<data_class>/<data_object_id>/<related_data_class>`)
	}

	mErr := multierror.Append(
		d.Set("workspace_id", parts[0]),
		d.Set("data_class", parts[1]),
		d.Set("data_object_id", parts[2]),
		d.Set("related_data_class", parts[3]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
