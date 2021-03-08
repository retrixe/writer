import React, { useEffect, useState } from 'react'
import { css } from '@emotion/react'
import ReactDOM from 'react-dom'
import Dialog from './dialog'

const App = () => {
  const [file, setFile] = useState('')
  const [dialog, setDialog] = useState('')
  const [devices, setDevices] = useState(['N/A'])
  const [progress, setProgress] = useState(0)
  const [selectedDevice, setSelectedDevice] = useState('N/A')
  // useEffect(() => window.setFileGo(file), [file])
  useEffect(() => window.refreshDevices(), [])
  window.setFileReact = setFile
  window.setDialogReact = setDialog
  window.setDevicesReact = setDevices
  window.setProgressReact = setProgress
  window.setSelectedDeviceReact = setSelectedDevice

  const onFlashButtonClick = () => {
    if (selectedDevice && selectedDevice !== 'N/A') window.flash(file, selectedDevice.split(' ')[0])
    else setDialog('Error: Select a device to flash the ISO to!')
  }
  const onFileInputChange = (event) => setFile(event.target.value.replace(/\n/g, ''))

  return (
    <>
      {dialog && <Dialog
        dismiss={() => setDialog('')}
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
            {devices.map(device => <option key={device} value={device}>{device}</option>)}
          </select>
          <button onClick={() => window.refreshDevices()} css={css`min-width: 69px;`}>
            Refresh
          </button>
        </div>
        <span>Step 3: Click the button below to begin flashing.</span>
        <div css={css`display: flex; align-items: center; padding-top: 0.4em;`}>
          <button onClick={onFlashButtonClick}>Flash</button>
          <div css={css`width: 5;`} />
          <span>Progress: {progress}</span>
        </div>
      </div>
    </>
  )
}

ReactDOM.render(<App />, document.getElementById('app'))
