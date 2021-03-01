import React, { useState, useEffect } from 'react'
import { css } from '@emotion/react'
import ReactDOM from 'react-dom'
import Dialog from './dialog'

const App = () => {
  const [file, setFile] = useState('')
  const [dialog, setDialog] = useState('')
  const [devices, setDevices] = useState(['N/A'])
  const [progress, setProgress] = useState(0)
  const [selectedDevice, setSelectedDevice] = useState('N/A')
  useEffect(() => window.setFileGo(file), [file])
  useEffect(() => window.setSelectedDeviceGo(selectedDevice), [selectedDevice])
  window.setFileReact = setFile
  window.setDialogReact = setDialog
  window.setDevicesReact = setDevices
  window.setProgressReact = setProgress
  window.setSelectedDeviceReact = setSelectedDevice

  return (
    <>
      {dialog && <Dialog
        dismiss={() => setDialog('')}
        message={dialog.startsWith('Error: ') ? dialog.substr(7) : dialog}
        error={dialog.startsWith('Error: ')}
      />}
      <div css={css`padding: 8;`}>
        <span>Step 1: Enter the path to the file.</span>
        <div css={css`display: flex;`}>
          <input css={css`width: 100%;`} value={file} onChange={e => setFile(e.target.value)} />
          <button onClick={() => window.promptForFile()}>Select ISO</button>
        </div>
        <br />
        <span>Step 2: Select the device to flash the ISO to.</span>
        <select
          css={css`width: 100%`}
          value={selectedDevice}
          onChange={e => setSelectedDevice(e.target.value)}
        >
          {devices.map(device => <option key={device} value={device}>{device}</option>)}
        </select>
        <br />
        <span>Step 3: Click the button below to begin flashing.</span>
        <div css={css`display: flex; align-items: center;`}>
          <button onClick={() => window.flash()}>Flash</button>
          <div css={css`width: 5;`} />
          <span>Progress: {progress}</span>
        </div>
      </div>
    </>
  )
}

ReactDOM.render(<App />, document.getElementById('app'))
