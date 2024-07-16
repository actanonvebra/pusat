package models

var CountryCode string
var Origin string
var IocingOS string

type IoCInformation struct {
	IP       string `json:"ip"`
	CveCount int    `json:"cve_count"`

	Asn         string `json:"asn"`
	CountryCode string `json:"country_code"`
	Os          string `json:"os"`
	PortData    []struct {
		Port     string `json:"port"`
		Protocol string `json:"protocol"`
	} `json:"port_data"`
}
