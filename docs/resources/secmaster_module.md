---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_module"
description: |-
  Manages a SecMaster module resource within HuaweiCloud.
---

# huaweicloud_secmaster_module

Manages a SecMaster module resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

resource "huaweicloud_secmaster_module" "test" {
  workspace_id = var.workspace_id
  name         = "test-module"
  description  = "A test SecMaster module"
  module_type  = "tab"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the SecMaster module resource. If omitted,
  the provider-level region will be used. Changing this will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID.

* `name` - (Required, String) Specifies the name of the module.

* `description` - (Optional, String) Specifies the description of the module.

* `module_type` - (Optional, String) Specifies the type of the module. Valid values are **tab** and **section**.

* `module_json` - (Optional, String) Specifies the module-related information in JSON format.

* `metric_ids` - (Optional, String) Specifies the metric IDs when the module type is **section**.

* `thumbnail` - (Optional, String) Specifies the thumbnail of the module.

* `data_query` - (Optional, String) Specifies the data query method for the module.

* `boa_version` - (Optional, String) Specifies the BOA base version.

* `cloud_pack_id` - (Optional, String) Specifies the subscription package ID.

* `cloud_pack_name` - (Optional, String) Specifies the name of the subscription package.

* `cloud_pack_version` - (Optional, String) Specifies the version of the subscription package.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The module ID (also module ID).

* `project_id` - The project ID.

* `en_name` - The English name of the module.

* `en_description` - The English description of the module.

* `creator_id` - The ID of the creator.

* `create_time` - The creation time.

* `update_time` - The update time.

* `tag` - The module tag.

* `is_built_in` - Whether the module is built-in.

* `version` - The version of the SecMaster module.

## Import

SecMaster module can be imported using the `workspace_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_secmaster_module.test <workspace_id>/<id>
```
