// Stock_List — lista de stocks do portfolio
// Requirements: 2.2

import Stock_Item from './Stock_Item';

export default function Stock_List({ stocks, onEdit, onDelete, onViewDetails }) {
  return (
    <div className="stock-list">
      {stocks.map((stock) => (
        <Stock_Item
          key={stock.ticker}
          stock={stock}
          onEdit={onEdit}
          onDelete={onDelete}
          onViewDetails={onViewDetails}
        />
      ))}
    </div>
  );
}
