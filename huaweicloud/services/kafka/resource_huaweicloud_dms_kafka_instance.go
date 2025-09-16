package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dms/v2/availablezones"
	"github.com/chnsz/golangsdk/openstack/dms/v2/kafka/instances"
	"github.com/chnsz/golangsdk/openstack/dms/v2/products"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type ctxType string

const engineKafka = "kafka"

// @API Kafka GET /v2/available-zones
// @API Kafka POST /v2/{project_id}/instances/{instance_id}/crossvpc/modify
// @API Kafka POST /v2/{project_id}/instances/{instance_id}/extend
// @API Kafka DELETE /v2/{project_id}/instances/{instance_id}
// @API Kafka GET /v2/{project_id}/instances/{instance_id}
// @API Kafka PUT /v2/{project_id}/instances/{instance_id}
// @API Kafka POST /v2/{engine}/{project_id}/instances
// @API Kafka GET /v2/{project_id}/kafka/{instance_id}/tags
// @API Kafka POST /v2/{project_id}/kafka/{instance_id}/tags/action
// @API Kafka POST /v2/{project_id}/instances/{instance_id}/autotopic
// @API Kafka GET /v2/{project_id}/instances/{instance_id}/tasks
// @API Kafka GET /v2/{project_id}/instances/{instance_id}/tasks/{task_id}
// @API Kafka POST /v2/{project_id}/instances/{instance_id}/password
// @API Kafka PUT /v2/{project_id}/instances/{instance_id}/configs
// @API Kafka GET /v2/{project_id}/instances/{instance_id}/configs
// @API Kafka POST /v2/{project_id}/instances/action
// @API Kafka POST /v2/{project_id}/{engine}/instances/{instance_id}/plain-ssl-switch
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourceDmsKafkaInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaInstanceCreate,
		ReadContext:   resourceDmsKafkaInstanceRead,
		UpdateContext: resourceDmsKafkaInstanceUpdate,
		DeleteContext: resourceDmsKafkaInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_spec_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipv6_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"availability_zones": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "schema: Required",
			},
			"arch_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"flavor_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"product_id"},
				RequiredWith: []string{"broker_num", "storage_space"},
			},
			"broker_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"new_tenant_ips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"product_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_space": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"access_user": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
				ConflictsWith: []string{
					"kms_encrypted_password",
				},
			},
			"kms_encrypted_password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
				ConflictsWith: []string{
					"password",
				},
				Description: "schema: Internal",
			},
			// The API return format is "HH:mm:ss" for `maintain_begin` and `maintain_end`.
			"maintain_begin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(_, o, n string, _ *schema.ResourceData) bool {
					log.Printf("[DEBUG] maintain_begin DiffSuppressFunc: %s, %s", o, n)
					log.Printf("[DEBUG] maintain_begin: %v", regexp.MustCompile(fmt.Sprintf("^%s", n)).MatchString(o))
					return regexp.MustCompile(fmt.Sprintf("^%s", n)).MatchString(o)
				},
			},
			"maintain_end": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(_, o, n string, _ *schema.ResourceData) bool {
					return regexp.MustCompile(fmt.Sprintf("^%s", n)).MatchString(o)
				},
			},
			"public_ip_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enabled_mechanisms": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"retention_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"produce_reject", "time_base",
				}, false),
			},
			"dumping": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enable_auto_topic": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssl_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_client_plain": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": common.TagsSchema(),
			"cross_vpc_accesses": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advertised_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"listener_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// Typo, it is only kept in the code, will not be shown in the docs.
						"lisenter_ip": {
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "typo in lisenter_ip, please use \"listener_ip\" instead.",
						},
					},
				},
			},
			"port_protocol": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private_plain_enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `Whether to enable private plaintext access.`,
						},
						"private_sasl_ssl_enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `Whether to enable private SASL SSL access.`,
						},
						"private_sasl_plaintext_enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `Whether to enable private SASL plaintext access.`,
						},
						"public_plain_enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `Whether to enable public plaintext access.`,
						},
						"public_sasl_ssl_enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `Whether to enable public SASL SSL access.`,
						},
						"public_sasl_plaintext_enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `Whether to enable public SASL plaintext access.`,
						},
						"private_plain_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The address of the private plaintext access.`,
						},
						"private_plain_domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain name of the private plaintext access.`,
						},
						"private_sasl_ssl_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The address of the private SASL SSL access.`,
						},
						"private_sasl_ssl_domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain name of the private SASL SSL access.`,
						},
						"private_sasl_plaintext_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The address of the private SASL plaintext access.`,
						},
						"private_sasl_plaintext_domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain name of the private SASL plaintext access.`,
						},
						"public_plain_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The address of the public plaintext access.`,
						},
						"public_plain_domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain name of the public plaintext access.`,
						},
						"public_sasl_ssl_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The address of the public SASL SSL access.`,
						},
						"public_sasl_ssl_domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain name of the public SASL SSL access.`,
						},
						"public_sasl_plaintext_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The address of the public SASL plaintext access.`,
						},
						"public_sasl_plaintext_domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain name of the public SASL plaintext access.`,
						},
					},
				},
				Description: `The port protocol information of the Kafka instance.`,
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),
			// Attributes.
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"partition_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enable_public_ip": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"public_ip_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"used_storage_space": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"connect_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extend_times": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ipv6_connect_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"connector_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connector_node_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"storage_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert_replaced": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_logical_volume": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"message_query_inst_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"node_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pod_connect_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_bandwidth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ssl_two_way_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Deprecated parameters.
			"manager_user": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Deprecated",
			},
			"manager_password": {
				Type:       schema.TypeString,
				Optional:   true,
				Sensitive:  true,
				ForceNew:   true,
				Deprecated: "Deprecated",
			},
			"available_zones": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
				AtLeastOneOf: []string{"available_zones", "availability_zones"},
				Deprecated:   "available_zones has deprecated, please use \"availability_zones\" instead.",
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Deprecated: "The bandwidth has been deprecated. " +
					"If you need to change the bandwidth, please update the product_id.",
			},
			// Deprecated attributes.
			"management_connect_address": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Deprecated",
			},
			// Typo, it is only kept in the code, will not be shown in the docs.
			"manegement_connect_address": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "typo in manegement_connect_address, please use \"management_connect_address\" instead.",
			},
			"port_protocols": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private_plain_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"private_plain_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_plain_domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_sasl_ssl_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"private_sasl_ssl_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_sasl_ssl_domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_sasl_plaintext_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"private_sasl_plaintext_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_sasl_plaintext_domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_plain_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"public_plain_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_plain_domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_sasl_ssl_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"public_sasl_ssl_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_sasl_ssl_domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_sasl_plaintext_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"public_sasl_plaintext_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_sasl_plaintext_domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Description: utils.SchemaDesc("Use port_protocol instead.",
					utils.SchemaDescInput{
						Deprecated: true,
					}),
			},
		},
	}
}

func validateAndBuildPublicIpIDParam(publicIpIDs []interface{}, bandwidth string) (string, error) {
	bandwidthAndIPMapper := map[string]int{
		"100MB":  3,
		"300MB":  3,
		"600MB":  4,
		"1200MB": 8,
	}
	needIpCount := bandwidthAndIPMapper[bandwidth]

	if needIpCount != len(publicIpIDs) {
		return "", fmt.Errorf("error creating Kafka instance: "+
			"%d public ip IDs needed when bandwidth is set to %s, but got %d",
			needIpCount, bandwidth, len(publicIpIDs))
	}
	return strings.Join(utils.ExpandToStringList(publicIpIDs), ","), nil
}

