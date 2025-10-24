package identitycenter

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IdentityCenter POST /v1/instances/{instance_id}/access-control-attribute-configuration
// @API IdentityCenter GET /v1/instances/{instance_id}/access-control-attribute-configuration
// @API IdentityCenter PUT /v1/instances/{instance_id}/access-control-attribute-configuration
// @API IdentityCenter DELETE /v1/instances/{instance_id}/access-control-attribute-configuration
func ResourceAccessControlAttributeConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccessControlAttributeConfigurationCreate,
		UpdateContext: resourceAccessControlAttributeConfigurationUpdate,
		ReadContext:   resourceAccessControlAttributeConfigurationRead,
		DeleteContext: resourceAccessControlAttributeConfigurationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the IAM Identity Center instance.`,
			},
			"access_control_attributes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the name of the attribute associated with the identity in your identity source.`,
						},
						"value": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the value used to map the specified attribute to the identity source.`,
						},
					},
				},
				Description: `Specifies the properties of ABAC configuration in IAM Identity Center instance.`,
			},
		},
	}
}

func resourceAccessControlAttributeConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	instanceId := d.Get("instance_id").(string)

	// createIdentityCenterAttrConfig: create IdentityCenter access control attribute configuration
	var (
		createAttrConfigHttpUrl = "v1/instances/{instance_id}/access-control-attribute-configuration"
		createProduct           = "identitycenter"
	)
	client, err := cfg.NewServiceClient(createProduct, region)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter client: %s", err)
	}

	createAttrConfigPath := client.Endpoint + createAttrConfigHttpUrl
	createAttrConfigPath = strings.ReplaceAll(createAttrConfigPath, "{instance_id}", instanceId)

	createAttrConfigPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createAttrConfigPathOpt.JSONBody = buildAccessControlAttributeConfigParams(d)
	_, err = client.Request("POST", createAttrConfigPath, &createAttrConfigPathOpt)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter access control attribute configuration: %s", err)
	}

	d.SetId(instanceId)

	return resourceAccessControlAttributeConfigurationRead(ctx, d, meta)
}

func buildAccessControlAttributeConfigParams(d *schema.ResourceData) map[string]interface{} {
	rawArray := d.Get("access_control_attributes").(*schema.Set).List()
	accessControlAttributes := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		raw := v.(map[string]interface{})
		accessControlAttributes[i] = map[string]interface{}{
			"key": raw["key"],
			"value": map[string]interface{}{
				"source": raw["value"],
			},
		}
	}

	bodyParams := map[string]interface{}{
		"instance_access_control_attribute_configuration": map[string]interface{}{
			"access_control_attributes": accessControlAttributes,
		},
	}

	return bodyParams
}

func resourceAccessControlAttributeConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// getIdentityCenterAttrConfig: Query Identity Center access control attribute configuration
	getAttrConfigClient, err := cfg.NewServiceClient("identitycenter", region)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter client: %s", err)
	}

	resp, err := GetAccessControlAttributeConfiguration(getAttrConfigClient, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying IdentityCenter access control attribute configuration")
	}

	status := utils.PathSearch("status", resp, "").(string)
	if status != "ENABLED" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("access_control_attributes", flattenAccessControlAttributes(
			utils.PathSearch("instance_access_control_attribute_configuration.access_control_attributes", resp, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAccessControlAttributes(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}
	curArray := resp.([]interface{})
	rst := make([]map[string]interface{}, len(curArray))
	for i, v := range curArray {
		rst[i] = map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value.source", v, nil),
		}
	}
	return rst
}

func resourceAccessControlAttributeConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateIdentityCenterAttrConfig: update IdentityCenter access control attribute configuration
	var (
		updateAttrConfigHttpUrl = "v1/instances/{instance_id}/access-control-attribute-configuration"
		updateProduct           = "identitycenter"
	)
	client, err := cfg.NewServiceClient(updateProduct, region)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter client: %s", err)
	}

	updateAttrConfigPath := client.Endpoint + updateAttrConfigHttpUrl
	updateAttrConfigPath = strings.ReplaceAll(updateAttrConfigPath, "{instance_id}", d.Id())

	updateAttrConfigPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateAttrConfigPathOpt.JSONBody = buildAccessControlAttributeConfigParams(d)
	_, err = client.Request("PUT", updateAttrConfigPath, &updateAttrConfigPathOpt)
	if err != nil {
		return diag.Errorf("error updating IdentityCenter access control attribute configuration: %s", err)
	}

	return resourceAccessControlAttributeConfigurationRead(ctx, d, meta)
}

func resourceAccessControlAttributeConfigurationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteIdentityCenterAttrConfig: Delete IdentityCenter access control attribute configuration
	var (
		deleteAttrConfigHttpUrl = "v1/instances/{instance_id}/access-control-attribute-configuration"
		deleteAttrConfigProduct = "identitycenter"
	)
	deleteAttrConfigClient, err := cfg.NewServiceClient(deleteAttrConfigProduct, region)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter client: %s", err)
	}

	deleteAttrConfigPath := deleteAttrConfigClient.Endpoint + deleteAttrConfigHttpUrl
	deleteAttrConfigPath = strings.ReplaceAll(deleteAttrConfigPath, "{instance_id}", d.Id())

	deleteAttrConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteAttrConfigClient.Request("DELETE", deleteAttrConfigPath, &deleteAttrConfigOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IdentityCenter access control attribute configuration")
	}

	return nil
}

func GetAccessControlAttributeConfiguration(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	getAttrConfigHttpUrl := "v1/instances/{instance_id}/access-control-attribute-configuration"
	getAttrConfigPath := client.Endpoint + getAttrConfigHttpUrl
	getAttrConfigPath = strings.ReplaceAll(getAttrConfigPath, "{instance_id}", id)

	getAttrConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAttrConfigResp, err := client.Request("GET", getAttrConfigPath, &getAttrConfigOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getAttrConfigResp)
}
