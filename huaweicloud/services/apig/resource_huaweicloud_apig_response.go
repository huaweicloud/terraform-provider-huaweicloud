package apig

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/responses"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG GET /v2/{project_id}/apigw/instances/{instanceId}/api-groups/{group_id}/gateway-responses/{response_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instanceId}/api-groups/{group_id}/gateway-responses/{response_id}
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instanceId}/api-groups/{group_id}/gateway-responses/{response_id}
// @API APIG GET /v2/{project_id}/apigw/instances/{instanceId}/api-groups/{group_id}/gateway-responses
// @API APIG POST /v2/{project_id}/apigw/instances/{instanceId}/api-groups/{group_id}/gateway-responses
func ResourceApigResponseV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResponseCreate,
		ReadContext:   resourceResponseRead,
		UpdateContext: resourceResponseUpdate,
		DeleteContext: resourceResponseDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCustomResponseImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the API custom response is located.",
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The ID of the dedicated instance to which the API group and the API custom response " +
					"belongs.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the API group to which the API custom response belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the API custom response.",
			},
			"rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"error_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The error type of the API custom response rule.",
						},
						"body": {
							Type: schema.TypeString,
							// If parameter body omitted, The API will return 'The parameters must be specified'.
							Required:    true,
							Description: "The body template of the API custom response rule.",
						},
						"status_code": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The HTTP status code of the API custom response rule.",
						},
						"headers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The key name of the response header.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The value for the specified response header key.",
									},
								},
							},
							Description: "The configuration of the custom response headers.",
						},
					},
				},
				Description: "The API custom response rules definition.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the API custom response.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the API custom response.",
			},
		},
	}
}

func buildCustomResponseHeaders(headers *schema.Set) []responses.ResponseInfoHeader {
	if headers.Len() < 1 {
		return nil
	}

	result := make([]responses.ResponseInfoHeader, 0, headers.Len())
	for _, v := range headers.List() {
		result = append(result, responses.ResponseInfoHeader{
			Key:   utils.PathSearch("key", v, "").(string),
			Value: utils.PathSearch("value", v, "").(string),
		})
	}
	return result
}

// 'error_type' is the key of the response mapping, and 'body' and 'status_code' are the structural elements of the
// mapping value.
func buildCustomResponses(respSet *schema.Set) map[string]responses.ResponseInfo {
	result := make(map[string]responses.ResponseInfo)

	for _, response := range respSet.List() {
		rule := response.(map[string]interface{})

		result[rule["error_type"].(string)] = responses.ResponseInfo{
			Body:    rule["body"].(string),
			Status:  rule["status_code"].(int),
			Headers: buildCustomResponseHeaders(rule["headers"].(*schema.Set)),
		}
	}

	return result
}

func resourceResponseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}
	opt := responses.ResponseOpts{
		InstanceId: d.Get("instance_id").(string),
		GroupId:    d.Get("group_id").(string),
		Name:       d.Get("name").(string),
		Responses:  buildCustomResponses(d.Get("rule").(*schema.Set)),
	}
	resp, err := responses.Create(client, opt).Extract()
	if err != nil {
		return diag.Errorf("error creating APIG custom response: %s", err)
	}
	d.SetId(resp.Id)

	return resourceResponseRead(ctx, d, meta)
}

func flattenCustomResponseHeaders(headers []responses.ResponseInfoHeader) []interface{} {
	if len(headers) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(headers))
	for _, v := range headers {
		result = append(result, map[string]interface{}{
			"key":   v.Key,
			"value": v.Value,
		})
	}
	return result
}

func flattenCustomResponses(respMap map[string]responses.ResponseInfo) []map[string]interface{} {
	if len(respMap) < 1 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(respMap))
	for errorType, rule := range respMap {
		if rule.IsDefault {
			// The IsDefault of the modified response will be marked as false,
			// record these responses and skip other unmodified responses (IsDefault is true).
			continue
		}
		result = append(result, map[string]interface{}{
			"error_type":  errorType,
			"body":        rule.Body,
			"status_code": rule.Status,
			"headers":     flattenCustomResponseHeaders(rule.Headers),
		})
	}
	return result
}

func resourceResponseRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		groupId    = d.Get("group_id").(string)
		responseId = d.Id()
	)
	resp, err := responses.Get(client, instanceId, groupId, responseId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "APIG custom response")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("rule", flattenCustomResponses(resp.Responses)),
		d.Set("created_at", resp.CreateTime),
		d.Set("updated_at", resp.UpdateTime),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving API custom response (%s) fields: %s", responseId, err)
	}
	return nil
}

func resourceResponseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}
	// Only updating the name will cause all the response rules that have been set to be reset, so no matter whether
	// the response rules are updated or not, the response rules must be carried in the update opt.
	opt := responses.ResponseOpts{
		InstanceId: d.Get("instance_id").(string),
		GroupId:    d.Get("group_id").(string),
		Name:       d.Get("name").(string),
		Responses:  buildCustomResponses(d.Get("rule").(*schema.Set)),
	}
	_, err = responses.Update(client, d.Id(), opt).Extract()
	if err != nil {
		return diag.Errorf("error updating APIG custom response: %s", err)
	}
	return resourceResponseRead(ctx, d, meta)
}

func resourceResponseDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		groupId    = d.Get("group_id").(string)
		responseId = d.Id()
	)
	err = responses.Delete(client, instanceId, groupId, responseId).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting APIG custom response (%s) from the dedicated group (%s)",
			responseId, groupId))
	}
	return nil
}

// Some resources of the APIG service are associated with dedicated instances and groups,
// but their IDs cannot be found on the console.
// This method is used to solve the above problem by importing resources by associating ID and their name.
func resourceCustomResponseImportState(_ context.Context, d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import IDs and name, " +
			"must be <instance_id>/<group_id>/<name>")
	}

	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = parts[0]
		groupId    = parts[1]
		name       = parts[2]

		opts = responses.ListOpts{
			InstanceId: instanceId,
			GroupId:    groupId,
		}
	)
	pages, err := responses.List(client, opts).AllPages()
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error getting custom response list from server: %s", err)
	}
	resp, err := responses.ExtractResponses(pages)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error extract custom responses: %s", err)
	}
	if len(resp) < 1 {
		return []*schema.ResourceData{d}, fmt.Errorf("unable to find any custom response from server")
	}
	// Since there are no parameters about custom responses in the query options, we need to get the response list and
	// filter by the response name.
	for _, val := range resp {
		if val.Name == name {
			d.SetId(val.Id)
			mErr := multierror.Append(nil,
				d.Set("instance_id", instanceId),
				d.Set("group_id", groupId),
			)
			return []*schema.ResourceData{d}, mErr.ErrorOrNil()
		}
	}
	return []*schema.ResourceData{d}, fmt.Errorf("unable to find the custom response (%s) from server", name)
}
