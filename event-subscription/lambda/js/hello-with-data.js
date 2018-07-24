module.exports = {
    main: function (event, context) {
        console.log(event);
        return event.data;
    }
}