func getProducts(cfg *config.Config, region, engine string) (*products.GetResponse, error) {
	dmsV2Client, err := cfg.DmsV2Client(region)
	if err != nil {
		return nil, fmt.Errorf("error getting DMS product client V2: %s", err)
	}
	v, err := products.Get(dmsV2Client, engine) // nolint: staticcheck
	return v, err
}

func getKafkaProductDetails(cfg *config.Config, d *schema.ResourceData) (*products.Detail, error) {
	productRsp, err := getProducts(cfg, cfg.GetRegion(d), engineKafka)
	if err != nil {
		return nil, fmt.Errorf("error querying Kafka product list: %s", err)
	}

	productID := d.Get("product_id").(string)
	engineVersion := d.Get("engine_version").(string)

	for _, ps := range productRsp.Hourly {
		if ps.Version != engineVersion {
			continue
		}
		for _, v := range ps.Values {
			for _, p := range v.Details {
				if p.ProductID == productID {
					return &p, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("can not found Kafka product details base on product_id: %s", productID)
}

func UpdateCrossVpcAccess(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	newVal := d.Get("cross_vpc_accesses")
	var crossVpcAccessArr []map[string]interface{}

	instance, err := instances.Get(client, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("error getting DMS Kafka instance: %v", err)
	}

	crossVpcAccessArr, err = FlattenCrossVpcInfo(instance.CrossVpcInfo)
	if err != nil {
		return fmt.Errorf("error retrieving details of the cross-VPC: %v", err)
	}

	newAccessArr := newVal.([]interface{})
	contentMap := make(map[string]string)
	for i, oldAccess := range crossVpcAccessArr {
		listenerIp := oldAccess["listener_ip"].(string)
		if listenerIp == "" {
			listenerIp = oldAccess["lisenter_ip"].(string)
		}
		// If we configure the advertised ip as ["192.168.0.19", "192.168.0.8"], the length of new accesses is 2,
		// and the length of old accesses is always 3.
		if len(newAccessArr) > i {
			// Make sure the index is valid.
			newAccess := newAccessArr[i].(map[string]interface{})
			// Since the "advertised_ip" is already a definition in the schema, the key name must exist.
			if advIp, ok := newAccess["advertised_ip"].(string); ok && advIp != "" {
				contentMap[listenerIp] = advIp
				continue
			}
		}
		contentMap[listenerIp] = listenerIp
	}

	log.Printf("[DEBUG} Update Kafka cross-vpc contentMap: %#v", contentMap)

	retryFunc := func() (interface{}, bool, error) {
		updateRst, err := instances.UpdateCrossVpc(client, d.Id(), instances.CrossVpcUpdateOpts{
			Contents: contentMap,
		})
		retry, err := handleMultiOperationsError(err)
		return updateRst, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     KafkaInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating advertised IP: %v", err)
	}
	updateRst := r.(*instances.CrossVpc)

	if !updateRst.Success {
		failedIps := make([]string, 0, len(updateRst.Connections))
		for _, conn := range updateRst.Connections {
			if !conn.Success {
				failedIps = append(failedIps, conn.ListenersIp)
			}
		}
		return fmt.Errorf("failed to update the advertised IPs corresponding to some listener IPs (%v)", failedIps)
	}
	return nil
}

func resourceDmsKafkaInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error initializing DMS Kafka(v2) client: %s", err)
	}

	var dErr diag.Diagnostics
	if _, ok := d.GetOk("flavor_id"); ok {
		dErr = createKafkaInstanceWithFlavor(ctx, d, meta)
	} else {
		dErr = createKafkaInstanceWithProductID(ctx, d, meta)
	}
	if dErr != nil {
		return dErr
	}

	// After the kafka instance is created, wait for the access port to complete the binding.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"BOUND"},
		Refresh:      kafkaInstanceCrossVpcInfoRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		dErr = diag.Errorf("Kafka instance is created, but failed to enable cross-VPC %s : %s", d.Id(), err)
		dErr[0].Severity = diag.Warning
		return dErr
	}

	if _, ok := d.GetOk("cross_vpc_accesses"); ok {
		if err = UpdateCrossVpcAccess(ctx, client, d); err != nil {
			return diag.Errorf("failed to update default advertised IP: %s", err)
		}
	}

	if parameters := d.Get("parameters").(*schema.Set); parameters.Len() > 0 {
		if err = initializeParameters(ctx, d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDmsKafkaInstanceRead(ctx, d, meta)
}

func buildKafkaPortProtocol(portProtocols []interface{}) *instances.PortProtocol {
	if len(portProtocols) == 0 {
		return nil
	}

	portProtocol := portProtocols[0].(map[string]interface{})
	return &instances.PortProtocol{
		PrivatePlainEnable:         utils.Bool(portProtocol["private_plain_enable"].(bool)),
		PrivateSaslSslEnable:       utils.Bool(portProtocol["private_sasl_ssl_enable"].(bool)),
		PrivateSaslPlaintextEnable: utils.Bool(portProtocol["private_sasl_plaintext_enable"].(bool)),
		PublicPlainEnable:          utils.Bool(portProtocol["public_plain_enable"].(bool)),
		PublicSaslSslEnable:        utils.Bool(portProtocol["public_sasl_ssl_enable"].(bool)),
		PublicSaslPlaintextEnable:  utils.Bool(portProtocol["public_sasl_plaintext_enable"].(bool)),
	}
}

func createKafkaInstanceWithFlavor(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DmsV2Client(region)
	if err != nil {
		return diag.Errorf("error initializing DMS Kafka(v2) client: %s", err)
	}

	createOpts := &instances.CreateOps{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		Engine:                engineKafka,
		EngineVersion:         d.Get("engine_version").(string),
		AccessUser:            d.Get("access_user").(string),
		VPCID:                 d.Get("vpc_id").(string),
		SecurityGroupID:       d.Get("security_group_id").(string),
		SubnetID:              d.Get("network_id").(string),
		ProductID:             d.Get("flavor_id").(string),
		ArchType:              d.Get("arch_type").(string),
		KafkaManagerUser:      d.Get("manager_user").(string),
		MaintainBegin:         d.Get("maintain_begin").(string),
		MaintainEnd:           d.Get("maintain_end").(string),
		RetentionPolicy:       d.Get("retention_policy").(string),
		ConnectorEnalbe:       d.Get("dumping").(bool),
		EnableAutoTopic:       d.Get("enable_auto_topic").(bool),
		StorageSpecCode:       d.Get("storage_spec_code").(string),
		StorageSpace:          d.Get("storage_space").(int),
		BrokerNum:             d.Get("broker_num").(int),
		EnterpriseProjectID:   cfg.GetEnterpriseProjectID(d),
		SslEnable:             d.Get("ssl_enable").(bool),
		KafkaSecurityProtocol: d.Get("security_protocol").(string),
		SaslEnabledMechanisms: utils.ExpandToStringList(d.Get("enabled_mechanisms").(*schema.Set).List()),
		Ipv6Enable:            d.Get("ipv6_enable").(bool),
		VpcClientPlain:        d.Get("vpc_client_plain").(bool),
		PortProtocol:          buildKafkaPortProtocol(d.Get("port_protocol").([]interface{})),
		TenantIps:             utils.ExpandToStringList(d.Get("new_tenant_ips").([]interface{})),
	}

	if chargingMode, ok := d.GetOk("charging_mode"); ok && chargingMode == "prePaid" {
		var autoRenew bool
		if d.Get("auto_renew").(string) == "true" {
			autoRenew = true
		}
		isAutoPay := true
		createOpts.BssParam = instances.BssParam{
			ChargingMode: d.Get("charging_mode").(string),
			PeriodType:   d.Get("period_unit").(string),
			PeriodNum:    d.Get("period").(int),
			IsAutoRenew:  &autoRenew,
			IsAutoPay:    &isAutoPay,
		}
	}

	if ids, ok := d.GetOk("public_ip_ids"); ok {
		createOpts.EnablePublicIP = true
		createOpts.PublicIpID = strings.Join(utils.ExpandToStringList(ids.(*schema.Set).List()), ",")
	}

	var availableZones []string
	if zoneIDs, ok := d.GetOk("available_zones"); ok {
		availableZones = utils.ExpandToStringList(zoneIDs.([]interface{}))
	} else {
		// convert the codes of the availability zone into ids
		azCodes := d.Get("availability_zones").(*schema.Set)
		availableZones, err = GetAvailableZoneIDByCode(cfg, region, azCodes.List())
		if err != nil {
			return diag.FromErr(err)
		}
	}
	createOpts.AvailableZones = availableZones

	// set tags
	if tagRaw := d.Get("tags").(map[string]interface{}); len(tagRaw) > 0 {
		createOpts.Tags = utils.ExpandResourceTags(tagRaw)
	}
	log.Printf("[DEBUG] Create DMS Kafka instance options: %#v", createOpts)
	// Add password here, so it wouldn't go in the above log entry
	password := d.Get("password").(string)
	if password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			password, err = decryptPasswordWithKmsID(ctx, d, meta)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	createOpts.Password = password
	createOpts.KafkaManagerPassword = d.Get("manager_password").(string)

	kafkaInstance, err := instances.CreateWithEngine(client, createOpts, engineKafka).Extract()
	if err != nil {
		return diag.Errorf("error creating Kafka instance: %s", err)
	}
	instanceID := kafkaInstance.InstanceID

	var delayTime time.Duration = 300
	if chargingMode, ok := d.GetOk("charging_mode"); ok && chargingMode == "prePaid" {
		err = waitForOrderComplete(ctx, d, cfg, client, instanceID)
		if err != nil {
			return diag.FromErr(err)
		}
		delayTime = 5
	}

	log.Printf("[INFO] Creating Kafka instance, ID: %s", instanceID)
	d.SetId(instanceID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"RUNNING"},
		Refresh:      KafkaInstanceStateRefreshFunc(client, instanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        delayTime * time.Second,
		PollInterval: 15 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for Kafka instance (%s) to be ready: %s", instanceID, err)
	}

	return nil
}

func createKafkaInstanceWithProductID(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DmsV2Client(region)
	if err != nil {
		return diag.Errorf("error initializing DMS Kafka(v2) client: %s", err)
	}

	product, err := getKafkaProductDetails(cfg, d)
	if err != nil {
		return diag.Errorf("Error querying product detail: %s", err)
	}

	bandwidth := product.Bandwidth
	defaultPartitionNum, _ := strconv.ParseInt(product.PartitionNum, 10, 64)
	defaultStorageSpace, _ := strconv.ParseInt(product.Storage, 10, 64)

	// check storage
	storageSpace, ok := d.GetOk("storage_space")
	if ok && storageSpace.(int) < int(defaultStorageSpace) {
		return diag.Errorf("storage capacity (storage_space) must be greater than the minimum capacity of the product, "+
			"product capacity is %v, got: %v", defaultStorageSpace, storageSpace)
	}

	var availableZones []string
	if zoneIDs, ok := d.GetOk("available_zones"); ok {
		availableZones = utils.ExpandToStringList(zoneIDs.([]interface{}))
	} else {
		// Convert AZ Codes to AZ IDs
		azCodes := d.Get("availability_zones").(*schema.Set)
		availableZones, err = GetAvailableZoneIDByCode(cfg, region, azCodes.List())
		if err != nil {
			return diag.FromErr(err)
		}
	}

	createOpts := &instances.CreateOps{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		Engine:                engineKafka,
		EngineVersion:         d.Get("engine_version").(string),
		Specification:         bandwidth,
		StorageSpace:          int(defaultStorageSpace),
		PartitionNum:          int(defaultPartitionNum),
		AccessUser:            d.Get("access_user").(string),
		VPCID:                 d.Get("vpc_id").(string),
		SecurityGroupID:       d.Get("security_group_id").(string),
		SubnetID:              d.Get("network_id").(string),
		AvailableZones:        availableZones,
		ArchType:              d.Get("arch_type").(string),
		ProductID:             d.Get("product_id").(string),
		KafkaManagerUser:      d.Get("manager_user").(string),
		MaintainBegin:         d.Get("maintain_begin").(string),
		MaintainEnd:           d.Get("maintain_end").(string),
		RetentionPolicy:       d.Get("retention_policy").(string),
		ConnectorEnalbe:       d.Get("dumping").(bool),
		EnableAutoTopic:       d.Get("enable_auto_topic").(bool),
		StorageSpecCode:       d.Get("storage_spec_code").(string),
		EnterpriseProjectID:   cfg.GetEnterpriseProjectID(d),
		SslEnable:             d.Get("ssl_enable").(bool),
		KafkaSecurityProtocol: d.Get("security_protocol").(string),
		SaslEnabledMechanisms: utils.ExpandToStringList(d.Get("enabled_mechanisms").(*schema.Set).List()),
		Ipv6Enable:            d.Get("ipv6_enable").(bool),
		VpcClientPlain:        d.Get("vpc_client_plain").(bool),
		PortProtocol:          buildKafkaPortProtocol(d.Get("port_protocol").([]interface{})),
		TenantIps:             utils.ExpandToStringList(d.Get("new_tenant_ips").([]interface{})),
	}

	if chargingMode, ok := d.GetOk("charging_mode"); ok && chargingMode == "prePaid" {
		var autoRenew bool
		if d.Get("auto_renew").(string) == "true" {
			autoRenew = true
		}
		isAutoPay := true
		createOpts.BssParam = instances.BssParam{
			ChargingMode: d.Get("charging_mode").(string),
			PeriodType:   d.Get("period_unit").(string),
			PeriodNum:    d.Get("period").(int),
			IsAutoRenew:  &autoRenew,
			IsAutoPay:    &isAutoPay,
		}
	}

	if pubIpIDs, ok := d.GetOk("public_ip_ids"); ok {
		publicIpIDs, err := validateAndBuildPublicIpIDParam(pubIpIDs.([]interface{}), bandwidth)
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.EnablePublicIP = true
		createOpts.PublicIpID = publicIpIDs
	}

	// set tags
	if tagsRaw := d.Get("tags").(map[string]interface{}); len(tagsRaw) > 0 {
		createOpts.Tags = utils.ExpandResourceTags(tagsRaw)
	}
	log.Printf("[DEBUG] Create DMS Kafka instance options: %#v", createOpts)

	// Add password here, so it wouldn't go in the above log entry
	password := d.Get("password").(string)
	if password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			password, err = decryptPasswordWithKmsID(ctx, d, meta)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	createOpts.Password = password
	createOpts.KafkaManagerPassword = d.Get("manager_password").(string)

	kafkaInstance, err := instances.CreateWithEngine(client, createOpts, engineKafka).Extract()
	if err != nil {
		return diag.Errorf("error creating DMS kafka instance: %s", err)
	}
	instanceID := kafkaInstance.InstanceID

	var delayTime time.Duration = 300
	if chargingMode, ok := d.GetOk("charging_mode"); ok && chargingMode == "prePaid" {
		err = waitForOrderComplete(ctx, d, cfg, client, instanceID)
		if err != nil {
			return diag.FromErr(err)
		}
		delayTime = 5
	}

	log.Printf("[INFO] Creating Kafka instance, ID: %s", instanceID)

	// Store the instance ID now
	d.SetId(instanceID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"RUNNING"},
		Refresh:      KafkaInstanceStateRefreshFunc(client, instanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        delayTime * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for Kafka instance (%s) to be ready: %s", instanceID, err)
	}

	// resize storage capacity of the instance
	if ok && storageSpace.(int) != int(defaultStorageSpace) {
		err = resizeKafkaInstanceStorage(ctx, d, client)
		if err != nil {
			dErrs := diag.Errorf("Kafka instance is created, but fails resize the storage capacity, "+
				"expected %v GB, but got %v GB, error: %s ", storageSpace.(int), defaultStorageSpace, err)
			dErrs[0].Severity = diag.Warning
			return dErrs
		}
	}

	return nil
}

func decryptPasswordWithKmsID(_ context.Context, d *schema.ResourceData, meta interface{}) (string, error) {
	var (
		cfg                = meta.(*config.Config)
		region             = cfg.GetRegion(d)
		dataDecryptHttpUrl = "v1.0/{project_id}/kms/decrypt-data"
		dataDecryptProduct = "kms"
	)

	client, err := cfg.NewServiceClient(dataDecryptProduct, region)
	if err != nil {
		return "", fmt.Errorf("error creating KMS client: %s", err)
	}

	dataDecryptPath := client.Endpoint + dataDecryptHttpUrl
	dataDecryptPath = strings.ReplaceAll(dataDecryptPath, "{project_id}", client.ProjectID)

	dataDecryptOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	bodyParams := map[string]interface{}{
		"cipher_text": utils.ValueIgnoreEmpty(d.Get("kms_encrypted_password")),
	}

	dataDecryptOpt.JSONBody = utils.RemoveNil(bodyParams)
	dataDecryptResp, err := client.Request("POST", dataDecryptPath, &dataDecryptOpt)
	if err != nil {
		return "", fmt.Errorf("error running kms decrypt operation: %s", err)
	}

	dataDecryptRespBody, err := utils.FlattenResponse(dataDecryptResp)
	if err != nil {
		return "", fmt.Errorf("err flatting response: %s", err)
	}

	plainText := utils.PathSearch("plain_text", dataDecryptRespBody, "").(string)
	if plainText == "" {
		return "", errors.New("unable to find the plain text from the API response")
	}

	return plainText, nil
}

func waitForOrderComplete(ctx context.Context, d *schema.ResourceData, conf *config.Config,
	client *golangsdk.ServiceClient, instanceID string) error {
	region := conf.GetRegion(d)
	orderId, err := getInstanceOrderId(ctx, d, client, instanceID)
	if err != nil {
		return err
	}
	if orderId == "" {
		log.Printf("[WARN] error get order id by instance ID: %s", instanceID)
		return nil
	}

	bssClient, err := conf.BssV2Client(region)
	if err != nil {
		return fmt.Errorf("error creating BSS v2 client: %s", err)
	}
	// wait for order success
	err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	_, err = common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error waiting for Kafka order resource %s complete: %s", orderId, err)
	}
	return nil
}

func getInstanceOrderId(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) (string, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"EMPTY"},
		Target:       []string{"CREATING"},
		Refresh:      kafkaInstanceCreatingFunc(client, instanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        500 * time.Millisecond,
		PollInterval: 500 * time.Millisecond,
	}
	orderId, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "", fmt.Errorf("error waiting for Kafka instance (%s) to creating: %s", instanceID, err)
	}
	return orderId.(string), nil
}

func kafkaInstanceCreatingFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res := instances.Get(client, instanceID)
		if res.Err != nil {
			actualCode := utils.PathSearch("Actual", res.Err, 0).(int)
			if actualCode == 0 {
				return nil, "", fmt.Errorf("unable to find status code from the API response")
			}
			if actualCode == 404 {
				return res, "EMPTY", nil
			}
		}
		instance, err := res.Extract()
		if err != nil {
			return nil, "", err
		}
		return instance.OrderID, "CREATING", nil
	}
}

func FlattenCrossVpcInfo(str string) (result []map[string]interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Recover panic when flattening Cross-VPC structure: %#v \nCrossVpcInfo: %s", r, str)
			err = fmt.Errorf("faield to flattening Cross-VPC structure: %#v", r)
		}
	}()

	return unmarshalFlattenCrossVpcInfo(str)
}

func unmarshalFlattenCrossVpcInfo(crossVpcInfoStr string) ([]map[string]interface{}, error) {
	if crossVpcInfoStr == "" {
		return nil, nil
	}

	crossVpcInfos := make(map[string]interface{})
	err := json.Unmarshal([]byte(crossVpcInfoStr), &crossVpcInfos)
	if err != nil {
		return nil, fmt.Errorf("failed to Unmarshal CrossVpcInfo, crossVpcInfo: %s, error: %s", crossVpcInfoStr, err)
	}

	ipArr := make([]string, 0, len(crossVpcInfos))
	for ip := range crossVpcInfos {
		ipArr = append(ipArr, ip)
	}
	sort.Strings(ipArr) // Sort by listeners IP.

	result := make([]map[string]interface{}, len(crossVpcInfos))
	for i, ip := range ipArr {
		crossVpcInfo := crossVpcInfos[ip].(map[string]interface{})
		result[i] = map[string]interface{}{
			"listener_ip":   ip,
			"lisenter_ip":   ip,
			"advertised_ip": crossVpcInfo["advertised_ip"],
			"port":          crossVpcInfo["port"],
			"port_id":       crossVpcInfo["port_id"],
		}
	}
	return result, nil
}

