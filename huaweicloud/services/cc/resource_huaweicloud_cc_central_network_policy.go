// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CC
// ---------------------------------------------------------------

package cc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CC POST /v3/{domain_id}/gcn/central-network/{central_network_id}/policies
// @API CC DELETE /v3/{domain_id}/gcn/central-network/{central_network_id}/policies/{id}
// @API CC GET /v3/{domain_id}/gcn/central-network/{central_network_id}/policies?id={id}
func ResourceCentralNetworkPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCentralNetworkPolicyCreate,
		ReadContext:   resourceCentralNetworkPolicyRead,
		DeleteContext: resourceCentralNetworkPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCentralNetworkPolicyImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"central_network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Central network ID.`,
			},
			"er_instances": {
				Type:        schema.TypeList,
				MinItems:    1,
				Elem:        centralNetworkPolicyAssociateErInstanceDocumentSchema(),
				Required:    true,
				ForceNew:    true,
				Description: `List of the enterprise routers on the central network policy.`,
			},
			"planes": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        centralNetworkPolicyCentralNetworkPolicyPlaneDocumentSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `List of the central network policy planes.`,
			},
			"document_template_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Central network policy document template version.`,
			},
			"is_applied": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the central network policy is applied.`,
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Central network policy version.`,
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Central network policy status.`,
			},
		},
	}
}

func centralNetworkPolicyCentralNetworkPolicyPlaneDocumentSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"associate_er_tables": {
				Type:        schema.TypeList,
				Elem:        centralNetworkPolicyCentralNetworkPolicyPlaneDocumentAssociateErTableDocumentSchema(),
				Required:    true,
				ForceNew:    true,
				Description: `List of route tables associated with the central network policy.`,
			},
			"exclude_er_connections": {
				Type:        schema.TypeList,
				Elem:        centralNetworkPolicyCentralNetworkPolicyPlaneDocumentExcludeErConnectionsSchema(),
				Optional:    true,
				ForceNew:    true,
				Description: `List of the enterprise router connections excluded from the central network policy.`,
			},
		},
	}
	return &sc
}

func centralNetworkPolicyCentralNetworkPolicyPlaneDocumentAssociateErTableDocumentSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Project ID.`,
			},
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Region ID.`,
			},
			"enterprise_router_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Enterprise router ID.`,
			},
			"enterprise_router_table_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Enterprise router table ID.`,
			},
		},
	}
	return &sc
}

func centralNetworkPolicyCentralNetworkPolicyPlaneDocumentExcludeErConnectionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"exclude_er_instances": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem:        centralNetworkPolicyAssociateErInstanceDocumentSchema(),
				Description: `List of enterprise routers that will not establish a connection.`,
			},
		},
	}
}

func centralNetworkPolicyAssociateErInstanceDocumentSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Project ID.`,
			},
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Region ID.`,
			},
			"enterprise_router_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Enterprise router ID.`,
			},
		},
	}
	return &sc
}

func resourceCentralNetworkPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createCentralNetworkPolicy: create a central network policy
	var (
		createCentralNetworkPolicyHttpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/policies"
		createCentralNetworkPolicyProduct = "cc"
	)
	createCentralNetworkPolicyClient, err := cfg.NewServiceClient(createCentralNetworkPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	createCentralNetworkPolicyPath := createCentralNetworkPolicyClient.Endpoint + createCentralNetworkPolicyHttpUrl
	createCentralNetworkPolicyPath = strings.ReplaceAll(createCentralNetworkPolicyPath, "{domain_id}", cfg.DomainID)
	createCentralNetworkPolicyPath = strings.ReplaceAll(createCentralNetworkPolicyPath, "{central_network_id}",
		d.Get("central_network_id").(string))

	createCentralNetworkPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createCentralNetworkPolicyOpt.JSONBody = utils.RemoveNil(buildCreateCentralNetworkPolicyBodyParams(d))
	createCentralNetworkPolicyResp, err := createCentralNetworkPolicyClient.Request("POST",
		createCentralNetworkPolicyPath, &createCentralNetworkPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating CentralNetworkPolicy: %s", err)
	}

	createCentralNetworkPolicyRespBody, err := utils.FlattenResponse(createCentralNetworkPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("central_network_policy.id", createCentralNetworkPolicyRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CentralNetworkPolicy: ID is not found in API response")
	}
	d.SetId(id)

	return resourceCentralNetworkPolicyRead(ctx, d, meta)
}

func buildCreateCentralNetworkPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"central_network_policy_document": map[string]interface{}{
			"default_plane": "default-plane",
			"planes":        buildCreateCentralNetworkPolicyRequestBodyCentralNetworkPolicyPlaneDocument(d.Get("planes")),
			"er_instances":  buildCreateCentralNetworkPolicyRequestBodyAssociateErInstanceDocument(d.Get("er_instances")),
		},
	}
	return bodyParams
}

