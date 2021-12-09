package dli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dli/v1/auth"
	"github.com/chnsz/golangsdk/openstack/dli/v1/flinkjob"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

const updateAction = "update"

func ResourceDliPermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDliPermissionCreate,
		ReadContext:   resourceDliPermissionRead,
		DeleteContext: resourceDliPermissionDelete,
		UpdateContext: resourceDliPermissionUpdate,
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
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"object": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"privileges": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_admin": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceDliPermissionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.DliV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud DLI v1 client: %s", err)
	}

	obj := d.Get("object").(string)
	userName := d.Get("user_name").(string)

	if strings.HasPrefix(obj, "queues.") {
		queueName := strings.Replace(obj, "queues.", "", 1)
		opts := auth.GrantQueuePermissionOpts{
			QueueName: queueName,
			UserName:  userName,
			Action:    updateAction,
		}

		opts.Privileges = utils.ExpandToStringList(d.Get("privileges").([]interface{}))

		rst, createErr := auth.GrantQueuePermission(client, opts)
		if createErr != nil {
			return fmtp.DiagErrorf("Error granting permission in DLI: %s", createErr)
		}

		if rst != nil && !rst.IsSuccess {
			return fmtp.DiagErrorf("Error granting permission in DLI: %s", rst.Message)
		}

	} else {
		opts := auth.GrantDataPermissionOpts{
			UserName: userName,
			Action:   updateAction,
		}

		ids := utils.ExpandToStringList(d.Get("privileges").([]interface{}))
		opts.Privileges = append(opts.Privileges, auth.DataPermission{
			Object:     obj,
			Privileges: ids,
		})

		rst, createErr := auth.GrantDataPermission(client, opts)
		if createErr != nil {
			return fmtp.DiagErrorf("Error granting permission in DLI: %s", createErr)
		}

		if rst != nil && !rst.IsSuccess {
			return fmtp.DiagErrorf("Error granting permission in DLI: %s", rst.Message)
		}
	}

	d.SetId(fmt.Sprintf("%s/%s", obj, userName))

	return resourceDliPermissionRead(ctx, d, meta)
}

func resourceDliPermissionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DliV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud DLI v1 client: %s", err)
	}

	obj, userName := ParseAuthInfoFromId(d.Id())

	permission, pErr := QueryPermission(client, obj, userName)
	if pErr != nil {
		return common.CheckDeletedDiag(d, err, "DLI")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("user_name", permission.UserName),
		d.Set("object", obj),
		d.Set("privileges", permission.Privileges),
		d.Set("is_admin", permission.IsAdmin),
	)
	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("An error occurred during resource parameter setting for DLI permission: %s", mErr)
	}
	return nil
}

func resourceDliPermissionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceDliPermissionCreate(ctx, d, meta)
}

func checkPrefixMatchDataPermession(obj string) bool {
	if strings.HasPrefix(obj, "jobs.flink.") || strings.HasPrefix(obj, "groups.") || strings.HasPrefix(obj, "groups.") {
		return true
	}
	return false
}

func queryDatabaseRelatePermission(client *golangsdk.ServiceClient, obj, userName string) (*auth.Privilege, error) {
	objArray := strings.Split(obj, ".")
	if len(objArray) == 2 {
		rst, err := auth.ListDatabasePermission(client, objArray[1])
		if err != nil {
			return nil, parseDliErrorToError404(err)
		}
		if rst != nil && !rst.IsSuccess {
			return nil, fmtp.Errorf("Error query DLI permission of database: %s", err)
		}

		for _, v := range rst.Privileges {
			if v.UserName == userName {
				return &v, nil
			}
		}
	} else if len(objArray) == 4 || len(objArray) == 6 {
		rst, err := auth.ListTablePermission(client, objArray[1], objArray[3])
		if err != nil {
			return nil, parseDliErrorToError404(err)
		}
		if rst != nil && !rst.IsSuccess {
			return nil, fmtp.Errorf("Error query DLI permission of table: %s", err)
		}
		for _, v := range rst.Privileges {
			if v.Object == obj && v.UserName == userName {
				privilege := auth.Privilege{
					IsAdmin:    v.IsAdmin,
					Privileges: v.Privileges,
					UserName:   v.UserName,
				}
				return &privilege, nil
			}
		}
	}

	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte("Error query DLI permission"),
		},
	}
}

