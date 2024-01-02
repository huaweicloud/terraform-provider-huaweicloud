/*
The Anti-DDoS traffic cleaning service (Anti-DDoS for short) defends resources (Elastic Cloud Servers (ECSs), Elastic Load Balance (ELB) instances, and Bare Metal Servers (BMSs)) on HUAWEI CLOUD against network- and application-layer distributed denial of service (DDoS) attacks and sends alarms immediately when detecting an attack. In addition, Anti-DDoS improves the utilization of bandwidth and ensures the stable running of users' services.

Example to query alarm configuration.

    actual, err := alarmreminding.WarnAlert(client.ServiceClient()).Extract()
    if err != nil {
      panic(err)
    }
*/
package alarmreminding
