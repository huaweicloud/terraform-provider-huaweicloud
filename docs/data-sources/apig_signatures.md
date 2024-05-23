---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_signatures"
description: |-
  Use this data source to query the signatures under the specified APIG instance within HuaweiCloud.
---

# huaweicloud_apig_signatures

Use this data source to query the signatures under the specified APIG instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "signature_name" {}

data "huaweicloud_apig_signatures" "test" {
  instance_id = var.instance_id
  name        = var.signature_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the signatrues belong.

* `signature_id` - (Optional, String) Specifies the ID of signature to be queried.

* `name` - (Optional, String) Specifies the name of signature to be queried.  
  The valid length is limited from `3` to `64`, only English letters, Chinese characters, digits and underscores (_) are
  allowed. The name must start with an English letter or Chinese character.

* `type` - (Optional, String) Specifies the type of signature to be queried.  
  The valid values are as follows:
  + **basic**: Basic auth type.
  + **hmac**: HMAC type.
  + **aes**: AES type

* `algorithm` - (Optional, String) Specifies the algorithm of the signature to be queried.  
 This parameter is only available when signature `type` is `aes`.  
  The valid values are as follows:
  + **aes-128-cfb**
  + **aes-256-cfb**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `signatures` - All signature key that match the filter parameters.
  The [signatures](#attrblock_signatures) structure is documented below.

<a name="attrblock_signatures"></a>
The `signatures` block supports:

* `id` - The ID of the signature.

* `name` - The name of the signature.

* `type` - The type of the signature.

* `key` - The key of the signature.

* `secret` - The secret of the signature.

* `algorithm` - The algorithm of the signature.

* `bind_num` - The number of bound APIs.

* `created_at` - The creation time of the signature, in RFC3339 format.

* `updated_at` - The latest update time of the signature, in RFC3339 format.
