package lvsnetwork

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIfaceVrrp() *schema.Resource {
	return &schema.Resource{
		Create: resourceIfaceVrrpCreate,
		Read:   resourceIfaceVrrpRead,
		Update: resourceIfaceVrrpUpdate,
		Delete: resourceIfaceVrrpDelete,

		Schema: map[string]*schema.Schema{
			"iface": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_vip": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_vip_only": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ip_master": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					testInput := net.ParseIP(value)
					if testInput.To16() == nil {
						errors = append(errors, fmt.Errorf("[ERROR] %q %v isn't an IPv4 or IPv6", k, value))
					}
					return
				},
			},
			"ip_slave": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					testInput := net.ParseIP(value)
					if testInput.To16() == nil {
						errors = append(errors, fmt.Errorf("[ERROR] %q %v isn't an IPv4 or IPv6", k, value))
					}
					return
				},
			},
			"mask": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 8 || value > 127 {
						errors = append(errors, fmt.Errorf("[ERROR] %q must be in the range from 8 to 127", k))
					}
					return
				},
			},
			"prio_master": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1 || value > 255 {
						errors = append(errors, fmt.Errorf("[ERROR] %q must be in the range from 1 to 255", k))
					}
					return
				},
			},
			"prio_slave": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1 || value > 255 {
						errors = append(errors, fmt.Errorf("[ERROR] %q must be in the range from 1 to 255", k))
					}
					return
				},
			},
			"vlan_device": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vrrp_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"iface_vrrp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id_vrrp": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1 || value > 255 {
						errors = append(errors, fmt.Errorf("[ERROR] %q must be in the range from 1 to 255", k))
					}
					return
				},
			},
			"auth_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "PASS" && value != "AH" {
						errors = append(errors, fmt.Errorf("[ERROR] %q must be PASS or AH", k))
					}
					return
				},
			},
			"auth_pass": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if strings.Count(value, "") > 8 {
						errors = append(errors, fmt.Errorf("[ERROR] %q %v too long", k, value))
					}
					return
				},
			},
			"post_up": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"default_gw": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					testInput := net.ParseIP(value)
					if testInput.To16() == nil {
						errors = append(errors, fmt.Errorf("[ERROR] %q %v isn't an IPv4 or IPv6", k, value))
					}
					return
				},
			},
			"lacp_slaves": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lacp_slaves_slave": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sync_iface": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"garp_m_delay": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1 || value > 10 {
						errors = append(errors, fmt.Errorf("[ERROR] %q must be in the range from 1 to 10", k))
					}
					return
				},
			},
			"advert_int": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1 || value > 10 {
						errors = append(errors, fmt.Errorf("[ERROR] %q must be in the range from 1 to 10", k))
					}
					return
				},
			},
			"garp_master_refresh": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 10 || value > 300 {
						errors = append(errors, fmt.Errorf("[ERROR] %q must be in the range from 10 to 300", k))
					}
					return
				},
			},
			"use_vmac": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceIfaceVrrpCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if len(d.Get("ip_vip").([]interface{})) != 0 {
		err := validateIPList(d)
		if err != nil {
			return err
		}
		setVrrpConfig(d, m)
	}
	if !d.Get("ip_vip_only").(bool) {
		if strings.Contains(d.Get("iface").(string), "vlan") && d.Get("vlan_device").(string) == "" {
			tfErr := d.Set("vlan_device", "bond1")
			if tfErr != nil {
				panic(tfErr)
			}
		}
		if (d.Get("lacp_slaves_slave").(string) == "") && (d.Get("lacp_slaves").(string) != "") {
			tfErr := d.Set("lacp_slaves_slave", d.Get("lacp_slaves").(string))
			if tfErr != nil {
				panic(tfErr)
			}
		}
		if len(d.Get("ip_vip").([]interface{})) != 0 {
			if d.Get("ip_master").(string) == "" {
				return fmt.Errorf("[ERROR] IP_vip_only = false so ip_master missing")
			}
			if d.Get("ip_slave").(string) == "" {
				return fmt.Errorf("[ERROR] IP_vip_only = false so ip_slave missing")
			}
			if d.Get("mask").(int) == 0 {
				return fmt.Errorf("[ERROR] IP_vip_only = false so mask missing")
			}
		}
	}
	IfaceVrrp := createStrucIfaceVrrp(d)
	_, err := client.requestAPI("ADD", &IfaceVrrp)
	if err != nil {
		return err
	}
	if len(d.Get("ip_vip").([]interface{})) == 0 {
		d.SetId(d.Get("iface").(string) + "_0")
	} else {
		d.SetId(d.Get("iface").(string) + "_" + strconv.Itoa(d.Get("id_vrrp").(int)))
	}
	return nil
}
func resourceIfaceVrrpRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	IfaceVrrp := createStrucIfaceVrrp(d)
	IfaceVrrpRead, err := client.requestAPI("CHECK", &IfaceVrrp)
	if err != nil {
		return err
	}
	if IfaceVrrpRead.Iface == "null" {
		d.SetId("")
	}
	if IfaceVrrpRead.IPMaster == "?" {
		tfErr := d.Set("ip_master", "")
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("mask", 0)
		if tfErr != nil {
			panic(tfErr)
		}
		_, exists := d.GetOk("post_up")
		if exists {
			tfErr = d.Set("post_up", []string{})
			if tfErr != nil {
				panic(tfErr)
			}
		}
		_, exists = d.GetOk("default_gw")
		if exists {
			tfErr = d.Set("default_gw", "")
			if tfErr != nil {
				panic(tfErr)
			}
		}
		_, exists = d.GetOk("lacp_slaves")
		if exists {
			tfErr = d.Set("lacp_slaves", "")
			if tfErr != nil {
				panic(tfErr)
			}
		}
		_, exists = d.GetOk("lacp_slaves_slave")
		if exists {
			tfErr = d.Set("lacp_slaves_slave", "")
			if tfErr != nil {
				panic(tfErr)
			}
		}
		_, exists = d.GetOk("vlan_device")
		if exists {
			tfErr = d.Set("vlan_device", "")
			if tfErr != nil {
				panic(tfErr)
			}
		}
	} else if len(IfaceVrrpRead.PostUp) == 1 {
		if IfaceVrrpRead.PostUp[0] == "?" {
			_, exists := d.GetOk("post_up")
			if exists {
				tfErr := d.Set("post_up", []string{})
				if tfErr != nil {
					panic(tfErr)
				}
			}
		}
	}
	if IfaceVrrpRead.IPSlave == "?" {
		tfErr := d.Set("ip_slave", "")
		if tfErr != nil {
			panic(tfErr)
		}
	}
	if IfaceVrrpRead.IDVrrp == "?" {
		tfErr := d.Set("vrrp_group", "")
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("ip_vip", []string{})
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("prio_master", 0)
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("auth_type", "")
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("auth_pass", "")
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("garp_m_delay", 0)
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("advert_int", 0)
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("garp_master_refresh", 0)
		if tfErr != nil {
			panic(tfErr)
		}
		_, exists := d.GetOk("sync_iface")
		if exists {
			tfErr = d.Set("sync_iface", "")
			if tfErr != nil {
				panic(tfErr)
			}
		}
	}
	if len(IfaceVrrpRead.IPVip) == 1 && IfaceVrrpRead.IPVip[0] == "?" {
		tfErr := d.Set("ip_vip", []string{})
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("prio_master", 0)
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("auth_type", "")
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("auth_pass", "")
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("garp_m_delay", 0)
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("advert_int", 0)
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("garp_master_refresh", 0)
		if tfErr != nil {
			panic(tfErr)
		}
		_, exists := d.GetOk("sync_iface")
		if exists {
			tfErr = d.Set("sync_iface", "")
			if tfErr != nil {
				panic(tfErr)
			}
		}
	}
	if IfaceVrrpRead.PrioSlave == "?" {
		tfErr := d.Set("prio_slave", 0)
		if tfErr != nil {
			panic(tfErr)
		}
	}
	return nil
}
func resourceIfaceVrrpUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	d.Partial(true)
	if d.HasChange("iface") {
		return fmt.Errorf("[ERROR] you can't change iface")
	}
	if d.HasChange("ip_master") {
		o, _ := d.GetChange("ip_master")
		if o.(string) != "" {
			return fmt.Errorf("[ERROR] you can't change ip_master")
		}
	}
	if d.HasChange("ip_slave") {
		o, _ := d.GetChange("ip_slave")
		if o.(string) != "" {
			return fmt.Errorf("[ERROR] you can't change ip_slave")
		}
	}
	if d.HasChange("mask") {
		o, _ := d.GetChange("mask")
		if o.(int) != 0 {
			return fmt.Errorf("[ERROR] you can't change mask")
		}
	}
	if d.HasChange("default_gw") {
		o, _ := d.GetChange("default_gw")
		if o.(string) != "" {
			return fmt.Errorf("[ERROR] you can't change default_gw")
		}
	}
	if d.HasChange("lacp_slaves") {
		o, _ := d.GetChange("lacp_slaves")
		if o.(string) != "" {
			return fmt.Errorf("[ERROR] you can't change lacp_slaves")
		}
	}
	if d.HasChange("lacp_slaves_slave") {
		o, _ := d.GetChange("lacp_slaves_slave")
		if o.(string) != "" {
			return fmt.Errorf("[ERROR] you can't change lacp_slaves_slave")
		}
	}
	if d.HasChange("vlan_device") {
		o, _ := d.GetChange("vlan_device")
		if o.(string) != "" {
			return fmt.Errorf("[ERROR] you can't change vlan_device")
		}
	}
	if len(d.Get("ip_vip").([]interface{})) != 0 {
		err := validateIPList(d)
		if err != nil {
			return err
		}
		setVrrpConfig(d, m)

	}
	if d.HasChange("ip_vip_only") {
		if !d.Get("ip_vip_only").(bool) {
			if strings.Contains(d.Get("iface").(string), "vlan") && d.Get("vlan_device").(string) == "" {
				tfErr := d.Set("vlan_device", "bond1")
				if tfErr != nil {
					panic(tfErr)
				}
			}
			if d.Get("ip_master").(string) == "" {
				return fmt.Errorf("[ERROR] ip_vip_only = false so ip_master missing")
			}
			if d.Get("ip_slave").(string) == "" {
				return fmt.Errorf("[ERROR] ip_vip_only = false so ip_slave missing")
			}
			if d.Get("mask").(int) == 0 {
				return fmt.Errorf("[ERROR] IP_vip_only = false so mask missing")
			}
		} else {
			tfErr := d.Set("ip_master", "")
			if tfErr != nil {
				panic(tfErr)
			}
			tfErr = d.Set("ip_slave", "")
			if tfErr != nil {
				panic(tfErr)
			}
			tfErr = d.Set("mask", 0)
			if tfErr != nil {
				panic(tfErr)
			}
			_, exists := d.GetOk("post_up")
			if exists {
				tfErr = d.Set("post_up", []string{})
				if tfErr != nil {
					panic(tfErr)
				}
			}
			_, exists = d.GetOk("default_gw")
			if exists {
				tfErr = d.Set("default_gw", "")
				if tfErr != nil {
					panic(tfErr)
				}
			}
			_, exists = d.GetOk("lacp_slaves")
			if exists {
				tfErr = d.Set("lacp_slaves", "")
				if tfErr != nil {
					panic(tfErr)
				}
			}
			_, exists = d.GetOk("lacp_slaves_slave")
			if exists {
				tfErr = d.Set("lacp_slaves_slave", "")
				if tfErr != nil {
					panic(tfErr)
				}
			}
			_, exists = d.GetOk("vlan_device")
			if exists {
				tfErr = d.Set("vlan_device", "")
				if tfErr != nil {
					panic(tfErr)
				}
			}
		}
	}
	IfaceVrrp := createStrucIfaceVrrp(d)
	if (len(d.Get("ip_vip").([]interface{})) != 0) && (d.HasChange("id_vrrp") || d.HasChange("iface_vrrp")) {
		oldID, newID := d.GetChange("id_vrrp")
		if oldID.(int) != 0 {
			err := client.requestAPIMove(&IfaceVrrp, oldID.(int))
			if err != nil {
				return err
			}
		} else {
			err := client.requestAPIMove(&IfaceVrrp, newID.(int))
			if err != nil {
				return err
			}
		}
		d.SetId(d.Get("iface").(string) + "_" + strconv.Itoa(d.Get("id_vrrp").(int)))
		d.SetPartial("id_vrrp")
		d.SetPartial("sync_iface")
	}
	_, err := client.requestAPI("CHANGE", &IfaceVrrp)
	if err != nil {
		return err
	}
	if len(d.Get("ip_vip").([]interface{})) == 0 {
		tfErr := d.Set("id_vrrp", 0)
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("prio_master", 0)
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("prio_slave", 0)
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("vrrp_group", "")
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("auth_type", "")
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("auth_pass", "")
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("garp_m_delay", 0)
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("advert_int", 0)
		if tfErr != nil {
			panic(tfErr)
		}
		tfErr = d.Set("garp_master_refresh", 0)
		if tfErr != nil {
			panic(tfErr)
		}
		d.SetId(d.Get("iface").(string) + "_0")
	}
	d.Partial(false)
	return nil
}
func resourceIfaceVrrpDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	IfaceVrrp := createStrucIfaceVrrp(d)
	_, err := client.requestAPI("REMOVE", &IfaceVrrp)
	if err != nil {
		return err
	}
	return nil
}

