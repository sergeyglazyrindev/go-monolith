hljs.initHighlightingOnLoad();
// TODO: Highlight Code for better visual

$(document).ready(function() {
    var parsed = {};
    $('pre code').each(function(i, block) {
        hljs.highlightBlock(block);
        if (block.result.language=="json") {
            var rawJSON = JSON.parse($(this).text() || '{}');
            $(this).text(JSON.stringify(rawJSON, undefined, 2));
        }
        hljs.highlightBlock(block);
    });
});
