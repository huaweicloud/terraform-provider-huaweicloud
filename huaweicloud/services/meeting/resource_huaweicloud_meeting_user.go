package meeting

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/meeting/v1/assignments"
	"github.com/chnsz/golangsdk/openstack/meeting/v1/users"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type UserStatus int

var (
	userStatusNormal  UserStatus = 0
	userStatusDisable UserStatus = 1
)

// @API Meeting POST /v1/usg/dcs/corp/admin/delete
// @API Meeting POST /v1/usg/dcs/corp/admin
// @API Meeting POST /v1/usg/dcs/corp/member/delete
// @API Meeting GET /v1/usg/dcs/corp/member/{account}
// @API Meeting PUT /v1/usg/dcs/corp/member/{account}
// @API Meeting POST /v1/usg/dcs/corp/member
// @API Meeting POST /v1/usg/acs/token/validate
// @API Meeting POST /v1/usg/acs/auth/account
// @API Meeting POST /v2/usg/acs/auth/appauth
func ResourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceUserManagementImportState,
		},

		Schema: map[string]*schema.Schema{
			// Authorization arguments
			"account_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"account_password"},
			},
			"account_password": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"app_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"app_key"},
				ExactlyOneOf: []string{"account_name"},
			},
			"app_key": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"corp_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			// Arguments
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"account": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"third_account": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"app_id"},
			},
			"department_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"english_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"country": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"phone": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"hide_phone": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_send_notify": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"signature": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.IntInSlice([]int{
					int(userStatusNormal), int(userStatusDisable),
				}),
			},
			"title": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sort_level": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"is_admin": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"sip_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"department_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"department_name_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func isSendNotify(send bool) string {
	// Send without filling in or other values, and send by default
	if send {
		return "1"
	}
	return "0" // The zero means do not send.
}

func buildUserCreateOpts(d *schema.ResourceData, token string) users.CreateOpts {
	return users.CreateOpts{
		Account:      d.Get("account").(string),
		ThirdAccount: d.Get("third_account").(string),
		Name:         d.Get("name").(string),
		Password:     d.Get("password").(string),
		Country:      d.Get("country").(string),
		DeptCode:     d.Get("department_code").(string),
		Description:  d.Get("description").(string),
		Email:        d.Get("email").(string),
		EnglishName:  d.Get("english_name").(string),
		Phone:        d.Get("phone").(string),
		HidePhone:    utils.Bool(d.Get("hide_phone").(bool)),
		SendNotify:   isSendNotify(d.Get("is_send_notify").(bool)),
		Signature:    d.Get("signature").(string),
		SortLevel:    d.Get("sort_level").(int),
		Status:       utils.Int(d.Get("status").(int)),
		Title:        d.Get("title").(string),
		// Authorization token.
		Token: token,
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	token, err := NewMeetingToken(conf, d)
	if err != nil {
		return diag.FromErr(err)
	}
	client := NewMeetingV1Client(conf)

	resp, err := users.Create(NewMeetingV1Client(conf), buildUserCreateOpts(d, token))
	if err != nil {
		return diag.Errorf("error creating cloud meeting user: %s", err)
	}

	d.SetId(resp.UserAccount)
	if d.Get("is_admin").(bool) {
		opt := assignments.CreateOpts{
			Account: d.Id(),
			// Authorization token.
			Token: token,
		}
		err = assignments.Create(client, opt)
		if err != nil {
			return diag.Errorf("error creating cloud meeting user: %s", err)
		}
	}
	return resourceUserRead(ctx, d, meta)
}

func parseUserNotFoundError(respErr error) error {
	var apiErr users.ErrResponse
	if errCode, ok := respErr.(golangsdk.ErrDefault400); ok && errCode.Body != nil {
		pErr := json.Unmarshal(errCode.Body, &apiErr)
		if pErr == nil && apiErr.Code == "USG.201040000" {
			return golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Body: []byte("the user does not exist"),
				},
			}
		}
	}
	return respErr
}

func resourceUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	token, err := NewMeetingToken(conf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	opt := users.GetOpts{
		Token:   token,
		Account: d.Id(),
	}
	resp, err := users.Get(NewMeetingV1Client(conf), opt)
	if err != nil {
		return common.CheckDeletedDiag(d, parseUserNotFoundError(err), "error retrieving cloud meeting user")
	}

	mErr := multierror.Append(nil,
		d.Set("account", resp.UserAccount),
		d.Set("name", resp.Name),
		d.Set("third_account", resp.ThirdAccount),
		d.Set("country", resp.Country),
		d.Set("department_code", resp.DeptCode),
		d.Set("description", resp.Description),
		d.Set("email", resp.Email),
		d.Set("english_name", resp.EnglishName),
		d.Set("phone", resp.Phone),
		d.Set("hide_phone", resp.HidePhone),
		d.Set("signature", resp.Signature),
		d.Set("status", resp.Status),
		d.Set("title", resp.Title),
		d.Set("sort_level", resp.SortLevel),
		// Computed parameters.
		d.Set("sip_number", resp.SipNum),
		d.Set("type", resp.UserType),
		d.Set("department_name", resp.DeptName),
		d.Set("department_name_path", resp.DeptNamePath),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUserUpdateOpts(d *schema.ResourceData, token string) users.UpdateOpts {
	result := users.UpdateOpts{
		Account:     d.Id(),
		Name:        d.Get("name").(string),
		Country:     d.Get("country").(string),
		Phone:       d.Get("phone").(string),
		DeptCode:    utils.String(d.Get("department_code").(string)),
		Description: utils.String(d.Get("description").(string)),
		Email:       utils.String(d.Get("email").(string)),
		EnglishName: utils.String(d.Get("english_name").(string)),
		HidePhone:   utils.Bool(d.Get("hide_phone").(bool)),
		Signature:   utils.String(d.Get("signature").(string)),
		SortLevel:   d.Get("sort_level").(int),
		Status:      utils.Int(d.Get("status").(int)),
		Title:       utils.String(d.Get("title").(string)),
		// Authorization token.
		Token: token,
	}

	return result
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	token, err := NewMeetingToken(conf, d)
	if err != nil {
		return diag.FromErr(err)
	}
	client := NewMeetingV1Client(conf)

	if d.HasChangeExcept("is_admin") {
		opt := buildUserUpdateOpts(d, token)
		_, err = users.Update(client, opt)
		if err != nil {
			return diag.Errorf("error updating cloud meeting user: %s", err)
		}
	}
	if d.HasChange("is_admin") {
		opt := assignments.DeleteOpts{
			Token:    token,
			Accounts: []string{d.Id()},
		}
		err = assignments.BatchDelete(client, opt)
		if err != nil {
			return diag.Errorf("error deleting cloud meeting user: %s", err)
		}
	}

	return resourceUserRead(ctx, d, meta)
}

func resourceUserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	token, err := NewMeetingToken(conf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	opt := users.DeleteOpts{
		Token: token,
	}
	err = users.BatchDelete(NewMeetingV1Client(conf), opt, []string{d.Id()})
	if err != nil {
		return diag.Errorf("error deleting cloud meeting user: %s", err)
	}
	return nil
}

func resourceUserManagementImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData,
	error) {
	var mErr *multierror.Error
	parts := strings.Split(d.Id(), "/")
	switch len(parts) {
	case 3:
		d.SetId(parts[0])
		mErr = multierror.Append(mErr,
			d.Set("account_name", parts[1]),
			d.Set("account_password", parts[2]),
		)
	case 5:
		d.SetId(parts[0])
		mErr = multierror.Append(mErr,
			d.Set("app_id", parts[1]),
			d.Set("app_key", parts[2]),
			d.Set("corp_id", parts[3]),
			d.Set("user_id", parts[4]),
		)
	default:
		return nil, fmt.Errorf("the imported ID specifies an invalid format, must be " +
			"<id>/<account_name>/<account_password> or <id>/<app_id>/<app_key>/<corp_id>/<user_id>.")
	}
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