// validateIPList : validate list of cidr in ip_vip
func validateIPList(d *schema.ResourceData) error {
	var errors []string
	VIPList := make([]string, len(d.Get("ip_vip").([]interface{})))
	for i, d := range d.Get("ip_vip").([]interface{}) {
		VIPList[i] = d.(string)
	}
	for _, IP := range VIPList {
		testInput := net.ParseIP(IP)
		if testInput.To16() == nil {
			errors = append(errors, fmt.Sprintf("This VIP %v isn't an IPv4 or IPv6", IP))
		}
	}
	if len(errors) != 0 {
		return fmt.Errorf("[ERROR] " + strings.Join(errors, ", "))
	}
	return nil
}

//  createStrucIfaceVrrp prepare IfaceVrrp before call API
func createStrucIfaceVrrp(d *schema.ResourceData) ifaceVrrp {
	VIPList := make([]string, len(d.Get("ip_vip").([]interface{})))
	for i, d := range d.Get("ip_vip").([]interface{}) {
		VIPList[i] = d.(string)
	}
	PostupList := make([]string, len(d.Get("post_up").([]interface{})))
	for i, d := range d.Get("post_up").([]interface{}) {
		PostupList[i] = d.(string)
	}
	mask := strconv.Itoa(d.Get("mask").(int))
	PrioMaster := strconv.Itoa(d.Get("prio_master").(int))
	PrioSlave := strconv.Itoa(d.Get("prio_slave").(int))
	IDVrrp := strconv.Itoa(d.Get("id_vrrp").(int))
	GarpMDelay := strconv.Itoa(d.Get("garp_m_delay").(int))
	AdvertInt := strconv.Itoa(d.Get("advert_int").(int))
	GarpMasterRefresh := strconv.Itoa(d.Get("garp_master_refresh").(int))

	IfaceVrrp := ifaceVrrp{
		Iface:             d.Get("iface").(string),
		IPVip:             VIPList,
		IPVipOnly:         d.Get("ip_vip_only").(bool),
		IPMaster:          d.Get("ip_master").(string),
		IPSlave:           d.Get("ip_slave").(string),
		Mask:              mask,
		PrioMaster:        PrioMaster,
		PrioSlave:         PrioSlave,
		VlanDevice:        d.Get("vlan_device").(string),
		VrrpGroup:         d.Get("vrrp_group").(string),
		IfaceVrrp:         d.Get("iface_vrrp").(string),
		IDVrrp:            IDVrrp,
		AuthType:          d.Get("auth_type").(string),
		AuthPass:          d.Get("auth_pass").(string),
		PostUp:            PostupList,
		DefaultGW:         d.Get("default_gw").(string),
		LACPSlavesMaster:  d.Get("lacp_slaves").(string),
		LACPSlavesSlave:   d.Get("lacp_slaves_slave").(string),
		SyncIface:         d.Get("sync_iface").(string),
		GarpMDelay:        GarpMDelay,
		AdvertInt:         AdvertInt,
		GarpMasterRefresh: GarpMasterRefresh,
		UseVmac:           d.Get("use_vmac").(bool),
	}
	return IfaceVrrp
}

