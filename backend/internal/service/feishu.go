package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
	larksheets "github.com/larksuite/oapi-sdk-go/v3/service/sheets/v3"
	"io/ioutil"
	"net/http"
	"shengcai/internal/model"
	"shengcai/internal/repository"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FeiShuService interface {
	GetTenantAccessToken(ctx context.Context, appID string, appSecret string) (string, error)
	GetFirstSheetIDAndRowCount(ctx context.Context, appID string, appSecret string, spreadsheetToken string) (string, int, error)
	SaveTableData(ctx context.Context, appID string, appSecret string, spreadsheetToken string) error
	GetDocumentData(ctx context.Context, appID string, appSecret string, documentMetaData map[string]string) (string, error)
	GetSheetLatestModifyTime(ctx context.Context, appID string, appSecret string, spreadsheetToken string) (string, error)
}

func NewFeiShuService(
	service *Service,
	sheetInfoRepo repository.SheetInfoRepository,
	cellDataRepo repository.CellDataRepository,
) FeiShuService {
	return &feiShuService{
		Service:       service,
		sheetInfoRepo: sheetInfoRepo,
		cellDataRepo:  cellDataRepo,
	}
}

type feiShuService struct {
	*Service
	sheetInfoRepo repository.SheetInfoRepository
	cellDataRepo  repository.CellDataRepository
}

func (s *feiShuService) GetTenantAccessToken(ctx context.Context, appID string, appSecret string) (string, error) {
	url := "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"

	// 请求体数据
	data := map[string]string{
		"app_id":     appID,
		"app_secret": appSecret,
	}

	// 将数据编码为 JSON 格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error encoding JSON: %w", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// 定义响应体的结构
	var response struct {
		Code              int    `json:"code"`
		Msg               string `json:"msg"`
		TenantAccessToken string `json:"tenant_access_token"`
		Expire            int    `json:"expire"`
	}

	// 解析 JSON 响应体
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error decoding JSON: %w", err)
	}

	// 判断请求是否成功
	if response.Code == 0 && response.Msg == "ok" {
		return response.TenantAccessToken, nil
	} else {
		return "", fmt.Errorf("request failed with code: %d, message: %s", response.Code, response.Msg)
	}
}

