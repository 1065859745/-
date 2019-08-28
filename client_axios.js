const axios = require('axios')
const express = require('express')
const myLog = require('./exports/log.js')
const bodyParser = require('body-parser')

// server port default 80
const port = 80

// dingtalk url
var url = []
// "https://oapi.dingtalk.com/robot/send?access_token=8138648881fc2be6cac96c4b98bb0ee612b24045759d5bc54b39b30c9c47a8d1"

var argv = process.argv
argv.includes('-u') ? check('-u') : (argv.includes('--url') ? check('--url') : console.error('Please enter your dingtalk url!'))
function check(param) {
  let temp = 0
  url = argv.slice(argv.indexOf(param) + 1)
  url.forEach(element => {
    /^http.?:\/\//g.test(element) ? null : temp++ 
  })
  temp != 0 ? console.log('Please check your dingtalk url!') : (url.length != 0 ? startServer() : console.error('Please enter your dingtalk url!'))
}
function startServer() {
  url.forEach(element => {
    console.log('dingtalk url: %s', element)
  });
  var app = express()
  var staticFile = __dirname + '/static';

  // parser post request
  app.use(bodyParser.urlencoded({ extended: true }))

  app.use(express.static('static'))

  // test DingTalk interface
  app.get('/test', (req, res) => {
    let messages = "Test dingtalk url"
    sendMessages(req, res, messages)
    myLog(req, res)
    res.end()
  })

  // send messages to dingtalk
  function sendMessages(req, res, messages) {
    url.forEach(element => {
      axios.post(
        element, '{"msgtype":"text","text":{"content":"' + messages + '"}}',
        { headers: { 'Content-Type': 'application/json' } }
      ).then(function (response) {

        // Whether the message was sent successfully
        response.data.errmsg == 'ok' ? console.log('messages: %s\nMessage has been successfully sent to DingTalk\nDingTalk url: %s', messages, element) : console.error('Messages sending failure: %s', response.data.errmsg)
        myLog(req, res)
      }).catch(function (error) {
        console.log(error)
      })
    })
    myLog(req, res)
  }

  app.post('/', bodyParser.json(), (req, res) => {
    // send post request to DingTalk
    let messages
    req.body.commonLabels.instance ? messages = req.body.commonLabels.instance + ' ' + req.body.commonLabels.alertname + '\n开始时间: ' + req.body.alerts[0].startsAt + '\n级别: ' + req.body.commonLabels.severity : messages = 'undefined' + req.body.commonLabels.alertname + '\n开始时间: ' + req.body.alerts[0].startsAt + '\n级别: ' + req.body.commonLabels.severity
    req.body.receiver ? sendMessages(req, res, messages) : myLog(req, res)
    res.end()
  })

  // start DingTalk server
  var server = app.listen(port, () => {
    const host = server.address().address
    console.log('server was started at %s:%s', host, port)
  })
}