func buildCreateCentralNetworkPolicyRequestBodyCentralNetworkPolicyPlaneDocument(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok && len(rawArray) > 0 {
		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name":                   "default-plane",
				"associate_er_tables":    buildCentralNetworkPolicyPlaneAssociateErTableDocument(raw["associate_er_tables"]),
				"exclude_er_connections": buildCentralNetworkPolicyPlaneExcludeErConnections(raw["exclude_er_connections"]),
			}
		}
		return rst
	}

	return []map[string]interface{}{{"name": "default-plane"}}
}

func buildCentralNetworkPolicyPlaneAssociateErTableDocument(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"project_id":                 utils.ValueIgnoreEmpty(raw["project_id"]),
				"region_id":                  utils.ValueIgnoreEmpty(raw["region_id"]),
				"enterprise_router_id":       utils.ValueIgnoreEmpty(raw["enterprise_router_id"]),
				"enterprise_router_table_id": utils.ValueIgnoreEmpty(raw["enterprise_router_table_id"]),
			}
		}
		return rst
	}
	return nil
}

func buildCentralNetworkPolicyPlaneExcludeErConnections(rawParams interface{}) [][]interface{} {
	if rawConnections, ok := rawParams.([]interface{}); ok {
		connections := make([][]interface{}, len(rawConnections))
		for i, rawConnection := range rawConnections {
			connection := rawConnection.(map[string]interface{})
			erInstancesRaw := connection["exclude_er_instances"].([]interface{})
			erInstances := make([]interface{}, len(erInstancesRaw))
			for j, erInstanceRaw := range erInstancesRaw {
				v := erInstanceRaw.(map[string]interface{})
				erInstances[j] = map[string]interface{}{
					"project_id":           utils.ValueIgnoreEmpty(v["project_id"]),
					"region_id":            utils.ValueIgnoreEmpty(v["region_id"]),
					"enterprise_router_id": utils.ValueIgnoreEmpty(v["enterprise_router_id"]),
				}
			}
			connections[i] = erInstances
		}
		return connections
	}
	return nil
}

func buildCreateCentralNetworkPolicyRequestBodyAssociateErInstanceDocument(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"project_id":           utils.ValueIgnoreEmpty(raw["project_id"]),
				"region_id":            utils.ValueIgnoreEmpty(raw["region_id"]),
				"enterprise_router_id": utils.ValueIgnoreEmpty(raw["enterprise_router_id"]),
			}
		}
		return rst
	}
	return nil
}

func resourceCentralNetworkPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getCentralNetworkPolicy: Query the central network policy
	var (
		getCentralNetworkPolicyHttpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/policies?id={id}"
		getCentralNetworkPolicyProduct = "cc"
	)
	getCentralNetworkPolicyClient, err := cfg.NewServiceClient(getCentralNetworkPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	getCentralNetworkPolicyPath := getCentralNetworkPolicyClient.Endpoint + getCentralNetworkPolicyHttpUrl
	getCentralNetworkPolicyPath = strings.ReplaceAll(getCentralNetworkPolicyPath, "{domain_id}", cfg.DomainID)
	getCentralNetworkPolicyPath = strings.ReplaceAll(getCentralNetworkPolicyPath, "{central_network_id}",
		d.Get("central_network_id").(string))
	getCentralNetworkPolicyPath = strings.ReplaceAll(getCentralNetworkPolicyPath, "{id}", d.Id())

	getCentralNetworkPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getCentralNetworkPolicyResp, err := getCentralNetworkPolicyClient.Request("GET", getCentralNetworkPolicyPath,
		&getCentralNetworkPolicyOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CentralNetworkPolicy")
	}

	respBodyJson, err := utils.FlattenResponse(getCentralNetworkPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("central_network_policies[?id =='%s']|[0]", d.Id())
	respBodyJson = utils.PathSearch(jsonPath, respBodyJson, nil)
	if respBodyJson == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("central_network_id", utils.PathSearch("central_network_id", respBodyJson, nil)),
		d.Set("document_template_version", utils.PathSearch("document_template_version", respBodyJson, nil)),
		d.Set("is_applied", utils.PathSearch("is_applied", respBodyJson, nil)),
		d.Set("version", utils.PathSearch("version", respBodyJson, nil)),
		d.Set("state", utils.PathSearch("state", respBodyJson, nil)),
		d.Set("planes", flattenGetCentralNetworkPolicyResponseBodyCentralNetworkPolicyPlaneDocument(respBodyJson)),
		d.Set("er_instances", flattenGetCentralNetworkPolicyResponseBodyAssociateErInstanceDocument(respBodyJson)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetCentralNetworkPolicyResponseBodyCentralNetworkPolicyPlaneDocument(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("document.planes", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"associate_er_tables":    flattenCentralNetworkPolicyPlaneDocumentAssociateErTables(v),
			"exclude_er_connections": flattenCentralNetworkPolicyPlaneDocumentExcludeErConnections(v),
		})
	}
	return rst
}

func flattenCentralNetworkPolicyPlaneDocumentAssociateErTables(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("associate_er_tables", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"project_id":                 utils.PathSearch("project_id", v, nil),
			"region_id":                  utils.PathSearch("region_id", v, nil),
			"enterprise_router_id":       utils.PathSearch("enterprise_router_id", v, nil),
			"enterprise_router_table_id": utils.PathSearch("enterprise_router_table_id", v, nil),
		})
	}
	return rst
}

func flattenCentralNetworkPolicyPlaneDocumentExcludeErConnections(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("exclude_er_connections", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, connectionRaw := range curArray {
		connection := connectionRaw.([]interface{})
		erInstances := make([]interface{}, 0, len(connection))
		for _, v := range connection {
			erInstances = append(erInstances, map[string]interface{}{
				"project_id":           utils.PathSearch("project_id", v, nil),
				"region_id":            utils.PathSearch("region_id", v, nil),
				"enterprise_router_id": utils.PathSearch("enterprise_router_id", v, nil),
			})
		}
		rst = append(rst, map[string]interface{}{"exclude_er_instances": erInstances})
	}
	return rst
}

func flattenGetCentralNetworkPolicyResponseBodyAssociateErInstanceDocument(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("document.er_instances", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"project_id":           utils.PathSearch("project_id", v, nil),
			"region_id":            utils.PathSearch("region_id", v, nil),
			"enterprise_router_id": utils.PathSearch("enterprise_router_id", v, nil),
		})
	}
	return rst
}

func resourceCentralNetworkPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteCentralNetworkPolicy: delete the central network policy
	var (
		deleteCentralNetworkPolicyHttpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/policies/{id}"
		deleteCentralNetworkPolicyProduct = "cc"
	)
	deleteCentralNetworkPolicyClient, err := cfg.NewServiceClient(deleteCentralNetworkPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	deleteCentralNetworkPolicyPath := deleteCentralNetworkPolicyClient.Endpoint + deleteCentralNetworkPolicyHttpUrl
	deleteCentralNetworkPolicyPath = strings.ReplaceAll(deleteCentralNetworkPolicyPath, "{domain_id}", cfg.DomainID)
	deleteCentralNetworkPolicyPath = strings.ReplaceAll(deleteCentralNetworkPolicyPath, "{id}", d.Id())
	deleteCentralNetworkPolicyPath = strings.ReplaceAll(deleteCentralNetworkPolicyPath, "{central_network_id}",
		d.Get("central_network_id").(string))

	deleteCentralNetworkPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err = deleteCentralNetworkPolicyClient.Request("DELETE", deleteCentralNetworkPolicyPath,
		&deleteCentralNetworkPolicyOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CentralNetworkPolicy")
	}

	return nil
}

func resourceCentralNetworkPolicyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <central_network_id>/<id>")
	}

	d.Set("central_network_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
