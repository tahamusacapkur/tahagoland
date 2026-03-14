/*
	ICI Tech Teknoloji A.Ş.
	User        : ICI
	Name        : Ibrahim COBANI
	Date        : 14.03.2026
	Time        : 08:00
	Notes       : MCP Server over HTTP with Gin + ngrok for ChatGPT integration

	Mimari:
	┌──────────┐    HTTPS     ┌────────┐    HTTP     ┌──────────────┐
	│ ChatGPT  │ ──────────►  │ ngrok  │ ─────────►  │ Gin + MCP    │
	│          │              │        │             │ localhost:50000│
	└──────────┘              └────────┘             └──────────────┘

	Endpoints:
	- /test           → Test MCP Server
	- /reservation    → Reservation MCP Server
	- /hotel-content  → Hotel Content MCP Server
	- /upsell         → Upsell MCP Server
*/

package main

import (
	"log"

	"MCPExample2/router"
)

func main() {
	r := router.Setup()

	log.Println("🚀 MCP Server başlatılıyor: http://localhost:50000")
	log.Println("📡 Endpoints: /test, /reservation, /hotel-content, /upsell")

	if err := r.Run(":50000"); err != nil {
		log.Fatal(err)
	}
}
