package vpc

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPC POST /v3/{project_id}/vpc/address-groups
// @API VPC DELETE /v3/{project_id}/vpc/address-groups/{address_group_id}/force
// @API VPC GET /v3/{project_id}/vpc/address-groups/{address_group_id}
// @API VPC PUT /v3/{project_id}/vpc/address-groups/{address_group_id}
// @API VPC DELETE /v3/{project_id}/vpc/address-groups/{address_group_id}
func ResourceVpcAddressGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcAddressGroupCreate,
		UpdateContext: resourceVpcAddressGroupUpdate,
		DeleteContext: resourceVpcAddressGroupDelete,
		ReadContext:   resourceVpcAddressGroupRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"addresses": {
				// the addresses will be sorted by cloud
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      4,
				ValidateFunc: validation.IntInSlice([]int{4, 6}),
			},
			"ip_extra_set": {
				// the addresses will be sorted by cloud
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"remarks": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_capacity": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceVpcAddressGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	createAddressGroupHttpUrl := "v3/{project_id}/vpc/address-groups"
	createAddressGroupPath := client.Endpoint + createAddressGroupHttpUrl
	createAddressGroupPath = strings.ReplaceAll(createAddressGroupPath, "{project_id}", client.ProjectID)

	createAddressGroupPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	addressGroupBody := buildCreateAddressGroupBodyParams(d, cfg)
	createAddressGroupPathOpt.JSONBody = utils.RemoveNil(addressGroupBody)

	log.Printf("[DEBUG] Create VPC address group options: %#v", addressGroupBody)
	response, err := client.Request("POST", createAddressGroupPath, &createAddressGroupPathOpt)
	if err != nil {
		return diag.Errorf("error creating VPC address group: %s", err)
	}

	responseBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("address_group.id", responseBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating VPC address group: ID is not found in API response")
	}

	d.SetId(id)
	return resourceVpcAddressGroupRead(ctx, d, meta)
}

func buildCreateAddressGroupBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"address_group": map[string]interface{}{
			"name":                  d.Get("name"),
			"ip_set":                utils.ValueIgnoreEmpty(d.Get("addresses").(*schema.Set).List()),
			"ip_version":            d.Get("ip_version"),
			"description":           utils.ValueIgnoreEmpty(d.Get("description")),
			"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
			"max_capacity":          utils.ValueIgnoreEmpty(d.Get("max_capacity").(int)),
			"ip_extra_set":          buildIpExtraSet(d),
		},
	}

	return bodyParams
}

func buildIpExtraSet(d *schema.ResourceData) []map[string]interface{} {
	ipExtraSet := d.Get("ip_extra_set").(*schema.Set).List()
	if len(ipExtraSet) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(ipExtraSet))
	for i, v := range ipExtraSet {
		res[i] = map[string]interface{}{
			"ip":      utils.PathSearch("ip", v, nil),
			"remarks": utils.PathSearch("remarks", v, nil),
		}
	}

	return res
}

func resourceVpcAddressGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	getAddressGroupHttpUrl := "v3/{project_id}/vpc/address-groups/{address_group_id}"
	getAddressGroupPath := client.Endpoint + getAddressGroupHttpUrl
	getAddressGroupPath = strings.ReplaceAll(getAddressGroupPath, "{project_id}", client.ProjectID)
	getAddressGroupPath = strings.ReplaceAll(getAddressGroupPath, "{address_group_id}", d.Id())

	getAddressGroupPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	response, err := client.Request("GET", getAddressGroupPath, &getAddressGroupPathOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error fetching VPC address group")
	}

	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("address_group.name", respBody, nil)),
		d.Set("description", utils.PathSearch("address_group.description", respBody, nil)),
		d.Set("addresses", utils.PathSearch("address_group.ip_set", respBody, nil)),
		d.Set("ip_version", utils.PathSearch("address_group.ip_version", respBody, nil)),
		d.Set("max_capacity", utils.PathSearch("address_group.max_capacity", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("address_group.enterprise_project_id", respBody, nil)),
		d.Set("ip_extra_set", flattenIpExtraSet(
			utils.PathSearch("address_group.ip_extra_set", respBody, []interface{}{}).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenIpExtraSet(ipExtraSet []interface{}) []map[string]interface{} {
	if len(ipExtraSet) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(ipExtraSet))

	for i, v := range ipExtraSet {
		res[i] = map[string]interface{}{
			"ip":      utils.PathSearch("ip", v, nil),
			"remarks": utils.PathSearch("remarks", v, nil),
		}
	}

	return res
}

func resourceVpcAddressGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	updateAddressGroupHttpUrl := "v3/{project_id}/vpc/address-groups/{address_group_id}"
	updateAddressGroupPath := client.Endpoint + updateAddressGroupHttpUrl
	updateAddressGroupPath = strings.ReplaceAll(updateAddressGroupPath, "{project_id}", client.ProjectID)
	updateAddressGroupPath = strings.ReplaceAll(updateAddressGroupPath, "{address_group_id}", d.Id())

	addressGroupBody := buildUpdateAddressGroupBodyParams(d)

	updateAddressGroupPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(addressGroupBody),
	}

	log.Printf("[DEBUG] Update VPC address group options: %#v", addressGroupBody)
	_, err = client.Request("PUT", updateAddressGroupPath, &updateAddressGroupPathOpt)
	if err != nil {
		return diag.Errorf("error updating VPC address group: %s", err)
	}

	return resourceVpcAddressGroupRead(ctx, d, meta)
}

func buildUpdateAddressGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	addressGroup := map[string]interface{}{}

	if d.HasChange("name") {
		addressGroup["name"] = d.Get("name")
	}

	if d.HasChange("description") {
		addressGroup["description"] = d.Get("description")
	}

	if d.HasChange("addresses") {
		addressGroup["ip_set"] = d.Get("addresses").(*schema.Set).List()
	}

	if d.HasChange("ip_extra_set") {
		addressGroup["ip_extra_set"] = buildIpExtraSet(d)
	}

	if d.HasChange("max_capacity") {
		addressGroup["max_capacity"] = d.Get("max_capacity").(int)
	}

	bodyParam := map[string]interface{}{
		"address_group": addressGroup,
	}

	return bodyParam
}

func resourceVpcAddressGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	var deleteAddressGroupHttpUrl string

	if !d.Get("force_destroy").(bool) {
		deleteAddressGroupHttpUrl = "v3/{project_id}/vpc/address-groups/{address_group_id}"
	} else {
		deleteAddressGroupHttpUrl = "v3/{project_id}/vpc/address-groups/{address_group_id}/force"
	}

	deleteAddressGroupPath := client.Endpoint + deleteAddressGroupHttpUrl
	deleteAddressGroupPath = strings.ReplaceAll(deleteAddressGroupPath, "{project_id}", client.ProjectID)
	deleteAddressGroupPath = strings.ReplaceAll(deleteAddressGroupPath, "{address_group_id}", d.Id())

	deleteAddressGroupPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deleteAddressGroupPath, &deleteAddressGroupPathOpt)

	if err != nil {
		return diag.Errorf("error deleting VPC address group: %s", err)
	}

	return nil
}
