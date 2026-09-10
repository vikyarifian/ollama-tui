package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
)

// Preset represents approved compliance templates matching internal requirements
type Preset struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	SystemPrompt string `json:"system_prompt"`
}

// Config holds the application level configuration including Ollama settings and Presets.
type Config struct {
	OllamaURL    string   `json:"ollama_url"`
	DefaultModel string   `json:"default_model"`
	Presets      []Preset `json:"presets"`
}

// GetDefaultConfig returns the hardcoded fallback configuration matching compliance templates
func GetDefaultConfig() *Config {
	return &Config{
		OllamaURL:    "http://localhost:11434",
		DefaultModel: "llama3",
		Presets: []Preset{
			{
				Name:        "Format Nota Dinas (Compliance)",
				Description: "Ubah draft kasar menjadi Nota Dinas resmi sesuai standard kepatuhan internal.",
				SystemPrompt: "Anda adalah asisten sekretariat korporat. Tugas Anda adalah menyusun Nota Dinas resmi bahasa Indonesia yang baku, sopan, dan terstruktur berdasarkan poin-poin yang diberikan. Patuhi aturan kerahasiaan data internal dan jangan sebutkan data nasabah spesifik.",
			},
			{
				Name:        "Temuan Audit Internal",
				Description: "Bantu menyusun draf temuan audit (Kondisi, Kriteria, Penyebab, Akibat, Rekomendasi).",
				SystemPrompt: "Anda adalah auditor internal senior. Bantu user merumuskan temuan audit menggunakan kerangka formal: Kondisi, Kriteria, Akibat, Penyebab, dan Rekomendasi. Gunakan bahasa Indonesia profesional dan objektif tanpa kesan menuduh.",
			},
			{
				Name:        "Filter Kepatuhan Data PII",
				Description: "Periksa dokumen untuk memastikan tidak ada data rahasia/PII bocor.",
				SystemPrompt: "Anda adalah kepatuhan (compliance) officer. Periksa teks berikut dan tunjukkan jika ada Personally Identifiable Information (PII) seperti NIK, nomor rekening, nama lengkap nasabah, atau nomor HP yang harus disamarkan sebelum diundgah ke publik.",
			},
			{
				Name:        "Penerjemah Dokumen Formal",
				Description: "Menerjemahkan instruksi atau draft memo dari Bahasa Inggris ke Bahasa Indonesia formal.",
				SystemPrompt: "Anda adalah penerjemah resmi korporat. Terjemahkan teks Bahasa Inggris berikut ke Bahasa Indonesia yang formal (sesuai EYD/PUEBI) dengan gaya penulisan korespondensi bisnis di Indonesia.",
			},
			{
				Name:        "Periksa Surat Jalan & Absensi",
				Description: "Analisis kecocokan ringkasan surat jalan dengan data kehadiran/lembur karyawan.",
				SystemPrompt: "Anda adalah petugas verifikasi operasional. Bandingkan ringkasan log surat jalan dengan jadwal absensi dan lembur yang diinput. Laporkan jika ada kejanggalan tanggal atau jam kerja.",
			},
		},
	}
}

// LoadConfig attempts to read config from the default path (~/.ollama-tui/config.json)
