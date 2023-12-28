# Package antiddos
    import "github.com/chnsz/golangsdk/openstack/antiddos/v2/alarmreminding"
**[Overview](#overview)**  

**[Index](#index)**  

**[API Mapping](#api-mapping)**  

**[Content](#content)**  

## Overview
The Anti-DDoS traffic cleaning service (Anti-DDoS for short) defends resources (Elastic Cloud Servers (ECSs), Elastic Load Balance (ELB) instances, and Bare Metal Servers (BMSs)) on HUAWEI CLOUD against network- and application-layer distributed denial of service (DDoS) attacks and sends alarms immediately when detecting an attack. In addition, Anti-DDoS improves the utilization of bandwidth and ensures the stable running of users' services.

Example to query alarm configuration.

    actual, err := alarmreminding.WarnAlert(client.ServiceClient()).Extract()
    if err != nil {
      panic(err)
    }
## Index
**[func WarnAlert(*golangsdk.ServiceClient) (WarnAlertResult)](#func-warnalert)**  
## API Mapping
|Catalog|API|EndPoint|
|----|---|--------|
|antiddos|func WarnAlert(*golangsdk.ServiceClient) (WarnAlertResult)|GET /v2/{project_id}/warnalert/alertconfig/query|
## Content
## func WarnAlert
    func WarnAlert(*golangsdk.ServiceClient) (WarnAlertResult)  
This API allows you to query alarm configuration, such as whether a certain type of alarms will be received, and whether alarms are received through SMS messages or emails.
