package scraper

import (
	"regexp"
	"strings"
)

// GetSourceConfigs retorna as configurações de todas as fontes disponíveis
func GetSourceConfigs() map[string]SourceConfig {
	return map[string]SourceConfig{
		"investidor10": getInvestidor10Config(),
		"auvp":         getAuvpConfig(),
		"fundamentus":  getFundamentusConfig(),
	}
}

// =============================================================================
// INVESTIDOR10
// =============================================================================

func getInvestidor10Config() SourceConfig {
	return SourceConfig{
		Source:    "investidor10",
		URLPrefix: "https://investidor10.com.br/acoes/",
		Selectors: map[string]string{
			"price": "#cards-ticker > div._card.cotacao > div._card-body > div > span",
			"pe":    "#cards-ticker > div._card.val > div._card-body > span",
			"psr":   "#table-indicators > div:nth-child(2) > div.value.d-flex.justify-content-between.align-items-center > span",
			"bvps":  "#table-indicators > div:nth-child(14) > div.value.d-flex.justify-content-between.align-items-center > span",
			"eps":   "#table-indicators > div:nth-child(15) > div.value.d-flex.justify-content-between.align-items-center > span",
			"dy":    "#cards-ticker > div._card.dy > div._card-body > span",
		},
		Cleaners: map[string]func(string) string{
			"price": cleanInvestidor10Price,
			"pe":    cleanInvestidor10Numeric,
			"psr":   cleanInvestidor10Numeric,
			"bvps":  cleanInvestidor10Price,
			"eps":   cleanInvestidor10Price,
			"dy":    cleanInvestidor10Percentage,
		},
	}
}

// Investidor10: remove R$, normaliza separadores decimais
func cleanInvestidor10Price(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "R$", "")
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, ",", ".")
	re := regexp.MustCompile(`[^\d.]`)
	value = re.ReplaceAllString(value, "")
	return value
}

// Investidor10: remove espaços e símbolos especiais
func cleanInvestidor10Numeric(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, ",", ".")
	value = strings.ReplaceAll(value, " ", "")
	re := regexp.MustCompile(`[^\d.]`)
	value = re.ReplaceAllString(value, "")
	return value
}

// Investidor10: remove % e normaliza percentual
func cleanInvestidor10Percentage(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "%", "")
	value = strings.ReplaceAll(value, ",", ".")
	value = strings.ReplaceAll(value, " ", "")
	re := regexp.MustCompile(`[^\d.]`)
	value = re.ReplaceAllString(value, "")
	return value
}

// =============================================================================
// AUVP
// =============================================================================

func getAuvpConfig() SourceConfig {
	return SourceConfig{
		Source:    "auvp",
		URLPrefix: "https://analitica.auvp.com.br/acoes/",
		Selectors: map[string]string{
			"price": "#asset-header > div > div.flex.w-full.items-center > div.w-fit.items-center.h-fit.ml-auto.flex > div.w-fit.h-fit.flex.items-center.justify-center.mr-2 > span",
			"pe":    "#PAGE_CONTAINER > main > div.mx-2.md\\:mx-0 > div > div:nth-child(3) > div.w-full.max-w-\\[var\\(--max-width\\)\\].mx-auto.flex.flex-col > div > div:nth-child(2) > div.flex.w-full.flex-col.p-2.gap-4 > div:nth-child(1) > p",
			"psr":   "#PAGE_CONTAINER > main > div.mx-2.md\\:mx-0 > div > div:nth-child(6) > div.w-full.max-w-\\[var\\(--max-width\\)\\].mx-auto.flex.flex-col > div.bg-card.rounded-xl.p-4.relative.group.max-w-\\[100vw\\].w-full.mt-12 > div.md\\:px-6.md\\:pb-6.relative.mt-4.p-0.overflow-y-hidden > div:nth-child(1) > div > div:nth-child(8) > div.mt-auto.mb-1.flex.justify-between.items-center > span",
			"bvps":  "#PAGE_CONTAINER > main > div.mx-2.md\\:mx-0 > div > div:nth-child(6) > div.w-full.max-w-\\[var\\(--max-width\\)\\].mx-auto.flex.flex-col > div.bg-card.rounded-xl.p-4.relative.group.max-w-\\[100vw\\].w-full.mt-12 > div.md\\:px-6.md\\:pb-6.relative.mt-4.p-0.overflow-y-hidden > div:nth-child(1) > div > div:nth-child(5) > div.mt-auto.mb-1.flex.justify-between.items-center > span",
			"eps":   "#PAGE_CONTAINER > main > div.mx-2.md\\:mx-0 > div > div:nth-child(6) > div.w-full.max-w-\\[var\\(--max-width\\)\\].mx-auto.flex.flex-col > div.bg-card.rounded-xl.p-4.relative.group.max-w-\\[100vw\\].w-full.mt-12 > div.md\\:px-6.md\\:pb-6.relative.mt-4.p-0.overflow-y-hidden > div:nth-child(1) > div > div:nth-child(4) > div.mt-auto.mb-1.flex.justify-between.items-center > span",
			"dy":    "#PAGE_CONTAINER > main > div.mx-2.md\\:mx-0 > div > div:nth-child(6) > div.w-full.max-w-\\[var\\(--max-width\\)\\].mx-auto.flex.flex-col > div.bg-card.rounded-xl.p-4.relative.group.max-w-\\[100vw\\].w-full.mt-12 > div.md\\:px-6.md\\:pb-6.relative.mt-4.p-0.overflow-y-hidden > div:nth-child(1) > div > div:nth-child(1) > div.mt-auto.mb-1.flex.justify-between.items-center > span",
		},
		Cleaners: map[string]func(string) string{
			"price": cleanAuvpPrice,
			"pe":    cleanAuvpRatio,
			"psr":   cleanAuvpRatio,
			"bvps":  cleanAuvpPrice,
			"eps":   cleanAuvpPrice,
			"dy":    cleanAuvpPercentage,
		},
	}
}

