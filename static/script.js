function connectSocketServer() {
  var support = "MozWebSocket" in window ? 'MozWebSocket' : ("WebSocket" in window ? 'WebSocket' : null);
  ws = new window[support](location.protocol.replace("http","ws") + "//" + location.host + '/socket');

  ws.onmessage = function (evt) {
    console.log(evt.data);
  };
}

var viewModel;

function init() {
    connectSocketServer()
    viewModel = {
        subtitles: ko.observableArray()
    };
    ko.applyBindings(viewModel);
    getSubtitles()
}

$(init);

function getSubtitles() {
    $.ajax({
        url: "/getSubtitles"
    })
    .done(function(data) {
        viewModel.subtitles(data.Lines)
    });
}

function play(event) {
    $.ajax({
        url: "/play",
        data: {Index: event.target.id}
    })
}