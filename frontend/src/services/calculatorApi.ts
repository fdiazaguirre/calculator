export type Operation =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'power'
  | 'sqrt'
  | 'percentage'

// Discriminated union: a success carries a numeric result, a failure carries a
// message. Impossible states (success with no result, etc.) cannot be typed.
export type CalculationResponse =
  | { success: true; result: number; error: null }
  | { success: false; result: null; error: string }

function isCalculationResponse(value: unknown): value is CalculationResponse {
  if (typeof value !== 'object' || value === null) {
    return false
  }
  const v = value as Record<string, unknown>
  if (v.success === true) {
    return typeof v.result === 'number' && v.error === null
  }
  if (v.success === false) {
    return v.result === null && typeof v.error === 'string'
  }
  return false
}

export interface CalculateOptions {
  timeoutMs?: number
}

/** Thrown when the backend cannot be reached or does not respond in time (FR-012). */
export class ServiceUnavailableError extends Error {
  constructor() {
    super('Calculator service is unavailable. Please try again.')
    this.name = 'ServiceUnavailableError'
  }
}

const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'
const DEFAULT_TIMEOUT_MS = 5000

/**
 * Sends a calculation request to the backend and returns the response envelope.
 * A `success: false` envelope is a normal (non-throwing) outcome — the caller
 * inspects `error` to decide what to display. A network failure or timeout
 * throws {@link ServiceUnavailableError}.
 */
export async function calculate(
  operation: Operation,
  a: number,
  b?: number,
  options: CalculateOptions = {},
): Promise<CalculationResponse> {
  const { timeoutMs = DEFAULT_TIMEOUT_MS } = options
  const controller = new AbortController()
  const timer = setTimeout(() => controller.abort(), timeoutMs)

  try {
    const response = await fetch(`${API_BASE}/api/v1/calculate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ operation, a, b: b ?? null }),
      signal: controller.signal,
    })

    // The backend returns the same envelope for 200 and 400, so status is not
    // inspected — the shape is. A non-conforming body (proxy error page, gateway
    // JSON) means the service is not answering correctly.
    const body: unknown = await response.json()
    if (!isCalculationResponse(body)) {
      throw new ServiceUnavailableError()
    }
    return body
  } catch (error) {
    if (error instanceof ServiceUnavailableError) {
      throw error
    }
    // AbortError (timeout), TypeError (network failure), and JSON parse errors
    // all mean the user cannot get a result right now.
    throw new ServiceUnavailableError()
  } finally {
    clearTimeout(timer)
  }
}