// Auvp: trata múltiplos formatos de separadores
func cleanAuvpPrice(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "R$", "")
	value = strings.TrimSpace(value)

	// Auvp pode usar ponto ou vírgula - converter para padrão
	if strings.Contains(value, ",") && strings.LastIndex(value, ",") > strings.LastIndex(value, ".") {
		value = strings.ReplaceAll(value, ".", "")
		value = strings.ReplaceAll(value, ",", ".")
	}
	value = strings.ReplaceAll(value, " ", "")
	re := regexp.MustCompile(`[^\d.]`)
	value = re.ReplaceAllString(value, "")
	return value
}

// Auvp: remove "x" e trata valores de ratio
func cleanAuvpRatio(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "x", "")
	value = strings.ReplaceAll(value, "X", "")

	if strings.Contains(value, ",") && strings.LastIndex(value, ",") > strings.LastIndex(value, ".") {
		value = strings.ReplaceAll(value, ".", "")
		value = strings.ReplaceAll(value, ",", ".")
	}
	value = strings.ReplaceAll(value, " ", "")
	re := regexp.MustCompile(`[^\d.]`)
	value = re.ReplaceAllString(value, "")
	return value
}

// Auvp: remove % e trata percentual
func cleanAuvpPercentage(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "%", "")

	if strings.Contains(value, ",") && strings.LastIndex(value, ",") > strings.LastIndex(value, ".") {
		value = strings.ReplaceAll(value, ".", "")
		value = strings.ReplaceAll(value, ",", ".")
	}
	value = strings.ReplaceAll(value, " ", "")
	re := regexp.MustCompile(`[^\d.]`)
	value = re.ReplaceAllString(value, "")
	return value
}

// =============================================================================
// FUNDAMENTUS
// =============================================================================

func getFundamentusConfig() SourceConfig {
	return SourceConfig{
		Source:    "fundamentus",
		URLPrefix: "https://www.fundamentus.com.br/detalhes.php?papel=",
		Selectors: map[string]string{
			"price": "body > div.center > div.conteudo.clearfix > table:nth-child(2) > tbody > tr:nth-child(1) > td.data.destaque.w3 > span",
			"pe":    "body > div.center > div.conteudo.clearfix > table:nth-child(4) > tbody > tr:nth-child(2) > td:nth-child(4) > span",
			"psr":   "body > div.center > div.conteudo.clearfix > table:nth-child(4) > tbody > tr:nth-child(5) > td:nth-child(4) > span",
			"bvps":  "body > div.center > div.conteudo.clearfix > table:nth-child(4) > tbody > tr:nth-child(3) > td:nth-child(6) > span",
			"eps":   "body > div.center > div.conteudo.clearfix > table:nth-child(4) > tbody > tr:nth-child(2) > td:nth-child(6) > span",
			"dy":    "body > div.center > div.conteudo.clearfix > table:nth-child(4) > tbody > tr:nth-child(9) > td:nth-child(4) > span",
		},
		Cleaners: map[string]func(string) string{
			"price": cleanFundamentusPrice,
			"pe":    cleanFundamentusNumeric,
			"psr":   cleanFundamentusNumeric,
			"bvps":  cleanFundamentusPrice,
			"eps":   cleanFundamentusPrice,
			"dy":    cleanFundamentusPercentage,
		},
	}
}

// Fundamentus: trata separadores de milhar e decimais
func cleanFundamentusPrice(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "R$", "")
	value = strings.ReplaceAll(value, "%", "")
	value = strings.TrimSpace(value)

	// Fundamentus pode usar ponto como separador de milhar
	// Lógica: se há mais de um ponto, remove os anteriores
	countDots := strings.Count(value, ".")
	if countDots > 1 {
		parts := strings.Split(value, ".")
		value = strings.Join(parts[:len(parts)-1], "") + "." + parts[len(parts)-1]
	}
	if strings.Contains(value, ",") {
		value = strings.ReplaceAll(value, ",", ".")
	}
	value = strings.ReplaceAll(value, " ", "")
	re := regexp.MustCompile(`[^\d.]`)
	value = re.ReplaceAllString(value, "")
	return value
}

// Fundamentus: normaliza valores numéricos genéricos
func cleanFundamentusNumeric(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "x", "")
	value = strings.ReplaceAll(value, "X", "")
	value = strings.ReplaceAll(value, "-", "")

	countDots := strings.Count(value, ".")
	if countDots > 1 {
		parts := strings.Split(value, ".")
		value = strings.Join(parts[:len(parts)-1], "") + "." + parts[len(parts)-1]
	}
	if strings.Contains(value, ",") {
		value = strings.ReplaceAll(value, ",", ".")
	}
	value = strings.ReplaceAll(value, " ", "")
	re := regexp.MustCompile(`[^\d.]`)
	value = re.ReplaceAllString(value, "")
	return value
}

// Fundamentus: remove % e trata percentual
func cleanFundamentusPercentage(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "%", "")

	countDots := strings.Count(value, ".")
	if countDots > 1 {
		parts := strings.Split(value, ".")
		value = strings.Join(parts[:len(parts)-1], "") + "." + parts[len(parts)-1]
	}
	if strings.Contains(value, ",") {
		value = strings.ReplaceAll(value, ",", ".")
	}
	value = strings.ReplaceAll(value, " ", "")
	re := regexp.MustCompile(`[^\d.]`)
	value = re.ReplaceAllString(value, "")
	return value
}
