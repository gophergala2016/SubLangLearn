function connectSocketServer() {
  var support = "MozWebSocket" in window ? 'MozWebSocket' : ("WebSocket" in window ? 'WebSocket' : null);
  socketClient = new window[support](location.protocol.replace("http","ws") + "//" + location.host + '/socket');

  socketClient.onmessage = function (evt) {
      var indexes = JSON.parse(evt.data);
      var isFutureIndex = false;
      for (var i=0; i < indexes.length; i++) {
          var index = indexes[i]
          if (index < 0) {
              index = -index - 1;
              isFutureIndex = true;
          }
          indexes[i] = "#" + index
      }
      var selector = indexes.join(", ")
      var newSelected = $(selector)
      selected.css("background-color", "")
      newSelected.css("background-color", isFutureIndex ? "LightYellow" : "GreenYellow")
      var offset = newSelected.offset();
      offset.top -= $(window).height() / 2;
      selected = newSelected
      $('html, body').animate({
        scrollTop: offset.top,
        scrollLeft: offset.left
      });
  };

  socketClient.onclose = function(event) {
      socketClient = null;
  }
}

var selected = $([]);

var viewModel;

var socketClient;

function init() {
    connectSocketServer()
    viewModel = {
        subtitles: ko.observableArray(),
        shifts: ko.observableArray(["-2", "-1", "0", "+1", "+2"]),
        selectedShift : ko.observable("0")
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
    if (!socketClient) {
        connectSocketServer()
        getSubtitles()
    }

    $.ajax({
        url: "/play",
        data: {Index: event.target.id, Shift: +viewModel.selectedShift()}
    })
}