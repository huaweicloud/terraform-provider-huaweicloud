---
subcategory: "LakeFormation"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lakeformation_instances"
description: |-
  Use this data source to get the list of LakeFormation instances within HuaweiCloud.
---

# huaweicloud_lakeformation_instances

Use this data source to get the list of LakeFormation instances within HuaweiCloud.

## Example Usage

### Query all instances

```hcl
data "huaweicloud_lakeformation_instances" "test" {}
```

### Query instances using name filter parameter

```hcl
variable "instance_name" {}

data "huaweicloud_lakeformation_instances" "test" {
  name = var.instance_name
}
```

### Query instances in the recycle bin

```hcl
data "huaweicloud_lakeformation_instances" "test" {
  in_recycle_bin = true
}
```

### Query instances using enterprise project ID filter parameter

```hcl
variable "enterprise_project_id" {}

data "huaweicloud_lakeformation_instances" "test" {
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `in_recycle_bin` - (Optional, Bool) Specifies whether to query instances in the recycle bin.
  Defaults to `false`.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the instances
  belong.  
  Defaults to query instances under all enterprise projects.

* `name` - (Optional, String) Specifies the name of the instance to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

`id` - The data source ID.

* `instances` - The list of instances that matched filter parameters.
  The [instances](#lakeformation_instances_attr) structure is documented below.

<a name="lakeformation_instances_attr"></a>
The `instances` block supports:

* `instance_id` - The ID of the instance.

* `name` - The name of the instance.

* `description` - The description of the instance.

* `enterprise_project_id` - The ID of the enterprise project to which the instance belongs.

* `shared` - Whether the instance is shared.

* `default_instance` - Whether the instance is the default instance.

* `create_time` - The creation timestamp of the instance.

* `update_time` - The update timestamp of the instance.

* `status` - The status of the instance.

* `in_recycle_bin` - Whether the instance is in the recycle bin.

* `tags` - The key/value pairs to associate with the instance.
