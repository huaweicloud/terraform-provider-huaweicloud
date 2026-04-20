---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_stack_set_deployment"
description: |-
  Manages a resource to deploy stack set within HuaweiCloud.
---

# huaweicloud_rfs_stack_set_deployment

Manages a resource to deploy stack set within HuaweiCloud.

-> This resource is a one-time action resource used to trigger stack set deployment. Deleting this resource will
  not cancel deployment operation, but will only remove resource information from tf state file.

## Example Usage

```hcl
variable "stack_set_name" {}
variable "region" {}
variable "domain_id" {}

resource "huaweicloud_rfs_stack_set_deployment" "test" {
  stack_set_name = var.stack_set_name

  deployment_targets {
    regions    = [var.region]
    domain_ids = [var.domain_id]
  }

  template_body = <<-EOT
    resource "huaweicloud_vpc" "example" {
      name = "example-vpc"
      cidr = "192.168.0.0/16"
    }
  EOT
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `stack_set_name` - (Required, String, NonUpdatable) Specifies the name of the stack set.

* `deployment_targets` - (Required, List, NonUpdatable) Specifies the deployment target information.

  The [deployment_targets](#deployment_targets_struct) structure is documented below.

* `stack_set_id` - (Optional, String, NonUpdatable) Specifies the unique ID of the stack set.
  For parallel development in a team, users may want to ensure that the stack set they operate is the one created by
  themselves, not the one with the same name created by other teammates after deleting the previous one.
  Therefore, they can use this ID for strong matching.
  RFS ensures that the ID of each stack set is different and does not change with updates. If the `stack_set_id` value
  is different from the current stack set ID, `400` is returned.

* `template_body` - (Optional, String, NonUpdatable) Specifies the HCL template that describes the target state of
  resources. The resource orchestration service will compare the differences between this template and the current
  remote resource state.

* `template_uri` - (Optional, String, NonUpdatable) Specifies the OBS address of an HCL template. The template describes
  the target status of a resource. RFS compares the difference between the statuses of this template and the current
  remote resource. Ensure that the OBS address is located in the same region as the RFS.
  The corresponding file must be a tf file or a zip package.
  A ".tf" file must be named with a ".tf" or ".tf.json" suffix, compatible with HCL, and UTF-8 encoded.
  Currently, only the ".zip" package is supported. The file name extension must be ".zip". The decompressed files cannot
  contain ".tfvars" files. The maximum size of the file is `1` MB before decompression and `1` MB after decompression.
  A maximum of `100` files can be archived to one ".zip" package.

-> Either `template_body` or `template_uri` must be specified, but not both.

* `vars_uri` - (Optional, String, NonUpdatable) Specifies the OBS address of the HCL parameter file.
  Transferring parameters is supported by the HCL template.
  The same template can use different parameters for different purposes.
  + The value of `vars_uri` needs to point to a pre-signed URL address of OBS; other addresses are not currently supported.
  + The resource orchestration service supports both `vars_body` and `vars_uri`. If the same variable is declared using
    either of these methods, a `400` error will occur.
  + The content in `vars_uri` uses HCL's tfvars format. Users can save the content in ".tfvars" to a file and upload it
    to OBS, passing the OBS pre-signed URL to `vars_uri`.
  + The resource stack does not support encryption of sensitive data; the resource orchestration service will directly
    use, log, display, and store the parameter file content corresponding to `vars_uri` in plaintext.

* `vars_body` - (Optional, String, NonUpdatable) Specifies the content of the HCL variable file. Transferring parameters
  is supported by the HCL template. The same template can use different parameters for different purposes.
  + The `vars_body` uses the tfvars format of HCL. You can submit the content in the ".tfvars" file to the `vars_body`.
  + RFS supports `vars_structure`, `vars_body`, and `vars_uri`. If they declare the same variable, error `400` will be
    reported.
  + If `vars_body` is too large, you can use `vars_uri`.
  + Stack sets do not encrypt sensitive data. RFS uses logs, displays, and stores `vars_body` as plaintext.

-> Either `vars_body` or `vars_uri` can be specified, but not both.

* `operation_preferences` - (Optional, List, NonUpdatable) Specifies the user-specified preferences for how to perform
  a stack set operation. This parameter takes effect only in a specified single operation.
  If this parameter is not specified, the default operation preferences is that only one stack is deployed at a time and
  after all stack instances in a region are deployed completely, the next region will be selected randomly for deployment.
  The default value of failure tolerance count in a region is `0`.

  The [operation_preferences](#operation_preferences_struct) structure is documented below.

* `call_identity` - (Optional, String, NonUpdatable) Specifies the call identity. The parameter is only supported when
  the stack set permission model is **SERVICE_MANAGED**. Specify whether you are acting as an account administrator in
  the organization's management account or as a delegated administrator in a member account. Defaults to **SELF**.
  When the resource stack permission mode is **SELF_MANAGED**, defaults to **SELF**.
  No matter what call identity is specified, the stack set involved in request is always belonging to management account.
  Valid values are:
  + **SELF**: Call as organization management account.
  + **DELEGATED_ADMIN**: Call as service delegated administrator.

<a name="deployment_targets_struct"></a>
The `deployment_targets` block supports:

* `regions` - (Required, List of String, NonUpdatable) Specifies the list of regions involved in the stack set operation.
  Stack instances in the stack set are selected for deployment. This operation applies to the Cartesian product of
  the `regions` and `domain_ids` input by the user. If a region that is not managed by the stack set is specified,
  an error is reported.

* `domain_ids` - (Optional, List of String, NonUpdatable) Specifies the list of tenant IDs involved in the operation.
  Stack instances in the stack set are selected for deployment. This operation applies to the Cartesian product of
  the `regions` and `domain_ids` input by the user. If a `domain_id` that is not managed by the stack set is specified,
  an error is reported. When the stack set permission model is **SERVICE_MANAGED**, this parameter needs to be used with
  `domain_id_filter_type`. It's used to specify, exclude or additionally deploy the domain IDs of member accounts from
  the organizational units in the deployment target.

* `domain_ids_uri` - (Optional, String, NonUpdatable) Specifies the OBS address of the file containing tenant IDs.
  Tenant IDs are separated by commas (,) and line breaks are supported. Currently, only CSV files are supported, and the
  files should be encoded in UTF-8. The file size cannot exceed `100` KB.
  Do not use Excel for operations on the CSV file to be uploaded. Otherwise, inconsistencies may occur in results read
  from the CSV file. You are advised to use Notepad to open the file and check whether the content complies with your
  expectation.
  If this parameter is specified in the DeployStackSet API, stack instances in the stack set are selected for deployment.
  This operation applies to the Cartesian product of the `domain_ids_uri` file and regions input by the user.
  If a `domain_id` that is not managed by the stack set is specified, an error is reported.
  When the stack set permission model is **SERVICE_MANAGED**, this parameter needs to be used with `domain_id_filter_type`.
  Used to specify, exclude or additionally deploy the domain IDs of member accounts from the organizational units in
  the deployment target.

  -> Either `domain_ids` or `domain_ids_uri` must be specified, but not both.

* `organizational_unit_ids` - (Optional, List of String, NonUpdatable) Specifies the list of organizational unit IDs.
  This parameter is only allowed to be specified when the stack set permission model is **SERVICE_MANAGED**.
  The list of `organizational_unit_ids`, it can be the root organization (Root) ID or the ID of organizational units.
  Only OU IDs that have been managed by the resource stack set are allowed to be specified. If you specify OU IDs that
  are not managed by the resource stack set records, an error will be reported.

* `domain_id_filter_type` - (Optional, String, NonUpdatable) Specifies the domain IDs filter type. This parameter is only
  supported when stack set permission model is **SERVICE_MANAGED**. Defaults to **NONE**.
  Valid values are:
  + **INTERSECTION**: Select specified accounts from the OUs in deployment target for deployment. You can specify either
    `domain_ids` or `domain_ids_uri`, but not both.
  + **DIFFERENCE**: Exclude specified accounts from the OUs in deployment target for deployment. You can specify either
    `domain_ids` or `domain_ids_uri`, but not both.
  + **UNION**: In addition to deploy all accounts from the OUs in deployment target, it will also deploy to the specified
    account. Users can deploy the OU and specific individual accounts in stack set operation by specifying both
    `organizational_unit_ids` and `domain_ids`/`domain_ids_uri`. You can specify either `domain_ids` or `domain_ids_uri`,
    but not both. CreateStackInstances does not allow using this type.
  + **NONE**: Only deploy to all accounts from the OUs in deployment target. You can not specify `domain_ids` or
    `domain_ids_uri`.

<a name="operation_preferences_struct"></a>
The `operation_preferences` block supports:

* `region_concurrency_type` - (Optional, String, NonUpdatable) Specifies the concurrency type of deploying stack
  instances in regions. The value is case-sensitive.
  Valid values are:
  + **SEQUENTIAL**: Stack instances are deployed in sequence among regions, that is, after all stack instances in a
    region are deployed completely, the next region will be selected for deployment.
  + **PARALLEL**: Stack instances are deployed in all specified regions concurrently.

* `region_order` - (Optional, List of String, NonUpdatable) Specifies the region deployment order. This parameter can be
  specified only when `region_concurrency_type` is set to **SEQUENTIAL**. The `region_order` must only contain all
  regions in this deployment target. If this parameter is not specified, the region deployment order is random.
  The `region_order` takes effect only during a specified single operation.

* `failure_tolerance_count` - (Optional, Int, NonUpdatable) Specifies the maximum number of failed stack instances in
  a region. The value must be `0` or a positive integer. Defaults to `0`.
  If the value of `region_concurrency_type` is **SEQUENTIAL**, when the number of stack instances that deploy failed in
  a region exceeds the `failure_tolerance_count`, all other instances that are still in **WAIT_IN_PROGRESS** status will
  be canceled. The status of the canceled instance changes to **CANCEL_COMPLETE**.
  If the value of `region_concurrency_type` is **PARALLEL**, when the number of stack instances that deploy failed in
  a region exceeds the `failure_tolerance_count`, the stack set only cancels all instances that are still in
  **WAIT_IN_PROGRESS** status in this region. The status of the canceled instance changes to **CANCEL_COMPLETE**.
  Stack instances that are in **OPERATION_IN_PROGRESS** status or have been deployed
  (that is, in **OPERATION_COMPLETE** or **OPERATION_FAILED** status) are not affected.

* `failure_tolerance_percentage` - (Optional, Int, NonUpdatable) Specifies the maximum percentage of failed stack
  instances in a region. The value must be `0` or a positive integer. Defaults to `0`.
  By multiplying the `failure_tolerance_percentage` by the number of stack instances in the region and rounding it down,
  the actual number of failure tolerance count can be obtained.

-> Either `failure_tolerance_count` or `failure_tolerance_percentage` can be specified, but not both.

* `max_concurrent_count` - (Optional, Int, NonUpdatable) Specifies the maximum number of concurrent accounts can be
  deployed in a region. The value must be a positive integer. The default value is `1`.
  `max_concurrent_count` is at most one more than the failure tolerance count. If `failure_tolerance_percentage` is
  specified, `max_concurrent_count` is at most one more than the result of `failure_tolerance_percentage` multiplied by
  the number of stack instances in a region to guarantee that the deployment stops at the required level of failure
  tolerance.

* `max_concurrent_percentage` - (Optional, Int, NonUpdatable) Specifies the maximum percentage of concurrent accounts
  can be deployed in a region. The value must be a positive integer. The default value is `1`.
  The RFS calculates the actual maximum number of concurrent accounts by rounding down the value obtained by multiplying
  the percentage by the number of stack instances in each region.
  This actual maximum number of concurrent accounts is at most one more than the failure tolerance count.
  If `failure_tolerance_percentage` is specified, this actual maximum number of concurrent accounts is at most one more
  than the result of `failure_tolerance_percentage` multiplied by the number of stack instances in a region to guarantee
  that the deployment stops at the required level of failure tolerance.

-> Either `max_concurrent_count` or `max_concurrent_percentage` can be specified, but not both.

* `failure_tolerance_mode` - (Optional, String, NonUpdatable) Specifies the failure tolerance mode of deploying stack
  instances in regions. The value can be **STRICT_FAILURE_TOLERANCE** or **SOFT_FAILURE_TOLERANCE**.
  Defaults to **STRICT_FAILURE_TOLERANCE**. The value is case-sensitive.
  Detailed introduction:
  + **STRICT_FAILURE_TOLERANCE**: This option dynamically lowers the concurrency level to ensure the number of failed
    stack instances never exceeds the value of failure_tolerance_count + `1`. If failure_tolerance_percentage is specified,
    this option ensures the number of failed stack instances never exceeds the result of failure_tolerance_percentage
    multiplied by the number of stack instances in a region.
  + The initial actual maximum number of concurrent accounts is max_concurrent_count. If `max_concurrent_percentage` is
    specified, the initial actual maximum number of concurrent accounts is the result of `max_concurrent_percentage`
    multiplied by the number of stack instances. The actual maximum number of concurrent accounts is then reduced
    proportionally by the number of failed stack instances.
  + **SOFT_FAILURE_TOLERANCE**: This option separates `failure_tolerance_count` (`failure_tolerance_percentage`) from
    the actual maximum number of concurrent accounts. This option allows actual maximum number of concurrent accounts
    to keep at the concurrency level set by the `max_concurrent_count`, or `max_concurrent_percentage`.
  + This option does not ensure the number of failed stack instances is less than failure_tolerance_count + `1`.
    If `failure_tolerance_percentage` is specified, this option does not ensure the number of failed stack instances is
    less than the result of `max_concurrent_percentage` multiplied by the number of stack instances.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the stack set operation ID.
