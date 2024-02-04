import math
from selenium.webdriver.common.by import By
from loguru import logger


class Acao:
    def __init__(self, ticker, driver):
        self.ticker = ticker
        self.driver = driver
        self.momento = 0
        self.indicadores_de_momento = {
            "cotacao": {
                "identificador": "/html/body/div[1]/div[2]/table[1]/tbody/tr[1]/td[4]/span",
                "valor": "",
            },
            "pvp": {
                "identificador": "/html/body/div[1]/div[2]/table[3]/tbody/tr[3]/td[4]/span",
                "valor": "",
            },
            "pl": {
                "identificador": "/html/body/div[1]/div[2]/table[3]/tbody/tr[2]/td[4]/span",
                "valor": "",
            },
            "dy": {
                "identificador": "/html/body/div[1]/div[2]/table[3]/tbody/tr[9]/td[4]/span",
                "valor": "",
            },
            "psr": {
                "identificador": "/html/body/div[1]/div[2]/table[3]/tbody/tr[5]/td[4]/span",
                "valor": "",
            },
            "lpa": {
                "identificador": "/html/body/div[1]/div[2]/table[3]/tbody/tr[2]/td[6]/span",
                "valor": "",
            },
            "vpa": {
                "identificador": "/html/body/div[1]/div[2]/table[3]/tbody/tr[3]/td[6]/span",
                "valor": "",
            },
        }

        self._setar_indicadores()
        self._calcular_momento_momentaneo()

    def _setar_indicadores(self):

        self.driver.get(
            f"https://www.fundamentus.com.br/detalhes.php?papel={self.ticker}"
        )

        for indicador in self.indicadores_de_momento.values():
            indicador["valor"] = self.driver.find_element(
                By.XPATH, indicador["identificador"]
            ).text.replace(",", ".")

    def get_momento(self):
        return self.momento

    def _calcular_momento_momentaneo(self):
        momento = 0
        momento += self._calcular_nota_pvp()
        momento += self._calcular_nota_psr()
        momento += self._calcular_nota_pl()
        momento += self._calcular_dividend_yield()
        momento += self._calcular_graham()

        self.momento = momento

    def _calcular_nota_pvp(self):
        try:
            pvp = float(self.indicadores_de_momento["pvp"]["valor"])
            if pvp < 2:
                return 1
        except Exception as error:
            logger.warning(
                f"Erro ao calcular nota do PVP do ticker {self.ticker}: {error}"
            )

        return 0

    def _calcular_nota_psr(self):
        try:
            psr = float(self.indicadores_de_momento["psr"]["valor"])
            if psr < 2:
                return 1
        except Exception as error:
            logger.warning(
                f"Erro ao calcular nota do PSR do ticker {self.ticker}: {error}"
            )

        return 0

    def _calcular_nota_pl(self):
        try:
            pl = float(self.indicadores_de_momento["pl"]["valor"])
            if pl <= 6 and pl > 0:
                return 1
        except Exception as error:
            logger.warning(
                f"Erro ao calcular nota do PL do ticker {self.ticker}: {error}"
            )

        return 0

    def _calcular_dividend_yield(self):
        try:
            dy = self.indicadores_de_momento["dy"]["valor"]
            dy = float(dy.replace("%", ""))
            if dy >= 4:
                return 1
        except Exception as error:
            logger.warning(
                f"Erro ao calcular nota do DY do ticker {self.ticker}: {error}"
            )

        return 0

    def _calcular_graham(self):
        nota = 0

        try:
            cotacao = float(self.indicadores_de_momento["cotacao"]["valor"])
            lpa = float(self.indicadores_de_momento["lpa"]["valor"])
            vpa = float(self.indicadores_de_momento["vpa"]["valor"])
            graham = math.sqrt(22.5 * lpa * vpa)

            if cotacao < graham:
                nota += 1
                if cotacao < graham * 0.7:
                    nota += 1

        except Exception as error:
            logger.warning(f"Erro ao calcular Graham do ticker {self.ticker}: {error}")
        finally:
            return nota
