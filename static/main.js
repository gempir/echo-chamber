$(".keyword").keypress(function (e) {
    if (e.which == 13) {
        fetchResults($(".keyword").val());
    }
});

function fetchResults(query) {
    $.get("/search?q=" + encodeURIComponent(query), function (data) {
        var text = data.hits.hits[0]._source.text;
        var username = data.hits.hits[0]._source.username;
        speak(text);
        showMessage(username, text);
    });
}

if ('speechSynthesis' in window) {
    speechSynthesis.onvoiceschanged = function () {
        var $voicelist = $('#voices');

        if ($voicelist.find('option').length == 0) {
            speechSynthesis.getVoices().forEach(function (voice, index) {
                var $option = $('<option>')
                    .val(index)
                    .html(voice.name + (voice.default ? ' (default)' : ''));

                $voicelist.append($option);
            });
        }
    }
}

function speak(text) {
    var msg = new SpeechSynthesisUtterance();
    var voices = window.speechSynthesis.getVoices();
    msg.voice = voices[$('#voices').val()];
    msg.text = text;
    
    speechSynthesis.speak(msg);
}


function showMessage(username, text) {
    $(".message").html(username + ": " + text);
}
