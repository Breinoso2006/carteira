// API_Client — todas as chamadas HTTP à carteira-api e scraping-api

const CARTEIRA_API = 'http://localhost:3002';
const SCRAPING_API = 'http://localhost:3001';

async function handleResponse(response) {
  if (!response.ok) {
    const body = await response.json().catch(() => ({ error: response.statusText }));
    throw new Error(body.error || `HTTP ${response.status}`);
  }
  return response.json();
}

export async function getPortfolio() {
  const res = await fetch(`${CARTEIRA_API}/portfolio`);
  return handleResponse(res);
}

export async function addStock(ticker, fundamentalistGrade) {
  const res = await fetch(`${CARTEIRA_API}/portfolio`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ ticker, fundamentalist_grade: fundamentalistGrade }),
  });
  return handleResponse(res);
}

export async function updateStock(ticker, fundamentalistGrade) {
  const res = await fetch(`${CARTEIRA_API}/portfolio`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ ticker, fundamentalist_grade: fundamentalistGrade }),
  });
  return handleResponse(res);
}

export async function removeStock(ticker) {
  const res = await fetch(`${CARTEIRA_API}/portfolio/${encodeURIComponent(ticker)}`, {
    method: 'DELETE',
  });
  return handleResponse(res);
}

export async function getStockFundamentals(ticker) {
  const res = await fetch(`${SCRAPING_API}/${encodeURIComponent(ticker)}`);
  const data = await handleResponse(res);
  // A scraping-api retorna campos em PascalCase e sem invalid_fields.
  // Normaliza para lowercase e garante que invalid_fields existe como array.
  return {
    symbol:         data.Symbol        ?? null,
    price:          data.Price         ?? null,
    pe:             data.PE            ?? null,
    pbv:            data.PBV           ?? null,
    psr:            data.PSR           ?? null,
    bvps:           data.BVps          ?? null,
    eps:            data.EPS           ?? null,
    dy:             data.DY            ?? null,
    source:         data.Source        ?? '',
    invalid_fields: Array.isArray(data.invalid_fields) ? data.invalid_fields : [],
  };
}
