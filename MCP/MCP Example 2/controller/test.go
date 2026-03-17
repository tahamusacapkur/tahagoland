package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"MCPExample2/service"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/xuri/excelize/v2"
)

const (
	ip2locationAPIKey = "F1055F5E6C8BA23789CA746A3D605785"
	excelFilePath     = "data/Aylık Hijyen ve Gıda Güvenliği Denetim Puanı.xlsx"
	sheetName         = "Rapor"
)

var monthColumns = map[string]int{
	"ocak": 4, "şubat": 5, "mart": 6, "nisan": 7,
	"mayıs": 8, "haziran": 9, "temmuz": 10, "ağustos": 11,
	"eylül": 12, "ekim": 13, "kasım": 14, "aralık": 15,
}

type IPLookupInput struct {
	IP string `json:"ip" jsonschema:"required,the IP address to look up. Ask the user for their IP address if not provided"`
}

type ipLocationResponse struct {
	IP          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionName  string  `json:"region_name"`
	CityName    string  `json:"city_name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ZipCode     string  `json:"zip_code"`
	TimeZone    string  `json:"time_zone"`
	ASN         string  `json:"asn"`
	AS          string  `json:"as"`
	IsProxy     bool    `json:"is_proxy"`
}

type EmptyInput struct{}

type AuditScoreInput struct {
	Department string `json:"department" jsonschema:"required,the main department name to search for, e.g. IST-İÇ OPERASYON"`
	Month      string `json:"month" jsonschema:"required,the month name in Turkish, e.g. Ocak, Şubat, Mart"`
}

// --- Room Price MCP Input (with jsonschema tags for MCP tool) ---

type RoomPriceInput struct {
	CheckInDate  string `json:"check_in_date" jsonschema:"required,check-in date in YYYY-MM-DD format"`
	CheckOutDate string `json:"check_out_date" jsonschema:"required,check-out date in YYYY-MM-DD format"`
	Adults       int    `json:"adults" jsonschema:"required,number of adult guests"`
	Children     int    `json:"children" jsonschema:"required,number of child guests (0 if none)"`
	ChildAges    []int  `json:"child_ages" jsonschema:"ages of each child guest. Required if children > 0. One age per child e.g. [3, 7] for 2 children aged 3 and 7"`
}

func NewTestServer() *mcp.Server {
	server := mcp.NewServer(
		&mcp.Implementation{Name: "test", Version: "v1.0.0"},
		nil,
	)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "greet",
		Description: "Test: Bir kişiye merhaba der",
	}, testGreet)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "ip-location",
		Description: "ALWAYS use this tool when the user asks about location, 'where am I', 'ben nerdeyim', or anything related to IP geolocation. Looks up the geographic location of an IP address. Ask the user for their IP if not provided.",
	}, ipLocation)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list-departments",
		Description: "Ünite (Ana Departman) ve Departman listesini getirir. 'Departmanları listele', 'Ana departmanlar neler', 'Hangi üniteler var' gibi sorularda kullan.",
	}, listDepartments)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "audit-score",
		Description: "Hijyen ve Gıda Güvenliği denetim puanlarını getirir. Departman adı ve ay bilgisi gereklidir. Örnek: 'IST-İÇ OPERASYON Ocak ayı puanları', 'AHL-ÜRETİM Mart denetim sonuçları'. Departmanı ve alt departmanlarının puanlarını döner.",
	}, auditScore)
	mcp.AddTool(server, &mcp.Tool{
		Name: "room-price",
		Description: `Otel oda fiyatlarını sorgular. Müşterinin giriş tarihi (check_in_date), çıkış tarihi (check_out_date), yetişkin sayısı (adults) ve çocuk sayısı (children) bilgilerini alır. Müsait odaları fiyatlarıyla birlikte listeler. Örnek: '1 Ekim - 3 Ekim arası 2 yetişkin için oda fiyatları', 'Yarın için uygun odalar'. Tarihler YYYY-MM-DD formatında olmalıdır.

CRITICAL: Bu tool'un döndüğü çıktı önceden formatlanmış zengin Markdown içerir (görseller, tablolar, linkler, butonlar). Çıktıyı KESİNLİKLE olduğu gibi kullanıcıya göster. Özetleme, sadeleştirme, yeniden formatlama yapma. Markdown görsellerini (![...](...)) kaldırma. Linkleri ([...](...)) kaldırma. Tabloları değiştirme. Çıktıyı birebir kullanıcıya ilet, hiçbir şeyi atlama veya değiştirme.`,
	}, roomPrice)
	return server
}

