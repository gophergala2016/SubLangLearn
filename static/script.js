function connectSocketServer() {
  var support = "MozWebSocket" in window ? 'MozWebSocket' : ("WebSocket" in window ? 'WebSocket' : null);
  ws = new window[support](location.protocol.replace("http","ws") + "//" + location.host + '/socket');

  ws.onmessage = function (evt) {
      var indexes = JSON.parse(evt.data);
      for (var i=0; i < indexes.length; i++) {
          indexes[i] = "#" + indexes[i]
      }
      var selector = indexes.join(", ")
      var newSelected = $(selector)
      selected.css("background-color", "")
      newSelected.css("background-color", "GreenYellow")
      var offset = newSelected.offset();
      offset.top -= $(window).height() / 2;
      selected = newSelected
      $('html, body').animate({
        scrollTop: offset.top,
        scrollLeft: offset.left
      });
  };
}

var selected = $([]);

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