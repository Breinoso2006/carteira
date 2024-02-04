from helpers import printar_momentos, calcular_pesos
from pprint import pprint as pp

tickers = {
    "ALUP11": {"nota fundamentalista": 70, "momento": 3, "nota final": 0, "peso": 0},
    "BMEB4": {"nota fundamentalista": 85, "momento": 5, "nota final": 0, "peso": 0},
    "CSAN3": {"nota fundamentalista": 75, "momento": 2, "nota final": 0, "peso": 0},
    "EGIE3": {"nota fundamentalista": 95, "momento": 1, "nota final": 0, "peso": 0},
    "FESA4": {"nota fundamentalista": 85, "momento": 5, "nota final": 0, "peso": 0},
    "ITSA4": {"nota fundamentalista": 100, "momento": 4, "nota final": 0, "peso": 0},
    "KLBN11": {"nota fundamentalista": 90, "momento": 3, "nota final": 0, "peso": 0},
    "SUZB3": {"nota fundamentalista": 75, "momento": 5, "nota final": 0, "peso": 0},
    "TUPY3": {"nota fundamentalista": 75, "momento": 4, "nota final": 0, "peso": 0},
    "UNIP6": {"nota fundamentalista": 85, "momento": 2, "nota final": 0, "peso": 0},
    "VIVT3": {"nota fundamentalista": 75, "momento": 3, "nota final": 0, "peso": 0},
    "WIZC3": {"nota fundamentalista": 75, "momento": 2, "nota final": 0, "peso": 0},
}

# printar_momentos(tickers)
pp(calcular_pesos(tickers))