func testGreet(ctx context.Context, req *mcp.CallToolRequest, input GreetInput) (*mcp.CallToolResult, GreetOutput, error) {
	return nil, GreetOutput{Greeting: "🧪 [Test] Merhaba " + input.Name + "!"}, nil
}

func ipLocation(ctx context.Context, req *mcp.CallToolRequest, input IPLookupInput) (*mcp.CallToolResult, GreetOutput, error) {
	url := fmt.Sprintf("https://api.ip2location.io/?key=%s", ip2locationAPIKey)
	if input.IP != "" {
		url += "&ip=" + input.IP
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, GreetOutput{Greeting: "IP sorgusu yapılamadı: " + err.Error()}, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, GreetOutput{Greeting: "Yanıt okunamadı: " + err.Error()}, nil
	}

	var loc ipLocationResponse
	if err := json.Unmarshal(body, &loc); err != nil {
		return nil, GreetOutput{Greeting: "Yanıt ayrıştırılamadı: " + err.Error()}, nil
	}

	result := fmt.Sprintf(
		"📍 IP: %s\n🌍 Ülke: %s (%s)\n🏙️ Şehir: %s, %s\n🗺️ Koordinat: %.5f, %.5f\n📮 Posta Kodu: %s\n🕐 Saat Dilimi: UTC%s\n🌐 ISP: %s (ASN: %s)\n🛡️ Proxy: %v",
		loc.IP, loc.CountryName, loc.CountryCode,
		loc.CityName, loc.RegionName,
		loc.Latitude, loc.Longitude,
		loc.ZipCode, loc.TimeZone,
		loc.AS, loc.ASN, loc.IsProxy,
	)

	return nil, GreetOutput{Greeting: result}, nil
}

func listDepartments(ctx context.Context, req *mcp.CallToolRequest, input EmptyInput) (*mcp.CallToolResult, GreetOutput, error) {
	f, err := excelize.OpenFile(excelFilePath)
	if err != nil {
		return nil, GreetOutput{Greeting: "Excel dosyası açılamadı: " + err.Error()}, nil
	}
	defer f.Close()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, GreetOutput{Greeting: "Rapor sayfası okunamadı: " + err.Error()}, nil
	}

	var sb strings.Builder
	sb.WriteString("📋 Ünite ve Departman Listesi\n")
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	for i := 8; i < len(rows); i++ {
		if len(rows[i]) == 0 {
			continue
		}
		cellVal := rows[i][0]
		trimmed := strings.TrimSpace(cellVal)
		if trimmed == "" {
			continue
		}

		// Indentli satırlar alt departman — atla
		isIndented := cellVal[0] == ' ' || cellVal[0] == '\t'
		if isIndented {
			continue
		}

		// Hücre rengini kontrol et — mavi (2E5090) = Ünite (Ana Departman)
		cellRef, _ := excelize.CoordinatesToCellName(1, i+1) // excelize 1-indexed
		styleID, _ := f.GetCellStyle(sheetName, cellRef)
		isUnit := false
		if style, err := f.GetStyle(styleID); err == nil && style != nil && style.Fill.Color != nil {
			for _, color := range style.Fill.Color {
				c := strings.ToUpper(strings.TrimPrefix(strings.TrimPrefix(color, "FF"), "#"))
				if c == "2E5090" {
					isUnit = true
					break
				}
			}
		}

		if isUnit {
			sb.WriteString(fmt.Sprintf("\n🏢 %s\n", trimmed))
		} else {
			sb.WriteString(fmt.Sprintf("  📌 %s\n", trimmed))
		}
	}

	return nil, GreetOutput{Greeting: sb.String()}, nil
}

