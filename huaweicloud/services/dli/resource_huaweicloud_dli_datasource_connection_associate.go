package dli

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

// @API DLI POST /v2.0/{project_id}/datasource/enhanced-connections/{connection_id}/associate-queue
// @API DLI GET /v2.0/{project_id}/datasource/enhanced-connections/{connection_id}
// @API DLI POST /v2.0/{project_id}/datasource/enhanced-connections/{connection_id}/disassociate-queue
func ResourceDatasourceConnectionAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatasourceConnectionAssociateCreate,
		ReadContext:   resourceDatasourceConnectionAssociateRead,
		UpdateContext: resourceDatasourceConnectionAssociateUpdate,
		DeleteContext: resourceDatasourceConnectionAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the datasource enhanced connection and elastic resource pools are located.`,
			},
			"connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the datasource enhanced connection to be associated.`,
			},
			"elastic_resource_pools": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of the elastic resource pool names.`,
			},
		},
	}
}

func buildCreateDatasourceConnectionAssociateBodyParams(poolNames *schema.Set) map[string]interface{} {
	return map[string]interface{}{
		"elastic_resource_pools": utils.ExpandToStringListBySet(poolNames),
	}
}

func associateResourcePoolsToConnection(client *golangsdk.ServiceClient, connectionId string, poolNames *schema.Set) error {
	httpUrl := "v2.0/{project_id}/datasource/enhanced-connections/{connection_id}/associate-queue"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{connection_id}", connectionId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDatasourceConnectionAssociateBodyParams(poolNames)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return fmt.Errorf("error associating the elastic resource pools to the enhanced connection: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return err
	}
	if !utils.PathSearch("is_success", respBody, true).(bool) {
		return fmt.Errorf("unable to associate the elastic resource pools to the enhanced connection: %s",
			utils.PathSearch("message", respBody, "Message Not Found"))
	}
	return nil
}

func resourceDatasourceConnectionAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		connectionId = d.Get("connection_id").(string)
	)

	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	err = associateResourcePoolsToConnection(client, connectionId, d.Get("elastic_resource_pools").(*schema.Set))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(connectionId)

	return resourceDatasourceConnectionAssociateRead(ctx, d, meta)
}

func GetDatasourceConnectionAssociatedPoolNames(client *golangsdk.ServiceClient, connectionId string) ([]interface{}, error) {
	respBody, err := getConnectionById(client, connectionId)
	if err != nil {
		return nil, err
	}

	associatedPoolNames := utils.PathSearch("elastic_resource_pools[*].name", respBody, make([]interface{}, 0)).([]interface{})
	if len(associatedPoolNames) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return associatedPoolNames, nil
}

func resourceDatasourceConnectionAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		connectionId = d.Id()
	)

	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	associatedPoolNames, err := GetDatasourceConnectionAssociatedPoolNames(client, connectionId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DLI Enhanced connection associate elastic resource pools")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("connection_id", connectionId),
		d.Set("elastic_resource_pools", associatedPoolNames),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCreateDatasourceConnectionDisassociateBodyParams(poolNames *schema.Set) map[string]interface{} {
	return map[string]interface{}{
		"elastic_resource_pools": utils.ExpandToStringListBySet(poolNames),
	}
}

func disassociateResourcePoolsToConnection(client *golangsdk.ServiceClient, connectionId string, poolNames *schema.Set) error {
	httpUrl := "v2.0/{project_id}/datasource/enhanced-connections/{connection_id}/disassociate-queue"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{connection_id}", connectionId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDatasourceConnectionDisassociateBodyParams(poolNames)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return fmt.Errorf("error disassociating the elastic resource pools from the enhanced connection: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return err
	}
	if !utils.PathSearch("is_success", respBody, true).(bool) {
		return fmt.Errorf("unable to disassociate the elastic resource pools from the enhanced connection: %s",
			utils.PathSearch("message", respBody, "Message Not Found"))
	}
	return nil
}

func resourceDatasourceConnectionAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		connectionId = d.Get("connection_id").(string)
	)

	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	oldRaws, newRaws := d.GetChange("elastic_resource_pools")
	associatedPoolNames := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
	disassociatedPoolNames := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))

	if disassociatedPoolNames.Len() > 0 {
		err = disassociateResourcePoolsToConnection(client, connectionId, disassociatedPoolNames)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if associatedPoolNames.Len() > 0 {
		err = associateResourcePoolsToConnection(client, connectionId, associatedPoolNames)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDatasourceConnectionAssociateRead(ctx, d, meta)
}

func resourceDatasourceConnectionAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		connectionId = d.Id()
	)

	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	disassociatedPoolNames, _ := d.GetChange("elastic_resource_pools")
	err = disassociateResourcePoolsToConnection(client, connectionId, disassociatedPoolNames.(*schema.Set))
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
