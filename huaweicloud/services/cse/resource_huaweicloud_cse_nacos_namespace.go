package cse

import (
	"context"
	"fmt"
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
		},
	}
}

func resourceNacosNamespaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/nacos/v1/console/namespaces?customNamespaceId={namespace_id}&namespaceName={namespace_name}"

		engineId = d.Get("engine_id").(string)
		name     = d.Get("name").(string)
	)

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate resource ID for CSE Nacos namespace (%s): %s", name, err)
	}

	client, err := cfg.NewServiceClient("cse", region)
	if err != nil {
		return diag.Errorf("error creating CSE client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{namespace_id}", resourceId)
	createPath = strings.ReplaceAll(createPath, "{namespace_name}", name)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"x-engine-id":  engineId,
		},
	}

	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault502); ok {
			// Bad gateway means Nacos engine does not exist.
			return diag.Errorf("unable to create the namespace because the Nacos engine (%s) does not exist", engineId)
		}
		return diag.Errorf("error creating namespace under Nacos microservice engine (%s): %s", engineId, err)
	}

	d.SetId(resourceId)

	return resourceNacosNamespaceRead(ctx, d, meta)
}

func listNacosNamespaces(client *golangsdk.ServiceClient, engineId string) ([]interface{}, error) {
	httpUrl := "v1/{project_id}/nacos/v1/console/namespaces"

	queryPath := client.Endpoint + httpUrl
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"x-engine-id":  engineId,
		},
	}

	requestResp, err := client.Request("GET", queryPath, &opt)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault502); ok {
			// Bad gateway means Nacos engine does not exist.
			return nil, golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Method:    "GET",
					URL:       httpUrl,
					RequestId: "NONE",
					Body:      []byte(fmt.Sprintf("the Nacos engine (%s) does not exist", engineId)),
				},
			}
		}
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func GetNacosNamespaceById(client *golangsdk.ServiceClient, engineId, namespaceId string) (interface{}, error) {
	namespaces, err := listNacosNamespaces(client, engineId)
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
					engineId, namespaceId)),
			},
		}
	}
	return namespace, nil
}

func resourceNacosNamespaceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		engineId    = d.Get("engine_id").(string)
		namespaceId = d.Id()
	)
	client, err := cfg.NewServiceClient("cse", region)
	if err != nil {
		return diag.Errorf("error creating CSE client: %s", err)
	}

	respBody, err := GetNacosNamespaceById(client, engineId, namespaceId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying Nacos namespace (%s)", namespaceId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("namespaceShowName", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceNacosNamespaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/nacos/v1/console/namespaces?namespace={namespace_id}&namespaceShowName={namespace_name}"

		engineId   = d.Get("engine_id").(string)
		name       = d.Get("name").(string)
		resourceId = d.Id()
	)

	client, err := cfg.NewServiceClient("cse", region)
	if err != nil {
		return diag.Errorf("error creating CSE client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{namespace_id}", resourceId)
	updatePath = strings.ReplaceAll(updatePath, "{namespace_name}", name)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"x-engine-id":  engineId,
		},
	}

	_, err = client.Request("PUT", updatePath, &opt)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault502); ok {
			// Bad gateway means Nacos engine does not exist.
			return diag.Errorf("unable to update the namespace because the Nacos engine (%s) does not exist", engineId)
		}
		return diag.Errorf("error updating namespace (%s) under the Nacos microservice engine (%s): %s", resourceId, engineId, err)
	}

	return resourceNacosNamespaceRead(ctx, d, meta)
}

func resourceNacosNamespaceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/nacos/v1/console/namespaces?namespaceId={namespace_id}"

		engineId    = d.Get("engine_id").(string)
		namespaceId = d.Id()
	)

	client, err := cfg.NewServiceClient("cse", region)
	if err != nil {
		return diag.Errorf("error creating CSE client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{namespace_id}", namespaceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"x-engine-id":  engineId,
		},
	}

	_, err = client.Request("DELETE", deletePath, &opt)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault502); ok {
			// Bad gateway means Nacos engine has been removed and returns override the error with 404 status code.
			err = golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Method:    "DELETE",
					URL:       httpUrl,
					RequestId: "NONE",
					Body:      []byte(fmt.Sprintf("unable to delete the namespace because the Nacos engine (%s) does not exist", engineId)),
				},
			}
		}
		// Delete method always return a 200 status code whether namespace is exist.
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error deleting namespace (%s) from the Nacos microserivce engine (%s)", namespaceId, engineId))
	}
	return nil
}

func resourceNacosNamespaceImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<engine_id>/<id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("engine_id", parts[0])
}
