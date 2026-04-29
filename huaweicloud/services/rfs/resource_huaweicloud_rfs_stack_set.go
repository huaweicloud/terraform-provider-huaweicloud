package rfs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RFS POST /v1/stack-sets
// @API RFS GET /v1/stack-sets/{stack_set_name}/metadata
// @API RFS PATCH /v1/stack-sets/{stack_set_name}
// @API RFS DELETE /v1/stack-sets/{stack_set_name}
func ResourceStackSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStackSetCreate,
		ReadContext:   resourceStackSetRead,
		UpdateContext: resourceStackSetUpdate,
		DeleteContext: resourceStackSetDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"stack_set_name",
			"permission_model",
			"template_body",
			"template_uri",
			"vars_body",
			"vars_uri",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"stack_set_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stack_set_description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"permission_model": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"administration_agency_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"managed_agency_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_body": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vars_body": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vars_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"initial_stack_description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"administration_agency_urn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"call_identity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"managed_operation": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_parallel_operation": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"stack_set_id": {
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vars_uri_content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organizational_unit_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceStackSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v1/stack-sets"
		requestId, _ = uuid.GenerateUUID()
		stackSetName = d.Get("stack_set_name").(string)
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateStackSetBodyParams(d)),
		MoreHeaders: map[string]string{
			"Client-Request-Id": requestId,
		},
	}
	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating RFS stack set: %s", err)
	}

	d.SetId(stackSetName)

	return resourceStackSetRead(ctx, d, meta)
}

func buildCreateStackSetBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"stack_set_name":             d.Get("stack_set_name"),
		"stack_set_description":      utils.ValueIgnoreEmpty(d.Get("stack_set_description")),
		"permission_model":           utils.ValueIgnoreEmpty(d.Get("permission_model")),
		"administration_agency_name": utils.ValueIgnoreEmpty(d.Get("administration_agency_name")),
		"managed_agency_name":        utils.ValueIgnoreEmpty(d.Get("managed_agency_name")),
		"template_body":              utils.ValueIgnoreEmpty(d.Get("template_body")),
		"template_uri":               utils.ValueIgnoreEmpty(d.Get("template_uri")),
		"vars_uri":                   utils.ValueIgnoreEmpty(d.Get("vars_uri")),
		"vars_body":                  utils.ValueIgnoreEmpty(d.Get("vars_body")),
		"initial_stack_description":  utils.ValueIgnoreEmpty(d.Get("initial_stack_description")),
		"administration_agency_urn":  utils.ValueIgnoreEmpty(d.Get("administration_agency_urn")),
		"managed_operation":          buildManagedOperationParams(d),
		"call_identity":              utils.ValueIgnoreEmpty(d.Get("call_identity")),
	}
}

func buildManagedOperationParams(d *schema.ResourceData) map[string]interface{} {
	raw := d.Get("managed_operation").([]interface{})
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}

	rawMap, ok := raw[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"enable_parallel_operation": rawMap["enable_parallel_operation"],
	}
}

func QueryStackSetMetaData(client *golangsdk.ServiceClient, stackSetName string) (interface{}, error) {
	var (
		httpUrl      = "v1/stack-sets/{stack_set_name}/metadata"
		requestId, _ = uuid.GenerateUUID()
	)
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{stack_set_name}", stackSetName)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Client-Request-Id": requestId,
		},
	}
	requestResp, err := client.Request("GET", requestPath, &opt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceStackSetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		stackSetName = d.Id()
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	resp, err := QueryStackSetMetaData(client, stackSetName)
	if err != nil {
		// If the resource does not exist, the response HTTP status code of the details API is `404`.
		return common.CheckDeletedDiag(d, err, "RFS resource stack")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("stack_set_name", utils.PathSearch("stack_set_name", resp, nil)),
		d.Set("stack_set_id", utils.PathSearch("stack_set_id", resp, nil)),
		d.Set("stack_set_description", utils.PathSearch("stack_set_description", resp, nil)),
		d.Set("initial_stack_description", utils.PathSearch("initial_stack_description", resp, nil)),
		d.Set("permission_model", utils.PathSearch("permission_model", resp, nil)),
		d.Set("administration_agency_name", utils.PathSearch("administration_agency_name", resp, nil)),
		d.Set("managed_agency_name", utils.PathSearch("managed_agency_name", resp, nil)),
		d.Set("status", utils.PathSearch("status", resp, nil)),
		d.Set("vars_uri_content", utils.PathSearch("vars_uri_content", resp, nil)),
		d.Set("vars_body", utils.PathSearch("vars_body", resp, nil)),
		d.Set("create_time", utils.PathSearch("create_time", resp, nil)),
		d.Set("update_time", utils.PathSearch("update_time", resp, nil)),
		d.Set("administration_agency_urn", utils.PathSearch("administration_agency_urn", resp, nil)),
		d.Set("organizational_unit_ids", utils.PathSearch("organizational_unit_ids", resp, nil)),
		d.Set("managed_operation", flattenManagedOperation(utils.PathSearch("managed_operation", resp, nil))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenManagedOperation(managedOp interface{}) []map[string]interface{} {
	if managedOp == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"enable_parallel_operation": utils.PathSearch("enable_parallel_operation", managedOp, nil),
		},
	}
}

func resourceStackSetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v1/stack-sets/{stack_set_name}"
		stackSetName = d.Get("stack_set_name").(string)
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request ID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{stack_set_name}", stackSetName)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Client-Request-Id": requestId,
		},
		JSONBody: utils.RemoveNil(buildUpdateStackSetBodyParams(d)),
	}

	_, err = client.Request("PATCH", requestPath, &opt)
	if err != nil {
		return diag.Errorf("error updating RFS stack set: %s", err)
	}

	return resourceStackSetRead(ctx, d, meta)
}

func buildUpdateStackSetBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"stack_set_id":               utils.ValueIgnoreEmpty(d.Get("stack_set_id")),
		"stack_set_description":      utils.ValueIgnoreEmpty(d.Get("stack_set_description")),
		"initial_stack_description":  utils.ValueIgnoreEmpty(d.Get("initial_stack_description")),
		"administration_agency_name": utils.ValueIgnoreEmpty(d.Get("administration_agency_name")),
		"managed_agency_name":        utils.ValueIgnoreEmpty(d.Get("managed_agency_name")),
		"administration_agency_urn":  utils.ValueIgnoreEmpty(d.Get("administration_agency_urn")),
		"managed_operation":          buildManagedOperationParams(d),
		"call_identity":              utils.ValueIgnoreEmpty(d.Get("call_identity")),
	}
}

func resourceStackSetDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v1/stack-sets/{stack_set_name}"
		stackSetName = d.Get("stack_set_name").(string)
		requestId, _ = uuid.GenerateUUID()
	)
	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{stack_set_name}", stackSetName)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Client-Request-Id": requestId,
		},
	}
	_, err = client.Request("DELETE", deletePath, &opt)
	if err != nil {
		return diag.Errorf("error deleting RFS stack set: %s", err)
	}

	return nil
}
