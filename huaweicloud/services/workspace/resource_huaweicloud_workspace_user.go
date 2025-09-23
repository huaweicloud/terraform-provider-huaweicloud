package workspace

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/workspace/v2/users"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	RFC3339NoT      = "2006-01-02T15:04:05Z"
	MilliRFC3339NoT = "2006-01-02T15:04:05.000Z"
)

// @API Workspace DELETE /v2/{project_id}/users/{user_id}
// @API Workspace GET /v2/{project_id}/users/{user_id}
// @API Workspace PUT /v2/{project_id}/users/{user_id}
// @API Workspace POST /v2/{project_id}/users
func ResourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,

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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"active_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The activation mode of the user.`,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The phone number of the user.`,
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The initial passowrd of user.`,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_expires": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},
			"password_never_expires": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_change_password": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"next_login_change_password": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"total_desktops": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func calculateExpireTime(timeStr string) (string, error) {
	if timeStr == "0" || timeStr == "" {
		return "0", nil
	}
	return translateUTC0Time(timeStr)
}

func translateUTC0Time(timeStr string) (string, error) {
	timestamp, err := time.Parse("2006-01-02T15:04:05Z", timeStr)
	if err != nil {
		return "", err
	}

	return utils.FormatTimeStampRFC3339(timestamp.Unix()-int64(utils.GetTimezoneCode()*3600), true, MilliRFC3339NoT), nil
}

func buildUserCreateOpts(d *schema.ResourceData) (users.CreateOpts, error) {
	result := users.CreateOpts{
		Name:                    d.Get("name").(string),
		ActiveType:              d.Get("active_type").(string),
		Email:                   d.Get("email").(string),
		Phone:                   d.Get("phone").(string),
		Password:                d.Get("password").(string),
		Description:             d.Get("description").(string),
		EnableChangePassword:    utils.Bool(d.Get("enable_change_password").(bool)),
		NextLoginChangePassword: utils.Bool(d.Get("next_login_change_password").(bool)),
	}

	expireTime, err := calculateExpireTime(d.Get("account_expires").(string))
	if err != nil {
		return result, err
	}
	result.AccountExpires = expireTime
	log.Printf("[DEBUG] The createOpts of Workspace user is: %#v", result)

	return result, nil
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	opts, err := buildUserCreateOpts(d)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := users.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating Workspace user: %s", err)
	}

	d.SetId(resp.ID)

	return resourceUserUpdate(ctx, d, meta)
}

func parseUserAccountExpires(expires int) string {
	if expires == 0 {
		return strconv.Itoa(expires)
	}
	return utils.FormatTimeStampRFC3339(int64(expires/1000), false, RFC3339NoT)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.WorkspaceV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	resp, err := users.Get(client, d.Id())
	if err != nil {
		// WKS.00170312: The tanant(Workspace service) not exist.
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "WKS.00170312"), "Workspace user")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("active_type", resp.ActiveType),
		d.Set("email", resp.Email),
		d.Set("phone", resp.Phone),
		d.Set("description", resp.Description),
		d.Set("account_expires", parseUserAccountExpires(resp.AccountExpires)),
		d.Set("enable_change_password", resp.EnableChangePassword),
		d.Set("next_login_change_password", resp.NextLoginChangePassword),
		d.Set("password_never_expires", resp.PasswordNeverExpires),
		d.Set("disabled", resp.Disabled),
		d.Set("locked", resp.Locked),
		d.Set("total_desktops", resp.TotalDesktops),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func buildUserUpdateOpts(d *schema.ResourceData) (users.UpdateOpts, error) {
	result := users.UpdateOpts{
		ActiveType:              d.Get("active_type").(string),
		Email:                   utils.String(d.Get("email").(string)),
		Phone:                   utils.String(d.Get("phone").(string)),
		Description:             utils.String(d.Get("description").(string)),
		EnableChangePassword:    utils.Bool(d.Get("enable_change_password").(bool)),
		NextLoginChangePassword: utils.Bool(d.Get("next_login_change_password").(bool)),
		PasswordNeverExpires:    utils.Bool(d.Get("password_never_expires").(bool)),
		Disabled:                utils.Bool(d.Get("disabled").(bool)),
	}

	expireTime, err := calculateExpireTime(d.Get("account_expires").(string))
	if err != nil {
		return result, err
	}
	result.AccountExpires = expireTime
	log.Printf("[DEBUG] The updateOpts of Workspace user is: %#v", result)

	return result, nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	opts, err := buildUserUpdateOpts(d)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = users.Update(client, d.Id(), opts)
	if err != nil {
		return diag.Errorf("error updating Workspace user (%s): %s", d.Id(), err)
	}
	return resourceUserRead(ctx, d, meta)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	userId := d.Id()
	err = users.Delete(client, userId)
	if err != nil {
		// WKS.00170312: The user has been deleted.
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "WKS.00170312"),
			fmt.Sprintf("error deleting Workspace user (%s)", d.Id()))
	}

	return nil
}
