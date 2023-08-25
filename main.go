package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	sdk "github.com/aziontech/azionapi-go-sdk/idns"
)

const intelligentDnsURL = "https://api.azionapi.net/"

type Client struct {
	apiClient sdk.APIClient
}

type ResponseZone struct {
	Response []Results `json:"results"`
}

type Results struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	var personalToken string

	fmt.Println("Hey, there! Welcome to iDNS helper")

	fmt.Println("Please provide your Personal Token:")
	fmt.Scanf("%s", &personalToken)

	err := IDNShandler(personalToken)
	if err != nil {
		log.Fatal(err)
	}
}

func IDNShandler(personalToken string) error {
	var domainName, dnsZone string
	active := true

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please provide the Domain Name:")
	domainName, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	domainName = strings.TrimSpace(domainName)

	fmt.Println("Enter a DNS zone:")
	dnsZone, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	dnsZone = strings.TrimSpace(dnsZone)

	personalToken = strings.TrimSpace(personalToken)

	client := NewClient(intelligentDnsURL, personalToken)
	err = client.NewIdns(&domainName, &dnsZone, &active)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) NewIdns(domainName, domain *string, active *bool) error {
	ctx := context.Background()
	idns := &sdk.Zone{
		Name:     domainName,
		Domain:   domain,
		IsActive: active,
	}

	fmt.Println("\n--------------------------------------------------\n")
	fmt.Println("Creating iDNS zone...")

	req := c.apiClient.ZonesApi.PostZone(ctx).Zone(*idns)
	res, httpResp, err := req.Execute()
	if err != nil {
		bytes, readErr := io.ReadAll(httpResp.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		fmt.Println("\nError")
		fmt.Println(string(bytes))
		return err
	}

	response, err := res.MarshalJSON()
	if err != nil {
		return err
	}

	var createdZone ResponseZone
	if unmarshalErr := json.Unmarshal(response, &createdZone); unmarshalErr != nil {
		return unmarshalErr
	}

	fmt.Println("Zone created")
	fmt.Println("Zone ID:", createdZone.Response[0].ID)
	fmt.Println("Zone name:", createdZone.Response[0].Name)
	return nil
}

func NewClient(url string, token string) *Client {
	conf := sdk.NewConfiguration()

	conf.HTTPClient = &http.Client{}
	conf.AddDefaultHeader("Authorization", "token "+token)
	conf.AddDefaultHeader("Accept", "application/json;version=3")
	conf.Servers = sdk.ServerConfigurations{
		{URL: url},
	}

	return &Client{
		apiClient: *sdk.NewAPIClient(conf),
	}
}
