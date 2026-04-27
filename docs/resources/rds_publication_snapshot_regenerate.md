---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_publication_snapshot_regenerate"
description: |-
  Manages an RDS publication snapshot regenerate resource within HuaweiCloud.
---

# huaweicloud_rds_publication_snapshot_regenerate

Manages an RDS publication snapshot regenerate resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "publication_id" {}

resource "huaweicloud_rds_publication_snapshot_regenerate" "test" {
  instance_id    = var.instance_id
  publication_id = var.publication_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS instance.

* `publication_id` - (Required, String, NonUpdatable) Specifies the ID of the publication.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
