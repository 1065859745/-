function log(req, res) {
  let myDate = new Date()
  let ip = req.ip.split(':')
  console.log('%s -- %s "%s %s %s" %s %s', ip[ip.length - 1], myDate.toLocaleString(), req.method, req.path, req.protocol, res.statusCode, req.get('user-Agent'))
}
module.exports = log