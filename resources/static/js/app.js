var searchOpen = false;


function openSearchForm(key) {
    if (!searchOpen) {
        searchOpen = true;

        $('.search-wrapper').show();

        if (key) {
            $('.search-form input').val(key).focus();
        }
        else {
            $('.search-form input').focus();
        }
    }
}// END openSearchForm

function closeSearchForm() {
    $('.search-wrapper').hide();
    $('.search-form input').val('');

    searchOpen = false;
}// END closeSearchForm


$(document).ready(function(){
    $(document).keyup(function (e) {
        var key = e.charCode || e.keyCode || 0;

        // Detect if keypress was 0-9 or a-z
        if ((key >= 48 && key <= 57) || (key >= 65 && key <= 90)){
            // Open search form
            openSearchForm(String.fromCharCode(key).toLowerCase());
        }

        // Detect ESC key press
        if (key == 27) {
            // Close search form
            closeSearchForm();
        }
    });

    $('.search-nav-btn').click(function(){
        // Open search form
        openSearchForm();
    });

    $('.search-form-close').click(function(){
        // Close search form
        closeSearchForm();
    });
});
