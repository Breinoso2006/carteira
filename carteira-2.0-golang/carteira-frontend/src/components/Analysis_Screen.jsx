// Analysis_Screen — tela de análise de pesos do portfolio
// Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6

import { useState, useEffect } from 'react';
import { getPortfolio } from '../api/client';

/**
 * Props:
 *   onNavigate: (screen: string) => void
 */
export default function Analysis_Screen({ onNavigate }) {
  const [stocks, setStocks] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    async function loadPortfolio() {
      setError(null);
      setLoading(true);
      try {
        const data = await getPortfolio();
        setStocks(data);
      } catch (e) {
        setError(e.message);
      } finally {
        setLoading(false);
      }
    }

    loadPortfolio();
  }, []);

  return (
    <div className="analysis-screen">
      <header className="analysis-screen__header">
        <h1>Análise do Portfolio</h1>
        <button onClick={() => onNavigate('portfolio')}>Voltar</button>
      </header>

      {error && (
        <div className="analysis-screen__error" role="alert">
          {error}
        </div>
      )}

      {loading && (
        <div className="analysis-screen__loading" aria-live="polite">
          Carregando...
        </div>
      )}

      {!loading && !error && stocks.length === 0 && (
        <p className="analysis-screen__empty">
          Não há dados de portfolio para analisar.
        </p>
      )}

      {!loading && !error && stocks.length > 0 && (
        <table className="analysis-screen__table">
          <thead>
            <tr>
              <th>Ticker</th>
              <th>Nota Fundamentalista</th>
              <th>Peso</th>
            </tr>
          </thead>
          <tbody>
            {stocks.map((stock) => (
              <tr key={stock.ticker}>
                <td>{stock.ticker}</td>
                <td>{stock.fundamentalist_grade}</td>
                <td>{Number(stock.weight).toFixed(2) + '%'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
