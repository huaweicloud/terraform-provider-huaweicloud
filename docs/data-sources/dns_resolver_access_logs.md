---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_resolver_access_logs"
description: |-
  Use this data source to get the list of DNS resolver access logs within HuaweiCloud.
---

# huaweicloud_dns_resolver_access_logs

Use this data source to get the list of DNS resolver access logs within HuaweiCloud.

## Example Usage

### Query all resolver access logs

```hcl
data "huaweicloud_dns_resolver_access_logs" "test" {}
```

### Query resolver access logs by the specified VPC ID

```hcl
variable "vpc_id" {}

data "huaweicloud_dns_resolver_access_logs" "test" {
  vpc_id = var.vpc_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the resolver access logs are located.  
  If omitted, the provider-level region will be used.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `access_logs` - The list of resolver access logs that match the filter parameters.  
  The [access_logs](#resolver_access_logs_struct) structure is documented below.

<a name="resolver_access_logs_struct"></a>
The `access_logs` block supports:

* `id` - The ID of the resolver access log.

* `lts_group_id` - The ID of the log group.

* `lts_topic_id` - The ID of the log stream.

* `vpc_ids` - The list of VPC IDs associated with the resolver access log.
