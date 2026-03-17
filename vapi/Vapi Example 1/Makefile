APP_NAME := vapi-example-1
CMD_DIR  := ./cmd
BUILD_DIR := ./bin

.PHONY: all build run clean tidy list create get delete help

## help: Bu yardım mesajını göster
help:
	@echo ""
	@echo "Vapi Voice AI - Go Demo"
	@echo "========================"
	@echo ""
	@echo "Kullanım:"
	@echo "  make build          Projeyi derle"
	@echo "  make run            Yardım mesajını göster (build + run)"
	@echo "  make tidy           go mod tidy çalıştır"
	@echo "  make clean          Build çıktılarını temizle"
	@echo ""
	@echo "  make list           Tüm asistanları listele"
	@echo "  make create         Yeni asistan oluştur (config.yaml'dan)"
	@echo "  make get ID=<id>    Asistan detayını getir"
	@echo "  make delete ID=<id> Asistanı sil"
	@echo ""
	@echo "Örnek:"
	@echo "  make create"
	@echo "  make list"
	@echo "  make get ID=abc-123-def"
	@echo "  make delete ID=abc-123-def"
	@echo ""

## build: Projeyi derle
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)
	@echo "Build tamamlandı: $(BUILD_DIR)/$(APP_NAME)"

## run: Build edip yardım mesajını göster
run: build
	$(BUILD_DIR)/$(APP_NAME)

## tidy: Bağımlılıkları düzenle
tidy:
	go mod tidy

## clean: Build çıktılarını sil
clean:
	rm -rf $(BUILD_DIR)
	@echo "Temizlendi."

## list: Tüm asistanları listele
list: build
	$(BUILD_DIR)/$(APP_NAME) list

## create: config.yaml'daki ayarlarla yeni asistan oluştur
create: build
	$(BUILD_DIR)/$(APP_NAME) create

## get: Asistan detayını getir (make get ID=xxx)
get: build
ifndef ID
	@echo "Hata: ID parametresi gerekli."
	@echo "Kullanım: make get ID=<assistant_id>"
	@exit 1
endif
	$(BUILD_DIR)/$(APP_NAME) get $(ID)

## delete: Asistanı sil (make delete ID=xxx)
delete: build
ifndef ID
	@echo "Hata: ID parametresi gerekli."
	@echo "Kullanım: make delete ID=<assistant_id>"
	@exit 1
endif
	$(BUILD_DIR)/$(APP_NAME) delete $(ID)

all: tidy build
