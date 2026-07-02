---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_classifier"
description: |-
  Manages a SecMaster classifier resource within HuaweiCloud.
---

# huaweicloud_secmaster_classifier

Manages a SecMaster classifier resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "dataclass_id" {}
variable "name" {}

resource "huaweicloud_secmaster_classifier" "test" {
  workspace_id = var.workspace_id
  name         = var.name
  dataclass_id = var.dataclass_id
  data_source  = "CFW"
  description  = "test description"

  classifier {
    direct_classifier = "false"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID.

* `name` - (Required, String) Specifies the name.

* `dataclass_id` - (Required, String, NonUpdatable) Specifies the data class ID.

* `data_source` - (Required, String, NonUpdatable) Specifies the data source.  

* `description` - (Required, String) Specifies the description.

* `classifier` - (Required, List) Specifies the classifier information.

  The [classifier](#secmaster_classifier_struct) structure is documented below.

<a name="secmaster_classifier_struct"></a>
The `classifier` block supports:

* `direct_classifier` - (Required, String) Specifies whether to classify directly.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (the classifier ID).

* `mapping_id` - The mapping ID.

* `project_id` - The project ID.

* `dataclass_name` - The data class name.

* `status` - The status.

* `complete_degree` - The completion degree.

* `instance_num` - The number of associated instances.

* `built_in` - Whether the data is built-in.

* `create_time` - The creation time.

* `creator_id` - The creator ID.

* `creator_name` - The creator name.

* `update_time` - The update time.

* `modifier_id` - The modifier ID.

* `modifier_name` - The modifier name.

## Import

The classifier can be imported using the `workspace_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_classifier.test <workspace_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `data_source`.
It is generally recommended running `terraform plan` after importing a classifier.
You can then decide if changes should be applied to the classifier, or the resource definition should be updated to
align with the classifier. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_secmaster_classifier" "test" {
  ...

  lifecycle {
    ignore_changes = [
      data_source,
    ]
  }
}
```
