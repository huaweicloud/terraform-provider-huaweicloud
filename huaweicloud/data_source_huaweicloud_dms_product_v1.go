package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/dms/v1/products"
)

func dataSourceDmsProductV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDmsProductV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
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
			},
			"partition_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_spec_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"io_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"node_num": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func dataSourceDmsProductV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dmsV1Client, err := config.dmsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error get HuaweiCloud dms product client: %s", err)
	}

	instance_engine := d.Get("engine").(string)
	if instance_engine != "rabbitmq" && instance_engine != "kafka" {
		return fmt.Errorf("The instance_engine value should be 'rabbitmq' or 'kafka', not: %s", instance_engine)
	}

	v, err := products.Get(dmsV1Client, instance_engine).Extract()
	if err != nil {
		return err
	}
	Products := v.Hourly

	log.Printf("[DEBUG] Dms get products : %+v", Products)
	instance_type := d.Get("instance_type").(string)
	if instance_type != "single" && instance_type != "cluster" {
		return fmt.Errorf("The instance_type value should be 'single' or 'cluster', not: %s", instance_type)
	}
	var FilteredPd []products.Detail
	var FilteredPdInfo []products.ProductInfo
	for _, pd := range Products {
		version := d.Get("version").(string)
		if version != "" && pd.Version != version {
			continue
		}

		for _, value := range pd.Values {
			if value.Name != instance_type {
				continue
			}
			for _, detail := range value.Details {
				vm_specification := d.Get("vm_specification").(string)
				if vm_specification != "" && detail.VMSpecification != vm_specification {
					continue
				}
				bandwidth := d.Get("bandwidth").(string)
				if bandwidth != "" && detail.Bandwidth != bandwidth {
					continue
				}

				if instance_type == "single" || instance_engine == "kafka" {
					storage := d.Get("storage").(string)
					if storage != "" && detail.Storage != storage {
						continue
					}
					IOs := getIObyIOtype(d, detail.IOs)
					if len(IOs) == 0 {
						continue
					}
					detail.IOs = IOs
				} else {
					for _, pdInfo := range detail.ProductInfos {
						storage := d.Get("storage").(string)
						if storage != "" && pdInfo.Storage != storage {
							continue
						}
						node_num := d.Get("node_num").(string)
						if node_num != "" && pdInfo.NodeNum != node_num {
							continue
						}

						IOs := getIObyIOtype(d, pdInfo.IOs)
						if len(IOs) == 0 {
							continue
						}
						pdInfo.IOs = IOs
						FilteredPdInfo = append(FilteredPdInfo, pdInfo)
					}
					if len(FilteredPdInfo) == 0 {
						continue
					}
					detail.ProductInfos = FilteredPdInfo
				}
				FilteredPd = append(FilteredPd, detail)
			}
		}
	}

	if len(FilteredPd) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your filters and try again.")
	}

	pd := FilteredPd[0]
	d.Set("vm_specification", pd.VMSpecification)
	if instance_type == "single" || instance_engine == "kafka" {
		d.SetId(pd.ProductID)
		d.Set("storage", pd.Storage)
		d.Set("partition_num", pd.PartitionNum)
		d.Set("bandwidth", pd.Bandwidth)
		d.Set("storage_spec_code", pd.IOs[0].StorageSpecCode)
		d.Set("io_type", pd.IOs[0].IOType)
		log.Printf("[DEBUG] Dms product : %+v", pd)
	} else {
		if len(pd.ProductInfos) < 1 {
			return fmt.Errorf("Your query returned no results. Please change your filters and try again.")
		}
		pdInfo := pd.ProductInfos[0]
		d.SetId(pdInfo.ProductID)
		d.Set("storage", pdInfo.Storage)
		d.Set("io_type", pdInfo.IOs[0].IOType)
		d.Set("node_num", pdInfo.NodeNum)
		d.Set("storage_spec_code", pdInfo.IOs[0].StorageSpecCode)
		log.Printf("[DEBUG] Dms product : %+v", pdInfo)
	}

	return nil
}
