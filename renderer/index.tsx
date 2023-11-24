import { createRoot } from 'react-dom/client'
import App from './App'

// TODO: Use SWC Emotion plugin in future once Parcel reads .swcrc files...
const el = document.getElementById('app')
if (el !== null) {
  createRoot(el).render(<App />)
}