func queryDataPermission(client *golangsdk.ServiceClient, obj, userName string) (*auth.Privilege, error) {
	rst, err := auth.ListDataPermission(client, auth.ListDataPermissionOpts{Object: obj})
	if err != nil {
		return nil, parseDliErrorToError404(err)
	}

	if rst != nil && !rst.IsSuccess {
		return nil, fmtp.Errorf("Error query DLI permission: %s", rst.Message)
	}

	for _, v := range rst.Privileges {
		if v.UserName == userName {
			return &v, nil
		}
	}

	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte("Error query DLI permission"),
		},
	}
}

func queryQueuePermission(client *golangsdk.ServiceClient, obj, userName string) (*auth.Privilege, error) {
	queueInfo := strings.SplitN(obj, ".", 2)
	rst, err := auth.ListQueuePermission(client, queueInfo[1])
	if err != nil {
		return nil, parseDliErrorToError404(err)
	}

	if rst != nil && !rst.IsSuccess {
		return nil, fmtp.Errorf("Error query DLI permission: %s", rst.Message)
	}

	for _, v := range rst.Privileges {
		if v.UserName == userName {
			return &v, nil
		}
	}

	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte("Error query DLI permission"),
		},
	}
}

func resourceDliPermissionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.DliV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud DLI v1 client: %s", err)
	}

	obj, userName := ParseTableInfoFromId(d.Id())
	if strings.HasPrefix(obj, "queues.") {
		queueName := strings.Replace(obj, "queues.", "", 1)
		opts := auth.GrantQueuePermissionOpts{
			QueueName:  queueName,
			UserName:   userName,
			Action:     updateAction,
			Privileges: []string{},
		}

		rst, createErr := auth.GrantQueuePermission(client, opts)
		if createErr != nil {
			return fmtp.DiagErrorf("Error granting permission in DLI: %s", createErr)
		}

		if rst != nil && !rst.IsSuccess {
			return fmtp.DiagErrorf("Error granting permission in DLI: %s", rst.Message)
		}
	} else {
		opts := auth.GrantDataPermissionOpts{
			UserName:   userName,
			Action:     updateAction,
			Privileges: []auth.DataPermission{{Object: obj, Privileges: []string{}}},
		}

		rst, createErr := auth.GrantDataPermission(client, opts)
		if createErr != nil {
			return common.CheckDeletedDiag(d, parseDliErrorToError404(createErr), "DLI")
		}

		if rst != nil && !rst.IsSuccess {
			return fmtp.DiagErrorf("Error delete DLI permission: %s", rst.Message)
		}
	}

	d.SetId("")
	return nil
}

func ParseAuthInfoFromId(id string) (object, userName string) {
	idArrays := strings.Split(id, "/")
	object = idArrays[0]
	userName = idArrays[1]
	return
}

// Object format:
// databases.Database_name
// databases.Database_name.tables.Table_name
// databases.Database_name.tables.Table_name.columns.Column_name
// jobs.flink.Flink_job_ID
// groups.Package_group_name
// resources.PackageName
// queues.queueName
func QueryPermission(client *golangsdk.ServiceClient, obj, userName string) (*auth.Privilege, error) {
	var permission *auth.Privilege
	var err error

	if strings.HasPrefix(obj, "databases") {
		permission, err = queryDatabaseRelatePermission(client, obj, userName)
	} else if checkPrefixMatchDataPermession(obj) {
		permission, err = queryDataPermission(client, obj, userName)
	} else if strings.HasPrefix(obj, "queues") {
		permission, err = queryQueuePermission(client, obj, userName)
	} else {
		return nil, fmtp.Errorf("The object is illegal:object=%s,userName=%s, error=%s", obj, userName, err)
	}

	return permission, err
}

func parseDliErrorToError404(respErr error) error {
	var apiError flinkjob.DliError

	if errCode, ok := respErr.(golangsdk.ErrDefault400); ok {
		pErr := json.Unmarshal(errCode.Body, &apiError)
		if pErr == nil && apiError.ErrorCode == "DLI.0002" {
			return golangsdk.ErrDefault404(errCode)
		}
	}
	return respErr
}
