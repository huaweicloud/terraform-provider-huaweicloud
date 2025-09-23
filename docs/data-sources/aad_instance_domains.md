---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_instance_domains"
description: |-
  Use this data source to get the list of Advanced Anti-DDos instance domain information within HuaweiCloud.
---

# huaweicloud_aad_instance_domains

Use this data source to get the list of Advanced Anti-DDos instance domains information within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_aad_instance_domains" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `instance_name` - The instance name.

* `domains` - The domain information list.  
  The [domains](#domains_struct) structure is documented below.

<a name="domains_struct"></a>
The `domains` block supports:

* `domain_id` - The domain ID.

* `domain_name` - The domain name.

* `cname` - The domain CNAME.

* `domain_status` - The domain status. `0` represents normal, `1` represents freeze.

* `cc_status` - The CC protection status.

* `https_cert_status` - The certificate status. `1` represents uploaded, `2` represents not uploaded.

* `cert_name` - The certificate name.

* `protocol_type` - The domain protocol list.

* `real_server_type` - The real server type.

* `real_servers` - The real servers.

* `waf_status` - The WAF protection status.
