from loguru import logger
from selenium import webdriver
from classes import Acao
from pprint import pprint as pp


def calculate_percentages(tickers: dict):
    nota_total = 0
    weights = []

    try:
        driver = webdriver.Safari()

        for ticker, values in tickers.items():
            values["moment"] += Acao(ticker, driver).get_moment()
            print(f"{ticker} moment =  {values["moment"]} \n")
            transiente_grade = round(
                values["fundamentalist grade"]
                * (values["fundamentalist grade"] / 100 + values["moment"] / 6)
                / 2
            )

            if transiente_grade >= 60:
                values["final grade"] = transiente_grade * 0.40
            elif transiente_grade >= 50:
                values["final grade"] = transiente_grade * 0.35
            else:
                values["final grade"] = transiente_grade * 0.25

            nota_total += values["final grade"]

        driver.close()

        for ticker, values in tickers.items():
            final_weight = round(values["final grade"] / nota_total * 100, 2)
            weights.append({f"{ticker}": f"{final_weight} %"})

        return pp(weights)

    except Exception as error:
        logger.error(error)
