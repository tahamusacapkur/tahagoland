package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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
