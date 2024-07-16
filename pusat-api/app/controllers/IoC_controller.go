package controllers

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/gofiber/fiber/v2"
	"github.com/likexian/whois"
)

func GetIOCs(c *fiber.Ctx) error {
	var wg sync.WaitGroup

	ioc := c.Query("ioc")
	fmt.Println(ioc)

	iocing := models.IoCInformation{}
	models.CountryCode = ""
	models.Origin = ""
	models.IocingOS = ""

	wg.Add(1)
	go osDetection(&wg, ioc)
	fmt.Println("osDetection ok")

	cmd := exec.Command("nmap", "-sV", "-Pn", "--top-port=50", "--script=vulscan/vulscan.nse", ioc, "--host-timeout", "600")

	out, err := cmd.CombinedOutput()
	str := string(out)
	fmt.Println(str)

	iocing.IP = ioc

	cve := regexp.MustCompile(`\[CVE-([0-9]{4}-[0-9]{4})\]`)
	port := regexp.MustCompile(`(\d+)/(tcp|udp)`)

	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		line := scanner.Text()

		matchesCVE := cve.FindAllStringSubmatch(line, -1)
		matchesPORT := port.FindAllStringSubmatch(line, -1)

		if matchesCVE != nil || matchesPORT != nil {
			for _, matchC := range matchesCVE {
				if len(matchC) >= 2 {
					fmt.Println(matchC[0])
					iocing.CveCount++
				}
				break
			}
			for _, matchP := range matchesPORT {
				if len(matchP) >= 2 {
					fmt.Println(matchP[0])
					splitedPortProtocol := strings.Split(matchP[0], "/")
					iocing.PortData = append(iocing.PortData, struct {
						Port     string `json:"port"`
						Protocol string `json:"protocol"`
					}{
						Port:     splitedPortProtocol[0],
						Protocol: splitedPortProtocol[1],
					})
				}
				break
			}

		}
	}
	fmt.Println("wg wait")
	wg.Wait()
	fmt.Println("wg done")

	if err = scanner.Err(); err != nil {
		return c.JSON(fiber.Map{
			"error": true,
			"msg":   "Scanner Error",
			"data":  nil,
		})
	}
	fmt.Println("scaning ok")
	getWhois(ioc)
	fmt.Println("whois ok")

	iocing.CountryCode = models.CountryCode
	iocing.Os = models.IocingOS
	iocing.Asn = models.Origin

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"data":  iocing,
	})

}

func getWhois(ioc string) (string, error) {
	result, err := whois.Whois(ioc)
	if err != nil {
		fmt.Println(err)
	}
	pattern := regexp.MustCompile(`(?m)^country:.*$`)
	originPattern := regexp.MustCompile(`(?m)^origin:.*$`)

	matches := pattern.FindAllString(result, -1)
	originMatches := originPattern.FindAllString(result, -1)

	for _, match := range matches {
		parts := strings.Fields(match)
		models.CountryCode = strings.TrimSpace(parts[1])
	}
	if models.CountryCode == "" {
		pattern = regexp.MustCompile(`(?m)^Country:.*$`)
		matches = pattern.FindAllString(result, -1)
		for _, match := range matches {
			parts := strings.Fields(match)
			models.CountryCode = strings.TrimSpace(parts[1])
		}
	}

	for _, match := range originMatches {
		parts := strings.Fields(match)
		models.Origin = strings.TrimSpace(parts[1])
	}
	if models.Origin == "" {
		originPattern = regexp.MustCompile(`(?m)^Origin:.*$`)
		originMatches = originPattern.FindAllString(result, -1)
		for _, match := range originMatches {
			parts := strings.Fields(match)
			models.Origin = strings.TrimSpace(parts[1])
		}
	}

	return result, err
}

func osDetection(wg *sync.WaitGroup, ioc string) {
	cmd := exec.Command("nmap", "-O", ioc, "--host-timeout", "300")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("CombinedOutput error:", err)
		return
	}

	str := string(out)

	regex := regexp.MustCompile(`\b(Linux|Windows|Unix)\b`)
	finds := regex.FindAllString(str, -1)

	if len(finds) > 0 {
		models.IocingOS = finds[0]
	}
	wg.Done()
}
