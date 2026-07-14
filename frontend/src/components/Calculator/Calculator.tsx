import { useState } from 'react'
import type { FormEvent } from 'react'
import {
  calculate,
  ServiceUnavailableError,
  type CalculationResponse,
  type Operation,
} from '../../services/calculatorApi'
import { Button } from '../ui/Button'
import { Display } from '../ui/Display'
import './calculator.css'

const OPERATIONS: { value: Operation; label: string; symbol: string }[] = [
  { value: 'add', label: 'Add', symbol: '+' },
  { value: 'subtract', label: 'Subtract', symbol: '−' },
  { value: 'multiply', label: 'Multiply', symbol: '×' },
  { value: 'divide', label: 'Divide', symbol: '÷' },
  { value: 'power', label: 'Power', symbol: '^' },
  { value: 'sqrt', label: 'Square root', symbol: '√' },
  { value: 'percentage', label: 'Percentage', symbol: '%' },
]

// Operations that consume a single operand; the second field is hidden for them.
const UNARY_OPERATIONS: ReadonlySet<Operation> = new Set(['sqrt'])

/** Parses an operand field, returning null for empty or non-numeric input. */
function parseOperand(raw: string): number | null {
  const trimmed = raw.trim()
  if (trimmed === '') {
    return null
  }
  const value = Number(trimmed)
  return Number.isFinite(value) ? value : null
}

export function Calculator() {
  const [a, setA] = useState('')
  const [b, setB] = useState('')
  const [operation, setOperation] = useState<Operation>('add')
  const [outcome, setOutcome] = useState<CalculationResponse | null>(null)
  const [pending, setPending] = useState(false)

  const isUnary = UNARY_OPERATIONS.has(operation)

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()

    const parsedA = parseOperand(a)
    if (parsedA === null) {
      setOutcome({ success: false, result: null, error: 'First operand must be a number' })
      return
    }

    let secondOperand: number | undefined
    if (!isUnary) {
      const parsedB = parseOperand(b)
      if (parsedB === null) {
        setOutcome({ success: false, result: null, error: 'Second operand must be a number' })
        return
      }
      secondOperand = parsedB
    }

    setPending(true)
    try {
      const response = await calculate(operation, parsedA, secondOperand)
      setOutcome(response)
    } catch (error) {
      const message =
        error instanceof ServiceUnavailableError
          ? error.message
          : 'Something went wrong. Please try again.'
      setOutcome({ success: false, result: null, error: message })
    } finally {
      setPending(false)
    }
  }

  return (
    <section className="calculator" aria-labelledby="calculator-heading">
      <header className="calculator__header">
        <h1 id="calculator-heading" className="calculator__title">
          Calculator
        </h1>
        <p className="calculator__subtitle">Precise arithmetic, computed server-side.</p>
      </header>

      <Display result={outcome?.result ?? null} error={outcome?.error ?? null} />

      <form className="calculator__form" onSubmit={handleSubmit}>
        <div className="calculator__field">
          <label htmlFor="operand-a">First operand</label>
          <input
            id="operand-a"
            inputMode="decimal"
            value={a}
            onChange={(e) => setA(e.target.value)}
            placeholder="0"
          />
        </div>

        <div className="calculator__field">
          <label htmlFor="operation">Operation</label>
          <select
            id="operation"
            value={operation}
            onChange={(e) => setOperation(e.target.value as Operation)}
          >
            {OPERATIONS.map((op) => (
              <option key={op.value} value={op.value}>
                {op.symbol} {op.label}
              </option>
            ))}
          </select>
        </div>

        {!isUnary && (
          <div className="calculator__field">
            <label htmlFor="operand-b">Second operand</label>
            <input
              id="operand-b"
              inputMode="decimal"
              value={b}
              onChange={(e) => setB(e.target.value)}
              placeholder="0"
            />
          </div>
        )}

        <Button type="submit" disabled={pending}>
          {pending ? 'Calculating…' : 'Calculate'}
        </Button>
      </form>
    </section>
  )
}
