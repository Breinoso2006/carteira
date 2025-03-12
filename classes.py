import math, re
from time import sleep
from selenium.webdriver.common.by import By
from loguru import logger
from termcolor import colored


class Acao:
    def __init__(self, ticker, driver):
        self.ticker = ticker
        self.driver = driver
        self.moment = 0
        self.moment_indicators = {
            "price": {
                "identifier": "/html/body/div[1]/div[2]/table[1]/tbody/tr[1]/td[4]/span",
                "identifier2": "//*[@id='asset-header']/div/div[1]/div[2]/div[1]/span",
            },
            "pvp": {
                "identifier": "/html/body/div[1]/div[2]/table[3]/tbody/tr[3]/td[4]/span",
                "identifier2": "//*[@id='PAGE_CONTAINER']/main/div[3]/div/div[1]/div[2]/div/div[2]/div[2]/div[2]/p",
            },
            "pl": {
                "identifier": "/html/body/div[1]/div[2]/table[3]/tbody/tr[2]/td[4]/span",
                "identifier2": "//*[@id='PAGE_CONTAINER']/main/div[3]/div/div[1]/div[2]/div/div[2]/div[2]/div[1]/p",
            },
            "dy": {
                "identifier": "/html/body/div[1]/div[2]/table[3]/tbody/tr[9]/td[4]/span",
                "identifier2": "//*[@id='PAGE_CONTAINER']/main/div[3]/div/div[1]/div[2]/div/div[2]/div[2]/div[3]/p",
            },
            "psr": {
                "identifier": "/html/body/div[1]/div[2]/table[3]/tbody/tr[5]/td[4]/span",
                "identifier2": "//*[@id='PAGE_CONTAINER']/main/div[3]/div/div[2]/div[2]/div/div[2]/div[1]/div/div[8]/div[2]/span",
            },
            "lpa": {
                "identifier": "/html/body/div[1]/div[2]/table[3]/tbody/tr[2]/td[6]/span",
                "identifier2": "//*[@id='PAGE_CONTAINER']/main/div[3]/div/div[2]/div[2]/div/div[2]/div[1]/div/div[4]/div[2]/span",
            },
            "vpa": {
                "identifier": "/html/body/div[1]/div[2]/table[3]/tbody/tr[3]/td[6]/span",
                "identifier2": "//*[@id='PAGE_CONTAINER']/main/div[3]/div/div[2]/div[2]/div/div[2]/div[1]/div/div[5]/div[2]/span",
            },
        }

        self.__set_indicators()
        self.__transient_moment_calculation()

    def __set_indicators(self):

        self.driver.get(
            f"https://www.fundamentus.com.br/detalhes.php?papel={self.ticker}"
        )

        for indicator in self.moment_indicators.values():
            try:
                value = self.driver.find_element(
                    By.XPATH, indicator["identifier"]
                ).text.replace(",", ".")

                if value == None or not re.search(r'\d', value):
                    self.driver.get(
                        f"https://analitica.auvp.com.br/ativos/{self.ticker}"
                    )
                    sleep(2)
                    value = self.driver.find_element(
                        By.XPATH, indicator["identifier2"]
                    ).text.replace(",", ".")
                    self.driver.get(
                        f"https://www.fundamentus.com.br/detalhes.php?papel={self.ticker}"
                    )
                    
                indicator["value"] = value
            except Exception as error:
                logger.warning(
                f"Error while getting information about {self.ticker}: {error}"
            )
                indicator["value"] = 0

    def get_moment(self):
        return self.moment

    def __transient_moment_calculation(self):
        moment = 0
        moment += self._pvp_grade_calculation()
        moment += self._psr_grade_calculation()
        moment += self._pl_grade_calculation()
        moment += self._dividend_yield_grade_calculation()
        moment += self._graham_calculation()

        self.moment = moment

    def _pvp_grade_calculation(self):
        try:
            pvp = float(self.moment_indicators["pvp"]["value"])
            if pvp < 2:
                print(colored(f"{self.ticker} PVP = {pvp}", "green"))
                return 1
            print(colored(f"{self.ticker} PVP = {pvp}", "red"))
        except Exception as error:
            logger.warning(
                f"Error calculating PVP indicator from ticker {self.ticker}: {error}"
            )

        return 0

    def _psr_grade_calculation(self):
        try:
            psr = float(self.moment_indicators["psr"]["value"])
            if psr < 2 and psr > 0:
                print(colored(f"{self.ticker} PSR = {psr}", "green"))
                return 1
            print(colored(f"{self.ticker} PSR = {psr}", "red"))
        except Exception as error:
            logger.warning(
                f"Error calculating PSR indicator from ticker {self.ticker}: {error}"
            )

        return 0

    def _pl_grade_calculation(self):
        try:
            pl = float(self.moment_indicators["pl"]["value"])
            if pl <= 8 and pl > 0:
                print(colored(f"{self.ticker} PL = {pl}", "green"))
                return 1
            print(colored(f"{self.ticker} PL = {pl}", "red"))
        except Exception as error:
            logger.warning(
                f"Error calculating PL indicator from ticker {self.ticker}: {error}"
            )

        return 0

    def _dividend_yield_grade_calculation(self):
        try:
            dy = self.moment_indicators["dy"]["value"]
            dy = float(dy.replace("%", ""))
            if dy >= 4:
                print(colored(f"{self.ticker} DY = {dy}", "green"))
                return 1
            print(colored(f"{self.ticker} DY = {dy}", "red"))
        except Exception as error:
            logger.warning(
                f"Error calculating DY indicator from ticker {self.ticker}: {error}"
            )

        return 0

    def _graham_calculation(self):
        nota = 0

        try:
            price = float(self.moment_indicators["price"]["value"])
            lpa = float(self.moment_indicators["lpa"]["value"])
            vpa = float(self.moment_indicators["vpa"]["value"])
            graham = math.sqrt(22.5 * lpa * vpa)
            graham_margin = graham * 0.7
            print(f"{self.ticker} Price = {price}")
            if price < graham:
                print(colored(f"{self.ticker} Graham = {graham}", "green"))
                nota += 1
                if price < graham_margin:
                    print(
                        colored(
                            f"{self.ticker} Graham Margin = {graham_margin}", "green"
                        )
                    )
                    nota += 1
                else:
                    print(
                        colored(f"{self.ticker} Graham Margin = {graham_margin}", "red")
                    )
            else:
                print(colored(f"{self.ticker} Graham = {graham}", "red"))
                print(colored(f"{self.ticker} Graham Margin = {graham_margin}", "red"))

        except Exception as error:
            logger.warning(
                f"Error calculating Graham indicator from ticker {self.ticker}: {error}"
            )
        finally:
            return nota
