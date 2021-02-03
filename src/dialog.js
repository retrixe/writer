import React from 'react'
import { css } from '@emotion/react'
import PropTypes from 'prop-types'

const Dialog = (props) => {
  return (
    <div css={css`
    background-color: rgba(0, 0, 0, 0.4);
    justify-content: center;
    align-items: center;
    position: fixed;
    display: flex;
    height: 100%;
    width: 100%;
    z-index: 1;
    `}>
      <div css={css`
      background-color: white;
      justify-content: flex-start;
      flex-direction: column;
      max-height: 180px;
      max-width: 270px;
      display: flex;
      padding: 8px;
      height: 80%;
      width: 60%;
      `}>
        <h2 css={css`color: ${props.error ? '#ff5555' : 'black'}; margin: 0px;`}>
          {props.error ? 'Error' : 'Message'}
        </h2>
        <p>{props.message}</p>
        <div css={css`flex: 1;`} />
        <button css={css`align-self: center;`} onClick={props.dismiss}>Dismiss</button>
      </div>
    </div>
  )
}

Dialog.propTypes = {
  message: PropTypes.string.isRequired,
  dismiss: PropTypes.func.isRequired,
  error: PropTypes.bool.isRequired
}

export default Dialog
