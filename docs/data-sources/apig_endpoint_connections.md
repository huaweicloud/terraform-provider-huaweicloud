---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_endpoint_connections"
description: |-
  Use this data source to get the VPC endpoint connections under specified the APIG instance within HuaweiCloud.
---

# huaweicloud_apig_endpoint_connections

Use this data source to get the VPC endpoint connections under specified the APIG instance within HuaweiCloud.

## Example Usage

```hcl
variable instance_id {}

data "huaweicloud_apig_endpoint_connections" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the dedicated instance ID to which the endpoint connections belong.

* `endpoint_id` - (Optional, String) Specifies the ID of the endpoint connection.

* `packet_id` - (Optional, Int) Specifies packet ID of endpoint connection.

* `status` - (Optional, String) Specifies status of endpoint connection.
  The valid values are as follows:
  + **pendingAcceptance**
  + **accepted**
  + **rejected**
  + **failed**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `connections` - All endpoint connections that match the filter parameters.

  The [connections](#connections_struct) structure is documented below.

<a name="connections_struct"></a>
The `connections` block supports:

* `id` - The ID of the endpoint connection.

* `packet_id` - The packet ID of the endpoint connection.

* `domain_id` - The IAM account ID of the endpoint connection creator.

* `status` - The current status of the endpoint connection.

* `created_at` - The creation time of the endpoint connection, in RFC3339 format.

* `updated_at` - The latest time of the endpoint connection, in RFC3339 format.
