const peersEl = document.getElementById('peers')
const logEl = document.getElementById('log')
const msgInput = document.getElementById('message')
const sendBtn = document.getElementById('send')

const url = (location.protocol === 'https:' ? 'wss' : 'ws') + '://' + location.host + '/ws'
const ws = new WebSocket(url)

function addLog(msg) {
  const el = document.createElement('div')
  el.textContent = msg
  logEl.appendChild(el)
  logEl.scrollTop = logEl.scrollHeight
}

ws.addEventListener('open', () => addLog('connected'))
ws.addEventListener('close', () => addLog('disconnected'))
ws.addEventListener('message', ev => {
  const text = ev.data
  addLog(text)
})

sendBtn.addEventListener('click', () => {
  const v = msgInput.value.trim()
  if (!v) return
  ws.send(v)
  msgInput.value = ''
})

msgInput.addEventListener('keydown', e => { if (e.key === 'Enter') sendBtn.click() })
