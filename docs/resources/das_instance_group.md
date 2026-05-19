---
subcategory: "DAS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_instance_group"
description: |-
  Manages a DAS instance group resource within HuaweiCloud.
---

# huaweicloud_das_instance_group

Manages a DAS instance group resource within HuaweiCloud.

## Example Usage

```hcl
variable "group_name" {}
variable "description" {}

resource "huaweicloud_das_instance_group" "test" {
  datastore_type = "MySQL"
  group_name     = var.group_name
  description    = var.description
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the instance group is located.  
  If omitted, the provider-level region will be used.

* `datastore_type` - (Required, String, NonUpdatable) Specifies the database type.  
  The valid values are as follows:
  + **MySQL**
  + **TaurusDB**
  + **GaussDB**
  + **MariaDB**

* `group_name` - (Required, String) Specifies the instance group name.  
  The length should not exceed `64` characters.

* `description` - (Required, String) Specifies the description of the instance group.  
  The length should not exceed `128` characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `group_id`.

## Import

The DAS instance group can be imported using `<datastore_type>/<id>`, e.g.

```bash
$ terraform import huaweicloud_das_instance_group.test <datastore_type>/<id>
```
