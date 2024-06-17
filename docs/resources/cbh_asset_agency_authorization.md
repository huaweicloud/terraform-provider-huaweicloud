---
subcategory: "Cloud Bastion Host (CBH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbh_asset_agency_authorization"
description: |-
  Manages a CBH asset agency authorization resource within HuaweiCloud.
---

# huaweicloud_cbh_asset_agency_authorization

Manages a CBH asset agency authorization resource within HuaweiCloud.

-> After you enable CSMS credentials and KMS key agency authorization, you need to wait about `10` minutes, the CBH
instance can obtain a token with agency permissions. Destroying resources will not change the current asset
agency authorization status.

## Example Usage

```hcl
resource "huaweicloud_cbh_asset_agency_authorization" "test" {
  csms = true
  kms  = true
} 
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CBH asset agency authorization.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `csms` - (Required, Bool) Specifies whether to enable CSMS credential agency authorization. The value can be **true**
  or **false**.  
  If set to **true** to enable agency authorization, the CBH service will have the permission to query your CSMS
  credential list. You can select credentials as resource accounts on the CBH instance.

* `kms` - (Required, Bool) Specifies whether to enable KMS key agency authorization. The value can be **true** or
  **false**.  
  If set to **true** to enable agency authorization, the CBH service will have the permission to use the KMS interface
  to obtain the CSMS credential value. You can use this credential value to log in to the managed host on the CBH
  instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