func auditScore(ctx context.Context, req *mcp.CallToolRequest, input AuditScoreInput) (*mcp.CallToolResult, GreetOutput, error) {
	col, ok := monthColumns[strings.ToLower(input.Month)]
	if !ok {
		return nil, GreetOutput{Greeting: "Geçersiz ay: " + input.Month + ". Türkçe ay adı girin (Ocak, Şubat, Mart, ...)"}, nil
	}

	f, err := excelize.OpenFile(excelFilePath)
	if err != nil {
		return nil, GreetOutput{Greeting: "Excel dosyası açılamadı: " + err.Error()}, nil
	}
	defer f.Close()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, GreetOutput{Greeting: "Rapor sayfası okunamadı: " + err.Error()}, nil
	}

	searchDept := strings.ToUpper(strings.TrimSpace(input.Department))
	mainRow := -1

	// Ana departmanı bul (indent olmayan satır)
	for i := 8; i < len(rows); i++ {
		if len(rows[i]) == 0 {
			continue
		}
		cellVal := rows[i][0]
		trimmed := strings.TrimSpace(cellVal)
		if trimmed == "" {
			continue
		}
		isIndented := len(cellVal) > 0 && (cellVal[0] == ' ' || cellVal[0] == '\t')
		if !isIndented && strings.ToUpper(trimmed) == searchDept {
			mainRow = i
			break
		}
	}

	if mainRow == -1 {
		return nil, GreetOutput{Greeting: "'" + input.Department + "' departmanı bulunamadı."}, nil
	}

	// Ana departman puanı
	mainScore := getCellValue(rows, mainRow, col)
	avgScore := getCellValue(rows, mainRow, 16) // Q kolonu = Yıl Ort.

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📊 %s — %s Ayı Denetim Puanları\n", strings.TrimSpace(rows[mainRow][0]), input.Month))
	sb.WriteString(fmt.Sprintf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
	sb.WriteString(fmt.Sprintf("🏢 Genel Puan: %s (Yıl Ort: %s)\n\n", mainScore, avgScore))
	sb.WriteString("📋 Alt Departmanlar:\n")

	// Alt departmanları topla (sonraki indentli satırlar)
	hasSubDepts := false
	for i := mainRow + 1; i < len(rows); i++ {
		if len(rows[i]) == 0 {
			continue
		}
		cellVal := rows[i][0]
		trimmed := strings.TrimSpace(cellVal)
		if trimmed == "" {
			continue
		}
		isIndented := len(cellVal) > 0 && (cellVal[0] == ' ' || cellVal[0] == '\t')
		if !isIndented {
			break // Yeni ana departmana geçtik
		}
		subScore := getCellValue(rows, i, col)
		subAvg := getCellValue(rows, i, 16)
		sb.WriteString(fmt.Sprintf("  ▸ %-30s %s (Yıl Ort: %s)\n", trimmed, subScore, subAvg))
		hasSubDepts = true
	}

	if !hasSubDepts {
		sb.WriteString("  (Alt departman bulunamadı)\n")
	}

	return nil, GreetOutput{Greeting: sb.String()}, nil
}

