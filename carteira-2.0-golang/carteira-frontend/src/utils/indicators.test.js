import { describe, it, expect } from 'vitest';
import * as fc from 'fast-check';
import { getIndicatorColor } from './indicators.js';

// Feature: carteira-frontend, Property 7: Coloração correta dos indicadores PE, PBV, PSR, DY
describe('Property 7: Coloração correta dos indicadores PE, PBV, PSR, DY', () => {
  it('PE: green se pe > 0 && pe <= 8, red caso contrário', () => {
    fc.assert(
      fc.property(
        fc.float({ min: -1000, max: 1000, noNaN: true }),
        (pe) => {
          const color = getIndicatorColor('pe', { pe, invalid_fields: [] });
          if (pe > 0 && pe <= 8) return color === 'green';
          return color === 'red';
        }
      ),
      { numRuns: 100 }
    );
  });

  it('PBV: green se pbv > 0 && pbv <= 2, red caso contrário', () => {
    fc.assert(
      fc.property(
        fc.float({ min: -1000, max: 1000, noNaN: true }),
        (pbv) => {
          const color = getIndicatorColor('pbv', { pbv, invalid_fields: [] });
          if (pbv > 0 && pbv <= 2) return color === 'green';
          return color === 'red';
        }
      ),
      { numRuns: 100 }
    );
  });

  it('PSR: green se psr > 0 && psr < 2, red caso contrário', () => {
    fc.assert(
      fc.property(
        fc.float({ min: -1000, max: 1000, noNaN: true }),
        (psr) => {
          const color = getIndicatorColor('psr', { psr, invalid_fields: [] });
          if (psr > 0 && psr < 2) return color === 'green';
          return color === 'red';
        }
      ),
      { numRuns: 100 }
    );
  });

  it('DY: green se dy >= 4, red caso contrário', () => {
    fc.assert(
      fc.property(
        fc.float({ min: -1000, max: 1000, noNaN: true }),
        (dy) => {
          const color = getIndicatorColor('dy', { dy, invalid_fields: [] });
          if (dy >= 4) return color === 'green';
          return color === 'red';
        }
      ),
      { numRuns: 100 }
    );
  });
});

// Feature: carteira-frontend, Property 8: Coloração correta do indicador Graham
describe('Property 8: Coloração correta do indicador Graham', () => {
  it('Graham: green se price < sqrt(22.5 * eps * bvps), red caso contrário (eps > 0, bvps > 0)', () => {
    fc.assert(
      fc.property(
        fc.record({
          price: fc.double({ min: 0.01, max: 10000, noNaN: true }),
          eps: fc.double({ min: 0.01, max: 1000, noNaN: true }),
          bvps: fc.double({ min: 0.01, max: 1000, noNaN: true }),
        }),
        ({ price, eps, bvps }) => {
          const color = getIndicatorColor('graham', { price, eps, bvps, invalid_fields: [] });
          const grahamValue = Math.sqrt(22.5 * eps * bvps);
          if (price < grahamValue) return color === 'green';
          return color === 'red';
        }
      ),
      { numRuns: 100 }
    );
  });

  it('Graham: red se eps <= 0 ou bvps <= 0', () => {
    fc.assert(
      fc.property(
        fc.record({
          price: fc.double({ min: 0.01, max: 10000, noNaN: true }),
          eps: fc.double({ min: -1000, max: 0, noNaN: true }),
          bvps: fc.double({ min: 0.01, max: 1000, noNaN: true }),
        }),
        ({ price, eps, bvps }) => {
          const color = getIndicatorColor('graham', { price, eps, bvps, invalid_fields: [] });
          return color === 'red';
        }
      ),
      { numRuns: 100 }
    );
  });
});

// Feature: carteira-frontend, Property 9: Campos inválidos omitem coloração
describe('Property 9: Campos inválidos omitem coloração', () => {
  it('retorna neutral quando o campo relevante está em invalid_fields', () => {
    const indicatorToInvalidField = {
      pe: 'PE',
      pbv: 'PBV',
      psr: 'PSR',
      dy: 'DY',
    };

    fc.assert(
      fc.property(
        fc.constantFrom('pe', 'pbv', 'psr', 'dy'),
        fc.float({ min: -1000, max: 1000, noNaN: true }),
        (indicator, value) => {
          const invalidField = indicatorToInvalidField[indicator];
          const fundamentals = {
            pe: value,
            pbv: value,
            psr: value,
            dy: value,
            invalid_fields: [invalidField],
          };
          const color = getIndicatorColor(indicator, fundamentals);
          return color === 'neutral';
        }
      ),
      { numRuns: 100 }
    );
  });

  it('Graham: retorna neutral quando EPS, BVps ou Price está em invalid_fields', () => {
    fc.assert(
      fc.property(
        fc.array(fc.constantFrom('EPS', 'BVps', 'Price'), { minLength: 1, maxLength: 3 }),
        fc.double({ min: 0.01, max: 1000, noNaN: true }),
        fc.double({ min: 0.01, max: 1000, noNaN: true }),
        fc.double({ min: 0.01, max: 10000, noNaN: true }),
        (invalidFields, eps, bvps, price) => {
          const color = getIndicatorColor('graham', { eps, bvps, price, invalid_fields: invalidFields });
          return color === 'neutral';
        }
      ),
      { numRuns: 100 }
    );
  });
});
