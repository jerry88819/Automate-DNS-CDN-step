# 使用騰訊雲自動化 DNS 和 CDN 步驟

## 將域名添加到騰訊雲 CDN

要自動將域名添加到騰訊雲 CDN，可以使用以下 API 端點：

/tencent/test


### 需要的參數

1. **Domain_name**：您在聚名網 (Juming.com) 上選擇的域名。
2. **Cdn_str**：CDN 要指向的位置。

## 清理緩存

要清理某個域名的緩存，可以使用以下 API 端點：

/tencent/purge

### 需要的參數

- **Domain_name**：需要清理緩存的域名。

---
