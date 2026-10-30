import { useEffect, useState } from 'react'

const apiBaseUrl = (window.__APP_CONFIG__?.apiBaseUrl || '/api').replace(/\/$/, '')

async function readResponse(response) {
  if (!response.ok) throw new Error(`Richiesta fallita (${response.status})`)
  return response.json()
}

export default function App() {
  const [items, setItems] = useState([])
  const [text, setText] = useState('')
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    const controller = new AbortController()
    fetch(`${apiBaseUrl}/items`, { signal: controller.signal })
      .then(readResponse)
      .then(setItems)
      .catch((cause) => {
        if (cause.name !== 'AbortError') setError('Impossibile caricare la lista.')
      })
      .finally(() => setLoading(false))
    return () => controller.abort()
  }, [])

  async function addItem(event) {
    event.preventDefault()
    const description = text.trim()
    if (!description || saving) return

    setSaving(true)
    setError('')
    try {
      const item = await fetch(`${apiBaseUrl}/items`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ description }),
      }).then(readResponse)
      setItems((current) => [...current, item])
      setText('')
    } catch {
      setError('Impossibile aggiungere l’elemento.')
    } finally {
      setSaving(false)
    }
  }

  return (
    <main className="card">
      <h1>La mia lista</h1>
      <form onSubmit={addItem}>
        <label htmlFor="description">Nuovo elemento</label>
        <div className="input-row">
          <input id="description" value={text} onChange={(event) => setText(event.target.value)} maxLength={255} placeholder="Per esempio: comprare il pane" required />
          <button type="submit" disabled={saving || !text.trim()}>{saving ? 'Salvataggio…' : 'Aggiungi'}</button>
        </div>
      </form>
      {error && <p className="error" role="alert">{error}</p>}
      {loading ? <p>Caricamento…</p> : items.length === 0 ? <p>La lista è vuota.</p> : (
        <ul>{items.map((item) => <li key={item.id}>{item.description}</li>)}</ul>
      )}
    </main>
  )
}
