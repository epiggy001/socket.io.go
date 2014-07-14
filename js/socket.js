(function(window, undefined) {
  "use strict";
  var socket = function(url) {
    this._conn = new WebSocket(url);
    this._handlers = {};

    var self = this;
    this._conn.onmessage = function(evt) {
      var data = JSON.parse(evt.data);
      var fn = self._handlers[data.name];
      if (fn) {
        fn(data.name, data.msg);
      }
    };
  };

  socket.prototype.onclose = function(fn) {
    this._conn.onclose = fn;
  };

  socket.prototype.on = function(evt, fn) {
    this._handlers[evt] = fn;
  };

  socket.prototype.emit = function(name, msg) {
    var m = {};
    m.name = name;
    m.msg = msg;
    this._conn.send(JSON.stringify(m));
  };

  window.Socket = socket;
})(window)
