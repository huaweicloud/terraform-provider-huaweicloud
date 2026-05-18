---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_instance_groups"
description: |-
  Use this data source to get the list of DAS instance groups.
---

# huaweicloud_das_instance_groups

Use this data source to get the list of DAS instance groups.

## Example Usage

```hcl
data "huaweicloud_das_instance_groups" "test" {
  datastore_type = "MySQL"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the DAS instance groups are located.  
  If omitted, the provider-level region will be used.

* `datastore_type` - (Required, String) Specifies the database type.  
  The valid values are as follows:
  + **MySQL**
  + **TaurusDB**
  + **GaussDB**
  + **MariaDB**

* `name` - (Optional, String) Specifies the instance group name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - The list of the instance groups that matched the filter parameters.  
  The [groups](#instance_group_attr) structure is documented below.

<a name="instance_group_attr"></a>
The `groups` block supports:

* `id` - The instance group ID.

* `name` - The instance group name.

* `description` - The description of the instance group.
