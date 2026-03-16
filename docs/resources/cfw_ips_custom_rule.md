---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_ips_custom_rule"
description: |-
  Manages a CFW IPS custom rule resource within HuaweiCloud.
---

# huaweicloud_cfw_ips_custom_rule

Manages a CFW IPS custom rule resource within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}

resource "huaweicloud_cfw_ips_custom_rule" "test" {
  fw_instance_id = var.fw_instance_id
  ips_name       = "test-name"
  action_type    = 0
  affected_os    = 0
  attack_type    = 3
  direction      = 1
  protocol       = 10
  severity       = 1
  software       = 3

  contents {
    content           = "vvvvvv"
    depth             = 65535
    is_hex            = false
    is_ignore         = true
    is_uri            = false
    offset            = 50
    relative_position = 0
  }

  contents {
    content           = "DF"
    depth             = 65535
    is_hex            = true
    is_ignore         = false
    is_uri            = false
    offset            = 200
    relative_position = 0
  }

  dst_port {
    port_type = 1
    ports     = "9008"
  }

  src_port {
    port_type = 0
    ports     = "5005"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `fw_instance_id` - (Required, String, NonUpdatable) Specifies the firewall ID.
  It is a unique ID generated after a firewall instance is created. You can obtain the firewall ID by referring to
  the data source `huaweicloud_cfw_firewalls`.

* `ips_name` - (Required, String, NonUpdatable) Specifies the IPS custom rule name.

* `protocol` - (Required, Int, NonUpdatable) Specifies the protocol type. Valid values are:
  + `1`: FTP
  + `2`: TELNET
  + `3`: SMTP
  + `4`: DNS_TCP
  + `5`: DNS_UDP
  + `6`: DHCP
  + `7`: TFTP
  + `8`: FINGER
  + `9`: HTTP
  + `10`: POP3
  + `11`: SUNRPC_TCP
  + `12`: SUNRPC_UDP
  + `13`: NNTP
  + `14`: MSRPC_TCP
  + `15`: MSRPC_UDP
  + `16`: NETBIOS_NAME_TCP
  + `17`: NETBIOS_NAME_UDP
  + `18`: NETBIOS_SMB
  + `19`: NETBIOS_DATAGRAM
  + `20`: IMAP4
  + `21`: SNMP
  + `22`: LDAP
  + `23`: MSSQL
  + `24`: ORACLE

* `action_type` - (Required, Int) Specifies the action type. Valid values are:
  + `0`: Used only for logging
  + `1`: Reset/Intercept

* `affected_os` - (Required, Int) Specifies the affected operating system. Valid values are:
  + `0`: ANY
  + `1`: Windows
  + `2`: Linux
  + `3`: FreeBSD
  + `4`: Solaris
  + `5`: Other Unix
  + `6`: Network equipment
  + `7`: Mac OS
  + `8`: IOS
  + `9`: Android
  + `10`: Others

* `attack_type` - (Required, Int) Specifies the attack type. Valid values are:
  + `1`: Access control
  + `2`: Vulnerability scan
  + `3`: Email attack
  + `4`: Vulnerability attack
  + `5`: Web attack
  + `6`: Password attack
  + `7`: Hijacking attack
  + `8`: Protocol anomaly
  + `9`: Trojan horse
  + `10`: Worm
  + `11`: Buffer overflow
  + `12`: Hacker tool
  + `13`: Spyware
  + `14`: DDoS flood
  + `15`: Application layer DDoS attack
  + `16`: Other suspicious behavior
  + `17`: Suspicious DNS activity
  + `18`: Network fishing
  + `19`: Spam email
  + `20`: Other attack

* `contents` - (Required, List) Specifies the the message content to match IPS attacks.

  The [contents](#contents_struct) structure is documented below.

* `direction` - (Required, Int) Specifies the direction. Valid values are:
  + `-1`: All directions
  + `0`: Client to server
  + `1`: Server to client

* `dst_port` - (Required, List) Specifies the destination port information.
  A maximum of one set of elements can be configured.

  The [dst_port](#port_struct) structure is documented below.

* `severity` - (Required, Int) Specifies the severity. Valid value are:
  + `0`: Critical
  + `1`: High
  + `2`: Medium
  + `3`: Low

* `software` - (Required, Int) Specifies the software affected. Valid value are:
  + `0`: ANY
  + `1`: ADOBE
  + `2`: APACHE
  + `3`: APPLE
  + `4`: CA
  + `5`: CISCO
  + `6`: GOOGLE_CHROME
  + `7`: HP
  + `8`: IBM
  + `9`: IE
  + `10`: IIS
  + `11`: MC_AFEE
  + `12`: MEDIA_PLAYER
  + `13`: MICROSOFT_NET
  + `14`: MICROSOFT_EDGE
  + `15`: MICROSOFT_EXCHANGE
  + `16`: MICROSOFT_OFFICE
  + `17`: MICROSOFT_OUTLOOK
  + `18`: MICROSOFT_SHARE_POINT
  + `19`: MICROSOFT_WINDOWS
  + `20`: MOZILLA
  + `21`: MSSQL
  + `22`: MYSQL
  + `23`: NOVELL
  + `24`: ORACLE
  + `25`: SAMBA
  + `26`: SAMSUNG
  + `27`: SAP
  + `28`: SCADA
  + `29`: SQUID
  + `30`: SUN
  + `31`: SYMANTEC
  + `32`: TREND_MICRO
  + `33`: VMWARE
  + `34`: WORD_PRESS
  + `35`: Others

* `src_port` - (Required, List) Specifies the source port information.
  A maximum of one set of elements can be configured.

  The [src_port](#port_struct) structure is documented below.

<a name="port_struct"></a>
The `dst_port` and `src_port` block supports:

* `port_type` - (Required, Int) Specifies the port type. Valid values are:
  + `-1`: All ports
  + `0`: Include ports
  + `1`: Exclude ports

* `ports` - (Optional, String) Specifies the port.

<a name="contents_struct"></a>
The `contents` block supports:

* `content` - (Required, String) Specifies the content.

* `depth` - (Required, Int) Specifies the position to end the matching when matching features. Valid values are
  from `1` to `65535`.

* `is_hex` - (Optional, Bool) Specifies whether to enable hexadecimal matching. Valid values are:
  + **true**: Enable hexadecimal matching
  + **false**: Disable hexadecimal matching

  Defaults to **false**.

* `is_ignore` - (Optional, Bool) Specifies whether to ignore case. Valid values are:
  + **true**: Ignore case
  + **false**: Do not ignore case

  Defaults to **false**.

* `is_uri` - (Optional, Bool) Specifies whether to match a field in the URL that is the same as `content`.
  Valid values are:
  + **true**: Match
  + **false**: Does not match

  Defaults to **false**.

* `offset` - (Optional, Int) Specifies the starting position when matching features.
  Valid values are from `0` to `65535`.

* `relative_position` - (Optional, Int) Specifies the starting position. Valid values are `0` and `1`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also the IPS custom rule ID).

* `config_status` - The configuration status. Valid values are:
  + `0`: Unknown
  + `1`: Configuring
  + `2`: Configured
  + `3`: Configuration failed

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The resource can be imported using `fw_instance_id`, `id` (IPS custom rule ID), separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cfw_ips_custom_rule.test <fw_instance_id>/<id>
```
