APP_NAME := GreeterMCP
PORT := 50000
NGROK_DOMAIN := mcp.stayos.ai

.PHONY: build run ngrok all stop clean

# Derle
build:
	go build -o $(APP_NAME) .

# Derle ve çalıştır
run: build
	./$(APP_NAME)

# ngrok tünelini aç
ngrok:
	ngrok http --url=$(NGROK_DOMAIN) --host-header=localhost:$(PORT) $(PORT)

# Derle, arka planda çalıştır, ngrok aç
all: build
	@lsof -ti :$(PORT) | xargs kill 2>/dev/null || true
	@echo "🚀 Server başlatılıyor..."
	./$(APP_NAME) & echo $$! > .pid
	@sleep 2
	@echo "📡 ngrok başlatılıyor..."
	ngrok http --url=$(NGROK_DOMAIN) --host-header=localhost:$(PORT) $(PORT)

# Arka plandaki server'ı durdur
stop:
	@if [ -f .pid ]; then \
		kill $$(cat .pid) 2>/dev/null && echo "Server durduruldu." || echo "Server zaten durmuş."; \
		rm -f .pid; \
	else \
		echo "PID dosyası bulunamadı."; \
	fi

# Temizle
clean: stop
	rm -f $(APP_NAME) .pid