func setKafkaFlavorId(d *schema.ResourceData, flavorId string) error {
	re := regexp.MustCompile(`^\d(\d|-)*\d$`)
	if re.MatchString(flavorId) {
		return d.Set("product_id", flavorId)
	}
	return d.Set("flavor_id", flavorId)
}

func resourceDmsKafkaInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.DmsV2Client(region)
	if err != nil {
		return diag.Errorf("error initializing DMS Kafka(v2) client: %s", err)
	}
	eipClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating networking client: %s", err)
	}

	v, err := instances.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DMS Kafka instance")
	}
	log.Printf("[DEBUG] Get Kafka instance: %+v", v)

	crossVpcAccess, err := FlattenCrossVpcInfo(v.CrossVpcInfo)
	if err != nil {
		return diag.Errorf("error parsing the cross-VPC information: %v", err)
	}

	partitionNum, _ := strconv.ParseInt(v.PartitionNum, 10, 64)
	// Convert the AZ ids to AZ codes.
	availableZoneIDs := v.AvailableZones
	availableZoneCodes, err := GetAvailableZoneCodeByID(cfg, region, availableZoneIDs)
	mErr := multierror.Append(nil, err)

	var chargingMode = "postPaid"
	if v.ChargingMode == 0 {
		chargingMode = "prePaid"
	}

	var publicIpIds []string
	var publicIpAddrs []string
	if v.PublicConnectionAddress != "" {
		addressList := getPublicIPAddresses(v.PublicConnectionAddress)
		publicIps, err := common.GetEipsbyAddresses(eipClient, addressList, "all_granted_eps")
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}
		publicIpIds = make([]string, len(publicIps))
		publicIpAddrs = make([]string, len(publicIps))
		for i, ip := range publicIps {
			publicIpIds[i] = ip.ID
			publicIpAddrs[i] = ip.PublicAddress
		}
	}
	createdAt, _ := strconv.Atoi(v.CreatedAt)

	mErr = multierror.Append(mErr,
		d.Set("region", cfg.GetRegion(d)),
		setKafkaFlavorId(d, v.ProductID), // Set flavor_id or product_id.
		d.Set("name", v.Name),
		d.Set("description", v.Description),
		d.Set("engine_version", v.EngineVersion),
		d.Set("bandwidth", v.Specification),
		// storage_space indicates total_storage_space while creating
		// set value of total_storage_space to storage_space to keep consistent
		d.Set("storage_space", v.TotalStorageSpace),
		d.Set("vpc_id", v.VPCID),
		d.Set("security_group_id", v.SecurityGroupID),
		d.Set("network_id", v.SubnetID),
		d.Set("ipv6_enable", v.Ipv6Enable),
		d.Set("available_zones", availableZoneIDs),
		d.Set("availability_zones", availableZoneCodes),
		d.Set("broker_num", v.BrokerNum),
		d.Set("maintain_begin", v.MaintainBegin),
		d.Set("maintain_end", v.MaintainEnd),
		d.Set("ssl_enable", v.SslEnable),
		d.Set("retention_policy", v.RetentionPolicy),
		d.Set("dumping", v.ConnectorEnalbe),
		d.Set("enable_auto_topic", v.EnableAutoTopic),
		d.Set("storage_spec_code", v.StorageSpecCode),
		d.Set("enterprise_project_id", v.EnterpriseProjectID),
		d.Set("manegement_connect_address", v.ManagementConnectAddress),
		d.Set("management_connect_address", v.ManagementConnectAddress),
		d.Set("access_user", v.AccessUser),
		d.Set("cross_vpc_accesses", crossVpcAccess),
		d.Set("charging_mode", chargingMode),
		d.Set("public_ip_ids", publicIpIds),
		d.Set("vpc_client_plain", v.VpcClientPlain),
		d.Set("port_protocols", flattenKafkaSecurityConfig(v.PortProtocols)),
		d.Set("port_protocol", flattenKafkaSecurityConfig(v.PortProtocols)),
		// Attributes.
		d.Set("engine", v.Engine),
		d.Set("partition_num", partitionNum),
		d.Set("enable_public_ip", v.EnablePublicIP),
		d.Set("public_ip_address", publicIpAddrs),
		d.Set("used_storage_space", v.UsedStorageSpace),
		d.Set("connect_address", v.ConnectAddress),
		d.Set("port", v.Port),
		d.Set("status", v.Status),
		d.Set("resource_spec_code", v.ResourceSpecCode),
		d.Set("user_id", v.UserID),
		d.Set("user_name", v.UserName),
		d.Set("extend_times", v.ExtendTimes),
		d.Set("ipv6_connect_addresses", v.Ipv6ConnectAddresses),
		d.Set("connector_id", v.ConnectorID),
		d.Set("connector_node_num", v.ConnectorNodeNum),
		d.Set("storage_resource_id", v.StorageResourceID),
		d.Set("storage_type", v.StorageType),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(createdAt)/1000, false)),
		d.Set("cert_replaced", v.CertReplaced),
		d.Set("is_logical_volume", v.IsLogicalVolume),
		d.Set("message_query_inst_enable", v.MessageQueryInstEnable),
		d.Set("node_num", v.NodeNum),
		d.Set("pod_connect_address", v.PodConnectAddress),
		d.Set("public_bandwidth", v.PublicBandWidth),
		d.Set("ssl_two_way_enable", v.SslTwoWayEnable),
		d.Set("type", v.Type),
	)

	// set tags
	if resourceTags, err := tags.Get(client, engineKafka, d.Id()).Extract(); err == nil {
		tagMap := utils.TagsToMap(resourceTags.Tags)
		if err = d.Set("tags", tagMap); err != nil {
			mErr = multierror.Append(mErr,
				fmt.Errorf("error saving tags to state for DMS kafka instance (%s): %s", d.Id(), err))
		}
	} else {
		log.Printf("[WARN] error fetching tags of DMS kafka instance (%s): %s", d.Id(), err)
	}

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("failed to set attributes for DMS kafka instance: %s", mErr)
	}

	return setKafkaInstanceParameters(ctx, d, client)
}

