import { createRoot } from 'react-dom/client'
import App from './App'

declare global {
  /* eslint-disable no-var */
  // Exports from Go app process.
  var flash: (filePath: string, devicePath: string) => void
  var cancelFlash: () => void
  var promptForFile: () => void
  var refreshDevices: () => void
  // Export React state to the global scope.
  var setFileReact: (file: string) => void
  var setSpeedReact: (speed: string) => void
  var setDialogReact: (dialog: string) => void
  var setDevicesReact: (devices: string[]) => void
  var setFileSizeReact: (fileSize: number) => void
  var setProgressReact: (progress: number | string | null) => void
  var setSelectedDeviceReact: (selectedDevice: string) => void
} /* eslint-enable no-var */

// LOW-TODO: Use SWC Emotion plugin in future once Parcel reads .swcrc files...
const el = document.getElementById('app')
if (el !== null) {
  createRoot(el).render(<App />)
}
