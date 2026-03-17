package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// RoomPriceInput is the shared input for room price queries.
type RoomPriceInput struct {
	CheckInDate  string `json:"check_in_date"`
	CheckOutDate string `json:"check_out_date"`
	Adults       int    `json:"adults"`
	Children     int    `json:"children"`
	ChildAges    []int  `json:"child_ages"`
}

// Validate checks required fields.
func (r *RoomPriceInput) Validate() error {
	if r.CheckInDate == "" {
		return fmt.Errorf("check_in_date zorunludur")
	}
	if r.CheckOutDate == "" {
		return fmt.Errorf("check_out_date zorunludur")
	}
	if r.Adults < 1 {
		return fmt.Errorf("en az 1 yetişkin olmalıdır")
	}
	if r.Children > 0 && len(r.ChildAges) != r.Children {
		return fmt.Errorf("%d çocuk için %d yaş belirtildi, her çocuk için yaş gereklidir", r.Children, len(r.ChildAges))
	}
	return nil
}

// --- API request/response types ---

type roomPriceRequest struct {
	LeadReservation   map[string]interface{} `json:"lead_reservation"`
	CheckInDate       string                 `json:"check_in_date"`
	CheckOutDate      string                 `json:"check_out_date"`
	RoomDistributions []roomDistribution     `json:"room_distributions"`
}

type roomDistribution map[string]interface{}

type RoomPriceResponse struct {
	CheckInDate       string      `json:"check_in_date"`
	CheckOutDate      string      `json:"check_out_date"`
	StayDays          int         `json:"stay_days"`
	CurrencyCode      string      `json:"currency_code"`
	Adult             int         `json:"adult"`
	TotalGuest        int         `json:"total_guest"`
	GuestCountSummary string      `json:"guest_count_summary"`
	BoardType         string      `json:"board_type"`
	ValidityDate      string      `json:"validity_date"`
	Hotel             HotelInfo   `json:"hotel"`
	Rooms             []RoomEntry `json:"rooms"`
}

type HotelInfo struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	CallCenter      string   `json:"call_center"`
	MasterColor     string   `json:"master_color"`
	SubColor        string   `json:"sub_color"`
	Address         string   `json:"address"`
	City            string   `json:"city"`
	Country         string   `json:"country"`
	WebSiteURL      string   `json:"web_site_url"`
	LatLng          string   `json:"lat_lng"`
	WhatsappNo      string   `json:"whatsapp_no"`
	SpReception     string   `json:"sp_reception"`
	SpRoomService   string   `json:"sp_room_service"`
	SpSpa           string   `json:"sp_spa"`
	WifiName        string   `json:"wifi_name"`
	WifiPassword    string   `json:"wifi_password"`
	GoogleMapAddr   string   `json:"google_map_address"`
	SliderImages    []string `json:"slider_images"`
	Logo            []string `json:"logo"`
	DarkLogo        []string `json:"dark_logo"`
	LightLogo       []string `json:"light_logo"`
	AppStoreLink    string   `json:"app_store_link"`
	GooglePlayLink  string   `json:"google_play_link"`
	SmInstagramURL  string   `json:"sm_instagram_url"`
	SmFacebookURL   string   `json:"sm_facebook_url"`
	SmTwitterURL    string   `json:"sm_twitter_url"`
	SmYoutubeURL    string   `json:"sm_youtube_url"`
	TripAdvisorLink string   `json:"trip_advisor_link"`
}

type RoomEntry struct {
	Quantity      int            `json:"quantity"`
	RoomType      RoomType       `json:"room_type"`
	BookingOffers []BookingOffer `json:"booking_offers"`
	CurrencyCode  string         `json:"currency_code"`
}

type RoomType struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	ImageURL     string      `json:"image_url"`
	AdultCount   int         `json:"adult_count"`
	RoomCapacity int         `json:"room_capacity"`
	Images       []RoomImage `json:"images"`
}

type RoomImage struct {
	URL string `json:"url"`
}

type BookingOffer struct {
	ContractID          int     `json:"contract_id"`
	ContractName        string  `json:"contract_name"`
	ContractDescription string  `json:"contract_description"`
	AvgRoomAmount       float64 `json:"avg_room_amount"`
	TotalAmount         float64 `json:"total_amount"`
	Taxes               []Tax   `json:"taxes"`
	TotalTaxAmount      float64 `json:"total_tax_amount"`
	NetAmount           float64 `json:"net_amount"`
	AvgNightPrice       float64 `json:"avg_night_price"`
	Night               int     `json:"night"`
}

type Tax struct {
	Name   string  `json:"name"`
	Rate   float64 `json:"rate"`
	Amount float64 `json:"amount"`
}

// GetRoomPrice calls the icibot API and returns the parsed response.
func GetRoomPrice(ctx context.Context, input RoomPriceInput) (*RoomPriceResponse, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	dist := roomDistribution{
		"adult": input.Adults,
		"child": input.Children,
	}
	for i, age := range input.ChildAges {
		dist[fmt.Sprintf("child_age_%d", i+1)] = age
	}

	reqBody := roomPriceRequest{
		LeadReservation:   map[string]interface{}{},
		CheckInDate:       input.CheckInDate,
		CheckOutDate:      input.CheckOutDate,
		RoomDistributions: []roomDistribution{dist},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("istek oluşturulamadı: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.icibot.net/exapi/room_sales/3/ask_for_web_room_price",
		strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("HTTP isteği oluşturulamadı: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API isteği başarısız: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("yanıt okunamadı: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API hatası (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var priceResp RoomPriceResponse
	if err := json.Unmarshal(body, &priceResp); err != nil {
		return nil, fmt.Errorf("yanıt ayrıştırılamadı: %w", err)
	}

	return &priceResp, nil
}

// StripHTMLTags removes HTML tags and decodes common entities.
func StripHTMLTags(s string) string {
	var result strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(r)
		}
	}
	out := result.String()
	out = strings.ReplaceAll(out, "&#39;", "'")
	out = strings.ReplaceAll(out, "&amp;", "&")
	out = strings.ReplaceAll(out, "&quot;", "\"")
	out = strings.ReplaceAll(out, "&lt;", "<")
	out = strings.ReplaceAll(out, "&gt;", ">")
	return strings.TrimSpace(out)
}

// ExtractLocaleText extracts a locale's text from a JSON object string.
func ExtractLocaleText(jsonStr string, locale string) string {
	if jsonStr == "" || jsonStr[0] != '{' {
		return jsonStr
	}
	var localeMap map[string]string
	if err := json.Unmarshal([]byte(jsonStr), &localeMap); err != nil {
		return ""
	}
	if text, ok := localeMap[locale]; ok && text != "" {
		return text
	}
	return ""
}
