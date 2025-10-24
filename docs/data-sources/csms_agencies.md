---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_csms_agencies"
description: |-
  Use this data source to check whether a agency exists.
---

# huaweicloud_csms_agencies

Use this data source to check whether a agency exists.

## Example Usage

```hcl
variable "secret_type" {}

data "huaweicloud_csms_agencies" "test" {
  secret_type = var.secret_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `secret_type` - (Required, String) Specifies the secret type.
  The valid values are as follows:
  + **RDS-FG**
  + **GaussDB-FG**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `agency_granted` - Whether the agency exists.
