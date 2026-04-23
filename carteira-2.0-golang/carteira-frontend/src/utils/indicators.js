/**
 * Retorna a cor de um indicador fundamentalista.
 * @param {string} indicator - 'pe' | 'pbv' | 'psr' | 'dy' | 'graham'
 * @param {object} fundamentals - objeto com os campos da scraping-api
 * @returns {'green' | 'red' | 'neutral'}
 */
export function getIndicatorColor(indicator, fundamentals) {
  const { pe, pbv, psr, dy, price, eps, bvps, invalid_fields = [] } = fundamentals;

  switch (indicator) {
    case 'pe':
      if (invalid_fields.includes('PE') || pe == null) return 'neutral';
      return pe > 0 && pe <= 8 ? 'green' : 'red';

    case 'pbv':
      if (invalid_fields.includes('PBV') || pbv == null) return 'neutral';
      return pbv > 0 && pbv <= 2 ? 'green' : 'red';

    case 'psr':
      if (invalid_fields.includes('PSR') || psr == null) return 'neutral';
      return psr > 0 && psr < 2 ? 'green' : 'red';

    case 'dy':
      if (invalid_fields.includes('DY') || dy == null) return 'neutral';
      return dy >= 4 ? 'green' : 'red';

    case 'graham': {
      const depsInvalid = ['EPS', 'BVps', 'Price'].some(f => invalid_fields.includes(f));
      if (depsInvalid || eps == null || bvps == null || price == null) return 'neutral';
      if (eps <= 0 || bvps <= 0) return 'red';
      const grahamValue = Math.sqrt(22.5 * eps * bvps);
      return price < grahamValue ? 'green' : 'red';
    }

    default:
      return 'neutral';
  }
}