func getPublicIPAddresses(rawParam string) []string {
	allAddressPortList := strings.Split(rawParam, ",")
	rst := make([]string, 0, len(allAddressPortList))
	for _, addressPort := range allAddressPortList {
		address := strings.Split(addressPort, ":")[0]
		rst = append(rst, address)
	}
	return rst
}

func setKafkaInstanceParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) diag.Diagnostics {
	// set parameters
	configs, err := instances.GetConfigurations(client, d.Id()).Extract()
	if err != nil {
		log.Printf("[WARN] error fetching parameters of the instance (%s): %s", d.Id(), err)
		return nil
	}

	var params []map[string]interface{}
	for _, parameter := range d.Get("parameters").(*schema.Set).List() {
		name := parameter.(map[string]interface{})["name"]
		for _, kafkaParam := range configs.KafkaConfigs {
			if kafkaParam.Name == name {
				p := map[string]interface{}{
					"name":  kafkaParam.Name,
					"value": kafkaParam.Value,
				}
				params = append(params, p)
				break
			}
		}
	}

	if len(params) > 0 {
		if err = d.Set("parameters", params); err != nil {
			log.Printf("[WARN] error saving parameters to the Kafka instance (%s): %s", d.Id(), err)
		}
		if ctx.Value(ctxType("staticParametersChanged")) == "true" {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Parameters Changed",
					Detail:   "Static parameters changed, the instance needs reboot to make parameters take effect.",
				},
			}
		}
	}
	return nil
}

func flattenKafkaSecurityConfig(kafkaSecurityConfig instances.PortProtocols) interface{} {
	return []map[string]interface{}{
		{
			"private_plain_enable":               kafkaSecurityConfig.PrivatePlainEnable,
			"private_plain_address":              kafkaSecurityConfig.PrivatePlainAddress,
			"private_plain_domain_name":          kafkaSecurityConfig.PrivatePlainDomainName,
			"private_sasl_ssl_enable":            kafkaSecurityConfig.PrivateSaslSslEnable,
			"private_sasl_ssl_address":           kafkaSecurityConfig.PrivateSaslSslAddress,
			"private_sasl_ssl_domain_name":       kafkaSecurityConfig.PrivateSaslSslDomainName,
			"private_sasl_plaintext_enable":      kafkaSecurityConfig.PrivateSaslPlaintextEnable,
			"private_sasl_plaintext_address":     kafkaSecurityConfig.PrivateSaslPlaintextAddress,
			"private_sasl_plaintext_domain_name": kafkaSecurityConfig.PrivateSaslPlaintextDomainName,
			"public_plain_enable":                kafkaSecurityConfig.PublicPlainEnable,
			"public_plain_address":               kafkaSecurityConfig.PublicPlainAddress,
			"public_plain_domain_name":           kafkaSecurityConfig.PublicPlainDomainName,
			"public_sasl_ssl_enable":             kafkaSecurityConfig.PublicSaslSslEnable,
			"public_sasl_ssl_address":            kafkaSecurityConfig.PublicSaslSslAddress,
			"public_sasl_ssl_domain_name":        kafkaSecurityConfig.PublicSaslSslDomainName,
			"public_sasl_plaintext_enable":       kafkaSecurityConfig.PublicSaslPlaintextEnable,
			"public_sasl_plaintext_address":      kafkaSecurityConfig.PublicSaslPlaintextAddress,
			"public_sasl_plaintext_domain_name":  kafkaSecurityConfig.PublicSaslPlaintextDomainName,
		},
	}
}

func resourceDmsKafkaInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error initializing DMS Kafka(v2) client: %s", err)
	}

	var mErr *multierror.Error
	if d.HasChanges("name", "description", "maintain_begin", "maintain_end",
		"security_group_id", "retention_policy", "enterprise_project_id") {
		description := d.Get("description").(string)
		updateOpts := instances.UpdateOpts{
			Description:         &description,
			MaintainBegin:       d.Get("maintain_begin").(string),
			MaintainEnd:         d.Get("maintain_end").(string),
			SecurityGroupID:     d.Get("security_group_id").(string),
			RetentionPolicy:     d.Get("retention_policy").(string),
			EnterpriseProjectID: d.Get("enterprise_project_id").(string),
		}

		if d.HasChange("name") {
			updateOpts.Name = d.Get("name").(string)
		}

		retryFunc := func() (interface{}, bool, error) {
			err = instances.Update(client, d.Id(), updateOpts).Err
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     KafkaInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"RUNNING"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			mErr = multierror.Append(mErr, fmt.Errorf("error updating Kafka Instance: %s", err))
		}
	}

	if d.HasChanges("product_id", "flavor_id", "storage_space", "broker_num") {
		err = resizeKafkaInstance(ctx, d, meta)
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}

	if d.HasChange("tags") {
		// update tags
		if err = utils.UpdateResourceTags(client, d, engineKafka, d.Id()); err != nil {
			mErr = multierror.Append(mErr, fmt.Errorf("error updating tags of Kafka instance: %s, err: %s",
				d.Id(), err))
		}
	}

	if d.HasChange("cross_vpc_accesses") {
		if err = UpdateCrossVpcAccess(ctx, client, d); err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
			return diag.Errorf("error updating the auto-renew of the Kafka instance (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("enable_auto_topic") {
		enableAutoTopic := d.Get("enable_auto_topic").(bool)
		autoTopicOpts := instances.AutoTopicOpts{
			EnableAutoTopic: &enableAutoTopic,
		}
		retryFunc := func() (interface{}, bool, error) {
			err = instances.UpdateAutoTopic(client, d.Id(), autoTopicOpts).Err
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     KafkaInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"RUNNING"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			mErr = multierror.Append(mErr, fmt.Errorf("error enabling or disabling automatic topic: %s", err))
		}

		// The enabling or disabling automatic topic is done if the status of its related task is SUCCESS
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"CREATED"},
			Target:       []string{"SUCCESS"},
			Refresh:      FilterTaskRefreshFunc(client, d.Id(), "kafkaConfigModify"),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        1 * time.Second,
			PollInterval: 5 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			mErr = multierror.Append(mErr,
				fmt.Errorf("error waiting for the automatic topic task of the instance (%s) to be done: %s", d.Id(), err))
		}
	}

	// This logic must be done before resetting the password, because resetting the password is not possible when SSL has never been enabled.
	if d.HasChange("port_protocol") {
		if err = updateInstancePortProtocol(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("password", "kms_encrypted_password") {
		password := d.Get("password").(string)
		if password == "" {
			if v := d.Get("kms_encrypted_password").(string); v != "" {
				password, err = decryptPasswordWithKmsID(ctx, d, meta)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
		resetPasswordOpts := instances.ResetPasswordOpts{
			NewPassword: password,
		}
		retryFunc := func() (interface{}, bool, error) {
			err = instances.ResetPassword(client, d.Id(), resetPasswordOpts).Err
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     KafkaInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"RUNNING"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			e := fmt.Errorf("error resetting password: %s", err)
			mErr = multierror.Append(mErr, e)
		}
	}

	if ctx, err = updateKafkaParameters(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error while updating DMS Kafka instances, %s", mErr)
	}

	return resourceDmsKafkaInstanceRead(ctx, d, meta)
}

func resizeKafkaInstance(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)
	client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error initializing DMS(v2) client: %s", err)
	}

	if d.HasChanges("product_id") {
		product, err := getKafkaProductDetails(cfg, d)
		if err != nil {
			return fmt.Errorf("failed to resize Kafka instance, query product details error: %s", err)
		}
		storageSpace := d.Get("storage_space").(int)
		resizeOpts := instances.ResizeInstanceOpts{
			NewSpecCode:     &product.SpecCode,
			NewStorageSpace: &storageSpace,
		}
		log.Printf("[DEBUG] Resize Kafka instance storage space options: %s", utils.MarshalValue(resizeOpts))

		if err := doKafkaInstanceResize(ctx, d, client, resizeOpts); err != nil {
			return err
		}
	}

	if d.HasChanges("flavor_id") {
		flavorID := d.Get("flavor_id").(string)
		operType := "vertical"
		resizeOpts := instances.ResizeInstanceOpts{
			OperType:     &operType,
			NewProductID: &flavorID,
		}
		log.Printf("[DEBUG] Resize Kafka instance flavor ID options: %s", utils.MarshalValue(resizeOpts))

		if err := doKafkaInstanceResize(ctx, d, client, resizeOpts); err != nil {
			return err
		}
	}

	if d.HasChanges("broker_num") {
		operType := "horizontal"
		brokerNum := d.Get("broker_num").(int)

		resizeOpts := instances.ResizeInstanceOpts{
			OperType:     &operType,
			NewBrokerNum: &brokerNum,
		}
		oldNum, newNum := d.GetChange("broker_num")
		if d.HasChange("public_ip_ids") {
			// precheck
			oldRaw, newRaw := d.GetChange("public_ip_ids")
			add := newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))
			sub := oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))
			allow := (sub.Len() == 0) && (add.Len() == newNum.(int)-oldNum.(int))
			if !allow {
				return fmt.Errorf("error resizing instance: the old EIP ID should not be changed, and the adding nums of " +
					"EIP ID should be same as the adding broker nums")
			}

			publicIpIds := strings.Join(utils.ExpandToStringList(add.List()), ",")
			resizeOpts.PublicIpID = &publicIpIds
		}
		if v, ok := d.GetOk("new_tenant_ips"); ok {
			// precheck
			if len(v.([]interface{})) > newNum.(int)-oldNum.(int) {
				return fmt.Errorf("error resizing instance: the nums of new tenant IP must be less than the adding broker nums")
			}
			resizeOpts.TenantIps = utils.ExpandToStringList(v.([]interface{}))
		}

		log.Printf("[DEBUG] Resize Kafka instance broker number options: %s", utils.MarshalValue(resizeOpts))

		if err := doKafkaInstanceResize(ctx, d, client, resizeOpts); err != nil {
			return err
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"PENDING"},
			Target:       []string{"BOUND"},
			Refresh:      kafkaInstanceBrokerNumberRefreshFunc(client, d.Id(), brokerNum),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        10 * time.Second,
			PollInterval: 10 * time.Second,
		}
		if _, err = stateConf.WaitForStateContext(ctx); err != nil {
			return err
		}
	}

	if d.HasChanges("storage_space") {
		if err = resizeKafkaInstanceStorage(ctx, d, client); err != nil {
			return err
		}
	}

	return nil
}

func resizeKafkaInstanceStorage(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	newStorageSpace := d.Get("storage_space").(int)
	operType := "storage"
	resizeOpts := instances.ResizeInstanceOpts{
		OperType:        &operType,
		NewStorageSpace: &newStorageSpace,
	}
	log.Printf("[DEBUG] Resize Kafka instance storage space options: %s", utils.MarshalValue(resizeOpts))

	return doKafkaInstanceResize(ctx, d, client, resizeOpts)
}

