interface DisplayProps {
  result: number | null
  error: string | null
}

/**
 * Shows exactly one outcome: the result, the error, or a neutral placeholder.
 * Uses aria-live so screen readers announce each new outcome (FR-010).
 */
export function Display({ result, error }: DisplayProps) {
  const state = error !== null ? 'error' : result !== null ? 'result' : 'empty'

  return (
    <div className={`display display--${state}`} role="status" aria-live="polite">
      {result !== null && error === null && (
        <span className="display__value">{formatNumber(result)}</span>
      )}
      {error !== null && <span className="display__error">{error}</span>}
      {result === null && error === null && (
        <span className="display__placeholder">Enter a calculation</span>
      )}
    </div>
  )
}

function formatNumber(value: number): string {
  // Values are already rounded server-side; String() gives the shortest exact
  // decimal (e.g. 0.3, not 0.30) without reintroducing float noise.
  return String(value)
}
