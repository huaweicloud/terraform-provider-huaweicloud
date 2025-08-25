---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_document_execute"
description: |-
  Manages a COC document execute resource within HuaweiCloud.
---

# huaweicloud_coc_document_execute

Manages a COC document execute resource within HuaweiCloud.

~> Deleting document execute resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "document_id" {}
variable "project_id" {}
variable "domain_id" {}

resource "huaweicloud_coc_document_execute" "test" {
  document_id = var.document_id
  parameters {
    key   = "project_id"
    value = var.project_id
  }
  parameters {
    key   = "agency_urn"
    value = "iam::${var.domain_id}:agency:ServiceAgencyForCOC"
  }
}
```

## Argument Reference

The following arguments are supported:

* `document_id` - (Required, String, NonUpdatable) Specifies the document ID.

* `version` - (Optional, String, NonUpdatable) Specifies the document version number. If not specified, the default is
  the latest version.

* `parameters` - (Optional, List, NonUpdatable) Specifies the global parameters.
  The [parameters](#key_value_struct) structure is documented below.

* `sys_tags` - (Optional, List, NonUpdatable) Specifies the list of system tags.
  The [sys_tags](#key_value_struct) structure is documented below.

* `target_parameter_name` - (Optional, String, NonUpdatable) Specifies the parameter name of the batch execution object
  in rate control mode.

* `targets` - (Optional, List, NonUpdatable) Specifies the method to use with `target_parameter_name`.
  The [targets](#targets_struct) structure is documented below.

* `document_type` - (Optional, String, NonUpdatable) Specifies the type of document to perform.
  The value can be **PRIVATE** or **PUBLIC**.

* `description` - (Optional, String, NonUpdatable) Specifies the execution description.

<a name="key_value_struct"></a>
The `parameters`, `sys_tags` blocks support:

* `key` - (Optional, String, NonUpdatable) Specifies the key.

* `value` - (Optional, String, NonUpdatable) Specifies the value.

<a name="targets_struct"></a>
The `targets` blocks support:

* `key` - (Optional, String, NonUpdatable) Specifies the dimension of the instantiated execution target.
  The enumeration value as **InstanceValues**, **BatchValues**.

* `values` - (Optional, String, NonUpdatable) Specifies the target instance to be executed based on the global parameter
  specified by `target_parameter_name`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The execution ID.

* `execution_parameters` - Indicates the global parameters for work order execution.

  The [execution_parameters](#data_key_value_struct) structure is documented below.

* `document_name` - Indicates the document name.

* `document_version_id` - Indicates the document version ID.

* `start_time` - Indicates the work order execution start time.

* `end_time` - Indicates the work order execution end time.

* `update_time` - Indicates the work order update time.

* `creator` - Indicates the work order creator.

* `status` - Indicates the work order status.

* `tags` - Indicates the key/value tags of the work order.

* `type` - Indicates the work order type.

<a name="data_key_value_struct"></a>
The `execution_parameters` block supports:

* `key` - Indicates the key.

* `value` - Indicates the value.

## Import

The COC document execution can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_coc_document_execute.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `parameters`, `sys_tags` and `document_type`.
It is generally recommended running `terraform plan` after importing a document execution.
You can then decide if changes should be applied to the document execution, or the resource definition should be updated
to align with the document execution. Also you can ignore changes as below.

```hcl
resource "huaweicloud_coc_document_execute" "test" {
    ...

  lifecycle {
    ignore_changes = [
      parameters, sys_tags, document_type,
    ]
  }
}
```
