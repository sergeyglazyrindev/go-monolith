$('body').delegate('select[name="DataType"]', 'change', function(ev) {
    var widgetType = $(ev.currentTarget).find('option:selected').text().toLowerCase();
    var url = location.pathname + '?widgetType=' + widgetType;
    location.href = url;
})