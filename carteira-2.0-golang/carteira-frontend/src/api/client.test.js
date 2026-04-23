import { describe, it, expect, vi, beforeEach } from 'vitest';
import { getPortfolio, addStock, updateStock, removeStock, getStockFundamentals } from './client.js';

beforeEach(() => {
  vi.restoreAllMocks();
});

// Feature: carteira-frontend, Property 14: API_Client usa URLs base corretas
describe('Property 14: API_Client usa URLs base corretas', () => {
  it('getPortfolio usa URL base da carteira-api (http://localhost:3002)', async () => {
    global.fetch = vi.fn().mockResolvedValue({ ok: true, json: async () => [] });
    await getPortfolio();
    expect(fetch).toHaveBeenCalledWith(expect.stringContaining('http://localhost:3002'));
  });

  it('addStock usa URL base da carteira-api (http://localhost:3002)', async () => {
    global.fetch = vi.fn().mockResolvedValue({ ok: true, json: async () => ({}) });
    await addStock('WEGE3', 80);
    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining('http://localhost:3002'),
      expect.any(Object)
    );
  });

  it('updateStock usa URL base da carteira-api (http://localhost:3002)', async () => {
    global.fetch = vi.fn().mockResolvedValue({ ok: true, json: async () => ({}) });
    await updateStock('WEGE3', 90);
    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining('http://localhost:3002'),
      expect.any(Object)
    );
  });

  it('removeStock usa URL base da carteira-api (http://localhost:3002)', async () => {
    global.fetch = vi.fn().mockResolvedValue({ ok: true, json: async () => ({}) });
    await removeStock('WEGE3');
    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining('http://localhost:3002'),
      expect.any(Object)
    );
  });

  it('getStockFundamentals usa URL base da scraping-api (http://localhost:3001)', async () => {
    global.fetch = vi.fn().mockResolvedValue({ ok: true, json: async () => ({}) });
    await getStockFundamentals('WEGE3');
    expect(fetch).toHaveBeenCalledWith(expect.stringContaining('http://localhost:3001'));
  });
});

// Feature: carteira-frontend, Property 15: API_Client propaga erros HTTP com mensagem da API
describe('Property 15: API_Client propaga erros HTTP com mensagem da API', () => {
  it('getPortfolio lança erro com mensagem da API em resposta 4xx', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 404,
      json: async () => ({ error: 'ticker não encontrado' }),
    });
    await expect(getPortfolio()).rejects.toThrow('ticker não encontrado');
  });

  it('getStockFundamentals lança erro com mensagem da API em resposta 4xx', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 404,
      json: async () => ({ error: 'ação não encontrada' }),
    });
    await expect(getStockFundamentals('INVALID')).rejects.toThrow('ação não encontrada');
  });

  it('addStock lança erro com mensagem da API em resposta 4xx', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 400,
      json: async () => ({ error: 'dados inválidos' }),
    });
    await expect(addStock('WEGE3', 80)).rejects.toThrow('dados inválidos');
  });

  it('updateStock lança erro com mensagem da API em resposta 5xx', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 500,
      json: async () => ({ error: 'erro interno do servidor' }),
    });
    await expect(updateStock('WEGE3', 90)).rejects.toThrow('erro interno do servidor');
  });

  it('removeStock lança erro com mensagem da API em resposta 4xx', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 404,
      json: async () => ({ error: 'stock não encontrada' }),
    });
    await expect(removeStock('WEGE3')).rejects.toThrow('stock não encontrada');
  });
});
