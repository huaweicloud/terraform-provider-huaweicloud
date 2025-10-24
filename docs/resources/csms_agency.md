---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_csms_agency"
description: |-
  Manages a resource to create agency within HuaweiCloud.
---

# huaweicloud_csms_agency

Manages a resource to create agency within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "secret_type" {}

resource "huaweicloud_csms_agency" "test" {
  secret_type = var.secret_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `secret_type` - (Required, String, NonUpdatable) Specifies the secret type.
  The valid values are as follows:
  + **RDS-FG**
  + **GaussDB-FG**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `agencies` - The List of agencies.
  The [agencies](#agency_struct) structure is documented below.

<a name="agency_struct"></a>
The `agencies` block supports:

* `agency_id` - The agency ID.

* `agency_name` - The agency name.
