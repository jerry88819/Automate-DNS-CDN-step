package handler

import (
	"Tencent_backstage_api/config"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tidwall/gjson"
) // import()

const ( // 生成密碼使用
	MixStr = "0123456789abcdefghijklmnopqrstuvwxyz"
) // const()

type VerifyRecordInfo struct {
	SubDomain  string `json:"SubDomain"`  // _cdnauth
	Record     string `json:"Record"`     // 212135465136....
	RecordType string `json:"RecordType"` // TXT
	RequestId  string `json:"RequestId"`
} // VerifyRecordInfo()

type VerifyRecord struct {
	Response VerifyRecordInfo `json:"Response"`
} // VerifyRecord()

func CreatePasswordSixDigit() string {
	rand.Seed(time.Now().UnixNano())

	var ans []byte = make([]byte, 6) // 6 為要生成之密碼的長度
	for i := 0; i < 6; i++ {
		index := rand.Intn(len(MixStr))
		ans[i] = MixStr[index]
	} // for()

	return string(ans)
} // CreatePasswordSixDigit()

func GetVerifyRecord(client *cdn.Client, domain string) string {

	fmt.Println("In Verify Record Function!!!")

	reqq := cdn.NewCreateVerifyRecordRequest()

	reqq.Domain = common.StringPtr(domain)

	rep, err := client.CreateVerifyRecord(reqq)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s \n", err)
		log.Fatal(err)
	}

	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", rep.ToJsonString())

	data := VerifyRecord{}
	json.Unmarshal([]byte(rep.ToJsonString()), &data)

	fmt.Println(data.Response.Record)

	return data.Response.Record
} // GetVerifyRecord()

func VerifyRecordFunc(client *cdn.Client, domain string) error {
	requestV := cdn.NewVerifyDomainRecordRequest()

	requestV.Domain = common.StringPtr(domain)

	responseV, err := client.VerifyDomainRecord(requestV)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
	}

	if err != nil {
		panic(err)
	}

	fmt.Println(responseV.ToJsonString())

	return err
} //  VerifyRecordFunc()

func AddDomainToCdn( client *cdn.Client, domain string, cdn_abbreviation string ) error {

	cdnstr := config.Config(cdn_abbreviation) + ":80"

	request := cdn.NewAddCdnDomainRequest()

	request.Domain = common.StringPtr(domain)
	request.ServiceType = common.StringPtr("web")
	request.Origin = &cdn.Origin{
		Origins:            common.StringPtrs([]string{cdnstr}),
		OriginType:         common.StringPtr("domain"),
		OriginPullProtocol: common.StringPtr("http"),
	}

	response, err := client.AddCdnDomain(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return err
	}

	if err != nil {
		return err
	}

	fmt.Printf("%s", response.ToJsonString())
	fmt.Println("成功添加網域至騰訊雲CDN!!!")

	return err
} // AddDomainToCdn()

func QueryCname( client *cdn.Client, domain string ) ( string, error ) {
	requestDomainInfo := cdn.NewDescribeDomainsRequest()

	requestDomainInfo.Filters = []*cdn.DomainFilter{
		&cdn.DomainFilter{
			Name:  common.StringPtr("domain"),
			Value: common.StringPtrs([]string{domain}), // mxw4kh.jifuzaixian.com
		},
	}

	responseDomianInfo, err := client.DescribeDomains(requestDomainInfo)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
	} // if()
	if err != nil {
		return "Query Cname failed!", err
	} // if()

	value := gjson.Get(responseDomianInfo.ToJsonString(), "Response.Domains.0.Cname")
	println(value.String())
	nowCname := value.String()

	return nowCname, err
} // QueryCname()

func PurgeDomain( client *cdn.Client, domain []string ) error {

	request := cdn.NewPurgePathCacheRequest()

	request.Paths = common.StringPtrs(domain)
	request.FlushType = common.StringPtr("flush") // flush：刷新产生更新的资源 	delete：刷新全部资源 

	response, err := client.PurgePathCache(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
	} // if()
	
	if err != nil {
		return err
	} // if()

	fmt.Println( response.ToJsonString() )
	return err
} // PurgeDomain()
