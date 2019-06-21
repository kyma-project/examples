module.exports = {
    main: function (event, context) {
        console.log(event);
        console.log("Cloud Events Attributes");
        console.log("ce-specversion: " + event.extensions.request.headers['ce-specversion']);
        console.log("ce-type: " + event.extensions.request.headers['ce-type']);
        console.log("ce-source: " + event.extensions.request.headers['ce-source']);
        console.log("ce-id: " + event.extensions.request.headers['ce-id']);
        console.log("ce-time: " + event.extensions.request.headers['ce-time']);
        console.log("ce-eventtypeversion: " + event.extensions.request.headers['ce-eventtypeversion']);
        console.log("ce-knativehistory: " + event.extensions.request.headers['ce-knativehistory']);
        console.log("content-type: " + event.extensions.request.headers['content-type']);
        return event.data;
    }
}
