---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_name_validation"
description: |-
  Manages a DRS job name validation resource.
---

# huaweicloud_drs_name_validation

Manages a DRS job name validation resource.

-> 1. This resource is a one-time action resource used to verify the availability of a DRS job name. Deleting this resource
  will not perform any operations on the cloud, but will only remove the resource information from the tf state file.
  <br/>2. This resource is used to verify whether the job name meets the format requirements and is available
  before creating a DRS job.

## Example Usage

```hcl
variable "name" {}
variable "type" {}

resource "huaweicloud_drs_name_validation" "test" {
  name = var.name
  type = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the DRS job name to be validated.
  The name must meet the following rules: 4 to 50 characters, starts with a letter, can contain letters, digits,
  hyphens (-) and underscores (_), and cannot contain other special characters.

* `type` - (Required, String, NonUpdatable) Specifies the type of the DRS job.
  The valid values are as follows:
  + **TRANS**
  + **SUBSCRIPTION**
  + **OFFLINE**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `is_valid` - Whether the job name is valid.
