package dms

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/dms/v2/products"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

const (
	engineKafka    = "kafka"
	engineRabbitMQ = "rabbitmq"

	instanceTypeSingle  = "single"
	instanceTypeCluster = "cluster"
)

func DataSourceDmsProduct() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsProductRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{engineKafka, engineRabbitMQ},
					false,
				),
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{instanceTypeSingle, instanceTypeCluster},
					false,
				),
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vm_specification": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"partition_num": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage_spec_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage_spec_codes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"node_num": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"io_type": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "io_type has deprecated, please use storage_spec_code",
			},
		},
	}
}

func getIObyIOtype(d *schema.ResourceData, IOs []products.IO) []products.IO {
	io_type := d.Get("io_type").(string)
	storage_spec_code := d.Get("storage_spec_code").(string)

	if io_type != "" || storage_spec_code != "" {
		var getIOs []products.IO
		for _, io := range IOs {
			if io_type == io.IOType || storage_spec_code == io.StorageSpecCode {
				getIOs = append(getIOs, io)
			}
		}
		return getIOs
	}

	return IOs
}

func getProducts(config *config.Config, region, engine string) (*products.GetResponse, error) {
	dmsV2Client, err := config.DmsV2Client(region)
	if err != nil {
		return nil, fmtp.Errorf("Error get HuaweiCloud DMS product client V2: %s", err)
	}
	v, err := products.Get(dmsV2Client, engine)
	return v, err
}

func dataSourceDmsProductRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)

	instanceEngine := d.Get("engine").(string)
	r, err := getProducts(config, config.GetRegion(d), instanceEngine)
	if err != nil {
		return diag.FromErr(err)
	}

	hourlyProducts := r.Hourly
	logp.Printf("[DEBUG] Get a list of DMS products, engine:%s, list: %+v", instanceEngine, hourlyProducts)

	instanceType := d.Get("instance_type").(string)
	vmSpecification := d.Get("vm_specification").(string)
	bandwidth := d.Get("bandwidth").(string)
	storage := d.Get("storage").(string)
	nodeNum := d.Get("node_num").(string)
	version := d.Get("version").(string)

	partitionNum := d.Get("partition_num").(string)
	filterAZs := d.Get("availability_zones").([]interface{})

	filteredProducts := make([]products.Detail, 0)
	isFound := false
	for _, pd := range hourlyProducts {
		if version != "" && pd.Version != version {
			continue
		}

		for _, value := range pd.Values {
			if value.Name != instanceType {
				continue
			}
			for _, detail := range value.Details {
				// The vm_specification has been removed and the evs_flavor_id return is the same as the
				// vm_specification.
				if vmSpecification != "" && detail.EcsFlavorId != vmSpecification {
					continue
				}

				if bandwidth != "" && detail.Bandwidth != bandwidth {
					continue
				}

				if partitionNum != "" && detail.PartitionNum != partitionNum {
					continue
				}

				if instanceType == instanceTypeSingle || instanceEngine == engineKafka {
					if storage != "" && detail.Storage != storage {
						continue
					}
					if !filterAZ(detail.AvailableZones, filterAZs) {
						continue
					}

					IOs := getIObyIOtype(d, detail.IOs)
					if len(IOs) == 0 {
						continue
					}
					detail.IOs = IOs
				} else {
					productInfos := make([]products.ProductInfo, 0)
					for _, productInfo := range detail.ProductInfos {
						if storage != "" && productInfo.Storage != storage {
							continue
						}
						if nodeNum != "" && productInfo.NodeNum != nodeNum {
							continue
						}
						if !filterAZ(productInfo.AvailableZones, filterAZs) {
							continue
						}

						IOs := getIObyIOtype(d, productInfo.IOs)
						if len(IOs) == 0 {
							continue
						}
						productInfo.IOs = IOs
						productInfos = append(productInfos, productInfo)
					}
					if len(productInfos) == 0 {
						continue
					}
					detail.ProductInfos = productInfos
				}
				filteredProducts = append(filteredProducts, detail)
				isFound = true
				break
			}
			if isFound {
				break
			}
		}
		if isFound {
			break
		}
	}

	if len(filteredProducts) < 1 {
		return fmtp.DiagErrorf("Your query returned no results. Please change your filters and try again.")
	}

	pd := filteredProducts[0]

	var mErr *multierror.Error
	if instanceType == instanceTypeSingle || instanceEngine == engineKafka {
		d.SetId(pd.ProductID)

		storageSpecCodes := make([]string, 0, len(pd.IOs))
		for _, v := range pd.IOs {
			storageSpecCodes = append(storageSpecCodes, v.StorageSpecCode)
		}
		mErr = multierror.Append(err,
			d.Set("vm_specification", pd.EcsFlavorId),
			d.Set("storage", pd.Storage),
			d.Set("partition_num", pd.PartitionNum),
			d.Set("bandwidth", pd.Bandwidth),
			d.Set("storage_spec_code", pd.IOs[0].StorageSpecCode),
			d.Set("storage_spec_codes", storageSpecCodes),
			d.Set("io_type", pd.IOs[0].IOType),
			d.Set("availability_zones", pd.AvailableZones),
		)
	} else {
		if len(pd.ProductInfos) < 1 {
			return fmtp.DiagErrorf("Your query returned no results. Please change your filters and try again.")
		}
		pdInfo := pd.ProductInfos[0]
		d.SetId(pdInfo.ProductID)

		storageSpecCodes := make([]string, 0, len(pd.IOs))
		for _, v := range pdInfo.IOs {
			storageSpecCodes = append(storageSpecCodes, v.StorageSpecCode)
		}
		mErr = multierror.Append(err,
			d.Set("vm_specification", pd.EcsFlavorId),
			d.Set("storage", pdInfo.Storage),
			d.Set("io_type", pdInfo.IOs[0].IOType),
			d.Set("node_num", pdInfo.NodeNum),
			d.Set("storage_spec_codes", storageSpecCodes),
			d.Set("storage_spec_code", pdInfo.IOs[0].StorageSpecCode),
			d.Set("availability_zones", pdInfo.AvailableZones),
		)
	}
	logp.Printf("[DEBUG] DMS product detail : %#v", pd)
	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("Error setting DMS product attributes: %s", mErr)
	}

	return nil
}

func filterAZ(azs []string, filterAZs []interface{}) bool {
	if len(azs) == 0 {
		return false
	}
	if len(filterAZs) == 0 {
		return true
	}

	validAZMap := map[string]bool{}
	for _, v := range azs {
		validAZMap[v] = true
	}

	for _, v := range filterAZs {
		if _, ok := validAZMap[v.(string)]; !ok {
			return false
		}
	}
	return true
}
