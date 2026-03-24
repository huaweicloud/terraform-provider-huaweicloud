package cse

import (
	"context"
	"fmt"
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

var nacosNamespaceNonUpdatableParams = []string{
	"engine_id",
	"enterprise_project_id",
}

// @API CSE POST /v1/{project_id}/nacos/v1/console/namespaces
// @API CSE GET /v1/{project_id}/nacos/v1/console/namespaces
// @API CSE PUT /v1/{project_id}/nacos/v1/console/namespaces
// @API CSE DELETE /v1/{project_id}/nacos/v1/console/namespaces
func ResourceNacosNamespace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNacosNamespaceCreate,
		ReadContext:   resourceNacosNamespaceRead,
		UpdateContext: resourceNacosNamespaceUpdate,
		DeleteContext: resourceNacosNamespaceDelete,

		CustomizeDiff: config.FlexibleForceNew(nacosNamespaceNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceNacosNamespaceImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the Nacos namespace is located.`,
			},

			// Required parameters.
			"engine_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Nacos microservice engine to which the namespace belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the Nacos namespace.`,
			},

			// Optional parameters.
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the enterprise project to which the Nacos namespace belongs.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					}),
			},
		},
	}
}

func buildNacosNamespaceRequestHeaders(engineId, enterpriseProjectId string) map[string]string {
	moreHeaders := map[string]string{
		"Content-Type": "application/json",
		"x-engine-id":  engineId,
	}

	if enterpriseProjectId != "" {
		moreHeaders["X-Enterprise-Project-ID"] = enterpriseProjectId
	}
	return moreHeaders
}

func createNacosNamespace(client *golangsdk.ServiceClient, engineId, enterpriseProjectId, namespaceId, namespaceName string) error {
	httpUrl := "v1/{project_id}/nacos/v1/console/namespaces?customNamespaceId={namespace_id}&namespaceName={namespace_name}"

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{namespace_id}", namespaceId)
	createPath = strings.ReplaceAll(createPath, "{namespace_name}", namespaceName)

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildNacosNamespaceRequestHeaders(engineId, enterpriseProjectId),
	}

	_, err := client.Request("POST", createPath, &createOpts)
	return err
}

func resourceNacosNamespaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                 = meta.(*config.Config)
		region              = cfg.GetRegion(d)
		engineId            = d.Get("engine_id").(string)
		enterpriseProjectId = cfg.GetEnterpriseProjectID(d)
		name                = d.Get("name").(string)
	)

	client, err := cfg.NewServiceClient("cse", region)
	if err != nil {
		return diag.Errorf("error creating CSE client: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate resource ID of Nacos namespace (%s): %s", name, err)
	}

	err = createNacosNamespace(client, engineId, enterpriseProjectId, resourceId, name)
	if err != nil {
		// When the Nacos engine does not exist, the error code returned is 400, and the error information is:
		// {"error_code": "SVCSTG.00401116", "error_message": "engine does not exist"}
		// So, we need to convert the error code to 404 error and specify the key of the error code to be "error_code".
		parsedErr := common.ConvertExpected400ErrInto404Err(err, "error_code", microserviceEngineNotFoundCodes...)
		switch parsedErr.(type) {
		case golangsdk.ErrDefault404, golangsdk.ErrDefault502:
			// Bad gateway means Nacos engine does not exist.
			return diag.Errorf("unable to create the namespace because the Nacos engine (%s) does not exist", engineId)
		default:
			return diag.Errorf("error creating namespace under Nacos microservice engine (%s): %s", engineId, err)
		}
	}

	d.SetId(resourceId)

	return resourceNacosNamespaceRead(ctx, d, meta)
}

func listNacosNamespaces(client *golangsdk.ServiceClient, engineId, enterpriseProjectId string) ([]interface{}, error) {
	httpUrl := "v1/{project_id}/nacos/v1/console/namespaces"

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildNacosNamespaceRequestHeaders(engineId, enterpriseProjectId),
	}

	requestResp, err := client.Request("GET", listPath, &listOpts)
	if err != nil {
		parsedErr := common.ConvertExpected400ErrInto404Err(err, "error_code", microserviceEngineNotFoundCodes...)
		switch parsedErr.(type) {
		case golangsdk.ErrDefault404, golangsdk.ErrDefault502:
			// Bad gateway means Nacos engine does not exist.
			return nil, golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Method:    "GET",
					URL:       httpUrl,
					RequestId: "NONE",
					Body:      []byte(fmt.Sprintf("the Nacos engine (%s) does not exist", engineId)),
				},
			}
		default:
			return nil, err
		}
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{}), nil
}