func (s *feiShuService) GetFirstSheetIDAndRowCount(ctx context.Context, appID string, appSecret string, spreadsheetToken string) (string, int, error) {
	// 创建 Client
	client := lark.NewClient(appID, appSecret)

	// 创建请求对象
	req := larksheets.NewQuerySpreadsheetSheetReqBuilder().
		SpreadsheetToken(spreadsheetToken).
		Build()

	// 发起请求
	resp, err := client.Sheets.SpreadsheetSheet.Query(ctx, req)
	// 处理错误
	if err != nil {
		return "", 0, err
	}
	// 服务端错误处理
	if !resp.Success() {
		return "", 0, fmt.Errorf("request failed with code: %d, message: %s", resp.Code, resp.Msg)
		//return "", 0, fmt.Errorf("code:"+strconv.Itoa(resp.Code), "msg:"+resp.Msg, "request_id:"+resp.RequestId())
	}

	sheetID := resp.Data.Sheets[0].SheetId
	rowCount := resp.Data.Sheets[0].GridProperties.RowCount

	return *sheetID, *rowCount, nil
}
func (s *feiShuService) SaveTableData(ctx context.Context, appID string, appSecret string, spreadsheetToken string) error {
	lastModifyTime, err := s.GetSheetLatestModifyTime(ctx, appID, appSecret, spreadsheetToken)
	if err != nil {
		return err
	}

	updateLog := ""
	// 将时间戳字符串转换为 int64 类型的时间戳
	timestamp, err := strconv.ParseInt(lastModifyTime, 10, 64)
	if err == nil {
		// 将时间戳转换为 time.Time 类型
		t := time.Unix(timestamp, 0)

		// 格式化时间为 yyyy-mm-dd hh:mm:ss
		updateLog = t.Format("2006-01-02 15:04:05")
	} else {
		fmt.Println("Error parsing timestamp:", err)
	}

	err = s.tm.Transaction(ctx, func(ctx context.Context) error {
		if err = s.sheetInfoRepo.Create(ctx, &model.SheetInfo{
			SheetID:      spreadsheetToken,
			SheetName:    "",
			UpdateLog:    updateLog,
			RuntimeState: "",
		}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	sheetID, rowCount, err := s.GetFirstSheetIDAndRowCount(ctx, appID, appSecret, spreadsheetToken)
	if err != nil {
		return err
	}
	fmt.Println("sheetID ==>", sheetID)
	fmt.Println("rowCount ==>", rowCount)

	tenantAccessToken, err := s.GetTenantAccessToken(ctx, appID, appSecret)
	if err != nil {
		return err
	}
	fmt.Println("tenantAccessToken ==>", tenantAccessToken)

	var tableData []map[string]string

	// 生成请求URL，范围从B{startRow}到D{endRow}
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/sheets/v2/spreadsheets/%s/values/%s!B%d:D%d?dateTimeRenderOption=FormattedString", spreadsheetToken, sheetID, 2, rowCount)

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tenantAccessToken)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	// 解析 JSON 响应体
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Revision         int    `json:"revision"`
			SpreadsheetToken string `json:"spreadsheetToken"`
			ValueRange       struct {
				MajorDimension string          `json:"majorDimension"`
				Range          string          `json:"range"`
				Revision       int             `json:"revision"`
				Values         [][]interface{} `json:"values"`
			} `json:"valueRange"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}

	// 判断请求是否成功
	if response.Code == 0 && response.Msg == "success" {
		for sortNumber, row := range response.Data.ValueRange.Values {
			rowData := make(map[string]string)
			for index, cell := range row {
				switch v := cell.(type) {
				case []interface{}:
					if index == 0 && len(v) > 0 {
						text, textOk := v[0].(map[string]interface{})["text"].(string)
						link, linkOk := v[0].(map[string]interface{})["link"].(string)

						// 确保 text 和 link 不为空字符串
						if textOk && text != "" && linkOk && link != "" {
							rowData["text"] = text
							rowData["link"] = link
							rowData["sortNumber"] = strconv.Itoa(sortNumber)
						}
					}
				case string:
					if index == 2 {
						// 处理单元格内容为字符串的情况
						rowData["date"] = v
					}
				}
			}
			// 检查 text 和 link 是否都存在
			if _, textExists := rowData["text"]; textExists {
				if _, linkExists := rowData["link"]; linkExists {
					// 检查 rowData 是否包含 date，如果没有则默认设置为空字符串
					if _, dateExists := rowData["date"]; !dateExists {
						rowData["date"] = ""
					}
					tableData = append(tableData, rowData)
				}
			}
		}

		// 输出 tableData 的内容
		fmt.Println("tableData[0] ==>", tableData[0])
		fmt.Println("tableData[len(tableData)-1] ==>", tableData[len(tableData)-1])

		var wg sync.WaitGroup
		concurrencyLimit := 5
		sem := make(chan struct{}, concurrencyLimit)

		for _, rowData := range tableData {
			sem <- struct{}{} // 将空结构体放入通道以限制并发
			wg.Add(1)

			go func(rowData map[string]string) {
				defer wg.Done()
				defer func() { <-sem }() // 从通道中移除空结构体以释放资源

				// 调用 GetDocumentData 方法
				content, err := s.GetDocumentData(ctx, appID, appSecret, rowData)
				if err != nil {
					content = fmt.Sprintf("Error processing link %s: %v\n", rowData["link"], err)
				}
				fmt.Println("text ==>", rowData["text"])
				fmt.Println("link ==>", rowData["link"])
				fmt.Println("date ==>", rowData["date"])
				fmt.Println("sortNumber ==>", rowData["sortNumber"])

				sortNumber, err := strconv.Atoi(rowData["sortNumber"])
				if err != nil {
					fmt.Printf("error converting index to int: %v", err)
					sortNumber = 9999
				}

				err = s.tm.Transaction(ctx, func(ctx context.Context) error {
					if err = s.cellDataRepo.Create(ctx, &model.CellData{
						Title:        rowData["text"],
						Link:         rowData["link"],
						Content:      content,
						Abstract:     "",
						RuntimeState: "",
						SheetID:      spreadsheetToken,
						ReleaseDate:  rowData["date"],
						Keyword:      "",
						SortNumber:   sortNumber,
					}); err != nil {
						return err
					}
					return nil
				})
				if err != nil {
					fmt.Printf("cell_data write failed: Title='%s', Link='%s', ReleaseDate='%s', Error: %v\n",
						rowData["text"], rowData["link"], rowData["date"], err)
				}

				// 每次调用后 sleep 1 秒
				time.Sleep(1000 * time.Millisecond)
			}(rowData)
		}

		// 等待所有协程完成
		wg.Wait()

		return nil
	} else {
		return fmt.Errorf("request failed with code: %d, message: %s", response.Code, response.Msg)
	}
}

func (s *feiShuService) GetDocumentData(ctx context.Context, appID string, appSecret string, documentMetaData map[string]string) (string, error) {
	//text := documentMetaData["text"]
	link := documentMetaData["link"]
	documentID := strings.Split(link, "/")[len(strings.Split(link, "/"))-1]
	//date := documentMetaData["date"]

	tenantAccessToken, err := s.GetTenantAccessToken(ctx, appID, appSecret)
	if err != nil {
		return "", err
	}

	// 生成请求URL，范围从B{startRow}到D{endRow}
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/docx/v1/documents/%s/raw_content", documentID)

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tenantAccessToken)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// 解析 JSON 响应体
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Content string `json:"content"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error decoding JSON: %w", err)
	}

	// 判断请求是否成功
	if response.Code == 0 && response.Msg == "success" {
		return response.Data.Content, nil
	} else {
		return "", fmt.Errorf("request failed with code: %d, message: %s", response.Code, response.Msg)
	}
}

func (s *feiShuService) GetSheetLatestModifyTime(ctx context.Context, appID string, appSecret string, spreadsheetToken string) (string, error) {
	// 创建 Client
	client := lark.NewClient(appID, appSecret)
	// 创建请求对象
	req := larkdrive.NewBatchQueryMetaReqBuilder().
		MetaRequest(larkdrive.NewMetaRequestBuilder().
			RequestDocs([]*larkdrive.RequestDoc{
				larkdrive.NewRequestDocBuilder().
					DocToken(spreadsheetToken).
					DocType(`sheet`).
					Build(),
			}).
			WithUrl(false).
			Build()).
		Build()

	// 发起请求
	resp, err := client.Drive.Meta.BatchQuery(ctx, req)

	// 处理错误
	if err != nil {
		return "", err
	}

	// 服务端错误处理
	if !resp.Success() {
		return "", fmt.Errorf("request failed with code: %d, message: %s", resp.Code, resp.Msg)
	}

	return *resp.Data.Metas[0].LatestModifyTime, nil
}
