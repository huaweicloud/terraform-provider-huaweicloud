---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_api_associated_signatures"
description: |-
  Use this data source to query the signatures associated with the specified API within HuaweiCloud.
---

# huaweicloud_apig_api_associated_signatures

Use this data source to query the signatures associated with the specified API within HuaweiCloud.

## Example Usage

### Query the contents of all signatures bound to the current API

```hcl
variable "instance_id" {}
variable "associated_api_id" {}

data "huaweicloud_apig_api_associated_signatures" "test" {
  instance_id = var.instance_id
  api_id      = var.associated_api_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the associated signatures.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the signatures belong.

* `api_id` - (Required, String) Specifies the ID of the API bound to the signature.

* `signature_id` - (Optional, String) Specifies the ID of the signature.

* `name` - (Optional, String) Specifies the name of the signature.

* `type` - (Optional, String) Specifies the type of the signature.  
  The valid values are as follows:
  + **basic**: Basic auth type.
  + **hmac**: HMAC type.
  + **aes**: AES type

* `env_id` - (Optional, String) Specifies the ID of the environment where the API is published.

* `env_name` - (Optional, String) Specifies the name of the environment where the API is published.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `signatures` - All signatures that match the filter parameters.
  The [signatures](#api_associated_throttling_signatures) structure is documented below.

<a name="api_associated_throttling_signatures"></a>
The `signatures` block supports:

* `id` - The ID of the signature.

* `name` - The name of the signature.

* `type` - The type of the signature.

* `key` - The signature key.

* `secret` - The signature secret.

* `env_id` - The ID of the environment where the API is published.

* `env_name` - The name of the environment where the API is published.

* `bind_id` - The bind ID.

* `bind_time` - The time that the signature is bound to the API.
