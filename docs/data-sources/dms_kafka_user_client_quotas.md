---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_user_client_quotas"
description: |-
  Use this data source to get the list of Kafka instance user client quotas.
---

# huaweicloud_dms_kafka_user_client_quotas

Use this data source to get the list of Kafka instance user client quotas.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_kafka_user_client_quotas" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `user` - (Optional, String) Specifies the user name.

* `client` - (Optional, String) Specifies the client ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - Indicates the client quotas.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `user` - Indicates the username.

* `user_default` - Indicates whether to use the default user settings.

* `client` - Indicates the client ID.

* `client_default` - Indicates whether to use the default client settings.

* `producer_byte_rate` - Indicates the production rate limit. The unit is byte/s.

* `consumer_byte_rate` - Indicates the consumption rate limit. The unit is byte/s.
