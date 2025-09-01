package config

type TokenScope string
type Language string

const (
	TokenScopeCreatePayment TokenScope = "payment-create"
	TokenScopeAll           TokenScope = "payment-all"

	// Czech language code
	CZECH Language = "CS"
	// English language code
	ENGLISH Language = "EN"
	// Slovak language code
	SLOVAK Language = "SK"
	// German language code
	GERMAN Language = "DE"
	// Russian language code
	RUSSIAN Language = "RU"
	// Polish language code
	POLISH Language = "PL"
	// Hungarian language code
	HUNGARIAN Language = "HU"
	// French language code
	FRENCH Language = "FR"
	// Romanian language code
	ROMANIAN Language = "RO"
	// Bulgarian language code
	BULGARIAN Language = "BG"
	// Croatian language code
	CROATIAN Language = "HR"
	// Italian language code
	ITALIAN Language = "IT"
	// Spanish language code
	SPANISH Language = "ES"
	// Ukrainian language code
	UKRAINIAN Language = "UK"
	// Estonian language code
	ESTONIAN Language = "ET"
	// Lithuanian language code
	LITHUANIAN Language = "LT"
	// Latvian language code
	LATVIAN Language = "LV"
	// Slovenian language code
	SLOVENIAN Language = "SL"
	// Portuguese language code
	PORTUGUESE Language = "PT"
)
