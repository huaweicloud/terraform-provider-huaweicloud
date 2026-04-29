---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_stack_set"
description: |-
  Manages a RFS stack set resource within HuaweiCloud.
---

# huaweicloud_rfs_stack_set

Manages a RFS stack set resource within HuaweiCloud.

## Example Usage

```hcl
variable "stack_set_name" {}
variable "stack_set_description" {}
variable "permission_model" {}
variable "template_uri" {}
variable "administration_agency_name" {}
variable "managed_agency_name" {}
variable "initial_stack_description" {}

resource "huaweicloud_rfs_stack_set" "test" {
  stack_set_name             = var.stack_set_name
  stack_set_description      = var.stack_set_description
  permission_model           = var.permission_model
  template_uri               = var.template_uri
  administration_agency_name = var.administration_agency_name
  managed_agency_name        = var.managed_agency_name
  initial_stack_description  = var.initial_stack_description

  managed_operation {
    enable_parallel_operation = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the resource is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `stack_set_name` - (Required, String, NonUpdatable) Specifies the name of a stack set. The name must be unique within
  the domain_id + region scope. It can contain Chinese characters, uppercase and lowercase English letters, digits,
  underscores (_), and hyphens (-). The first character must be a Chinese character or an English letter. The name is
  case-sensitive and can contain a maximum of `128` characters.

* `stack_set_description` - (Optional, String) Specifies the description of a stack set. The description can contain
  a maximum of `1,024` characters.

* `permission_model` - (Optional, String, NonUpdatable) Specifies the permission model that defines how RFS creates the
  required delegation when operating stack sets. To update the permission model, users can delete and recreate a stack
  set with the same name. The valid values are as follows:
  + **SELF_MANAGED**: Self-managed. Users need to manually create delegations in advance, including both the
    delegation from the management account to RFS and the delegation from member accounts to the management account.
    If the delegation does not exist or has insufficient permissions, creating a stack set will not fail; errors will
    only occur when creating stack instances.
  + **SERVICE_MANAGED**: Service-managed. RFS will automatically create all IAM delegations required for deploying
    Organization member accounts. Users need to enable "Resource Orchestration Stack Set Service" in the Organization
    trusted services list in advance. Only the Organization's management account or the delegated administrator for
    "Resource Orchestration Stack Set Service" can specify SERVICE_MANAGED to create stack sets.  
    Defaults to **SELF_MANAGED**.

* `administration_agency_name` - (Optional, String) Specifies the name of the administration agency. RFS uses this
  delegation to obtain the permissions that member accounts delegate to the management account. This delegation must
  contain the `iam:tokens:assume` permission to subsequently obtain credentials from the managed delegation. If it does
  not contain this permission, errors will occur when adding or deploying instances. When using the **SELF_MANAGED**
  permission type, either `administration_agency_name` or `administration_agency_urn` must be specified, but not both.
  It is recommended to use `administration_agency_urn` when using trusted delegations. The `administration_agency_name`
  parameter only supports receiving delegation names. If a trusted delegation name is provided, template deployment will
  fail. When using the **SERVICE_MANAGED** permission type, specifying this parameter will result in a `400` error. The
  maximum length is `64` characters.

* `managed_agency_name` - (Optional, String) Specifies the name of the managed agency. RFS uses this delegation to
  obtain the permissions required for actual resource deployment. The delegation names from different member accounts
  to the management account must be consistent. Currently, it does not support defining different delegation permissions
  based on different providers. When using the **SELF_MANAGED** permission type, this parameter must be specified. When
  using the **SERVICE_MANAGED** permission type, specifying this parameter will result in a `400` error. The maximum
  length is `64` characters.

* `template_body` - (Optional, String, NonUpdatable) Specifies the HCL template that describes the target state of
  resources. RFS will compare the differences between this template and the current state of remote resources. Either
  `template_body` or `template_uri` must be specified, but not both. Note that stack sets do not support sensitive data
  encryption. RFS will directly use, log, display, and store the corresponding `template_body` in plaintext. The maximum
  length is `51,200` characters.

  -> Exactly one of `template_body` and `template_uri` must be specified.

* `template_uri` - (Optional, String, NonUpdatable) Specifies the OBS URL of the HCL template that describes the target
  state of resources. RFS will compare the differences between this template and the current state of remote resources.
  Ensure that the OBS bucket location is consistent with the RFS service location. The corresponding file should be
  either a pure Terraform file or a ZIP archive. Pure Terraform files must end with `.tf` or `.tf.json`, comply with HCL
  syntax, and use UTF-8 encoding. ZIP archives must end with `.zip`. The extracted files must not contain `.tfvars`
  files. The maximum size before and after extraction is `1MB`. The number of files in the ZIP archive cannot exceed
  `100`. Either `template_body` or `template_uri` must be specified, but not both. Note that stack sets do not support
  sensitive data encryption. RFS will directly use, log, display, and store the template file content corresponding to
  `template_uri` in plaintext. If the template file is a ZIP archive, the length of internal file or folder names must
  not exceed `255` bytes, the deepest path length must not exceed `2,048` bytes, and the ZIP package size must not
  exceed `1MB`. The maximum length is `2,048` characters.

  -> Exactly one of `template_body` and `template_uri` must be specified.

* `vars_body` - (Optional, String, NonUpdatable) Specifies the content of the HCL parameter file. HCL templates support
  parameter input, meaning that the same template can achieve different effects with different parameters. The
  `vars_body` uses the HCL tfvars format. Users can submit the content from `.tfvars` files to `vars_body`. RFS supports
  both `vars_body` and `vars_uri`. If the same variable is declared in both methods, a `400` error will be reported. If
  `vars_body` is too large, users can use `vars_uri` instead. Stack sets do not support sensitive data encryption. RFS
  will directly use, log, display, and store the corresponding `vars_body` in plaintext. The maximum length is `51,200`
  characters.

* `vars_uri` - (Optional, String, NonUpdatable) Specifies the OBS URL of the HCL parameter file. HCL templates support
  parameter input, meaning that the same template can achieve different effects with different parameters. Ensure that
  the OBS bucket location is consistent with the RFS service location. The `vars_uri` must point to an OBS pre-signed
  URL address; other addresses are not currently supported. RFS supports both `vars_body` and `vars_uri`. If the same
  variable is declared in both methods, a `400` error will be reported. The content in `vars_uri` uses the HCL tfvars
  format. Users can save the content from `.tfvars` files, upload them to OBS, and pass the OBS pre-signed URL to
  `vars_uri`. Stack sets do not support sensitive data encryption. RFS will directly use, log, display, and store the
  parameter file content corresponding to `vars_uri` in plaintext. The maximum length is `2,048` characters.

* `initial_stack_description` - (Optional, String) Specifies the description of the initial stack. It can be used to
  help customers identify stacks managed by the stack set. Stacks under the stack set will only use this description
  during creation. Customers can update the initial stack description through the UpdateStackSet API. Subsequent updates
  to the stack set description will not synchronize to update the descriptions of already managed stacks. The maximum
  length is `1,024` characters.

* `administration_agency_urn` - (Optional, String) Specifies the URN of the administration agency. RFS uses this
  delegation to obtain the permissions that member accounts delegate to the management account. This delegation must
  contain the `sts:tokens:assume` permission to subsequently obtain credentials from the managed delegation. If it does
  not contain this permission, errors will occur when adding or deploying instances. When using the **SELF_MANAGED**
  permission type, either `administration_agency_name` or `administration_agency_urn` must be specified, but not both.
  It is recommended to use `administration_agency_urn` when using trusted delegations. The `administration_agency_name`
  parameter only supports receiving delegation names. If a trusted delegation name is provided, template deployment will
  fail. When using the **SERVICE_MANAGED** permission type, specifying this parameter will result in a `400` error.

* `call_identity` - (Optional, String) Specifies the identity that calls the stack set API. This parameter can only be
  specified when the stack set permission model is **SERVICE_MANAGED**. It is used to specify whether the user calls the
  stack set as the organization management account or as a delegated administrator in a member account. When the stack
  set permission model is **SELF_MANAGED**, the default value is **SELF**. Regardless of the specified user identity,
  the stack set involved in the operation always belongs to the organization management account.
  The valid values are as follows:
  + **SELF**: Calls as the organization management account.
  + **DELEGATED_ADMIN**: Calls as a delegated administrator. The Huawei Cloud account must have been registered as a
    delegated administrator for "Resource Orchestration Stack Set Service" in the organization.
    Defaults to **SELF**.

* `managed_operation` - (Optional, List) Specifies the managed operation configuration.  
  The [suggestion](#stackset_managed_operation) structure is documented below.

<a name="stackset_managed_operation"></a>
The `managed_operation` block supports:

* `enable_parallel_operation` - (Optional, Bool) Specifies whether to enable parallel operations for stack instances.  
  Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the same as `stack_set_name`.

* `stack_set_id` - The ID of the stack set.

* `create_time` - The creation time of the stack set. It is represented in UTC format (YYYY-MM-DDTHH:mm:ss.SSSZ),
  such as **1970-01-01T00:00:00.000Z**.

* `update_time` - The update time of the stack set. It is represented in UTC format (YYYY-MM-DDTHH:mm:ss.SSSZ),
  such as **1970-01-01T00:00:00.000Z**.

* `status` - The status of the stack set.

* `vars_uri_content` - The content of the variable file referenced by `vars_uri`.

* `organizational_unit_ids` - The list of organizational unit IDs associated with the stack set.

## Import

Stack set can be imported using the `stack_set_name`, e.g.

```bash
$ terraform import huaweicloud_rfs_stack_set.test <stack_set_name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include `template_uri`,`template_body`,`vars_uri`,
`call_identity`. It is generally recommended running `terraform plan` after importing a resource. You can then decide
if changes should be applied to the resource, or the resource definition should be updated to align with the resource.
Also you can ignore changes as below.

```hcl
resource "huaweicloud_rfs_stack_set" "test" {
  ...

  lifecycle {
    ignore_changes = [
      template_uri,
      template_body,
      vars_uri,
      call_identity
    ]
  }
}
