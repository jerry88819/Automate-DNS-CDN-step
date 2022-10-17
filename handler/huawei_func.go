package handler

import (
	"fmt"
	"log"

	"github.com/tidwall/gjson"

	dns "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/model"
)

func CreatePublicZone( clientW *dns.DnsClient, huaweiDomain string ) { // 這就是 origin domain 

	// 在華為雲裡面創建公共網域 創建過的話不會再次創建

	requestW := &model.CreatePublicZoneRequest{}
	requestW.Body = &model.CreatePublicZoneReq{
		Name: huaweiDomain,
	}

	responseW, err := clientW.CreatePublicZone(requestW)
	if err == nil {
		fmt.Printf("%+v\n", responseW)
	} else {
		fmt.Println(err)
	}
} // CreatePublicZone()

func FindPublicZoneId( clientW *dns.DnsClient, huaweiDomain string ) ( string ) {
	requestFindID := &model.ListPublicZonesRequest{}
	responseFindID, err := clientW.ListPublicZones(requestFindID)
	if err == nil {
		// fmt.Printf("%+v\n", responseFindID)
		fmt.Println("Find Public Zones Success!!!")
	} else {
		fmt.Println(err)
		log.Fatal(err)
	} // else()

	tempQuery := `zones.#(name=="` + huaweiDomain + `.").id`
	valueFindID := gjson.Get(responseFindID.String(), tempQuery )
	fmt.Println("這是目前公共網域的 id : ", valueFindID)

	return valueFindID.String()
} // FindPublicZoneId()

func CreateRecordSet( clientW *dns.DnsClient, public_zone_id string, txt_record string, huaweiDomain string ) error {
	requestRecordT := &model.CreateRecordSetRequest{}
	requestRecordT.ZoneId = public_zone_id // 公共網域 zone_id  域名 id

	var listRecordsbody1 = []string{
		"\"" + txt_record + "\"",
    }

	tttt := "_cdnauth" + "."+ huaweiDomain + "."

	fmt.Println(tttt)

	requestRecordT.Body = &model.CreateRecordSetReq{
		// Records: listRecordsbody1,
		Records: listRecordsbody1,
		// Type: data.Response.RecordType,
		Type: "TXT",
		Name: "_cdnauth" + "."+ huaweiDomain + ".",
		// Name: data.Response.SubDomain + "."+ huaweiDomain + ".",
	}

	response_TXT, err := clientW.CreateRecordSet(requestRecordT)
	if err == nil {
        fmt.Printf("%+v\n", response_TXT)
    } else {
        fmt.Println(err) // 代表此域名以驗證過 txt 紀錄以存在 所以才生成失敗
    }

	return err
} // CreateRecordSet()

func AddCname( clientW *dns.DnsClient, nowCname string, domain string, public_zone_id string ) error {
	requestRecord := &model.CreateRecordSetRequest{}
	requestRecord.ZoneId = public_zone_id // 公共網域 zone_id  域名 id

	var listRecordsbody = []string{
		nowCname,
	}

	requestRecord.Body = &model.CreateRecordSetReq{
		Records: listRecordsbody,
		Type:    "CNAME",
		Name:    domain+".",
	}
	responseRecord, err := clientW.CreateRecordSet(requestRecord)
	if err == nil {
		fmt.Printf("%+v\n", responseRecord)
	} else {
		fmt.Println(err)
	} // else()

	return err 
} // AddCname()