package lvsnetwork

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider lvsnetwork for terraform
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"firewall_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  8080,
			},
			"https": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"insecure": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"login": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vault_enable": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"login", "password"},
			},
			"vault_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "lvs",
			},
			"vault_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"default_id_vrrp": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1 || value > 255 {
						errors = append(errors, fmt.Errorf("%q must be in the range from 1 to 255", k))
					}
					return
				},
			},
			"default_vrrp_group": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "VG_1",
			},
			"default_advert_int": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1 || value > 10 {
						errors = append(errors, fmt.Errorf("[ERROR] %q must be in the range from 1 to 10", k))
					}
					return
				},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"lvsnetwork_ifacevrrp": resourceIfaceVrrp(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		firewallIP:       d.Get("firewall_ip").(string),
		firewallPortAPI:  d.Get("port").(int),
		https:            d.Get("https").(bool),
		insecure:         d.Get("insecure").(bool),
		logname:          os.Getenv("USER"),
		login:            d.Get("login").(string),
		password:         d.Get("password").(string),
		vaultEnable:      d.Get("vault_enable").(bool),
		vaultPath:        d.Get("vault_path").(string),
		vaultKey:         d.Get("vault_key").(string),
		defaultIDVrrp:    d.Get("default_id_vrrp").(int),
		defaultVrrpGroup: d.Get("default_vrrp_group").(string),
		defaultAdvertInt: d.Get("default_advert_int").(int),
	}
	return config.Client()
}
