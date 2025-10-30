---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_access_codes"
description: |-
  Use this data source to get the list of AOM prometheus instance access codes.
---

# huaweicloud_aom_access_codes

Use this data source to get the list of AOM prometheus instance access codes.

## Example Usage

```hcl
data "huaweicloud_aom_access_codes" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `access_codes` - Indicates the access codes.

  The [access_codes](#access_codes_struct) structure is documented below.

<a name="access_codes_struct"></a>
The `access_codes` block supports:

* `access_code_id` - Indicates the access code ID.

* `access_code` - Indicates the access code.

* `status` - Indicates the status.

* `create_at` - Indicates the creation time.
