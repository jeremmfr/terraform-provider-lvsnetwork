package lvsnetwork

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceIfaceVrrp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIfaceVrrpCreate,
		ReadContext:   resourceIfaceVrrpRead,
		UpdateContext: resourceIfaceVrrpUpdate,
		DeleteContext: resourceIfaceVrrpDelete,

		Schema: map[string]*schema.Schema{
			"iface": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"ip_slave": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"mask": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(8, 127),
			},
			"prio_master": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 255),
			},
			"prio_slave": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 255),
			},
			"vlan_device": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 255),
			},
			"auth_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PASS", "AH"}, false),
			},
			"auth_pass": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 7),
			},
			"post_up": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"default_gw": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"lacp_slaves": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"lacp_slaves_slave": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"sync_iface": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"garp_m_delay": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"advert_int": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"garp_master_refresh": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(10, 300),
			},
			"use_vmac": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"track_script": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func resourceIfaceVrrpCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if len(d.Get("ip_vip").([]interface{})) != 0 {
		err := validateIPList(d)
		if err != nil {
			return diag.FromErr(err)
		}
		setVrrpConfig(d, m)
	} else if tfErr := d.Set("track_script", []string{}); tfErr != nil {
		panic(tfErr)
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
				return diag.FromErr(fmt.Errorf("[ERROR] IP_vip_only = false so ip_master missing"))
			}
			if d.Get("ip_slave").(string) == "" {
				return diag.FromErr(fmt.Errorf("[ERROR] IP_vip_only = false so ip_slave missing"))
			}
			if d.Get("mask").(int) == 0 {
				return diag.FromErr(fmt.Errorf("[ERROR] IP_vip_only = false so mask missing"))
			}
		}
	}
	IfaceVrrp := createStrucIfaceVrrp(d)
	_, err := client.requestAPIIFaceVrrp(ctx, ADD, &IfaceVrrp)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(d.Get("ip_vip").([]interface{})) == 0 {
		d.SetId(d.Get("iface").(string) + "_0")
	} else {
		d.SetId(d.Get("iface").(string) + "_" + strconv.Itoa(d.Get("id_vrrp").(int)))
	}

	return nil
}

func resourceIfaceVrrpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	IfaceVrrp := createStrucIfaceVrrp(d)
	IfaceVrrpRead, err := client.requestAPIIFaceVrrp(ctx, CHECK, &IfaceVrrp)
	if err != nil {
		return diag.FromErr(err)
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
	} else if len(IfaceVrrpRead.PostUp) == 1 && IfaceVrrpRead.PostUp[0] == "?" {
		_, exists := d.GetOk("post_up")
		if exists {
			tfErr := d.Set("post_up", []string{})
			if tfErr != nil {
				panic(tfErr)
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
		tfErr = d.Set("track_script", []string{})
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
		tfErr = d.Set("track_script", []string{})
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

func resourceIfaceVrrpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	d.Partial(true)
	if len(d.Get("ip_vip").([]interface{})) != 0 {
		err := validateIPList(d)
		if err != nil {
			return diag.FromErr(err)
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
				return diag.FromErr(fmt.Errorf("[ERROR] ip_vip_only = false so ip_master missing"))
			}
			if d.Get("ip_slave").(string) == "" {
				return diag.FromErr(fmt.Errorf("[ERROR] ip_vip_only = false so ip_slave missing"))
			}
			if d.Get("mask").(int) == 0 {
				return diag.FromErr(fmt.Errorf("[ERROR] IP_vip_only = false so mask missing"))
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
			err := client.requestAPIIFaceVrrpMove(ctx, &IfaceVrrp, oldID.(int))
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			err := client.requestAPIIFaceVrrpMove(ctx, &IfaceVrrp, newID.(int))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		d.SetId(d.Get("iface").(string) + "_" + strconv.Itoa(d.Get("id_vrrp").(int)))
		d.Partial(false)
	}
	_, err := client.requestAPIIFaceVrrp(ctx, CHANGE, &IfaceVrrp)
	if err != nil {
		return diag.FromErr(err)
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
		tfErr = d.Set("track_script", []string{})
		if tfErr != nil {
			panic(tfErr)
		}
		d.SetId(d.Get("iface").(string) + "_0")
	}
	d.Partial(false)

	return nil
}

func resourceIfaceVrrpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	IfaceVrrp := createStrucIfaceVrrp(d)
	_, err := client.requestAPIIFaceVrrp(ctx, REMOVE, &IfaceVrrp)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// validateIPList : validate list of cidr in ip_vip.
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

//  createStrucIfaceVrrp prepare IfaceVrrp before call API.
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
	TrackScript := make([]string, 0)
	for _, elem := range d.Get("track_script").([]interface{}) {
		TrackScript = append(TrackScript, elem.(string))
	}

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
		TrackScript:       TrackScript,
	}

	return IfaceVrrp
}

// setVrrpConfig : set default parameters (computed).
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
		tfErr := d.Set("auth_pass", client.getDefaultAuthPass())
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
	if len(d.Get("track_script").([]interface{})) == 0 {
		tfErr := d.Set("track_script", client.getDefaultTrackScript())
		if tfErr != nil {
			panic(tfErr)
		}
	}
}
