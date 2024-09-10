package meeting

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/meeting/v1/assignments"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API Meeting POST /v1/usg/dcs/corp/admin/delete
// @API Meeting GET /v1/usg/dcs/corp/admin/{account}
// @API Meeting POST /v1/usg/dcs/corp/admin
func ResourceAdminAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdminAssignmentCreate,
		ReadContext:   resourceAdminAssignmentRead,
		DeleteContext: resourceAdminAssignmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAdminAssignmentImportState,
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

			// Arguments
			"account": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAdminAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	token, err := NewMeetingToken(conf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	account := d.Get("account").(string)
	opt := assignments.CreateOpts{
		Account: d.Get("account").(string),
		// Authorization token.
		Token: token,
	}
	err = assignments.Create(NewMeetingV1Client(conf), opt)
	if err != nil {
		return diag.Errorf("error assign the administrator role to a cloud meeting user: %s", err)
	}

	d.SetId(account)
	return resourceAdminAssignmentRead(ctx, d, meta)
}

func parseAdministratorNotFoundError(respErr error) error {
	var apiErr assignments.ErrResponse
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

func resourceAdminAssignmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	token, err := NewMeetingToken(conf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	opt := assignments.GetOpts{
		Token:   token,
		Account: d.Id(),
	}
	resp, err := assignments.Get(NewMeetingV1Client(conf), opt)
	if err != nil {
		return common.CheckDeletedDiag(d, parseAdministratorNotFoundError(err),
			"error retrieving administrator information")
	}

	mErr := multierror.Append(nil,
		d.Set("account", resp.Account),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAdminAssignmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	token, err := NewMeetingToken(conf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	opt := assignments.DeleteOpts{
		Token:    token,
		Accounts: []string{d.Id()},
	}
	err = assignments.BatchDelete(NewMeetingV1Client(conf), opt)
	if err != nil {
		return diag.Errorf("unassigned admin role failed: %s", err)
	}
	return nil
}

func resourceAdminAssignmentImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData,
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
		)
	default:
		return nil, fmt.Errorf("the imported ID specifies an invalid format, must be " +
			"<id>/<account_name>/<account_password> or <id>/<app_id>/<app_key>/<corp_id>/<user_id>.")
	}
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
