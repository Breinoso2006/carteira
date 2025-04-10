from loguru import logger
from selenium import webdriver
from classes import Acao
from pprint import pprint as pp
import math


def calculate_percentages(tickers: dict, with_moment: bool = False):
    print(f"\nCalculating percentages for each Ticker, with moment {with_moment}:\n")
    nota_total = 0
    weights = []

    try:
        if with_moment:
            driver = webdriver.Safari()

        for ticker, values in tickers.items():
            if with_moment:
                values["moment"] += Acao(ticker, driver).get_moment()
                print(f"{ticker} moment =  {values["moment"]} \n")
            else:
                values["moment"] += 6
            transiente_grade = round(
                values["fundamentalist grade"]
                * (values["fundamentalist grade"] / 100 + math.sqrt(values["moment"] / 6))
                / 2
            )

            if transiente_grade >= 70:
                values["final grade"] = transiente_grade * 0.40
            elif transiente_grade >= 60:
                values["final grade"] = transiente_grade * 0.30
            elif transiente_grade >= 50:
                values["final grade"] = transiente_grade * 0.15
            else:
                values["final grade"] = transiente_grade * 0.05

            nota_total += values["final grade"]

        if with_moment:
            driver.close()

        weights = {}
        for ticker, values in tickers.items():
            final_weight = round(values["final grade"] / nota_total * 100, 2)
            weights.update({f"{ticker}": final_weight / 100})
        
        pp(weights)
        
        return weights

    except Exception as error:
        driver.close()
        logger.error(error)

def calculate_money_to_insert(money: int, weights: dict):
    print("\nCalculating money to insert in each Ticker:\n")
    for ticker, weight in weights.items():
        print(f"{ticker} = ", money * weight)