---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_remediation_configuration"
description: |-
  Manages a RMS remediation configuration resource within HuaweiCloud.
---

# huaweicloud_rms_remediation_configuration

Manages a RMS remediation configuration resource within HuaweiCloud.

## Example Usage

```hcl
variable "policy_assignment_id" {}
variable "target_id" {}
variable "resource_id " {}
variable "auth_type" {}
variable "auth_value" {}
variable "var_key" {}
variable "var_value" {}

resource "huaweicloud_rms_remediation_configuration" "test" {
  policy_assignment_id = var.policy_assignment_id
  target_type          = "rfs"
  target_id            = var.target_id

  resource_parameter {
    resource_id = var.resource_id 
  }

  static_parameter {
    var_key   = var.var_key
    var_value = var.var_value
  }
  
  auth_type             = var.auth_type
  auth_value            = var.auth_value
  maximum_attempts      = 6  
  retry_attempt_seconds = 60  
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `policy_assignment_id` - (Required, String, NonUpdatable) Specifies the policy assignment ID.

* `target_type` - (Required, String) Specifies the execution method of remediation.
  The valid value can be **fgs** or **rfs**.

* `target_id` - (Required, String) Specifies the ID of a remediation object.
  + If the execution method is **fgs**, the value is a function URN.
  + If the execution method is **rfs**, the value is the name and version ID that separated by a slash (/).
  If the version is not specified, V1 is used by default.

* `resource_parameter` - (Required, List) Specifies the dynamic parameter of remediation.

  The [resource_parameter](#resource_parameter_struct) structure is documented below.

* `automatic` - (Optional, Bool) Specifies whether remediation is automatic.
  The default value is **false**.

* `static_parameter` - (Optional, List) Specifies the static parameters for the remediation execution.

  The [static_parameter](#static_parameter_struct) structure is documented below.

* `auth_type` - (Optional, String) Specifies the authorization type for remediation configurations.
  The valid value can be **agency** or **trustAgency**.

* `auth_value` - (Optional, String) Specifies the information of dependent service authorization.

* `maximum_attempts` - (Optional, Int) Specifies the maximum number of retries allowed within a specified period.
  The maximum value is **25**. The minimum value is **1**. The default value is **5**.

* `retry_attempt_seconds` - (Optional, Int) Specifies the time period during which the number of attempts specified
  in the `maximum_attempts` can be tried.
  The maximum value is **43200**. The minimum value is **60**. The default value is **3600**.
  If remediation retries exceed the limit, corresponding resources will be classified as exceptions of remediation.

<a name="resource_parameter_struct"></a>
The `resource_parameter` block supports:

* `resource_id` - (Required, String) Specifies the parameter name for passing the resource ID.

<a name="static_parameter_struct"></a>
The `static_parameter` block supports:

* `var_key` - (Optional, String) Specifies the static parameter name.

* `var_value` - (Optional, String) Specifies the static parameter value in JSON format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The time when the remediation configuration was created.

* `updated_at` - The time when the remediation configuration was updated.

* `created_by` - The user who created the remediation configuration.

## Import

The RMS remediation configuration can be imported by using the `id`, e.g.

```bash
$ terraform import huaweicloud_rms_remediation_configuration.test <id>
```
