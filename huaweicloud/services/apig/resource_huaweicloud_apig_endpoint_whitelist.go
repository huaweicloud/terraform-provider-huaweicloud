package apig

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/endpoints"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG POST /v2/{project_id}/apigw/instances/{instanceId}/vpc-endpoint/permissions/batch-add
// @API APIG POST /v2/{project_id}/apigw/instances/{instanceId}/vpc-endpoint/permissions/batch-delete
// @API APIG GET /v2/{project_id}/apigw/instances/{instanceId}/vpc-endpoint/permissions
func ResourceEndpointWhiteList() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointWhiteListCreate,
		ReadContext:   resourceEndpointWhiteListRead,
		UpdateContext: resourceEndpointWhiteListUpdate,
		DeleteContext: resourceEndpointWhiteListDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceEndpointWhiteListImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the endpoint service is located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the endpoint service belongs.",
			},
			"whitelists": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The whitelist records of the endpoint service.",
			},
		},
	}
}

func resourceEndpointWhiteListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		whitelists = d.Get("whitelists").(*schema.Set)
	)
	err = createEndpointWhiteListForApis(client, instanceId, utils.ExpandToStringListBySet(whitelists))
	if err != nil {
		return diag.Errorf("error creating whitelist records: %s", err)
	}
	d.SetId(instanceId)

	return resourceEndpointWhiteListRead(ctx, d, meta)
}

func createEndpointWhiteListForApis(client *golangsdk.ServiceClient, instanceId string, whitelists []string) error {
	opt := endpoints.BatchOpts{
		InstanceId:  instanceId,
		Permissions: whitelists,
	}
	_, err := endpoints.AddPermissions(client, opt)
	if err != nil {
		return err
	}

	return nil
}

func resourceEndpointWhiteListRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)

		opt = endpoints.ListOpts{
			InstanceId: instanceId,
		}
	)
	resp, err := endpoints.ListPermissions(client, opt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying whitelist record")
	}

	var (
		domainId          = cfg.DomainID
		initialPermission = "iam:domain::" + domainId
	)
	var whiteLists []string
	for _, endpointPermission := range resp {
		if endpointPermission.Permission != initialPermission {
			whiteLists = append(whiteLists, endpointPermission.Permission)
		}
	}

	// Check the whitelists list except for the current account.
	if len(whiteLists) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "whitelists of endpoint service")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("whitelists", whiteLists),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving whitelist fields for the endpoint service of the dedicated instance "+
			"(%s): %s", instanceId, err)
	}

	return nil
}

func resourceEndpointWhiteListUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

		oldRaw, newRaw = d.GetChange("whitelists")

		addSet = newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))
		rmSet  = oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))
	)

	if rmSet.Len() > 0 {
		whitelists := utils.ExpandToStringListBySet(rmSet)
		err := deleteEndpointWhiteListFromApis(client, instanceId, whitelists)
		if err != nil {
			return diag.Errorf("error deleting whitelist records: %s", err)
		}
	}
	if addSet.Len() > 0 {
		whitelists := utils.ExpandToStringListBySet(addSet)
		err = createEndpointWhiteListForApis(client, instanceId, whitelists)
		if err != nil {
			return diag.Errorf("error creating whitelist records: %s", err)
		}
	}

	return resourceEndpointWhiteListRead(ctx, d, meta)
}

func resourceEndpointWhiteListDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		whitelists = d.Get("whitelists").(*schema.Set)
	)
	err = deleteEndpointWhiteListFromApis(client, instanceId, utils.ExpandToStringListBySet(whitelists))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting whitelists for the endpoint service")
	}

	return nil
}

func deleteEndpointWhiteListFromApis(client *golangsdk.ServiceClient, instanceId string, whitelists []string) error {
	opt := endpoints.BatchOpts{
		InstanceId:  instanceId,
		Permissions: whitelists,
	}

	err := endpoints.DeletePermissions(client, opt)
	if err != nil {
		return err
	}

	return nil
}

func resourceEndpointWhiteListImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	return []*schema.ResourceData{d}, d.Set("instance_id", d.Id())
}
