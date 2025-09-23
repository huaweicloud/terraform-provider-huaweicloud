---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_ips_custom_rules"
description: |-
  Use the data source to get the list of CFW IPS custom rules.
---

# huaweicloud_cfw_ips_custom_rules

Use the data source to get the list of CFW IPS custom rules.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "object_id" {}

data "huaweicloud_cfw_ips_custom_rules" "test" {
  fw_instance_id = var.fw_instance_id
  object_id      = var.object_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall ID.

* `object_id` - (Required, String) Specifies the protected object ID.

* `action_type` - (Optional, Int) Specifies the action type.
  The valid value can be **0** (log only) or **1** (reset/block).

* `affected_os` - (Optional, Int) Specifies the affected OS.
  The valid values are as follows:
  + **1**: Windows;
  + **2**: Linux;
  + **3**: FreeBSD;
  + **4**: Solaris;
  + **5**: Other Unix;
  + **6**: Network device;
  + **7**: MAC OS;
  + **8**: IOS;
  + **9**: Android;
  + **10**: Other;

* `attack_type` - (Optional, Int) Specifies the attack type.
  The valid values are as follows:
  + **1**: access control;
  + **2**: vulnerability scan;
  + **3**: email phishing;
  + **4**: vulnerability exploits;
  + **5**: web attack;
  + **6**: password cracking;
  + **7**: hijacking attack;
  + **8**: protocol exception;
  + **9**: trojan;
  + **10**: worm;
  + **11**: buffer overflow;
  + **12**: hacker tool;
  + **13**: spyware;
  + **14**: DDoS flood;
  + **15**: application-layer DDoS attack;
  + **16**: other suspicious behavior;
  + **17**: suspicious DNS activity;
  + **18**: phishing;
  + **19**: spam;

* `ips_name` - (Optional, String) Specifies the IPS custom rule name.

* `protocol` - (Optional, Int) Specifies the protocol.
  The valid values are as follows:
  + **1**: FTP;
  + **2**: TELNET;
  + **3**: SMTP;
  + **4**: DNS-TCP;
  + **5**: DNS-UDP;
  + **6**: DHCP;
  + **7**: TFTP;
  + **8**: FINGER;
  + **9**: HTTP;
  + **10**: POP3;
  + **11**: SUNRPC-TCP;
  + **12**: SUNRPC-UDP;
  + **13**: NNTP;
  + **14**: MSRPC-TCP;
  + **15**: MSRPC-UDP;
  + **16**: NETBIOS-NAME_TCP;
  + **17**: NETBIOS-NAME_UDP;
  + **18**: NETBIOS-SMB;
  + **19**: NETBIOS-DATAGRAM;
  + **20**: IMAP4;
  + **21**: SNMP;
  + **22**: LDAP;
  + **23**: MSSQL;
  + **24**: ORACLE;
  + **25**: MYSQL;
  + **26**: VOIP-SIP-TCP;
  + **27**: VOIP-SIP-UDP;
  + **28**: VOIP-H245;
  + **29**: VOIP-Q931;
  + **30**: OTHER-TCP;
  + **31**: OTHER-UDP;

* `severity` - (Optional, Int) Specifies the severity.
  The valid values are as follows:
  + **0**: critical;
  + **1**: high;
  + **2**: medium;
  + **3**: low;

* `software` - (Optional, Int) Specifies the affected software.
  The valid values are as follows:
  + **1**: ADOBE;
  + **2**: APACHE;
  + **3**: APPLE;
  + **4**: CA;
  + **5**: CISCO;
  + **6**: GOOGLE CHROME;
  + **7**: HP;
  + **8**: IBM;
  + **9**: IE;
  + **10**: IIS;
  + **11**: MCAFEE;
  + **12**: MEDIAPLAYER;
  + **13**: MICROSOFT.NET;
  + **14**: MICROSOFT EDGE;
  + **15**: MICROSOFT EXCHANGE;
  + **16**: MICROSOFT OFFICE;
  + **17**: MICROSOFT OUTLOOK;
  + **18**: MICROSOFT SHAREPOINT;
  + **19**: MICROSOFT WINDOWS;
  + **20**: MOZILLA;
  + **21**: MSSQL;
  + **22**: MYSQL;
  + **23**: NOVELL;
  + **24**: ORACLE;
  + **25**: SAMBA;
  + **26**: SAMSUNG;
  + **27**: SAP;
  + **28**: SCADA;
  + **29**: SQUID;
  + **30**: SUN;
  + **31**: SYMANTEC;
  + **32**: TREND MICRO;
  + **33**: VMWARE;
  + **34**: WORDPRESS;
  + **35**: OTHER;

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The custom IPS rule records.

  The [records](#data_records_struct) structure is documented below.

<a name="data_records_struct"></a>
The `records` block supports:

* `severity` - The severity.

* `software` - The affected software.

* `src_ports` - The source port.

* `dst_port_type` - The destination port type.

* `content` - The content storage in JSON format.

* `dst_ports` - The destination port.

* `attack_type` - The attack type.

* `src_port_type` - The source port type.

* `protocol` - The protocol.

* `affected_os` - The affected OS.

* `config_status` - The rule status.

* `group_id` - The firewall cluster ID.

* `ips_cfw_id` - The ID of a custom IPS rule in CFW.

* `ips_id` - The ID of a rule in Hillstone.

* `ips_name` - The IPS rule name.

* `action` - The action.
