import { afterAll, afterEach, beforeAll, describe, expect, it } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { Calculator } from './Calculator'

const API_URL = 'http://localhost:8080/api/v1/calculate'

const server = setupServer(
  http.post(API_URL, async ({ request }) => {
    const body = (await request.json()) as { operation: string; a: number; b: number | null }
    const ops: Record<string, number> = {
      add: body.a + (body.b ?? 0),
      subtract: body.a - (body.b ?? 0),
      multiply: body.a * (body.b ?? 0),
      divide: body.a / (body.b ?? 1),
      power: body.a ** (body.b ?? 0),
      sqrt: Math.sqrt(body.a),
      percentage: (body.a * (body.b ?? 0)) / 100,
    }
    // Mirror the backend contract: results arrive rounded to 12 significant
    // digits (FR-008), so 0.1 + 0.2 is delivered as 0.3.
    const result = Number(ops[body.operation].toPrecision(12))
    return HttpResponse.json({ success: true, result, error: null })
  }),
)

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
afterEach(() => server.resetHandlers())
afterAll(() => server.close())

async function compute(a: string, operation: string, b: string) {
  const user = userEvent.setup()
  await user.clear(screen.getByLabelText(/first operand/i))
  await user.type(screen.getByLabelText(/first operand/i), a)
  await user.selectOptions(screen.getByLabelText(/operation/i), operation)
  await user.clear(screen.getByLabelText(/second operand/i))
  await user.type(screen.getByLabelText(/second operand/i), b)
  await user.click(screen.getByRole('button', { name: /calculate/i }))
}

describe('Calculator', () => {
  it('renders the core controls', () => {
    render(<Calculator />)
    expect(screen.getByLabelText(/first operand/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/second operand/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /calculate/i })).toBeInTheDocument()
  })

  it('displays the sum when the user computes 7 + 5', async () => {
    render(<Calculator />)
    await compute('7', 'add', '5')
    expect(await screen.findByText('12')).toBeInTheDocument()
  })

  it('displays 0.3 for 0.1 + 0.2 without floating-point noise', async () => {
    render(<Calculator />)
    await compute('0.1', 'add', '0.2')
    expect(await screen.findByText('0.3')).toBeInTheDocument()
  })

  it('displays the quotient when the user computes 9 ÷ 4', async () => {
    render(<Calculator />)
    await compute('9', 'divide', '4')
    expect(await screen.findByText('2.25')).toBeInTheDocument()
  })

  it('displays the error message for division by zero', async () => {
    server.use(
      http.post(API_URL, () =>
        HttpResponse.json({
          success: false,
          result: null,
          error: 'Division by zero is not allowed',
        }),
      ),
    )
    render(<Calculator />)
    await compute('5', 'divide', '0')
    expect(await screen.findByText('Division by zero is not allowed')).toBeInTheDocument()
  })

  it('rejects a non-numeric first operand locally without calling the backend', async () => {
    let requested = false
    server.use(
      http.post(API_URL, () => {
        requested = true
        return HttpResponse.json({ success: true, result: 0, error: null })
      }),
    )
    render(<Calculator />)
    await compute('abc', 'add', '5')
    expect(await screen.findByText(/first operand must be a number/i)).toBeInTheDocument()
    expect(requested).toBe(false)
  })

  it('rejects an empty second operand locally', async () => {
    render(<Calculator />)
    const user = userEvent.setup()
    await user.clear(screen.getByLabelText(/first operand/i))
    await user.type(screen.getByLabelText(/first operand/i), '5')
    await user.selectOptions(screen.getByLabelText(/operation/i), 'add')
    await user.clear(screen.getByLabelText(/second operand/i))
    await user.click(screen.getByRole('button', { name: /calculate/i }))
    expect(await screen.findByText(/second operand must be a number/i)).toBeInTheDocument()
  })

  it('shows a service-unavailable message when the backend cannot be reached', async () => {
    server.use(http.post(API_URL, () => HttpResponse.error()))
    render(<Calculator />)
    await compute('2', 'add', '2')
    expect(await screen.findByText(/service is unavailable/i)).toBeInTheDocument()
  })

  it('offers the advanced operations in the dropdown', () => {
    render(<Calculator />)
    const select = screen.getByLabelText(/operation/i)
    expect(select).toContainHTML('value="power"')
    expect(select).toContainHTML('value="sqrt"')
    expect(select).toContainHTML('value="percentage"')
  })

  it('computes exponentiation 2 ^ 10 = 1024', async () => {
    render(<Calculator />)
    await compute('2', 'power', '10')
    expect(await screen.findByText('1024')).toBeInTheDocument()
  })

  it('computes 15% of 200 = 30', async () => {
    render(<Calculator />)
    await compute('15', 'percentage', '200')
    expect(await screen.findByText('30')).toBeInTheDocument()
  })

  it('hides the second operand and computes the root when sqrt is selected', async () => {
    const user = userEvent.setup()
    render(<Calculator />)

    await user.selectOptions(screen.getByLabelText(/operation/i), 'sqrt')
    expect(screen.queryByLabelText(/second operand/i)).not.toBeInTheDocument()

    await user.clear(screen.getByLabelText(/first operand/i))
    await user.type(screen.getByLabelText(/first operand/i), '144')
    await user.click(screen.getByRole('button', { name: /calculate/i }))
    expect(await screen.findByText('12')).toBeInTheDocument()
  })

  it('clears a previous error after a subsequent valid calculation', async () => {
    server.use(
      http.post(API_URL, () =>
        HttpResponse.json({
          success: false,
          result: null,
          error: 'Division by zero is not allowed',
        }),
      ),
    )
    render(<Calculator />)
    await compute('5', 'divide', '0')
    expect(await screen.findByText('Division by zero is not allowed')).toBeInTheDocument()

    server.resetHandlers()
    await compute('2', 'add', '2')
    expect(await screen.findByText('4')).toBeInTheDocument()
    expect(screen.queryByText('Division by zero is not allowed')).not.toBeInTheDocument()
  })
})
