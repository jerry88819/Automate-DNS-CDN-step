package handler

import (
	"fmt"
	"time"

	"Tencent_backstage_api/config"

	"github.com/gofiber/fiber/v2"

	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	dns "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2"

	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/region"
) // import()

// 在騰訊雲的 cdn 下建立 domain name  
func CreateTencentCdn( c *fiber.Ctx ) error {

	fmt.Println("In \"Create Tencent Cdn Function\"")

	str := new(Domain_Info) 
	err := c.BodyParser(str) // BodyParser can only parse into struct
	if err != nil {
		fmt.Println(err)
		return err
	} // if()

	fmt.Println("json domain : ", str.Domain_name)
	fmt.Println("json cdn : ", str.Cdn_str)

	random_password := CreatePasswordSixDigit() //  6位數亂碼
	origin_domain := str.Domain_name // 為添加任何元素的 domain name
	domain_with_pass := random_password + "." + origin_domain // ex : 1q2d5a.google.com

	fmt.Println("Domain with pass : ", domain_with_pass )

	// Tencent Client

	credential := common.NewCredential(
		config.Config("tencent_account_key"),
		config.Config("tencent_password_key"),
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cdn.tencentcloudapi.com"
	client, _ := cdn.NewClient(credential, "", cpf)

	// Huawei ClientW

	auth := basic.NewCredentialsBuilder().
		WithAk(config.Config("ak")).
		WithSk(config.Config("sk")).
		WithProjectId(config.Config("RegionProjectId")).
		Build()

	clientW := dns.NewDnsClient(
		dns.DnsClientBuilder().
			WithRegion(region.ValueOf("cn-east-3")).
			WithCredential(auth).
			Build())


	// 先去找認證的 txt record
	verify_string := GetVerifyRecord( client, domain_with_pass )

	// 在華為新增 Create Public Zone, 重複的話沒關係
	CreatePublicZone( clientW, origin_domain )

	//  找尋此 public zone id
	public_zone_id := FindPublicZoneId( clientW, origin_domain )

	// 用此 public zone id 去添加 txt 紀錄 
	exist := false
	err = CreateRecordSet( clientW, public_zone_id, verify_string, origin_domain )
	if err != nil {
		fmt.Println("The txt record has already exist!!! Don't need to add again!!!")
		exist = true
	} // if()

	// 到 tencent cloud cdn 去做驗證 

	if !exist { //  假如 txt 驗證原本不存在的話, 才需要驗證

		success := false

		for !success {

			time.Sleep(20 * time.Second)
			err = VerifyRecordFunc( client, domain_with_pass )
			if err == nil {
				success = true
			} else { // if()
				fmt.Println("Verify record again! Please wait...")
			} // else()

		} // for()

	} // if()

	// Add domain into Tencent cloud cdn
	err = AddDomainToCdn( client, domain_with_pass, str.Cdn_str )
	if err != nil {
		return c.JSON("Add Domain into cdn failed!")
	} // if()

	// 搜尋剛剛所添加網域的 CNAME 
	cname_str, err := QueryCname(client, domain_with_pass)
	if err != nil {
		return c.JSON("Failed to quering cname!")
	} // if()

	// 新增 Cname 至 huawei cloud dns public zone
	err = AddCname( clientW, cname_str, domain_with_pass, public_zone_id )
	if err != nil {
		return c.JSON("Failed to adding cname in public zone!")
	} // if()

	return c.JSON("Create Tencent Cdn success!")
} // CreateTencentCdn()

// 這是用來清緩存的 單一清理
func PurgeSingleDomain( c *fiber.Ctx ) error { 
	fmt.Println("In \"Purge Single Domain name function\" !")
	domain_name := c.Query("domain") // 給出要緩存的域名

	var domain_array []string
	domain_array = append(domain_array, domain_name) 

	// Tencent Client

	credential := common.NewCredential(
		config.Config("tencent_account_key"),
		config.Config("tencent_password_key"),
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cdn.tencentcloudapi.com"
	client, _ := cdn.NewClient(credential, "", cpf)

	err := PurgeDomain( client, domain_array )
	if err != nil {
		return c.JSON("Failed to purge domain.")
	} // if()

	return c.JSON("Purge domain success!")
} // PurgeSingleDomain()

type Domain_Info struct {
	Domain_name		string		`json:"domain_name"`
	Cdn_str			string		`json:"cdn"`
} // Domain_Info()