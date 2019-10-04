package lvsnetwork

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Client = provider configuration
type Client struct {
	HTTPS            bool
	Insecure         bool
	DefaultIDVrrp    int
	DefaultAdvertInt int
	Port             int
	DefaultVrrpGroup string
	FirewallIP       string
	Logname          string
	Login            string
	Password         string
}
type ifaceVrrp struct {
	IPVipOnly         bool     `json:"IP_vip_only"`
	UseVmac           bool     `json:"Use_vmac"`
	Iface             string   `json:"iface"`
	IPMaster          string   `json:"IP_master"`
	IPSlave           string   `json:"IP_slave"`
	Mask              string   `json:"Mask"`
	PrioMaster        string   `json:"Prio_master"`
	PrioSlave         string   `json:"Prio_slave"`
	VlanDevice        string   `json:"Vlan_device"`
	VrrpGroup         string   `json:"Vrrp_group"`
	IfaceVrrp         string   `json:"Iface_vrrp"`
	IDVrrp            string   `json:"Id_vrrp"`
	AuthType          string   `json:"Auth_type"`
	AuthPass          string   `json:"Auth_pass"`
	DefaultGW         string   `json:"Default_GW"`
	LACPSlavesMaster  string   `json:"LACP_slaves_master"`
	LACPSlavesSlave   string   `json:"LACP_slaves_slave"`
	SyncIface         string   `json:"Sync_iface"`
	GarpMDelay        string   `json:"Garp_m_delay"`
	AdvertInt         string   `json:"Advert_int"`
	GarpMasterRefresh string   `json:"Garp_master_refresh"`
	IPVip             []string `json:"IP_vip"`
	PostUp            []string `json:"Post_up"`
}

// NewClient configure
func NewClient(firewallIP string, firewallPortAPI int, https bool, insecure bool, logname string,
	login string, password string, defaultIDVrrp int, defaultVrrpGroup string, defaultAdvertInt int) *Client {
	client := &Client{
		FirewallIP:       firewallIP,
		Port:             firewallPortAPI,
		HTTPS:            https,
		Insecure:         insecure,
		Logname:          logname,
		Login:            login,
		Password:         password,
		DefaultIDVrrp:    defaultIDVrrp,
		DefaultVrrpGroup: defaultVrrpGroup,
		DefaultAdvertInt: defaultAdvertInt,
	}
	return client
}

// getDefaultIDVrrp : get provider config for computed id_vrrp resource parameter
func (client *Client) getDefaultIDVrrp() int {
	return client.DefaultIDVrrp
}

// getDefaultVrrpGroup : get provider config for computed vrrp_group resource parameter
func (client *Client) getDefaultVrrpGroup() string {
	return client.DefaultVrrpGroup
}

// getDefaultAdvertInt : get provider config for computed advert_int resource parameter
func (client *Client) getDefaultAdvertInt() int {
	return client.DefaultAdvertInt
}

// newRequest : call API
func (client *Client) newRequest(uri string, ifaceVrrpReq *ifaceVrrp) (int, string, error) {
	urlString := "http://" + client.FirewallIP + ":" + strconv.Itoa(client.Port) + uri + "?&logname=" + client.Logname
	if client.HTTPS {
		urlString = strings.Replace(urlString, "http://", "https://", -1)
	}
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(ifaceVrrpReq)
	if err != nil {
		return 500, "", err
	}
	req, err := http.NewRequest("POST", urlString, body)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	if client.Login != "" && client.Password != "" {
		req.SetBasicAuth(client.Login, client.Password)
	}
	if err != nil {
		return 500, "", err
	}
	tr := &http.Transport{
		DisableKeepAlives: true,
	}
	if client.Insecure {
		tr = &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		}
	}
	httpClient := &http.Client{Transport: tr}
	log.Printf("[DEBUG] Request API (%v) %v", urlString, body)
	resp, err := httpClient.Do(req)
	if err != nil {
		return 500, "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 500, "", err
	}
	log.Printf("[DEBUG] Response API (%v) %v => %v", urlString, resp.StatusCode, string(respBody))
	return resp.StatusCode, string(respBody), nil
}

// requestAPI : prepare request to API and call api with newRequest()
func (client *Client) requestAPI(action string, ifaceVrrpReq *ifaceVrrp) (ifaceVrrp, error) {
	var ifaceVrrpReturn ifaceVrrp
	switch action {
	case "ADD":
		uriString := "/add_iface_vrrp/" + ifaceVrrpReq.Iface + "/"
		statuscode, body, err := client.newRequest(uriString, ifaceVrrpReq)
		if err != nil {
			return ifaceVrrpReturn, err
		}
		if statuscode == 401 {
			return ifaceVrrpReturn, fmt.Errorf("you are Unauthorized")
		}
		if statuscode != 200 {
			return ifaceVrrpReturn, fmt.Errorf(body)
		}
		return ifaceVrrpReturn, nil
	case "REMOVE":
		uriString := "/remove_iface_vrrp/" + ifaceVrrpReq.Iface + "/"
		statuscode, body, err := client.newRequest(uriString, ifaceVrrpReq)
		if err != nil {
			return ifaceVrrpReturn, err
		}
		if statuscode == 401 {
			return ifaceVrrpReturn, fmt.Errorf("you are Unauthorized")
		}
		if statuscode != 200 {
			return ifaceVrrpReturn, fmt.Errorf(body)
		}
		return ifaceVrrpReturn, nil
	case "CHECK":
		uriString := "/check_iface_vrrp/" + ifaceVrrpReq.Iface + "/"
		statuscode, body, err := client.newRequest(uriString, ifaceVrrpReq)
		if err != nil {
			return ifaceVrrpReturn, err
		}
		if statuscode == 401 {
			return ifaceVrrpReturn, fmt.Errorf("you are Unauthorized")
		}
		if statuscode == 404 {
			ifaceVrrpReturn.Iface = "null"
			return ifaceVrrpReturn, nil
		}

		errDecode := json.Unmarshal([]byte(body), &ifaceVrrpReturn)
		if errDecode != nil {
			return ifaceVrrpReturn, fmt.Errorf("[ERROR] decode json API response (%v) %v", errDecode, body)
		}
		return ifaceVrrpReturn, nil
	case "CHANGE":
		uriString := "/change_iface_vrrp/" + ifaceVrrpReq.Iface + "/"
		statuscode, body, err := client.newRequest(uriString, ifaceVrrpReq)
		if err != nil {
			return ifaceVrrpReturn, err
		}
		if statuscode == 401 {
			return ifaceVrrpReturn, fmt.Errorf("you are Unauthorized")
		}
		if statuscode != 200 {
			return ifaceVrrpReturn, fmt.Errorf(body)
		}
		return ifaceVrrpReturn, nil
	default:
		return ifaceVrrpReturn, fmt.Errorf("internal error => unknown action for requestAPI")
	}
}

// requestAPIMove : call /moveid_iface_vrrp/ on api
func (client *Client) requestAPIMove(ifaceVrrpReq *ifaceVrrp, oldID int) error {
	uriString := "/moveid_iface_vrrp/" + ifaceVrrpReq.Iface + "/" + strconv.Itoa(oldID) + "/"
	statuscode, body, err := client.newRequest(uriString, ifaceVrrpReq)
	if err != nil {
		return err
	}
	if statuscode == 401 {
		return fmt.Errorf("you are Unauthorized")
	}
	if statuscode != 200 {
		return fmt.Errorf(body)
	}
	return nil
}
