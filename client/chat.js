var Ractive = require('ractive')
var websocket = require('websocket-stream')
var template = require('./chatbox_template.html')

var ws = websocket('ws://localhost:5678')

var messages = ["Welcome to chat!"]

var ractive = new Ractive({
  el: 'chat-container',
  template: template,
  data: {messages: messages}
})

ws.on('data', function(data) {
    messages.push(data)
    ractive.set('messages', messages)
})

var submit_el = document.getElementById('submit-form')
var text_el = document.getElementById('submit-text')
submit_el.addEventListener('submit', function(ev) {
    ev.preventDefault()
    message = text_el.value
    ws.write(message)
    text_el.value = ''
})
