import React, { useEffect, useState } from 'react'
import { css } from '@emotion/react'
import ReactDOM from 'react-dom'
import Dialog from './dialog'

const floor = num => Math.floor(num * 100) / 100
// const varToString = varObj => Object.keys(varObj)[0]; const s = (setObj, value) => {
// const name = varToString(setObj); setObj[name](value); window[name + 'Go'](value)}

const App = () => {
  const [file, setFile] = useState('')
  const [speed, setSpeed] = useState('')
  const [dialog, setDialog] = useState('')
  const [confirm, setConfirm] = useState(false)
  const [devices, setDevices] = useState(['N/A'])
  const [fileSize, setFileSize] = useState(0)
  const [progress, setProgress] = useState(null)
  const [selectedDevice, setSelectedDevice] = useState('N/A')
  // useEffect(() => window.setFileGo(file), [file])
  useEffect(() => window.refreshDevices(), [])
  window.setFileReact = setFile
  window.setSpeedReact = setSpeed
  window.setDialogReact = setDialog
  window.setDevicesReact = setDevices
  window.setProgressReact = setProgress
  window.setFileSizeReact = setFileSize
  window.setSelectedDeviceReact = setSelectedDevice

  const inProgress = typeof progress === 'number'
  useEffect(() => setConfirm(false), [inProgress])
  const onFlashButtonClick = () => {
    if (inProgress) { // TODO: A dialog would be better.
      if (confirm) {
        setConfirm(false)
        window.cancelFlash()
      } else setConfirm(true)
      return
    }
    setProgress(null)
    if (selectedDevice === 'N/A') return setDialog('Error: Select a device to flash the ISO to!')
    if (!file) return setDialog('Error: Select an ISO to flash to a device!')
    if (BigInt(fileSize) > BigInt(selectedDevice.split(' ')[0])) {
      return setDialog('Error: The ISO file is too big to fit on the selected drive!')
    }
    if (!confirm) return setConfirm(true)
    setConfirm(false)
    window.flash(file, selectedDevice.split(' ')[1])
  }
  const onFileInputChange = (event) => setFile(event.target.value.replace(/\n/g, ''))

  return (
    <>
      {dialog && <Dialog
        handleDismiss={() => setDialog('')}
        message={dialog.startsWith('Error: ') ? dialog.substr(7) : dialog}
        error={dialog.startsWith('Error: ')}
                 />}
      <div css={css`padding: 8;`}>
        <span>Step 1: Enter the path to the file.</span>
        <div css={css`display: flex; padding-bottom: 0.4em;`}>
          <textarea css={css`width: 100%;`} value={file} onChange={onFileInputChange} />
          <button onClick={() => window.promptForFile()}>Select ISO</button>
        </div>
        <span>Step 2: Select the device to flash the ISO to.</span>
        <div css={css`display: flex; padding-bottom: 0.4em; padding-top: 0.4em;`}>
          <select
            css={css`width: 100%`}
            value={selectedDevice}
            onChange={e => setSelectedDevice(e.target.value)}
          >
            {devices.map(device => (
              <option key={device} value={device}>{device.substr(device.indexOf(' ') + 1)}</option>
            ))}
          </select>
          <button onClick={() => window.refreshDevices()} css={css`min-width: 69px;`}>
            Refresh
          </button>
        </div>
        <span>Step 3: Click the button below to begin flashing.</span>
        <div css={css`display: flex; align-items: center; padding-top: 0.4em;`}>
          <button onClick={onFlashButtonClick}>
            {confirm ? 'Confirm?' : (inProgress ? 'Cancel' : 'Flash')}
          </button>
          <div css={css`width: 5;`} />
          {inProgress && <span>Progress: {floor(progress * 100 / fileSize)}% | Speed: {speed}</span>}
          {typeof progress === 'string' && <span>{progress}</span>}
        </div>
      </div>
    </>
  )
}

ReactDOM.render(<App />, document.getElementById('app'))