func roomPrice(ctx context.Context, req *mcp.CallToolRequest, input RoomPriceInput) (*mcp.CallToolResult, GreetOutput, error) {
	// Service'e delege et
	svcInput := service.RoomPriceInput{
		CheckInDate:  input.CheckInDate,
		CheckOutDate: input.CheckOutDate,
		Adults:       input.Adults,
		Children:     input.Children,
		ChildAges:    input.ChildAges,
	}

	priceResp, err := service.GetRoomPrice(ctx, svcInput)
	if err != nil {
		return nil, GreetOutput{Greeting: err.Error()}, nil
	}

	if len(priceResp.Rooms) == 0 {
		return nil, GreetOutput{Greeting: "Seçtiğiniz tarihler için müsait oda bulunamadı. Farklı tarihler deneyebilirsiniz."}, nil
	}

	hotel := priceResp.Hotel
	var contents []mcp.Content

	// Helper: görseli indir ve ImageContent (base64) olarak ekle
	addImage := func(imgURL string) {
		if imgURL == "" || strings.HasSuffix(strings.ToLower(imgURL), ".mp4") {
			return
		}
		imgResp, err := http.Get(imgURL)
		if err != nil || imgResp.StatusCode != http.StatusOK {
			return
		}
		defer imgResp.Body.Close()
		data, err := io.ReadAll(imgResp.Body)
		if err != nil || len(data) == 0 {
			return
		}
		mimeType := "image/jpeg"
		ext := strings.ToLower(filepath.Ext(imgURL))
		switch ext {
		case ".png":
			mimeType = "image/png"
		case ".webp":
			mimeType = "image/webp"
		case ".gif":
			mimeType = "image/gif"
		}
		contents = append(contents, &mcp.ImageContent{Data: data, MIMEType: mimeType})
	}

	addText := func(text string) {
		contents = append(contents, &mcp.TextContent{Text: text})
	}

	// ===== OTEL HEADER =====
	addText(fmt.Sprintf("# %s\n\n%s\n\n", hotel.Name, hotel.Address))

	// ===== REZERVASYON ÖZETİ =====
	var sb strings.Builder
	sb.WriteString("---\n\n## Arama Detayları\n\n")
	sb.WriteString("| | |\n|---|---|\n")
	sb.WriteString(fmt.Sprintf("| **Giriş Tarihi** | %s |\n", priceResp.CheckInDate))
	sb.WriteString(fmt.Sprintf("| **Çıkış Tarihi** | %s |\n", priceResp.CheckOutDate))
	sb.WriteString(fmt.Sprintf("| **Konaklama** | %d Gece |\n", priceResp.StayDays))
	sb.WriteString(fmt.Sprintf("| **Misafirler** | %s |\n", priceResp.GuestCountSummary))
	sb.WriteString(fmt.Sprintf("| **Pansiyon Tipi** | %s |\n", priceResp.BoardType))
	sb.WriteString(fmt.Sprintf("| **Fiyat Geçerliliği** | %s |\n", priceResp.ValidityDate))
	sb.WriteString("\n---\n\n## Müsait Odalar\n\n")
	addText(sb.String())

	// ===== ODALAR =====
	for i, room := range priceResp.Rooms {
		addText(fmt.Sprintf("### %d. %s\n", i+1, room.RoomType.Name))

		// Sadece ana oda görseli
		addImage(room.RoomType.ImageURL)

		var roomSB strings.Builder
		desc := service.StripHTMLTags(room.RoomType.Description)
		if desc != "" {
			roomSB.WriteString(fmt.Sprintf("\n> %s\n\n", desc))
		}
		roomSB.WriteString(fmt.Sprintf("**Kapasite:** %d yetişkin, maksimum %d kişi\n\n", room.RoomType.AdultCount, room.RoomType.RoomCapacity))

		for j, offer := range room.BookingOffers {
			contractDesc := service.ExtractLocaleText(offer.ContractDescription, "tr")
			roomSB.WriteString(fmt.Sprintf("#### Teklif %d: %s\n\n", j+1, offer.ContractName))
			if contractDesc != "" {
				roomSB.WriteString(fmt.Sprintf("> %s\n\n", contractDesc))
			}
			roomSB.WriteString("| | |\n|---|---|\n")
			roomSB.WriteString(fmt.Sprintf("| Gecelik Fiyat | **%.2f %s** |\n", offer.AvgNightPrice, room.CurrencyCode))
			roomSB.WriteString(fmt.Sprintf("| %d Gece Toplam | %.2f %s |\n", offer.Night, offer.TotalAmount, room.CurrencyCode))
			for _, t := range offer.Taxes {
				roomSB.WriteString(fmt.Sprintf("| %s | %.2f %s |\n", t.Name, t.Amount, room.CurrencyCode))
			}
			roomSB.WriteString(fmt.Sprintf("| **Vergiler Dahil Toplam** | **%.2f %s** |\n", offer.NetAmount, room.CurrencyCode))
			roomSB.WriteString("\n")
			roomSB.WriteString(fmt.Sprintf("👉 **[SATIN AL — %s — %.2f %s](https://google.com)**\n\n", room.RoomType.Name, offer.NetAmount, room.CurrencyCode))
		}
		roomSB.WriteString("---\n\n")
		addText(roomSB.String())
	}

	// ===== İLETİŞİM =====
	var contactSB strings.Builder
	contactSB.WriteString("## Otel Bilgileri & İletişim\n\n| | |\n|---|---|\n")
	if hotel.CallCenter != "" {
		contactSB.WriteString(fmt.Sprintf("| **Çağrı Merkezi** | %s |\n", hotel.CallCenter))
	}
	if hotel.SpReception != "" {
		contactSB.WriteString(fmt.Sprintf("| **Resepsiyon** | %s |\n", hotel.SpReception))
	}
	if hotel.WhatsappNo != "" {
		contactSB.WriteString(fmt.Sprintf("| **WhatsApp** | %s |\n", hotel.WhatsappNo))
	}
	if hotel.WebSiteURL != "" {
		contactSB.WriteString(fmt.Sprintf("| **Web Sitesi** | %s |\n", hotel.WebSiteURL))
	}
	contactSB.WriteString("\n*Hayalinizdeki tatil sizi bekliyor! Hemen rezervasyonunuzu oluşturabilirsiniz.*\n")
	addText(contactSB.String())

	return &mcp.CallToolResult{Content: contents}, GreetOutput{}, nil
}

func getCellValue(rows [][]string, row, col int) string {
	if row >= len(rows) || col >= len(rows[row]) {
		return "-"
	}
	val := strings.TrimSpace(rows[row][col])
	if val == "" {
		return "-"
	}
	// Uzun ondalık sayıları kısalt
	if strings.Contains(val, ".") {
		parts := strings.Split(val, ".")
		if len(parts[1]) > 2 {
			return parts[0] + "." + parts[1][:2]
		}
	}
	return val
}
