---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_domain_v1"
sidebar_current: "docs-huaweicloud-resource-cdn-domain-v1"
description: |-
  cdn domain management
---

# huaweicloud\_cdn\_domain\_v1

cdn domain management

## Example Usage

### create a cdn domain

```hcl
resource "huaweicloud_cdn_domain_v1" "domain_1" {
  name = "${var.domain_name}"
  type = "web"

  sources {
    domain      = "${var.origin_server}"
    domain_type = "ipaddr"
    active      = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The acceleration domain name.
    Changing this parameter will create a new resource.

* `type` - (Required) The service type. The valid values are  'web', 'download' and 'video'.
    Changing this parameter will create a new resource.

* `sources` - (Required) An array of one or more objects specifies the domain name of the origin server.
    The sources object structure is documented below. Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional) The enterprise project id.
    Changing this parameter will create a new resource.


The `sources` block supports:

* `domain` - (Required) The domain name or IP address of the origin server.
    Changing this parameter will create a new resource.

* `domain_type` - (Required) The origin server type. The valid values are 'ipaddr', 'domain', and 'obs_bucket'.
    Changing this parameter will create a new resource.

* `active` - (Optional) Whether an origin server is active or standby (1: active; 0: standby).
    The default value is 1. Changing this parameter will create a new resource.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.

* `type` - See Argument Reference above.

* `enterprise_project_id` - See Argument Reference above.

* `sources/domain` - See Argument Reference above.

* `sources/domain_type` - See Argument Reference above.

* `sources/active` - See Argument Reference above.

* `id` - The acceleration domain name ID.

* `cname` - The CNAME of the acceleration domain name.

* `domain_status` - The status of the acceleration domain name. The available values are
    'online', 'offline', 'configuring', 'configure_failed', 'checking', 'check_failed'  and 'deleting.'

* `service_area` - The area covered by the acceleration service.


## Import

Domains can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_cdn_domain_v1.domain_1 fe2462fac09a4a42a76ecc4a1ef542f1
```