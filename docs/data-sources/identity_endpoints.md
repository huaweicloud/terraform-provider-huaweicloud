---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_endpoints"
description: |-
  Use this data source to get information of endpoint within HuaweiCloud.
---

# huaweicloud_identity_endpoints

Use this data source to get information of endpoint within HuaweiCloud.

## Example Usage

### query by endpoint ID

```hcl
variable "endpoint_id" {}

data "huaweicloud_identity_endpoints" "test" {
  endpoint_id = var.endpoint_id
}
```

### query all endpoints

```hcl
data "huaweicloud_identity_endpoints" "test" {}
```

## Argument Reference

The following arguments are supported:

* `endpoint_id` - (Optional, String) Specifies the terminal node ID to be queried.

* `interface` - (Optional, String) Specifies terminal node plane. It's not allowed if `endpoint_id` is specified.

* `service_id` - (Optional, String) Specifies the ID of service. It's not allowed if `endpoint_id` is specified.

## Attribute Reference

* `endpoints` - Indicates service endpoint information.
  The [endpoint](#IdentityEndpoint_endpoints) structure is documented below.

<a name="IdentityEndpoint_endpoints"></a>
The `endpoints` block contains:

* `service_id` - Indicates the service ID.

* `region_id` - Indicates the region ID.

* `link` - Indicates the resource URL information.

* `id` - Indicates the endpoint ID.

* `interface` - Indicates the terminal node plane.

* `region` - Indicates the region where the endpoint is located.

* `url` - Indicates the URL or network address of the endpoint.

* `enabled` - Indicates whether the endpoint is accessible and operational. The value can be:
  + **true**: accessible;
  + **false**: operational;
