# Package antiddos
    import "github.com/chnsz/golangsdk/openstack/antiddos/v2/alarmreminding"
**[概述](#概述)**  

**[目录](#目录)**  

**[API对应表](#API对应表)**  

**[开始](#开始)**  

## 概述
Anti-DDoS流量清洗（Anti-DDoS Service）是通过专业的防DDoS设备来为客户互联网应用提供精细化的抵御DDoS攻击能力（包括CC、SYN flood、UDP flood等所有DDoS攻击方式）。可根据租用带宽及业务模型自助配置防护阈值参数，系统检测到攻击后通知用户进行网站防御。

示例代码，查询用户配置信息。

    actual, err := alarmreminding.WarnAlert(client.ServiceClient()).Extract()
    if err != nil {
      panic(err)
    }
## 目录
**[func WarnAlert(*golangsdk.ServiceClient) (WarnAlertResult)](#func-warnalert)**  
## API对应表
|类别|API|EndPoint|
|----|---|--------|
|antiddos|func WarnAlert(*golangsdk.ServiceClient) (WarnAlertResult)|GET /v2/{project_id}/warnalert/alertconfig/query|
## 开始
## func WarnAlert
    func WarnAlert(*golangsdk.ServiceClient) (WarnAlertResult)  
查询用户配置信息，用户可以通过此接口查询是否接收某类告警，同时可以配置是手机短信还是电子邮件接收告警信息。
