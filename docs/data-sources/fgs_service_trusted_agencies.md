---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_service_trusted_agencies"
description: |-
  Use this data source to query the list of service trusted agencies within HuaweiCloud.
---

# huaweicloud_fgs_service_trusted_agencies

Use this data source to query the list of service trusted agencies within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_fgs_service_trusted_agencies" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the service trusted agencies are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

The following attributes are exported:

* `id` - The data source ID.

* `agencies` - The list of service trusted agencies.  
  The [agencies](#fgs_service_trusted_agencies_attr) structure is documented below.

<a name="fgs_service_trusted_agencies_attr"></a>
The `agencies` block supports:

* `name` - The name of the trusted agency.

* `expire_time` - The expiration time of the trusted agency.  
  When the agency never expires, the value is empty string.
