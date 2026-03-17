# Vapi Voice AI - Go Demo

[Vapi.ai](https://vapi.ai) REST API ile sesli yapay zeka asistanları oluşturmak, listelemek, detay görüntülemek ve
silmek için basit bir Go CLI uygulaması.

## Proje Yapısı

```
Vapi Example 1/
├── cmd/
│   └── main.go              # Entry point — CLI komutları
├── internal/
│   ├── client/
│   │   └── client.go        # VapiClient — HTTP istekleri
│   ├── config/
│   │   └── config.go        # Viper ile config.yaml okuma
│   └── model/
│       └── model.go         # Struct tanımları (Assistant, Model, Voice...)
├── config.yaml              # Senin API key'in (git'e COMMIT ETME)
├── config.yaml.example      # Template — paylaşılabilir
├── go.mod
├── go.sum
├── Makefile
├── .gitignore
└── README.md
```

## Ön Gereksinimler

- Go 1.22+
- Vapi.ai hesabı → [dashboard.vapi.ai](https://dashboard.vapi.ai)

## Kurulum

1. **Bağımlılıkları indir:**

```bash
make tidy
```

2. **Config dosyasını hazırla:**

```bash
cp config.yaml.example config.yaml
```

3. **API Key'ini ekle:**

[Vapi Dashboard](https://dashboard.vapi.ai) → Organization Settings → API Keys → **Private Key**'i kopyala.

`config.yaml` dosyasını aç ve `api_key` alanına yapıştır:

```yaml
vapi:
  api_key: "BURAYA_PRIVATE_KEY"
  base_url: "https://api.vapi.ai"
```

> Alternatif olarak environment variable da kullanabilirsin:
> ```bash
> export VAPI_API_KEY="senin-key-in"
> ```

## Kullanım

### Makefile ile (önerilen)

```bash
make help               # Tüm komutları göster
make list               # Asistanları listele
make create             # Yeni asistan oluştur (config.yaml'dan)
make get ID=abc-123     # Asistan detayı
make delete ID=abc-123  # Asistanı sil
```

### Doğrudan Go ile

```bash
go run ./cmd list
go run ./cmd create
go run ./cmd get <assistant_id>
go run ./cmd delete <assistant_id>
```

## Nasıl Test Ederim?

### Adım 1 — API key'ini doğrula

```bash
make list
```

Eğer `[]` veya asistan listesi dönüyorsa key'in doğru demektir. Hata alırsan key'i kontrol et.

### Adım 2 — Asistan oluştur

```bash
make create
```

Çıktıda `ID`, `Name` ve `Created` bilgisi gelecek. Bu ID'yi not al.

### Adım 3 — Dashboard'dan kontrol et

[dashboard.vapi.ai](https://dashboard.vapi.ai) → Assistants sayfasına git. Oluşturduğun asistanı orada göreceksin.

### Adım 4 — Asistanı test et

Dashboard'da asistanın yanındaki **"Talk"** butonuna tıkla. Mikrofon izni ver ve konuşmaya başla. Asistan sana
`config.yaml`'daki `first_message` ile karşılık verecek.

### Adım 5 — Temizlik

```bash
make delete ID=<yukarıdaki-id>
```

## Config Açıklaması

| Alan                             | Açıklama                                                 |
|----------------------------------|----------------------------------------------------------|
| `vapi.api_key`                   | Dashboard'dan aldığın Private API Key                    |
| `vapi.base_url`                  | Vapi API adresi (`https://api.vapi.ai`)                  |
| `assistant.name`                 | Asistanın adı                                            |
| `assistant.first_message`        | Asistanın söyleyeceği ilk mesaj                          |
| `assistant.model.provider`       | LLM sağlayıcı (`openai`, `anthropic`)                    |
| `assistant.model.model`          | Model adı (`gpt-4o-mini`, `gpt-4o`, `claude-3-5-sonnet`) |
| `assistant.model.system_prompt`  | Asistanın kişiliğini belirleyen system prompt            |
| `assistant.voice.provider`       | TTS sağlayıcı (`11labs`, `deepgram`, `playht`)           |
| `assistant.voice.voice_id`       | Ses ID'si (11labs dashboard'dan alınır)                  |
| `assistant.transcriber.provider` | STT sağlayıcı (`deepgram`)                               |
| `assistant.transcriber.language` | Konuşma dili (`tr`, `en`)                                |

## Sıkça Sorulan Sorular

**OpenAI API key'i lazım mı?**
Hayır. Vapi kendi LLM bağlantısını yönetiyor, dakika başı ücretlendirmede dahil. İstersen dashboard'da Provider Keys
bölümünden kendi key'ini ekleyebilirsin — opsiyonel.

**Public Key mi Private Key mi?**
Go backend yazdığın için **Private Key** (Server-side API access). Public Key tarayıcı/mobil SDK entegrasyonları
içindir.

**Ücretlendirme nasıl?**
Platform: $0.05/dk + LLM + TTS + STT maliyetleri. Toplamda ~$0.15–0.36/dk arası.
