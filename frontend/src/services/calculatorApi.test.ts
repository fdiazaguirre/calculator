import { afterAll, afterEach, beforeAll, describe, expect, it } from 'vitest'
import { delay, http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { calculate, ServiceUnavailableError } from './calculatorApi'

const API_URL = 'http://localhost:8080/api/v1/calculate'

const server = setupServer(
  http.post(API_URL, async ({ request }) => {
    const body = (await request.json()) as { operation: string; a: number; b: number | null }
    const table: Record<string, number> = {
      add: body.a + (body.b ?? 0),
      subtract: body.a - (body.b ?? 0),
      multiply: body.a * (body.b ?? 0),
      divide: body.a / (body.b ?? 1),
      power: body.a ** (body.b ?? 0),
      sqrt: Math.sqrt(body.a),
      percentage: (body.a * (body.b ?? 0)) / 100,
    }
    return HttpResponse.json({ success: true, result: table[body.operation], error: null })
  }),
)

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
afterEach(() => server.resetHandlers())
afterAll(() => server.close())

describe('calculate', () => {
  it('returns the sum for an add request', async () => {
    const res = await calculate('add', 7, 5)
    expect(res.success).toBe(true)
    expect(res.result).toBe(12)
    expect(res.error).toBeNull()
  })

  it('returns the difference for a subtract request', async () => {
    const res = await calculate('subtract', 10, 4)
    expect(res.result).toBe(6)
  })

  it('returns the product for a multiply request', async () => {
    const res = await calculate('multiply', 6, 7)
    expect(res.result).toBe(42)
  })

  it('returns the quotient for a divide request', async () => {
    const res = await calculate('divide', 9, 4)
    expect(res.result).toBe(2.25)
  })

  it('returns the power for an exponentiation request', async () => {
    const res = await calculate('power', 2, 10)
    expect(res.result).toBe(1024)
  })

  it('returns the root for a unary square-root request (no b)', async () => {
    const res = await calculate('sqrt', 144)
    expect(res.result).toBe(12)
  })

  it('returns the percentage for a percentage request', async () => {
    const res = await calculate('percentage', 15, 200)
    expect(res.result).toBe(30)
  })

  it('sends b as null for a unary square-root request', async () => {
    let received: unknown
    server.use(
      http.post(API_URL, async ({ request }) => {
        received = await request.json()
        return HttpResponse.json({ success: true, result: 3, error: null })
      }),
    )
    await calculate('sqrt', 9)
    expect(received).toEqual({ operation: 'sqrt', a: 9, b: null })
  })

  it('returns a success:false envelope without throwing for a calculation error', async () => {
    server.use(
      http.post(API_URL, () =>
        HttpResponse.json({
          success: false,
          result: null,
          error: 'Division by zero is not allowed',
        }),
      ),
    )
    const res = await calculate('divide', 5, 0)
    expect(res.success).toBe(false)
    expect(res.result).toBeNull()
    expect(res.error).toBe('Division by zero is not allowed')
  })

  it('throws ServiceUnavailableError when the network fails', async () => {
    server.use(http.post(API_URL, () => HttpResponse.error()))
    await expect(calculate('add', 1, 2)).rejects.toBeInstanceOf(ServiceUnavailableError)
  })

  it('throws ServiceUnavailableError when the body is not a valid envelope', async () => {
    // e.g. an nginx/proxy error page returning unexpected JSON.
    server.use(http.post(API_URL, () => HttpResponse.json({ message: 'Bad Gateway' })))
    await expect(calculate('add', 1, 2)).rejects.toBeInstanceOf(ServiceUnavailableError)
  })

  it('throws ServiceUnavailableError when the request times out', async () => {
    server.use(
      http.post(API_URL, async () => {
        await delay(50)
        return HttpResponse.json({ success: true, result: 3, error: null })
      }),
    )
    await expect(calculate('add', 1, 2, { timeoutMs: 10 })).rejects.toBeInstanceOf(
      ServiceUnavailableError,
    )
  })
})
