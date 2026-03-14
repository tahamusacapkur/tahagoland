# Chromem Sample 1

Vektör veritabanı ve arama özelliği örneği.

## Kurulum

### .env Dosyası

Bu proje OpenAI API key'e ihtiyaç duyar. Projeyi çalıştırmadan önce bu dizinde bir `.env` dosyası oluşturmanız gerekir.

1. Bu dizinde `.env` adında bir dosya oluşturun:

```bash
touch .env
```

2. Dosyanın içine OpenAI API key'inizi ekleyin:

```
OPENAI_API_KEY : "your-api-key-here"
```

3. API key'inizi [OpenAI Platform](https://platform.openai.com/api-keys) adresinden alabilirsiniz.

> **Not:** `.env` dosyası `.gitignore` ile git takibinden hariç tutulmuştur. API key'inizi asla commit etmeyin.
