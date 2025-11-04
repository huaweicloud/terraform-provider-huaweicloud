---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_key_replicate"
description: |-
  Manages a KMS key replication resource within HuaweiCloud.
---

# huaweicloud_kms_key_replicate

Manages a KMS key replication resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not delete the replicated key,
  but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "key_id" {}
variable "replica_region" {}
variable "key_alias" {}
variable "replica_project_id" {}

resource "huaweicloud_kms_key_replicate" "test" {
  key_id             = var.key_id
  replica_region     = var.replica_region
  key_alias          = var.key_alias
  replica_project_id = var.replica_project_id
  key_description    = "Replicated key for cross-region usage"
  
  tags = {
    environment = "production"
    owner       = "security"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `key_id` - (Required, String, NonUpdatable) Specifies the ID of the key to be replicated.
  The valid length is `36` bytes, meeting regular match **^[0-9a-z]{8}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{12}$**.
  For example: **0d0466b0-e727-4d9c-b35d-f84bb474a37f**.

* `replica_region` - (Required, String, NonUpdatable) Specifies the target region to which the key is replicated.
  The value should be a valid region name, such as **cn-north-4**.
  The value of this field cannot be the same as the region in the source key.

* `key_alias` - (Required, String, NonUpdatable) Specifies the alias of the replicated key.
  The alias must be unique in the target region and project.

* `replica_project_id` - (Required, String, NonUpdatable) Specifies the target project ID to which the key is replicated.

* `key_description` - (Optional, String, NonUpdatable) Specifies the description of the replicated key.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID of the replicated key.
  If omitted, the default enterprise project is used.

* `tags` - (Optional, Map, NonUpdatable) Specifies the key/value pairs to associate with the replicated key.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the replicated key.
