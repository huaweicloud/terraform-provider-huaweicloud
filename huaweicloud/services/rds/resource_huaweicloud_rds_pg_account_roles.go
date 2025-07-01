package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var pgAccountRolesNonUpdatableParams = []string{"instance_id", "user"}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/db-user-role
// @API RDS GET /v3/{project_id}/instances
// @API RDS GET /v3/{project_id}/instances/{instance_id}/db_user/detail
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/db-user-role
func ResourcePgAccountRoles() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePgAccountRolesCreate,
		UpdateContext: resourcePgAccountRolesUpdate,
		ReadContext:   resourcePgAccountRolesRead,
		DeleteContext: resourcePgAccountRolesDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(pgAccountRolesNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
				Description: `Specifies the ID of the RDS PostgreSQL instance.`,
			},
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the username of the account.`,
			},
			"roles": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of roles.`,
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

func resourcePgAccountRolesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	roles := d.Get("roles").(*schema.Set).List()
	requestBody := buildUpdatePgAccountMemberOfBodyParams(d.Get("user").(string), roles)
	err = updateMemberOf(ctx, d, client, "POST", requestBody)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("instance_id").(string), d.Get("user").(string)))

	return resourcePgAccountRolesRead(ctx, d, meta)
}

func resourcePgAccountRolesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/db_user/detail?page=1&limit=100"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and user from resource id
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, must be <instance_id>/<name>")
	}
	instanceId := parts[0]
	accountName := parts[1]

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getPgAccountResp, err := pagination.ListAllItems(
		client,
		"page",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS PostgreSQL account roles")
	}

	respJson, err := json.Marshal(getPgAccountResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
	if err != nil {
		return diag.FromErr(err)
	}

	roles := utils.PathSearch(fmt.Sprintf("users[?name=='%s']|[0].memberof", accountName), respBody, nil)

	if roles == nil || len(roles.([]interface{})) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("user", accountName),
		d.Set("roles", roles),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePgAccountRolesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	oldRaws, newRaws := d.GetChange("roles")
	addRoles := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
	deleteRoles := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))

	if deleteRoles.Len() > 0 {
		requestBody := buildUpdatePgAccountMemberOfBodyParams(d.Get("user").(string), deleteRoles.List())
		err = updateMemberOf(ctx, d, client, "DELETE", requestBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if addRoles.Len() > 0 {
		requestBody := buildUpdatePgAccountMemberOfBodyParams(d.Get("user").(string), addRoles.List())
		err = updateMemberOf(ctx, d, client, "POST", requestBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourcePgAccountRolesRead(ctx, d, meta)
}

func resourcePgAccountRolesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	roles := d.Get("roles").(*schema.Set).List()
	requestBody := buildUpdatePgAccountMemberOfBodyParams(d.Get("user").(string), roles)
	err = updateMemberOf(ctx, d, client, "DELETE", requestBody)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
