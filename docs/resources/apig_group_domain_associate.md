---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_group_domain_associate"
description: |-
  Use this resource to associate a domain name with an API group.
---

# huaweicloud_apig_group_domain_associate

Use this resource to associate a domain name with an API group.

~> This resource cannot be used simultaneously with the domain name management parameter (`url_domains`) of the group
   resource.<br>Using `lifecycle.ignore_changes` to ignore changes for the corresponding group when using this resource.

-> Different groups under the same dedicated instance cannot be bound to the same independent domain name.
   Each API group can be associated with up to `5` domain names.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_id" {}
variable "domain_name" {}

resource "huaweicloud_apig_group_domain_associate" "test" {
  instance_id = var.instance_id
  group_id    = var.group_id
  url_domain  = var.domain_name

  min_ssl_version           = "TLSv1.1"
  ingress_http_port         = 80
  ingress_https_port        = -1
  is_http_redirect_to_https = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the APIG (API) group is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the dedicated instance to which the group belongs.

* `group_id` - (Required, String, NonUpdatable) Specifies the ID of the group to associate with the domain name.

* `url_domain` - (Required, String, NonUpdatable) Specifies the associated domain name.
  The maximum valid length is `255` characters and the value must conform to the domain name specifications (the regular
  expression **^(a-zA-Z0-9?\.){1,7}\[a-zA-Z\]{2,64}\.?$** or the regular expression
  **^*{1,6}\.\[a-zA-Z\]{2,64}.?$**).

* `min_ssl_version` - (Optional, String) Specifies the minimum SSL protocol version.  
  The valid values are as follows:
  + **TLSv1.1**
  + **TLSv1.2**

  The default value is different in different regions. For the specific default value, please consult the corresponding
  service OnCall through the [service ticket](https://support.huaweicloud.com/intl/en-us/usermanual-ticket/topic_0065264094.html).

* `ingress_http_port` - (Optional, Int) Specifies the HTTP protocol inbound access port bound to the domain name.  
  The valid value is `-1` (disable), `80` or custom HTTP inbound access port (must be opened by the instance, ranges
  from `1,024` to `49,151`).

* `ingress_https_port` - (Optional, Int) Specifies the HTTPS protocol inbound access port bound to the domain name.  
  The valid value is `-1` (disable), `443` or custom HTTPS inbound access port (must be opened by the instance, ranges
  from `1,024` to `49,151`).

-> `ingress_http_port` and `ingress_https_port` cannot be disabled at the same time.

* `is_http_redirect_to_https` - (Optional, Bool) Specifies whether to enable redirection from `HTTP` to `HTTPS`.
  Defaults to `false`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, the format is `<instance_id>/<group_id>/<url_domain>`.

## Import

Associated information for the specified domain and group can be imported using resource `id` (consists of
`instance_id`, `group_id` and `url_domain`, separated by the slashes (/)), e.g.

```bash
$ terraform import huaweicloud_apig_group_domain_associate.test <instance_id>/<group_id>/<url_domain>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response.
The missing attributes includes: `is_http_redirect_to_https`.
It is generally recommended running `terraform plan` after importing this resource.
You can then decide if changes should be applied to the resource, or the definition should be updated to align with the
resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_apig_group_domain_associate" "test" {
  ...

  lifecycle {
    ignore_changes = [
      is_http_redirect_to_https,
    ]
  }
}
```