// GetNacosNamespaceById retrieves the Nacos namespace by its ID.
func GetNacosNamespaceById(client *golangsdk.ServiceClient, engineId, enterpriseProjectId, namespaceId string) (interface{}, error) {
	namespaces, err := listNacosNamespaces(client, engineId, enterpriseProjectId)
	if err != nil {
		return nil, err
	}

	namespace := utils.PathSearch(fmt.Sprintf("[?namespace=='%s']|[0]", namespaceId), namespaces, nil)
	if namespace == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "v1/{project_id}/nacos/v1/console/namespaces",
				RequestId: "NONE",
				Body: []byte(fmt.Sprintf("the namespace (%s) has been removed from the Nacos microservice engine (%s)",
					namespaceId, engineId)),
			},
		}
	}
	return namespace, nil
}

func resourceNacosNamespaceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                 = meta.(*config.Config)
		region              = cfg.GetRegion(d)
		engineId            = d.Get("engine_id").(string)
		enterpriseProjectId = cfg.GetEnterpriseProjectID(d)
		namespaceId         = d.Id()
	)
	client, err := cfg.NewServiceClient("cse", region)
	if err != nil {
		return diag.Errorf("error creating CSE client: %s", err)
	}

	respBody, err := GetNacosNamespaceById(client, engineId, enterpriseProjectId, namespaceId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying Nacos namespace (%s)", namespaceId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		// Required parameters.
		d.Set("engine_id", engineId),
		d.Set("name", utils.PathSearch("namespaceShowName", respBody, nil)),
		// Optional parameters.
		d.Set("enterprise_project_id", enterpriseProjectId),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func updateNacosNamespace(client *golangsdk.ServiceClient, engineId, enterpriseProjectId, namespaceId, namespaceName string) error {
	httpUrl := "v1/{project_id}/nacos/v1/console/namespaces?namespace={namespace_id}&namespaceShowName={namespace_name}"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{namespace_id}", namespaceId)
	updatePath = strings.ReplaceAll(updatePath, "{namespace_name}", namespaceName)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildNacosNamespaceRequestHeaders(engineId, enterpriseProjectId),
	}

	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func resourceNacosNamespaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.HasChange("name") {
		return resourceNacosNamespaceRead(ctx, d, meta)
	}

	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		engineId   = d.Get("engine_id").(string)
		name       = d.Get("name").(string)
		resourceId = d.Id()
	)

	client, err := cfg.NewServiceClient("cse", region)
	if err != nil {
		return diag.Errorf("error creating CSE client: %s", err)
	}

	err = updateNacosNamespace(client, engineId, cfg.GetEnterpriseProjectID(d), resourceId, name)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault502); ok {
			// Bad gateway means Nacos engine does not exist.
			return diag.Errorf("unable to update the namespace because the Nacos engine (%s) does not exist", engineId)
		}
		return diag.Errorf("error updating namespace (%s) under the Nacos microservice engine (%s): %s", resourceId, engineId, err)
	}

	return resourceNacosNamespaceRead(ctx, d, meta)
}

func deleteNacosNamespace(client *golangsdk.ServiceClient, engineId, enterpriseProjectId, namespaceId string) error {
	httpUrl := "v1/{project_id}/nacos/v1/console/namespaces?namespaceId={namespace_id}"

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{namespace_id}", namespaceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildNacosNamespaceRequestHeaders(engineId, enterpriseProjectId),
	}

	_, err := client.Request("DELETE", deletePath, &opt)
	return err
}

func resourceNacosNamespaceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                 = meta.(*config.Config)
		region              = cfg.GetRegion(d)
		engineId            = d.Get("engine_id").(string)
		enterpriseProjectId = cfg.GetEnterpriseProjectID(d)
		namespaceId         = d.Id()
	)

	client, err := cfg.NewServiceClient("cse", region)
	if err != nil {
		return diag.Errorf("error creating CSE client: %s", err)
	}

	// Delete method always return a 200 status code whether namespace is exist.
	err = deleteNacosNamespace(client, engineId, enterpriseProjectId, namespaceId)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault502); ok {
			// Bad gateway means Nacos engine has been removed and returns override the error with 404 status code.
			err = golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Method:    "DELETE",
					URL:       "v1/{project_id}/nacos/v1/console/namespaces?namespaceId={namespace_id}",
					RequestId: "NONE",
					Body: []byte(fmt.Sprintf("unable to delete the namespace because the Nacos microservice engine (%s) does not exist",
						engineId)),
				},
			}
		}
		// When the Nacos engine does not exist, the error code returned is 400, and the error information is:
		// {"error_code": "SVCSTG.00401116", "error_message": "engine does not exist"}
		// So, we need to convert the error code to 404 error and specify the key of the error code to be "error_code".
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", microserviceEngineNotFoundCodes...),
			fmt.Sprintf("error deleting namespace (%s) from the Nacos microservice engine (%s)", namespaceId, engineId))
	}
	return nil
}

func resourceNacosNamespaceImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")

	switch len(parts) {
	case 2:
		d.Set("enterprise_project_id", nil)
	case 3:
		d.Set("enterprise_project_id", parts[2])
	default:
		return nil, fmt.Errorf("invalid format specified for import ID, want '<engine_id>/<id>/<enterprise_project_id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("engine_id", parts[0])
}
