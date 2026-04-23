// Stock_Details_Modal — modal de dados fundamentalistas de uma stock
// Requirements: 3.2, 3.3, 3.4, 3.5, 3.6, 3.7, 3.8, 3.9, 3.10, 3.11, 3.12, 3.13

import { useState, useEffect } from 'react';
import { createPortal } from 'react-dom';
import { getStockFundamentals } from '../api/client.js';
import { getIndicatorColor } from '../utils/indicators.js';

/**
 * Props:
 *   ticker: string
 *   onClose: () => void
 */
export default function Stock_Details_Modal({ ticker, onClose }) {
  const [fundamentals, setFundamentals] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Busca os dados fundamentalistas no mount (Requirement 3.2, 3.3)
  useEffect(() => {
    async function fetchFundamentals() {
      setError(null);
      setLoading(true);
      try {
        const data = await getStockFundamentals(ticker);
        setFundamentals(data);
      } catch (e) {
        setError(e.message);
      } finally {
        setLoading(false);
      }
    }

    fetchFundamentals();
  }, [ticker]);

  // Fecha ao clicar no overlay externo (Requirement 3.7)
  function handleOverlayClick(e) {
    if (e.target === e.currentTarget) {
      onClose();
    }
  }

  // Retorna a classe CSS para um indicador (Requirement 3.8–3.13)
  function indicatorClass(indicator) {
    if (!fundamentals) return '';
    const color = getIndicatorColor(indicator, fundamentals);
    return color === 'neutral' ? '' : color;
  }

  // Verifica se um campo está em invalid_fields (Requirement 3.5)
  function isInvalid(fieldKey) {
    if (!fundamentals) return false;
    return fundamentals.invalid_fields.includes(fieldKey);
  }

  // Formata um valor numérico ou exibe '—' se nulo
  function fmt(value) {
    if (value == null) return '—';
    return value;
  }

  return createPortal(
    <div
      className="modal-overlay"
      onClick={handleOverlayClick}
      role="dialog"
      aria-modal="true"
      aria-label={`Dados fundamentalistas de ${ticker}`}
    >
      <div className="modal-container">
        <header className="modal-header">
          <h2 className="modal-title">
            {fundamentals ? fundamentals.symbol : ticker}
          </h2>
          {/* Botão de fechar — Requirement 3.7 */}
          <button
            className="modal-close"
            onClick={onClose}
            aria-label="Fechar modal"
          >
            ✕
          </button>
        </header>

        <div className="modal-body">
          {/* Indicador de carregamento — Requirement 3.3 */}
          {loading && (
            <div className="modal-loading" aria-live="polite">
              Carregando dados de {ticker}...
            </div>
          )}

          {/* Mensagem de erro — Requirement 3.6 */}
          {!loading && error && (
            <div className="modal-error" role="alert">
              {error}
            </div>
          )}

          {/* Dados fundamentalistas — Requirement 3.4, 3.5, 3.8–3.13 */}
          {!loading && !error && fundamentals && (
            <div className="modal-fundamentals">
              {/* Campos informativos */}
              <div className="fundamentals-row">
                <span className="fundamentals-label">Símbolo</span>
                <span className="fundamentals-value">{fundamentals.symbol}</span>
              </div>

              <div className="fundamentals-row">
                <span className="fundamentals-label">
                  Preço
                  {isInvalid('Price') && (
                    <span className="invalid-warning" title="Campo inválido ou indisponível"> ⚠</span>
                  )}
                </span>
                <span className={`fundamentals-value ${isInvalid('Price') ? 'invalid-field' : ''}`}>
                  {fmt(fundamentals.price)}
                </span>
              </div>

              {/* Indicadores com coloração — Requirement 3.8–3.13 */}
              <div className="fundamentals-row">
                <span className="fundamentals-label">
                  P/E
                  {isInvalid('PE') && (
                    <span className="invalid-warning" title="Campo inválido ou indisponível"> ⚠</span>
                  )}
                </span>
                <span className={`fundamentals-value ${indicatorClass('pe')} ${isInvalid('PE') ? 'invalid-field' : ''}`}>
                  {fmt(fundamentals.pe)}
                </span>
              </div>

              <div className="fundamentals-row">
                <span className="fundamentals-label">
                  P/BV
                  {isInvalid('PBV') && (
                    <span className="invalid-warning" title="Campo inválido ou indisponível"> ⚠</span>
                  )}
                </span>
                <span className={`fundamentals-value ${indicatorClass('pbv')} ${isInvalid('PBV') ? 'invalid-field' : ''}`}>
                  {fmt(fundamentals.pbv)}
                </span>
              </div>

              <div className="fundamentals-row">
                <span className="fundamentals-label">
                  PSR
                  {isInvalid('PSR') && (
                    <span className="invalid-warning" title="Campo inválido ou indisponível"> ⚠</span>
                  )}
                </span>
                <span className={`fundamentals-value ${indicatorClass('psr')} ${isInvalid('PSR') ? 'invalid-field' : ''}`}>
                  {fmt(fundamentals.psr)}
                </span>
              </div>

              <div className="fundamentals-row">
                <span className="fundamentals-label">
                  BVps
                  {isInvalid('BVps') && (
                    <span className="invalid-warning" title="Campo inválido ou indisponível"> ⚠</span>
                  )}
                </span>
                <span className={`fundamentals-value ${isInvalid('BVps') ? 'invalid-field' : ''}`}>
                  {fmt(fundamentals.bvps)}
                </span>
              </div>

              <div className="fundamentals-row">
                <span className="fundamentals-label">
                  EPS
                  {isInvalid('EPS') && (
                    <span className="invalid-warning" title="Campo inválido ou indisponível"> ⚠</span>
                  )}
                </span>
                <span className={`fundamentals-value ${isInvalid('EPS') ? 'invalid-field' : ''}`}>
                  {fmt(fundamentals.eps)}
                </span>
              </div>

              <div className="fundamentals-row">
                <span className="fundamentals-label">
                  DY
                  {isInvalid('DY') && (
                    <span className="invalid-warning" title="Campo inválido ou indisponível"> ⚠</span>
                  )}
                </span>
                <span className={`fundamentals-value ${indicatorClass('dy')} ${isInvalid('DY') ? 'invalid-field' : ''}`}>
                  {fmt(fundamentals.dy)}
                </span>
              </div>

              {/* Graham — Requirement 3.12, 3.13 */}
              <div className="fundamentals-row">
                <span className="fundamentals-label">
                  Graham
                  {(['EPS', 'BVps', 'Price'].some(f => isInvalid(f))) && (
                    <span className="invalid-warning" title="Um ou mais campos necessários são inválidos"> ⚠</span>
                  )}
                </span>
                <span className={`fundamentals-value ${indicatorClass('graham')}`}>
                  {fundamentals.eps != null && fundamentals.bvps != null && fundamentals.eps > 0 && fundamentals.bvps > 0
                    ? Math.sqrt(22.5 * fundamentals.eps * fundamentals.bvps).toFixed(2)
                    : '—'}
                </span>
              </div>

              <div className="fundamentals-row">
                <span className="fundamentals-label">Fonte</span>
                <span className="fundamentals-value">{fundamentals.source}</span>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  , document.body);
}