// setVrrpConfig : set default parameters (computed)
func setVrrpConfig(d *schema.ResourceData, m interface{}) {
	client := m.(*Client)
	if d.Get("id_vrrp").(int) == 0 {
		tfErr := d.Set("id_vrrp", client.getDefaultIDVrrp())
		if tfErr != nil {
			panic(tfErr)
		}
	}
	if d.Get("prio_master").(int) == 0 {
		tfErr := d.Set("prio_master", 150)
		if tfErr != nil {
			panic(tfErr)
		}
	}
	if d.Get("prio_slave").(int) == 0 {
		tfErr := d.Set("prio_slave", 100)
		if tfErr != nil {
			panic(tfErr)
		}
	}
	if d.Get("vrrp_group").(string) == "" {
		tfErr := d.Set("vrrp_group", client.getDefaultVrrpGroup())
		if tfErr != nil {
			panic(tfErr)
		}
	}
	if d.Get("auth_type").(string) == "" {
		tfErr := d.Set("auth_type", "PASS")
		if tfErr != nil {
			panic(tfErr)
		}
	}
	if d.Get("auth_pass").(string) == "" {
		tfErr := d.Set("auth_pass", "word")
		if tfErr != nil {
			panic(tfErr)
		}
	}
	if d.Get("garp_m_delay").(int) == 0 {
		tfErr := d.Set("garp_m_delay", 5)
		if tfErr != nil {
			panic(tfErr)
		}
	}
	if d.Get("advert_int").(int) == 0 {
		tfErr := d.Set("advert_int", client.getDefaultAdvertInt())
		if tfErr != nil {
			panic(tfErr)
		}
	}
	if d.Get("garp_master_refresh").(int) == 0 {
		tfErr := d.Set("garp_master_refresh", 60)
		if tfErr != nil {
			panic(tfErr)
		}
	}
}