func doKafkaInstanceResize(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, opts instances.ResizeInstanceOpts) error {
	retryFunc := func() (interface{}, bool, error) {
		_, err := instances.Resize(client, d.Id(), opts)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     KafkaInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("resize Kafka instance failed: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING", "EXTENDING"},
		Target:       []string{"RUNNING"},
		Refresh:      kafkaResizeStateRefresh(client, d, opts.OperType),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        180 * time.Second,
		PollInterval: 15 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for instance (%s) to resize: %v", d.Id(), err)
	}
	return nil
}

func kafkaResizeStateRefresh(client *golangsdk.ServiceClient, d *schema.ResourceData, operType *string) resource.StateRefreshFunc {
	flavorID := d.Get("flavor_id").(string)
	if flavorID == "" {
		flavorID = d.Get("product_id").(string)
	}
	storageSpace := d.Get("storage_space").(int)
	brokerNum := d.Get("broker_num").(int)

	return func() (interface{}, string, error) {
		v, err := instances.Get(client, d.Id()).Extract()
		if err != nil {
			return nil, "failed", err
		}

		if ((operType == nil || *operType == "vertical") && v.ProductID != flavorID) || // change flavor
			(operType != nil && *operType == "storage" && v.TotalStorageSpace != storageSpace) || // expansion
			(operType != nil && *operType == "horizontal" && v.BrokerNum != brokerNum) { // expand broker number
			return v, "PENDING", nil
		}

		return v, v.Status, nil
	}
}

func buildSwitchInstancePortProtocolBodyParams(d *schema.ResourceData, protocol string, enable bool) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"protocol":  protocol,
		"enable":    enable,
		"user_name": utils.ValueIgnoreEmpty(d.Get("access_user").(string)),
		"pass_word": utils.ValueIgnoreEmpty(d.Get("password").(string)),
	}

	// sasl_enabled_mechanisms only support encrypted protocol.
	encryptedProtocolFields := []string{
		"private_sasl_ssl_enable",
		"private_sasl_plaintext_enable",
		"public_sasl_ssl_enable",
		"public_sasl_plaintext_enable",
	}
	if utils.StrSliceContains(encryptedProtocolFields, protocol) {
		bodyParams["sasl_enabled_mechanisms"] = utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("enabled_mechanisms").(*schema.Set).List()))
	}

	return bodyParams
}

func switchInstancePortProtocol(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	protocol string, enable bool) error {
	var (
		httpUrl       = "v2/{project_id}/kafka/instances/{instance_id}/plain-ssl-switch"
		instanceId    = d.Id()
		updateTimeout = d.Timeout(schema.TimeoutUpdate)
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildSwitchInstancePortProtocolBodyParams(d, protocol, enable)),
	}

	retryFunc := func() (interface{}, bool, error) {
		resp, err := client.Request("POST", updatePath, &requestOpt)
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	resp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     KafkaInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      updateTimeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating instance port protocol (%s), %s", protocol, err)
	}

	respBody, err := utils.FlattenResponse(resp.(*http.Response))
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("unable to find job ID of the instance port protocol (%s) update", protocol)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATED"},
		Target:       []string{"SUCCESS"},
		Refresh:      kafkaInstanceTaskStatusRefreshFunc(client, instanceId, jobId),
		Timeout:      updateTimeout,
		Delay:        15 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)

	return err
}

func updateInstancePortProtocol(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	oldRaw, newRaw := d.GetChange("port_protocol")
	newPortProtocol := newRaw.([]interface{})[0].(map[string]interface{})
	oldPortProtocol := oldRaw.([]interface{})[0].(map[string]interface{})
	parsedPortProtocol := make([]map[string]interface{}, 0, len(newPortProtocol))
	for k, v := range newPortProtocol {
		newValue, ok := v.(bool)
		if !ok {
			continue
		}

		// Compare the new and old port protocol, if the value is different, add it to the parsedPortProtocol.
		if newValue != oldPortProtocol[k].(bool) {
			item := map[string]interface{}{
				"key":   k,
				"value": newValue,
			}

			// At least one of the port protocol should be enabled, so enable them first and then disable.
			if newValue {
				parsedPortProtocol = append([]map[string]interface{}{item}, parsedPortProtocol...)
				continue
			}
			parsedPortProtocol = append(parsedPortProtocol, item)
		}
	}

	for _, item := range parsedPortProtocol {
		err := switchInstancePortProtocol(ctx, client, d, utils.PathSearch("key", item, "").(string),
			utils.PathSearch("value", item, "").(bool))
		if err != nil {
			return fmt.Errorf("error updating port protocol instance (%s), %s", d.Id(), err)
		}
	}

	return nil
}

func resourceDmsKafkaInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error initializing DMS Kafka(v2) client: %s", err)
	}

	if d.Get("charging_mode") == "prePaid" {
		retryFunc := func() (interface{}, bool, error) {
			err = common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()})
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     KafkaInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"RUNNING"},
			Timeout:      d.Timeout(schema.TimeoutDelete),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.Errorf("error unsubscribe Kafka instance: %s", err)
		}
	} else {
		retryFunc := func() (interface{}, bool, error) {
			err = instances.Delete(client, d.Id()).ExtractErr()
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     KafkaInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"RUNNING"},
			Timeout:      d.Timeout(schema.TimeoutDelete),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return common.CheckDeletedDiag(d, err, "failed to delete Kafka instance")
		}
	}

	// Wait for the instance to delete before moving on.
	log.Printf("[DEBUG] Waiting for Kafka instance (%s) to be deleted", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"DELETING", "RUNNING", "ERROR"}, // Status may change to ERROR on deletion.
		Target:       []string{"DELETED"},
		Refresh:      KafkaInstanceStateRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        120 * time.Second,
		PollInterval: 15 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for DMS Kafka instance (%s) to be deleted: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] DMS Kafka instance %s has been deleted", d.Id())
	d.SetId("")
	return nil
}

func kafkaInstanceBrokerNumberRefreshFunc(client *golangsdk.ServiceClient, instanceID string, brokerNum int) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			return nil, "QUERY ERROR", err
		}

		if brokerNum == resp.BrokerNum && resp.CrossVpcInfo != "" {
			crossVpcInfoMap, err := FlattenCrossVpcInfo(resp.CrossVpcInfo)
			if err != nil {
				return resp, "ParseError", err
			}

			if len(crossVpcInfoMap) == brokerNum {
				return resp, "BOUND", nil
			}
		}
		return resp, "PENDING", nil
	}
}

func kafkaInstanceCrossVpcInfoRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			return nil, "QUERY ERROR", err
		}
		if resp.CrossVpcInfo != "" {
			return resp, "BOUND", nil
		}
		return resp, "PENDING", nil
	}
}

func KafkaInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "DELETED", nil
			}
			return nil, "QUERY ERROR", err
		}

		return v, v.Status, nil
	}
}

func GetAvailableZoneIDByCode(cfg *config.Config, region string, azCodes []interface{}) ([]string, error) {
	if len(azCodes) == 0 {
		return nil, fmt.Errorf(`arguments "azCodes" is required`)
	}

	availableZones, err := getAvailableZones(cfg, region)
	if err != nil {
		return nil, err
	}

	codeIDMapping := make(map[string]string)
	for _, v := range availableZones {
		codeIDMapping[v.Code] = v.ID
	}

	azIDs := make([]string, 0, len(azCodes))
	for _, code := range azCodes {
		if id, ok := codeIDMapping[code.(string)]; ok {
			azIDs = append(azIDs, id)
		}
	}
	log.Printf("[DEBUG] DMS converts the AZ codes to AZ IDs: \n%#v => \n%#v", azCodes, azIDs)
	return azIDs, nil
}

