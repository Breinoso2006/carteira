// Stock_Item — item individual da lista de stocks
// Requirements: 3.1, 5.1, 6.1

export default function Stock_Item({ stock, onEdit, onDelete, onViewDetails }) {
  return (
    <div className="stock-item">
      <div className="stock-item__info">
        <span className="stock-item__ticker">{stock.ticker}</span>
        <span className="stock-item__grade">Nota: {stock.fundamentalist_grade}</span>
      </div>
      <div className="stock-item__actions">
        <button onClick={() => onEdit(stock)}>Editar</button>
        <button onClick={() => onDelete(stock.ticker)}>Remover</button>
        <button onClick={() => onViewDetails(stock.ticker)}>Ver Detalhes</button>
      </div>
    </div>
  );
}
