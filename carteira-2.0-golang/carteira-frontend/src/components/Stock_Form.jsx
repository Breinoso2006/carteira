// Stock_Form — formulário de adição e edição de stocks
// Requirements: 4.1, 4.4, 4.5, 5.2

import { useState, useEffect } from 'react';

/**
 * Props:
 *   mode: 'add' | 'edit'
 *   initialValues?: { ticker: string, fundamentalist_grade: number }
 *   onSubmit: (ticker: string, grade: number) => void
 *   onCancel?: () => void
 */
export default function Stock_Form({ mode, initialValues, onSubmit, onCancel }) {
  const [ticker, setTicker] = useState('');
  const [fundamentalistGrade, setFundamentalistGrade] = useState('');
  const [errors, setErrors] = useState({});

  // Sincroniza os campos quando o modo ou os valores iniciais mudam.
  // Necessário porque o componente permanece montado entre modo 'add' e 'edit'.
  useEffect(() => {
    if (mode === 'edit' && initialValues) {
      setTicker(initialValues.ticker);
      setFundamentalistGrade(String(initialValues.fundamentalist_grade));
    } else {
      setTicker('');
      setFundamentalistGrade('');
    }
    setErrors({});
  }, [mode, initialValues]);

  function validate() {
    const newErrors = {};

    if (!ticker || ticker.trim() === '') {
      newErrors.ticker = 'O ticker é obrigatório.';
    }

    const gradeNum = parseFloat(fundamentalistGrade);
    if (
      fundamentalistGrade.trim() === '' ||
      isNaN(gradeNum) ||
      gradeNum <= 0 ||
      gradeNum > 100
    ) {
      newErrors.grade = 'A nota deve ser um número real no intervalo (0, 100].';
    }

    return newErrors;
  }

  function handleSubmit(e) {
    e.preventDefault();

    const validationErrors = validate();
    if (Object.keys(validationErrors).length > 0) {
      setErrors(validationErrors);
      return;
    }

    setErrors({});
    onSubmit(ticker.trim(), parseFloat(fundamentalistGrade));
  }

  return (
    <form className="stock-form" onSubmit={handleSubmit} noValidate>
      <div className="stock-form__field">
        <label htmlFor="stock-form-ticker">Ticker</label>
        <input
          id="stock-form-ticker"
          type="text"
          value={ticker}
          onChange={(e) => setTicker(e.target.value)}
          readOnly={mode === 'edit'}
          aria-invalid={!!errors.ticker}
          aria-describedby={errors.ticker ? 'stock-form-ticker-error' : undefined}
        />
        {errors.ticker && (
          <span id="stock-form-ticker-error" className="stock-form__error" role="alert">
            {errors.ticker}
          </span>
        )}
      </div>

      <div className="stock-form__field">
        <label htmlFor="stock-form-grade">Nota Fundamentalista</label>
        <input
          id="stock-form-grade"
          type="number"
          step="any"
          value={fundamentalistGrade}
          onChange={(e) => setFundamentalistGrade(e.target.value)}
          aria-invalid={!!errors.grade}
          aria-describedby={errors.grade ? 'stock-form-grade-error' : undefined}
        />
        {errors.grade && (
          <span id="stock-form-grade-error" className="stock-form__error" role="alert">
            {errors.grade}
          </span>
        )}
      </div>

      <div className="stock-form__actions">
        <button type="submit">
          {mode === 'add' ? 'Adicionar' : 'Salvar'}
        </button>
        {onCancel && (
          <button type="button" onClick={onCancel}>
            Cancelar
          </button>
        )}
      </div>
    </form>
  );
}
