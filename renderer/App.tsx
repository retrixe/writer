import { useEffect, useState } from 'react'
import { css } from '@emotion/react'
import JSBI from 'jsbi'
import Dialog from './Dialog'

declare global {
  /* eslint-disable no-var */
  // Flash and cancel flash.
  var flash: (filePath: string, devicePath: string) => void
  var cancelFlash: () => void
  // UI update prompts.
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

// const floor = (num: number) => Math.floor(num * 100) / 100
// const varToString = varObj => Object.keys(varObj)[0]; const s = (setObj, value) => {
// const name = varToString(setObj); setObj[name](value); window[name + 'Go'](value)}

const App = (): JSX.Element => {
  const [file, setFile] = useState('')
  const [speed, setSpeed] = useState('')
  const [dialog, setDialog] = useState('')
  const [confirm, setConfirm] = useState(false)
  const [devices, setDevices] = useState(['N/A'])
  const [fileSize, setFileSize] = useState(0)
  const [progress, setProgress] = useState<number | string | null>(null)
  const [selectedDevice, setSelectedDevice] = useState('N/A')
  // useEffect(() => globalThis.setFileGo(file), [file])
  useEffect(() => globalThis.refreshDevices(), [])
  globalThis.setFileReact = setFile
  globalThis.setSpeedReact = setSpeed
  globalThis.setDialogReact = setDialog
  globalThis.setDevicesReact = setDevices
  globalThis.setProgressReact = setProgress
  globalThis.setFileSizeReact = setFileSize
  globalThis.setSelectedDeviceReact = setSelectedDevice

  const inProgress = typeof progress === 'number'
  useEffect(() => setConfirm(false), [inProgress])
  const onFlashButtonClick = (): void => {
    if (inProgress) {
      // TODO: A dialog would be better.
      if (confirm) {
        setConfirm(false)
        globalThis.cancelFlash()
      } else setConfirm(true)
      return
    }
    setProgress(null)
    if (selectedDevice === 'N/A') return setDialog('Error: Select a device to flash the ISO to!')
    if (file === '') return setDialog('Error: Select an ISO to flash to a device!')
    if (JSBI.greaterThan(JSBI.BigInt(fileSize), JSBI.BigInt(selectedDevice.split(' ')[0]))) {
      return setDialog('Error: The ISO file is too big to fit on the selected drive!')
    }
    if (!confirm) return setConfirm(true)
    setConfirm(false)
    globalThis.flash(file, selectedDevice.split(' ')[1])
  }
  const onFileInputChange: React.ChangeEventHandler<HTMLTextAreaElement> = event =>
    setFile(event.target.value.replace(/\n/g, ''))

  const progressPercent = inProgress
    ? JSBI.divide(JSBI.multiply(JSBI.BigInt(progress), JSBI.BigInt(100)), JSBI.BigInt(fileSize))
    : JSBI.BigInt(0)
  return (
    <>
      {dialog !== '' && (
        <Dialog
          handleDismiss={() => setDialog('')}
          message={dialog.startsWith('Error: ') ? dialog.substring(7) : dialog}
          error={dialog.startsWith('Error: ')}
        />
      )}
      <div
        css={css`
          padding: 8;
        `}
      >
        <span>Step 1: Enter the path to the file.</span>
        <div
          css={css`
            display: flex;
            padding-bottom: 0.4em;
          `}
        >
          <textarea
            css={css`
              width: 100%;
            `}
            value={file}
            onChange={onFileInputChange}
          />
          <button onClick={() => globalThis.promptForFile()}>Select ISO</button>
        </div>
        <span>Step 2: Select the device to flash the ISO to.</span>
        <div
          css={css`
            display: flex;
            padding-bottom: 0.4em;
            padding-top: 0.4em;
          `}
        >
          <select
            css={css`
              width: 100%;
            `}
            value={selectedDevice}
            onChange={e => setSelectedDevice(e.target.value)}
          >
            {devices.map(device => (
              <option key={device} value={device}>
                {device.substr(device.indexOf(' ') + 1)}
              </option>
            ))}
          </select>
          <button
            onClick={() => globalThis.refreshDevices()}
            css={css`
              min-width: 69px;
            `}
          >
            Refresh
          </button>
        </div>
        <span>Step 3: Click the button below to begin flashing.</span>
        <div
          css={css`
            display: flex;
            align-items: center;
            padding-top: 0.4em;
          `}
        >
          <button onClick={onFlashButtonClick}>
            {confirm ? 'Confirm?' : inProgress ? 'Cancel' : 'Flash'}
          </button>
          <div
            css={css`
              width: 5;
            `}
          />
          {inProgress && (
            <span>
              Progress: {progressPercent.toString()}% | Speed: {speed}
            </span>
          )}
          {typeof progress === 'string' && <span>{progress}</span>}
        </div>
      </div>
    </>
  )
}

export default App
