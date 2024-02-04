from loguru import logger
from selenium import webdriver
from classes import Acao


def printar_momentos(tickers: dict):
    driver = webdriver.Safari()
    for ticker in tickers:
        try:
            momento = Acao(ticker, driver).get_momento()
            print(ticker, "=", momento)
        except Exception as error:
            logger.error(f"Erro ticker {ticker}: {error}")

    driver.close()
    return tickers


def calcular_pesos(tickers: dict):
    nota_total = 0
    pesos = []

    try:
        for ticker in tickers.values():
            nota_momentanea = round(
                ticker["nota fundamentalista"]
                * (ticker["nota fundamentalista"] / 100 + ticker["momento"] / 6)
                / 2
            )

            if nota_momentanea >= 60:
                ticker["nota final"] = nota_momentanea * 0.40
            elif nota_momentanea >= 50:
                ticker["nota final"] = nota_momentanea * 0.35
            else:
                ticker["nota final"] = nota_momentanea * 0.25

            nota_total += ticker["nota final"]

        for ticker, dados in tickers.items():
            peso_final = round(dados["nota final"] / nota_total * 100, 2)
            pesos.append({f"{ticker}": f"{peso_final} %"})

        return pesos

    except Exception as error:
        logger.error(error)
