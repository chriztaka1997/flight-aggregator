package utils

// GetTimezone returns the timezone for Indonesian airports
// Maps airport codes to their respective timezone locations
func GetTimezone(airportCode string) string {
	// WIB (Western Indonesian Time): UTC+7
	wibAirports := map[string]bool{
		"CGK": true, // Jakarta - Soekarno-Hatta International
		"SUB": true, // Surabaya - Juanda International
		"KNO": true, // Medan - Kuala Namu International
		"PLM": true, // Palembang - Sultan Mahmud Badaruddin II
		"PDG": true, // Padang - Minangkabau International
		"BTJ": true, // Banda Aceh - Sultan Iskandar Muda
		"PKU": true, // Pekanbaru - Sultan Syarif Kasim II
		"BDO": true, // Bandung - Husein Sastranegara
		"SRG": true, // Semarang - Achmad Yani
		"JOG": true, // Yogyakarta - Adisucipto
	}

	// WITA (Central Indonesian Time): UTC+8
	witaAirports := map[string]bool{
		"DPS": true, // Denpasar (Bali) - Ngurah Rai International
		"UPG": true, // Makassar - Sultan Hasanuddin International
		"BPN": true, // Balikpapan - Sultan Aji Muhammad Sulaiman
		"LOP": true, // Lombok - Lombok International
		"BDJ": true, // Banjarmasin - Syamsudin Noor
		"SOC": true, // Solo - Adisumarmo International
		"PKY": true, // Palangkaraya - Tjilik Riwut
		"MDC": true, // Manado - Sam Ratulangi
	}

	// WIT (Eastern Indonesian Time): UTC+9
	witAirports := map[string]bool{
		"DJJ": true, // Jayapura - Sentani International
		"AMQ": true, // Ambon - Pattimura
		"TIM": true, // Timika - Moses Kilangin
		"SRR": true, // Sorong - Domine Eduard Osok
	}

	if wibAirports[airportCode] {
		return "Asia/Jakarta" // WIB (UTC+7)
	}
	if witaAirports[airportCode] {
		return "Asia/Makassar" // WITA (UTC+8)
	}
	if witAirports[airportCode] {
		return "Asia/Jayapura" // WIT (UTC+9)
	}

	// Default to WIB (most common in Indonesia)
	return "Asia/Jakarta"
}
