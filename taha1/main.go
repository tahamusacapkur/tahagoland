package main

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/kenshaw/escpos"
)

func main() {
	var buf bytes.Buffer
	p := escpos.New(&buf)

	// Read 1400.pj file content
	content, err := os.ReadFile("1400.pj")
	if err != nil {
		panic(err)
	}

	p.Init()
	p.Write(string(content))
	p.Formfeed()
	p.Cut()
	p.End()

	// Yazıcı adı tam olarak bu olmalı (kontrol et: "lpstat -p")
	cmd := exec.Command("lp", "-d", "ZHU_HAI_SUNCSW_Receipt_Printer_Co__Ltd__Gprinter_GP_L80160", "-o", "raw")

	// ESC/POS komutlarını stdin'den geçir
	cmd.Stdin = &buf

	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
