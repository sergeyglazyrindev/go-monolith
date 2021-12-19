$(document).ready(function() {
    $('select[name="Type"]').change(function(ev) {
        var select = $(ev.currentTarget);
        var selectVal = select.val();
        if (selectVal === "0") {
            $('select[name="Field"]').parent().parent().parent().parent().hide();
            $('select[name="ContentType"]').parent().parent().parent().parent().hide();
            $('input[name="StaticPath"]').parent().parent().parent().parent().hide();
        } else if (selectVal === "2") {
            $('select[name="Field"]').parent().parent().parent().parent().show();
            $('select[name="ContentType"]').parent().parent().parent().parent().show();
            $('input[name="StaticPath"]').parent().parent().parent().parent().hide();
        } else {
            $('select[name="Field"]').parent().parent().parent().parent().hide();
            $('select[name="ContentType"]').parent().parent().parent().parent().hide();
            $('input[name="StaticPath"]').parent().parent().parent().parent().show();
        }
    });
    $('select[name="ContentType"]').change(function(ev) {
        var select = $(ev.currentTarget);
        var contentTypeIden = select.find('option:selected').data('iden');
        $('select[name="Field"]').find("option").remove();
        if (!contentTypeIden) {
            $('select[name="Field"]').append('<option value="">Choose model</option>');
        } else {
            $('select[name="Field"]').append('<option value="">Choose field</option>');
            var fields = ProjectFields[contentTypeIden];
            fields.forEach(function(field) {
                $('select[name="Field"]').append('<option value="' + field + '">' + field +'</option>');
            });
        }
        var selected = $('select[name="Field"]').data('selected');
        if (selected) {
            $('select[name="Field"]').find('option[value="' + selected +'"]').attr('selected', 'selected');
        }
    });
    $('select[name="ContentType"]').trigger('change');
    $('select[name="Type"]').trigger('change');
});