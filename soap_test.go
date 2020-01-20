package gosoap

import (
	"net/http"
	"testing"
)

var (
	scts = []struct {
		URL string
		Err bool
	}{
		{
			URL: "://www.server",
			Err: false,
		},
		{
			URL: "",
			Err: false,
		},
		{
			URL: "http://ec.europa.eu/taxation_customs/vies/checkVatService.wsdl",
			Err: true,
		},
	}
)

func TestSoapClient(t *testing.T) {
	for _, sct := range scts {
		_, err := SoapClient(sct.URL, http.Client{})
		if err != nil && sct.Err {
			t.Errorf("URL: %s - error: %s", sct.URL, err)
		}
	}
}

type CheckVatRequest struct {
	CountryCode string
	VatNumber   string
}

func (r CheckVatRequest) SoapBuildRequest() *Request {
	return NewRequest("checkVat", Params{
		"countryCode": r.CountryCode,
		"vatNumber":   r.VatNumber,
	})
}

type CheckVatResponse struct {
	CountryCode string `xml:"countryCode"`
	VatNumber   string `xml:"vatNumber"`
	RequestDate string `xml:"requestDate"`
	Valid       string `xml:"valid"`
	Name        string `xml:"name"`
	Address     string `xml:"address"`
}

type CapitalCityResponse struct {
	CapitalCityResult string
}

type NumberToWordsResponse struct {
	NumberToWordsResult string
}

type WhoisResponse struct {
	WhoisResult string
}

var (
	rv CheckVatResponse
	rc CapitalCityResponse
	rn NumberToWordsResponse
	rw WhoisResponse

	params = Params{}
)

func TestClient_Call(t *testing.T) {
	soap, err := SoapClient("http://ec.europa.eu/taxation_customs/vies/checkVatService.wsdl", http.Client{})
	if err != nil {
		t.Errorf("error not expected: %s", err)
	}

	var res *Response

	params["vatNumber"] = "6388047V"
	params["countryCode"] = "IE"
	res, err = soap.Call("", params)
	if err == nil {
		t.Errorf("method is empty")
	}

	if res != nil {
		t.Errorf("body is empty")
	}

	res, err = soap.Call("checkVat", params)
	if err != nil {
		t.Errorf("error in soap call: %s", err)
	}

	res.Unmarshal(&rv)
	if rv.CountryCode != "IE" {
		t.Errorf("error: %+v", rv)
	}

	soap, err = SoapClient("http://webservices.oorsprong.org/websamples.countryinfo/CountryInfoService.wso?WSDL", http.Client{})
	if err != nil {
		t.Errorf("error not expected: %s", err)
	}

	res, err = soap.Call("CapitalCity", Params{"sCountryISOCode": "GB"})
	if err != nil {
		t.Errorf("error in soap call: %s", err)
	}

	res.Unmarshal(&rc)

	if rc.CapitalCityResult != "London" {
		t.Errorf("error: %+v", rc)
	}

	soap, err = SoapClient("http://www.dataaccess.com/webservicesserver/numberconversion.wso?WSDL", http.Client{})
	if err != nil {
		t.Errorf("error not expected: %s", err)
	}

	res, err = soap.Call("NumberToWords", Params{"ubiNum": "23"})
	if err != nil {
		t.Errorf("error in soap call: %s", err)
	}

	res.Unmarshal(&rn)

	if rn.NumberToWordsResult != "twenty three " {
		t.Errorf("error: %+v", rn)
	}

	soap, err = SoapClient("https://domains.livedns.co.il/API/DomainsAPI.asmx?WSDL", http.Client{})
	if err != nil {
		t.Errorf("error not expected: %s", err)
	}

	res, err = soap.Call("Whois", Params{"DomainName": "google.com"})
	if err != nil {
		t.Errorf("error in soap call: %s", err)
	}

	res.Unmarshal(&rw)

	if rw.WhoisResult != "0" {
		t.Errorf("error: %+v", rw)
	}

	c := &Client{}
	res, err = c.Call("", Params{})
	if err == nil {
		t.Errorf("error expected but nothing got.")
	}

	c.SetWSDL("://test.")

	res, err = c.Call("checkVat", params)
	if err == nil {
		t.Errorf("invalid WSDL")
	}
}

func TestClient_CallByStruct(t *testing.T) {
	soap, err := SoapClient("http://ec.europa.eu/taxation_customs/vies/checkVatService.wsdl", http.Client{})
	if err != nil {
		t.Errorf("error not expected: %s", err)
	}

	var res *Response

	res, err = soap.CallByStruct(CheckVatRequest{
		CountryCode: "IE",
		VatNumber:   "6388047V",
	})
	if err != nil {
		t.Errorf("error in soap call: %s", err)
	}

	res.Unmarshal(&rv)
	if rv.CountryCode != "IE" {
		t.Errorf("error: %+v", rv)
	}
}

func TestClient_Call_NonUtf8(t *testing.T) {
	soap, err := SoapClient("https://demo.ilias.de/webservice/soap/server.php?wsdl", http.Client{})
	if err != nil {
		t.Errorf("error not expected: %s", err)
	}

	_, err = soap.Call("login", Params{"client": "demo", "username": "robert", "password": "iliasdemo"})
	if err != nil {
		t.Errorf("error in soap call: %s", err)
	}
}

func TestProcess_doRequest(t *testing.T) {
	c := &process{
		Client: &Client{
			HttpClient: &http.Client{},
		},
	}

	_, err := c.doRequest("")
	if err == nil {
		t.Errorf("body is empty")
	}

	_, err = c.doRequest("://teste.")
	if err == nil {
		t.Errorf("invalid WSDL")
	}
}
