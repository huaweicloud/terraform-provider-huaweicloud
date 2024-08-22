---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_domain"
description: ""
---

# huaweicloud_live_domain

Manages a Live domain within HuaweiCloud.

## Example Usage

### Create an ingest domain name and a streaming domain name

```hcl
variable "ingest_domain_name" {}
variable "streaming_domain_name" {}

resource "huaweicloud_live_domain" "ingestDomain" {
  name = var.ingest_domain_name
  type = "push"
}

resource "huaweicloud_live_domain" "streamingDomain" {
  name               = var.streaming_domain_name
  type               = "pull"
  ingest_domain_name = huaweicloud_live_domain.ingestDomain.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the Live domain resource. If omitted,
the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the domain name. Changing this parameter will create a new resource.

-> A level-1 domain name cannot be used as an ingest domain or streaming domain. If your domain name is **example.com**,
you can use sub-domain names, for example, **test-push.example.com** and **test-play.example.com**,
as the ingest domain name and streaming domain name.

* `type` - (Required, String, ForceNew) Specifies the type of domain name. The options are as follows:
  + **pull**: streaming domain name.
  + **push**: ingest domain name.

  Changing this parameter will create a new resource.

* `ingest_domain_name` - (Optional, String) Specifies the ingest domain name, which associates with the streaming
domain name to push streams to nearby CDN nodes.

* `status` - (Optional, String) Specifies status of the domain name. The options are as follows:
  + **on**: enable the domain name.
  + **off**: disable the domain name.

  The default value is `on`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals to domain name.

* `cname` - CNAME record of the domain name.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `update` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

Domains can be imported using the `name`, e.g.

```bash
$ terraform import huaweicloud_live_domain.test domainName
```
