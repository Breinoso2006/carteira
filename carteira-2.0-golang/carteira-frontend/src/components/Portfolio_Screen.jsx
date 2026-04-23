// Portfolio_Screen — tela principal de gerenciamento do portfolio
// Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 4.2, 4.3, 4.6, 5.3, 5.4, 5.5, 6.2, 6.3, 6.4

import { useState, useEffect } from 'react';
import { getPortfolio, addStock, updateStock, removeStock } from '../api/client';
import Stock_List from './Stock_List';
import Stock_Form from './Stock_Form';
import Stock_Details_Modal from './Stock_Details_Modal';

/**
 * Props:
 *   onNavigate: (screen: string) => void
 */
export default function Portfolio_Screen({ onNavigate }) {
  const [stocks, setStocks] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [formMode, setFormMode] = useState('add');
  const [editingStock, setEditingStock] = useState(null);
  const [modalTicker, setModalTicker] = useState(null);

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

  // Carrega o portfolio no mount
  useEffect(() => {
    loadPortfolio();
  }, []);

  async function handleFormSubmit(ticker, grade) {
    setError(null);
    try {
      if (formMode === 'add') {
        await addStock(ticker, grade);
      } else {
        await updateStock(ticker, grade);
      }
      // Volta para modo add e limpa o stock em edição
      setFormMode('add');
      setEditingStock(null);
      // Recarrega a lista após mutação bem-sucedida
      await loadPortfolio();
    } catch (e) {
      setError(e.message);
    }
  }

  function handleEdit(stock) {
    setFormMode('edit');
    setEditingStock(stock);
  }

  function handleCancelEdit() {
    setFormMode('add');
    setEditingStock(null);
  }

  async function handleDelete(ticker) {
    setError(null);
    try {
      await removeStock(ticker);
      // Recarrega a lista após remoção bem-sucedida
      await loadPortfolio();
    } catch (e) {
      setError(e.message);
    }
  }

  function handleViewDetails(ticker) {
    setModalTicker(ticker);
  }

  function handleCloseModal() {
    setModalTicker(null);
  }

  return (
    <div className="portfolio-screen">
      <header className="portfolio-screen__header">
        <h1>Carteira</h1>
        <button onClick={() => onNavigate('analysis')}>Ver Análise</button>
      </header>

      {error && (
        <div className="portfolio-screen__error" role="alert">
          {error}
        </div>
      )}

      <section className="portfolio-screen__form">
        <Stock_Form
          mode={formMode}
          initialValues={
            formMode === 'edit' && editingStock
              ? {
                  ticker: editingStock.ticker,
                  fundamentalist_grade: editingStock.fundamentalist_grade,
                }
              : undefined
          }
          onSubmit={handleFormSubmit}
          onCancel={formMode === 'edit' ? handleCancelEdit : undefined}
        />
      </section>

      <section className="portfolio-screen__list">
        {loading && (
          <div className="portfolio-screen__loading" aria-live="polite">
            Carregando...
          </div>
        )}

        {!loading && !error && stocks.length === 0 && (
          <p className="portfolio-screen__empty">
            Nenhuma stock cadastrada no portfolio.
          </p>
        )}

        {!loading && stocks.length > 0 && (
          <Stock_List
            stocks={stocks}
            onEdit={handleEdit}
            onDelete={handleDelete}
            onViewDetails={handleViewDetails}
          />
        )}
      </section>

      {modalTicker !== null && (
        <Stock_Details_Modal
          ticker={modalTicker}
          onClose={handleCloseModal}
        />
      )}
    </div>
  );
}
