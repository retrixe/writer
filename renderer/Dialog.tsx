import * as styles from './Dialog.module.scss'

const Dialog = (props: {
  handleDismiss: () => void
  message: string
  error: boolean
}): JSX.Element => {
  return (
    <div className={styles.dialog}>
      <div className={styles['dialog-contents']}>
        <h2 className={`${styles.header} ${props.error ? styles.error : ''}`}>
          {props.error ? 'Error' : 'Message'}
        </h2>
        <p>{props.message}</p>
        <div className={styles['flex-spacer']} />
        <button className={styles['dismiss-button']} onClick={props.handleDismiss}>
          Dismiss
        </button>
      </div>
    </div>
  )
}

export default Dialog
