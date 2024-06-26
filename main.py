from helpers import calculate_percentages, calculate_money_to_insert

# If there is an error in the reading, manually analyze the data and increase the moment if necessary
tickers = {
    "ALUP11": {"fundamentalist grade": 70, "moment": 0, "final grade": 0, "weight": 0},
    "BMEB4": {"fundamentalist grade": 85, "moment": 1, "final grade": 0, "weight": 0},
    "CSAN3": {"fundamentalist grade": 75, "moment": 0, "final grade": 0, "weight": 0},
    "EGIE3": {"fundamentalist grade": 95, "moment": 0, "final grade": 0, "weight": 0},
    "FESA4": {"fundamentalist grade": 70, "moment": 0, "final grade": 0, "weight": 0},
    "ITSA4": {"fundamentalist grade": 100, "moment": 0, "final grade": 0, "weight": 0},
    "KLBN11": {"fundamentalist grade": 80, "moment": 0, "final grade": 0, "weight": 0},
    "SUZB3": {"fundamentalist grade": 75, "moment": 0, "final grade": 0, "weight": 0},
    "TUPY3": {"fundamentalist grade": 75, "moment": 0, "final grade": 0, "weight": 0},
    "UNIP6": {"fundamentalist grade": 80, "moment": 0, "final grade": 0, "weight": 0},
    "VIVT3": {"fundamentalist grade": 75, "moment": 0, "final grade": 0, "weight": 0},
    "WIZC3": {"fundamentalist grade": 75, "moment": 0, "final grade": 0, "weight": 0},
    "FLRY3": {"fundamentalist grade": 70, "moment": 0, "final grade": 0, "weight": 0},
}
money = 1800

weights = calculate_percentages(tickers)
calculate_money_to_insert(money, weights)
