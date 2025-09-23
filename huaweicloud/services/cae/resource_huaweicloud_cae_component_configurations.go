package cae

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

var ConfigRelatedResourcesNotFoundCodes = []string{
	"CAE.01500208", // Application or component does not found.
	"CAE.01500404", // Environment does not found.
}

// ListNode is the structure of linked list node.
type ListNode struct {
	Timestamp int64
	Val       interface{} // The stored data for node.
	Next      *ListNode   // Pointer to next node.
}

// LinkedList is the structure of Singly Linked linked list.
type LinkedList struct {
	head *ListNode // Head node of linked list.
	size int       // Size of linked list.
}

// @API CAE POST /v1/{project_id}/cae/applications/{application_id}/components/{component_id}/configurations
// @API CAE GET /v1/{project_id}/cae/applications/{application_id}/components/{component_id}/configurations
// @API CAE DELETE /v1/{project_id}/cae/applications/{application_id}/components/{component_id}/configurations
func ResourceComponentConfigurations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComponentConfigurationsCreate,
		ReadContext:   resourceComponentConfigurationsRead,
		DeleteContext: resourceComponentConfigurationsDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceComponentConfigurationsImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region in which to create the resource.`,
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the environment where the application is located.`,
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the application where the component is located.`,
			},
			"component_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the component to which the configurations belong.`,
			},
			"items": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: `The type of the configuration.`,
						},
						"data": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  `The configuration detail.`,
						},
					},
				},
				Description: `The list of configurations for component.`,
			},
		},
	}
}

func buildCreateComponentConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"api_version": "v1",
		"kind":        "ComponentConfiguration",
		"items":       buildConfigurationItemsBodyParams(d.Get("items").(*schema.Set)),
	}
}

func buildConfigurationItemsBodyParams(items *schema.Set) []interface{} {
	if items.Len() < 1 {
		return nil
	}

	result := make([]interface{}, 0, items.Len())
	for _, v := range items.List() {
		result = append(result, map[string]interface{}{
			"type": utils.PathSearch("type", v, nil),
			"data": utils.ValueIgnoreEmpty(unmarshalJsonFormatParamster("Component configuration",
				utils.PathSearch("data", v, "").(string))),
		})
	}
	return result
}

func resourceComponentConfigurationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/cae/applications/{application_id}/components/{component_id}/configurations"
		componentId = d.Get("component_id").(string)
	)
	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE Client: %s", err)
	}

	modifyPath := client.Endpoint + httpUrl
	modifyPath = strings.ReplaceAll(modifyPath, "{project_id}", client.ProjectID)
	modifyPath = strings.ReplaceAll(modifyPath, "{application_id}", d.Get("application_id").(string))
	modifyPath = strings.ReplaceAll(modifyPath, "{component_id}", componentId)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Environment-ID": d.Get("environment_id").(string),
		},
	}
	opts.JSONBody = utils.RemoveNil(buildCreateComponentConfigurationBodyParams(d))
	// The creation API will return an 'operation_id' field, but this field is not helpful for resource logic, so it's
	// ignored.
	_, err = client.Request("POST", modifyPath, &opts)
	if err != nil {
		return diag.Errorf("error creating (overriding) the configurations for a component (%s): %s", componentId, err)
	}
	d.SetId(componentId)

	return resourceComponentConfigurationsRead(ctx, d, meta)
}

// NewLinkedList is a method that used to create a new linked list.
func NewLinkedList() *LinkedList {
	return &LinkedList{
		head: nil,
		size: 0,
	}
}

// Add is a method that used to add a new node to the head of the linked list.
func (l *LinkedList) Add(timestamp int64, value interface{}) {
	newNode := &ListNode{
		Timestamp: timestamp,
		Val:       value,
	}
	if l.head == nil {
		l.head = newNode
		l.size++
		return
	}
	current := l.head

	for current.Next != nil {
		current = current.Next
	}

	current.Next = newNode
	l.size++
}

// GetLatestEditingResult is a method that used to obtain the last editing results of specified type of configuration.
func (l *LinkedList) GetLatestEditingResult() interface{} {
	var (
		latestUpdateTimestamp int64
		value                 interface{}
	)
	current := l.head
	for current != nil {
		if current.Timestamp > latestUpdateTimestamp {
			latestUpdateTimestamp = current.Timestamp
			value = current.Val
		}
		current = current.Next
	}

	return value
}

func FilterActivatedConfigurations(configurations []interface{}) map[string]*LinkedList {
	activatedMap := make(map[string]*LinkedList)

	for _, v := range configurations {
		configDetail := utils.RemoveNil(utils.PathSearch("data", v, make(map[string]interface{})).(map[string]interface{}))
		if configDetail == nil {
			continue
		}
		configType := utils.PathSearch("type", v, "").(string)
		_, ok := activatedMap[configType]
		if !ok {
			activatedMap[configType] = NewLinkedList()
		}
		activatedMap[configType].Add(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("operated_at", v, "").(string)),
			marshalJsonFormatParamster("component configuration", configDetail))
	}

	return activatedMap
}

func resourceComponentConfigurationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/cae/applications/{application_id}/components/{component_id}/configurations"
		componentId = d.Get("component_id").(string)
	)
	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE Client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{application_id}", d.Get("application_id").(string))
	getPath = strings.ReplaceAll(getPath, "{component_id}", componentId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Environment-ID": d.Get("environment_id").(string),
		},
	}
	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", ConfigRelatedResourcesNotFoundCodes...),
			fmt.Sprintf("error querying configurations for the specified component (%s): %s", componentId, err))
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.Errorf("error retrieving configurations of the specified component (%s): %s", componentId, err)
	}
	items := FilterActivatedConfigurations(utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}))
	if len(items) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "CAE component configurations")
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("items", flattenConfigurationItems(items)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenConfigurationItems(items map[string]*LinkedList) []interface{} {
	result := make([]interface{}, 0, len(items))
	for k, v := range items {
		result = append(result, map[string]interface{}{
			"type": k,
			"data": v.GetLatestEditingResult(),
		})
	}
	return result
}

func resourceComponentConfigurationsDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/cae/applications/{application_id}/components/{component_id}/configurations"
		componentId = d.Get("component_id").(string)
	)
	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE Client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{application_id}", d.Get("application_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{component_id}", componentId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Environment-ID": d.Get("environment_id").(string),
		},
	}

	// Note: This API will delete all configurations under this component, regardless of whether it's configured
	// through other channels.
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting configurations under the specific component (%s): %s", componentId, err)
	}
	return nil
}

func resourceComponentConfigurationsImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want "+
			"'<environment_id>/<application_id>/<component_id>', but got '%s'", d.Id())
	}

	d.SetId(parts[2])
	mErr := multierror.Append(nil,
		d.Set("environment_id", parts[0]),
		d.Set("application_id", parts[1]),
		d.Set("component_id", parts[2]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
