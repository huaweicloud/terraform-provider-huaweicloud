package identitycenter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var provisionPermissionSetNonUpdatableParams = []string{"instance_id", "permission_set_id", "account_id"}

// @API IdentityCenter POST /v1/instances/{instance_id}/permission-sets/{permission_set_id}/provision
// @API IdentityCenter GET /v1/instances/{instance_id}/permission-sets/provisioning-status/{request_id}
func ResourceProvisionPermissionSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProvisionPermissionSetCreate,
		UpdateContext: resourceProvisionPermissionSetUpdate,
		ReadContext:   resourceProvisionPermissionSetRead,
		DeleteContext: resourceProvisionPermissionSetDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceProvisionPermissionSetImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(provisionPermissionSetNonUpdatableParams),

		Description: "schema: Internal",
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
				Description: `Specifies the ID of the IAM Identity Center instance.`,
			},
			"permission_set_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the permission set ID of the IAM Identity Center.`,
			},
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the account ID.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The authorization status of a permission set.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceProvisionPermissionSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	instanceId := d.Get("instance_id").(string)
	permissionSetId := d.Get("permission_set_id").(string)
	// createIdentityCenterProvisionPermissionSet: create IdentityCenter provision permission set
	var (
		createProvisionPermissionSetHttpUrl = "v1/instances/{instance_id}/permission-sets/{permission_set_id}/provision"
		createProduct                       = "identitycenter"
	)
	client, err := cfg.NewServiceClient(createProduct, region)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter client: %s", err)
	}

	createProvisionPermissionSetPath := client.Endpoint + createProvisionPermissionSetHttpUrl
	createProvisionPermissionSetPath = strings.ReplaceAll(createProvisionPermissionSetPath, "{instance_id}", instanceId)
	createProvisionPermissionSetPath = strings.ReplaceAll(createProvisionPermissionSetPath, "{permission_set_id}", permissionSetId)

	createProvisionPermissionSetPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createProvisionPermissionSetPathOpt.JSONBody = map[string]interface{}{
		"target_type": "ACCOUNT",
		"target_id":   d.Get("account_id").(string),
	}

	resp, err := client.Request("POST", createProvisionPermissionSetPath, &createProvisionPermissionSetPathOpt)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter provision permission set: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error flattening IdentityCenter provision permission set: %s", err)
	}

	requestId := utils.PathSearch("permission_set_provisioning_status.request_id", respBody, "").(string)
	if requestId == "" {
		return diag.Errorf("unable to find the request ID from the API response")
	}
	d.SetId(requestId)

	err = checkProvisionPermissionSetStatus(ctx, client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceProvisionPermissionSetRead(ctx, d, meta)
}

func resourceProvisionPermissionSetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	instanceId := d.Get("instance_id").(string)

	client, err := cfg.NewServiceClient("identitycenter", region)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter client: %s", err)
	}

	resp, err := getProvisionPermissionSetStatus(client, instanceId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying IdentityCenter provision permission set")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("permission_set_id", utils.PathSearch("permission_set_provisioning_status.permission_set_id", resp, nil)),
		d.Set("account_id", utils.PathSearch("permission_set_provisioning_status.account_id", resp, nil)),
		d.Set("status", utils.PathSearch("permission_set_provisioning_status.status", resp, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceProvisionPermissionSetUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceProvisionPermissionSetDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting IdentityCenter provision permission set resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func checkProvisionPermissionSetStatus(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	instanceId := d.Get("instance_id").(string)
	timeout := d.Timeout(schema.TimeoutCreate)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      provisionPermissionSetStateRefreshFunc(client, instanceId, d.Id()),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for IdentityCenter provision permission set to be completed: %s", err)
	}
	return nil
}

func provisionPermissionSetStateRefreshFunc(client *golangsdk.ServiceClient, instanceId, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := getProvisionPermissionSetStatus(client, instanceId, id)
		if err != nil {
			return nil, "ERROR", err
		}

		status := utils.PathSearch("permission_set_provisioning_status.status", respBody, "").(string)
		if status == "SUCCEEDED" {
			return respBody, "COMPLETED", nil
		}

		if status == "FAILED" {
			return respBody, "ERROR", fmt.Errorf("failed to provision IdentityCenter permission set")
		}

		return respBody, "PENDING", nil
	}
}

func getProvisionPermissionSetStatus(client *golangsdk.ServiceClient, instanceId, id string) (interface{}, error) {
	getProvisionPermissionSetHttpUrl := "v1/instances/{instance_id}/permission-sets/provisioning-status/{request_id}"
	getProvisionPermissionSetPath := client.Endpoint + getProvisionPermissionSetHttpUrl
	getProvisionPermissionSetPath = strings.ReplaceAll(getProvisionPermissionSetPath, "{instance_id}", instanceId)
	getProvisionPermissionSetPath = strings.ReplaceAll(getProvisionPermissionSetPath, "{request_id}", id)

	getProvisionPermissionSetPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getProvisionPermissionSetResp, err := client.Request("GET", getProvisionPermissionSetPath, &getProvisionPermissionSetPathOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getProvisionPermissionSetResp)
}

func resourceProvisionPermissionSetImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format: the format must be <instance_id>/<request_id>")
		return nil, err
	}

	instanceID := parts[0]
	requestId := parts[1]

	d.Set("instance_id", instanceID)
	d.SetId(requestId)

	return []*schema.ResourceData{d}, nil
}
