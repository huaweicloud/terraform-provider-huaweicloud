package dli

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// The error code corresponding to errCodeDbNotFound is an important sign that the related resource does not exist.
// When the object is database and the database does not exist, the API return this error:
// + {"error_code": "DLI.0002", "error_msg": "There is no database named xxx"}
// When the object is data table and the data table (even containing the database) does not exist, the API return this error:
// + {"error_code": "DLI.0002", "error_msg": "There is no table named named xxx"}
const errCodeDbNotFound string = "DLI.0002"

// @API DLI PUT /v1.0/{project_id}/authorization
// @API DLI GET /v1.0/{project_id}/authorization/privileges
func ResourceDatabasePrivilege() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabasePrivilegeCreate,
		ReadContext:   resourceDatabasePrivilegeRead,
		UpdateContext: resourceDatabasePrivilegeUpdate,
		DeleteContext: resourceDatabasePrivilegeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDatabasePrivilegeImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region in which to create the resource.`,
			},
			"object": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The authorization object definition.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the authorized (IAM) user.`,
			},
			"privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of permissions granted to the database or data table.`,
				DiffSuppressFunc: func(_, _, _ string, d *schema.ResourceData) bool {
					if _, ok := d.GetOk("privileges"); !ok {
						// If the value of the permission is nil, the Computed behavior will prevent changes.
						// If the value of the permission is not nil and the length of the array is zero, need a
						// DiffSuppress function to prevent changes.
						return true
					}
					return false
				},
			},
		},
	}
}

func buildModifyDatabasePrivilegeBodyParams(d *schema.ResourceData, action string) map[string]interface{} {
	// Default permission for database and data table.
	privileges := []string{"SELECT"}
	if d.Get("privileges").(*schema.Set).Len() > 0 {
		privileges = utils.ExpandToStringListBySet(d.Get("privileges").(*schema.Set))
	}

	return map[string]interface{}{
		"user_name": utils.ValueIgnoreEmpty(d.Get("user_name")),
		"action":    action,
		"privileges": []map[string]interface{}{
			{
				"object":     d.Get("object"),
				"privileges": privileges,
			},
		},
	}
}

func modifyDatabasePrivilege(client *golangsdk.ServiceClient, d *schema.ResourceData, action string) error {
	httpUrl := "v1.0/{project_id}/authorization"

	modifyPath := client.Endpoint + httpUrl
	modifyPath = strings.ReplaceAll(modifyPath, "{project_id}", client.ProjectID)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildModifyDatabasePrivilegeBodyParams(d, action)),
	}

	requestResp, err := client.Request("PUT", modifyPath, &opts)
	if err != nil {
		return err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return err
	}
	if !utils.PathSearch("is_success", respBody, true).(bool) {
		return fmt.Errorf("unable to %s the privileges: %s", action,
			utils.PathSearch("message", respBody, "Message Not Found"))
	}
	return nil
}

func resourceDatabasePrivilegeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	err = modifyDatabasePrivilege(client, d, "grant")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%v/%v", d.Get("object"), d.Get("user_name")))

	return resourceDatabasePrivilegeRead(ctx, d, meta)
}

func GetObjectPrivilegesForSpecifiedUser(client *golangsdk.ServiceClient, object, userName string) (objectResp, privilege interface{}, err error) {
	var (
		filterExpression string
		httpUrl          = "v1.0/{project_id}/authorization/privileges?object={object}"
	)
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{object}", object)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var requestResp *http.Response
	requestResp, err = client.Request("GET", getPath, &getOpts)
	if err != nil {
		return
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return
	}

	if !utils.PathSearch("is_success", respBody, true).(bool) {
		err = golangsdk.ErrDefault500{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1.0/{project_id}/authorization/privileges",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("unable to query the privileges: %s", utils.PathSearch("message", respBody, "Message Not Found"))),
			},
		}
		return
	}
	filterExpression = fmt.Sprintf("privileges[?user_name=='%s']|[0]", userName)
	privilege = utils.PathSearch(filterExpression, respBody, nil)
	if privilege != nil {
		objectResp = utils.PathSearch("object", privilege, nil)
		if objectResp == nil && utils.PathSearch("object_type", respBody, "type_not_found").(string) == "database" {
			// The object value will not be included in the structure 'privileges' if the grant object type is database,
			// but the strings that make up the object will be returned in the upper structure, which are 'object_type'
			// and 'object_name'.
			objectResp = fmt.Sprintf("databases.%v", utils.PathSearch("object_name", respBody, "object_not_found"))
		}
	} else {
		err = golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1.0/{project_id}/authorization/privileges",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the privileges not found for the object (%s) and user (%s)", object, userName)),
			},
		}
	}
	return
}

func resourceDatabasePrivilegeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		object   = d.Get("object").(string)
		userName = d.Get("user_name").(string)
	)
	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	objectResp, privilege, err := GetObjectPrivilegesForSpecifiedUser(client, object, userName)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", errCodeDbNotFound),
			"error retrieving privileges")
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("object", objectResp),
		d.Set("user_name", utils.PathSearch("user_name", privilege, nil)),
		d.Set("privileges", utils.PathSearch("privileges", privilege, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDatabasePrivilegeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	err = modifyDatabasePrivilege(client, d, "update")
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDatabasePrivilegeRead(ctx, d, meta)
}

func resourceDatabasePrivilegeDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	err = modifyDatabasePrivilege(client, d, "revoke")
	return diag.FromErr(err)
}

func resourceDatabasePrivilegeImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	var (
		importId = d.Id()
		parts    = strings.Split(importId, "/")
	)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid resource ID format for privilege management, want '<object>/<user_name>', "+
			"but got '%s'", importId)
	}

	mErr := multierror.Append(nil,
		d.Set("object", parts[0]),
		d.Set("user_name", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
