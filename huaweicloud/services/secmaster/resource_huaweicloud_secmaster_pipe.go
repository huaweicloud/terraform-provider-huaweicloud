package secmaster

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

var nonUpdatableParamsPipe = []string{"workspace_id", "dataspace_id", "pipe_name", "mapping", "timestamp_field"}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/siem/pipes
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}/index
func ResourcePipe() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipeCreate,
		ReadContext:   resourcePipeRead,
		UpdateContext: resourcePipeUpdate,
		DeleteContext: resourcePipeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePipeImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsPipe),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dataspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipe_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"shards": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"storage_period": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mapping": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(_, o, n string, _ *schema.ResourceData) bool {
					equal, _ := utils.CompareJsonTemplateAreEquivalent(o, n)
					return equal
				},
				Description: "The mapping configuration for the pipe. This is a JSON object where keys are dynamic field names. " +
					"Each field must contain 'type' (string), 'is_chinese_exist' (boolean), and 'properties' (object). " +
					"Example: {\"field1\": {\"type\": \"text\", \"is_chinese_exist\": true, \"properties\": {}}, \"field2\": {...}}",
			},
			"timestamp_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"pipe_alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pipe_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dataspace_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"process_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreatePipeBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	bodyParams := map[string]interface{}{
		"dataspace_id":   d.Get("dataspace_id"),
		"pipe_name":      d.Get("pipe_name"),
		"shards":         d.Get("shards"),
		"storage_period": d.Get("storage_period"),
	}

	if v, ok := d.GetOk("description"); ok {
		bodyParams["description"] = v
	}

	if v, ok := d.GetOk("mapping"); ok {
		var mappingObj interface{}
		if err := json.Unmarshal([]byte(v.(string)), &mappingObj); err != nil {
			return nil, fmt.Errorf("error parsing mapping parameter: %s", err)
		}
		bodyParams["mapping"] = mappingObj
	}

	if v, ok := d.GetOk("timestamp_field"); ok {
		bodyParams["timestamp_field"] = v
	}

	return bodyParams, nil
}

func resourcePipeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		workspaceId   = d.Get("workspace_id").(string)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/siem/pipes"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)

	bodyParams, err := buildCreatePipeBodyParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         bodyParams,
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster pipe: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	pipeId := utils.PathSearch("pipe_id", respBody, "").(string)
	if pipeId == "" {
		return diag.Errorf("error creating SecMaster pipe: unable to find pipe ID")
	}

	d.SetId(pipeId)

	return resourcePipeRead(ctx, d, meta)
}

func GetPipeInfo(client *golangsdk.ServiceClient, workspaceId, pipeId string) (interface{}, error) {
	getPath := client.Endpoint + "v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath = strings.ReplaceAll(getPath, "{pipe_id}", pipeId)
	getOpts := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func GetPipeIndexInfo(client *golangsdk.ServiceClient, workspaceId, pipeId string) (interface{}, error) {
	getPath := client.Endpoint + "v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}/index"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath = strings.ReplaceAll(getPath, "{pipe_id}", pipeId)
	getOpts := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourcePipeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	respBody, err := GetPipeInfo(client, workspaceId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "code", "SecMaster.20030001"),
			"error retrieving SecMaster pipe",
		)
	}

	indexRespBody, err := GetPipeIndexInfo(client, workspaceId, d.Id())
	if err != nil {
		log.Printf("[WARN] error retrieving SecMaster pipe index: %s", err)
	}

	mappingValue := utils.PathSearch("mapping", indexRespBody, nil)
	var mappingStr string
	if mappingValue != nil {
		mappingBytes, err := json.Marshal(mappingValue)
		if err != nil {
			log.Printf("[WARN] error marshaling mapping field: %s", err)
		} else {
			mappingStr = string(mappingBytes)
		}
	}

	timestampField := utils.PathSearch("timestamp_field", indexRespBody, nil)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("pipe_name", utils.PathSearch("pipe_name", respBody, nil)),
		d.Set("pipe_alias", utils.PathSearch("pipe_alias", respBody, nil)),
		d.Set("pipe_type", utils.PathSearch("pipe_type", respBody, nil)),
		d.Set("dataspace_id", utils.PathSearch("dataspace_id", respBody, nil)),
		d.Set("dataspace_name", utils.PathSearch("dataspace_name", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("mapping", mappingStr),
		d.Set("shards", utils.PathSearch("shards", respBody, nil)),
		d.Set("storage_period", utils.PathSearch("storage_period", respBody, nil)),
		d.Set("category", utils.PathSearch("category", respBody, nil)),
		d.Set("owner_type", utils.PathSearch("owner_type", respBody, nil)),
		d.Set("process_status", utils.PathSearch("process_status", respBody, nil)),
		d.Set("create_by", utils.PathSearch("create_by", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
		d.Set("update_by", utils.PathSearch("update_by", respBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", respBody, nil)),
		d.Set("domain_id", utils.PathSearch("domain_id", respBody, nil)),
		d.Set("project_id", utils.PathSearch("project_id", respBody, nil)),
		d.Set("timestamp_field", timestampField),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdatePipeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{}

	if v, ok := d.GetOk("description"); ok {
		bodyParams["description"] = v
	}

	if v, ok := d.GetOk("shards"); ok {
		bodyParams["shards"] = v
	}
	if v, ok := d.GetOk("storage_period"); ok {
		bodyParams["storage_period"] = v
	}

	return bodyParams
}

func resourcePipeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{workspace_id}", workspaceId)
	updatePath = strings.ReplaceAll(updatePath, "{pipe_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         buildUpdatePipeBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating SecMaster pipe: %s", err)
	}

	return resourcePipeRead(ctx, d, meta)
}

func resourcePipeDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{workspace_id}", workspaceId)
	deletePath = strings.ReplaceAll(deletePath, "{pipe_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SecMaster pipe")
	}

	return nil
}

func resourcePipeImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<pipe_id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])

	mErr := multierror.Append(nil,
		d.Set("workspace_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
