package config


import (
	"encoding/xml"
	"os"
	"path/filepath"
	"log"
	"strings"
	"io/ioutil"
)

type DefaultNames struct {
	ThemeName          string `xml:"DefaultNames>ThemeName"`
	NickName           string `xml:"DefaultNames>NickName"`
	Email              string `xml:"DefaultNames>Email"`
	AccountName        string `xml:"DefaultNames>AccountName"`
	Mtime              string `xml:"DefaultNames>Mtime"`
	DesignerUrl        string `xml:"DefaultNames>DesignerUrl"`
	HaveDesignerInfo   string `xml:"DefaultNames>HaveDesignerInfo"`
	NoDesignerInfo     string `xml:"DefaultNames>NoDesignerInfo"`
	EnHaveDesignerInfo string `xml:"DefaultNames>EnHaveDesignerInfo"`
	EnNoDesignerInfo   string `xml:"DefaultNames>EnNoDesignerInfo"`
}

type Globalconf struct {
	XMLName xml.Name `xml:"Config"`
	CoptTo []  string `xml:"CopyTo>Value"`
	CnTemplet  string `xml:"CnTemplet"`
	EnTemplet  string `xml:"EnTemplet"`
	Username   string `xml:"Username"`
	Password   string `xml:"Password"`
	Port       string `xml:"Port"`
	Threshold  int    `xml:"Threshold"`
	Subject    string `xml:"Subject"`
	EnSubject  string `xml:"EnSubject"`
	SmtpHost   string `xml:"SmtpHost"`
	From       string `xml:"From"`
	DefaultNames
}
var RunPath string
var Gf Globalconf
var CnTemplet string
var EnTemplet string

func getCurrentDirectory() string {  
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))  
    if err != nil {  
        log.Fatal(err)  
    }  
    return strings.Replace(dir, "\\", "/", -1)  
}  

func init () {
	RunPath = getCurrentDirectory()
	RunPath += "/"

	config, err := ioutil.ReadFile(RunPath+"conf/mail.xml")
	if err != nil {
		log.Fatal("ReadFile Conf:", err)
	}

	err = xml.Unmarshal(config, &Gf)
	if err != nil {
		log.Fatal("xml.Unmarsha1 err:", err)
	}

	cnbyte, err := ioutil.ReadFile(RunPath+Gf.CnTemplet)
	if err != nil {
		log.Fatal("ReadFile CnTemplet:", err)
	}
	CnTemplet = string(cnbyte)

	enbyte, err := ioutil.ReadFile(RunPath+Gf.EnTemplet)
	if err != nil {
		log.Fatal("ReadFile CnTemplet:", err)
	}
	EnTemplet = string(enbyte)
}