func GetAvailableZoneCodeByID(cfg *config.Config, region string, azIDs []string) ([]string, error) {
	if len(azIDs) == 0 {
		return nil, fmt.Errorf(`arguments "azIDs" is required`)
	}

	availableZones, err := getAvailableZones(cfg, region)
	if err != nil {
		return nil, err
	}

	idCodeMapping := make(map[string]string)
	for _, v := range availableZones {
		idCodeMapping[v.ID] = v.Code
	}

	azCodes := make([]string, 0, len(azIDs))
	for _, id := range azIDs {
		if code, ok := idCodeMapping[id]; ok {
			azCodes = append(azCodes, code)
		}
	}
	log.Printf("[DEBUG] DMS converts the AZ IDs to AZ codes: \n%#v => \n%#v", azIDs, azCodes)
	return azCodes, nil
}

func getAvailableZones(cfg *config.Config, region string) ([]availablezones.AvailableZone, error) {
	client, err := cfg.DmsV2Client(region)
	if err != nil {
		return nil, fmt.Errorf("error initializing DMS(v2) client: %s", err)
	}

	r, err := availablezones.Get(client)
	if err != nil {
		return nil, fmt.Errorf("error querying available Zones: %s", err)
	}

	return r.AvailableZones, nil
}

func handleMultiOperationsError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode := utils.PathSearch("error_code", apiError, "").(string)
		if errorCode == "" {
			return false, fmt.Errorf("unable to find error code from the API response")
		}

		// CBC.99003651: unsubscribe fail, another operation is being performed
		if errorCode == "DMS.00400026" || errorCode == "CBC.99003651" {
			return true, err
		}
	}
	return false, err
}

func FilterTaskRefreshFunc(client *golangsdk.ServiceClient, instanceID string, taskName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// getAutoTopicTask: query automatic topic task
		getTasksHttpUrl := "v2/{project_id}/instances/{instance_id}/tasks"
		getTasksPath := client.Endpoint + getTasksHttpUrl
		getTasksPath = strings.ReplaceAll(getTasksPath, "{project_id}",
			client.ProjectID)
		getTasksPath = strings.ReplaceAll(getTasksPath, "{instance_id}", instanceID)
		getTasksPathOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getTasksPathResp, err := client.Request("GET", getTasksPath, &getTasksPathOpt)
		if err != nil {
			return nil, "QUERY ERROR", err
		}
		getTasksRespBody, err := utils.FlattenResponse(getTasksPathResp)
		if err != nil {
			return nil, "PARSE ERROR", err
		}

		task := utils.PathSearch(fmt.Sprintf("tasks|[?name=='%s']|[0]", taskName), getTasksRespBody, nil)
		if task == nil {
			return nil, "NIL ERROR", fmt.Errorf("failed to find the task of the name(%s)", taskName)
		}

		status := utils.PathSearch("status", task, nil)
		return task, fmt.Sprint(status), nil
	}
}

func buildKafkaInstanceParameters(params *schema.Set) instances.KafkaConfigs {
	paramList := make([]instances.ConfigParam, 0, params.Len())
	for _, v := range params.List() {
		paramList = append(paramList, instances.ConfigParam{
			Name:  v.(map[string]interface{})["name"].(string),
			Value: v.(map[string]interface{})["value"].(string),
		})
	}
	configOpts := instances.KafkaConfigs{
		KafkaConfigs: paramList,
	}
	return configOpts
}

func initializeParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	parametersRaw := d.Get("parameters").(*schema.Set)
	configOpts := buildKafkaInstanceParameters(parametersRaw)
	restartRequired, err := modifyParameters(ctx, client, d.Timeout(schema.TimeoutCreate), d.Id(), &configOpts)
	if err != nil {
		return err
	}

	if *restartRequired {
		return restartKafkaInstance(ctx, d.Timeout(schema.TimeoutCreate), client, d.Id())
	}
	return nil
}

func modifyParameters(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	instanceID string, configOpts *instances.KafkaConfigs) (*bool, error) {
	// modify configs
	retryFunc := func() (interface{}, bool, error) {
		resp, err := instances.ModifyConfiguration(client, instanceID, *configOpts).Extract()
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     KafkaInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("error modifying parameters for the Kafka instance (%s): %s", instanceID, err)
	}

	modifyConfigurationResp := r.(*instances.ModifyConfigurationResp)

	// wait for task complete
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATED"},
		Target:       []string{"SUCCESS"},
		Refresh:      kafkaInstanceTaskStatusRefreshFunc(client, instanceID, modifyConfigurationResp.JobId),
		Timeout:      timeout,
		Delay:        2 * time.Second,
		PollInterval: 2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return nil, fmt.Errorf("error waiting for the Kafka instance (%s) parameter to be updated: %s ", instanceID, err)
	}

	restartRequired := modifyConfigurationResp.StaticConfig > 0

	return &restartRequired, nil
}

func restartKafkaInstance(ctx context.Context, timeout time.Duration, client *golangsdk.ServiceClient,
	instanceID string) error {
	restartInstanceOpts := instances.RestartInstanceOpts{
		Action:    "restart",
		Instances: []string{instanceID},
	}

	retryFunc := func() (interface{}, bool, error) {
		_, err := instances.RebootInstance(client, restartInstanceOpts).Extract()
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     KafkaInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error rebooting the Kafka instance (%s): %s", instanceID, err)
	}

	// wait for the instance state to be 'RUNNING'.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RESTARTING"},
		Target:       []string{"RUNNING"},
		Refresh:      KafkaInstanceStateRefreshFunc(client, instanceID),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for the Kafka instance (%s) become RUNNING status: %s", instanceID, err)
	}
	return nil
}

func updateKafkaParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) (context.Context, error) {
	if !d.HasChange("parameters") {
		return ctx, nil
	}

	o, n := d.GetChange("parameters")
	os, ns := o.(*schema.Set), n.(*schema.Set)
	change := ns.Difference(os).List()
	paramList := make([]instances.ConfigParam, 0, len(change))
	if len(change) > 0 {
		for _, v := range change {
			configOpts := instances.ConfigParam{
				Name:  v.(map[string]interface{})["name"].(string),
				Value: v.(map[string]interface{})["value"].(string),
			}
			paramList = append(paramList, configOpts)
		}

		configOpts := instances.KafkaConfigs{
			KafkaConfigs: paramList,
		}

		restartRequired, err := modifyParameters(ctx, client, d.Timeout(schema.TimeoutCreate), d.Id(), &configOpts)
		if err != nil {
			return ctx, err
		}
		if *restartRequired {
			// Sending staticParametersChanged to Read to warn users the instance needs a reboot.
			ctx = context.WithValue(ctx, ctxType("staticParametersChanged"), "true")
		}
	}

	return ctx, nil
}

func kafkaInstanceTaskStatusRefreshFunc(client *golangsdk.ServiceClient, instanceID, taskID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		taskResp, err := instances.GetTask(client, instanceID, taskID).Extract()
		if err != nil {
			return nil, "QUERY ERROR", err
		}
		if len(taskResp.Tasks) == 0 {
			return nil, "NIL ERROR", fmt.Errorf("failed to find task(%s)", taskID)
		}
		return taskResp.Tasks[0], taskResp.Tasks[0].Status, nil
	}
}
